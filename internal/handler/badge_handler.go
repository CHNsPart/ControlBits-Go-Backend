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

// GET /badges
func (h *BadgeHandler) ListAll(c *gin.Context) {
	badges, err := h.service.ListAllBadges()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, badges)
}

// GET /users/me/badges
func (h *BadgeHandler) ListUserBadges(c *gin.Context) {
	userID := c.GetString("userID")
	badges, err := h.service.ListUserBadges(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, badges)
}

// GET /habits/:id/badges
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
