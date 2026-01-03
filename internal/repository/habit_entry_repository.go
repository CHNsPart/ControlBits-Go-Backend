package repository

import (
	"database/sql"
	"time"
)

type HabitEntryRepository struct {
	db *sql.DB
}

func NewHabitEntryRepository(db *sql.DB) *HabitEntryRepository {
	return &HabitEntryRepository{db: db}
}

// Check if entry exists for today
func (r *HabitEntryRepository) ExistsForDate(habitID string, date time.Time) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM habit_entries
			WHERE habit_id = $1 AND entry_date = $2
		)
	`

	var exists bool
	err := r.db.QueryRow(query, habitID, date).Scan(&exists)
	return exists, err
}

// Create habit entry
func (r *HabitEntryRepository) Create(habitID string, date time.Time, status string) error {
	query := `
		INSERT INTO habit_entries (id, habit_id, entry_date, status)
		VALUES (gen_random_uuid(), $1, $2, $3)
	`
	_, err := r.db.Exec(query, habitID, date, status)
	return err
}
