package main

import (
	"chat_server/shared"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

const HOSTNAME = "localhost"

const PORT = 8000

type UserData struct {
	Id         string
	Room       string
	Connection *websocket.Conn
	Channel    chan string
}

type RoomData struct {
	Id           string
	Channel      chan string
	Lock         sync.Mutex
	Participants []string
}

type Server struct {
	NextId             int
	Rooms              map[string]RoomData
	RoomsLock          sync.Mutex
	Users              map[string]UserData
	UsersLock          sync.Mutex
	ConnectionUpgrader websocket.Upgrader
}

func ReadJSON(connection *websocket.Conn) (shared.SystemMessage, error) {
	var Msg shared.SystemMessage
	err := connection.ReadJSON(&Msg)
	return Msg, err
}

func AcceptConnection(server *Server, w http.ResponseWriter, r *http.Request) (string, error) {
	server.UsersLock.Lock()
	defer server.UsersLock.Unlock()

	Id := strconv.Itoa(server.NextId)
	server.NextId += 1

	Data := UserData{
		Id:      Id,
		Channel: make(chan string),
	}

	var err error
	Data.Connection, err = server.ConnectionUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return Id, err
	}

	server.Users[Id] = Data
	return Id, nil
}

func ManageConnection(server *Server, w http.ResponseWriter, r *http.Request) error {
	Id, err := AcceptConnection(server, w, r)
	if err != nil {
		return err
	}
	defer server.Users[Id].Connection.Close()

	return nil
}

func NewServer() *Server {
	server := new(Server)
	server.Rooms = make(map[string]RoomData)
	server.Users = make(map[string]UserData)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	http.HandleFunc(shared.ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		ManageConnection(server, w, r)
	})
	return server
}

func (server *Server) Run() {
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", HOSTNAME, PORT), nil)
	log.Println(err)
}

func main() {
	server := NewServer()
	server.Run()
}
