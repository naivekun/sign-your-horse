package chaoxing

import (
	"encoding/json"
	"fmt"
	"sign-your-horse/provider"
	"strings"

	"github.com/imroc/req"
)

type ChaoxingProvider struct {
	Alias     string `json:"-"`
	Cookie    string `json:"cookie"`
	UserAgent string `json:"useragent"`
	UserID    string `json:"uid"`
	CourseID  string `json:"courseid"`
	ClassID   string `json:"classid"`

	PushMessageCallback func(string, string) error `json:"-"`
}

func (c *ChaoxingProvider) Init(alias string, configBytes json.RawMessage) error {
	c.Alias = alias
	return json.Unmarshal(configBytes, c)
}

func (c *ChaoxingProvider) PushMessageWithAlias(msg string) error {
	return c.PushMessageCallback(c.Alias, msg)
}

func (c *ChaoxingProvider) Task(params map[string]string) {
	taskID := params["aid"]
	r := req.New()
	resp, err := r.Get(
		fmt.Sprintf(
			"https://mobilelearn.chaoxing.com/pptSign/stuSignajax?name=&activeId=%s&uid=%s&clientip=&useragent=&latitude=-1&longitude=-1&fid=0&appType=15&enc=%s",
			taskID,
			c.UserID,
			params["enc"],
		),
		req.Header{
			"Cookie":     c.Cookie,
			"User-Agent": c.UserAgent,
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
	kvSlice := strings.Split(strings.TrimPrefix(QRMessage, "SIGNIN:"), "&")
	params := make(map[string]string)
	for _, kv := range kvSlice {
		e := strings.Split(kv, "=")
		if len(e) != 2 {
			continue
		}
		params[e[0]] = e[1]
	}
	c.Task(params)
}

func init() {
	provider.RegisterProvider("chaoxing_cloud", &ChaoxingProvider{})
}
