package config

import (
	"os"
)

// Config is the main configuration of app.
type Config struct {
	App   App
	MySQL MySQL
	Redis Redis
}

// App contains app configuration.
type App struct {
	FrontendURL string
	JWTSecret   string
	WithCache   string
	Port        string
}

// MySQL contains mysql configuration.
type MySQL struct {
	Username string
	Password string
	Read     string
	Write    string
	Port     string
	DB       string
}

// Redis contains redis configuration.
type Redis struct {
	Host     string
	Port     string
	Password string
	Db       string
}

// New constructor
func New() *Config {
	cfg := &Config{}
	cfg.setAppConfig()
	cfg.setRedisConfig()
	cfg.setMySQLConfig()

	return cfg
}

// SetAppConfig creates an App struct.
func (cfg *Config) setAppConfig() {
	cfg.App = App{
		FrontendURL: os.Getenv("FRONTEND_URL"),
		JWTSecret:   os.Getenv("APP_JWT_SECRET"),
		WithCache:   os.Getenv("APP_USE_CACHE"),
		Port:        os.Getenv("APP_PORT"),
	}
}

// SetMySQLConfig creates a MySQL config struct.
func (cfg *Config) setMySQLConfig() {
	cfg.MySQL = MySQL{
		Username: os.Getenv("MYSQL_USERNAME"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Read:     os.Getenv("MYSQL_READ"),
		Write:    os.Getenv("MYSQL_WRITE"),
		Port:     os.Getenv("MYSQL_PORT"),
		DB:       os.Getenv("MYSQL_DB"),
	}
}

// SetRedisConfig creates a Redis config struct.
func (cfg *Config) setRedisConfig() {
	cfg.Redis = Redis{
		Password: os.Getenv("REDIS_PASSWORD"),
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Db:       os.Getenv("REDIS_DB"),
	}
}
