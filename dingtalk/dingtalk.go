package dingtalk

// reference:
//	https://github.com/directus/directus/discussions/11881
//	https://open.dingtalk.com/document/orgapp/logon-free-third-party-websites
//	https://open.dingtalk.com/document/personalapp/obtain-user-token
//	https://open.dingtalk.com/document/personalapp/tutorial-on-how-to-obtain-logon-user-information-for-third-party

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-zoox/fetch"
	"github.com/go-zoox/oauth2"
)

type DingTalkConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *DingTalkConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "snsapi_login"
	}

	config := oauth2.Config{
		Name:         "DingTalk",
		AuthURL:      "https://login.dingtalk.com/oauth2/auth",
		TokenURL:     "https://api.dingtalk.com/v1.0/oauth2/userAccessToken", // return uid
		UserInfoURL:  "https://api.dingtalk.com/v1.0/contact/users/me",       // need uid
		LogoutURL:    "",
		Scope:        scope,
		RedirectURI:  cfg.RedirectURI,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		//
		AccessTokenAttributeName:  "accessToken",
		RefreshTokenAttributeName: "refreshToken",
		ExpiresInAttributeName:    "expireIn",
		TokenTypeAttributeName:    "",
		//
		EmailAttributeName:    "email",
		IDAttributeName:       "openId", // unionid 全钉钉统一，openid 每个公众号各一个
		NicknameAttributeName: "nick",
		AvatarAttributeName:   "avatarUrl",
		HomepageAttributeName: "url",
	}

	config.GetLoginURL = func(cfg *oauth2.Config, state string) string {
		return strings.Join([]string{
			cfg.AuthURL,
			fmt.Sprintf("?%s=", cfg.ClientIDAttributeName), cfg.ClientID,
			fmt.Sprintf("&%s=", cfg.RedirectURIAttributeName), url.QueryEscape(cfg.RedirectURI),
			fmt.Sprintf("&%s=", cfg.ResponseTypeAttributeName), "code",
			fmt.Sprintf("&%s=", cfg.ScopeAttributeName), url.QueryEscape(scope),
			fmt.Sprintf("&%s=", cfg.StateAttributeName), url.QueryEscape(state),
			fmt.Sprintf("&%s=", "prompt"), "consent",
		}, "")
	}

	config.GetAccessTokenResponse = func(cfg *oauth2.Config, code, state string) (*fetch.Response, error) {
		return fetch.Get(config.TokenURL, &fetch.Config{
			Headers: map[string]string{
				"Accept": "application/json",
			},
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
			Headers: map[string]string{
				"Accept":                      "application/json",
				"x-acs-dingtalk-access-token": token.AccessToken,
			},
		})
	}

	return oauth2.New(config)
}
