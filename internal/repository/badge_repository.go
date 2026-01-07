package repository

import (
	"database/sql"

	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
)

type BadgeRepository struct {
	db *sql.DB
}

func NewBadgeRepository(db *sql.DB) *BadgeRepository {
	return &BadgeRepository{db: db}
}

// AwardBadgesForStreak inserts newly earned badges and returns them.
func (r *BadgeRepository) AwardBadgesForStreak(userID, habitID string, streak int) ([]models.Badge, error) {
	if streak <= 0 {
		return nil, nil
	}

	query := `
		WITH inserted AS (
			INSERT INTO user_badges (user_id, habit_id, badge_id)
			SELECT $1, $2, b.id
			FROM badges b
			WHERE b.streak_required <= $3
			ON CONFLICT (user_id, habit_id, badge_id) DO NOTHING
			RETURNING badge_id
		)
		SELECT b.id, b.name, b.description, b.streak_required, b.icon
		FROM badges b
		JOIN inserted i ON i.badge_id = b.id
		ORDER BY b.streak_required
	`

	rows, err := r.db.Query(query, userID, habitID, streak)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	badges := make([]models.Badge, 0)
	for rows.Next() {
		var badge models.Badge
		if err := rows.Scan(
			&badge.ID,
			&badge.Name,
			&badge.Description,
			&badge.StreakRequired,
			&badge.Icon,
		); err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return badges, nil
}
