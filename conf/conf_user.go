package conf

import (
	"encoding/json"
	"sign-your-horse/common"
	"sign-your-horse/users"
)

type UserConfigBlock struct {
	Name     string          `json:"name"`
	UserType string          `json:"user_type"`
	Config   json.RawMessage `json:"config"`
}

func DumpUserConfigBlock() []UserConfigBlock {
	userTypeNameList, userInstanceList := users.GetAllUser()
	var ret []UserConfigBlock
	for i, userTypeName := range userTypeNameList {
		confJson, _ := json.Marshal(userInstanceList[i])
		ret = append(ret, UserConfigBlock{
			Name:     userTypeName + "_sample",
			UserType: userTypeName,
			Config:   confJson,
		})
	}
	return ret
}

func UpdateUserConfig(config *Config) {
	for _, u := range config.User {
		err := users.CreateUserWithConfig(u.UserType, u.Name, u.Config)
		common.Must(err)
	}
}
