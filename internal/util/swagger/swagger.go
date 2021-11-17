package swagger

import (
	httputil "ghoul/pkg/util/http"
	_ "ghoul/pkg/util/swagger" // Swagger stuffs
)

// ListRequest holds data of listing request from react-admin
// swagger:parameters usersList countriesList
type ListRequest struct {
	httputil.ListRequest
}
