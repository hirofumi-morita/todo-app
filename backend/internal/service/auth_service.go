package service

import (
	"errors"
	"time"
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
	hasura    *HasuraClient
	jwtSecret string
}

func NewAuthService(hasura *HasuraClient, jwtSecret string) *AuthService {
	return &AuthService{
		hasura:    hasura,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *AuthService) Register(email, password string) (*model.LoginResponse, error) {
	// Check if user already exists
	var existsResp struct {
		Users []struct {
			ID uuid.UUID `json:"id"`
		} `json:"users"`
	}
	err := s.hasura.execute(`
        query ($email: String!) {
          users(where: {email: {_eq: $email}}, limit: 1) {
            id
          }
        }
        `, map[string]interface{}{"email": email}, &existsResp)
	if err != nil {
		return nil, err
	}

	if len(existsResp.Users) > 0 {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	var userResp struct {
		InsertUsersOne model.User `json:"insert_users_one"`
	}
	err = s.hasura.execute(`
        mutation ($email: String!, $password_hash: String!) {
          insert_users_one(object: {email: $email, password_hash: $password_hash, role: "user"}) {
            id
            email
            role
            created_at
            updated_at
          }
        }
        `, map[string]interface{}{"email": email, "password_hash": hashedPassword}, &userResp)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := auth.GenerateToken(userResp.InsertUsersOne.ID, userResp.InsertUsersOne.Email, userResp.InsertUsersOne.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  userResp.InsertUsersOne,
	}, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(email, password string) (*model.LoginResponse, error) {
	// Get user
	var response struct {
		Users []struct {
			ID           uuid.UUID `json:"id"`
			Email        string    `json:"email"`
			PasswordHash string    `json:"password_hash"`
			Role         string    `json:"role"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
		} `json:"users"`
	}

	err := s.hasura.execute(`
        query ($email: String!) {
          users(where: {email: {_eq: $email}}, limit: 1) {
            id
            email
            password_hash
            role
            created_at
            updated_at
          }
        }
        `, map[string]interface{}{"email": email}, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Users) == 0 {
		return nil, ErrInvalidCredentials
	}

	userRecord := response.Users[0]

	// Check password
	if !auth.CheckPassword(password, userRecord.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	user := model.User{
		ID:        userRecord.ID,
		Email:     userRecord.Email,
		Role:      userRecord.Role,
		CreatedAt: userRecord.CreatedAt,
		UpdatedAt: userRecord.UpdatedAt,
	}

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
	var response struct {
		UsersByPk *model.User `json:"users_by_pk"`
	}

	err := s.hasura.execute(`
        query ($id: uuid!) {
          users_by_pk(id: $id) {
            id
            email
            role
            created_at
            updated_at
          }
        }
        `, map[string]interface{}{"id": userID}, &response)
	if err != nil {
		return nil, err
	}

	if response.UsersByPk == nil {
		return nil, ErrUserNotFound
	}

	return response.UsersByPk, nil
}
