package cloudscan

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
)

var MessageChan chan string

func Init(config json.RawMessage) (*CloudScanAPIServer, error) {
	ret := &CloudScanAPIServer{}
	err := json.Unmarshal(config, ret)
	ret.APIMessageOutputChan = make(chan string)
	MessageChan = ret.APIMessageOutputChan
	return ret, err
}

func (t *CloudScanAPIServer) Run() {
	box := packr.NewBox("static")
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.POST("/url/add", add)
	server.GET("/url/raw", raw)
	server.GET("/url/redirect", redirect)
	server.GET("/url/", urlinfo)
	server.GET("/url/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	server.StaticFS("/static", box)
	log.Println("server is listening at " + t.ServerAddr)
	if t.UseHTTPS {
		server.RunTLS(t.ServerAddr, t.ServerCert, t.ServerKey)
	} else {
		server.Run(t.ServerAddr)
	}
}
