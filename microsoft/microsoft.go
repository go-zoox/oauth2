package microsoft

import (
	"github.com/go-zoox/oauth2"
)

type MicrosoftConfig struct {
	// config.Config
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *MicrosoftConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "openid offline_access user.read"
	}

	config := oauth2.Config{
		Name:        "Microsoft",
		AuthURL:     "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		TokenURL:    "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		UserInfoURL: "https://graph.microsoft.com/v1.0/me",
		// LogoutURL:    "https://login.microsoftonline.com/logout",
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
		EmailAttributeName:    "mail",
		IDAttributeName:       "id",
		NicknameAttributeName: "displayName",
		// AvatarAttributeName:   "picture",
		// HomepageAttributeName: "web_url",
	}

	return oauth2.New(config)
}
