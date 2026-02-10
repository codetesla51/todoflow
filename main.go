package main

import (
	"github.com/codetesla51/todoapi/internal/config"
	"github.com/codetesla51/todoapi/internal/controllers"
	"github.com/codetesla51/todoapi/internal/database"
	"github.com/codetesla51/todoapi/internal/middleware"
	"github.com/codetesla51/todoapi/internal/models"
	"github.com/codetesla51/todoapi/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	database.Connect()
	services.ConnectRedis()
	database.Migrate(&models.User{}, &models.Todo{})
	r := gin.Default()
	middleware.InitRateLimiter()

	// Global protection for simple DOS
	r.Use(middleware.RateLimitByIP())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.CreateUser)
		auth.POST("/login", controllers.LoginUser)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	// Use User-based limiting for logged-in users
	protected.Use(middleware.RateLimitByUser())
	{
		protected.GET("/profile", controllers.GetUser)
		protected.POST("/todos", controllers.CreateTodo)
		protected.GET("/todos", controllers.GetMyTodos)
		protected.GET("/todos/:id", controllers.GetTodo)
		protected.PUT("/todos/:id", controllers.UpdateTodo)
		protected.DELETE("/todos/:id", controllers.DeleteTodo)
	}
	r.Run(":8080")
}
