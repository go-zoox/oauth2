package oauth2

import (
	"errors"
)

type Client struct {
	Config
	StepCallback
}

var logger Logger = &DefaultLogger{}

func New(config Config, options ...interface{}) (*Client, error) {
	for _, op := range options {
		switch op.(type) {
		case Logger:
			logger = op.(Logger)
		}
	}

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

// => authorize
func (oa *Client) Authorize(state string, callback func(loginUrl string)) {
	callback(oa.GetLoginUrl(state))
}

// <= callback
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

	user, err := oa.GetUser(&oa.Config, token)
	if err != nil {
		cb(nil, token, err)
		return
	}

	cb(user, token, nil)
}

// logout
func (oa *Client) Logout(state string, callback func(loginUrl string)) {
	callback(oa.GetLogoutUrl())
}
