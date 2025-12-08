package service

import (
	"database/sql"
	"errors"
	"todo-app/backend/internal/model"

	"github.com/google/uuid"
)

var (
	ErrCannotDeleteSelf = errors.New("cannot delete your own account")
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// GetAllUsers retrieves all users (admin function)
func (s *UserService) GetAllUsers() ([]model.User, error) {
	rows, err := s.db.Query(
		`SELECT id, email, role, created_at, updated_at
		FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUser retrieves a specific user by ID (admin function)
func (s *UserService) GetUser(userID uuid.UUID) (*model.User, error) {
	var user model.User
	err := s.db.QueryRow(
		`SELECT id, email, role, created_at, updated_at
		FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserRole updates a user's role (admin function)
func (s *UserService) UpdateUserRole(userID uuid.UUID, role string) (*model.User, error) {
	var user model.User
	err := s.db.QueryRow(
		`UPDATE users SET role = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, email, role, created_at, updated_at`,
		role, userID,
	).Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser deletes a user (admin function)
func (s *UserService) DeleteUser(currentUserID, targetUserID uuid.UUID) error {
	// Prevent self-deletion
	if currentUserID == targetUserID {
		return ErrCannotDeleteSelf
	}

	result, err := s.db.Exec(`DELETE FROM users WHERE id = $1`, targetUserID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// GetAllTodos retrieves all todos from all users (admin function)
func (s *UserService) GetAllTodos() ([]model.Todo, error) {
	rows, err := s.db.Query(
		`SELECT t.id, t.user_id, t.title, t.description, t.completed, t.created_at, t.updated_at
		FROM todos t
		ORDER BY t.created_at DESC`,
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
