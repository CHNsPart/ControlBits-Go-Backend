package service

import (
	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
	"github.com/mjubayerquanfinca/habit-tracker/internal/repository"
	"github.com/mjubayerquanfinca/habit-tracker/pkg/utils"
)

type HabitService struct {
	repo *repository.HabitRepository
}

func NewHabitService(repo *repository.HabitRepository) *HabitService {
	return &HabitService{repo: repo}
}

func (s *HabitService) Create(userID, name, description string) error {
	habit := &models.Habit{
		ID:          utils.GenerateUUID(),
		UserID:      userID,
		Name:        name,
		Description: description,
	}
	return s.repo.Create(habit)
}

func (s *HabitService) GetAll(userID string) ([]models.Habit, error) {
	return s.repo.FindByUser(userID)
}

func (s *HabitService) Update(userID, habitID, name, description string) error {
	habit := &models.Habit{
		ID:          habitID,
		UserID:      userID,
		Name:        name,
		Description: description,
	}
	return s.repo.Update(habit)
}

func (s *HabitService) GetByID(userID, habitID string) (*models.Habit, error) {
	return s.repo.GetByID(habitID, userID)
}

func (s *HabitService) Delete(userID, habitID string) error {
	return s.repo.Delete(habitID, userID)
}

func (s *HabitService) Archive(userID, habitID string) error {
	return s.repo.Archive(habitID, userID)
}

func (s *HabitService) Unarchive(userID, habitID string) error {
	return s.repo.Unarchive(habitID, userID)
}
