package weibo

// reference: https://open.weibo.com/wiki/%E5%BE%AE%E5%8D%9AAPI#OAuth2

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-zoox/fetch"
	"github.com/go-zoox/oauth2"
)

type WeiboConfig struct {
	// config.Config
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *WeiboConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "user:email"
	}

	config := oauth2.Config{
		Name:         "Weibo",
		AuthURL:      "https://api.weibo.com/oauth2/authorize",
		TokenURL:     "https://api.weibo.com/oauth2/access_token", // return uid
		UserInfoURL:  "https://api.weibo.com/2/users/show.json",   // need uid
		LogoutURL:    "https://open.weibo.com/logout.php",
		Scope:        scope,
		RedirectURI:  cfg.RedirectURI,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		//
		AccessTokenAttributeName:  "access_token",
		RefreshTokenAttributeName: "",
		ExpiresInAttributeName:    "expires_in",
		TokenTypeAttributeName:    "",
		//
		EmailAttributeName:    "",
		IDAttributeName:       "id",
		NicknameAttributeName: "name",
		AvatarAttributeName:   "profile_image_url",
		HomepageAttributeName: "url",
	}

	config.GenerateLogoutURL = func(cfg *oauth2.Config) string {
		return strings.Join([]string{
			cfg.LogoutURL,
			fmt.Sprintf("?%s=", "backurl"), url.QueryEscape(cfg.RedirectURI),
		}, "")
	}

	config.GetUserResponse = func(config *oauth2.Config, token *oauth2.Token, code string) (*fetch.Response, error) {
		response, err := fetch.Get("https://api.weibo.com/2/account/get_uid.json", &fetch.Config{
			// Headers: map[string]string{
			// 	"Authorization": "Bearer " + token.AccessToken,
			// },
			Query: fetch.ConfigQuery{
				"access_token": token.AccessToken,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch weibo uid with access_token")
		}

		uid := response.Get("uid").String()
		if uid == "" {
			return nil, fmt.Errorf("get weibo uid empty string, response: %s", response.String())
		}

		return fetch.Get(config.UserInfoURL, &fetch.Config{
			// Headers: map[string]string{
			// 	"Authorization": "Bearer " + token.AccessToken,
			// },
			Query: fetch.ConfigQuery{
				"access_token": token.AccessToken,
				"uid":          uid,
			},
		})
	}

	return oauth2.New(config)
}
