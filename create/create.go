package create

import (
	"fmt"

	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/oauth2/doreamon"
	"github.com/go-zoox/oauth2/feishu"
	"github.com/go-zoox/oauth2/github"
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
	default:
		return nil, fmt.Errorf("oauth2: provider(%s) not supported", provider)
	}
}
