# Chat Server

## Overview 

This is a simple chat server where users can join and message in rooms or message each other directly. 

TODO: add high level implementation

It offers the two primary features: 

1. Clients can create or join chat rooms where they can broadcast messages to other clients connected to the room.
2. Clients can private message other clients connected to the chat server in any room.

There are additional features for viewing rooms, participants, and user movement notifications.

See [DESIGN.md](./DESIGN.md) for my full design.

## Running the Server

### Compilation Method 

To build the client executable: 

1. Go into `src/`.
2. Run `go build -o ./client.exe ./client`. 
3. Run `./client`.

To build the server executable: 

1. Go into `src/`.
2. Run `go build -o ./server.exe ./server`.
3. Run `./server`.
4. Run `Ctrl + C` to stop the server.

### Direct Execution 

To run the client: 

1. Go into `src/`.
2. Run `go run ./client`.

To run the server: 

1. Go into `src/`.
2. Run `go run ./server`.
3. Run `Ctrl + C` to stop the server. 

    > Note you will see `exit status 0xc000013a` (Windows) or `^Csignal: interrupt` (Linux) which is expected. This just means the program was aborted by a manual `Ctrl + C`.

## Example Interaction

TODO: add example user inputs to clients

## Platforms

This was implemented on Windows 10 using `go version go1.23.0 windows/amd64` and [Cygwin](https://www.cygwin.com/). 

It was tested on the above Windows platform as well as on Ubuntu 22.04.4 LTS using `go version go1.23.0 linux/amd64`.

For more details on testing, see [TESTING.md](./TESTING.md)

## Acknowledgements 

- [Stack Overflow](https://stackoverflow.com/)
- [Go documentation](https://pkg.go.dev/std)
- [tour-of-go](https://go.dev/tour/)
- [ChatGPT](https://chatgpt.com/)

## Authors 

- [Chami Lamelas](https://sites.google.com/brandeis.edu/chamilamelas)