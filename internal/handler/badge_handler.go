package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjubayerquanfinca/habit-tracker/internal/service"
)

type BadgeHandler struct {
	service *service.BadgeService
}

func NewBadgeHandler(service *service.BadgeService) *BadgeHandler {
	return &BadgeHandler{service: service}
}

// ListAll godoc
// @Summary List all available badges
// @Tags Badges
// @Security BearerAuth
// @Produce json
// @Success 200 {array} BadgeResponse
// @Failure 500 {object} ErrorResponse
// @Router /badges [get]
func (h *BadgeHandler) ListAll(c *gin.Context) {
	badges, err := h.service.ListAllBadges()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, badges)
}

// ListUserBadges godoc
// @Summary List badges earned by current user
// @Tags Badges
// @Security BearerAuth
// @Produce json
// @Success 200 {array} EarnedBadgeResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me/badges [get]
func (h *BadgeHandler) ListUserBadges(c *gin.Context) {
	userID := c.GetString("userID")
	badges, err := h.service.ListUserBadges(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, badges)
}

// ListHabitBadges godoc
// @Summary List badges earned for a habit
// @Tags Badges
// @Security BearerAuth
// @Produce json
// @Param id path string true "Habit ID"
// @Success 200 {array} EarnedBadgeResponse
// @Failure 500 {object} ErrorResponse
// @Router /habits/{id}/badges [get]
func (h *BadgeHandler) ListHabitBadges(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")
	badges, err := h.service.ListHabitBadges(userID, habitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, badges)
}
