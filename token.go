package oauth2

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func GetToken(config *Config, code string, state string) (*Token, error) {
	token := &Token{}

	// oauth2_provider_token_url := "https://httpbin.zcorky.com/post"
	oauth2_provider_token_url := config.TokenUrl
	oauth2_client_id := config.ClientId
	oauth2_client_secret := config.ClientSecret
	oauth2_redirect_uri := config.RedirectUri
	//
	oauth2_access_token_attribute_name := config.AccessTokenAttributeName
	oauth2_refresh_token_attribute_name := config.RefreshTokenAttributeName
	oauth2_expires_in_attribute_name := config.ExpiresInAttributeName
	oauth2_token_type_attribute_name := config.TokenTypeAttributeName

	client := &http.Client{}
	req, err := http.NewRequest("POST", oauth2_provider_token_url, strings.NewReader(url.Values{
		"client_id":     {oauth2_client_id},
		"client_secret": {oauth2_client_secret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {oauth2_redirect_uri},
		"code":          {code},
		"state":         {state},
	}.Encode()))
	if err != nil {
		return nil, errors.New("get access token error by code (1): " + err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("get access token error by code (2): " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	logger.Info("[getToken]:", string(body))

	error_code := gjson.Get(string(body), "code").Int()
	error_message := gjson.Get(string(body), "message").String()
	if error_code == 5003002 {
		return nil, errors.New("code is expired: " + error_message)
	} else if error_code != 0 {
		return nil, errors.New("get access token error by code (3): " + err.Error())
	}

	//
	access_token := gjson.Get(string(body), oauth2_access_token_attribute_name).String()
	refresh_token := gjson.Get(string(body), oauth2_refresh_token_attribute_name).String()
	expires_in := gjson.Get(string(body), oauth2_expires_in_attribute_name).Int()
	token_type := gjson.Get(string(body), oauth2_token_type_attribute_name).String()

	token.AccessToken = access_token
	token.RefreshToken = refresh_token
	token.ExpiresIn = expires_in
	token.TokenType = token_type

	return token, nil
}
