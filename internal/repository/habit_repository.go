package repository

import (
	"database/sql"
	"time"

	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
)

type HabitRepository struct {
	db *sql.DB
}

func NewHabitRepository(db *sql.DB) *HabitRepository {
	return &HabitRepository{db: db}
}

// CREATE
func (r *HabitRepository) Create(habit *models.Habit) (string, error) {
       query := `
	       INSERT INTO habits (id, user_id, name, description)
	       VALUES ($1, $2, $3, $4)
	       RETURNING id
       `
       var id string
       err := r.db.QueryRow(
	       query,
	       habit.ID,
	       habit.UserID,
	       habit.Name,
	       habit.Description,
       ).Scan(&id)
       return id, err
}

// READ (all habits of a user)
func (r *HabitRepository) FindByUser(userID string) ([]models.Habit, error) {
	query := `
		SELECT id, user_id, name, description,
		       current_streak, longest_streak, is_archived,
		       created_at, updated_at
		FROM habits
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []models.Habit
	for rows.Next() {
		var h models.Habit
		err := rows.Scan(
			&h.ID,
			&h.UserID,
			&h.Name,
			&h.Description,
			&h.CurrentStreak,
			&h.LongestStreak,
			&h.IsArchived,
			&h.CreatedAt,
			&h.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		habits = append(habits, h)
	}

	return habits, nil
}

// READ (single habit by id)
func (r *HabitRepository) GetByID(habitID, userID string) (*models.Habit, error) {
	query := `
		SELECT id, user_id, name, description,
			   current_streak, longest_streak, is_archived,
			   created_at, updated_at
		FROM habits
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`
	var h models.Habit
	err := r.db.QueryRow(query, habitID, userID).Scan(
		&h.ID,
		&h.UserID,
		&h.Name,
		&h.Description,
		&h.CurrentStreak,
		&h.LongestStreak,
		&h.IsArchived,
		&h.CreatedAt,
		&h.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &h, nil
}

// UPDATE
func (r *HabitRepository) Update(habit *models.Habit) error {
	query := `
		UPDATE habits
		SET name = $1,
		    description = $2,
		    updated_at = NOW()
		WHERE id = $3 AND user_id = $4
	`
	_, err := r.db.Exec(
		query,
		habit.Name,
		habit.Description,
		habit.ID,
		habit.UserID,
	)
	return err
}

// DELETE
func (r *HabitRepository) Delete(habitID, userID string) error {
	query := `
		DELETE FROM habits
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(query, habitID, userID)
	return err
}

// UPDATE STREAK ON COMPLETION
func (r *HabitRepository) UpdateStreakOnCompletion(habitID, userID string, previousDate time.Time) (int, int, error) {
	query := `
		UPDATE habits
		SET current_streak = CASE
		        WHEN EXISTS (
		            SELECT 1 FROM habit_entries
		            WHERE habit_id = $1 AND entry_date = $2 AND status = 'completed'
		        )
		        THEN current_streak + 1
		        ELSE 1
		    END,
		    longest_streak = GREATEST(longest_streak, CASE
		        WHEN EXISTS (
		            SELECT 1 FROM habit_entries
		            WHERE habit_id = $1 AND entry_date = $2 AND status = 'completed'
		        )
		        THEN current_streak + 1
		        ELSE 1
		    END),
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $3
		RETURNING current_streak, longest_streak
	`
	var currentStreak int
	var longestStreak int
	err := r.db.QueryRow(query, habitID, previousDate, userID).Scan(&currentStreak, &longestStreak)
	if err != nil {
		return 0, 0, err
	}

	return currentStreak, longestStreak, nil
}

// RESET STREAK
func (r *HabitRepository) ResetStreak(habitID, userID string) error {
	query := `
		UPDATE habits
		SET current_streak = 0,
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(query, habitID, userID)
	return err
}

// ARCHIVE
func (r *HabitRepository) Archive(habitID, userID string) error {
	query := `
		UPDATE habits
		SET is_archived = TRUE, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(query, habitID, userID)
	return err
}

// UNARCHIVE
func (r *HabitRepository) Unarchive(habitID, userID string) error {
	query := `
		UPDATE habits
		SET is_archived = FALSE, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(query, habitID, userID)
	return err
}
