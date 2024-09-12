# Design Doc 

This is a design document that provides the background for implementation of my Chat Server.

## Task 

### Overview
 
Build a simple chat server using Go's networking libraries. Clients can connect to the server, join chat rooms, and send messages to each other. 

### Features 

- Clients must connect to the server 

    > Note: this will not be user-initiated, but is done automatically 

- Clients can join a specific chat room. If the room does not exist, then it must be created with that single user. If the client is already in another room, they are automatically removed from that room.
- Clients can leave a chat room.
- Clients can private message another client. 
- Clients can send a message that is broadcast to the chat room.
- Clients can see all the available rooms.
- Clients can see all the users in a room.
- Clients can gracefully disconnect from the server.
- Clients are notified of people leaving/joining a room.

### Limitations 

#### Stability / Fault-Tolerance 

- The server is not persistent. It will not remember connected users, created chat rooms, etc. between boots. This could be potentially solved with an on-file SQL database with tables that serve as user and room registries.
- The client-server model puts a lot of stress on the server as all communications are funnelled through it. This is done to avoid having clients maintain/manage connections to other clients and have the server maintain all the overhead with a simple client implementation.

#### Security 

- There is no client authorization. When clients submit their IDs to perform tasks (e.g. sending messages/joining rooms) these are not authenticated. Hence, someone could easily pretend to be someone else.
- Since all communication is done over HTTP room requests, messages, etc. are all easily visible and modifiable (e.g. a user could manipulate what one user is sending to another).

### Components 

- Client-server model: clients (running identical program) send all information to a central server where they are then forwarded to other clients 
- Communication will be done over TCP using HTTP and JSON. 
- User IDs are assigned sequentially (first user gets ID 0, second gets 1, etc.).

## Design 

### Server Behavior 

#### Connecting to Server 

- Success: a client gets an assigned ID and then is put into some overall list

#### Joining a Room

- Success: user must be associated with the room so that they will then receive broadcasts to that room
    - If there are any other users in the room then they need to be notified 
- Failure: a room must be created then success behavior is performed

#### Leaving a Room 

- Success: user is no longer associated with the room so that they will no longer receive broadcasts to that room 
    - If there are any other users in the room then they need to be notified 
- Failure: if a user is not in a room, nothing should happen but they should be notified 

#### Private Message 

- Success: receiver sees the message on their end 
- Failure: if receiver does not exist/aren't connected then sender should be notified 

#### Messaging a Room 

- Success: all other users in room see message on their end 
- Failure: if sender is not in a room then they should be modified 

#### Seeing Rooms 

- Success: all rooms are sent to the client 

#### Seeing Room Participants 

- Success: all room's users are sent to the client (including requesting user)
- Failure: if room does not exist, the client should be notified 

### Communication API 

For specific failure cases see the failures listed in the [previous section](#server-behavior).

#### Connecting to Server 

Route: `/chatserver/connect`

Method: `POST`

Request format: (empty) 

Response format: 

```json 
{
    "id": "0"
}
```

#### Joining a Room

Route: `/chatserver/join/<room>`

Method: `POST`

Request format: 

```json 
{
    "id": "0"
}
```

Response format (success and failure): 

```json 
{
    "room": "<room>",
    "users": [ 
        "0", 
        "1" 
    ]
}
```

#### Leaving a Room 

Route: `/chatserver/leave/<room>`

Method: `POST`

Request format: 

```json 
{
    "id": "0"
}
```

Response format (success): 

```json 
{
    "room": "<room>",
    "users": [ 
        "1" 
    ]
}
```

In the case of failure, a bad request error is sent to the client.

#### Private Message 

Route: `/chatserver/message/direct`

Method: `POST`

Request format: 

```json 
{
    "sender_id": "0",
    "receiver_id": "1",
    "message": "Hello, World!" 
}
```

Response format (success): 

```json 
{
    "sender_id": "0",
    "receiver_id": "1",
    "message": "Hello, World!"
}
```

In the case of failure, a bad request error is sent to the client.

#### Messaging a Room 

Route: `/chatserver/message/room`

Method: `POST`

Request format: 

```json 
{
    "id": "0",
    "message": "Hello, World!" 
}
```

Response format (success): 

```json 
{
    "id": "0", 
    "message": "Hello, World!",
    "room": "<room>",
    "users": [
        "0",
        "1"
    ]
}
```

In the case of failure, a bad request error is sent to the client.

#### Seeing Rooms 

Route: `/chatserver/rooms`

Method: `GET`

Request format: (empty)

Response format (success): 

```json 
{
    "rooms": [ 
        "room1",
        "room2"
    ]
}
```

#### Seeing Room Participants 

Route: `/chatserver/users/<room>`

Method: `GET`

Request format: (empty)

Response format (success): 

```json 
{
    "room": "<room>",
    "users": [ 
        "0",
        "1"
    ]
}
```

In the case of failure, a bad request error is sent to the client.

### Code 

- `client.go` : code for the client which will connect to the server, accept user input from the console to perform various outputs, and display interactions on the terminal.
- `server.go` : maintains a set of client connections and implements the chatting behaviors such as managing rooms and forwarding/relaying messages. 
- `api.go` : defines the API endpoints and API request/response types to match above JSON.
