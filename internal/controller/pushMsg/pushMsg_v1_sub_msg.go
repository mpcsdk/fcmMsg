package pushMsg

import (
	"context"

	v1 "fcmMsg/api/pushMsg/v1"
	"fcmMsg/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/mpcsdk/mpcCommon/mpccode"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
)

func (c *ControllerV1) SubMsg(ctx context.Context, req *v1.SubMsgReq) (*v1.SubMsgRes, error) {
	////
	//trace
	ctx, span := gtrace.NewSpan(ctx, "SubMsg")
	defer span.End()
	g.Log().Debug(ctx, "SubMsg : ", req)
	//
	_, err := service.UserInfo().GetUserInfo(ctx, req.Token)
	if err != nil {
		g.Log().Error(ctx, "SubMsg : ", req, err)
		g.RequestFromCtx(ctx).Response.WriteStatusExit(500)
		return nil, mpccode.CodeTokenInvalid()
	}
	////checke
	err = service.DB().Fcm().InsertFcmToken(ctx, &entity.FcmToken{
		Token:    req.Token,
		FcmToken: req.FcmToken,
		Address:  req.Address,
	})
	if err != nil {
		return nil, mpccode.CodeInternalError()
	}
	////
	service.Fcm().SubFcmToken(req.Address, req.FcmToken)
	return nil, nil
}
