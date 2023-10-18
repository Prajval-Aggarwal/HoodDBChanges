package server

import (
	_ "main/docs"
	"main/server/gateway"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {

	//Allowing CORS
	server.engine.Use(gateway.CORSMiddleware())

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
