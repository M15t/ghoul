package auth

import (
	"net/http"
	"time"

	"github.com/M15t/ghoul/internal/model"
	"github.com/M15t/ghoul/pkg/server"
	"github.com/labstack/echo/v4"
)

// Custom errors
var (
	ErrInvalidCredentials  = server.NewHTTPError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Username or password is incorrect")
	ErrUserBlocked         = server.NewHTTPError(http.StatusUnauthorized, "USER_BLOCKED", "Your account has been blocked and may not login")
	ErrInvalidRefreshToken = server.NewHTTPError(http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Invalid refresh token")
)

// LoginUser logs in the given user, returns access token
func (s *Auth) LoginUser(u *model.User) (*model.AuthToken, error) {
	claims := map[string]interface{}{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"role":     u.Role,
	}
	token, expiresin, err := s.jwt.GenerateToken(claims, nil)
	if err != nil {
		return nil, server.NewHTTPInternalError("Error generating token").SetInternal(err)
	}

	refreshToken := s.cr.UID()
	err = s.udb.Update(s.db, map[string]interface{}{"refresh_token": refreshToken, "last_login": time.Now()}, u.ID)
	if err != nil {
		return nil, server.NewHTTPInternalError("Error updating user").SetInternal(err)
	}

	return &model.AuthToken{AccessToken: token, TokenType: "bearer", ExpiresIn: expiresin, RefreshToken: refreshToken}, nil
}

// Authenticate tries to authenticate the user provided by given credentials
func (s *Auth) Authenticate(c echo.Context, data Credentials) (*model.AuthToken, error) {
	usr, err := s.udb.FindByUsername(s.db, data.Username)
	if err != nil || usr == nil {
		return nil, ErrInvalidCredentials.SetInternal(err)
	}
	if !s.cr.CompareHashAndPassword(usr.Password, data.Password) {
		return nil, ErrInvalidCredentials
	}
	if usr.Blocked {
		return nil, ErrUserBlocked
	}

	return s.LoginUser(usr)
}

// RefreshToken returns the new access token with expired time extended
func (s *Auth) RefreshToken(c echo.Context, data RefreshTokenData) (*model.AuthToken, error) {
	usr, err := s.udb.FindByRefreshToken(s.db, data.RefreshToken)
	if err != nil || usr == nil {
		return nil, ErrInvalidRefreshToken.SetInternal(err)
	}
	return s.LoginUser(usr)
}

// User returns user data stored in jwt token
func (s *Auth) User(c echo.Context) *model.AuthUser {
	id, _ := c.Get("id").(float64)
	user, _ := c.Get("username").(string)
	email, _ := c.Get("email").(string)
	role, _ := c.Get("role").(string)
	return &model.AuthUser{
		ID:       int(id),
		Username: user,
		Email:    email,
		Role:     role,
	}
}
