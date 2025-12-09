package handler

import (
	"errors"
	"net/http"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	userService *service.UserService
}

func NewAdminHandler(userService *service.UserService) *AdminHandler {
	return &AdminHandler{userService: userService}
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userService.GetUser(userUUID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required,oneof=user admin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUserRole(userUUID, req.Role)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	currentUserID := c.MustGet("user_id").(uuid.UUID)

	err = h.userService.DeleteUser(currentUserID, userUUID)
	if err != nil {
		if errors.Is(err, service.ErrCannotDeleteSelf) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete your own account"})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *AdminHandler) GetAllTodos(c *gin.Context) {
	todos, err := h.userService.GetAllTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}
