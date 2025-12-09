package service

import (
	"errors"
	"todo-app/backend/internal/model"

	"github.com/google/uuid"
)

var (
	ErrCannotDeleteSelf = errors.New("cannot delete your own account")
)

type UserService struct {
	hasura *HasuraClient
}

func NewUserService(hasura *HasuraClient) *UserService {
	return &UserService{hasura: hasura}
}

// GetAllUsers retrieves all users (admin function)
func (s *UserService) GetAllUsers() ([]model.User, error) {
	var response struct {
		Users []model.User `json:"users"`
	}

	err := s.hasura.execute(`
        query {
          users(order_by: {created_at: desc}) {
            id
            email
            role
            created_at
            updated_at
          }
        }
        `, nil, &response)
	if err != nil {
		return nil, err
	}

	return response.Users, nil
}

// GetUser retrieves a specific user by ID (admin function)
func (s *UserService) GetUser(userID uuid.UUID) (*model.User, error) {
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

// UpdateUserRole updates a user's role (admin function)
func (s *UserService) UpdateUserRole(userID uuid.UUID, role string) (*model.User, error) {
	var response struct {
		UpdateUsersByPk *model.User `json:"update_users_by_pk"`
	}

	err := s.hasura.execute(`
        mutation ($id: uuid!, $role: String!) {
          update_users_by_pk(pk_columns: {id: $id}, _set: {role: $role}) {
            id
            email
            role
            created_at
            updated_at
          }
        }
        `, map[string]interface{}{"id": userID, "role": role}, &response)
	if err != nil {
		return nil, err
	}

	if response.UpdateUsersByPk == nil {
		return nil, ErrUserNotFound
	}

	return response.UpdateUsersByPk, nil
}

// DeleteUser deletes a user (admin function)
func (s *UserService) DeleteUser(currentUserID, targetUserID uuid.UUID) error {
	// Prevent self-deletion
	if currentUserID == targetUserID {
		return ErrCannotDeleteSelf
	}

	var response struct {
		DeleteUsersByPk *struct {
			ID uuid.UUID `json:"id"`
		} `json:"delete_users_by_pk"`
	}

	err := s.hasura.execute(`
        mutation ($id: uuid!) {
          delete_users_by_pk(id: $id) {
            id
          }
        }
        `, map[string]interface{}{"id": targetUserID}, &response)
	if err != nil {
		return err
	}

	if response.DeleteUsersByPk == nil {
		return ErrUserNotFound
	}

	return nil
}

// GetAllTodos retrieves all todos from all users (admin function)
func (s *UserService) GetAllTodos() ([]model.Todo, error) {
	var response struct {
		Todos []model.Todo `json:"todos"`
	}

	err := s.hasura.execute(`
        query {
          todos(order_by: {created_at: desc}) {
            id
            user_id
            title
            description
            completed
            created_at
            updated_at
          }
        }
        `, nil, &response)
	if err != nil {
		return nil, err
	}

	return response.Todos, nil
}
