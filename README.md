# OAuth2 - Open Auth 2.0 Client

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-zoox/oauth2)](https://pkg.go.dev/github.com/go-zoox/oauth2)
[![Build Status](https://github.com/go-zoox/oauth2/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/go-zoox/oauth2/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-zoox/oauth2)](https://goreportcard.com/report/github.com/go-zoox/oauth2)
[![Coverage Status](https://coveralls.io/repos/github/go-zoox/oauth2/badge.svg?branch=master)](https://coveralls.io/github/go-zoox/oauth2?branch=master)
[![GitHub issues](https://img.shields.io/github/issues/go-zoox/oauth2.svg)](https://github.com/go-zoox/oauth2/issues)
[![Release](https://img.shields.io/github/tag/go-zoox/oauth2.svg?label=Release)](https://github.com/go-zoox/oauth2/tags)

## Installation
To install the package, run:
```bash
go get github.com/go-zoox/oauth2
```

## Supported Providers

This library supports many OAuth2 providers, including:

- **Supabase** - Full-featured authentication platform
- **GitHub** - Version control and collaboration
- **Google** - Google services authentication
- **Auth0** - Identity and access management
- **Microsoft Azure** - Microsoft cloud services
- **Slack** - Team communication platform
- **Discord** - Gaming and community communication
- **Facebook** - Social networking
- **GitLab** - DevOps platform
- **Twitter** - Social media platform
- **WeChat** - Chinese messaging platform
- **Doreamon** - Custom authentication provider
- And many more...

### Supabase Provider

The Supabase provider offers seamless integration with Supabase Auth:

```go
import "github.com/go-zoox/oauth2/supabase"

// Create Supabase client
client, err := supabase.New(&supabase.SupabaseConfig{
    BaseURL:      "https://your-project.supabase.co",
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    RedirectURI:  "http://localhost:8080/auth/callback",
    Scope:        "openid email profile",
})
```

For detailed Supabase setup instructions, see the [Supabase provider documentation](supabase/README.md).

## Getting Started

### Example 1: Using only one oauth2 provider => doreamon

```go
// step1: create oauth2 middleware/handler
// file: oauth2.go
import (
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/go-zoox/logger"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/oauth2/doreamon"
)

type CreateOAuth2DoreamonHandlerConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func CreateOAuth2DoreamonHandler(cfg *CreateOAuth2DoreamonHandlerConfig) func(
	w http.ResponseWriter,
	r *http.Request,
	CheckUser func(r *http.Request) error,
	RemeberUser func(user *oauth2.User, token *oauth2.Token) error,
	Next func() error,
) error {
	originPathCookieKey := "login_from"

	client, err := doreamon.New(&doreamon.DoreamonConfig{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURI:  cfg.RedirectURI,
		Scope:        "using_doreamon",
		Version:      "2",
	})
	if err != nil {
		panic(err)
	}

	return func(
		w http.ResponseWriter,
		r *http.Request,
		RestoreUser func(r *http.Request) error,
		SaveUser func(user *oauth2.User, token *oauth2.Token) error,
		Next func() error,
	) error {
		if r.Method != "GET" {
			return Next()
		}
		path := r.URL.Path

		if path == "/login" {
			client.Authorize("memos", func(loginUrl string) {
				http.Redirect(w, r, loginUrl, http.StatusFound)
			})
			return nil
		}

		if path == "/logout" {
			client.Logout(func(logoutUrl string) {
				http.Redirect(w, r, logoutUrl, http.StatusFound)
			})
			return nil
		}

		if path == "/login/doreamon/callback" {
			code := r.FormValue("code")
			state := r.FormValue("state")

			client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
				if err != nil {
					log.Println("[OAUTH2] Login Callback Error", err)
					time.Sleep(3 * time.Second)
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}

				if err := SaveUser(user, token); err != nil {
					logger.Info("failed to save user: %#v", err)
					time.Sleep(1)

					w.WriteHeader(500)
					w.Write([]byte("Failed to create user: " + user.Email))
					return
				}

				http.Redirect(w, r, "/", http.StatusFound)
			})

			return nil
		}

		if matched, _ := regexp.MatchString("\\.(js|css|json)$", path); err == nil && matched {
			return Next()
		}

		if err := RestoreUser(r); err != nil {
			logger.Info("failed to restart user: %#v", err)
			time.Sleep(1)
			http.SetCookie(w, &http.Cookie{
				Name:  "OriginPath",
				Value: path,
			})

			http.Redirect(w, r, "/login", http.StatusFound)
			return nil
		}

		// success
		if OriginPath, err := r.Cookie(originPathCookieKey); err == nil && OriginPath.Value != "" {
			time.Sleep(1)

			http.SetCookie(w, &http.Cookie{
				Name:    originPathCookieKey,
				Value:   "",
				Expires: time.Unix(0, 0),
			})

			http.Redirect(w, r, OriginPath.Value, http.StatusFound)
			return nil
		}

		return Next()
	}
}
```

```go
// step 2: use as go http middleware
//  here is memos/echo
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	if os.Getenv("DOREAMON_CLIENT_ID") == "" {
		panic("env DOREAMON_CLIENT_ID is required")
	}
	if os.Getenv("DOREAMON_CLIENT_SECRET") == "" {
		panic("env DOREAMON_CLIENT_SECRET is required")
	}
	if os.Getenv("DOREAMON_REDIRECT_URI") == "" {
		panic("env DOREAMON_REDIRECT_URI is required")
	}

	handler := CreateOAuth2DoreamonHandler(&CreateOAuth2DoreamonHandlerConfig{
		ClientID:     os.Getenv("DOREAMON_CLIENT_ID"),
		ClientSecret: os.Getenv("DOREAMON_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("DOREAMON_REDIRECT_URI"),
	})

	return func(c echo.Context) error {
		return handler(
			c.Response().Writer,
			c.Request(),
			func(r *http.Request) error {
				userID, ok := getUserSession(c)
				if !ok {
					return fmt.Errorf("no user session found")
				}

				c.Set(getUserIDContextKey(), userID)
				userFind := &api.UserFind{
					ID: &userID,
				}
				_, err := s.Store.FindUser(c.Request().Context(), userFind)
				if err != nil {
					return err
				}

				return nil
			},
			func(user *oauth2.User, token *oauth2.Token) error {
				ctx := c.Request().Context()
				// Get Or Create User
				userFind := &api.UserFind{
					Username: &user.Email,
				}
				dbUser, err := s.Store.FindUser(ctx, userFind)
				if err != nil || dbUser == nil {
					role := api.Host
					hostUserFind := api.UserFind{
						Role: &role,
					}
					hostUser, err := s.Store.FindUser(ctx, &hostUserFind)
					if err != nil {
						return err
					}
					if hostUser != nil {
						role = api.NormalUser
					}

					userCreate := &api.UserCreate{
						Username: user.Email,
						Role:     api.Role(role),
						Nickname: user.Nickname,
						Password: random.String(32),
						OpenID:   common.GenUUID(),
					}
					dbUser, err = s.Store.CreateUser(ctx, userCreate)
					if err != nil {
						return err
					}
				}

				if err = setUserSession(c, dbUser); err != nil {
					return err
				}

				return nil
			},
			func() error {
				return next(c)
			},
		)
	}
})
```

### Example 2: Support multiple oauth2 providers: github, wechat, gitee, doreamon

```go
// @TODO connect
```

## License
GoZoox is released under the [MIT License](./LICENSE).
