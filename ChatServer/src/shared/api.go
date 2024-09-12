package shared

const CONNECT_ENDPOINT = "/chatserver/connect"

const JOIN_ENDPOINT = "/chatserver/join/"

const LEAVE_ENDPOINT = "/chatserver/leave/"

const PRIVATE_MESSAGE_ENDPOINT = "/chatserver/message/direct"

const ROOM_MESSAGE_ENDPOINT = "/chatserver/message/room"

const ROOMS_ENDPOINT = "/chatserver/rooms"

const PARTICIPANTS_ENDPOINT = "/chatserver/users/"

type ConnectResponse struct {
	Id string `json:"id"`
}

type JoinRequest struct {
	Id string `json:"id"`
}

type JoinResponse struct {
	Room  string   `json:"room"`
	Users []string `json:"users"`
}

type LeaveRequest struct {
	Id string `json:"id"`
}

type LeaveResponse struct {
	Room  string   `json:"room"`
	Users []string `json:"users"`
}

type PrivateMessageRequest struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Message    string `json:"message"`
}

type PrivateMessageResponse struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Message    string `json:"message"`
}

type RoomMessageRequest struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type RoomMessageResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Room    string `json:"room"`
	Users   string `json:"users"`
}

type RoomResponse struct {
	Rooms []string `json:"rooms"`
}

type ParicipantsResponse struct {
	Room  string   `json:"room"`
	Users []string `json:"users"`
}
