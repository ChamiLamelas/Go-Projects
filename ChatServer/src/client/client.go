package main

import (
	"chat_server/shared"
	"fmt"
)

func main() {
	fmt.Println("Client")
	resp := shared.SystemMessage{UserId: "0"}
	fmt.Println(resp)
}
