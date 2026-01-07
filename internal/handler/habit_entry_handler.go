package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjubayerquanfinca/habit-tracker/internal/service"
)

type HabitEntryHandler struct {
	service *service.HabitEntryService
}

func NewHabitEntryHandler(service *service.HabitEntryService) *HabitEntryHandler {
	return &HabitEntryHandler{service: service}
}

// POST /habits/:id/complete
func (h *HabitEntryHandler) Complete(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")

	badges, err := h.service.CompleteHabit(userID, habitID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "habit marked as completed",
		"new_badges": badges,
	})
}

// POST /habits/:id/miss
func (h *HabitEntryHandler) Miss(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")

	if err := h.service.MissHabit(userID, habitID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "habit marked as missed"})
}
