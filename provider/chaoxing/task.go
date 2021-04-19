package chaoxing

import (
	"fmt"
	"regexp"
	"sign-your-horse/common"
	"sign-your-horse/users/chaoxing"
	"strings"
	"time"

	"github.com/imroc/req"
)

func (c *ChaoxingProvider) Task(user *chaoxing.ChaoxingUser) {
	extractTasksRegex := regexp.MustCompile(`activeDetail\(\d+,\d+,`)
	extrackTaskInfoRegex := regexp.MustCompile(`(?P<TaskID>\d+),(?P<TaskType>\d+)`)

	r := req.New()
	tasks, err := r.Get(
		fmt.Sprintf("https://mobilelearn.chaoxing.com/widget/pcpick/stu/index?courseId=%s&jclassId=%s",
			c.CourseID,
			c.ClassID),
		req.Header{
			"Cookie":     user.Cookie,
			"User-Agent": user.UserAgent,
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
		if len(tasksString) == 0 && c.Verbose {
			common.LogWithModule(c.Alias, " no task in list at %s", time.Now().String())
		} else {
			for _, task := range tasksString {
				taskInfo := extrackTaskInfoRegex.FindStringSubmatch(string(task))
				if len(taskInfo) != 2 {
					//something wrong
					continue
				}
				if taskInfo[1] != "2" {
					//not a sign-in event
					continue
				}
				taskID := taskInfo[0]
				resp, err := r.Get(
					fmt.Sprintf("https://mobilelearn.chaoxing.com/pptSign/stuSignajax?name=&activeId=%s&uid=%s&clientip=&useragent=&latitude=-1&longitude=-1&fid=0&appType=15", string(taskID), user.UserID),
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
		}
	}
}

func (c *ChaoxingProvider) PushMessageWithAlias(msg string) error {
	return c.PushMessageCallback(c.Alias, msg)
}
