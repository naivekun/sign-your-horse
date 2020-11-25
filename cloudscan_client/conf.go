package cloudscan_client

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type CloudScanMessageClient struct {
	EnableClient      bool   `json:"client_enable"`
	APIServerAddress  string `json:"server"`
	APIServerPassword string `json:"ws_api_password"`
	Verbose           bool   `json:"verbose"`
	HeartBeatInterval int    `json:"heartbeat_interval"`

	heartBeatChan     chan string     `json:"-"`
	messageInputChan  chan string     `json:"-"`
	MessageOutputChan chan string     `json:"-"`
	wsConn            *websocket.Conn `json:"-"`
}

func DumpCloudScanMessageClientConfig() json.RawMessage {
	ret, _ := json.Marshal(&CloudScanMessageClient{
		EnableClient:      false,
		Verbose:           true,
		HeartBeatInterval: 10,
		APIServerAddress:  "wss://www.pornhub.com/url/ws",
		APIServerPassword: "1145141919810",
	})
	return ret
}
