package main

import (
	"log"
	"main/server"

	"github.com/joho/godotenv"

	"main/server/db"
	"main/server/handler"
	admin "main/server/handler/admin"
	player "main/server/handler/player"
	"main/server/services/auth"
	"main/server/utils"

	"main/server/socket"
	"os"
)

// @title Gin Demo App
// @version 1.0
// @description This is a demo version of Gin app.
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connection := db.InitDB()
	db.Transfer(connection)
	socketServer := socket.SocketInit()
	utils.SocketServerInstance = socketServer
	handler.StartCron(socketServer)
	defer socketServer.Close()
	app := server.NewServer(connection)
	server.ConfigureRoutes(app, utils.SocketServerInstance)

	// //by default insertion
	go admin.AdminSignUpHandler()
	go handler.AddDummyDataHandler()
	go auth.AddAiToDB()
	go player.AddPlayerLevel()
	go handler.AddShopDataToDB()

	if err := app.Run(os.Getenv("PORT")); err != nil {
		log.Print(err)
	}
}
