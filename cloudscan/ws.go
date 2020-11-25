package cloudscan

import (
	"context"
	"net/http"
	"sign-your-horse/cloudscan_client"
	"sign-your-horse/common"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clientChanMap = make(map[string]chan string)

func (c *CloudScanAPIServer) handleWebSocketClient(clientConn *websocket.Conn, clientID string) {
	clientCtx, cancel := context.WithCancel(context.Background())
	hbTicker := time.NewTicker(time.Second * time.Duration(60))
	defer hbTicker.Stop()
	defer delete(clientChanMap, clientID)
	defer clientConn.Close()
	lastHeartBeatTime := time.Now()
	pushMessageChan := clientChanMap[clientID]

	//input dispatcher
	go func() {
		for {
			incommingMessage := &cloudscan_client.WSMessage{}
			err := clientConn.ReadJSON(incommingMessage)
			if err != nil {
				common.LogWithModule(moduleName, "Read data from %s failed: %s", clientID, err.Error())
				cancel()
				return
			}
			if incommingMessage.MessageType != cloudscan_client.WS_HEARTBEAT_CLIENT {
				common.LogWithModule(moduleName, "Invalid data received from %s failed", clientID)
				cancel()
				return
			}
			common.LogWithModule(moduleName, "accept heartbeat from %s[%s]", clientConn.RemoteAddr().String(), clientID)
			retMessage := &cloudscan_client.WSMessage{
				MessageType: cloudscan_client.WS_HEARTBEAT_SERVER,
				MessageData: incommingMessage.MessageData,
			}
			err = clientConn.WriteJSON(retMessage)
			if err != nil {
				common.LogWithModule(moduleName, "Write data to %s failed", clientID)
				cancel()
				return
			}
			lastHeartBeatTime = time.Now()
		}
	}()

	for {
		select {
		case <-clientCtx.Done():
			return

		case <-hbTicker.C:
			if lastHeartBeatTime.Add(time.Second * time.Duration(60)).Before(time.Now()) {
				common.LogWithModule(moduleName, "client %s does not response for 60s, kill connection", clientID)
				cancel()
				return
			}

		case m := <-pushMessageChan:
			common.LogWithModule(moduleName, "push message %s to %s", m, clientID)
			err := clientConn.WriteJSON(&cloudscan_client.WSMessage{
				MessageType: cloudscan_client.WS_DATA,
				MessageData: m,
			})
			if err != nil {
				common.LogWithModule(moduleName, "Write data to %s failed", clientID)
				cancel()
				return
			}
		}
	}
}

func (c *CloudScanAPIServer) handleWebSocket(gCtx *gin.Context) {
	if gCtx.GetHeader("X-Auth") != c.APIPassword {
		common.LogWithModule(moduleName, "invalid password from client: %s", gCtx.Request.RemoteAddr)
		gCtx.Status(http.StatusUnauthorized)
		return
	}
	clientID := gCtx.GetHeader("X-Client-ID")
	if clientID == "" {
		common.LogWithModule(moduleName, "invalid client ID from client: %s", gCtx.Request.RemoteAddr)
		gCtx.String(http.StatusBadRequest, "invalid client ID")
	}

	conn, err := wsupgrader.Upgrade(gCtx.Writer, gCtx.Request, nil)
	if err != nil {
		common.LogWithModule(moduleName, "Failed to upgrade websocket: ", err.Error())
		return
	}
	common.LogWithModule(moduleName, "accept connection from client: %s[%s]", gCtx.Request.RemoteAddr, clientID)
	clientChanMap[clientID] = make(chan string)
	c.handleWebSocketClient(conn, clientID)
}
