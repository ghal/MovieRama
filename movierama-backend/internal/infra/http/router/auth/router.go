package auth

import (
	"github.com/labstack/echo/v4"
	"movierama/internal/app/auth"
	"net/http"
)

// Router infrastructure definition.
type Router struct {
	asSvc auth.Service
}

// NewRouter returns an HTTP component to serve all the routes for the user auth.
func NewRouter(asSvc auth.Service) *Router {
	return &Router{
		asSvc: asSvc,
	}
}

// AppendRoutes adds auth routes to router.
func (r *Router) AppendRoutes(e *echo.Echo) {
	e.POST("/api/v1/auth/login", r.Login)
	e.POST("/api/v1/auth/register", r.Register)
}

// Login log's in the user.
func (r *Router) Login(c echo.Context) error {
	ul := new(auth.UserLogin)
	err := c.Bind(ul)
	if err != nil {
		return err
	}

	res, err := r.asSvc.Login(c, ul)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// Register creates the user.
func (r *Router) Register(c echo.Context) error {
	userInput := new(auth.UserRegister)
	err := c.Bind(userInput)
	if err != nil {
		return err
	}
	err = r.asSvc.Register(c.Request().Context(), userInput)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, nil)
}
