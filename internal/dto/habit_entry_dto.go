package dto

// CreateHabitEntryRequest is the request body for creating a habit entry
type CreateHabitEntryRequest struct {
	Date   string `json:"date" binding:"required"`
	Status string `json:"status" binding:"required"`
	Note   string `json:"note"`
}

// CreateHabitEntryResponse is the response for successful habit entry creation
type CreateHabitEntryResponse struct {
	ID        string   `json:"id"`
	Message   string   `json:"message"`
	NewBadges []string `json:"new_badges"`
}

// ListHabitEntriesResponse is the response for listing habit entries
type ListHabitEntriesResponse struct {
	Entries []HabitEntryResponse `json:"entries"`
}

type HabitEntryResponse struct {
	ID        string `json:"id"`
	HabitID   string `json:"habit_id"`
	EntryDate string `json:"entry_date"`
	Status    string `json:"status"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
}

// DeleteHabitEntryResponse is the response for deleting a habit entry
type DeleteHabitEntryResponse struct {
	Message string `json:"message"`
}
