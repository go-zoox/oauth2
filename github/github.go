package github

import (
	"fmt"
	"net/url"

	"github.com/go-zoox/oauth2"
)

type GitHubConfig struct {
	// config.Config
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *GitHubConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "user:email"
	}

	config := oauth2.Config{
		Name:         "GitHub",
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		UserInfoURL:  "https://api.github.com/user",
		LogoutURL:    "https://github.com/logout",
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
		IDAttributeName:       "login",
		NicknameAttributeName: "name",
		AvatarAttributeName:   "avatar_url",
		HomepageAttributeName: "html_url",
	}

	config.GetRegisterURL = func(oac *oauth2.Config) string {
		returnTo := fmt.Sprintf("https://github.com/login?client_id=%s", cfg.ClientID)
		return fmt.Sprintf("https://github.com/signup?return_to=%s", url.QueryEscape(returnTo))
	}

	return oauth2.New(config)
}
