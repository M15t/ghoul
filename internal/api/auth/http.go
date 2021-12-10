package auth

import (
	"net/http"

	"github.com/M15t/ghoul/internal/model"

	"github.com/labstack/echo/v4"
)

// HTTP represents auth http service
type HTTP struct {
	svc Service
}

// Service represents auth service interface
type Service interface {
	Authenticate(echo.Context, Credentials) (*model.AuthToken, error)
	RefreshToken(echo.Context, RefreshTokenData) (*model.AuthToken, error)
}

// NewHTTP creates new auth http service
func NewHTTP(svc Service, e *echo.Echo) {
	h := HTTP{svc}

	// swagger:operation POST /login auth authLogin
	// ---
	// summary: Logs in user by username and password
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/Credentials"
	// responses:
	//   "200":
	//     description: Access token
	//     schema:
	//       "$ref": "#/definitions/AuthToken"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	e.POST("/login", h.login)

	// swagger:operation POST /refresh-token auth authRefreshToken
	// ---
	// summary: Refresh access token
	// parameters:
	// - name: token
	//   in: body
	//   description: The given `refresh_token` when login
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RefreshTokenData"
	// responses:
	//   "200":
	//     description: New access token
	//     schema:
	//       "$ref": "#/definitions/AuthToken"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	e.POST("/refresh-token", h.refreshToken)
}

// Credentials represents login request data
// swagger:model
type Credentials struct {
	// example: superadmin
	Username string `json:"username" validate:"required"`
	// example: superadmin123!@#
	Password string `json:"password" validate:"required"`
}

// RefreshTokenData represents refresh token request data
// swagger:model
type RefreshTokenData struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *HTTP) login(c echo.Context) error {
	r := Credentials{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	resp, err := h.svc.Authenticate(c, r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) refreshToken(c echo.Context) error {
	r := RefreshTokenData{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	resp, err := h.svc.RefreshToken(c, r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
