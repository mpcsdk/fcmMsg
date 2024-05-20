package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type SendMsgReq struct {
	g.Meta   `path:"/sendMsg" tags:"sendMsg" method:"post" summary:"You first hello api"`
	Token    string `json:"token"`
	FcmToken string `json:"fcmToken"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Data     string `json:"data"`
}
type SendMsgRes struct {
	g.Meta   `mime:"text/html" example:"string"`
	Response string `json:"response"`
}

type SubMsgReq struct {
	g.Meta   `path:"/subMsg" tags:"subMsg" method:"post" summary:"You first hello api"`
	Token    string `json:"token"`
	FcmToken string `json:"fcmToken"`
	Address  string `json:"address"`
}
type SubMsgRes struct {
	g.Meta `mime:"text/html" example:"string"`
}
