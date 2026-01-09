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

func (s *HabitEntryService) ListEntries(habitID string, startDate, endDate *time.Time) ([]models.HabitEntry, error) {
	return s.entryRepo.ListByHabit(habitID, startDate, endDate)
}

func (s *HabitEntryService) CreateEntry(userID, habitID string, date time.Time, status, note string) ([]models.Badge, error) {
	normalizedDate := startOfDay(date)

	exists, err := s.entryRepo.ExistsForDate(habitID, normalizedDate)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("entry already exists for date")
	}

	if err := s.entryRepo.CreateEntry(habitID, normalizedDate, status, note); err != nil {
		return nil, err
	}

	switch status {
	case "completed":
		previousDate := normalizedDate.AddDate(0, 0, -1)
		currentStreak, _, err := s.habitRepo.UpdateStreakOnCompletion(habitID, userID, previousDate)
		if err != nil {
			return nil, err
		}
		badges, err := s.badgeRepo.AwardBadgesForStreak(userID, habitID, currentStreak)
		if err != nil {
			return nil, err
		}
		return badges, nil
	case "missed":
		if err := s.habitRepo.ResetStreak(habitID, userID); err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, errors.New("invalid status")
	}
}

func (s *HabitEntryService) DeleteEntry(entryID, habitID string) error {
	return s.entryRepo.DeleteEntry(entryID, habitID)
}
