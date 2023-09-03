package create

import (
	"fmt"

	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/oauth2/auth0"
	"github.com/go-zoox/oauth2/doreamon"
	"github.com/go-zoox/oauth2/feishu"
	"github.com/go-zoox/oauth2/github"
	"github.com/go-zoox/oauth2/gitlab"
	"github.com/go-zoox/oauth2/google"
	"github.com/go-zoox/oauth2/kakao"
	"github.com/go-zoox/oauth2/slack"
)

func Create(provider string, cfg *oauth2.Config) (oauth2.Client, error) {
	switch provider {
	case "doreamon":
		return doreamon.New(&doreamon.DoreamonConfig{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURI:  cfg.RedirectURI,
			Scope:        cfg.Scope,
			Version:      "v2",
		})
	case "github":
		return github.New(&github.GitHubConfig{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURI:  cfg.RedirectURI,
			Scope:        cfg.Scope,
		})
	case "feishu":
		return feishu.New(&feishu.FeishuConfig{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURI:  cfg.RedirectURI,
			Scope:        cfg.Scope,
		})
	case "gitlab":
		return gitlab.New(&gitlab.GitLabConfig{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURI:  cfg.RedirectURI,
			Scope:        cfg.Scope,
		})
	case "slack":
		return slack.New(&slack.SlackConfig{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURI:  cfg.RedirectURI,
			Scope:        cfg.Scope,
		})
	case "kakao":
		return kakao.New(&kakao.KakaoConfig{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURI:  cfg.RedirectURI,
			Scope:        cfg.Scope,
		})
	case "google":
		return google.New(&google.GoogleConfig{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURI:  cfg.RedirectURI,
			Scope:        cfg.Scope,
		})
	//
	case "auth0":
		return auth0.New(&auth0.Auth0Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURI:  cfg.RedirectURI,
			Scope:        cfg.Scope,
			BaseURL:      cfg.BaseURL,
		})
	default:
		return nil, fmt.Errorf("oauth2: provider(%s) not supported", provider)
	}
}
