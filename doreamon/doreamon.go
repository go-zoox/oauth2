package doreamon

import (
	"github.com/go-zoox/oauth2"
)

type DoreamonConfig struct {
	// config.Config
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	Version      string `json:"version"`
}

func New(cfg *DoreamonConfig) (oauth2.Client, error) {
	authURL := "https://login.zcorky.com/authorize"
	logoutURL := "https://login.zcorky.com/logout"
	if cfg.Version != "" {
		authURL = "https://login.zcorky.com/v2/authorize"
		logoutURL = "https://login.zcorky.com/v2/logout"
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
		LogoutURL:    logoutURL,
		RegisterURL:  "https://login.zcorky.com/register",
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
