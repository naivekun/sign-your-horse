package chaoxing

import (
	"encoding/json"
	"fmt"
	"sign-your-horse/common"
	"sign-your-horse/provider"
	"sign-your-horse/users"
	"sign-your-horse/users/chaoxing"
	"time"
)

type ActiveTime struct {
	Weekday  int    `json:"weekday"`
	Time     string `json:"time"`
	Duration int    `json:"duration"`
}

type ChaoxingProvider struct {
	Users        []string     `json:"users"`
	TaskInterval int          `json:"interval"`
	TaskTime     []ActiveTime `json:"tasktime"`
	Verbose      bool         `json:"verbose"`

	Alias               string                     `json:"-"`
	PushMessageCallback func(string, string) error `json:"-"`
}

func (c *ChaoxingProvider) Init(alias string, configBytes json.RawMessage) error {
	c.Alias = alias
	ret := json.Unmarshal(configBytes, c)
	if ret != nil {
		return ret
	}

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

	if c.TaskTime == nil {
		common.LogWithModule(alias, "no tasktime specified, module %s will work at all time", alias)
	} else {
		for i, activeTime := range c.TaskTime {
			if !checkActiveTime(&activeTime) {
				return common.Raise(fmt.Sprintf("invalid date format in tasktime entry #%d", i))
			}
		}
	}
	common.LogWithModule(alias, "Local time is "+time.Now().String()+". Check your time and timezone carefully!")
	return nil
}

func checkActiveTime(a *ActiveTime) bool {
	_, err := time.Parse("15:04", a.Time)
	if err != nil {
		return false
	}
	if a.Duration > 60*24 || a.Duration < 0 {
		return false
	}
	if a.Weekday < 0 || a.Weekday > 6 {
		return false
	}
	return true
}

func parseActiveTime(a *ActiveTime) time.Time {
	ft, _ := time.Parse("15:04", a.Time)
	t := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), ft.Hour(), ft.Minute(), time.Now().Second(), 0, time.Now().Location())
	return t
}

func (c *ChaoxingProvider) Run(pushMessage func(string, string) error) {
	c.PushMessageCallback = pushMessage
	for {
		isAtTaskTime := false
		for _, activeTime := range c.TaskTime {
			if int(time.Now().Weekday()) == activeTime.Weekday {
				if activeTime.Duration == 0 {
					isAtTaskTime = true
					break
				}
				t := parseActiveTime(&activeTime)
				if t.Before(time.Now()) && t.Add(time.Minute*time.Duration(activeTime.Duration)).After(time.Now()) {
					isAtTaskTime = true
					break
				}
			}

		}
		if isAtTaskTime || c.TaskTime == nil {
			for _, user := range c.Users {
				c.Task(users.GetUserInstance(user).(*chaoxing.ChaoxingUser))
			}
		} else {
			if c.Verbose {
				common.LogWithModule(c.Alias, "no task to do at %s because it is not task time", time.Now().String())
			}
		}
		time.Sleep(time.Duration(c.TaskInterval) * time.Second)
	}
}

func (c *ChaoxingProvider) Push(_ string) {}

func init() {
	provider.RegisterProvider("chaoxing", &ChaoxingProvider{
		Users: []string{
			"chaoxing_user_sample",
		},
		TaskInterval: 5,
		TaskTime: []ActiveTime{
			{
				Weekday:  int(time.Monday),
				Time:     "07:50",
				Duration: 20,
			},
			{
				Weekday:  int(time.Thursday),
				Time:     "13:50",
				Duration: 20,
			},
		},
		Verbose: true,
	})
}
