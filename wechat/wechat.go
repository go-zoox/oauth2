package wechat

// reference:
//	https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
//	https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html

import (
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/oauth2"
)

type WechatConfig struct {
	// config.Config
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *WechatConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "snsapi_login"
	}

	config := oauth2.Config{
		Name:            "Wechat",
		AuthURL:         "https://open.weixin.qq.com/connect/qrconnect",
		TokenURL:        "https://api.weixin.qq.com/sns/oauth2/access_token",
		RefershTokenURL: "https://api.weixin.qq.com/sns/oauth2/refresh_token",
		UserInfoURL:     "https://api.weixin.qq.com/sns/userinfo",
		LogoutURL:       "",
		Scope:           scope,
		RedirectURI:     cfg.RedirectURI,
		ClientID:        cfg.ClientID,
		ClientSecret:    cfg.ClientSecret,
		//
		ClientIDAttributeName:     "appid",
		AccessTokenAttributeName:  "access_token",
		RefreshTokenAttributeName: "refresh_token",
		ExpiresInAttributeName:    "expires_in",
		TokenTypeAttributeName:    "",
		//
		EmailAttributeName:    "",
		IDAttributeName:       "openid", // unionid 全微信统一，openid 每个公众号各一个
		NicknameAttributeName: "name",
		AvatarAttributeName:   "profile_image_url",
		HomepageAttributeName: "url",
	}

	config.GetAccessTokenResponse = func(cfg *oauth2.Config, code, state string) (*fetch.Response, error) {
		return fetch.Get(config.TokenURL, &fetch.Config{
			// Headers: map[string]string{
			// 	"Authorization": "Bearer " + token.AccessToken,
			// },
			Query: fetch.ConfigQuery{
				"appid":      cfg.ClientID,
				"secret":     cfg.ClientSecret,
				"code":       code,
				"grant_type": "authorization_code",
			},
		})
	}

	config.GetUserResponse = func(config *oauth2.Config, token *oauth2.Token, code string) (*fetch.Response, error) {
		return fetch.Get(config.UserInfoURL, &fetch.Config{
			// Headers: map[string]string{
			// 	"Authorization": "Bearer " + token.AccessToken,
			// },
			Query: fetch.ConfigQuery{
				"access_token": token.AccessToken,
				"openid":       token.Raw.Get("openid").String(),
			},
		})
	}

	return oauth2.New(config)
}
