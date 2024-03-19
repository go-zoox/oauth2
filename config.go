package oauth2

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-zoox/fetch"
)

// Config is the OAuth2 config.
type Config struct {
	Name        string
	AuthURL     string
	TokenURL    string
	UserInfoURL string
	//
	RefershTokenURL string
	//
	LogoutURL string
	//
	RegisterURL string
	// callback url = server url + callback path, example: https://example.com/login/callback
	RedirectURI string
	Scope       string
	//
	ClientID     string
	ClientSecret string

	//
	ClientIDAttributeName     string
	ClientSecretAttributeName string
	RedirectURIAttributeName  string
	ResponseTypeAttributeName string
	ScopeAttributeName        string
	StateAttributeName        string

	// Token.access_token, default: access_token
	AccessTokenAttributeName string
	// Token.refresh_token, default: refresh_token
	RefreshTokenAttributeName string
	// Token.expires_in, default: expires_in
	ExpiresInAttributeName string
	// Token.id_token, default: id_token
	TokenTypeAttributeName string

	// User.username, default: username
	UsernameAttributeName string
	// User.email, default: email
	EmailAttributeName string
	// User.id, default: id
	IDAttributeName string
	// User.nickname, default: nickname
	NicknameAttributeName string
	// User.avatar, default: avatar
	AvatarAttributeName string
	// User.homepage, default: homepage
	HomepageAttributeName string
	// User.permissions, default: permissions
	PermissionsAttributeName string
	// User.groups, default: groups
	GroupsAttributeName string

	// url: login(authorize) + logout
	GetLoginURL    func(cfg *Config, state string) string
	GetLogoutURL   func(cfg *Config) string
	GetRegisterURL func(cfg *Config) string

	// token
	GetAccessTokenResponse func(cfg *Config, code string, state string) (*fetch.Response, error)
	// user
	GetUserResponse func(cfg *Config, token *Token, code string) (*fetch.Response, error)
	//
	RefreshToken func(cfg *Config, refreshToken string) (*fetch.Response, error)

	// base url for identity providers, such as auth0, authing
	BaseURL string
}

// generateLoginURL gets the authorize url.
//
// Example: https://login.example.com/authorize?client_id=CLIENT_ID&redirect_uri=https%3A%2F%2Fabc.com%2Flogin%2Fcallback&response_type=code&scope=openid&state=anything
func (oac *Config) generateLoginURL(state string) string {
	if state == "" {
		state = "anything"
	}

	if oac.GetLoginURL != nil {
		return oac.GetLoginURL(oac, state)
	}

	clientID := oac.ClientID
	redirectURI := oac.RedirectURI // oac.ServerUrl + "/login/callback"
	responseType := "code"
	scope := oac.Scope

	if scope == "" {
		scope = "openid"
	}

	return strings.Join([]string{
		oac.AuthURL,
		fmt.Sprintf("?%s=", oac.ClientIDAttributeName), clientID,
		fmt.Sprintf("&%s=", oac.RedirectURIAttributeName), url.QueryEscape(redirectURI),
		fmt.Sprintf("&%s=", oac.ResponseTypeAttributeName), responseType,
		fmt.Sprintf("&%s=", oac.ScopeAttributeName), url.QueryEscape(scope),
		fmt.Sprintf("&%s=", oac.StateAttributeName), url.QueryEscape(state),
	}, "")
}

// generateLogoutURL gets the logout url.
//
// Example: https://login.example.com/logout?client_id=CLIENT_ID&redirect_uri=https%3A%2F%2Fabc.com%2Flogin/callback
func (oac *Config) generateLogoutURL() string {
	if oac.GetLogoutURL != nil {
		return oac.GetLogoutURL(oac)
	}

	clientID := oac.ClientID
	redirectURI := oac.RedirectURI // oac.ServerUrl + "/login/callback"

	if oac.LogoutURL == "" {
		return ""
	}

	return strings.Join([]string{
		oac.LogoutURL,
		fmt.Sprintf("?%s=", oac.ClientIDAttributeName), clientID,
		fmt.Sprintf("&%s=", oac.RedirectURIAttributeName), url.QueryEscape(redirectURI),
	}, "")
}

// generateRegisterURL gets the register url.
//
// Exmaple: https://login.example.com/register?client_id=CLIENT_ID&invitation_code=xxxx
func (oac *Config) generateRegisterURL() string {
	if oac.GetRegisterURL != nil {
		return oac.GetRegisterURL(oac)
	}

	// @TODO
	if oac.RegisterURL == "" {
		message := url.QueryEscape(fmt.Sprintf("oauth2 %s does not support register", oac.Name))
		return fmt.Sprintf("/error?code=%d&message=%s", 500, message)
	}

	return strings.Join([]string{
		oac.RegisterURL,
		fmt.Sprintf("?%s=", oac.ClientIDAttributeName), oac.ClientID,
	}, "")
}

// ApplyDefaultConfig applies the default config.
func ApplyDefaultConfig(config *Config) (err error) {
	if config.ClientIDAttributeName == "" {
		config.ClientIDAttributeName = "client_id"
	}

	if config.ClientSecretAttributeName == "" {
		config.ClientSecretAttributeName = "client_secret"
	}

	if config.RedirectURIAttributeName == "" {
		config.RedirectURIAttributeName = "redirect_uri"
	}

	if config.ResponseTypeAttributeName == "" {
		config.ResponseTypeAttributeName = "response_type"
	}

	if config.ScopeAttributeName == "" {
		config.ScopeAttributeName = "scope"
	}

	if config.StateAttributeName == "" {
		config.StateAttributeName = "state"
	}

	if config.AccessTokenAttributeName == "" {
		config.AccessTokenAttributeName = "access_token"
	}

	if config.RefreshTokenAttributeName == "" {
		config.RefreshTokenAttributeName = "refresh_token"
	}

	if config.ExpiresInAttributeName == "" {
		config.ExpiresInAttributeName = "expires_in"
	}

	if config.TokenTypeAttributeName == "" {
		config.TokenTypeAttributeName = "token_type"
	}

	if config.UsernameAttributeName == "" {
		config.UsernameAttributeName = "username"
	}

	if config.EmailAttributeName == "" {
		config.EmailAttributeName = "email"
	}

	if config.IDAttributeName == "" {
		config.IDAttributeName = "id"
	}

	if config.NicknameAttributeName == "" {
		config.NicknameAttributeName = "nickname"
	}

	if config.AvatarAttributeName == "" {
		config.AvatarAttributeName = "avatar"
	}

	if config.HomepageAttributeName == "" {
		config.HomepageAttributeName = "homepage"
	}

	if config.PermissionsAttributeName == "" {
		config.PermissionsAttributeName = "permissions"
	}

	if config.GroupsAttributeName == "" {
		config.GroupsAttributeName = "groups"
	}

	return
}

// ValidateConfig validates the config.
func ValidateConfig(config *Config) error {
	if config.AuthURL == "" {
		panic(ErrConfigAuthURLEmpty)
	}

	if config.TokenURL == "" {
		panic(ErrConfigTokenURLEmpty)
	}

	if config.UserInfoURL == "" {
		panic(ErrConfigUserInfoURLEmpty)
	}

	if config.RedirectURI == "" {
		panic(ErrConfigRedirectURIEmpty)
	}

	if config.ClientID == "" {
		panic(ErrConfigClientIDEmpty)
	}

	if config.ClientSecret == "" {
		panic(ErrConfigClientSecretEmpty)
	}

	return nil
}

// ErrConfigAuthURLEmpty is the error of AuthURL is empty.
var ErrConfigAuthURLEmpty = errors.New("oauth2: config auth url is empty")

// ErrConfigTokenURLEmpty is the error of TokenURL is empty.
var ErrConfigTokenURLEmpty = errors.New("oauth2: config token url is empty")

// ErrConfigUserInfoURLEmpty is the error of UserInfoURL is empty.
var ErrConfigUserInfoURLEmpty = errors.New("oauth2: config user info url is empty")

// ErrConfigRedirectURIEmpty is the error of RedirectURI is empty.
var ErrConfigRedirectURIEmpty = errors.New("oauth2: config redirect uri is empty")

// ErrConfigClientIDEmpty is the error of ClientID is empty.
var ErrConfigClientIDEmpty = errors.New("oauth2: config client id is empty")

// ErrConfigClientSecretEmpty is the error of ClientSecret is empty.
var ErrConfigClientSecretEmpty = errors.New("oauth2: config client secret is empty")

// var ErrConfigScopeEmpty = errors.New("oauth2: config scope is empty")
