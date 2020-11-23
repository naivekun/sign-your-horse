package cloudscan

import (
	"encoding/json"
)

type CloudScanAPIServer struct {
	ServerAddr string `json:"server"`
	UseHTTPS   bool   `json:"usehttps"`
	ServerCert string `json:"srvcert"`
	ServerKey  string `json:"srvkey"`

	APIMessageOutputChan chan string `json:"-"`
}

func DumpCloudScanAPIServerConfig() json.RawMessage {
	ret, _ := json.Marshal(&CloudScanAPIServer{
		ServerAddr: "0.0.0.0:3000",
		UseHTTPS:   false,
	})
	return ret
}
