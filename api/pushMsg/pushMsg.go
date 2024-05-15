// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package pushMsg

import (
	"context"
	
	"fcmMsg/api/pushMsg/v1"
)

type IPushMsgV1 interface {
	PushMsg(ctx context.Context, req *v1.PushMsgReq) (res *v1.PushMsgRes, err error)
	SubMsg(ctx context.Context, req *v1.SubMsgReq) (res *v1.SubMsgRes, err error)
}


