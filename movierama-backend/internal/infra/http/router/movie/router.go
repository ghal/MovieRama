package movie

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
	"movierama/internal/app/auth"
	"movierama/internal/app/movie"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Router infrastructure definition.
type Router struct {
	asSvc  movie.Service
	jwtCfg middleware.JWTConfig
}

// NewRouter returns an HTTP component to serve all the routes for the movies.
func NewRouter(asSvc movie.Service, jwtCfg middleware.JWTConfig) *Router {
	return &Router{
		asSvc:  asSvc,
		jwtCfg: jwtCfg,
	}
}

// AppendRoutes adds movies routes to router.
func (r *Router) AppendRoutes(e *echo.Echo) {
	e.GET("/movies", r.GetMoviesPublic)
	e.GET("/users/:user_id/movies", r.GetUserMoviesPublic)

	rg := e.Group("/api/v1")
	rg.Use(middleware.JWTWithConfig(r.jwtCfg))
	rg.GET("/movies", r.GetMovies)
	rg.GET("/users/:user_id/movies", r.GetUserMovies)
	rg.POST("/movies", r.CreateMovie)
	rg.POST("/movies/:movie_id/action/:action", r.MakeAction)
	rg.POST("/movies/:movie_id/remove_action/:action", r.RemoveAction)
}

// GetMoviesPublic gets the list of movies without auth.
func (r *Router) GetMoviesPublic(c echo.Context) error {
	sortType := c.QueryParam("sort")

	movies, err := r.asSvc.GetMoviesPublic(c.Request().Context(), sortType)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movies)
}

// GetMovies gets the list of movies.
func (r *Router) GetMovies(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), movie.AuthUserIDContextKey, GetAuthUserID(c))
	sortType := c.QueryParam("sort")

	movies, err := r.asSvc.GetMovies(ctx, sortType)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movies)
}

// GetUserMovies gets the user movies.
func (r *Router) GetUserMovies(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), movie.AuthUserIDContextKey, GetAuthUserID(c))
	sortType := c.QueryParam("sort")

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return err
	}
	movies, err := r.asSvc.GetUserMovies(ctx, userID, sortType)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movies)
}

// GetUserMoviesPublic gets the public user movies.
func (r *Router) GetUserMoviesPublic(c echo.Context) error {
	sortType := c.QueryParam("sort")

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return err
	}
	movies, err := r.asSvc.GetUserMoviesPublic(c.Request().Context(), userID, sortType)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movies)
}

// NewMovie contains the newly created movie payload struct.
type NewMovie struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// CreateMovie creates a movie.
func (r *Router) CreateMovie(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), movie.AuthUserIDContextKey, GetAuthUserID(c))

	m := new(NewMovie)
	err := c.Bind(m)
	if err != nil {
		return err
	}
	err = r.asSvc.CreateMovie(ctx, movie.NewMovie(*m))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, nil)
}

// MakeAction makes movie actions.
func (r *Router) MakeAction(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), movie.AuthUserIDContextKey, GetAuthUserID(c))

	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		return err
	}
	action := c.Param("action")

	err = r.asSvc.Action(ctx, movieID, action)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

// RemoveAction removes user movie action.
func (r *Router) RemoveAction(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), movie.AuthUserIDContextKey, GetAuthUserID(c))

	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		return err
	}
	action := c.Param("action")

	err = r.asSvc.RemoveAction(ctx, movieID, action)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func GetAuthUserID(c echo.Context) int {
	// Get user_id from token.
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	return claims.UserID
}
