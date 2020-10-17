package jwt_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/M15t/ghoul/pkg/mock"
	"github.com/M15t/ghoul/pkg/server"
	"github.com/M15t/ghoul/pkg/server/middleware/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func echoHandler(mw ...echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = server.NewErrorHandler(e).Handle
	for _, v := range mw {
		e.Use(v)
	}
	e.GET("/hello", hwHandler)
	return e
}

func hwHandler(c echo.Context) error {
	return c.String(200, "Hello World")
}

func TestMWFunc(t *testing.T) {
	cases := []struct {
		name       string
		wantStatus int
		header     string
		signMethod string
	}{
		{
			name:       "Empty header",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Header not containing Bearer",
			header:     "notBearer",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid header",
			header:     mock.HeaderInvalid(),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Success",
			header:     mock.HeaderValid(),
			wantStatus: http.StatusOK,
		},
	}
	jwtMW := jwt.New("HS256", "jwtsecret", 60)
	ts := httptest.NewServer(echoHandler(jwtMW.MWFunc()))
	defer ts.Close()
	path := ts.URL + "/hello"
	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", path, nil)
			req.Header.Set("Authorization", tt.header)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal("Cannot create http request")
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	testExp := time.Date(2019, 4, 24, 0, 0, 0, 0, time.Local)
	type args struct {
		claims map[string]interface{}
		expire *time.Time
	}
	tests := []struct {
		name    string
		algo    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Invalid algo",
			algo: "invalid",
		},
		{
			name: "Success with expire",
			algo: "HS256",
			args: args{
				claims: map[string]interface{}{"id": 1, "username": "superadmin", "role": 1},
				expire: &testExp,
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTYwMzg4MDAsImlkIjoxLCJyb2xlIjoxLCJ1c2VybmFtZSI6InN1cGVyYWRtaW4ifQ.jFXLTxxiV4Bs7cenIGxfcPkwJvVSDhlUz78qf4l_IqE",
			wantErr: false,
		},
		{
			name: "Success without expire",
			algo: "HS256",
			args: args{
				claims: map[string]interface{}{"id": 1, "username": "superadmin", "role": 1},
				expire: nil,
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.algo != "HS256" {
				assert.Panics(t, func() {
					jwt.New(tt.algo, "jwtsecret", 60)
				}, "The code did not panic")
				return
			}
			j := jwt.New(tt.algo, "jwtsecret", 60)
			got, _, err := j.GenerateToken(tt.args.claims, tt.args.expire)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				if tt.args.expire == nil && strings.Split(got, ".")[0] == tt.want {
					return
				}
				t.Errorf("GenerateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
