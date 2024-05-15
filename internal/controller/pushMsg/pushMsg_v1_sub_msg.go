package pushMsg

import (
	"context"

	v1 "fcmMsg/api/pushMsg/v1"
	"fcmMsg/internal/service"

	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

func (c *ControllerV1) SubMsg(ctx context.Context, req *v1.SubMsgReq) (res *v1.SubMsgRes, err error) {
	////checke
	err = service.DB().Fcm().InsertFcmToken(ctx, &entity.FcmToken{
		Token:    req.Token,
		FcmToken: req.FcmToken,
		Address:  req.Address,
	})
	if err != nil {
		return nil, err
	}
	////
	service.Fcm().SubFcmToken(req.Address, req.FcmToken)
	return nil, err
}
