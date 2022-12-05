package auth_test

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
					Password: "pass",
				}
				repo := sqluser.Mock{}
				repo.On("GetUserAuthDetails", context.TODO(), "test_username").Return(userDetails, nil)

				return &repo
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/")

				return c
			}(),
			expRes: &auth.LoginRes{
				Username: "test_username",
				Token:    "",
			},
			expErr: nil,
		},
		//"Should return error on repo error": {
		//	sqlRepo: func() *sqluser.Mock {
		//		repo := sqluser.Mock{}
		//		repo.On("GetUserAuthDetails", context.TODO(), "date").Return(&user.AuthUserDetails{}, errors.New("random error"))
		//
		//		return &repo
		//	}(),
		//	expErr: errors.New("random error"),
		//},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cfg := &config.Config{
				App: config.App{
					JWTSecret: "secret",
				},
			}
			app := auth.NewService(tt.sqlRepo, cfg)

			res, err := app.Login(tt.echoCtx, tt.userLogin)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expRes, res)
		})
	}
}
