package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codetesla51/todoapi/internal/config"
	"github.com/codetesla51/todoapi/internal/handlers"

	"github.com/codetesla51/todoapi/internal/database"
	"github.com/codetesla51/todoapi/internal/middleware"
	"github.com/codetesla51/todoapi/internal/models"
	"github.com/codetesla51/todoapi/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	database.Connect()
	services.ConnectRedis()
	database.Migrate(&models.User{}, &models.Todo{})
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"https://todoflow-black.vercel.app/", "http://localhost:5173"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Length", "X-RateLimit-Limit", "X-RateLimit-Remaining", "Retry-After"},
	}))

	middleware.InitRateLimiter()

	// Global protection for simple DOS
	r.Use(middleware.RateLimitByIP())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.CreateUser)
		auth.POST("/login", handlers.LoginUser)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	// Use User-based limiting for logged-in users
	protected.Use(middleware.RateLimitByUser())
	{
		protected.GET("/profile", handlers.GetUser)
		protected.POST("/todos", handlers.CreateTodo)
		protected.GET("/todos", handlers.GetMyTodos)
		protected.GET("/todos/:id", handlers.GetTodo)
		protected.PUT("/todos/:id", handlers.UpdateTodo)
		protected.PATCH("/todos/:id/status", handlers.UpdateTodoStatus)
		protected.DELETE("/todos/:id", handlers.DeleteTodo)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server started on :8080")

	<-ctx.Done()
	log.Println("Shutting down gracefully")

	stop()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
