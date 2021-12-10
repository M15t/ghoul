package user

import (
	"net/http"
	"strings"

	"github.com/M15t/ghoul/internal/model"
	"github.com/M15t/ghoul/pkg/server"
	dbutil "github.com/M15t/ghoul/pkg/util/db"
	httputil "github.com/M15t/ghoul/pkg/util/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents user http service
type HTTP struct {
	svc  Service
	auth model.Auth
}

// Service represents user application interface
type Service interface {
	Create(*model.AuthUser, CreationData) (*model.User, error)
	View(*model.AuthUser, int) (*model.User, error)
	List(*model.AuthUser, *dbutil.ListQueryCondition, *int64) ([]*model.User, error)
	Update(*model.AuthUser, int, UpdateData) (*model.User, error)
	Delete(*model.AuthUser, int) error
	Me(*model.AuthUser) (*model.User, error)
	ChangePassword(*model.AuthUser, PasswordChangeData) error
}

// NewHTTP creates new user http service
func NewHTTP(svc Service, auth model.Auth, eg *echo.Group) {
	h := HTTP{svc, auth}

	// swagger:operation POST /v1/users users usersCreate
	// ---
	// summary: Creates new user
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UserCreationData"
	// responses:
	//   "200":
	//     description: The new user
	//     schema:
	//       "$ref": "#/definitions/User"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.POST("", h.create)

	// swagger:operation GET /v1/users/{id} users usersView
	// ---
	// summary: Returns a single user
	// parameters:
	// - name: id
	//   in: path
	//   description: id of user
	//   type: integer
	//   required: true
	// responses:
	//   "200":
	//     description: The user
	//     schema:
	//       "$ref": "#/definitions/User"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.GET("/:id", h.view)

	// swagger:operation GET /v1/users users usersList
	// ---
	// summary: Returns list of users
	// responses:
	//   "200":
	//     description: List of users
	//     schema:
	//       "$ref": "#/definitions/UserListResp"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.GET("", h.list)

	// swagger:operation PATCH /v1/users/{id} users usersUpdate
	// ---
	// summary: Updates user information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of user
	//   type: integer
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UserUpdateData"
	// responses:
	//   "200":
	//     description: The updated user
	//     schema:
	//       "$ref": "#/definitions/User"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.PATCH("/:id", h.update)

	// swagger:operation DELETE /v1/users/{id} users usersDelete
	// ---
	// summary: Deletes an user
	// parameters:
	// - name: id
	//   in: path
	//   description: id of user
	//   type: integer
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.DELETE("/:id", h.delete)

	// swagger:operation GET /v1/users/me users usersMe
	// ---
	// summary: Returns authenticated user
	// responses:
	//   "200":
	//     description: Authenticated user
	//     schema:
	//       "$ref": "#/definitions/User"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.GET("/me", h.me)

	// swagger:operation PATCH /v1/users/me/password users usersChangePwd
	// ---
	// summary: Changes authenticated user password
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/PasswordChangeData"
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.PATCH("/me/password", h.changePassword)
}

// CreationData contains user data from json request
// swagger:model UserCreationData
type CreationData struct {
	Username  string `json:"username" validate:"required,min=3"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Mobile    string `json:"mobile" validate:"required,mobile"`
	Role      string `json:"role" validate:"required"`
	Blocked   bool   `json:"blocked"`
}

// UpdateData contains user data from json request
// swagger:model UserUpdateData
type UpdateData struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Mobile    *string `json:"mobile,omitempty" validate:"omitempty,mobile"`
	Role      *string `json:"role,omitempty"`
	Blocked   *bool   `json:"blocked,omitempty"`
}

// PasswordChangeData contains password change request
// swagger:model
type PasswordChangeData struct {
	OldPassword        string `json:"old_password" validate:"required"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required,eqfield=NewPassword"`
}

// ListResp contains list of users and current page number response
// swagger:model UserListResp
type ListResp struct {
	Data       []*model.User `json:"data"`
	TotalCount int64         `json:"total_count"`
}

func (h *HTTP) create(c echo.Context) error {
	r := CreationData{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	r.Email = strings.TrimSpace(r.Email)
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Mobile = strings.TrimSpace(strings.Replace(r.Mobile, " ", "", -1))
	r.Role = strings.TrimSpace(r.Role)

	if err := validateRole(&r.Role); err != nil {
		return err
	}

	resp, err := h.svc.Create(h.auth.User(c), r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) view(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	resp, err := h.svc.View(h.auth.User(c), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) list(c echo.Context) error {
	lq, err := httputil.ReqListQuery(c)
	if err != nil {
		return err
	}
	var count int64 = 0
	resp, err := h.svc.List(h.auth.User(c), lq, &count)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ListResp{resp, count})
}

func (h *HTTP) update(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	r := UpdateData{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	r.Email = httputil.TrimSpacePointer(r.Email)
	r.FirstName = httputil.TrimSpacePointer(r.FirstName)
	r.LastName = httputil.TrimSpacePointer(r.LastName)
	r.Mobile = httputil.RemoveSpacePointer(r.Mobile)
	r.Role = httputil.RemoveSpacePointer(r.Role)

	if err := validateRole(r.Role); err != nil {
		return err
	}

	resp, err := h.svc.Update(h.auth.User(c), id, r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(h.auth.User(c), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *HTTP) me(c echo.Context) error {
	resp, err := h.svc.Me(h.auth.User(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) changePassword(c echo.Context) error {
	r := PasswordChangeData{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	if err := h.svc.ChangePassword(h.auth.User(c), r); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func validateRole(input *string) error {
	if input == nil {
		return nil
	}

	validRole := false
	for _, role := range model.AvailableRoles {
		if role == *input {
			validRole = true
			break
		}
	}
	if !validRole {
		return server.NewHTTPValidationError("Invalid role")
	}
	return nil
}
