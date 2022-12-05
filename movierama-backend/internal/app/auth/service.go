package auth

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"movierama/internal/config"
	"movierama/internal/infra/repository/sql/user"
	"time"
)

const (
	accessTokenExpirationHours = 72
)

// JwtCustomClaims contains the custom jwt claims data.
type JwtCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// Service handles auth information.
type Service interface {
	Login(c echo.Context, user *UserLogin) (*LoginRes, error)
	Register(ctx context.Context, registerUser *UserRegister) error
}

type authService struct {
	ur  Repository
	cfg *config.Config
}

// NewService constructor.
func NewService(userRepo Repository, cfg *config.Config) Service {
	return &authService{
		ur:  userRepo,
		cfg: cfg,
	}
}

// UserLogin contains the payload struct for login.
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a authService) Login(c echo.Context, user *UserLogin) (*LoginRes, error) {
	dbUser, err := a.ur.GetUserAuthDetails(c.Request().Context(), user.Username)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		// If the two passwords don't match, return an error.
		return nil, err
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaims{
		UserID: dbUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * accessTokenExpirationHours).Unix(),
		},
	})

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(a.cfg.App.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &LoginRes{Username: dbUser.Username, Token: t}, nil
}

// UserRegister contains the payload struct for registration.
type UserRegister struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (a authService) Register(ctx context.Context, registerUser *UserRegister) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	u := &user.SQLUser{
		Username:  registerUser.Username,
		Password:  string(passwordHash),
		FirstName: registerUser.FirstName,
		LastName:  registerUser.LastName,
	}

	err = a.ur.CreateUser(ctx, u)
	if err != nil {
		return err
	}

	return nil
}
