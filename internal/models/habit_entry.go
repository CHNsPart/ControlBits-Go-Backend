package models

import "time"

type HabitEntry struct {
	ID        string    `db:"id" json:"id"`
	HabitID   string    `db:"habit_id" json:"habit_id"`
	EntryDate time.Time `db:"entry_date" json:"entry_date"`
	Status    string    `db:"status" json:"status"` // completed / missed
	Note      string    `db:"note" json:"note"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
