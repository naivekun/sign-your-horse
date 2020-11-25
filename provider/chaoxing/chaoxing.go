package chaoxing

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sign-your-horse/common"
	"sign-your-horse/provider"
	"strings"
	"time"

	"github.com/imroc/req"
)

type ChaoxingProvider struct {
	Alias           string `json:"-"`
	Cookie          string `json:"cookie"`
	UserAgent       string `json:"useragent"`
	UserID          string `json:"uid"`
	CourseID        string `json:"courseid"`
	ClassID         string `json:"classid"`
	TaskInterval    int    `json:"interval"`
	ShowLoopMessage bool   `json:"verbose"`

	PushMessageCallback func(string, string) error `json:"-"`
}

func (c *ChaoxingProvider) Init(alias string, configBytes json.RawMessage) error {
	c.Alias = alias
	return json.Unmarshal(configBytes, c)
}

func (c *ChaoxingProvider) PushMessageWithAlias(msg string) error {
	return c.PushMessageCallback(c.Alias, msg)
}

func (c *ChaoxingProvider) Task() {
	extractTasksRegex := regexp.MustCompile(`activeDetail\(\d+,`)
	extrackTaskIDRegex := regexp.MustCompile(`\d+`)

	r := req.New()
	tasks, err := r.Get(
		fmt.Sprintf("https://mobilelearn.chaoxing.com/widget/pcpick/stu/index?courseId=%s&jclassId=%s",
			c.CourseID,
			c.ClassID),
		req.Header{
			"Cookie":     c.Cookie,
			"User-Agent": c.UserAgent,
		},
	)
	if err != nil {
		c.PushMessageWithAlias("get task list failed: " + err.Error())
	} else {
		taskListString := tasks.String()
		if len(taskListString) == 0 {
			c.PushMessageWithAlias("get task list failed: empty page")
			return
		}
		finishedSepIndex := strings.Index(taskListString, "已结束")
		if finishedSepIndex == -1 {
			c.PushMessageWithAlias("invalid task page, maybe you need login?")
			return
		}
		taskListString = taskListString[:finishedSepIndex]
		tasksString := extractTasksRegex.FindAll([]byte(taskListString), -1)
		if len(tasksString) == 0 && c.ShowLoopMessage {
			common.LogWithModule(c.Alias, " no task to do at "+time.Now().String())
		} else {
			for _, task := range tasksString {
				taskID := extrackTaskIDRegex.Find(task)
				resp, err := r.Get(
					fmt.Sprintf("https://mobilelearn.chaoxing.com/pptSign/stuSignajax?name=&activeId=%s&uid=%s&clientip=&useragent=&latitude=-1&longitude=-1&fid=0&appType=15", string(taskID), c.UserID),
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
		}
	}
}

func (c *ChaoxingProvider) Run(pushMessage func(string, string) error) {
	c.PushMessageCallback = pushMessage
	for {
		c.Task()
		time.Sleep(time.Duration(c.TaskInterval) * time.Second)
	}
}

func (c *ChaoxingProvider) Push(_ string) {}

func init() {
	provider.RegisterProvider("chaoxing", &ChaoxingProvider{
		TaskInterval:    5,
		ShowLoopMessage: true,
	})
}
