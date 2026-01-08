package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjubayerquanfinca/habit-tracker/internal/service"
)

type HabitHandler struct {
	service *service.HabitService
}

func NewHabitHandler(service *service.HabitService) *HabitHandler {
	return &HabitHandler{service: service}
}

// POST /habits
func (h *HabitHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.Create(userID, req.Name, req.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "habit created"})
}

// GET /habits
func (h *HabitHandler) GetAll(c *gin.Context) {
	userID := c.GetString("userID")

	habits, err := h.service.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, habits)
}

// GET /habits/:id
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
	c.JSON(http.StatusOK, habit)
}

// PUT /habits/:id
func (h *HabitHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.Update(userID, habitID, req.Name, req.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "habit updated"})
}

// DELETE /habits/:id
func (h *HabitHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")

	if err := h.service.Delete(userID, habitID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "habit deleted"})
}

// POST /habits/:id/archive
func (h *HabitHandler) Archive(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")
	if err := h.service.Archive(userID, habitID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "habit archived"})
}

// POST /habits/:id/unarchive
func (h *HabitHandler) Unarchive(c *gin.Context) {
	userID := c.GetString("userID")
	habitID := c.Param("id")
	if err := h.service.Unarchive(userID, habitID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "habit unarchived"})
}
