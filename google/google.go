package google

import (
	"github.com/go-zoox/oauth2"
)

type GoogleConfig struct {
	// config.Config
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *GoogleConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"
	}

	config := oauth2.Config{
		Name:         "Google",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		UserInfoURL:  "https://www.googleapis.com/oauth2/v1/userinfo",
		LogoutURL:    "https://accounts.google.com/logout",
		Scope:        scope,
		RedirectURI:  cfg.RedirectURI,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		//
		AccessTokenAttributeName:  "access_token",
		RefreshTokenAttributeName: "refresh_token",
		ExpiresInAttributeName:    "expires_in",
		TokenTypeAttributeName:    "token_type",
		//
		EmailAttributeName:    "email",
		IDAttributeName:       "id",
		NicknameAttributeName: "name",
		AvatarAttributeName:   "picture",
		// HomepageAttributeName: "web_url",
	}

	return oauth2.New(config)
}
