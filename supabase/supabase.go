package supabase

import (
	"fmt"
	"net/url"

	"github.com/go-zoox/oauth2"
)

type SupabaseConfig struct {
	// Basic OAuth2 configuration
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	// Supabase specific configuration
	BaseURL string `json:"base_url"` // e.g., "https://your-project.supabase.co"
}

func New(cfg *SupabaseConfig) (oauth2.Client, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("supabase: base_url is required")
	}

	scope := cfg.Scope
	if scope == "" {
		scope = "openid email profile"
	}

	// Ensure BaseURL doesn't have trailing slash
	baseURL := cfg.BaseURL
	if baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	config := oauth2.Config{
		Name:         "Supabase",
		AuthURL:      baseURL + "/auth/v1/authorize",
		TokenURL:     baseURL + "/auth/v1/token",
		UserInfoURL:  baseURL + "/auth/v1/user",
		LogoutURL:    baseURL + "/auth/v1/logout",
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
		EmailAttributeName:    "email",
		IDAttributeName:       "id",
		NicknameAttributeName: "user_metadata.full_name",
		AvatarAttributeName:   "user_metadata.avatar_url",
		HomepageAttributeName: "user_metadata.website",
		//
		BaseURL: baseURL,
	}

	// Custom register URL for Supabase
	config.GetRegisterURL = func(oac *oauth2.Config) string {
		// Supabase doesn't have a standard register endpoint, redirect to auth
		return fmt.Sprintf("%s/auth/v1/signup", baseURL)
	}

	// Custom login URL to handle Supabase's OAuth flow
	config.GetLoginURL = func(oac *oauth2.Config, state string) string {
		clientID := oac.ClientID
		redirectURI := oac.RedirectURI
		scope := oac.Scope

		params := url.Values{}
		params.Add("client_id", clientID)
		params.Add("redirect_uri", redirectURI)
		params.Add("response_type", "code")
		params.Add("scope", scope)
		params.Add("state", state)

		return fmt.Sprintf("%s?%s", oac.AuthURL, params.Encode())
	}

	return oauth2.New(config)
}