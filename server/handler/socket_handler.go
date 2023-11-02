package handler

import (
	"fmt"
	"main/server/gateway"
	"main/server/services/socket"
	"main/server/services/socket/car"
	"main/server/services/socket/shop"

	socketio "github.com/googollee/go-socket.io"
)

func SocketHandler(server *socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		fmt.Println()
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
	server.OnEvent("/", "carUpgrade", gateway.SocketAuthMiddleware(car.UpgradeCarService))
	server.OnEvent("/", "carBuy", gateway.SocketAuthMiddleware(car.BuyCarService))
	server.OnEvent("/", "carRepair", gateway.SocketAuthMiddleware(car.RepairCarService))

	//Player car customise
	server.OnEvent("/", "colorCustomise", gateway.SocketAuthMiddleware(car.ColorCustomization))
	server.OnEvent("/", "wheelCustomise", gateway.SocketAuthMiddleware(car.WheelCustomize))
	server.OnEvent("/", "interiorCustomise", gateway.SocketAuthMiddleware(car.InteriorCustomize))
	server.OnEvent("/", "licenseCustomise", gateway.SocketAuthMiddleware(car.LicenseCustomize))

	server.OnEvent("/", "checkSum", gateway.SocketAuthMiddleware(socket.AmountCheckSum))

	//store connection
	server.OnEvent("/", "buyStore", gateway.SocketAuthMiddleware(shop.BuyFromShop))

	//reward socket
	server.OnEvent("/", "open", gateway.SocketAuthMiddleware(socket.Open))
	server.OnEvent("/", "close", gateway.SocketAuthMiddleware(socket.Close))

}
