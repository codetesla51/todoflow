package handlers

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/codetesla51/todoapi/internal/database"
	"github.com/codetesla51/todoapi/internal/models"
	"github.com/codetesla51/todoapi/internal/services"
	"github.com/gin-gonic/gin"
)

type TodoRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"max=2000"`
}
type StatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending completed"`
}

func CreateTodo(c *gin.Context) {
	var req TodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	todo := models.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "pending",
	}
	if err := database.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}
	services.DeletePattern(fmt.Sprintf("todos:user:%d:*", userID))
	services.DeleteCache(fmt.Sprintf("user:%d", userID))
	c.JSON(http.StatusCreated, gin.H{"message": "Todo created successfully", "todo": todo})
}

func GetMyTodos(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	var todos []models.Todo
	cacheKey := fmt.Sprintf("todos:user:%d:limit:%d:page:%d", userID, limit, page)
	err := services.GetCache(cacheKey, &todos)
	if err == nil {
		c.JSON(200, todos)
		return
	}

	database.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&todos)
	jitter := time.Duration(rand.IntN(60)) * time.Second
	services.SetCache(cacheKey, todos, 10*time.Minute+jitter)
	c.JSON(200, todos)
}

func GetTodo(c *gin.Context) {
	userID := c.GetUint("user_id")
	todoID := c.Param("id")

	var todo models.Todo
	cacheKey := fmt.Sprintf("todo:%s:user:%d", todoID, userID)
	err := services.GetCache(cacheKey, &todo)
	if err == nil {
		c.JSON(200, todo)
		return
	}
	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	jitter := time.Duration(rand.IntN(60)) * time.Second
	services.SetCache(cacheKey, todo, 10*time.Minute+jitter)
	c.JSON(200, todo)
}

func UpdateTodo(c *gin.Context) {
	userID := c.GetUint("user_id")
	todoID := c.Param("id")

	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	var req TodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todo.Title = req.Title
	todo.Description = req.Description

	if err := database.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}
	services.DeletePattern(fmt.Sprintf("todos:user:%d:*", userID))
	services.DeleteCache(fmt.Sprintf("user:%d", userID))
	services.DeleteCache(fmt.Sprintf("todo:%d:user:%d", todo.ID, userID))
	c.JSON(200, todo)
}
func UpdateTodoStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	todoID := c.Param("id")

	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todo.Status = req.Status

	if err := database.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}
	services.DeletePattern(fmt.Sprintf("todos:user:%d:*", userID))
	services.DeleteCache(fmt.Sprintf("user:%d", userID))
	services.DeleteCache(fmt.Sprintf("todo:%d:user:%d", todo.ID, userID))

	c.JSON(200, todo)
}

func DeleteTodo(c *gin.Context) {
	userID := c.GetUint("user_id")
	todoID := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", todoID, userID).Delete(&models.Todo{})
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}
	services.DeletePattern(fmt.Sprintf("todos:user:%d:*", userID))
	services.DeleteCache(fmt.Sprintf("user:%d", userID))
	services.DeleteCache(fmt.Sprintf("todo:%s:user:%d", todoID, userID))

	c.JSON(200, gin.H{"message": "Todo deleted"})
}
