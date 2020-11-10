package teachermate

import (
	"encoding/json"
	"log"
	"sign-your-horse/provider"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
)

type TeacherMateProvider struct {
	ServerAddr string `json:"server"`
	UseHTTPS   bool   `json:"usehttps"`
	Alias      string `json:"-"`
	ServerCert string `json:"srvcert"`
	ServerKey  string `json:"srvkey"`
}

var pushMessage func(string, string) error

func (t *TeacherMateProvider) Init(alias string, configBytes json.RawMessage) error {
	t.Alias = alias
	return json.Unmarshal(configBytes, t)
}

func (t *TeacherMateProvider) Run(pushMessage_ func(string, string) error) {
	box := packr.NewBox("static")
	pushMessage = pushMessage_
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.POST("/url/add", add)
	server.GET("/url/raw", raw)
	server.GET("/url/redirect", redirect)
	server.GET("/url/", urlinfo)
	server.StaticFS("/static", box)
	log.Println("server is listening at " + t.ServerAddr)
	if t.UseHTTPS {
		server.RunTLS(t.ServerAddr, t.ServerCert, t.ServerKey)
	} else {
		server.Run(t.ServerAddr)
	}
}

func init() {
	provider.RegisterProvider("teachermate", &TeacherMateProvider{
		ServerAddr: "0.0.0.0:3000",
		UseHTTPS:   false,
	})
}
