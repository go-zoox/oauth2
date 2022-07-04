package oauth2

import (
	"errors"
)

// Client is the OAuth2 Client.
type Client struct {
	Config
	StepCallback
}

// New creates a OAuth2 Client.
func New(config Config, options ...interface{}) (*Client, error) {
	if err := ValidateConfig(&config); err != nil {
		return nil, err
	}

	if err := ApplyDefaultConfig(&config); err != nil {
		return nil, err
	}

	return &Client{
		Config: config,
	}, nil
}

// Authorize is the first step of login
// means redirect to oauth server authorize page
func (oa *Client) Authorize(state string, callback func(loginUrl string)) {
	callback(oa.GetLoginURL(state))
}

// Callback is the second step of login,
// means oauth server visit callback url with code.
// And we will get access_token and refresh_token with the code.
// Then we can use access_token to get user info.
func (oa *Client) Callback(code, state string, cb func(user *User, token *Token, err error)) {
	if len(code) == 0 || len(state) == 0 {
		cb(nil, nil, errors.New("invalid oauth2 login callback, code or state are required"))
		return
	}

	token, err := oa.GetToken(&oa.Config, code, state)
	if err != nil {
		cb(nil, nil, err)
		return
	}

	user, err := oa.GetUser(&oa.Config, token, code)
	if err != nil {
		cb(nil, token, err)
		return
	}

	cb(user, token, nil)
}

// Logout just to logout the user
func (oa *Client) Logout(callback func(logoutUrl string)) {
	callback(oa.GetLogoutURL())
}
