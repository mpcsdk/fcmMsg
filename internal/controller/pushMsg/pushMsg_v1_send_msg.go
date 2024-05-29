package pushMsg

import (
	"context"

	v1 "fcmMsg/api/pushMsg/v1"
	"fcmMsg/internal/service"
)

func (c *ControllerV1) SendMsg(ctx context.Context, req *v1.SendMsgReq) (res *v1.SendMsgRes, err error) {
	////
	response, err := service.Fcm().Send(ctx, req.FcmToken, req.Title, req.Body, req.Data)
	return &v1.SendMsgRes{
		Response: response,
	}, err
}
