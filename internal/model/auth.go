package model

import (
	"github.com/labstack/echo/v4"
)

// AuthToken holds authentication token details with refresh token
// swagger:model
type AuthToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// AuthUser represents data stored in JWT token for user
type AuthUser struct {
	ID       int
	Username string
	Email    string
	Role     string
}

// Auth represents auth interface
type Auth interface {
	User(echo.Context) *AuthUser
}
