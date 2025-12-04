package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"todo-app/backend/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	db *sql.DB
}

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	rows, err := h.db.Query(
		`SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch todos"})
		return
	}
	defer rows.Close()

	todos := []model.Todo{}
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(
			&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
			&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan todo"})
			return
		}
		todos = append(todos, todo)
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

	var todo model.Todo
	err = h.db.QueryRow(
		`SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos WHERE id = $1 AND user_id = $2`,
		todoUUID, userID,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	if err != nil {
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

	var todo model.Todo
	err := h.db.QueryRow(
		`INSERT INTO todos (user_id, title, description)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, title, description, completed, created_at, updated_at`,
		userID, req.Title, req.Description,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

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

	// Check if todo exists and belongs to user
	var exists bool
	err = h.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM todos WHERE id = $1 AND user_id = $2)`,
		todoUUID, userID,
	).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	// Build update query dynamically
	query := "UPDATE todos SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 1

	if req.Title != nil {
		query += ", title = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.Title)
		argCount++
	}

	if req.Description != nil {
		query += ", description = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.Description)
		argCount++
	}

	if req.Completed != nil {
		query += ", completed = $" + fmt.Sprintf("%d", argCount)
		args = append(args, *req.Completed)
		argCount++
	}

	query += " WHERE id = $" + fmt.Sprintf("%d", argCount) + " AND user_id = $" + fmt.Sprintf("%d", argCount+1)
	args = append(args, todoUUID, userID)

	// Execute update
	_, err = h.db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo"})
		return
	}

	// Fetch updated todo
	var todo model.Todo
	err = h.db.QueryRow(
		`SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos WHERE id = $1`,
		todoUUID,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated todo"})
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

	result, err := h.db.Exec(
		`DELETE FROM todos WHERE id = $1 AND user_id = $2`,
		todoUUID, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete todo"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "todo deleted successfully"})
}
