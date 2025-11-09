package main

import (
	"video-conference-sdk/backend/db"
	"video-conference-sdk/backend/api"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitPostgres()
	r := gin.Default()

	r.POST("/auth/signup", api.SignupHandler)
	r.POST("/auth/login", api.LoginHandler)
	r.POST("/room/create", api.CreateRoomHandler)
	r.POST("/room/join", api.JoinQueueHandler)
	r.GET("/ws/signaling", api.SignalingHandler)

	r.Run(":8080")
}