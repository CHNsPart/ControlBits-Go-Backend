package models

import "time"

type UserBadge struct {
	ID       string    `db:"id"`
	UserID   string    `db:"user_id"`
	HabitID  string    `db:"habit_id"`
	BadgeID  string    `db:"badge_id"`
	EarnedAt time.Time `db:"earned_at"`
}
