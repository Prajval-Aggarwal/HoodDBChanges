package handler

import (
	"fmt"
	"main/server/gateway"
	"main/server/services/socket"
	"main/server/services/socket/car"

	socketio "github.com/googollee/go-socket.io"
)

func SocketHandler(server *socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		fmt.Println("Connecting host is:", s.URL().Host)
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		s.LeaveAll()
		fmt.Println("Disconnection Id:", s.ID())
		fmt.Println("closed", reason)
	})

	//Event for joining the player in a room
	server.OnEvent("/", "join", gateway.SocketAuthMiddleware(socket.JoinRoomEvent))

	//Handler for fetching player details
	server.OnEvent("/", "playerDetails", gateway.SocketAuthMiddleware(socket.GetPlayerDetails))

	//Player Car socket handler
	server.OnEvent("/", "carBuy", gateway.SocketAuthMiddleware(car.BuyCarService))

}
