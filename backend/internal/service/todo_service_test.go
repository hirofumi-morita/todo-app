package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"todo-app/backend/internal/model"
)

type mockResponse struct {
	status int
	body   string
}

func newMockHasuraClient(t *testing.T, responses []mockResponse) (*HasuraClient, func()) {
	t.Helper()

	var idx int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		if idx >= len(responses) {
			t.Fatalf("received more requests than expected")
		}

		resp := responses[idx]
		idx++

		status := resp.status
		if status == 0 {
			status = http.StatusOK
		}

		w.WriteHeader(status)
		if _, err := w.Write([]byte(resp.body)); err != nil {
			t.Fatalf("failed to write mock response: %v", err)
		}
	}))

	hasura := &HasuraClient{
		endpoint:   server.URL,
		httpClient: server.Client(),
	}

	return hasura, server.Close
}

func TestTodoService_GetTodos(t *testing.T) {
	userID := uuid.New()
	todoID1 := uuid.New()
	todoID2 := uuid.New()
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	client, shutdown := newMockHasuraClient(t, []mockResponse{
		{
			body: fmt.Sprintf(`{"data":{"todos":[{"id":"%s","user_id":"%s","title":"First","description":"desc1","completed":false,"created_at":"%s","updated_at":"%s"},{"id":"%s","user_id":"%s","title":"Second","description":null,"completed":true,"created_at":"%s","updated_at":"%s"}]}}`, todoID1, userID, now, now, todoID2, userID, now, now),
		},
	})
	defer shutdown()

	service := NewTodoService(client)
	todos, err := service.GetTodos(userID)
	if err != nil {
		t.Fatalf("GetTodos returned error: %v", err)
	}

	if len(todos) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(todos))
	}

	if todos[0].ID != todoID1 || todos[0].UserID != userID || todos[0].Title != "First" || todos[0].Completed {
		t.Fatalf("unexpected first todo: %+v", todos[0])
	}

	if todos[1].ID != todoID2 || !todos[1].Completed {
		t.Fatalf("unexpected second todo: %+v", todos[1])
	}
}

func TestTodoService_GetTodo(t *testing.T) {
	userID := uuid.New()
	todoID := uuid.New()
	now := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	t.Run("found", func(t *testing.T) {
		client, shutdown := newMockHasuraClient(t, []mockResponse{
			{
				body: fmt.Sprintf(`{"data":{"todos":[{"id":"%s","user_id":"%s","title":"Item","description":"detail","completed":false,"created_at":"%s","updated_at":"%s"}]}}`, todoID, userID, now, now),
			},
		})
		defer shutdown()

		service := NewTodoService(client)
		todo, err := service.GetTodo(userID, todoID)
		if err != nil {
			t.Fatalf("GetTodo returned error: %v", err)
		}

		if todo.ID != todoID || todo.UserID != userID || todo.Title != "Item" {
			t.Fatalf("unexpected todo: %+v", todo)
		}
	})

	t.Run("not found", func(t *testing.T) {
		client, shutdown := newMockHasuraClient(t, []mockResponse{
			{body: `{"data":{"todos":[]}}`},
		})
		defer shutdown()

		service := NewTodoService(client)
		_, err := service.GetTodo(userID, todoID)
		if !errors.Is(err, ErrTodoNotFound) {
			t.Fatalf("expected ErrTodoNotFound, got %v", err)
		}
	})
}

func TestTodoService_CreateTodo(t *testing.T) {
	userID := uuid.New()
	todoID := uuid.New()
	now := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	client, shutdown := newMockHasuraClient(t, []mockResponse{
		{
			body: fmt.Sprintf(`{"data":{"insert_todos_one":{"id":"%s","user_id":"%s","title":"New","description":"created","completed":false,"created_at":"%s","updated_at":"%s"}}}`, todoID, userID, now, now),
		},
	})
	defer shutdown()

	service := NewTodoService(client)
	desc := "created"
	todo, err := service.CreateTodo(userID, "New", &desc)
	if err != nil {
		t.Fatalf("CreateTodo returned error: %v", err)
	}

	if todo.ID != todoID || todo.UserID != userID || todo.Title != "New" || todo.Description == nil || *todo.Description != desc {
		t.Fatalf("unexpected todo: %+v", todo)
	}
}

func TestTodoService_UpdateTodo(t *testing.T) {
	userID := uuid.New()
	todoID := uuid.New()
	now := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	t.Run("with changes", func(t *testing.T) {
		client, shutdown := newMockHasuraClient(t, []mockResponse{
			{
				body: fmt.Sprintf(`{"data":{"update_todos":{"returning":[{"id":"%s","user_id":"%s","title":"Updated","description":"after","completed":true,"created_at":"%s","updated_at":"%s"}]}}}`, todoID, userID, now, now),
			},
		})
		defer shutdown()

		service := NewTodoService(client)
		desc := "after"
		req := model.UpdateTodoRequest{
			Title:       strPtr("Updated"),
			Description: &desc,
			Completed:   boolPtr(true),
		}

		todo, err := service.UpdateTodo(userID, todoID, req)
		if err != nil {
			t.Fatalf("UpdateTodo returned error: %v", err)
		}

		if todo.Title != "Updated" || todo.Description == nil || *todo.Description != "after" || !todo.Completed {
			t.Fatalf("unexpected todo: %+v", todo)
		}
	})

	t.Run("no changes returns current todo", func(t *testing.T) {
		client, shutdown := newMockHasuraClient(t, []mockResponse{
			{
				body: fmt.Sprintf(`{"data":{"todos":[{"id":"%s","user_id":"%s","title":"Existing","description":"keep","completed":false,"created_at":"%s","updated_at":"%s"}]}}`, todoID, userID, now, now),
			},
		})
		defer shutdown()

		service := NewTodoService(client)
		todo, err := service.UpdateTodo(userID, todoID, model.UpdateTodoRequest{})
		if err != nil {
			t.Fatalf("UpdateTodo returned error: %v", err)
		}

		if todo.Title != "Existing" || todo.Description == nil || *todo.Description != "keep" {
			t.Fatalf("unexpected todo: %+v", todo)
		}
	})

	t.Run("not found", func(t *testing.T) {
		client, shutdown := newMockHasuraClient(t, []mockResponse{
			{
				body: `{"data":{"update_todos":{"returning":[]}}}`,
			},
		})
		defer shutdown()

		service := NewTodoService(client)
		_, err := service.UpdateTodo(userID, todoID, model.UpdateTodoRequest{Title: strPtr("missing")})
		if !errors.Is(err, ErrTodoNotFound) {
			t.Fatalf("expected ErrTodoNotFound, got %v", err)
		}
	})
}

func TestTodoService_DeleteTodo(t *testing.T) {
	userID := uuid.New()
	todoID := uuid.New()

	t.Run("delete succeeds", func(t *testing.T) {
		client, shutdown := newMockHasuraClient(t, []mockResponse{
			{body: `{"data":{"delete_todos":{"affected_rows":1}}}`},
		})
		defer shutdown()

		service := NewTodoService(client)
		if err := service.DeleteTodo(userID, todoID); err != nil {
			t.Fatalf("DeleteTodo returned error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		client, shutdown := newMockHasuraClient(t, []mockResponse{
			{body: `{"data":{"delete_todos":{"affected_rows":0}}}`},
		})
		defer shutdown()

		service := NewTodoService(client)
		if err := service.DeleteTodo(userID, todoID); !errors.Is(err, ErrTodoNotFound) {
			t.Fatalf("expected ErrTodoNotFound, got %v", err)
		}
	})
}

func strPtr(v string) *string { return &v }
func boolPtr(v bool) *bool    { return &v }
