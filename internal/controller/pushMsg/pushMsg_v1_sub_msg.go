package pushMsg

import (
	"context"

	v1 "fcmMsg/api/pushMsg/v1"
	"fcmMsg/internal/service"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/mpcsdk/mpcCommon/mpccode"
	"github.com/mpcsdk/mpcCommon/userInfoGeter"
)

func (c *ControllerV1) SubMsg(ctx context.Context, req *v1.SubMsgReq) (*v1.SubMsgRes, error) {
	////
	//trace
	ctx, span := gtrace.NewSpan(ctx, "SubMsg")
	defer span.End()
	g.Log().Debug(ctx, "SubMsg : ", req)
	//
	if req.Address == "" || req.FcmToken == "" || req.Token == "" {
		return nil, mpccode.CodeParamInvalid()
	}
	//
	// info, err := service.UserInfo().GetUserInfo(ctx, req.Token)
	// if err != nil {
	// 	g.Log().Error(ctx, "SubMsg : ", req, err)
	// 	g.RequestFromCtx(ctx).Response.WriteStatusExit(500)
	// 	return nil, mpccode.CodeTokenInvalid()
	// }
	info := &userInfoGeter.UserInfo{
		UserId: "testId",
	}
	////check addr
	req.Address = common.HexToAddress(req.Address).String()
	////
	err := service.Fcm().SubFcmToken(ctx, info.UserId, req.Address, req.FcmToken, req.Token)
	if err != nil {
		g.Log().Warning(ctx, "SubMsg : ", req, err)
		return nil, mpccode.CodeInternalError()
	}
	return nil, nil
}
