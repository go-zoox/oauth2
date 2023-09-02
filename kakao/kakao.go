package kakao

import (
	"github.com/go-zoox/oauth2"
)

type KakaoConfig struct {
	// config.Config
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *KakaoConfig) (oauth2.Client, error) {
	config := oauth2.Config{
		Name:        "Kakao",
		AuthURL:     "https://kauth.kakao.com/oauth/authorize",
		TokenURL:    "https://kauth.kakao.com/oauth/token",
		UserInfoURL: "https://kapi.kakao.com/v2/user/me",
		LogoutURL:   "https://kapi.kakao.com/logout",
		// Scope:        cfg.Scope,
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
		NicknameAttributeName: "properties.nickname",
		AvatarAttributeName:   "properties.profile_image",
		// HomepageAttributeName: "web_url",
	}

	return oauth2.New(config)
}
