# Design Doc 

This is a design document that provides the background for implementation of my URL-shortener.

## Task 

### Overview 

Build a URL shortening service similar to Bit.ly. This will involve creating a web server, handling HTTP requests, storing the shortened URLs in a database, and redirecting to the original URL when the shortened link is accessed.

### Features 

- Shorten a given URL.
- Redirect from a shortened URL to the original URL.
- Store URL mappings in a database.
- Implement basic analytics (e.g., how many times a short URL has been accessed).

## Design  

### Components 

- HTTP server written in Go to service requests.
- On disk sqllite database to store mappings.

### HTTP Server Behavior

### Shorten

User can request a URL to be shortened to an alias. There are two modes, automatic and custom aliasing.

Automatic Aliasing: 
- Success: an alias is automatically generated.
- Failure: URL has already been shortened, and existing alias is provided.

Custom Aliasing: 
- Success: the custom alias is registered.
- Failures: 
    1. URL has already been shortened, and as a result, existing alias is provided.
    2. Alias is already in use for another URL. No mapping is created.

### Expand Alias 

User can expand an alias into the correct URL. 

- Success: an alias can be converted into a URL.
- Failure: an alias cannot be converted into a URL (no mapping exists).

### Analytics

User can request usage analytics of an alias.

- Success: analytics (request #) can be returned for an alias.
- Failure: alias does not exist (and no analytics can be returned).

### HTTP Server Endpoints

#### Shorten

Route: `/urlshortener/shorten`

Method: `POST`

Request formats:

- Automatic aliasing 
    ```json
    {
        "url": "https://www.google.com/"
    }
    ```

- Custom aliasing
    ```json
    {
        "url": "https://www.google.com/",
        "alias": "google"
    }
    ```

Response formats: 

- Automatic aliasing (success)
    ```json
    {
        "url": "https://www.google.com/",
        "alias": "123"
    }
    ```

- Automatic aliasing (failure): no JSON response, bad request error (400)

- Custom aliasing (success)
    ```json
    {
        "url": "https://www.google.com/",
        "alias": "custom"
    }
    ```

- Custom aliasing (failure #1): no JSON response, bad request error (400)

- Custom aliasing (failure #2): no JSON response, bad request error (400)

#### Expand Alias

Route: `/urlshortener/expand/123`

Method: `GET`

Request format: empty body

Response formats:

- Success: 
    ```json
    {
        "url": "https://www.google.com/",
        "alias": "123"
    }
    ```

- Failure: no JSON response, bad request error (400)

#### Analytics 

Route: `/urlshortener/analytics/123`

Method: `GET`

Request format: empty body

Response formats:

- Success: 
    ```json
    {
        "url": "https://www.google.com/",
        "alias": "123",
        "expansions": 100
    }
    ```

- Failure: no JSON response, bad request error (400)

### Computing Aliases

A more complex strategy to compute aliases would be to use some sort of hash. Instead, I will just maintain a counter that is incremented with each alias. 

To provide consistency between server restarts, we can get the maximum alias from the database used previously.

### Database

The database will have one table called `aliases`. The table will have the following schema. 

|Column|Type|Attributes|Description|Notes|
|-|-|-|-|-|
|`URL`|`TEXT`|Unique, non-null|Represents a long (real) URL.|None|
|`Alias`|`VARCHAR(12)`|Primary key|Represents an alias.|This is chosen as the primary key because if one were to split off analytics into another table, you would `JOIN` on this key.|
|`Expansions`|`INT`|None|Number of times an alias has been expanded to its URL.|None|
|`Automatic`|`BOOL`|None|Whether or not alias was automatically generated.|This is used to determine the maximum alias for initializing the counter upon server reboot.|

> Note: if deployed to Postgres/MySQL it may be better to use `VARCHAR` in place of `TEXT` for `URL` and `Alias`. However, `VARCHAR` is treated like `TEXT` by sqllite. See [here](https://www.sqlite.org/datatype3.html).

### Code 

`main.go`
- Initializes and starts a `Server`. 

`server.go` (used by `main.go`)
- Defines the `Server` type and its methods.
    - Methods include route handling methods as well as starting/closing the server.

`api.go` (used by `server.go`)
- Defines the API endpoints.
- Defines the following request, response types to match the JSON formats outlined above: 
    - `ShortenRequest`
    - `ShortenResponse`
    - `ExpandResponse`
    - `AnalyticsResponse`

`queries.go` (used by `server.go`)
- Defines database configurations.
- Defines the queries used by `Server` to interact with the database.




