package models

type Badge struct {
	ID             string `db:"id"`
	Name           string `db:"name"`
	Description    string `db:"description"`
	StreakRequired int    `db:"streak_required"`
	Icon           string `db:"icon"`
}
