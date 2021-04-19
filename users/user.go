package users

import (
	"encoding/json"
	"sign-your-horse/common"
)

type User interface {
	Type() string
}

var userMap = make(map[string]User)
var userInstanceMap = make(map[string]User)

func RegisterUser(name string, user User) {
	userMap[name] = user
}

func GetUser(name string) User {
	if user, found := userMap[name]; found {
		return user
	}
	return nil
}

func GetUserInstance(name string) User {
	if user, found := userInstanceMap[name]; found {
		return user
	}
	return nil
}

func GetAllUser() ([]string, []User) {
	var nameList []string
	var userList []User
	for name, user := range userMap {
		nameList = append(nameList, name)
		userList = append(userList, user)
	}
	return nameList, userList
}

func CreateUserWithConfig(userType string, alias string, config json.RawMessage) error {
	if _, found := userInstanceMap[alias]; found {
		return common.Raise("duplicate user " + alias)
	}
	user := GetUser(userType)
	if user == nil {
		return common.Raise("no user type named: " + userType)
	}
	userInstance := common.CloneEmpty(user).(User)
	err := json.Unmarshal(config, userInstance)
	if err != nil {
		return err
	}
	userInstanceMap[alias] = userInstance
	return nil
}
