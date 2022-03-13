package oauth2

import (
	"errors"
)

type OAuth2 struct {
	Config
	StepCallback
}

var logger Logger = &DefaultLogger{}

func New(config Config, options ...interface{}) (*OAuth2, error) {
	for _, op := range options {
		switch op.(type) {
		case Logger:
			op.(Logger).Info("resign logger")
			logger = op.(Logger)
		}
	}

	if err := ValidateConfig(&config); err != nil {
		return nil, err
	}

	if err := ApplyDefaultConfig(&config); err != nil {
		return nil, err
	}

	return &OAuth2{
		Config: config,
	}, nil
}

// => authorize
func (oa *OAuth2) Authorize(state string, callback func(loginUrl string)) {
	callback(oa.GetLoginUrl(state))
}

// <= callback
func (oa *OAuth2) Callback(code, state string, cb func(user *User, token *Token, err error)) {
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

// parts
