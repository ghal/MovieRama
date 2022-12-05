package movie_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"movierama/internal/app/auth"
	movieservice "movierama/internal/app/movie"
	"movierama/internal/infra/http/router/movie"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_GetMoviesPublic(t *testing.T) {
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc   *movieservice.SvcMock
		jwtConfig middleware.JWTConfig
		echoCtx   echo.Context
		expRes    string
		expErr    error
	}{
		"Should succeed on GetMoviesPublic call": {
			mockSvc: func() *movieservice.SvcMock {
				mockSvc := &movieservice.SvcMock{}
				mockSvc.On("GetMoviesPublic", context.TODO(), "date").
					Return(&movieservice.GetmoviesRes{
						Movies: []movieservice.Movie{
							{
								ID:          1,
								Title:       "movie title",
								Description: "movie description",
								UserID:      2,
								PostedBy:    "Full Name",
								Likes:       3,
								Hates:       4,
								UserLiked:   true,
								UserHated:   false,
								IsSameUser:  false,
								TimeAgo:     "1 minute ago",
							},
						},
					}, nil)

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/?sort=date", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/movies")

				return c
			}(),
			jwtConfig: middleware.JWTConfig{},
			expRes: func() string {
				b, err := json.Marshal(&movieservice.GetmoviesRes{
					Movies: []movieservice.Movie{
						{
							ID:          1,
							Title:       "movie title",
							Description: "movie description",
							UserID:      2,
							PostedBy:    "Full Name",
							Likes:       3,
							Hates:       4,
							UserLiked:   true,
							UserHated:   false,
							IsSameUser:  false,
							TimeAgo:     "1 minute ago",
						},
					},
				})
				assert.Nil(t, err)

				return string(b) + "\n"
			}(),
		},
		"Should return error on GetMoviesPublic error ": {
			mockSvc: func() *movieservice.SvcMock {
				mockSvc := &movieservice.SvcMock{}
				mockSvc.On("GetMoviesPublic", context.TODO(), "likes").
					Return(&movieservice.GetmoviesRes{}, errors.New("random error"))

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/?sort=likes", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/movies")

				return c
			}(),
			expErr: errors.New("random error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			r := movie.NewRouter(tt.mockSvc, tt.jwtConfig)
			err := r.GetMoviesPublic(tt.echoCtx)

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

func TestRouter_GetUserMoviesPublic(t *testing.T) {
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc   *movieservice.SvcMock
		jwtConfig middleware.JWTConfig
		echoCtx   echo.Context
		expRes    string
		expErr    error
	}{
		"Should succeed on GetUserMoviesPublic call": {
			mockSvc: func() *movieservice.SvcMock {
				mockSvc := &movieservice.SvcMock{}
				mockSvc.On("GetUserMoviesPublic", context.TODO(), 3, "date").
					Return(&movieservice.GetmoviesRes{
						Movies: []movieservice.Movie{
							{
								ID:          1,
								Title:       "movie title",
								Description: "movie description",
								UserID:      2,
								PostedBy:    "Full Name",
								Likes:       3,
								Hates:       4,
								UserLiked:   true,
								UserHated:   false,
								IsSameUser:  false,
								TimeAgo:     "1 minute ago",
							},
						},
					}, nil)

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/?sort=date", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/users/:user_id/movies")
				c.SetParamNames("user_id")
				c.SetParamValues("3")

				return c
			}(),
			jwtConfig: middleware.JWTConfig{},
			expRes: func() string {
				b, err := json.Marshal(&movieservice.GetmoviesRes{
					Movies: []movieservice.Movie{
						{
							ID:          1,
							Title:       "movie title",
							Description: "movie description",
							UserID:      2,
							PostedBy:    "Full Name",
							Likes:       3,
							Hates:       4,
							UserLiked:   true,
							UserHated:   false,
							IsSameUser:  false,
							TimeAgo:     "1 minute ago",
						},
					},
				})
				assert.Nil(t, err)

				return string(b) + "\n"
			}(),
		},
		"Should return error on GetUserMoviesPublic error ": {
			mockSvc: func() *movieservice.SvcMock {
				mockSvc := &movieservice.SvcMock{}
				mockSvc.On("GetUserMoviesPublic", context.TODO(), 3, "likes").
					Return(&movieservice.GetmoviesRes{}, errors.New("random error"))

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/?sort=likes", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/users/:user_id/movies")
				c.SetParamNames("user_id")
				c.SetParamValues("3")

				return c
			}(),
			expErr: errors.New("random error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			r := movie.NewRouter(tt.mockSvc, tt.jwtConfig)
			err := r.GetUserMoviesPublic(tt.echoCtx)

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

func TestRouter_GetMovies(t *testing.T) {
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc   *movieservice.SvcMock
		jwtConfig middleware.JWTConfig
		expRes    string
		expErr    error
	}{
		"Should succeed on GetMovies call": {
			mockSvc: func() *movieservice.SvcMock {
				mockSvc := &movieservice.SvcMock{}
				mockSvc.On("GetMovies", context.WithValue(context.TODO(), movieservice.AuthUserIDContextKey, 1), "date").
					Return(&movieservice.GetmoviesRes{
						Movies: []movieservice.Movie{
							{
								ID:          1,
								Title:       "movie title",
								Description: "movie description",
								UserID:      2,
								PostedBy:    "Full Name",
								Likes:       3,
								Hates:       4,
								UserLiked:   true,
								UserHated:   false,
								IsSameUser:  false,
								TimeAgo:     "1 minute ago",
							},
						},
					}, nil)

				return mockSvc
			}(),
			jwtConfig: middleware.JWTConfig{
				Claims: &auth.JwtCustomClaims{
					UserID:         1,
					StandardClaims: jwt.StandardClaims{},
				},
				SigningKey: []byte("secret"),
			},
			expRes: func() string {
				b, err := json.Marshal(&movieservice.GetmoviesRes{
					Movies: []movieservice.Movie{
						{
							ID:          1,
							Title:       "movie title",
							Description: "movie description",
							UserID:      2,
							PostedBy:    "Full Name",
							Likes:       3,
							Hates:       4,
							UserLiked:   true,
							UserHated:   false,
							IsSameUser:  false,
							TimeAgo:     "1 minute ago",
						},
					},
				})
				assert.Nil(t, err)

				return string(b) + "\n"
			}(),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?sort=date", nil)
			rec = httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user", &jwt.Token{
				Claims: tt.jwtConfig.Claims,
			})

			r := movie.NewRouter(tt.mockSvc, tt.jwtConfig)
			r.AppendRoutes(e)
			err := r.GetMovies(c)
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

func TestRouter_GetUserMovies(t *testing.T) {
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc   *movieservice.SvcMock
		jwtConfig middleware.JWTConfig
		expRes    string
		expErr    error
	}{
		"Should succeed on GetUserMovies call": {
			mockSvc: func() *movieservice.SvcMock {
				mockSvc := &movieservice.SvcMock{}
				mockSvc.On("GetUserMovies", context.WithValue(context.TODO(), movieservice.AuthUserIDContextKey, 1), 2, "date").
					Return(&movieservice.GetmoviesRes{
						Movies: []movieservice.Movie{
							{
								ID:          1,
								Title:       "movie title",
								Description: "movie description",
								UserID:      2,
								PostedBy:    "Full Name",
								Likes:       3,
								Hates:       4,
								UserLiked:   true,
								UserHated:   false,
								IsSameUser:  false,
								TimeAgo:     "1 minute ago",
							},
						},
					}, nil)

				return mockSvc
			}(),
			jwtConfig: middleware.JWTConfig{
				Claims: &auth.JwtCustomClaims{
					UserID:         1,
					StandardClaims: jwt.StandardClaims{},
				},
				SigningKey: []byte("secret"),
			},
			expRes: func() string {
				b, err := json.Marshal(&movieservice.GetmoviesRes{
					Movies: []movieservice.Movie{
						{
							ID:          1,
							Title:       "movie title",
							Description: "movie description",
							UserID:      2,
							PostedBy:    "Full Name",
							Likes:       3,
							Hates:       4,
							UserLiked:   true,
							UserHated:   false,
							IsSameUser:  false,
							TimeAgo:     "1 minute ago",
						},
					},
				})
				assert.Nil(t, err)

				return string(b) + "\n"
			}(),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?sort=date", nil)
			rec = httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("user_id")
			c.SetParamValues("2")
			c.Set("user", &jwt.Token{
				Claims: tt.jwtConfig.Claims,
			})

			r := movie.NewRouter(tt.mockSvc, tt.jwtConfig)
			r.AppendRoutes(e)
			err := r.GetUserMovies(c)
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

func TestRouter_CreateMovie(t *testing.T) {
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc   *movieservice.SvcMock
		jwtConfig middleware.JWTConfig
		expErr    error
	}{
		"Should succeed on CreateMovie call": {
			mockSvc: func() *movieservice.SvcMock {
				mockSvc := &movieservice.SvcMock{}
				m := movie.NewMovie{
					Title:       "Titanic",
					Description: "Titanic Description",
				}
				mockSvc.On("CreateMovie", context.WithValue(context.TODO(),
					movieservice.AuthUserIDContextKey, 1), movieservice.NewMovie(m)).
					Return(nil)

				return mockSvc
			}(),
			jwtConfig: middleware.JWTConfig{
				Claims: &auth.JwtCustomClaims{
					UserID:         1,
					StandardClaims: jwt.StandardClaims{},
				},
				SigningKey: []byte("secret"),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			body, err := json.Marshal(movie.NewMovie{
				Title:       "Titanic",
				Description: "Titanic Description",
			})
			req := httptest.NewRequest(http.MethodPost, "/?sort=date", strings.NewReader(string(body)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec = httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("user_id")
			c.SetParamValues("2")
			c.Set("user", &jwt.Token{
				Claims: tt.jwtConfig.Claims,
			})

			r := movie.NewRouter(tt.mockSvc, tt.jwtConfig)
			r.AppendRoutes(e)
			err = r.CreateMovie(c)
			if tt.expErr == nil {
				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expErr, err)
			}
		})
	}
}

func TestRouter_MakeAction(t *testing.T) {
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc   *movieservice.SvcMock
		jwtConfig middleware.JWTConfig
		expErr    error
	}{
		"Should succeed on MakeAction call": {
			mockSvc: func() *movieservice.SvcMock {
				mockSvc := &movieservice.SvcMock{}
				mockSvc.On("Action", context.WithValue(context.TODO(),
					movieservice.AuthUserIDContextKey, 1), 2, "like").
					Return(nil)

				return mockSvc
			}(),
			jwtConfig: middleware.JWTConfig{
				Claims: &auth.JwtCustomClaims{
					UserID:         1,
					StandardClaims: jwt.StandardClaims{},
				},
				SigningKey: []byte("secret"),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			body, err := json.Marshal(movie.NewMovie{
				Title:       "Titanic",
				Description: "Titanic Description",
			})
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec = httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("movie_id", "action")
			c.SetParamValues("2", "like")
			c.Set("user", &jwt.Token{
				Claims: tt.jwtConfig.Claims,
			})

			r := movie.NewRouter(tt.mockSvc, tt.jwtConfig)
			r.AppendRoutes(e)
			err = r.MakeAction(c)
			if tt.expErr == nil {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expErr, err)
			}
		})
	}
}

func TestRouter_RemoveAction(t *testing.T) {
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc   *movieservice.SvcMock
		jwtConfig middleware.JWTConfig
		expErr    error
	}{
		"Should succeed on RemoveAction call": {
			mockSvc: func() *movieservice.SvcMock {
				mockSvc := &movieservice.SvcMock{}
				mockSvc.On("RemoveAction", context.WithValue(context.TODO(),
					movieservice.AuthUserIDContextKey, 1), 2, "like").
					Return(nil)

				return mockSvc
			}(),
			jwtConfig: middleware.JWTConfig{
				Claims: &auth.JwtCustomClaims{
					UserID:         1,
					StandardClaims: jwt.StandardClaims{},
				},
				SigningKey: []byte("secret"),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			body, err := json.Marshal(movie.NewMovie{
				Title:       "Titanic",
				Description: "Titanic Description",
			})
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec = httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("movie_id", "action")
			c.SetParamValues("2", "like")
			c.Set("user", &jwt.Token{
				Claims: tt.jwtConfig.Claims,
			})

			r := movie.NewRouter(tt.mockSvc, tt.jwtConfig)
			r.AppendRoutes(e)
			err = r.RemoveAction(c)
			if tt.expErr == nil {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expErr, err)
			}
		})
	}
}
