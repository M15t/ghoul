package swaggerutil

import (
	"ghoul/pkg/server"
)

// Success empty response
// swagger:response ok
type swaggOKResp struct{}

// Error empty response
// swagger:response err
type swaggErrResp struct{}

// Error response with details
// swagger:response errDetails
type swaggErrDetailsResp struct {
	//in: body
	Body server.ErrorResponse
}
