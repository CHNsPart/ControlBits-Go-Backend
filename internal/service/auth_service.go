package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/mjubayerquanfinca/habit-tracker/internal/models"
	"github.com/mjubayerquanfinca/habit-tracker/internal/repository"
	"github.com/mjubayerquanfinca/habit-tracker/pkg/utils"
)

type AuthService struct {
	userRepo      *repository.UserRepository
	jwtSecret     string
	refreshTokens map[string]string // userID -> refreshToken
	resetTokens   map[string]string // resetToken -> userID
}

// constructor
func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		refreshTokens: make(map[string]string),
		resetTokens:   make(map[string]string),
	}
}

// RequestPasswordReset generates a reset token and (simulates) sending it to the user's email
func (s *AuthService) RequestPasswordReset(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}
	// Generate a simple random token (use a secure method in production)
	resetToken := utils.GenerateUUID()
	s.resetTokens[resetToken] = user.ID
	// Simulate sending email by returning the token (in production, send via email)
	return resetToken, nil
}

// ConfirmPasswordReset validates the reset token and updates the user's password
func (s *AuthService) ConfirmPasswordReset(resetToken, newPassword string) error {
	userID, ok := s.resetTokens[resetToken]
	if !ok {
		return errors.New("invalid or expired reset token")
	}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}
	// Invalidate the token after use
	delete(s.resetTokens, resetToken)
	return nil
}

// REGISTER
func (s *AuthService) Register(email, password, name string) error {
	// password hash
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:           utils.GenerateUUID(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
	}

	return s.userRepo.Create(user)
}

// LOGIN
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) Login(email, password string) (*TokenPair, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// password match
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	// create refresh token (longer expiry)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		"type": "refresh",
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	// store refresh token in-memory (for demo; use DB in production)
	s.refreshTokens[user.ID] = refreshTokenString

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// ValidateRefreshToken validates a refresh token and returns a new access token if valid
func (s *AuthService) RefreshAccessToken(refreshTokenString string) (string, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return "", errors.New("invalid refresh token claims")
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid refresh token subject")
	}
	// check if refresh token matches stored one
	if s.refreshTokens[userID] != refreshTokenString {
		return "", errors.New("refresh token does not match")
	}
	// issue new access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return accessTokenString, nil
}

// UpdateUser updates user info (email, name, password)
func (s *AuthService) UpdateUser(userID, email, name string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if email != "" {
		user.Email = email
	}
	if name != "" {
		user.Name = name
	}

	return s.userRepo.Update(user)
}

// DeleteUser deletes a user by ID
func (s *AuthService) DeleteUser(userID string) error {
	return s.userRepo.Delete(userID)
}

// GetUserByID returns user by ID
func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

// ChangePassword changes the password for a user after verifying the old password
func (s *AuthService) ChangePassword(userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(user)
}

// Logout invalidates the refresh token for a user
func (s *AuthService) Logout(userID string) {
	delete(s.refreshTokens, userID)
}
