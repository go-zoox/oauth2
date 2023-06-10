package oauth2

import (
	"errors"

	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
	"github.com/tidwall/gjson"
)

// User is the oauth2 user.
type User struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	Avatar      string   `json:"avatar"`
	Nickname    string   `json:"nickname"`
	Groups      []string `json:"groups"`
	Permissions []string `json:"permissions"`

	raw *fetch.Response
}

// Raw gets raw data with *fetch.Response.
func (u *User) Raw() *fetch.Response {
	return u.raw
}

// GetUser gets the user by token.
func GetUser(config *Config, token *Token, code string) (*User, error) {
	user := &User{}

	var response *fetch.Response
	var err error
	if config.GetUserResponse != nil {
		response, err = config.GetUserResponse(config, token, code)
	} else {
		response, err = fetch.Get(config.UserInfoURL, &fetch.Config{
			Headers: map[string]string{
				"Authorization": "Bearer " + token.AccessToken,
			},
		})
	}
	if err != nil {
		return nil, errors.New("get user info error: " + err.Error())
	}

	logger.Debugf("[oauth2][user]: %s", response.String())

	errorCode := response.Get("code").Int()
	errorMessage := response.Get("message").String()
	if errorCode != 0 {
		return nil, errors.New("get user info error(4): " + errorMessage)
	}

	user.ID = response.Get(config.IDAttributeName).String()
	user.Email = response.Get(config.EmailAttributeName).String()
	user.Nickname = response.Get(config.NicknameAttributeName).String()
	user.Avatar = response.Get(config.AvatarAttributeName).String()
	user.Permissions = make([]string, 0)

	permissionsResult := response.Get(config.PermissionsAttributeName)
	permissionsResult.ForEach(func(key, value gjson.Result) bool {
		user.Permissions = append(user.Permissions, value.String())
		return true
	})

	groupsResult := response.Get(config.GroupsAttributeName)
	groupsResult.ForEach(func(key, value gjson.Result) bool {
		user.Groups = append(user.Groups, value.String())
		return true
	})

	user.raw = response

	return user, nil
}
