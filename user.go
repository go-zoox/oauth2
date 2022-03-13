package oauth2

import (
	"errors"
	"io/ioutil"
	"net/http"

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

	client := &http.Client{}
	req, err := http.NewRequest("GET", oauth2_provider_userinfo_url, nil)
	if err != nil {
		return nil, errors.New("get user info error(1): " + err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("get user info error(2): " + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("get user info error(3): " + err.Error())
	}

	logger.Info("[getUser]:", string(body))

	error_code := gjson.Get(string(body), "code").Int()
	error_message := gjson.Get(string(body), "message").String()
	if error_code != 0 {
		return nil, errors.New("get user info error(4): " + error_message)
	}

	oauth2_email_attribute_name := config.EmailAttributeName
	oauth2_id_attribute_name := config.IdAttributeName
	oauth2_nickname_attribute_name := config.NicknameAttributeName
	oauth2_avatar_attribute_name := config.AvatarAttributeName
	oauth2_permissions_attribute_name := config.PermissionsAttributeName
	oauth2_groups_attribute_name := config.GroupsAttributeName

	user.Id = gjson.Get(string(body), oauth2_email_attribute_name).String()
	user.Email = gjson.Get(string(body), oauth2_id_attribute_name).String()
	user.Nickname = gjson.Get(string(body), oauth2_nickname_attribute_name).String()
	user.Avatar = gjson.Get(string(body), oauth2_avatar_attribute_name).String()
	user.Permissions = make([]string, 0)

	permissionsResult := gjson.Get(string(body), oauth2_permissions_attribute_name)
	permissionsResult.ForEach(func(key, value gjson.Result) bool {
		user.Permissions = append(user.Permissions, value.String())
		return true
	})

	groupsResult := gjson.Get(string(body), oauth2_groups_attribute_name)
	groupsResult.ForEach(func(key, value gjson.Result) bool {
		user.Groups = append(user.Groups, value.String())
		return true
	})

	return user, nil
}
