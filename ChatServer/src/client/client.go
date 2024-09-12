package main

import (
	"chat_server/shared"
	"fmt"
)

func main() {
	fmt.Println("Client")
	resp := shared.ConnectResponse{Id: "0"}
	fmt.Println(resp)
}
