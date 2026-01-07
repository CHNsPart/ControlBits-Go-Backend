package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/mjubayerquanfinca/habit-tracker/internal/config"
	"github.com/mjubayerquanfinca/habit-tracker/internal/handler"
	"github.com/mjubayerquanfinca/habit-tracker/internal/middleware"
	"github.com/mjubayerquanfinca/habit-tracker/internal/repository"
	"github.com/mjubayerquanfinca/habit-tracker/internal/service"
)

func main() {
	// ---------- DB ----------
	db := config.ConnectDB()
	defer db.Close()

	// ---------- GIN ----------
	r := gin.Default()

	// ---------- AUTH ----------
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "super-secret-key")
	authHandler := handler.NewAuthHandler(authService)

	// ---------- HABITS ----------
	habitRepo := repository.NewHabitRepository(db)
	habitService := service.NewHabitService(habitRepo)
	habitHandler := handler.NewHabitHandler(habitService)

	// ---------- HABIT ENTRIES ----------
	entryRepo := repository.NewHabitEntryRepository(db)
	badgeRepo := repository.NewBadgeRepository(db)
	entryService := service.NewHabitEntryService(entryRepo, habitRepo, badgeRepo)
	entryHandler := handler.NewHabitEntryHandler(entryService)

	// ---------- ROUTES ----------
	api := r.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/forgot-password", authHandler.RequestPasswordReset)
		auth.POST("/reset-password", authHandler.ConfirmPasswordReset)
	}

	// Protected routes
	protected := api.Group("/users")
	protected.Use(middleware.JWTAuthMiddleware("super-secret-key"))
	{

		// User profile
		protected.GET("/me", authHandler.GetMe)
		protected.PUT("/me", authHandler.Update)
		protected.PUT("/me/password", authHandler.ChangePassword)
		protected.DELETE("/me", authHandler.Delete)

		// Habits CRUD
		protected.POST("/habits", habitHandler.Create)
		protected.GET("/habits", habitHandler.GetAll)
		protected.PUT("/habits/:id", habitHandler.Update)
		protected.DELETE("/habits/:id", habitHandler.Delete)

		// Habit complete / miss
		protected.POST("/habits/:id/complete", entryHandler.Complete)
		protected.POST("/habits/:id/miss", entryHandler.Miss)
	}

	log.Println("API server running on :8080")
	r.Run(":8080")
}
