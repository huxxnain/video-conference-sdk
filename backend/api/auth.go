package api

import (
	"net/http"
	"time"
	"strings"
	"os"
	"video-conference-sdk/backend/db"
	"video-conference-sdk/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// JWT Claims
type Claims struct {
	UserID uint `json:"user_id"`
	OrgID  uint `json:"org_id"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a JWT token for a user
func GenerateJWT(userID, orgID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}
	
	claims := Claims{
		UserID: userID,
		OrgID:  orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Authorization header required"})
			c.Abort()
			return
		}
		
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Bearer token required"})
			c.Abort()
			return
		}
		
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "default-secret-key-change-in-production"
		}
		
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token"})
			c.Abort()
			return
		}
		
		if claims, ok := token.Claims.(*Claims); ok {
			c.Set("user_id", claims.UserID)
			c.Set("org_id", claims.OrgID)
		}
		
		c.Next()
	}
}