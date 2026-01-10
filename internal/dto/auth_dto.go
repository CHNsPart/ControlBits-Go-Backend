package dto

// RegisterRequest is the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// RegisterResponse is the response for successful registration
// (could be just a message or user info)
type RegisterResponse struct {
	Message string `json:"message"`
}

// LoginRequest is the request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the response for successful login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}

// RefreshTokenRequest is the request body for refreshing tokens
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse is the response for a new access token
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// RequestPasswordResetRequest is the request body for password reset request
type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// RequestPasswordResetResponse is the response for password reset request
type RequestPasswordResetResponse struct {
	ResetToken string `json:"reset_token"`
	Message    string `json:"message"`
}

// ConfirmPasswordResetRequest is the request body for confirming password reset
type ConfirmPasswordResetRequest struct {
	ResetToken  string `json:"reset_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ConfirmPasswordResetResponse is the response for confirming password reset
type ConfirmPasswordResetResponse struct {
	Message string `json:"message"`
}

// ChangePasswordRequest is the request body for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePasswordResponse is the response for changing password
type ChangePasswordResponse struct {
	Message string `json:"message"`
}
