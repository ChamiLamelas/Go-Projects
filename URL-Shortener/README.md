# URL Shortener

## Overview 

This is a simple URL shortener similar to [bitly](https://bitly.com/). 

At a high level, it is implemented via a (local) HTTP server via a simple JSON RESTful API. Server state is maintained between server boots via an on-disk [SQLite](https://www.sqlite.org/) database.

> Note: the on-disk database is stored in a folder `data/` next to `src/` in a file called `database.db`.

It offers the following features: 

1. A user can provide a URL to be shortened to an alias. By leaving the alias blank, an alias is automatically assigned. 
2. A user can expand an alias to a URL. 
3. A user can see how many times a URL has been expanded.

See [DESIGN.md](./DESIGN.md) for my full design.

## Running the Server

### Compilation Method 

1. Go into `src/`.
2. Run `go build .` 
3. Run `./url-shortener.exe`
4. Run `Ctrl + C` to stop the server.

### Direct Execution 

1. Go into `src/`.
2. Run `go run .`
3. Run `Ctrl + C` to stop the server. 

    > Note you will see `exit status 0xc000013a` which is expected. This just means the program was aborted by a manual `Ctrl + C`.

## Using the Server 

The easiest way to use the server is to make requests with curl. On Windows, use Cygwin. I've given some sample interactions below.

1. Shorten a URL to an automatically assigned alias:

    ```bash
    curl -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com"}'
    ```

    The response looks like: 

    ```json
    {
        "url":"https://www.google.com",
        "alias":"0"
    }
    ```

2. Shorten a URL to a custom alias: 

    ```bash
    curl -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com", "alias":"google"}'
    ```

    The response looks like: 

    ```json
    {
        "url":"https://www.google.com",
        "alias":"google"
    }
    ```

3. Expand an alias: 

    ```bash
    curl -X GET http://localhost:8000/urlshortener/expand/google
    ```

    The response looks like: 

    ```json
    {
        "url":"https://www.google.com",
        "alias":"google"
    }
    ```

4. Get analytics on an alias: 

    ```bash
    curl -X GET http://localhost:8000/urlshortener/analytics/google
    ```

    The response looks like: 

    ```json
    {
        "url":"https://www.google.com",
        "alias":"google",
        "expansions":1
    }
    ```

## Platform

This was implemented and tested on Windows 10 using `go version go1.23.0 windows/amd64` and [Cygwin](https://www.cygwin.com/). For more details on testing, see [TESTING.md](./TESTING.md)

## Acknowledgements 

- [Stack Overflow](https://stackoverflow.com/)
- [Go documentation](https://pkg.go.dev/std)
- [tour-of-go](https://go.dev/tour/)
- [ChatGPT](https://chatgpt.com/)

## Authors 

- [Chami Lamelas](https://sites.google.com/brandeis.edu/chamilamelas)