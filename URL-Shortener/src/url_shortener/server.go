/*
Package url_shortener serves as a library of utilities for the URL-Shortener
application. This includes the definition of our API, database configuration,
and HTTP server implementation. This is used by the main package to instantiate
and run a server easily. This library could be used in other applications
that do more than just initializing and booting a server.

This file provides our HTTP server implementation. There are two functions
an invoking file is meant to use: to setup a server and start it. The
rest of the functions provide under the hood implementation details such
as setting up the database (using the configuration) and implementing the
route handling (i.e. the application operations).
*/

package url_shortener

/*
The first imports are all regular package imports, but the last starts with
an underscore. This is used to say we are importing the package but not
actually using it in the code. Without the underscore, compilation fails.
However, the package still needs to be imported in order to include the
SQLite driver that's used by Go's sql package to instantiate a connection.
*/
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

/*
Specifies the hostname/network interface where we will be listening
for HTTP connections. Setting this to localhost only allows local
connections. Leaving this as the empty string means all interfaces
(including external) which triggers a Firewall warning (due to
lack of rule) on each unique execution on Windows.
*/
const HOSTNAME = "localhost"

// Port to listen on
const PORT = 8000

// Error message for body of internal server errors sent to user
const INTERNAL_ERROR_MESSAGE = "Unexpected Internal Server Error"

// Represents our server type
type Server struct {
	db        *sql.DB
	nextAlias int
}

////////////////////////// PRIVATE FUNCTIONS ///////////////////////

/*
Initializes the SQLite database used by the server. In particular,
the database is loaded from a file (if it exists). If not, the
single table used by our server (aliases) is initialized.

Note that this function takes a *Server as an argument, not as
the receiver. This is because the Server has not been set up
yet and it would be strange to have an uninitialized object
as the receiver. This is the case for the majority of functions
in this file as they are operating on unintialized or partially
initialized Server objects.

Parameters:

	s: Pointer to Server for which we initialize the database

Returns:

	If initialization failed, an error is returned, otherwise if
	all goes well, nil is returned.
*/
func InitializeDatabase(s *Server) error {
	/*
		Makes folder for database file if it doesn't exist,
		basically a mkdir -p followed by a chmod 0x777
	*/
	err := os.MkdirAll(DATABASE_FOLDER, os.ModePerm)
	if err != nil {
		return err
	}

	// Opens a connection to the database (must be closed)
	s.db, err = sql.Open(SQL_DRIVER, DATABASE_FILE)
	if err != nil {
		return err
	}

	// Creates the table if it doesn't exist
	_, err = s.db.Exec(QUERY_CREATE_TABLE)
	return err
}

/*
This sets the next alias value maintained by our server. To
allow our server to work between boots, we cannot restart the
automatic alias assignment from 0 each time. It starts with
1 more than the maximum alias assigned already by the server.

Parameters:

	s: Pointer to Server for which we set the next alias

Returns:

	If initialization failed, an error is returned, otherwise if
	all goes well, nil is returned.
*/
func SetNextAlias(s *Server) error {
	/*
		Parse maximum previously assigned alias. Note, MAX
		will always return an element. A MAX on an empty row
		selection will return NULL. Hence, we use the special
		NullString type which is a type that can represent
		null or a string.
	*/
	row := s.db.QueryRow(QUERY_GET_NEXT_ALIAS)
	var maybe_max_alias sql.NullString
	err := row.Scan(&maybe_max_alias)
	if err != nil {
		return err
	}

	/*
		This indicates MAX returned NULL, so we start with an
		alias of 0.
	*/
	if !maybe_max_alias.Valid {
		s.nextAlias = 0
		return nil
	}

	/*
		Otherwise we convert alias to int and set the next alias
		to 1 beyond it.
	*/
	s.nextAlias, err = strconv.Atoi(maybe_max_alias.String)
	if err != nil {
		return err
	}
	s.nextAlias += 1
	return nil
}

/*
Reports an unexpected internal error back to the user and logs it.

For this and the other two report error functions, we generally log
more detailed information related to some internal failure (e.g.
a SQL violation) while providing more vague or user friendly
messages to the user.

Parameters:

	w: Where we write response for user
	err: Unexpected error that occurred which is logged
*/
func ReportUnexpectedInternalServerError(w http.ResponseWriter, err error) {
	log.Println(err)

	/*
		http.Error will automatically set the provided error code
		whereas writing regularly to w would have the OK code
	*/
	http.Error(w, INTERNAL_ERROR_MESSAGE, http.StatusInternalServerError)
}

/*
Reports an invalid request method error back to the user and logs it.
By request method we mean POST, GET, DELETE, etc.

Parameters:

	w: Where we write response for user
	method: The method that was determined to be incorrect
*/
func ReportInvalidMethodError(w http.ResponseWriter, method string) {
	log.Printf("Received method: %s\n", method)
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

/*
Reports a bad request error back to the user and logs it. A bad request
can contain any issue (e.g. duplicate aliases, URLs, etc.).

Parameters:

	w: Where we write response for user
	log_err_msg: Message we only log related to bad request
	user_err_msg: Message we both log and send to user for bad request
*/
func ReportBadRequestError(w http.ResponseWriter, log_err_msg string, user_err_msg string) {
	log.Printf("Internal Error: %s, Error sent to User: %s", log_err_msg, user_err_msg)
	http.Error(w, user_err_msg, http.StatusBadRequest)
}

/*
Sends a JSON response to a user. This includes specifying the response type
to JSON and then converting a Go type to JSON.

Parameters:

	w: Where we write response for user
	v: Go type to convert to JSON, in our case, this will be one of the
		Response structs specified in api.go. It is any type to match
		the parameter type of Encode( ).
*/
func RespondAsJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

/*
Queries the server's database for the alias that is associated with
a given URL.

Parameters:

	s: Pointer to Server for which we set the next alias
	url: Given URL for which we're finding the associated alias

Returns:

	The alias associated with the URL and, if an error occurred
	during lookup, the lookup error. If no error occurred, nil
	is returned and it is safe to assume alias holds the
	associated URL.
*/
func GetAliasByURL(s *Server, url string) (string, error) {
	row := s.db.QueryRow(QUERY_GET_ALIAS_BY_URL_TEMPLATE, url)
	var alias string
	err := row.Scan(&alias)
	return alias, err
}

/*
Shortens the URL provided by assigning it an automatic alias.

Parameters:

	s: Pointer to HTTP server that will be updated/used to make
		new URL <-> alias mapping
	request: Pointer to struct that represents contents of shorten
		request. For this function, only request.Url is used.

Returns:

	The created alias, error message that is meant to be sent to
	the user (also used to indicate whether an internal or
	request error occurred in Shorten( )), and the internal
	error that occurred. If successful, the error message and
	error are the empty string and nil respectively. If the
	error is nil, it is assumed the returned alias is not empty.
*/
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

/*
Shortens the URL provided by assigning it a provided custom alias.

Parameters:

	s: Pointer to HTTP server that will be updated/used to make
		new URL <-> alias mapping
	request: Pointer to struct that represents contents of shorten
		request.

Returns:

	The provided alias, error message that is meant to be sent to
	the user (also used to indicate whether an internal or
	request error occurred in Shorten( )), and the internal
	error that occurred. If successful, the error message and
	error are the empty string and nil respectively. If the
	error is nil, it is assumed the returned alias is not empty.
*/
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

/*
Handles requests on the /shorten endpoint.

Parameters:

	s: Pointer to HTTP server that will be updated/used to make
		new URL <-> alias mapping
	request: Pointer to struct that represents contents of HTTP
		request
	w: Where we write response for user
*/
func Shorten(s *Server, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ReportInvalidMethodError(w, r.Method)
		return
	}

	var request ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ReportBadRequestError(w, err.Error(), "Invalid JSON format")
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
			ReportUnexpectedInternalServerError(w, err)
		} else {
			ReportBadRequestError(w, err.Error(), err_msg)
		}
		return
	}

	RespondAsJSON(w, ShortenResponse{
		Url:   request.Url,
		Alias: alias,
	})
}

/*
Handles requests on the /expand/ endpoint.

Parameters:

	s: Pointer to HTTP server that will be used to expand an
		alias and record the expansion
	request: Pointer to struct that represents contents of HTTP
		request
	w: Where we write response for user
*/
func Expand(s *Server, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ReportInvalidMethodError(w, r.Method)
		return
	}

	alias := strings.TrimPrefix(r.URL.Path, EXPAND_ENDPOINT)

	row := s.db.QueryRow(QUERY_GET_URL_BY_ALIAS_TEMPLATE, alias)
	var url string
	err := row.Scan(&url)

	if err == sql.ErrNoRows {
		ReportBadRequestError(w, "No mapping exists for alias", fmt.Sprintf("Cannot expand %s, not mapped", alias))
		return
	} else if err != nil {
		ReportUnexpectedInternalServerError(w, err)
		return
	}

	_, err = s.db.Exec(QUERY_UPDATE_ANALYTICS_BY_ALIAS_TEMPLATE, alias)
	if err != nil {
		ReportUnexpectedInternalServerError(w, err)
		return
	}

	RespondAsJSON(w, ExpandResponse{
		Url:   url,
		Alias: alias,
	})
}

/*
Handles requests on the /analytics/ endpoint.

Parameters:

	s: Pointer to HTTP server that will be used to provide
		analytics on particular alias
	request: Pointer to struct that represents contents of HTTP
		request
	w: Where we write response for user
*/
func Analytics(s *Server, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ReportInvalidMethodError(w, r.Method)
		return
	}

	alias := strings.TrimPrefix(r.URL.Path, ANALYTICS_ENDPOINT)

	row := s.db.QueryRow(QUERY_GET_ANALYTICS_BY_ALIAS_TEMPLATE, alias)

	var url string
	var expansions int
	err := row.Scan(&url, &expansions)

	if err == sql.ErrNoRows {
		ReportBadRequestError(w, "No mapping exists for alias", fmt.Sprintf("Cannot get analytics for %s, not mapped", alias))
		return
	} else if err != nil {
		ReportUnexpectedInternalServerError(w, err)
		return
	}

	RespondAsJSON(w, AnalyticsResponse{
		Url:        url,
		Alias:      alias,
		Expansions: expansions,
	})
}

// Sets up the route handling for the server
func SetUpRoutes(s *Server) {
	/*
		HandleFunc sets up the functions that will operate on each
		endpoint. It takes a function that takes only two parameters
		(the response writer and request). However, we need our
		Server object to properly respond to requests, so we create
		an anonymous function that calls versions of each route
		handling function that takes the Server pointer (Shorten,
		Expand, Analytics).
	*/
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

//////////////// PUBLIC FUNCTIONS AND METHODS ///////////////////////

/*
Sets up a new Server object and returns it to the invoking code. This
initializes the server database, next alias, and the route handling.
It returns a pointer to the Server object if setup was successful.
nil is returned if setup failed.
*/
func NewServer() *Server {
	/*
		Go apparently doesn't distinguish between stack and heap in
		its spec, but new( ) does force a heap allocation under the
		hood. In this case, this makes sense as a pointer to our
		Server object should work outside of this function.

		However, if we didn't use new( ) and just did a Server
		literal declaration (i.e. Server{ }) then Go would automatically
		do a heap allocation under the hood once it detects that
		a pointer to a stack allocated struct would be invalidated
		upon function exit.

		See here for more: https://stackoverflow.com/a/10866871
	*/
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

/*
Runs the server by having it start listening on a particular
interface and port. Once it has been closed, the database
connection is closed.

Note, because this function operates on an initialized
Server, it is made a method with a Server receiver.
*/
func (s *Server) Run() {
	/*
		The nil parameter specifies we are using the default request
		multiplexer. In particular, it will try to match the endpoint
		to the routes that have been registered in SetupRoutes( ).

		Also, this function always returns an error even when Ctrl+C'd
		by the user. Regardless, when this function terminates, we have
		to make sure the database connection is closed for proper
		resource cleanup.
	*/
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", HOSTNAME, PORT), nil)
	log.Println(err)
	s.db.Close()
}
