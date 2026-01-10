package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mjubayerquanfinca/habit-tracker/internal/dto"
	"github.com/mjubayerquanfinca/habit-tracker/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

// constructor
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// REGISTER API
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	err := h.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.RegisterResponse{Message: "user registered successfully"})
}

// LOGIN API
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	tokenPair, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		Message:      "login successful",
	})
}

// PASSWORD RESET REQUEST API
// POST /api/v1/auth/request-password-reset
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req dto.RequestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	token, err := h.authService.RequestPasswordReset(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.RequestPasswordResetResponse{
		ResetToken: token,
		Message:    "Password reset token generated. (Simulated email)",
	})
}

// PASSWORD RESET CONFIRM API
// POST /api/v1/auth/confirm-password-reset
func (h *AuthHandler) ConfirmPasswordReset(c *gin.Context) {
	var req dto.ConfirmPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.authService.ConfirmPasswordReset(req.ResetToken, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.ConfirmPasswordResetResponse{Message: "Password has been reset successfully."})
}

// REFRESH TOKEN API
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	accessToken, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.RefreshTokenResponse{AccessToken: accessToken})
}

// PUT /users/me/password - Change password for current user
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("userID")
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.OldPassword == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "old and new password required"})
		return
	}
	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.ChangePasswordResponse{Message: "password changed successfully"})
}

// POST /auth/logout - Invalidate refresh token
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("userID")
	h.authService.Logout(userID)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
