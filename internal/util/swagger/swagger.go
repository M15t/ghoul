package swagger

import (
	httputil "github.com/M15t/ghoul/pkg/util/http"
	_ "github.com/M15t/ghoul/pkg/util/swagger" // Swagger stuffs
)

// ListRequest holds data of listing request from react-admin
// swagger:parameters usersList countriesList
type ListRequest struct {
	httputil.ListRequest
}
