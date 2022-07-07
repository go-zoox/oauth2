package doreamon

import (
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/oauth2/config"
)

type DoreamonConfig struct {
	config.Config
	Version string
}

func New(cfg *DoreamonConfig) (*oauth2.Client, error) {
	authURL := "https://login.zcorky.com/authorize"
	if cfg.Version != "" {
		authURL = "https://login.zcorky.com/v2/authorize"
	}
	scope := cfg.Scope
	if scope == "" {
		scope = "openid email profile"
	}

	config := oauth2.Config{
		Name:         "哆啦A梦",
		AuthURL:      authURL,
		TokenURL:     "https://login.zcorky.com/token",
		UserInfoURL:  "https://login.zcorky.com/user",
		LogoutURL:    "https://login.zcorky.com/logout",
		Scope:        scope,
		RedirectURI:  cfg.RedirectURI,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		//
		AccessTokenAttributeName:  "access_token",
		RefreshTokenAttributeName: "refresh_token",
		ExpiresInAttributeName:    "expires_in",
		TokenTypeAttributeName:    "token_type",
		EmailAttributeName:        "email",
		IDAttributeName:           "email",
		NicknameAttributeName:     "nickname",
		AvatarAttributeName:       "avatar",
		PermissionsAttributeName:  "permissions",
		GroupsAttributeName:       "groups",
	}

	return oauth2.New(config)
}
