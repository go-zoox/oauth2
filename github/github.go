package github

import (
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/oauth2/config"
)

type GitHubConfig struct {
	config.Config
}

func New(cfg *GitHubConfig) (*oauth2.Client, error) {
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

	return oauth2.New(config)
}
