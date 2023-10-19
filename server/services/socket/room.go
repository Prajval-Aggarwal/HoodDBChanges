package socket

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

func JoinRoomEvent(s socketio.Conn, reqData map[string]interface{}) {
	playerId := s.Context().(string)

	s.Join(playerId)

	fmt.Println("Player joined the room sucessfuuly")
}
