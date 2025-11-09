package api

import (
	"net/http"
	"time"
	"strings"
	"video-conference-sdk/backend/db"
	"video-conference-sdk/backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Signup payload
type SignupReq struct {
	OrgName  string `json:"org_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login payload
type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// POST /auth/signup
func SignupHandler(c *gin.Context) {
	var req SignupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid data", "err": err.Error()})
		return
	}
	var org models.Organization
	if err := db.DB.Where("name = ?", req.OrgName).First(&org).Error; err != nil {
		org = models.Organization{Name: req.OrgName}
		db.DB.Create(&org)
	}
	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "User/email already exists"})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user = models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		OrgID:        org.ID,
		Role:         "user",
	}
	db.DB.Create(&user)
	token, _ := GenerateJWT(user.ID, org.ID)
	c.JSON(200, gin.H{"token": token})
}

// POST /auth/login
func LoginHandler(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid data", "err": err.Error()})
		return
	}
	var user models.User
	if err := db.DB.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid credentials"})
		return
	}
	token, _ := GenerateJWT(user.ID, user.OrgID)
	c.JSON(200, gin.H{"token": token})
}