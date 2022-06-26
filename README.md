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

## Getting Started

```go
// config: oauth2/oauth2.go
package oauth2

import (
	"errors"
	"log"
	"strings"
)

type Client struct {
	*oauth2.Client
}

func New() *Client {
	client, err := oauth2.New(oauth2.Config{
		Name:                      conf.Oauth2.Name,
		AuthUrl:                   conf.Oauth2.AuthUrl,
		TokenUrl:                  conf.Oauth2.TokenUrl,
		UserInfoUrl:               conf.Oauth2.UserInfoUrl,
		LogoutUrl:                 conf.Oauth2.LogoutUrl,
		RedirectUri:               conf.Oauth2.ServerUrl + "/login/callback",
		Scope:                     conf.Oauth2.Scope,
		ClientId:                  conf.Oauth2.ClientId,
		ClientSecret:              conf.Oauth2.ClientSecret,
		AccessTokenAttributeName:  conf.Oauth2.AccessTokenAttributeName,
		RefreshTokenAttributeName: conf.Oauth2.RefreshTokenAttributeName,
		EmailAttributeName:        conf.Oauth2.EmailAttributeName,
		IdAttributeName:           conf.Oauth2.IdAttributeName,
		NicknameAttributeName:     conf.Oauth2.NicknameAttributeName,
		AvatarAttributeName:       conf.Oauth2.AvatarAttributeName,
		PermissionsAttributeName:  conf.Oauth2.PermissionsAttributeName,
	})
	if err != nil {
		panic("oauth2 init error, invalid config")
	}

	return &Client{client}
}
```

```go
// main logic
// on login
func login(w http.ResponseWriter, r *http.Request) {
  client := oauth2.New()
  client.Authorize(func(loginUrl string) {
    http.Redirect(w, r, loginUrl, http.StatusFound)
  })
}

// on login callback
func loginCallback(w http.ResponseWriter, r *http.Request) {
  code := r.FormValue("code")
  state := r.FormValue("state")

  client := oauth2.New()
  client.Callback(code, state, func(user *oauth2.User, err error) {
    if err != nil {
      log.Println("[OAUTH2] Login Callback Error", err)
      http.Redirect(w, r, "/login", http.StatusFound)
      return
    }

		// Check Permission
		if err := validatePermission(user, token); err != nil {
			log.Println("[OAUTH2] Permission Denied", user.Email)
			cb(nil, errors.New("permission denied"))
			return
		}

		// Get Or Create User
		dbUser, err := db.GetOrCreateUserByEmail(user.Email, user)

		isAdmin := false
		log.Println("[OAUTH2] Permissions", user.Permissions)
		if user.Permissions != nil {
			for _, p := range user.Permissions {
				if strings.ToUpper(p) == "ADMIN" {
					isAdmin = true
					break
				}
			}
		}

		if isAdmin != dbUser.IsAdmin {
			log.Println("[OAUTH2] Admin Change: ", dbUser.IsAdmin, " -> ", isAdmin)
			dbUser.IsAdmin = isAdmin
			if err := db.UpdateUser(dbUser); err != nil {
				log.Println("[OAUTH2] Update User Error", user.Email)
				cb(nil, errors.New("update user error"))
				return
			}
		}

    // login success
    session := sessions.Default(r)
    session.Set("user_id", user.Id)
    session.Save(r, w)

    http.Redirect(w, r, "/", http.StatusFound)
  })
}

// on logout
func logout(w http.ResponseWriter, r *http.Request) {
  client := oauth2.New()
  client.Logout(func(logoutUrl string) {
    http.Redirect(w, r, logoutUrl, http.StatusFound)
  })
}

func validatePermission(user *oauth2.User, token *oauth2.Token) error {
	if len(conf.Oauth2.AllowPermissions) == 0 {
		return nil
	}

	oauth2_allow_permissions := strings.Split(conf.Oauth2.AllowPermissions, ",")

	for _, p := range user.Permissions {
		for _, ap := range oauth2_allow_permissions {
			if p == ap {
				return nil
			}
		}
	}

	return errors.New("permission denied")
}
```

## License
GoZoox is released under the [MIT License](./LICENSE).