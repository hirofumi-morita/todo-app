package main

import (
	"fmt"
	"log"
	"todo-app/backend/internal/config"
	"todo-app/backend/internal/handler"
	"todo-app/backend/internal/middleware"
	"todo-app/backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Hasura GraphQL client
	hasuraClient := service.NewHasuraClient(cfg.HasuraEndpoint, cfg.HasuraAdminSecret)

	// Initialize services
	authService := service.NewAuthService(hasuraClient, cfg.JWTSecret)
	todoService := service.NewTodoService(hasuraClient)
	userService := service.NewUserService(hasuraClient)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)
	adminHandler := handler.NewAdminHandler(userService)

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://frontend:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// User profile
		protected.GET("/profile", authHandler.GetProfile)

		// Todo routes
		protected.GET("/todos", todoHandler.GetTodos)
		protected.GET("/todos/:id", todoHandler.GetTodo)
		protected.POST("/todos", todoHandler.CreateTodo)
		protected.PUT("/todos/:id", todoHandler.UpdateTodo)
		protected.DELETE("/todos/:id", todoHandler.DeleteTodo)
	}

	// Admin routes
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	admin.Use(middleware.AdminMiddleware())
	{
		// User management
		admin.GET("/users", adminHandler.GetAllUsers)
		admin.GET("/users/:id", adminHandler.GetUser)
		admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
		admin.DELETE("/users/:id", adminHandler.DeleteUser)

		// All todos
		admin.GET("/todos", adminHandler.GetAllTodos)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
