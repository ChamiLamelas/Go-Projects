/*
Package url_shortener serves as a library of utilities for the URL-Shortener
application. This includes the definition of our API, database configuration,
and HTTP server implementation. This is used by the main package to instantiate
and run a server easily. This library could be used in other applications
that do more than just initializing and booting a server.

This file provides the definition of our API. The first component of our API
are the endpoints a client sends/receives HTTP traffic to shorten URLs,
expand aliases, and run alias analytics. The second component are the
(struct) types used to represent the JSON requests and responses sent back
and forth to clients.
*/

package url_shortener

// Endpoint for shorten operation (map URL <-> alias)
const SHORTEN_ENDPOINT = "/urlshortener/shorten"

// Endpoint for expand operation (get URL from alias)
const EXPAND_ENDPOINT = "/urlshortener/expand/"

// Endpoint for analytics operation (get # expansions for alias)
const ANALYTICS_ENDPOINT = "/urlshortener/analytics/"

/*
Specifies the JSON structure for body of an HTTP request to
shorten/ endpoint. A user must provide a URL to shorten and
optionally, an alias for shortening the URL to.

Note about type definitions: the strings in backticks (`)
are known as a struct tag. A struct tag is used to associate
metadata with a struct field. The Go JSON encoding and
decoding packages use this metadata to map from the struct
field to a JSON key. For example, the Url field is mapped
to the url key in a JSON object.

It is worth pointing out the additional string 'omitempty'
that's used below. This serves two purposes. The first
purpose (which is not used in this application) is that
if that a field has the zero value of its type (e.g.
the empty string for strings) then it will not be included
in the JSON representation of this struct. The second
purpose (which is used in this application) is that if
the JSON is missing the particular key, the zero value
is substituted into the struct's field. In our application,
if the alias is not provided, it is set to the empty string
signally an alias must be automatically assigned.
*/
type ShortenRequest struct {
	Url   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

/*
Specifies the JSON structure for body of an HTTP response from
shorten/ endpoint. A user will receive the URL <-> alias mapping
that was created.
*/
type ShortenResponse struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}

/*
Specifies the JSON structure for body of an HTTP response from
expand/ endpoint. A user will receive the URL <-> alias mapping
that corresponds to the expansion.
*/
type ExpandResponse struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}

/*
Specifies the JSON structure for body of an HTTP response from
analytics/ endpoint. A user will receive the URL <-> alias
mapping and the number of times it was expanded.
*/
type AnalyticsResponse struct {
	Url        string `json:"url"`
	Alias      string `json:"alias"`
	Expansions int    `json:"expansions"`
}
