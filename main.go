package main

import (
	"log"
	"main/server"
	"main/server/db"
	"main/server/socket"
	"os"

	"github.com/joho/godotenv"
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
	defer socketServer.Close()
	app := server.NewServer(connection)
	server.ConfigureRoutes(app)

	//by default insertion
	go admin.AdminSignUpHandler()
	go admin.AddGarageTypesHandler()
	go admin.AddArenaTypesHandler()
	go handler.AddDummyDataHandler()

	go admin.AddRarityHandler()
	go auth.AddAiToDB()

	go player.AddPlayerLevel()

	go handler.AddShopDataToDB()

	if err := app.Run(os.Getenv("PORT")); err != nil {
		log.Print(err)
	}
}
