package handler

import (
	"time"

	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
)

// ErrorResponse documents the standard error response payload.
type ErrorResponse struct {
	Error string `json:"error"`
}

// UserResponse documents the user profile response payload.
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserRequest documents the update user request payload.
type UpdateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UpdateUserResponse documents the update user response payload.
type UpdateUserResponse struct {
	Message string `json:"message"`
}

// DeleteUserResponse documents the delete user response payload.
type DeleteUserResponse struct {
	Message string `json:"message"`
}

// LogoutResponse documents the logout response payload.
type LogoutResponse struct {
	Message string `json:"message"`
}

// CompleteHabitResponse documents the habit completion response payload.
type CompleteHabitResponse struct {
	Message   string         `json:"message"`
	NewBadges []models.Badge `json:"new_badges"`
}

// MissHabitResponse documents the habit miss response payload.
type MissHabitResponse struct {
	Message string `json:"message"`
}

// BadgeResponse documents the badge payload.
type BadgeResponse = models.Badge

// EarnedBadgeResponse documents the earned badge payload.
type EarnedBadgeResponse = models.EarnedBadge
