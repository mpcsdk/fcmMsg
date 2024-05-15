package pushMsg

import (
	"context"

	v1 "fcmMsg/api/pushMsg/v1"
	"fcmMsg/internal/service"
)

func (c *ControllerV1) PushMsg(ctx context.Context, req *v1.PushMsgReq) (res *v1.PushMsgRes, err error) {
	///check token

	/////
	/////
	response, err := service.Fcm().Send(ctx, req.FcmToken, req.Title, req.Body, req.Data)
	return &v1.PushMsgRes{
		Response: response,
	}, err
}
