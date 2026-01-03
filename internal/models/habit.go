package models

import "time"

type Habit struct {
	ID            string    `db:"id"`
	UserID        string    `db:"user_id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	CurrentStreak int       `db:"current_streak"`
	LongestStreak int       `db:"longest_streak"`
	IsArchived    bool      `db:"is_archived"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
