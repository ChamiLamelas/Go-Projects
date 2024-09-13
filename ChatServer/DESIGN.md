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
- Communication will be done by sending JSON messages over WebSockets. WebSockets are used in place of HTTP because HTTP is made to do single request-response interactions and is usually stateless. WebSockets is designed for bidirectional communication where client doesn't have to initiate. This is useful in chats where a client can be sent a message without sending anything themself.
- User IDs are assigned sequentially (first user gets ID 0, second gets 1, etc.).
- goroutines and locks are used throughout to implement concurrent processing of operations for users and rooms.

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

Instead of having an RESTful API over HTTP, we don't have clients communicating with different endpoints. Instead, they connect to a single endpoint and a WebSocket is established where messages are sent back and forth between clients and server. 

With this being said, by "messages" we do not mean messages between user but overall system messages that indicate the type of action a client is performing or what the server is communicating to the client. 

#### Connecting to Server 

In this case the client makes a connection to the following endpoint: `/chatserver/connect`. 

Request format: (empty) 

Response format: 

```json 
{
    "action": "connection",
    "user_id": "0"
}
```

#### Joining a Room

To join a room, a message is sent with the following format: 

```json 
{
    "action": "join",
    "room": "0"
}
```

To which the server will respond: 

```json 
{
    "action": "join",
    "room": "0"
}
```

#### Leaving a Room 

To leave a room, a message is sent with the following format: 

```json 
{
    "action": "leave"
}
```

To which the server will respond (in success case): 

```json 
{
    "action": "leave",
    "room": "0"
}
```

In failure case, the response will be: 

```json
{
    "action": "error", 
    "message": "user was not in a room"
}
```

#### Private Message 

To send a private message, a system message is sent with the following format: 

```json 
{
    "action": "private_message",
    "user_id": "1",
    "message": "Hello!"
}
```

To which the server will respond (in success case): 

```json 
{
    "action": "private_message",
    "user_id": "1",
    "message": "Hello!"
}
```

In failure case, the response will be: 

```json
{
    "action": "error", 
    "message": "user 1 is not on the server"
}
```

#### Messaging a Room 

To send a message to the room, a system message is sent with the following format: 

```json 
{
    "action": "room_message",
    "message": "Hello!"
}
```

To which the server will respond (in success case): 

```json 
{
    "action": "room_message",
    "room": "0",
    "message": "Hello!"
}
```

In failure case, the response will be: 

```json
{
    "action": "error", 
    "message": "user is not in a room"
}
```

#### Seeing Rooms 

To see all the rooms, a system message is sent with the following format: 

```json 
{
    "action": "rooms"
}
```

To which the server will respond: 

```json 
{
    "action": "rooms",
    "rooms": [
        "0",
        "1"
    ]
}
```

#### Seeing Room Participants 

To see all the participants in a room, a system message is sent with the following format: 

```json 
{
    "action": "participants",
    "room": "0"
}
```

To which the server will respond (in the success case): 

```json 
{
    "action": "participants",
    "room": "0",
    "participants": [
        "0", 
        "1"
    ]
}
```

In failure case, the response will be: 

```json
{
    "action": "error", 
    "message": "room 0 does not exist"
}
```

### Code 

The combined system message formats allows us to condense them all into one system message with the following fields: 

- `action` which is always present and is one of: 
    - `"connection"`
    - `"join"`
    - `"leave"`
    - `"error"`
    - `"private_message"`
    - `"room_message"`
    - `"rooms"`
    - `"participants"`
- `user_id` which is optional and used to provide a connected user with their ID and for specifying the recipient of a private message.
- `room` which is optional and used to join a room and to specify the joined or left room in a response.
- `message` which is optional and used to provide user and error messages.
- `rooms` which is optional and used to provide a list of rooms.
- `participants` which is optional and used to provide a list of participants in a room.  

- `client.go` : code for the client which will connect to the server, accept user input from the console to perform various outputs, and display interactions on the terminal.
- `server.go` : maintains a set of client connections and implements the chatting behaviors such as managing rooms and forwarding/relaying messages. 
- `api.go` : defines the API endpoints and API request/response types to match above JSON.
