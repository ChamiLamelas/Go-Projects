/*
Package url_shortener serves as a library of utilities for the URL-Shortener
application. This includes the definition of our API, database configuration,
and HTTP server implementation. This is used by the main package to instantiate
and run a server easily. This library could be used in other applications
that do more than just initializing and booting a server.

This file provides our database configuration. The first part of the configuration
are the database SQL settings (driver, database file). The second part are the
queries (and query templates) used to define the table and operations used by
the HTTP server to implement the URL-Shortener application. The third part
defines some SQL error messages that are used in the server implementation.
*/

package url_shortener

/*
The go sql package requires you to specify a driver for the SQL engine
you want to use. I chose sqlite as it allows you to have on-disk database
(for free) after discovering Heroku no longer provides an entirely free
database service.
*/
const SQL_DRIVER = "sqlite3"

// The folder where we put the database file
const DATABASE_FOLDER = "../data/"

// The path to the database file
const DATABASE_FILE = DATABASE_FOLDER + "database.db"

// Table creation query
const QUERY_CREATE_TABLE = `
CREATE TABLE IF NOT EXISTS aliases (
	URL TEXT UNIQUE NOT NULL,
	Alias TEXT PRIMARY KEY,
	Expansions INT,
	Automatic BOOL
);
`

/*
Query for getting the next alias upon server boot. In particular, gets
the maximum automatically assigned alias currently in the database.

Note: to avoid doing a string max (which would put "2" over "11"),
we cast to integer first. Casting is done slightly differently based
on SQL engine, so this query is not necessarily portable.
*/
const QUERY_GET_NEXT_ALIAS = `
SELECT MAX(CAST(Alias AS INTEGER))
FROM aliases 
WHERE Automatic
`

/*
This is a query template for inserting a new row (representing an
alias <-> URL mapping) into our table. The # expansions is not
templated as it always starts at 0 upon insert.

Note: Go's sql package allows for query templates where placeholders
are specified by a ?. Then, when query is used (either in a Query()
or Exec() call) a substitution is done that takes care to avoid
SQL injections. It will not do a naive string substitution.

Also, backticks here define a raw string in Go akin to a pre-formatted
triple quote string in Python which allows one to forgo having to
put in newlines manually while still preserving code readability.
*/
const QUERY_MAKE_MAPPING_TEMPLATE = `
INSERT INTO aliases (URL, Alias, Expansions, Automatic) 
VALUES (?, ?, 0, ?)
`

// Query to get the alias associated with a URL
const QUERY_GET_ALIAS_BY_URL_TEMPLATE = `
SELECT Alias
FROM aliases 
WHERE URL = ?
`

// Query to get the URL associated with an alias
const QUERY_GET_URL_BY_ALIAS_TEMPLATE = `
SELECT URL
FROM aliases
WHERE Alias = ?
`

// Query to increment number of expansions for an alias
const QUERY_UPDATE_ANALYTICS_BY_ALIAS_TEMPLATE = `
UPDATE aliases 
SET Expansions = Expansions + 1
WHERE Alias = ?
`

// Query to get the number of expansions for an alias
const QUERY_GET_ANALYTICS_BY_ALIAS_TEMPLATE = `
SELECT URL, Expansions
FROM aliases
WHERE Alias = ?
`

// Violation reported when an insert fails due to duplicate URLs
const DUPLICATE_URL_VIOLATION = "UNIQUE constraint failed: aliases.URL"

// Violation reported when an insert fails due to duplicate aliases
const DUPLICATE_ALIAS_VIOLATION = "UNIQUE constraint failed: aliases.Alias"
