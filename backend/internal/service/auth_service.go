package service

import (
	"database/sql"
	"errors"
	"todo-app/backend/internal/auth"
	"todo-app/backend/internal/model"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	db        *sql.DB
	jwtSecret string
}

func NewAuthService(db *sql.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *AuthService) Register(email, password string) (*model.LoginResponse, error) {
	// Check if user already exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	var user model.User
	err = s.db.QueryRow(
		`INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, 'user')
		RETURNING id, email, role, created_at, updated_at`,
		email, hashedPassword,
	).Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(email, password string) (*model.LoginResponse, error) {
	// Get user
	var user model.User
	err := s.db.QueryRow(
		`SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredentials
	}

	if err != nil {
		return nil, err
	}

	// Check password
	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetProfile retrieves user profile information
func (s *AuthService) GetProfile(userID uuid.UUID) (*model.User, error) {
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
