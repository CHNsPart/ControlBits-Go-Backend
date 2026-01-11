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

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
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

// Login godoc
// @Summary Login and receive access/refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
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

// RequestPasswordReset godoc
// @Summary Request a password reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RequestPasswordResetRequest true "Password reset request"
// @Success 200 {object} dto.RequestPasswordResetResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/forgot-password [post]
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

// ConfirmPasswordReset godoc
// @Summary Confirm password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ConfirmPasswordResetRequest true "Confirm password reset"
// @Success 200 {object} dto.ConfirmPasswordResetResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/reset-password [post]
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

// RefreshToken godoc
// @Summary Refresh access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
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

// ChangePassword godoc
// @Summary Change current user password
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Change password request"
// @Success 200 {object} dto.ChangePasswordResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/me/password [put]
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

// Logout godoc
// @Summary Logout and invalidate refresh token
// @Tags Auth
// @Produce json
// @Success 200 {object} LogoutResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("userID")
	h.authService.Logout(userID)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
