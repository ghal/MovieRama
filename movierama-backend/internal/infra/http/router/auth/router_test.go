package auth_test

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authservice "movierama/internal/app/auth"
	"movierama/internal/infra/http/router/auth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_Login(t *testing.T) {
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc   *authservice.SvcMock
		jwtConfig middleware.JWTConfig
		echoCtx   echo.Context
		expRes    string
		expErr    error
	}{
		"Should succeed on Login call": {
			mockSvc: func() *authservice.SvcMock {
				mockSvc := &authservice.SvcMock{}
				ul := &authservice.UserLogin{
					Username: "test_username",
					Password: "secret_password",
				}
				mockSvc.On("Login", mock.Anything, ul).Return(&authservice.LoginRes{
					Username: "test_username",
					Token:    "test_token",
				}, nil)

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				body, _ := json.Marshal(&authservice.UserLogin{
					Username: "test_username",
					Password: "secret_password",
				})
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath("/api/v1/auth/login")

				return c
			}(),
			jwtConfig: middleware.JWTConfig{},
			expRes: func() string {
				b, err := json.Marshal(&authservice.LoginRes{
					Username: "test_username",
					Token:    "test_token",
				})
				assert.Nil(t, err)

				return string(b) + "\n"
			}(),
		},
		"Should return error on Login error ": {
			mockSvc: func() *authservice.SvcMock {
				mockSvc := &authservice.SvcMock{}
				ul := &authservice.UserLogin{
					Username: "test_username",
					Password: "secret_password",
				}
				mockSvc.On("Login", mock.Anything, ul).
					Return(&authservice.LoginRes{}, errors.New("random error"))

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				body, _ := json.Marshal(&authservice.UserLogin{
					Username: "test_username",
					Password: "secret_password",
				})
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath("/api/v1/auth/login")

				return c
			}(),
			expErr: errors.New("random error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			r := auth.NewRouter(tt.mockSvc)
			err := r.Login(tt.echoCtx)

			if tt.expErr == nil {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, tt.expRes, rec.Body.String())
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expErr, err)
			}
		})
	}
}

func TestRouter_Register(t *testing.T) {
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc   *authservice.SvcMock
		jwtConfig middleware.JWTConfig
		echoCtx   echo.Context
		expRes    string
		expErr    error
	}{
		"Should succeed on Register call": {
			mockSvc: func() *authservice.SvcMock {
				mockSvc := &authservice.SvcMock{}
				ur := &authservice.UserRegister{
					FirstName: "test_firstname",
					LastName:  "test_lastname",
					Username:  "test_username",
					Password:  "secret_password",
				}
				mockSvc.On("Register", mock.Anything, ur).Return(nil)

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				body, _ := json.Marshal(&authservice.UserRegister{
					FirstName: "test_firstname",
					LastName:  "test_lastname",
					Username:  "test_username",
					Password:  "secret_password",
				})
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath("/api/v1/auth/register")

				return c
			}(),
			jwtConfig: middleware.JWTConfig{},
		},
		"Should return error on Register error ": {
			mockSvc: func() *authservice.SvcMock {
				mockSvc := &authservice.SvcMock{}
				ur := &authservice.UserRegister{
					FirstName: "test_firstname",
					LastName:  "test_lastname",
					Username:  "test_username",
					Password:  "secret_password",
				}
				mockSvc.On("Register", mock.Anything, ur).
					Return(errors.New("random error"))

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				body, _ := json.Marshal(&authservice.UserRegister{
					FirstName: "test_firstname",
					LastName:  "test_lastname",
					Username:  "test_username",
					Password:  "secret_password",
				})
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath("/api/v1/auth/register")

				return c
			}(),
			expErr: errors.New("random error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			r := auth.NewRouter(tt.mockSvc)
			err := r.Register(tt.echoCtx)

			if tt.expErr == nil {
				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expErr, err)
			}
		})
	}
}
