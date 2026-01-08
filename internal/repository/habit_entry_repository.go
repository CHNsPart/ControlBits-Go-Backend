package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
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

// List entries for a habit with optional date range
func (r *HabitEntryRepository) ListByHabit(habitID string, startDate, endDate *time.Time) ([]models.HabitEntry, error) {
	var args []interface{}
	var conditions []string
	args = append(args, habitID)
	conditions = append(conditions, "habit_id = $1")

	idx := 2
	if startDate != nil {
		conditions = append(conditions, fmt.Sprintf("entry_date >= $%d", idx))
		args = append(args, *startDate)
		idx++
	}
	if endDate != nil {
		conditions = append(conditions, fmt.Sprintf("entry_date <= $%d", idx))
		args = append(args, *endDate)
		idx++
	}

	query := `SELECT id, habit_id, entry_date, status, note, created_at FROM habit_entries WHERE ` + strings.Join(conditions, " AND ") + ` ORDER BY entry_date DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.HabitEntry
	for rows.Next() {
		var e models.HabitEntry
		err := rows.Scan(&e.ID, &e.HabitID, &e.EntryDate, &e.Status, &e.Note, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Create entry for arbitrary date/status
func (r *HabitEntryRepository) CreateEntry(habitID string, date time.Time, status, note string) error {
	query := `
		INSERT INTO habit_entries (id, habit_id, entry_date, status, note, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW())
	`
	_, err := r.db.Exec(query, habitID, date, status, note)
	return err
}

// Delete entry by ID
func (r *HabitEntryRepository) DeleteEntry(entryID, habitID string) error {
	query := `
		DELETE FROM habit_entries WHERE id = $1 AND habit_id = $2
	`
	_, err := r.db.Exec(query, entryID, habitID)
	return err
}
