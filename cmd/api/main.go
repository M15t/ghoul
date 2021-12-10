// ghoul - Golang API Starter Kits
//
// API documents for ghoul.
//
// ## Authentication
// All API endpoints started with version, ex: `/v1/*`, require authentication token.
// Firstly, grab the **access_token** from the response of `/login`. Then include this header in all API calls:
// ```
// Authorization: Bearer ${access_token}
// ```
//
// For testing directly on this Swagger page, use the `Authorize` button right here bellow.
//
// Terms Of Service: N/A
//
//     Host: %{HOST}
//     Version: 1.0.0
//     Contact: M15t <khanhnguyen1411@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - login: []
//     - bearer: []
//
//     SecurityDefinitions:
//     login:
//         type: oauth2
//         tokenUrl: /login
//         refreshUrl: /refresh
//         flow: password
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"github.com/M15t/ghoul/config"
	"github.com/M15t/ghoul/internal/api/auth"
	"github.com/M15t/ghoul/internal/api/country"
	"github.com/M15t/ghoul/internal/api/user"
	"github.com/M15t/ghoul/internal/rbac"
	dbutil "github.com/M15t/ghoul/internal/util/db"
	_ "github.com/M15t/ghoul/internal/util/swagger" // Swagger stuffs
	"github.com/M15t/ghoul/pkg/server"
	"github.com/M15t/ghoul/pkg/server/middleware/jwt"
	"github.com/M15t/ghoul/pkg/util/crypter"
)

func main() {
	cfg, err := config.Load()
	checkErr(err)

	db, err := dbutil.New(cfg.DbPsn, cfg.DbLog)
	checkErr(err)
	// connection.Close() is not available for GORM 1.20.0
	// defer db.Close()

	// Initialize HTTP server
	e := server.New(&server.Config{
		Stage:        cfg.Stage,
		Port:         cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		AllowOrigins: cfg.AllowOrigins,
		Debug:        cfg.Debug,
	})

	// Static page for Swagger API specs
	e.Static("/swaggerui", "swaggerui")

	// Initialize DB interfaces
	userDB := user.NewDB()
	countryDB := country.NewDB()

	// Initialize services
	crypterSvc := crypter.New()
	rbacSvc := rbac.New(cfg.Debug)
	jwtSvc := jwt.New(cfg.JwtAlgorithm, cfg.JwtSecret, cfg.JwtDuration)
	authSvc := auth.New(db, userDB, jwtSvc, crypterSvc)
	userSvc := user.New(db, userDB, rbacSvc, crypterSvc)
	countrySvc := country.New(db, countryDB, rbacSvc)

	// Initialize root API
	auth.NewHTTP(authSvc, e)

	// Initialize v1 API
	v1Router := e.Group("/v1")
	v1Router.Use(jwtSvc.MWFunc())

	user.NewHTTP(userSvc, authSvc, v1Router.Group("/users"))
	country.NewHTTP(countrySvc, authSvc, v1Router.Group("/countries"))

	// Start the HTTP server
	server.Start(e, cfg.Stage == "development")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
