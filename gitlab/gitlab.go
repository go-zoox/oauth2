package gitlab

import (
	"github.com/go-zoox/oauth2"
)

type GitLabConfig struct {
	// config.Config
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *GitLabConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "read_user profile"
	}

	config := oauth2.Config{
		Name:         "GitLab",
		AuthURL:      "https://gitlab.com/oauth/authorize",
		TokenURL:     "https://gitlab.com/oauth/token",
		UserInfoURL:  "https://gitlab.com/api/v4/user",
		LogoutURL:    "https://gitlab.com/logout",
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
		AvatarAttributeName:   "avatar_url",
		HomepageAttributeName: "web_url",
	}

	return oauth2.New(config)
}
