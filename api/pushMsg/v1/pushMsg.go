package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PushMsgReq struct {
	g.Meta   `path:"/pushMsg" tags:"pushMsg" method:"post" summary:"You first hello api"`
	Token    string `json:"token"`
	FcmToken string `json:"fcmToken"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Data     string `json:"data"`
}
type PushMsgRes struct {
	g.Meta   `mime:"text/html" example:"string"`
	Response string `json:"response"`
}

type SubMsgReq struct {
	g.Meta   `path:"/pushMsg" tags:"pushMsg" method:"post" summary:"You first hello api"`
	Token    string `json:"token"`
	FcmToken string `json:"fcmToken"`
	Address  string `json:"address"`
}
type SubMsgRes struct {
	g.Meta `mime:"text/html" example:"string"`
}
