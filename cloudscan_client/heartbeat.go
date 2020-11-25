package cloudscan_client

import (
	"sign-your-horse/common"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

func (c *CloudScanMessageClient) HeartBeatClient(ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(time.Second * time.Duration(c.HeartBeatInterval))
	defer ticker.Stop()
	for {
		hbCtx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(3))
		select {
		case t := <-ticker.C:
			err := c.wsConn.WriteJSON(&WSMessage{
				MessageType: WS_HEARTBEAT_CLIENT,
				MessageData: strconv.FormatInt(t.Unix(), 10),
			})
			if err != nil {
				common.LogWithModule(moduleName, "write heartbeat failed: "+err.Error())
				cancel()
				return
			}
			select {
			case hbRet := <-c.heartBeatChan:
				//hb response
				t2, err := strconv.ParseInt(hbRet, 10, 64)
				if err != nil {
					common.LogWithModule(moduleName, "invalid heartbeat response")
					cancel()
					return
				}
				if c.Verbose {
					common.LogWithModule(moduleName, "heartbeat latency "+strconv.FormatInt(t2-time.Now().Unix(), 10)+"ms")
				}

			case <-ctx.Done():
				return
			case <-hbCtx.Done():
				common.LogWithModule(moduleName, "heartbeat timeout")
				cancel()
				return
				//timeout
			}
		case <-ctx.Done():
			return
		case <-hbCtx.Done():
		}
	}
}
