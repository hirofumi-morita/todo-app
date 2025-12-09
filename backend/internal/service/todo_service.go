package service

import (
	"database/sql"
	"errors"
	"fmt"
	"todo-app/backend/internal/model"

	"github.com/google/uuid"
)

var (
	ErrTodoNotFound  = errors.New("todo not found")
	ErrInvalidTodoID = errors.New("invalid todo id")
)

type TodoService struct {
	db *sql.DB
}

func NewTodoService(db *sql.DB) *TodoService {
	return &TodoService{db: db}
}

// GetTodos retrieves all todos for a user
func (s *TodoService) GetTodos(userID uuid.UUID) ([]model.Todo, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// GetTodo retrieves a specific todo for a user
func (s *TodoService) GetTodo(userID, todoID uuid.UUID) (*model.Todo, error) {
	var todo model.Todo
	err := s.db.QueryRow(
		`SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos WHERE id = $1 AND user_id = $2`,
		todoID, userID,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrTodoNotFound
	}

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// CreateTodo creates a new todo for a user
func (s *TodoService) CreateTodo(userID uuid.UUID, title string, description *string) (*model.Todo, error) {
	var todo model.Todo
	err := s.db.QueryRow(
		`INSERT INTO todos (user_id, title, description)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, title, description, completed, created_at, updated_at`,
		userID, title, description,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// UpdateTodo updates a todo for a user
func (s *TodoService) UpdateTodo(userID, todoID uuid.UUID, req model.UpdateTodoRequest) (*model.Todo, error) {
	// Check if todo exists and belongs to user
	var exists bool
	err := s.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM todos WHERE id = $1 AND user_id = $2)`,
		todoID, userID,
	).Scan(&exists)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrTodoNotFound
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
	args = append(args, todoID, userID)

	// Execute update
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// Fetch updated todo
	var todo model.Todo
	err = s.db.QueryRow(
		`SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos WHERE id = $1`,
		todoID,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// DeleteTodo deletes a todo for a user
func (s *TodoService) DeleteTodo(userID, todoID uuid.UUID) error {
	result, err := s.db.Exec(
		`DELETE FROM todos WHERE id = $1 AND user_id = $2`,
		todoID, userID,
	)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}
