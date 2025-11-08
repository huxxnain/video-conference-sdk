package main

import (
	"github.com/gin-gonic/gin"
	"video-conference-sdk/backend/config"
	"video-conference-sdk/backend/db"
	"video-conference-sdk/backend/api"
	"video-conference-sdk/backend/ws"
)

func main() {
	config.Load()
	db.SetupDatabase()
	r := gin.Default()
	api.RegisterRoutes(r)
	ws.RegisterWebSocketRoutes(r)
	r.Run(":" + config.ServerPort)
}