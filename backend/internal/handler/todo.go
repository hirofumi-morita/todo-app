package handler

import (
	"errors"
	"net/http"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	todoService *service.TodoService
}

func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	todos, err := h.todoService.GetTodos(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetTodo(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	todoID := c.Param("id")

	todoUUID, err := uuid.Parse(todoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	todo, err := h.todoService.GetTodo(userID, todoUUID)
	if err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch todo"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req model.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.todoService.CreateTodo(userID, req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	todoID := c.Param("id")

	todoUUID, err := uuid.Parse(todoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	var req model.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.todoService.UpdateTodo(userID, todoUUID, req)
	if err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	todoID := c.Param("id")

	todoUUID, err := uuid.Parse(todoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	err = h.todoService.DeleteTodo(userID, todoUUID)
	if err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "todo deleted successfully"})
}
