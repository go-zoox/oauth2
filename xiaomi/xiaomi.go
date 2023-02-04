package xiaomi

// reference:
//	https://dev.mi.com/docs/passport/oauth2/
//	https://dev.mi.com/distribute/doc/details?pId=1515
//	https://dev.mi.com/docs/passport/open-api/

import (
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/oauth2"
)

type XiaoMiConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *XiaoMiConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "snsapi_login"
	}

	config := oauth2.Config{
		Name:         "XiaoMi",
		AuthURL:      "https://account.xiaomi.com/oauth2/authorize",
		TokenURL:     "https://account.xiaomi.com/oauth2/token",
		UserInfoURL:  "https://open.account.xiaomi.com/user/profile",
		LogoutURL:    "",
		Scope:        scope,
		RedirectURI:  cfg.RedirectURI,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		//
		AccessTokenAttributeName:  "accessToken",
		RefreshTokenAttributeName: "refreshToken",
		ExpiresInAttributeName:    "expireIn",
		TokenTypeAttributeName:    "token_type",
		//
		EmailAttributeName:    "",
		IDAttributeName:       "userId",
		NicknameAttributeName: "miliaoNick",
		AvatarAttributeName:   "miliaoIcon",
		HomepageAttributeName: "",
	}

	config.GetUserResponse = func(config *oauth2.Config, token *oauth2.Token, code string) (*fetch.Response, error) {
		return fetch.Get(config.UserInfoURL, &fetch.Config{
			Query: fetch.Query{
				"clientId": cfg.ClientID,
				"token":    token.AccessToken,
			},
		})
	}

	return oauth2.New(config)
}
