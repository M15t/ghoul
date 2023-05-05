package config_test

import (
	"os"
	"testing"

	"github.com/M15t/ghoul/config"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	type args struct {
		customize func()
	}
	cases := []struct {
		name     string
		args     args
		wantData *config.Configuration
		wantErr  bool
	}{
		{
			name: "Success",
			wantData: &config.Configuration{
				Stage:        "test",
				Host:         "localhost",
				Port:         8080,
				ReadTimeout:  10,
				WriteTimeout: 5,
				Debug:        true,
				DbLog:        true,
				DbDsn:        "postgres://ghoul:ghoul123@localhost:5432/ghoul?sslmode=disable&connect_timeout=5",
				JwtSecret:    "jwtsecret",
				JwtDuration:  31536000,
				JwtAlgorithm: "HS256",
			},
		},
		{
			name: "Failure",
			args: args{
				customize: func() {
					os.Setenv("PORT", "asdasd")
				},
			},
			wantData: nil,
			wantErr:  true,
		},
	}

	os.Setenv("UP_STAGE", "test")

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.customize != nil {
				tt.args.customize()
			}
			cfg, err := config.Load()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantData, cfg)
		})
	}
}
