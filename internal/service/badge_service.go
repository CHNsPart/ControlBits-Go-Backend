package service

import (
	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
	"github.com/mjubayerquanfinca/habit-tracker/internal/repository"
)

type BadgeService struct {
	badgeRepo *repository.BadgeRepository
}

func NewBadgeService(badgeRepo *repository.BadgeRepository) *BadgeService {
	return &BadgeService{badgeRepo: badgeRepo}
}

func (s *BadgeService) ListAllBadges() ([]models.Badge, error) {
	return s.badgeRepo.ListAllBadges()
}

func (s *BadgeService) ListUserBadges(userID string) ([]models.EarnedBadge, error) {
	return s.badgeRepo.ListUserBadges(userID)
}

func (s *BadgeService) ListHabitBadges(userID, habitID string) ([]models.EarnedBadge, error) {
	return s.badgeRepo.ListHabitBadges(userID, habitID)
}
