package models

type Badge struct {
	ID             string `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Description    string `db:"description" json:"description"`
	StreakRequired int    `db:"streak_required" json:"streak_required"`
	Icon           string `db:"icon" json:"icon"`
}
