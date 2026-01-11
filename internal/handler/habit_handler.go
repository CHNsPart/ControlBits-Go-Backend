package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjubayerquanfinca/habit-tracker/internal/dto"
	"github.com/mjubayerquanfinca/habit-tracker/internal/service"
)

type HabitHandler struct {
	service *service.HabitService
}

func NewHabitHandler(service *service.HabitService) *HabitHandler {
	return &HabitHandler{service: service}
}

// Create godoc
// @Summary Create a new habit
// @Tags Habits
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateHabitRequest true "Create habit request"
// @Success 201 {object} dto.CreateHabitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /habits [post]
func (h *HabitHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.CreateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	id, err := h.service.Create(userID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.CreateHabitResponse{ID: id, Message: "habit created"})
}

// GetAll godoc
// @Summary List habits for current user
// @Tags Habits
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.ListHabitsResponse
// @Failure 500 {object} ErrorResponse
// @Router /habits [get]
func (h *HabitHandler) GetAll(c *gin.Context) {
	userID := c.GetString("userID")

	habits, err := h.service.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var resp dto.ListHabitsResponse
	for _, h := range habits {
		resp.Habits = append(resp.Habits, dto.GetHabitResponse{
			ID:            h.ID,
			Name:          h.Name,
			Description:   h.Description,
			CurrentStreak: h.CurrentStreak,
			LongestStreak: h.LongestStreak,
			CreatedAt:     h.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     h.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Archived:      h.IsArchived,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// GetByID godoc
// @Summary Get a habit by ID
// @Tags Habits
// @Security BearerAuth
// @Produce json
// @Param id path string true "Habit ID"
// @Success 200 {object} dto.GetHabitResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /habits/{id} [get]
func (h *HabitHandler) GetByID(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")

	habit, err := h.service.GetByID(userID, habitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if habit == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "habit not found"})
		return
	}
	resp := dto.GetHabitResponse{
		ID:            habit.ID,
		Name:          habit.Name,
		Description:   habit.Description,
		CurrentStreak: habit.CurrentStreak,
		LongestStreak: habit.LongestStreak,
		CreatedAt:     habit.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     habit.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Archived:      habit.IsArchived,
	}
	c.JSON(http.StatusOK, resp)
}

// Update godoc
// @Summary Update a habit by ID
// @Tags Habits
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Habit ID"
// @Param request body dto.UpdateHabitRequest true "Update habit request"
// @Success 200 {object} dto.UpdateHabitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /habits/{id} [put]
func (h *HabitHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")

	var req dto.UpdateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.service.Update(userID, habitID, req.Name, req.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.UpdateHabitResponse{Message: "habit updated"})
}

// Delete godoc
// @Summary Delete a habit by ID
// @Tags Habits
// @Security BearerAuth
// @Produce json
// @Param id path string true "Habit ID"
// @Success 200 {object} dto.DeleteHabitResponse
// @Failure 500 {object} ErrorResponse
// @Router /habits/{id} [delete]
func (h *HabitHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")

	if err := h.service.Delete(userID, habitID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.DeleteHabitResponse{Message: "habit deleted"})
}

// Archive godoc
// @Summary Archive a habit
// @Tags Habits
// @Security BearerAuth
// @Produce json
// @Param id path string true "Habit ID"
// @Success 200 {object} dto.ArchiveHabitResponse
// @Failure 500 {object} ErrorResponse
// @Router /habits/{id}/archive [post]
func (h *HabitHandler) Archive(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")
	if err := h.service.Archive(userID, habitID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.ArchiveHabitResponse{Message: "habit archived"})
}

// Unarchive godoc
// @Summary Unarchive a habit
// @Tags Habits
// @Security BearerAuth
// @Produce json
// @Param id path string true "Habit ID"
// @Success 200 {object} dto.ArchiveHabitResponse
// @Failure 500 {object} ErrorResponse
// @Router /habits/{id}/unarchive [post]
func (h *HabitHandler) Unarchive(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")
	if err := h.service.Unarchive(userID, habitID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.ArchiveHabitResponse{Message: "habit unarchived"})
}
