package service

import (
	"errors"
	"time"

	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
	"github.com/mjubayerquanfinca/habit-tracker/internal/repository"
)

type HabitEntryService struct {
	entryRepo *repository.HabitEntryRepository
	habitRepo *repository.HabitRepository
	badgeRepo *repository.BadgeRepository
}

func NewHabitEntryService(
	entryRepo *repository.HabitEntryRepository,
	habitRepo *repository.HabitRepository,
	badgeRepo *repository.BadgeRepository,
) *HabitEntryService {
	return &HabitEntryService{
		entryRepo: entryRepo,
		habitRepo: habitRepo,
		badgeRepo: badgeRepo,
	}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// COMPLETE habit for today
func (s *HabitEntryService) CompleteHabit(userID, habitID string) ([]models.Badge, error) {
	today := startOfDay(time.Now())

	exists, err := s.entryRepo.ExistsForDate(habitID, today)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("habit already marked for today")
	}

	// create entry
	if err := s.entryRepo.Create(habitID, today, "completed"); err != nil {
		return nil, err
	}

	yesterday := today.AddDate(0, 0, -1)
	currentStreak, _, err := s.habitRepo.UpdateStreakOnCompletion(habitID, userID, yesterday)
	if err != nil {
		return nil, err
	}

	badges, err := s.badgeRepo.AwardBadgesForStreak(userID, habitID, currentStreak)
	if err != nil {
		return nil, err
	}

	return badges, nil
}

// MISS habit
func (s *HabitEntryService) MissHabit(userID, habitID string) error {
	today := startOfDay(time.Now())

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
