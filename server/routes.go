package server

import (
	_ "main/docs"
	"main/server/gateway"
	"main/server/handler"

	admin "main/server/handler/admin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {

	//Allowing CORS
	server.engine.Use(gateway.CORSMiddleware())

	//Auth routes
	server.engine.POST("/guest-login", admin.GuestLoginHandler)
	server.engine.POST("/login", admin.LoginHandler)
	server.engine.POST("/player-login", admin.PlayerLoginHandler)
	server.engine.PUT("/update-email", gateway.AdminAuthorization, admin.UpdateEmailHandler)
	server.engine.PATCH("/update-pass", gateway.AdminAuthorization, admin.UpdatePasswordHandler)
	server.engine.GET("/admin", admin.GetAdminHandler)
	server.engine.DELETE("/logout", gateway.AdminAuthorization, admin.LogutHandler)
	server.engine.DELETE("/delete-account", gateway.AdminAuthorization, admin.DeleteAccountHandler)

	//Admin garage routes
	server.engine.POST("/admin/garage/add", gateway.AdminAuthorization, admin.AddGarageHandler)
	server.engine.DELETE("/admin/garage/delete", gateway.AdminAuthorization, admin.DeleteGarageHandler)
	server.engine.PUT("/admin/garage/update", gateway.AdminAuthorization, admin.UpdateGarageHandler)
	server.engine.GET("/garage/types", admin.GetGarageTypesHandler)
	server.engine.GET("/garage/rarity", admin.GetRarityHandler)
	server.engine.GET("/garages/get-all", admin.GetAllGarageListHandler)

	//Admin Battle Arena Routes
	server.engine.POST("/admin/arena", gateway.AdminAuthorization, admin.AddArenaHandler)
	server.engine.DELETE("/admin/arena", gateway.AdminAuthorization, admin.DeleteArenaHandler)
	server.engine.PUT("/admin/arena", gateway.AdminAuthorization, admin.UpdateArenaHandler)
	server.engine.GET("/arena/get", admin.GetArenaListHandler)
	server.engine.GET("/arena/types", admin.GetArenaTypeHandler)

	//Shop routes
	server.engine.GET("/get-shop", handler.GetShopHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
