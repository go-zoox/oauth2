package oauth2

import "errors"

type StepCallback struct {
}

//
func (oa *StepCallback) GetToken(config *Config, code, state string) (*Token, error) {
	if len(code) == 0 || len(state) == 0 {
		return nil, errors.New("invalid oauth2 login callback, code or state are required")
	}

	token, err := GetToken(config, code, state)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (oa *StepCallback) GetUser(config *Config, token *Token) (*User, error) {
	user, err := GetUser(config, token)
	if err != nil {
		return nil, err
	}

	return user, nil
}
