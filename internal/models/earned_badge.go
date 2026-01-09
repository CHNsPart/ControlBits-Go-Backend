package models

import "time"

type EarnedBadge struct {
	ID             string    `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Description    string    `db:"description" json:"description"`
	StreakRequired int       `db:"streak_required" json:"streak_required"`
	Icon           string    `db:"icon" json:"icon"`
	EarnedAt       time.Time `db:"earned_at" json:"earned_at"`
	HabitID        string    `db:"habit_id" json:"habit_id"`
}
