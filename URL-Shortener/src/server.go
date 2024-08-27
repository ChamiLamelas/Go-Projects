package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const HOSTNAME = "localhost"

const PORT = 8000

const INTERNAL_ERROR_MESSAGE = "Unexpected Internal Server Error"

type Server struct {
	db        *sql.DB
	nextAlias int
}

func InitializeDatabase(s *Server) error {
	err := os.MkdirAll(DATABASE_FOLDER, os.ModePerm)
	if err != nil {
		return err
	}
	s.db, err = sql.Open(SQL_DRIVER, DATABASE_FILE)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(QUERY_CREATE_TABLE)
	return err
}

func SetNextAlias(s *Server) error {
	row := s.db.QueryRow(QUERY_GET_NEXT_ALIAS)

	var maybe_max_alias sql.NullString
	err := row.Scan(&maybe_max_alias)

	if err != nil {
		return err
	}

	if !maybe_max_alias.Valid {
		s.nextAlias = 0
		return nil
	}

	s.nextAlias, err = strconv.Atoi(maybe_max_alias.String)
	if err != nil {
		return err
	}
	s.nextAlias += 1
	return nil
}

func GetAliasByURL(s *Server, url string) (string, error) {
	row := s.db.QueryRow(QUERY_GET_ALIAS_BY_URL_TEMPLATE, url)
	var alias string
	err := row.Scan(&alias)
	return alias, err
}

func ShortenAutomatic(s *Server, request *ShortenRequest) (string, string, error) {
	var alias string
	for {
		alias = strconv.Itoa(s.nextAlias)
		_, err := s.db.Exec(QUERY_MAKE_MAPPING_TEMPLATE, request.Url, alias, 0, true)
		if err == nil {
			s.nextAlias += 1
			break
		} else if err.Error() == DUPLICATE_URL_VIOLATION {
			duplicate_url_err := err
			alias, err = GetAliasByURL(s, request.Url)
			if err != nil {
				return "", INTERNAL_ERROR_MESSAGE, err
			}
			return "", fmt.Sprintf("URL already has an alias %s.", alias), duplicate_url_err
		} else if err.Error() == DUPLICATE_ALIAS_VIOLATION {
			s.nextAlias += 1
		} else {
			return "", INTERNAL_ERROR_MESSAGE, err
		}
	}
	return alias, "", nil
}

func ShortenCustom(s *Server, request *ShortenRequest) (string, string, error) {
	_, err := s.db.Exec(QUERY_MAKE_MAPPING_TEMPLATE, request.Url, request.Alias, 0, false)
	if err == nil {
		return request.Alias, "", nil
	} else if err.Error() == DUPLICATE_URL_VIOLATION {
		duplicate_url_err := err
		alias, err := GetAliasByURL(s, request.Url)
		if err != nil {
			return "", INTERNAL_ERROR_MESSAGE, err
		}
		return "", fmt.Sprintf("URL already has an alias %s.", alias), duplicate_url_err
	} else if err.Error() == DUPLICATE_ALIAS_VIOLATION {
		return "", "Alias is already in use", err
	} else {
		return "", INTERNAL_ERROR_MESSAGE, err
	}
}

func RespondAsJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func Shorten(s *Server, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleInvalidMethodError(w, r.Method)
		return
	}

	var request ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		HandleBadRequestError(w, err.Error(), "Invalid JSON format")
		return
	}

	var alias string
	var err_msg string

	if request.Alias == "" {
		alias, err_msg, err = ShortenAutomatic(s, &request)
	} else {
		alias, err_msg, err = ShortenCustom(s, &request)
	}

	if err != nil {
		if err_msg == INTERNAL_ERROR_MESSAGE {
			HandleUnexpectedInternalServerError(w, err)
		} else {
			HandleBadRequestError(w, err.Error(), err_msg)
		}
		return
	}

	RespondAsJSON(w, ShortenResponse{
		Url:   request.Url,
		Alias: alias,
	})
}

func HandleUnexpectedInternalServerError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, INTERNAL_ERROR_MESSAGE, http.StatusInternalServerError)
}

func HandleInvalidMethodError(w http.ResponseWriter, method string) {
	log.Printf("Received method: %s\n", method)
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func HandleBadRequestError(w http.ResponseWriter, log_err_msg string, user_err_msg string) {
	log.Printf("Internal Error: %s, Error sent to User: %s", log_err_msg, user_err_msg)
	http.Error(w, user_err_msg, http.StatusBadRequest)
}

func Expand(s *Server, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleInvalidMethodError(w, r.Method)
		return
	}

	alias := strings.TrimPrefix(r.URL.Path, EXPAND_ENDPOINT)

	row := s.db.QueryRow(QUERY_GET_URL_BY_ALIAS_TEMPLATE, alias)
	var url string
	err := row.Scan(&url)

	if err == sql.ErrNoRows {
		HandleBadRequestError(w, "No mapping exists for alias", fmt.Sprintf("Cannot expand %s, not mapped", alias))
		return
	} else if err != nil {
		HandleUnexpectedInternalServerError(w, err)
		return
	}

	_, err = s.db.Exec(QUERY_UPDATE_ANALYTICS_BY_ALIAS_TEMPLATE, alias)
	if err != nil {
		HandleUnexpectedInternalServerError(w, err)
		return
	}

	RespondAsJSON(w, ExpandResponse{
		Url:   url,
		Alias: alias,
	})
}

func Analytics(s *Server, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleInvalidMethodError(w, r.Method)
		return
	}

	alias := strings.TrimPrefix(r.URL.Path, ANALYTICS_ENDPOINT)

	row := s.db.QueryRow(QUERY_GET_ANALYTICS_BY_ALIAS_TEMPLATE, alias)

	var url string
	var expansions int
	err := row.Scan(&url, &expansions)

	if err == sql.ErrNoRows {
		HandleBadRequestError(w, "No mapping exists for alias", fmt.Sprintf("Cannot get analytics for %s, not mapped", alias))
		return
	} else if err != nil {
		HandleUnexpectedInternalServerError(w, err)
		return
	}

	RespondAsJSON(w, AnalyticsResponse{
		Url:        url,
		Alias:      alias,
		Expansions: expansions,
	})
}

func SetUpRoutes(s *Server) {
	http.HandleFunc(SHORTEN_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		Shorten(s, w, r)
	})
	http.HandleFunc(EXPAND_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		Expand(s, w, r)
	})
	http.HandleFunc(ANALYTICS_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		Analytics(s, w, r)
	})
}

func NewServer() *Server {
	server := new(Server)
	err := InitializeDatabase(server)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = SetNextAlias(server)
	if err != nil {
		server.db.Close()
		log.Println(err)
		return nil
	}
	SetUpRoutes(server)
	return server
}

func (s *Server) Run() {
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", HOSTNAME, PORT), nil)
	log.Println(err)
	s.db.Close()
}
