package feishu

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-zoox/fetch"
	"github.com/go-zoox/oauth2"
)

type FeishuConfig struct {
	// config.Config
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func New(cfg *FeishuConfig) (oauth2.Client, error) {
	scope := cfg.Scope
	if scope == "" {
		scope = "user:email"
	}

	config := oauth2.Config{
		Name:         "飞书",
		AuthURL:      "https://open.feishu.cn/open-apis/authen/v1/index",
		TokenURL:     "https://open.feishu.cn/open-apis/authen/v1/access_token",
		UserInfoURL:  "https://open.feishu.cn/open-apis/authen/v1/user_info",
		LogoutURL:    "https://open.feishu.cn/open-apis/authen/v1/logout",
		Scope:        scope,
		RedirectURI:  cfg.RedirectURI,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		//
		ClientIDAttributeName:     "app_id",
		ClientSecretAttributeName: "app_secret",
		RedirectURIAttributeName:  "redirect_uri",
		ResponseTypeAttributeName: "response_type",
		ScopeAttributeName:        "scope",
		StateAttributeName:        "state",
		//
		AccessTokenAttributeName:  "data.access_token",
		RefreshTokenAttributeName: "data.refresh_token",
		ExpiresInAttributeName:    "data.expires_in",
		TokenTypeAttributeName:    "data.token_type",
		//
		EmailAttributeName:    "data.enterprise_email",
		IDAttributeName:       "data.union_id",
		NicknameAttributeName: "data.name",
		AvatarAttributeName:   "data.avatar_url",
	}

	config.GetAccessTokenResponse = func(cfg *oauth2.Config, code string, state string) (*fetch.Response, error) {
		response, err := fetch.Post("https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal", &fetch.Config{
			Body: map[string]string{
				"app_id":     cfg.ClientID,
				"app_secret": cfg.ClientSecret,
			},
		})
		if err != nil {
			return nil, err
		}

		app_access_token := response.Get("app_access_token").String()

		return fetch.Post(cfg.TokenURL, &fetch.Config{
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", app_access_token),
				"Content-Type":  "application/json; charset=utf-8",
			},
			Body: map[string]string{
				"grant_type": "authorization_code",
				"code":       code,
			},
		})
	}

	config.GetRegisterURL = func(oac *oauth2.Config) string {
		loginURL := strings.Join([]string{
			oac.AuthURL,
			fmt.Sprintf("?%s=", oac.ClientIDAttributeName), oac.ClientID,
			fmt.Sprintf("&%s=", oac.RedirectURIAttributeName), url.QueryEscape(oac.RedirectURI),
			fmt.Sprintf("&%s=", oac.ResponseTypeAttributeName), "code",
			fmt.Sprintf("&%s=", oac.ScopeAttributeName), url.QueryEscape(oac.Scope),
			fmt.Sprintf("&%s=", oac.StateAttributeName), url.QueryEscape("anything"),
		}, "")

		return fmt.Sprintf("https://www.feishu.cn/accounts/page/ug_register?redirect_uri=%s", url.QueryEscape(loginURL))
	}

	return oauth2.New(config)
}
