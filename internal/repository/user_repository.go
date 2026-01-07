package repository

import (
	"database/sql"
	"errors"

	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

// constructor
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create user (Register)
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
	)

	return err
}

// Find user by email (Login)
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

	// Find user by ID
	func (r *UserRepository) FindByID(userID string) (*models.User, error) {
		query := `
			SELECT id, email, password_hash, name, created_at, updated_at
			FROM users
			WHERE id = $1
		`

		var user models.User
		err := r.db.QueryRow(query, userID).Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		if err != nil {
			return nil, err
		}
		return &user, nil
	}

// Update user by ID
func (r *UserRepository) Update(user *models.User) error {
	query := `
			UPDATE users
			SET email = $1, password_hash = $2, name = $3, updated_at = NOW()
			WHERE id = $4
		`

	_, err := r.db.Exec(
		query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.ID,
	)
	return err
}

// Delete user by ID
func (r *UserRepository) Delete(userID string) error {
	query := `
			DELETE FROM users WHERE id = $1
		`
	_, err := r.db.Exec(query, userID)
	return err
}
