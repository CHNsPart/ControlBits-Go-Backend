package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mjubayerquanfinca/habit-tracker/internal/dto"
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

// GET /habits/:id/entries
func (h *HabitEntryHandler) ListEntries(c *gin.Context) {
	habitID := c.Param("id")
	var startDatePtr, endDatePtr *time.Time

	start := c.Query("start_date")
	end := c.Query("end_date")
	if start != "" {
		t, err := time.Parse("2006-01-02", start)
		if err == nil {
			startDatePtr = &t
		}
	}
	if end != "" {
		t, err := time.Parse("2006-01-02", end)
		if err == nil {
			endDatePtr = &t
		}
	}

	entries, err := h.service.ListEntries(habitID, startDatePtr, endDatePtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var resp dto.ListHabitEntriesResponse
	for _, e := range entries {
		resp.Entries = append(resp.Entries, dto.HabitEntryResponse{
			ID:        e.ID,
			HabitID:   e.HabitID,
			EntryDate: e.EntryDate.Format("2006-01-02"),
			Status:    e.Status,
			Note:      e.Note,
			CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	c.JSON(http.StatusOK, resp)
}

// POST /habits/:id/entries
func (h *HabitEntryHandler) CreateEntry(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")
	var req dto.CreateHabitEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" || req.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	status := strings.ToLower(req.Status)
	if status != "completed" && status != "missed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be completed or missed"})
		return
	}
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (use YYYY-MM-DD)"})
		return
	}
	entryID, badges, err := h.service.CreateEntry(userID, habitID, date, status, req.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var badgeIDs []string
	for _, b := range badges {
		badgeIDs = append(badgeIDs, b.ID)
	}
	c.JSON(http.StatusCreated, dto.CreateHabitEntryResponse{
		ID:        entryID,
		Message:   "entry created",
		NewBadges: badgeIDs,
	})
}

// DELETE /habits/:id/entries/:entryId
func (h *HabitEntryHandler) DeleteEntry(c *gin.Context) {
	habitID := c.Param("id")
	entryID := c.Param("entryId")
	if err := h.service.DeleteEntry(entryID, habitID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.DeleteHabitEntryResponse{Message: "entry deleted"})
}
