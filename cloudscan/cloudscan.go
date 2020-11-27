package cloudscan

import (
	"encoding/json"
	"sign-your-horse/common"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
)

const moduleName = "cloudscan"

//MessageInputChan receive message from HTTP API, main.go will call every provider with message received from this channel
var MessageInputChan chan string

func Init(config json.RawMessage) (*CloudScanAPIServer, error) {
	ret := &CloudScanAPIServer{}
	err := json.Unmarshal(config, ret)
	if ret.APIPassword == "" && ret.EnableAPIServer {
		return nil, common.Raise("empty ws_api_password is not allowd")
	}
	ret.APIMessageInputChan = make(chan string)
	MessageInputChan = ret.APIMessageInputChan
	return ret, err
}

func (t *CloudScanAPIServer) Run() {
	if !t.EnableServer {
		return
	}
	defer close(t.APIMessageInputChan)

	box := packr.NewBox("static")
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.POST("/url/add", add)
	server.GET("/url/raw", raw)
	server.GET("/url/redirect", redirect)
	server.GET("/url/", urlinfo)
	if t.EnableAPIServer {
		server.GET("/url/ws", t.handleWebSocket)
	}
	server.StaticFS("/static", box)
	common.LogWithModule(moduleName, "server is listening at "+t.ServerAddr)
	if t.UseHTTPS {
		common.Must(server.RunTLS(t.ServerAddr, t.ServerCert, t.ServerKey))
	} else {
		common.LogWithModule(moduleName, "server with http is not recommand! WebRTC and Websocket will not work")
		common.Must(server.Run(t.ServerAddr))
	}
}

func (t *CloudScanAPIServer) Push(msg string) {
	if t.EnableAPIServer {
		for clientID := range clientChanMap {
			clientChanMap[clientID] <- msg
		}
	}
}
