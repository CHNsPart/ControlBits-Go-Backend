package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/mjubayerquanfinca/habit-tracker/docs"
	"github.com/mjubayerquanfinca/habit-tracker/internal/config"
	"github.com/mjubayerquanfinca/habit-tracker/internal/handler"
	"github.com/mjubayerquanfinca/habit-tracker/internal/middleware"
	"github.com/mjubayerquanfinca/habit-tracker/internal/repository"
	"github.com/mjubayerquanfinca/habit-tracker/internal/service"
)

// @title ControlBits API
// @version 1.0.0
// @description Backend API for ControlBits habit tracking.
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// ---------- DB ----------
	db := config.ConnectDB()
	defer db.Close()

	// ---------- GIN ----------
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	// Swagger UI + spec
	r.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
	))

	// ---------- AUTH ----------
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, "super-secret-key")
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// ---------- HABITS ----------
	habitRepo := repository.NewHabitRepository(db)
	habitService := service.NewHabitService(habitRepo)
	habitHandler := handler.NewHabitHandler(habitService)

	// ---------- HABIT ENTRIES ----------
	entryRepo := repository.NewHabitEntryRepository(db)
	badgeRepo := repository.NewBadgeRepository(db)
	entryService := service.NewHabitEntryService(entryRepo, habitRepo, badgeRepo)
	entryHandler := handler.NewHabitEntryHandler(entryService)
	badgeService := service.NewBadgeService(badgeRepo)
	badgeHandler := handler.NewBadgeHandler(badgeService)

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
		protected.GET("/me", userHandler.GetMe)
		protected.PUT("/me", userHandler.Update)
		protected.PUT("/me/password", authHandler.ChangePassword)
		protected.DELETE("/me", userHandler.Delete)

		// User badges
		protected.GET("/me/badges", badgeHandler.ListUserBadges)
	}

	// Protected routes without /users prefix
	protectedAPI := api.Group("/")
	protectedAPI.Use(middleware.JWTAuthMiddleware("super-secret-key"))
	{
		// Habits CRUD
		protectedAPI.POST("/habits", habitHandler.Create)
		protectedAPI.GET("/habits", habitHandler.GetAll)
		protectedAPI.GET("/habits/:id", habitHandler.GetByID)
		protectedAPI.PUT("/habits/:id", habitHandler.Update)
		protectedAPI.DELETE("/habits/:id", habitHandler.Delete)
		protectedAPI.POST("/habits/:id/archive", habitHandler.Archive)
		protectedAPI.POST("/habits/:id/unarchive", habitHandler.Unarchive)

		// Habit complete / miss
		protectedAPI.POST("/habits/:id/complete", entryHandler.Complete)
		protectedAPI.POST("/habits/:id/miss", entryHandler.Miss)
		protectedAPI.GET("/habits/:id/entries", entryHandler.ListEntries)
		protectedAPI.POST("/habits/:id/entries", entryHandler.CreateEntry)
		protectedAPI.DELETE("/habits/:id/entries/:entryId", entryHandler.DeleteEntry)

		// Badges
		protectedAPI.GET("/badges", badgeHandler.ListAll)
		protectedAPI.GET("/habits/:id/badges", badgeHandler.ListHabitBadges)
	}

	// Start HTTPS server
	go func() {
		if err := r.RunTLS(":443", "cert.pem", "key.pem"); err != nil {
			log.Fatal(err)
		}
	}()

	// Redirect HTTP to HTTPS
	if err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "https://"+req.Host+req.RequestURI, http.StatusMovedPermanently)
	})); err != nil {
		log.Fatal(err)
	}
}
