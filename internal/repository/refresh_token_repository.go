package repository

import (
	"database/sql"
	"time"
)

type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Save(userID, token string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, userID, token, expiresAt)
	return err
}

func (r *RefreshTokenRepository) GetByUserID(userID string) (*RefreshToken, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`
	row := r.db.QueryRow(query, userID)
	var rt RefreshToken
	err := row.Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) DeleteByUserID(userID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
