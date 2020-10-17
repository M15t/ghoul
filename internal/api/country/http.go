package country

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/M15t/ghoul/internal/model"
	"github.com/M15t/ghoul/pkg/server"
	dbutil "github.com/M15t/ghoul/pkg/util/db"
	httputil "github.com/M15t/ghoul/pkg/util/http"
	"github.com/labstack/echo/v4"
)

// HTTP represents country http service
type HTTP struct {
	svc  Service
	auth model.Auth
}

// Service represents country application interface
type Service interface {
	Create(*model.AuthUser, CreationData) (*model.Country, error)
	View(*model.AuthUser, int) (*model.Country, error)
	List(*model.AuthUser, *dbutil.ListQueryCondition, *int) ([]*model.Country, error)
	Update(*model.AuthUser, int, UpdateData) (*model.Country, error)
	Delete(*model.AuthUser, int) error
}

// NewHTTP creates new country http service
func NewHTTP(svc Service, auth model.Auth, eg *echo.Group) {
	h := HTTP{svc, auth}

	// swagger:operation POST /v1/countries countries countriesCreate
	// ---
	// summary: Creates new country
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CountryCreationData"
	// responses:
	//   "200":
	//     description: The new country
	//     schema:
	//       "$ref": "#/definitions/Country"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.POST("", h.create)

	// swagger:operation GET /v1/countries/{id} countries countriesView
	// ---
	// summary: Returns a single country
	// parameters:
	// - name: id
	//   in: path
	//   description: id of country
	//   type: integer
	//   required: true
	// responses:
	//   "200":
	//     description: The country
	//     schema:
	//       "$ref": "#/definitions/Country"
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

	// swagger:operation GET /v1/countries countries countriesList
	// ---
	// summary: Returns list of countries
	// responses:
	//   "200":
	//     description: List of countries
	//     schema:
	//       "$ref": "#/definitions/CountryListResp"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.GET("", h.list)

	// swagger:operation PATCH /v1/countries/{id} countries countriesUpdate
	// ---
	// summary: Updates country information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of country
	//   type: integer
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CountryUpdateData"
	// responses:
	//   "200":
	//     description: The updated country
	//     schema:
	//       "$ref": "#/definitions/Country"
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

	// swagger:operation DELETE /v1/countries/{id} countries countriesDelete
	// ---
	// summary: Deletes a country
	// parameters:
	// - name: id
	//   in: path
	//   description: id of country
	//   type: integer
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.DELETE("/:id", h.delete)
}

// CreationData contains country data from json request
// swagger:model CountryCreationData
type CreationData struct {
	// example: Vietnam
	Name string `json:"name" validate:"required,min=3"`
	// example: vn
	Code string `json:"code" validate:"required,min=2,max=10"`
	// example: +84
	PhoneCode string `json:"phone_code" validate:"required,min=2,max=10"`
}

// UpdateData contains country data from json request
// swagger:model CountryUpdateData
type UpdateData struct {
	// example: Vietnam
	Name *string `json:"name,omitempty" validate:"omitempty,min=3"`
	// example: vn
	Code *string `json:"code,omitempty" validate:"omitempty,min=2,max=10"`
	// example: +84
	PhoneCode *string `json:"phone_code,omitempty" validate:"omitempty,min=2,max=10"`
}

// ListResp contains list of paginated countries and total numbers of countries
// swagger:model CountryListResp
type ListResp struct {
	// example: [{"id": 1, "created_at": "2020-01-14T10:03:41Z", "updated_at": "2020-01-14T10:03:41Z", "name": "Singapore", "code": "SG", "phone_code": "+65"}]
	Data []*model.Country `json:"data"`
	// example: 1
	TotalCount int `json:"total_count"`
}

func (h *HTTP) create(c echo.Context) error {
	r := CreationData{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	r.Name = strings.TrimSpace(r.Name)
	r.Code = strings.ToUpper(strings.TrimSpace(r.Code))
	r.PhoneCode = strings.ReplaceAll(r.PhoneCode, " ", "")

	if regexp.MustCompile(`^\+\d+$`).Match([]byte(r.PhoneCode)) == false {
		return server.NewHTTPValidationError("PhoneCode is invalid")
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
	count := 0
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
	r.Name = httputil.TrimSpacePointer(r.Name)
	r.Code = httputil.TrimSpacePointer(r.Code)
	r.PhoneCode = httputil.RemoveSpacePointer(r.PhoneCode)

	usr, err := h.svc.Update(h.auth.User(c), id, r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
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
