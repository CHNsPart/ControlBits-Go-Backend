package models

import "time"

type HabitEntry struct {
	ID        string    `db:"id"`
	HabitID  string    `db:"habit_id"`
	EntryDate time.Time `db:"entry_date"`
	Status    string    `db:"status"` // completed / missed
	Note      string    `db:"note"`
	CreatedAt time.Time `db:"created_at"`
}
