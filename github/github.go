package github

import "github.com/go-zoox/oauth2"

func New(clientId, clientSecret, redirectUri string) (*oauth2.Client, error) {
	config := oauth2.Config{
		Name:         "哆啦A梦",
		AuthUrl:      "https://github.com/login/oauth/authorize",
		TokenUrl:     "https://github.com/login/oauth/access_token",
		UserInfoUrl:  "https://api.github.com/user",
		LogoutUrl:    "https://github.com/logout",
		Scope:        "user",
		RedirectUri:  redirectUri,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		//
		AccessTokenAttributeName:  "access_token",
		RefreshTokenAttributeName: "refresh_token",
		ExpiresInAttributeName:    "expires_in",
		TokenTypeAttributeName:    "token_type",
		//
		EmailAttributeName:    "email",
		IdAttributeName:       "login",
		NicknameAttributeName: "name",
		AvatarAttributeName:   "avatar_url",
		HomepageAttributeName: "html_url",
	}

	return oauth2.New(config)
}
