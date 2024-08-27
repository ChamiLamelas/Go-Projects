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
	"sync"

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
	// Connection to SQLite database that holds mapping/analytics table
	db *sql.DB

	/*
		The first alias we try when assigning an alias automatically.
		Note, as shown in ShortenAutomatic( ) that multiple aliases
		may have to be tried.
	*/
	nextAlias int

	/*
		Mutex lock that ensures synchronized (consistent) updates to
		nextAlias in the event that multiple requests come in at the
		same time to automatically assign an alias. For more, see
		ShortenAutomatic( ).
	*/
	nextAliasLock sync.Mutex
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

	// Uncomment for testing concurrency robustness
	// log.Printf("Beginning to service shorten request for %s", request.Url)
	// defer log.Printf("Finished servicing shorten request for %s", request.Url)

	/*
		It is possible that two shorten/ requests come in back to back that
		require automatic alias assignment. As a result of Go's route handling
		this will spawn two concurrently running goroutines.

		For the purpose of illustrating the lock motivation, let us break
		this down to two goroutines/threads that are doing an increment
		operation. Let TMP be the result of the operation s.nextAlias + 1
		that occurs on the right hand side of the expanded +=. Let A be
		our counter (nextAlias). Let A start at 0.

		Here is a possible order of execution subject to context switching.

				goroutine #1					goroutine #2

		1.		TMP = A + 1
		2.										TMP = A + 1
		3.										A = TMP
		4.		A = TMP

		As a result, both routines set A to 1, even though two increments
		are done. Therefore, this is a race condition and the increments
		need to be protected by a Mutex lock.

		The code below is an expansion of the above scenario that involves
		potentially multiple updates to nextAlias. Hence, the whole function
		is locked off as a critical section.
	*/
	s.nextAliasLock.Lock()
	defer s.nextAliasLock.Unlock()

	var alias string

	// Note for { } is proper Go syntax for a while (true) { }
	for {
		// Convert current next alias to string and try to insert
		alias = strconv.Itoa(s.nextAlias)
		_, err := s.db.Exec(QUERY_MAKE_MAPPING_TEMPLATE, request.Url, alias, true)
		if err == nil {
			// Insertion successful -- return after we increase nextAlias
			s.nextAlias += 1
			break
		} else if err.Error() == DUPLICATE_URL_VIOLATION {
			// Insertion failed because the URL already has an alias

			/*
				Get the alias for the URL that we are trying to make
				a mapping for. This is so that a user would know how
				to visit their desired URL via the URL-Shortener
				application.

				Note, this query can fail so we overwrite err after
				saving the original duplicate error for logging
				in Shorten( ).
			*/
			duplicate_url_err := err
			alias, err = GetAliasByURL(s, request.Url)

			/*
				Don't expect this query to fail (because insertion failed
				due to duplicated URLs)
			*/
			if err != nil {
				return "", INTERNAL_ERROR_MESSAGE, err
			}
			return "", fmt.Sprintf("URL already has an alias %s.", alias), duplicate_url_err
		} else if err.Error() == DUPLICATE_ALIAS_VIOLATION {
			// Insertion failed because the alias is in use for another URL

			/*
				Unlike in the custom alias case, if insertion fails due to duplicate
				alias we can't just give up (as we are supposed to be assigning an
				alias automatically).

				First off, the reason we even check for this is because the user could
				make an alias be an automatic alias. For example, suppose now shorten/
				requests have been made and the user shortens with a custom alias of "0".
				This would conflict as this is the first nextAlias value.

				Therefore, we keep incrementing nextAlias until we find one that does
				not have a conflict and use that one as the automatic alias.
			*/
			s.nextAlias += 1
		} else {
			// Insertion failed for unexpected reason
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
	// Insert custom mapping into database
	_, err := s.db.Exec(QUERY_MAKE_MAPPING_TEMPLATE, request.Url, request.Alias, false)

	if err == nil {
		// Insertion successful -- return immediately
		return request.Alias, "", nil
	} else if err.Error() == DUPLICATE_URL_VIOLATION {
		// Insertion failed because the URL already has an alias

		/*
			Get the alias for the URL that we are trying to make
			a mapping for. This is so that a user would know how
			to visit their desired URL via the URL-Shortener
			application.

			Note, this query can fail so we overwrite err after
			saving the original duplicate error for logging
			in Shorten( ).
		*/
		duplicate_url_err := err
		alias, err := GetAliasByURL(s, request.Url)

		/*
			Don't expect this query to fail (because insertion failed
			due to duplicated URLs)
		*/
		if err != nil {
			return "", INTERNAL_ERROR_MESSAGE, err
		}

		return "", fmt.Sprintf("URL already has an alias %s.", alias), duplicate_url_err
	} else if err.Error() == DUPLICATE_ALIAS_VIOLATION {
		// Insertion failed because alias is being used for another URL
		return "", "Alias is already in use", err
	} else {
		// Insertion failed for unexpected reason
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
	// Only POST requests are allowed on the shorten/ endpoint
	if r.Method != http.MethodPost {
		ReportInvalidMethodError(w, r.Method)
		return
	}

	// Decode provided JSON string into appropriate request type
	var request ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ReportBadRequestError(w, err.Error(), "Invalid JSON format")
		return
	}

	/*
		If decoding (as specified in api.go) results in an empty
		alias we must automatically assign an alias.
	*/
	var alias string
	var err_msg string
	if request.Alias == "" {
		alias, err_msg, err = ShortenAutomatic(s, &request)
	} else {
		alias, err_msg, err = ShortenCustom(s, &request)
	}

	/*
		If an error occurred during shortening, we report it. Any
		internal errors always have the INTERNAL_ERROR_MESSAGE
		message so that's how we determine what type of error
		to report back to the user.
	*/
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
	// Only GET requests are allowed on the expand/ endpoint
	if r.Method != http.MethodGet {
		ReportInvalidMethodError(w, r.Method)
		return
	}

	/*
		Strip off the analytics/ endpoint from URL where request
		was made at to get the alias that was provided in the
		request.
	*/
	alias := strings.TrimPrefix(r.URL.Path, EXPAND_ENDPOINT)

	// Get the URL for the provided alias
	row := s.db.QueryRow(QUERY_GET_URL_BY_ALIAS_TEMPLATE, alias)
	var url string
	err := row.Scan(&url)

	/*
		sql.ErrNoRows is the error provided by Scan in the event that QueryRow( )
		returned nothing (i.e. there was no row with provided alias)

		We don't expect any other errors
	*/
	if err == sql.ErrNoRows {
		ReportBadRequestError(w, "No mapping exists for alias", fmt.Sprintf("Cannot expand %s, not mapped", alias))
		return
	} else if err != nil {
		ReportUnexpectedInternalServerError(w, err)
		return
	}

	/*
		Increase the number of expansions done on alias. Note because UPDATE internally
		does an increment, there's no need to provide the current number of expansions.

		In addition, one might ask why is this not done in a locked/transaction state.
		Suppose a user makes an expand/ request quickly followed by an analytics/
		request. As a result, two goroutines start running conccurently.

		Suppose the following order of operations occur. In the left column,
		SQL SELECT represents the QueryRow( ) call done above (which does
		the expansion) and SQL UPDATE represents the Exec( ) call done below
		(which does the expansion count update). In the right column, SQL
		SELECT represents the QueryRow( ) call done in the function below
		which gets the # of expansions.

				expand/ goroutine				analytics/ goroutine

		1.		SQL SELECT
		2.										SQL SELECT
		3.		SQL UPDATE

		First off, SQL operations in threads (which presumably includes goroutines)
		are default synchronized (and thread safe) as noted here:
		https://www.sqlite.org/draft/faq.html#q6. Therefore, each individual
		SQL operation in the columns above will execute properly without
		inconsitencies.

		Second, it is true that the analytics/ goroutine would report that
		the alias has not been expanded yet even though the expand/ SELECT
		happened first. However, we deem this is okay. This is because the
		alias has not truly been expanded in practice. This is because the
		expanded alias has not yet been returned to the user which would
		only happen after the UPDATE in the expand/ goroutine. We therefore
		do not believe that maintaining this particular consistency is
		worth the overhead of maintaining a locked state.
	*/
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
	// Only GET requests are allowed on the analytics/ endpoint
	if r.Method != http.MethodGet {
		ReportInvalidMethodError(w, r.Method)
		return
	}

	/*
		Strip off the analytics/ endpoint from URL where request
		was made at to get the alias that was provided in the
		request.
	*/
	alias := strings.TrimPrefix(r.URL.Path, ANALYTICS_ENDPOINT)

	// Get the URL, # expansions for the provided alias
	row := s.db.QueryRow(QUERY_GET_ANALYTICS_BY_ALIAS_TEMPLATE, alias)
	var url string
	var expansions int
	err := row.Scan(&url, &expansions)

	/*
		sql.ErrNoRows is the error provided by Scan in the event that QueryRow( )
		returned nothing (i.e. there was no row with provided alias)

		We don't expect any other errors
	*/
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

	// Default log granularity is seconds -- lowering to microseconds
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
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

		It is important to note that when a request comes in, it will
		result in a goroutine spawning where request/route handling is
		done.

		Also, this function always returns an error even when Ctrl+C'd
		by the user. Regardless, when this function terminates, we have
		to make sure the database connection is closed for proper
		resource cleanup.
	*/
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", HOSTNAME, PORT), nil)
	log.Println(err)
	s.db.Close()
}
