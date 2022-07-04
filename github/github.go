package github

import "github.com/go-zoox/oauth2"

func New(clientID, clientSecret, redirectURI string) (*oauth2.Client, error) {
	config := oauth2.Config{
		Name:         "GitHub",
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		UserInfoURL:  "https://api.github.com/user",
		LogoutURL:    "https://github.com/logout",
		Scope:        "user",
		RedirectURI:  redirectURI,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		//
		AccessTokenAttributeName:  "access_token",
		RefreshTokenAttributeName: "refresh_token",
		ExpiresInAttributeName:    "expires_in",
		TokenTypeAttributeName:    "token_type",
		//
		EmailAttributeName:    "email",
		IDAttributeName:       "login",
		NicknameAttributeName: "name",
		AvatarAttributeName:   "avatar_url",
		HomepageAttributeName: "html_url",
	}

	return oauth2.New(config)
}
