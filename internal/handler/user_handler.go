package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjubayerquanfinca/habit-tracker/internal/service"
)

// UserHandler handles user-related HTTP requests
// (implementations will be added in the next step)
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// PUT /api/v1/user
func (h *UserHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	if req.Email == "" && req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fields to update",
		})
		return
	}
	if err := h.userService.UpdateUser(userID, req.Email, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
	})
}

// DELETE /api/v1/user
func (h *UserHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	if err := h.userService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}

// GET /users/me
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}
