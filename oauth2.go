package oauth2

import (
	"errors"
)

// Client is the oauth2 client interface.
type Client interface {
	Authorize(state string, callback func(loginUrl string))
	Callback(code, state string, cb func(user *User, token *Token, err error))
	Logout(callback func(logoutUrl string))
}

// client is the OAuth2 client.
type client struct {
	Config
	StepCallback
}

// New creates a OAuth2 client.
func New(config Config, options ...interface{}) (Client, error) {
	if err := ValidateConfig(&config); err != nil {
		return nil, err
	}

	if err := ApplyDefaultConfig(&config); err != nil {
		return nil, err
	}

	return &client{
		Config: config,
	}, nil
}

// Authorize is the first step of login
// means redirect to oauth server authorize page
func (oa *client) Authorize(state string, callback func(loginUrl string)) {
	callback(oa.generateLoginURL(state))
}

// Callback is the second step of login,
// means oauth server visit callback url with code.
// And we will get access_token and refresh_token with the code.
// Then we can use access_token to get user info.
func (oa *client) Callback(code, state string, cb func(user *User, token *Token, err error)) {
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
func (oa *client) Logout(callback func(logoutUrl string)) {
	callback(oa.generateLogoutURL())
}
