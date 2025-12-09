package service

import (
	"errors"
	"todo-app/backend/internal/model"

	"github.com/google/uuid"
)

var (
	ErrTodoNotFound  = errors.New("todo not found")
	ErrInvalidTodoID = errors.New("invalid todo id")
)

type TodoService struct {
	hasura *HasuraClient
}

func NewTodoService(hasura *HasuraClient) *TodoService {
	return &TodoService{hasura: hasura}
}

// GetTodos retrieves all todos for a user
func (s *TodoService) GetTodos(userID uuid.UUID) ([]model.Todo, error) {
	var response struct {
		Todos []model.Todo `json:"todos"`
	}

	err := s.hasura.execute(`
        query ($userId: uuid!) {
          todos(where: {user_id: {_eq: $userId}}, order_by: {created_at: desc}) {
            id
            user_id
            title
            description
            completed
            created_at
            updated_at
          }
        }
        `, map[string]interface{}{"userId": userID}, &response)
	if err != nil {
		return nil, err
	}

	return response.Todos, nil
}

// GetTodo retrieves a specific todo for a user
func (s *TodoService) GetTodo(userID, todoID uuid.UUID) (*model.Todo, error) {
	var response struct {
		Todos []model.Todo `json:"todos"`
	}

	err := s.hasura.execute(`
        query ($id: uuid!, $userId: uuid!) {
          todos(where: {id: {_eq: $id}, user_id: {_eq: $userId}}, limit: 1) {
            id
            user_id
            title
            description
            completed
            created_at
            updated_at
          }
        }
        `, map[string]interface{}{"id": todoID, "userId": userID}, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Todos) == 0 {
		return nil, ErrTodoNotFound
	}

	return &response.Todos[0], nil
}

// CreateTodo creates a new todo for a user
func (s *TodoService) CreateTodo(userID uuid.UUID, title string, description *string) (*model.Todo, error) {
	var response struct {
		InsertTodosOne model.Todo `json:"insert_todos_one"`
	}

	err := s.hasura.execute(`
        mutation ($userId: uuid!, $title: String!, $description: String) {
          insert_todos_one(object: {user_id: $userId, title: $title, description: $description}) {
            id
            user_id
            title
            description
            completed
            created_at
            updated_at
          }
        }
        `, map[string]interface{}{"userId": userID, "title": title, "description": description}, &response)
	if err != nil {
		return nil, err
	}

	return &response.InsertTodosOne, nil
}

// UpdateTodo updates a todo for a user
func (s *TodoService) UpdateTodo(userID, todoID uuid.UUID, req model.UpdateTodoRequest) (*model.Todo, error) {
	changes := map[string]interface{}{}

	if req.Title != nil {
		changes["title"] = *req.Title
	}

	if req.Description != nil {
		changes["description"] = req.Description
	}

	if req.Completed != nil {
		changes["completed"] = *req.Completed
	}

	if len(changes) == 0 {
		return s.GetTodo(userID, todoID)
	}

	var response struct {
		UpdateTodos struct {
			Returning []model.Todo `json:"returning"`
		} `json:"update_todos"`
	}

	err := s.hasura.execute(`
        mutation ($id: uuid!, $userId: uuid!, $changes: todos_set_input!) {
          update_todos(where: {id: {_eq: $id}, user_id: {_eq: $userId}}, _set: $changes) {
            returning {
              id
              user_id
              title
              description
              completed
              created_at
              updated_at
            }
          }
        }
        `, map[string]interface{}{"id": todoID, "userId": userID, "changes": changes}, &response)
	if err != nil {
		return nil, err
	}

	if len(response.UpdateTodos.Returning) == 0 {
		return nil, ErrTodoNotFound
	}

	return &response.UpdateTodos.Returning[0], nil
}

// DeleteTodo deletes a todo for a user
func (s *TodoService) DeleteTodo(userID, todoID uuid.UUID) error {
	var response struct {
		DeleteTodos struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"delete_todos"`
	}

	err := s.hasura.execute(`
        mutation ($id: uuid!, $userId: uuid!) {
          delete_todos(where: {id: {_eq: $id}, user_id: {_eq: $userId}}) {
            affected_rows
          }
        }
        `, map[string]interface{}{"id": todoID, "userId": userID}, &response)
	if err != nil {
		return err
	}

	if response.DeleteTodos.AffectedRows == 0 {
		return ErrTodoNotFound
	}

	return nil
}
