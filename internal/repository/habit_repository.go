package repository

import (
	"database/sql"

	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
)

type HabitRepository struct {
	db *sql.DB
}

func NewHabitRepository(db *sql.DB) *HabitRepository {
	return &HabitRepository{db: db}
}

// CREATE
func (r *HabitRepository) Create(habit *models.Habit) error {
	query := `
		INSERT INTO habits (id, user_id, name, description)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(
		query,
		habit.ID,
		habit.UserID,
		habit.Name,
		habit.Description,
	)
	return err
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

// INCREMENT STREAK
func (r *HabitRepository) IncrementStreak(habitID, userID string) error {
	query := `
		UPDATE habits
		SET current_streak = current_streak + 1,
		    longest_streak = GREATEST(longest_streak, current_streak + 1),
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(query, habitID, userID)
	return err
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
