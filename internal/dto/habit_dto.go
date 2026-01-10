package dto

// CreateHabitRequest is the request body for creating a habit
type CreateHabitRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateHabitResponse is the response for successful habit creation
type CreateHabitResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// UpdateHabitRequest is the request body for updating a habit
type UpdateHabitRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateHabitResponse is the response for successful habit update
type UpdateHabitResponse struct {
	Message string `json:"message"`
}

// GetHabitResponse is the response for getting a single habit
type GetHabitResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	CurrentStreak int    `json:"current_streak"`
	LongestStreak int    `json:"longest_streak"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Archived      bool   `json:"archived"`
}

// ListHabitsResponse is the response for listing all habits
type ListHabitsResponse struct {
	Habits []GetHabitResponse `json:"habits"`
}

// ArchiveHabitResponse is the response for archiving/unarchiving a habit
type ArchiveHabitResponse struct {
	Message string `json:"message"`
}

// DeleteHabitResponse is the response for deleting a habit
type DeleteHabitResponse struct {
	Message string `json:"message"`
}
