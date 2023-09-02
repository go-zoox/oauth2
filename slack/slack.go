package slack

import (
	"github.com/go-zoox/oauth2"
)

type SlackConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *SlackConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "identity.basic,identity.email,identity.avatar"
	}

	config := oauth2.Config{
		Name:         "Slack",
		AuthURL:      "https://slack.com/oauth/v2/authorize",
		TokenURL:     "https://slack.com/api/oauth.v2.access",
		UserInfoURL:  "https://slack.com/api/users.identity",
		LogoutURL:    "https://slack.com/api/auth.signout",
		Scope:        scope,
		RedirectURI:  cfg.RedirectURI,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		//
		ScopeAttributeName: "user_scope",
		//
		AccessTokenAttributeName: "authed_user.access_token",
		TokenTypeAttributeName:   "authed_user.token_type",
		//
		EmailAttributeName:    "user.email",
		IDAttributeName:       "user.id",
		NicknameAttributeName: "user.name",
		AvatarAttributeName:   "user.image_48",
	}

	return oauth2.New(config)
}
