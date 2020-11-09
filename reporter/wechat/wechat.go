package wechat

import (
	"encoding/json"
	"fmt"
	"sign-your-horse/reporter"

	"github.com/imroc/req"
)

type WechatReporter struct {
	CorpID     string `json:"corpID"`
	CorpSecret string `json:"corpSecret`
}

func (w *WechatReporter) Init(config json.RawMessage) error {
	return json.Unmarshal(config, w)
}

func (w *WechatReporter) Report(msg string) error {
	r := req.New()
	resp, err := r.Get(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", w.CorpID, w.CorpSecret))
	if err != nil {
		return err
	}

}

func init() {
	reporter.RegisterReporter("wechat", &WechatReporter{})
}
