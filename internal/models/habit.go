package models

import "time"

type Habit struct {
	ID            string    `db:"id" json:"id"`
	UserID        string    `db:"user_id" json:"user_id"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	CurrentStreak int       `db:"current_streak" json:"current_streak"`
	LongestStreak int       `db:"longest_streak" json:"longest_streak"`
	IsArchived    bool      `db:"is_archived" json:"is_archived"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
