package teachermate_cloud

import (
	"encoding/json"
	"log"
	"sign-your-horse/provider"
	"strings"
)

type TeacherMateProvider struct {
	Alias string `json:"-"`
}

var pushMessage func(string, string) error

func (t *TeacherMateProvider) Init(alias string, configBytes json.RawMessage) error {
	t.Alias = alias
	return json.Unmarshal(configBytes, t)
}

func (t *TeacherMateProvider) Run(pushMessage_ func(string, string) error) {
	pushMessage = pushMessage_
}

func (t *TeacherMateProvider) Push(msg string) {
	if strings.HasPrefix(msg, "http") {
		err := pushMessage(t.Alias, "Teachermate new URL submitted: "+msg)
		if err != nil {
			log.Println("pushMessage failed: ", err.Error())
		}
	}
}

func init() {
	provider.RegisterProvider("teachermate_cloud", &TeacherMateProvider{})
}
