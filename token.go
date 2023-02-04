package oauth2

import (
	"errors"

	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
)

// Token is the oauth2 token.
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	//
	Raw *fetch.Response `json:"response"`
}

// GetToken gets the token by code and state.
func GetToken(config *Config, code string, state string) (*Token, error) {
	token := &Token{}

	oauth2ProviderTokenURL := config.TokenURL
	oauth2ClientID := config.ClientID
	oauth2ClientSecret := config.ClientSecret
	oauth2RedirectURI := config.RedirectURI
	//
	oauth2AccessTokenAttributeName := config.AccessTokenAttributeName
	oauth2RefreshTokenAttributeName := config.RefreshTokenAttributeName
	oauth2ExpiresInAttributeName := config.ExpiresInAttributeName
	oauth2TokenTypeAttributeName := config.TokenTypeAttributeName

	var response *fetch.Response
	var err error
	if config.GetAccessTokenResponse != nil {
		response, err = config.GetAccessTokenResponse(config, code, state)
	} else {
		response, err = fetch.Post(oauth2ProviderTokenURL, &fetch.Config{
			Headers: map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
				"Accept":       "application/json",
			},
			Body: map[string]string{
				"client_id":     oauth2ClientID,
				"client_secret": oauth2ClientSecret,
				"grant_type":    "authorization_code",
				"redirect_uri":  oauth2RedirectURI,
				"code":          code,
				"state":         state,
			},
		})
	}
	if err != nil {
		return nil, errors.New("get access token error by code (3): " + err.Error())
	}

	logger.Info("[oauth2][getToken]: %s", response.String())

	errorCode := response.Get("code").Int()
	errorMessage := response.Get("message").String()
	if errorCode == 5003002 {
		return nil, errors.New("code is expired: " + errorMessage)
	} else if errorCode != 0 {
		return nil, errors.New("get access token error by code (3): " + err.Error())
	}

	//
	accessToken := response.Get(oauth2AccessTokenAttributeName).String()
	refreshToken := response.Get(oauth2RefreshTokenAttributeName).String()
	expiresIn := response.Get(oauth2ExpiresInAttributeName).Int()
	tokenType := response.Get(oauth2TokenTypeAttributeName).String()

	token.AccessToken = accessToken
	token.RefreshToken = refreshToken
	token.ExpiresIn = expiresIn
	token.TokenType = tokenType

	token.Raw = response

	return token, nil
}
