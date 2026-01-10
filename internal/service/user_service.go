package service

import (
	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
	"github.com/mjubayerquanfinca/habit-tracker/internal/repository"
)

// UserService handles business logic for users
// (implementations will be added in the next step)
type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetUserByID returns user by ID
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

// UpdateUser updates user info (email, name)
func (s *UserService) UpdateUser(userID, email, name string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if email != "" {
		user.Email = email
	}
	if name != "" {
		user.Name = name
	}
	return s.userRepo.Update(user)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(userID string) error {
	return s.userRepo.Delete(userID)
}
