package wechat

import (
	"encoding/json"
	"fmt"
	"sign-your-horse/common"
	"sign-your-horse/reporter"
	"time"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

type WechatReporter struct {
	CorpID     string `json:"corpID"`
	CorpSecret string `json:"corpSecret"`
	ToParty    int    `json:"toparty"`
	AgentID    int    `json:"agentid"`
}

type WechatPushBody struct {
	ToParty int                `json:"toparty"`
	MsgType string             `json:"msgtype"`
	AgentID int                `json:"agentid"`
	Text    WechatPushBodyText `json:"text"`
	Safe    int                `json:"safe"`
}

type WechatPushBodyText struct {
	Content string `json:"content"`
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
	token := gjson.Get(resp.String(), "access_token").String()
	if token == "" {
		return common.Raise("invalid wechat token")
	}
	pushBody := &WechatPushBody{
		ToParty: w.ToParty,
		MsgType: "text",
		AgentID: w.AgentID,
		Text: WechatPushBodyText{
			Content: "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + msg,
		},
		Safe: 0,
	}
	r.SetJSONEscapeHTML(false)
	pushResp, err := r.Post(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token), req.BodyJSON(pushBody))
	if err != nil {
		return err
	}
	common.LogWithModule("wechat_push", "wechat pusher: "+pushResp.String())
	return nil
}

func init() {
	reporter.RegisterReporter("wechat", &WechatReporter{})
}
