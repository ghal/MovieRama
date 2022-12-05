package config_test

import (
	"movierama/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_New(t *testing.T) {
	os.Setenv("FRONTEND_URL", "http://localhost")
	os.Setenv("APP_JWT_SECRET", "secret_key")
	os.Setenv("APP_PORT", "appport")
	os.Setenv("MYSQL_USERNAME", "mysql_username")
	os.Setenv("MYSQL_PASSWORD", "mysql_password")
	os.Setenv("MYSQL_READ", "mysql_read")
	os.Setenv("MYSQL_WRITE", "mysql_write")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DB", "movierama_test")

	expCfg := &config.Config{
		App: config.App{
			FrontendURL: "http://localhost",
			JWTSecret:   "secret_key",
			Port:        "appport",
		},
		MySQL: config.MySQL{
			Username: "mysql_username",
			Password: "mysql_password",
			Read:     "mysql_read",
			Write:    "mysql_write",
			Port:     "3306",
			DB:       "movierama_test",
		},
	}

	cfg := config.New()
	assert.Equal(t, expCfg, cfg)
}
