/*
Package main contains only a single file: main.go. This is the entry point
to the URL-Shortener program. It initializes and starts a URL-Shortener
server. This file makes use of the url_shortener package.

This file consists of only a single function: main( ) which is the entry
point for the program and performs the above operations. Clients can
now connect to the server and shorten/expand URLs.

Note: You are meant to put a file that is meant to be built into an
executable (or a client program/entry point/etc.) into the main package.
*/

package main

/*
Note about this import: in Go you don't import headers or files like in C.
You import entire packages. Here, I am importing the url_shortener
package (which functions as a library of functions for our URL-Shortener
application). The url_shortener at the start is the module (created
via a go mod init url_shortener command).

Finally, observe in the directory structure in src/: the files within
the url_shortener package are placed into a url_shortener folder.
*/
import "url_shortener/url_shortener"

func main() {
	/*
		Note about this call: here we are invoking the NewServer() function
		on the package. It is important to note that you cannot have two
		identically named functions within the files of a package.
	*/
	server := url_shortener.NewServer()
	if server != nil {
		server.Run()
	}
}
