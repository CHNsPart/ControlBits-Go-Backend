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

func (r *BadgeRepository) ListAllBadges() ([]models.Badge, error) {
	query := `
		SELECT id, name, description, streak_required, icon
		FROM badges
		ORDER BY streak_required
	`

	rows, err := r.db.Query(query)
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

func (r *BadgeRepository) ListUserBadges(userID string) ([]models.EarnedBadge, error) {
	query := `
		SELECT b.id, b.name, b.description, b.streak_required, b.icon,
		       ub.earned_at, ub.habit_id
		FROM user_badges ub
		JOIN badges b ON b.id = ub.badge_id
		WHERE ub.user_id = $1
		ORDER BY ub.earned_at DESC, b.streak_required
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	badges := make([]models.EarnedBadge, 0)
	for rows.Next() {
		var badge models.EarnedBadge
		if err := rows.Scan(
			&badge.ID,
			&badge.Name,
			&badge.Description,
			&badge.StreakRequired,
			&badge.Icon,
			&badge.EarnedAt,
			&badge.HabitID,
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

func (r *BadgeRepository) ListHabitBadges(userID, habitID string) ([]models.EarnedBadge, error) {
	query := `
		SELECT b.id, b.name, b.description, b.streak_required, b.icon,
		       ub.earned_at, ub.habit_id
		FROM user_badges ub
		JOIN badges b ON b.id = ub.badge_id
		WHERE ub.user_id = $1 AND ub.habit_id = $2
		ORDER BY ub.earned_at DESC, b.streak_required
	`

	rows, err := r.db.Query(query, userID, habitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	badges := make([]models.EarnedBadge, 0)
	for rows.Next() {
		var badge models.EarnedBadge
		if err := rows.Scan(
			&badge.ID,
			&badge.Name,
			&badge.Description,
			&badge.StreakRequired,
			&badge.Icon,
			&badge.EarnedAt,
			&badge.HabitID,
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
