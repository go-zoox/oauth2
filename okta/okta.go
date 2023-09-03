package okta

import (
	"fmt"

	"github.com/go-zoox/oauth2"
)

type OktaConfig struct {
	// config.Config
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	//
	BaseURL string `json:"base_url"`
}

func New(cfg *OktaConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "profile email openid"
	}

	config := oauth2.Config{
		Name:         "Okta",
		AuthURL:      fmt.Sprintf("%s/oauth2/default/v1/authorize", cfg.BaseURL),
		TokenURL:     fmt.Sprintf("%s/oauth2/default/v1/token", cfg.BaseURL),
		UserInfoURL:  fmt.Sprintf("%s/oauth2/default/v1/userinfo", cfg.BaseURL),
		LogoutURL:    fmt.Sprintf("%s/logout", cfg.BaseURL),
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
		IDAttributeName:       "sub",
		NicknameAttributeName: "name",
		// AvatarAttributeName:   "picture",
		// HomepageAttributeName: "html_url",
	}

	return oauth2.New(config)
}
