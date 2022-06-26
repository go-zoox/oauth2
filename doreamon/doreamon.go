package doreamon

import "github.com/go-zoox/oauth2"

func New(clientID, clientSecret, redirectURI string) (*oauth2.Client, error) {
	config := oauth2.Config{
		Name:         "哆啦A梦",
		AuthURL:      "https://login.zcorky.com/authorize",
		TokenURL:     "https://login.zcorky.com/token",
		UserInfoURL:  "https://login.zcorky.com/user",
		LogoutURL:    "https://login.zcorky.com/logout",
		Scope:        "openid email profile",
		RedirectURI:  redirectURI,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		//
		AccessTokenAttributeName:  "access_token",
		RefreshTokenAttributeName: "refresh_token",
		ExpiresInAttributeName:    "expires_in",
		TokenTypeAttributeName:    "token_type",
		EmailAttributeName:        "email",
		IDAttributeName:           "email",
		NicknameAttributeName:     "nickname",
		AvatarAttributeName:       "avatar",
		PermissionsAttributeName:  "permissions",
		GroupsAttributeName:       "groups",
	}

	return oauth2.New(config)
}
