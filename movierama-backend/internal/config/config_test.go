package config_test

import (
	"movierama/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_New(t *testing.T) {
	os.Setenv("APP_USE_CACHE", "appusecache")
	os.Setenv("APP_PORT", "appport")
	os.Setenv("MYSQL_USERNAME", "mysql_username")
	os.Setenv("MYSQL_PASSWORD", "mysql_password")
	os.Setenv("MYSQL_READ", "mysql_read")
	os.Setenv("MYSQL_WRITE", "mysql_write")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DB", "movierama_test")
	os.Setenv("REDIS_PASSWORD", "redispass")
	os.Setenv("REDIS_HOST", "redishost")
	os.Setenv("REDIS_PORT", "redisport")
	os.Setenv("REDIS_DB", "redisdb")

	expCfg := &config.Config{
		App: config.App{
			WithCache: "appusecache",
			Port:      "appport",
		},
		MySQL: config.MySQL{
			Username: "mysql_username",
			Password: "mysql_password",
			Read:     "mysql_read",
			Write:    "mysql_write",
			Port:     "3306",
			DB:       "movierama_test",
		},
		Redis: config.Redis{
			Host:     "redishost",
			Port:     "redisport",
			Password: "redispass",
			Db:       "redisdb",
		},
	}

	cfg := config.New()
	assert.Equal(t, expCfg, cfg)
}
