package controllers

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/codetesla51/todoapi/internal/database"
	"github.com/codetesla51/todoapi/internal/models"
	"github.com/codetesla51/todoapi/internal/services"
	"github.com/codetesla51/todoapi/internal/utils"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func CreateUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		UserName: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists or failed to create"})
		return
	}
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(201, gin.H{
		"message": "User created successfully",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"email":    user.Email,
		},
	})
}
func GetUser(c *gin.Context) {
	id := c.GetUint("user_id")
	var user models.User
	cacheKey := fmt.Sprintf("user:%d", id)
	err := services.GetCache(cacheKey, &user)
	if err == nil {
		c.JSON(200, user)
		return
	}
	if err := database.DB.Preload("Todos").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	jitter := time.Duration(rand.IntN(60)) * time.Second

	err = services.SetCache(cacheKey, user, 10*time.Minute+jitter)

	c.JSON(200, user)
}
func LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"email":    user.Email,
		},
	})
}
