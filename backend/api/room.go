package api

import (
	"net/http"
	"time"
	"video-conference-sdk/backend/db"
	"video-conference-sdk/backend/models"
	"github.com/gin-gonic/gin"
)

// Create room payload
type CreateRoomReq struct {
	OrgID uint   `json:"org_id" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

// Join queue payload
type QueueJoinReq struct {
	UserID uint `json:"user_id" binding:"required"`
	RoomID uint `json:"room_id" binding:"required"`
}

// POST /room/create
func CreateRoomHandler(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid data", "err": err.Error()})
		return
	}
	room := models.Room{
		OrgID: req.OrgID,
		Name:  req.Name,
		Active: true,
	}
	db.DB.Create(&room)
	c.JSON(200, room)
}

// POST /room/join
func JoinQueueHandler(c *gin.Context) {
	var req QueueJoinReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid data", "err": err.Error()})
		return
	}
	entry := models.QueueEntry{
		UserID:   req.UserID,
		RoomID:   req.RoomID,
		JoinedAt: time.Now(),
	}
	db.DB.Create(&entry)
	c.JSON(200, entry)
}