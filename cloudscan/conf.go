package cloudscan

import (
	"encoding/json"
)

type CloudScanAPIServer struct {
	EnableServer    bool   `json:"server_enable"`
	ServerAddr      string `json:"server"`
	UseHTTPS        bool   `json:"usehttps"`
	ServerCert      string `json:"srvcert"`
	ServerKey       string `json:"srvkey"`
	EnableAPIServer bool   `json:"ws_api_enable"`
	APIPassword     string `json:"ws_api_password"`

	APIMessageInputChan chan string `json:"-"`
}

func DumpCloudScanAPIServerConfig() json.RawMessage {
	ret, _ := json.Marshal(&CloudScanAPIServer{
		EnableServer:    true,
		ServerAddr:      "0.0.0.0:3000",
		UseHTTPS:        false,
		ServerCert:      "cert.pem",
		ServerKey:       "key.pem",
		EnableAPIServer: true,
		APIPassword:     "1145141919810",
	})
	return ret
}
