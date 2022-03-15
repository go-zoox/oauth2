package oauth2

import (
	"errors"

	"github.com/go-zoox/fetch"
	"github.com/tidwall/gjson"
)

type User struct {
	Id          string   `json:"id"`
	Email       string   `json:"email"`
	Avatar      string   `json:"avatar"`
	Nickname    string   `json:"nickname"`
	Groups      []string `json:"groups"`
	Permissions []string `json:"permissions"`
}

func GetUser(config *Config, token *Token) (*User, error) {
	user := &User{}

	// oauth2_provider_userinfo_url := "https://httpbin.zcorky.com/get"
	oauth2_provider_userinfo_url := config.UserInfoUrl

	response, err := fetch.Get(oauth2_provider_userinfo_url, &fetch.Config{
		Headers: map[string]string{
			"Authorization": "Bearer " + token.AccessToken,
		},
	})
	if err != nil {
		return nil, errors.New("get user info error: " + err.Error())
	}

	logger.Info("[getUser]:", response.String())

	error_code := response.Get("code").Int()
	error_message := response.Get("message").String()
	if error_code != 0 {
		return nil, errors.New("get user info error(4): " + error_message)
	}

	oauth2_email_attribute_name := config.EmailAttributeName
	oauth2_id_attribute_name := config.IdAttributeName
	oauth2_nickname_attribute_name := config.NicknameAttributeName
	oauth2_avatar_attribute_name := config.AvatarAttributeName
	oauth2_permissions_attribute_name := config.PermissionsAttributeName
	oauth2_groups_attribute_name := config.GroupsAttributeName

	user.Id = response.Get(oauth2_email_attribute_name).String()
	user.Email = response.Get(oauth2_id_attribute_name).String()
	user.Nickname = response.Get(oauth2_nickname_attribute_name).String()
	user.Avatar = response.Get(oauth2_avatar_attribute_name).String()
	user.Permissions = make([]string, 0)

	permissionsResult := response.Get(oauth2_permissions_attribute_name)
	permissionsResult.ForEach(func(key, value gjson.Result) bool {
		user.Permissions = append(user.Permissions, value.String())
		return true
	})

	groupsResult := response.Get(oauth2_groups_attribute_name)
	groupsResult.ForEach(func(key, value gjson.Result) bool {
		user.Groups = append(user.Groups, value.String())
		return true
	})

	return user, nil
}
