package chaoxing

import (
	"encoding/json"
	"fmt"
	"sign-your-horse/common"
	"sign-your-horse/provider"
	"sign-your-horse/users"
	"sign-your-horse/users/chaoxing"
	"strings"

	"github.com/imroc/req"
)

type ChaoxingProvider struct {
	Users []string `json:"users"`

	Alias               string                     `json:"-"`
	PushMessageCallback func(string, string) error `json:"-"`
}

func (c *ChaoxingProvider) Init(alias string, configBytes json.RawMessage) error {
	c.Alias = alias

	// check user exist and user type
	for _, user := range c.Users {
		found := users.GetUserInstance(user)
		if found == nil {
			return common.Raise(fmt.Sprintf("user %s not found", user))
		}
		if found.Type() != chaoxing.USER_TYPE {
			return common.Raise(fmt.Sprintf("user %s must be %s", user, chaoxing.USER_TYPE))
		}
	}

	return json.Unmarshal(configBytes, c)
}

func (c *ChaoxingProvider) PushMessageWithAlias(msg string) error {
	return c.PushMessageCallback(c.Alias, msg)
}

func (c *ChaoxingProvider) Task(user *chaoxing.ChaoxingUser, params map[string]string) {
	taskID := params["aid"]
	r := req.New()
	resp, err := r.Get(
		fmt.Sprintf(
			"https://mobilelearn.chaoxing.com/pptSign/stuSignajax?name=&activeId=%s&uid=%s&clientip=&useragent=&latitude=-1&longitude=-1&fid=0&appType=15&enc=%s",
			taskID,
			user.UserID,
			params["enc"],
		),
		req.Header{
			"Cookie":     user.Cookie,
			"User-Agent": user.UserAgent,
		},
	)
	if err != nil {
		c.PushMessageWithAlias("task " + string(taskID) + " sign in failed: " + err.Error())
	} else {
		c.PushMessageWithAlias("task " + string(taskID) + " sign in result: " + resp.String())
	}
}

func (c *ChaoxingProvider) Run(pushMessage func(string, string) error) {
	c.PushMessageCallback = pushMessage
}

func (c *ChaoxingProvider) Push(QRMessage string) {
	if !strings.HasPrefix(QRMessage, "SIGNIN:") {
		return
	}
	c.PushMessageWithAlias("new URL submitted: " + QRMessage)
	kvSlice := strings.Split(strings.TrimPrefix(QRMessage, "SIGNIN:"), "&")
	params := make(map[string]string)
	for _, kv := range kvSlice {
		e := strings.Split(kv, "=")
		if len(e) != 2 {
			continue
		}
		params[e[0]] = e[1]
	}
	for _, user := range c.Users {
		c.Task(users.GetUserInstance(user).(*chaoxing.ChaoxingUser), params)
	}
}

func init() {
	provider.RegisterProvider("chaoxing_cloud", &ChaoxingProvider{
		Users: []string{
			"chaoxing_user_sample",
		},
	})
}
