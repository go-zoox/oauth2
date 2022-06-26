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
}

// GetUser gets the user by token.
func GetUser(config *Config, token *Token) (*User, error) {
	user := &User{}

	// oauth2ProviderUserinfoURL := "https://httpbin.zcorky.com/get"
	oauth2ProviderUserinfoURL := config.UserInfoURL

	response, err := fetch.Get(oauth2ProviderUserinfoURL, &fetch.Config{
		Headers: map[string]string{
			"Authorization": "Bearer " + token.AccessToken,
		},
	})
	if err != nil {
		return nil, errors.New("get user info error: " + err.Error())
	}

	logger.Info("[oauth2][getUser]: %s", response.String())

	errorCode := response.Get("code").Int()
	errorMessage := response.Get("message").String()
	if errorCode != 0 {
		return nil, errors.New("get user info error(4): " + errorMessage)
	}

	oauth2EmailAttributeName := config.EmailAttributeName
	oauth2IDAttributeName := config.IDAttributeName
	oauth2NicknameAttributeName := config.NicknameAttributeName
	oauth2AvatarAttributeName := config.AvatarAttributeName
	oauth2PermissionsAttributeName := config.PermissionsAttributeName
	oauth2GroupsAttributeName := config.GroupsAttributeName

	user.ID = response.Get(oauth2EmailAttributeName).String()
	user.Email = response.Get(oauth2IDAttributeName).String()
	user.Nickname = response.Get(oauth2NicknameAttributeName).String()
	user.Avatar = response.Get(oauth2AvatarAttributeName).String()
	user.Permissions = make([]string, 0)

	permissionsResult := response.Get(oauth2PermissionsAttributeName)
	permissionsResult.ForEach(func(key, value gjson.Result) bool {
		user.Permissions = append(user.Permissions, value.String())
		return true
	})

	groupsResult := response.Get(oauth2GroupsAttributeName)
	groupsResult.ForEach(func(key, value gjson.Result) bool {
		user.Groups = append(user.Groups, value.String())
		return true
	})

	return user, nil
}
