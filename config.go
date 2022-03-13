package oauth2

import (
	"errors"
	"net/url"
	"strings"
)

type Config struct {
	Name        string
	AuthUrl     string
	TokenUrl    string
	UserInfoUrl string
	LogoutUrl   string
	// callback url = server url + callback path, example: https://example.com/login/callback
	RedirectUri string
	Scope       string
	//
	ClientId     string
	ClientSecret string

	// Token.access_token, default: access_token
	AccessTokenAttributeName string
	// Token.refresh_token, default: refresh_token
	RefreshTokenAttributeName string
	// Token.expires_in, default: expires_in
	ExpiresInAttributeName string
	// Token.id_token, default: id_token
	TokenTypeAttributeName string

	// User.email, default: email
	EmailAttributeName string
	// User.id, default: id
	IdAttributeName string
	// User.nickname, default: nickname
	NicknameAttributeName string
	// User.avatar, default: avatar
	AvatarAttributeName string
	// User.permissions, default: permissions
	PermissionsAttributeName string
	// User.groups, default: groups
	GroupsAttributeName string
}

// Get The Authorize Url
//
// Example: https://login.example.com/authorize?client_id=CLIENT_ID&redirect_uri=https%3A%2F%2Fabc.com%2Flogin%2Fcallback&response_type=code&scope=openid&state=anything
func (oac *Config) GetLoginUrl(state string) string {
	if state == "" {
		state = "anything"
	}

	clientId := oac.ClientId
	redirectUri := oac.RedirectUri // oac.ServerUrl + "/login/callback"
	responseType := "code"
	scope := oac.Scope

	if scope == "" {
		scope = "openid"
	}

	return strings.Join([]string{
		oac.AuthUrl,
		"?client_id=", clientId,
		"&redirect_uri=", url.QueryEscape(redirectUri),
		"&response_type=", responseType,
		"&scope=", url.QueryEscape(scope),
		"&state=", url.QueryEscape(state),
	}, "")
}

// Get The Logout Url
//
// Example: https://login.example.com/logout?client_id=CLIENT_ID&redirect_uri=https%3A%2F%2Fabc.com%2Flogin/callback
func (oac *Config) GetLogoutUrl() string {
	clientId := oac.ClientId
	redirectUri := oac.RedirectUri // oac.ServerUrl + "/login/callback"

	return strings.Join([]string{
		oac.LogoutUrl,
		"?client_id=", clientId,
		"&redirect_uri=", url.QueryEscape(redirectUri),
	}, "")
}

func ApplyDefaultConfig(config *Config) (err error) {
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

	if config.EmailAttributeName == "" {
		config.EmailAttributeName = "email"
	}

	if config.IdAttributeName == "" {
		config.IdAttributeName = "id"
	}

	if config.NicknameAttributeName == "" {
		config.NicknameAttributeName = "nickname"
	}

	if config.AvatarAttributeName == "" {
		config.AvatarAttributeName = "avatar"
	}

	if config.PermissionsAttributeName == "" {
		config.PermissionsAttributeName = "permissions"
	}

	if config.GroupsAttributeName == "" {
		config.GroupsAttributeName = "groups"
	}

	return
}

func ValidateConfig(config *Config) error {
	if config.AuthUrl == "" {
		panic(ErrConfigAuthUrlEmpty)
	}

	if config.TokenUrl == "" {
		panic(ErrConfigTokenUrlEmpty)
	}

	if config.UserInfoUrl == "" {
		panic(ErrConfigUserInfoUrlEmpty)
	}

	if config.RedirectUri == "" {
		panic(ErrConfigRedirectUriEmpty)
	}

	if config.ClientId == "" {
		panic(ErrConfigClientIdEmpty)
	}

	if config.ClientSecret == "" {
		panic(ErrConfigClientSecretEmpty)
	}

	return nil
}

var ErrConfigAuthUrlEmpty = errors.New("oauth2: config auth url is empty")
var ErrConfigTokenUrlEmpty = errors.New("oauth2: config token url is empty")
var ErrConfigUserInfoUrlEmpty = errors.New("oauth2: config user info url is empty")
var ErrConfigRedirectUriEmpty = errors.New("oauth2: config redirect uri is empty")
var ErrConfigClientIdEmpty = errors.New("oauth2: config client id is empty")
var ErrConfigClientSecretEmpty = errors.New("oauth2: config client secret is empty")

// var ErrConfigScopeEmpty = errors.New("oauth2: config scope is empty")
