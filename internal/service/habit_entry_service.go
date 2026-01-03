package service

import (
	"errors"
	"time"

	"github.com/mjubayerquanfinca/habit-tracker/internal/repository"
)

type HabitEntryService struct {
	entryRepo *repository.HabitEntryRepository
	habitRepo *repository.HabitRepository
}

func NewHabitEntryService(
	entryRepo *repository.HabitEntryRepository,
	habitRepo *repository.HabitRepository,
) *HabitEntryService {
	return &HabitEntryService{
		entryRepo: entryRepo,
		habitRepo: habitRepo,
	}
}

// COMPLETE habit for today
func (s *HabitEntryService) CompleteHabit(userID, habitID string) error {
	today := time.Now().Truncate(24 * time.Hour)

	exists, err := s.entryRepo.ExistsForDate(habitID, today)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("habit already marked for today")
	}

	// create entry
	if err := s.entryRepo.Create(habitID, today, "completed"); err != nil {
		return err
	}

	// increase streak
	return s.habitRepo.IncrementStreak(habitID, userID)
}

// MISS habit
func (s *HabitEntryService) MissHabit(userID, habitID string) error {
	today := time.Now().Truncate(24 * time.Hour)

	exists, err := s.entryRepo.ExistsForDate(habitID, today)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("habit already marked for today")
	}

	// create entry
	if err := s.entryRepo.Create(habitID, today, "missed"); err != nil {
		return err
	}

	// reset streak
	return s.habitRepo.ResetStreak(habitID, userID)
}
