package chaoxing

import "sign-your-horse/users"

var USER_TYPE = "chaoxing_user"

type ChaoxingUser struct {
	Cookie    string `json:"cookie"`
	UserAgent string `json:"useragent"`
	UserID    string `json:"uid"`
}

func (c ChaoxingUser) Type() string {
	return USER_TYPE
}

func init() {
	users.RegisterUser(USER_TYPE, &ChaoxingUser{})
}
