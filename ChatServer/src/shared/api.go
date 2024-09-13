package shared

const ENDPOINT = "/chatserver"

var ACTIONS = map[string]bool{
	"connection":      true,
	"join":            true,
	"leave":           true,
	"error":           true,
	"private_message": true,
	"room_message":    true,
	"rooms":           true,
	"participants":    true,
}

type SystemMessage struct {
	Action       string   `json:"action"`
	UserId       string   `json:"user_id,omitempty"`
	Room         string   `json:"room,omitempty"`
	Message      string   `json:"message,omitempty"`
	Rooms        []string `json:"rooms,omitempty"`
	Participants []string `json:"participants,omitempty"`
}
