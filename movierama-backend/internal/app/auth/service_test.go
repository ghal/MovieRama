package auth_test

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"movierama/internal/app/auth"
	"movierama/internal/config"
	"movierama/internal/infra/repository/sql/user"
	sqluser "movierama/internal/infra/repository/sql/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Login(t *testing.T) {
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		sqlRepo   *sqluser.Mock
		userLogin *auth.UserLogin
		echoCtx   echo.Context
		expRes    *auth.LoginRes
		expErr    error
	}{
		"should login": {
			sqlRepo: func() *sqluser.Mock {
				userDetails := &user.AuthUserDetails{
					ID:       1,
					Username: "test_username",
					Password: "$2a$04$ShttrFUsTWC/aW1ahlr3rO2zLoQpuGfSHjKRB83f.8dJBebnBTOpG",
				}
				repo := sqluser.Mock{}
				repo.On("GetUserAuthDetails", context.TODO(), "test_username").
					Return(userDetails, nil)

				return &repo
			}(),
			userLogin: &auth.UserLogin{
				Username: "test_username",
				Password: "secret_pass",
			},
			expRes: &auth.LoginRes{
				Username: "test_username",
			},
			expErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cfg := &config.Config{
				App: config.App{
					JWTSecret: "secret",
				},
			}
			app := auth.NewService(tt.sqlRepo, cfg)
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec = httptest.NewRecorder()
			c := e.NewContext(req, rec)

			res, err := app.Login(c, tt.userLogin)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expRes.Username, res.Username)
			assert.NotEmpty(t, res.Token)
		})
	}
}

func Test_Register(t *testing.T) {
	tests := map[string]struct {
		sqlRepo      *sqluser.Mock
		userRegister *auth.UserRegister
		expErr       error
	}{
		"Should get public movies": {
			sqlRepo: func() *sqluser.Mock {
				repo := sqluser.Mock{}
				repo.On("CreateUser", context.TODO(), mock.Anything).Return(nil)

				return &repo
			}(),
			userRegister: &auth.UserRegister{
				Username:  "test_username",
				Password:  "test_password",
				FirstName: "test_firstname",
				LastName:  "test_lastname",
			},
			expErr: nil,
		},
		"Should return error on repo error": {
			sqlRepo: func() *sqluser.Mock {
				repo := sqluser.Mock{}
				repo.On("CreateUser", context.TODO(), mock.Anything).Return(errors.New("random error"))

				return &repo
			}(),
			userRegister: &auth.UserRegister{
				Username:  "test_username",
				Password:  "test_password",
				FirstName: "test_firstname",
				LastName:  "test_lastname",
			},
			expErr: errors.New("random error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			app := auth.NewService(tt.sqlRepo, config.New())

			err := app.Register(context.TODO(), tt.userRegister)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
