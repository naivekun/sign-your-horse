package cloudscan_client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sign-your-horse/common"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

const moduleName = "cloudscan_client"

const (
	WS_HEARTBEAT_CLIENT = iota
	WS_HEARTBEAT_SERVER
	WS_NODATA
	WS_DATA
)

func Init(config json.RawMessage) (*CloudScanMessageClient, error) {
	ret := &CloudScanMessageClient{}
	err := json.Unmarshal(config, ret)
	if ret.APIServerPassword == "" {
		return nil, common.Raise("empty ws_api_password is not allowed")
	}
	if !strings.HasPrefix(ret.APIServerAddress, "wss") {
		return nil, common.Raise("API connection must use tls")
	}
	if ret.HeartBeatInterval > 30 {
		return nil, common.Raise("heartbeat interval should not longer than 30s")
	}
	ret.heartBeatChan = make(chan string)
	ret.messageInputChan = make(chan string)
	ret.MessageOutputChan = make(chan string)
	return ret, err
}

type WSMessage struct {
	MessageType byte   `json:"type"`
	MessageData string `json:"data"`
}

func (c *CloudScanMessageClient) Dispatcher(ctx context.Context, cancel context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			incommingMsg := &WSMessage{}
			err := c.wsConn.ReadJSON(incommingMsg)
			if err != nil {
				common.LogWithModule(moduleName, "websocket ReadMessage failed: "+err.Error())
				cancel()
				break
			}
			switch incommingMsg.MessageType {
			case WS_HEARTBEAT_SERVER:
				c.heartBeatChan <- incommingMsg.MessageData
			case WS_DATA:
				c.messageInputChan <- incommingMsg.MessageData
			case WS_NODATA:
			}
		}
	}
}

func getNextReconnectInterval(lastInterval int) int {
	if lastInterval < 30*1000 {
		return lastInterval * 2
	}
	return 30 * 1000
}

func (c *CloudScanMessageClient) Run() {
	if !c.EnableClient {
		return
	}
	defer close(c.heartBeatChan)
	defer close(c.messageInputChan)
	defer close(c.MessageOutputChan)
	reconnectInterval := 100
	wsDialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	header := http.Header{}
	header.Add("X-Auth", c.APIServerPassword)
	header.Add("X-Client-ID", uuid.New().String())

	for {
		ctx, cancel := context.WithCancel(context.Background())
		conn, _, err := wsDialer.Dial(c.APIServerAddress, header)
		if err != nil {
			common.LogWithModule(moduleName, "websocket dial failed: "+err.Error())
		} else {
			c.wsConn = conn
			reconnectInterval = 100
			common.LogWithModule(moduleName, "websocket connection established: "+c.APIServerAddress)

			go c.Dispatcher(ctx, cancel)
			go c.HeartBeatClient(ctx, cancel)

		lMessageInput:
			for {
				select {
				case d := <-c.messageInputChan:
					c.MessageOutputChan <- d

				case <-ctx.Done():
					break lMessageInput
				}
			}
		}
		if conn != nil {
			err = conn.Close()
			if err != nil {
				common.LogWithModule(moduleName, "connection error: "+err.Error())
			}
		}
		reconnectInterval = getNextReconnectInterval(reconnectInterval)
		common.LogWithModule(moduleName, fmt.Sprintf("reconnect after %dms", reconnectInterval))
		time.Sleep(time.Millisecond * time.Duration(reconnectInterval))
	}
}
