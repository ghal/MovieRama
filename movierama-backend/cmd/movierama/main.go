package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"movierama/internal/app/auth"
	movieapp "movierama/internal/app/movie"
	"movierama/internal/config"
	authroute "movierama/internal/infra/http/router/auth"
	movieroute "movierama/internal/infra/http/router/movie"
	sqlrepo "movierama/internal/infra/repository/sql"
	"movierama/internal/infra/repository/sql/movie"
	"movierama/internal/infra/repository/sql/user"
)

func main() {
	e := echo.New()

	if err := run(e); err != nil {
		e.Logger.Fatalf("Error: %v")
	}
}

func run(e *echo.Echo) error {
	err := godotenv.Load("../../.env")
	if err != nil {
		e.Logger.Warnf("no .env file exists: %v", err)
	}

	// Setup app config
	cfg := config.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cfg.App.FrontendURL},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	reader, writer := setUpDB(cfg)
	defer func() {
		reader.Close()
		writer.Close()
	}()

	// Initialise repositories.
	mr, err := movie.NewRepository(reader, writer)
	if err != nil {
		return err
	}
	ur, err := user.NewRepository(reader, writer)
	if err != nil {
		return err
	}

	// Initialise Services.
	as := auth.NewService(ur, cfg)
	ms := movieapp.NewService(mr)

	// Configure middleware with the custom claims type
	jwtCfg := middleware.JWTConfig{
		Claims:     &auth.JwtCustomClaims{},
		SigningKey: []byte(cfg.App.JWTSecret),
	}
	// Initialise Routes.
	authroute.NewRouter(as).AppendRoutes(e)
	movieroute.NewRouter(ms, jwtCfg).AppendRoutes(e)

	// Run the app.
	return e.Start(":" + cfg.App.Port)
}

func setUpDB(cfg *config.Config) (read *sql.DB, write *sql.DB) {
	serviceStoreConfig := sqlrepo.DBConfig{
		Username: cfg.MySQL.Username,
		Pass:     cfg.MySQL.Password,
		Reader:   cfg.MySQL.Read,
		Writer:   cfg.MySQL.Write,
		Port:     cfg.MySQL.Port,
		DB:       cfg.MySQL.DB,
	}

	read, write = sqlrepo.NewDB(serviceStoreConfig)

	return read, write
}
