package doreamon

import "github.com/go-zoox/oauth2"

func New(clientId, clientSecret, redirectUri string) (*oauth2.OAuth2, error) {
	config := oauth2.Config{
		Name:         "哆啦A梦",
		AuthUrl:      "https://login.zcorky.com/authorize",
		TokenUrl:     "https://login.zcorky.com/token",
		UserInfoUrl:  "https://login.zcorky.com/user",
		LogoutUrl:    "https://login.zcorky.com/logout",
		Scope:        "openid email profile",
		RedirectUri:  redirectUri,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		//
		AccessTokenAttributeName:  "access_token",
		RefreshTokenAttributeName: "refresh_token",
		ExpiresInAttributeName:    "expires_in",
		TokenTypeAttributeName:    "token_type",
		EmailAttributeName:        "email",
		IdAttributeName:           "email",
		NicknameAttributeName:     "nickname",
		AvatarAttributeName:       "avatar",
		PermissionsAttributeName:  "permissions",
		GroupsAttributeName:       "groups",
	}

	return oauth2.New(config)
}
