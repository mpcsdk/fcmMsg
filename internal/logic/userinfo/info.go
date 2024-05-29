package userInfo

import (
	"context"
	"fcmMsg/internal/conf"
	"fcmMsg/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpccode"
	"github.com/mpcsdk/mpcCommon/userInfoGeter"
)

type sUserInfo struct {
	url string
	///
	userGeter *userInfoGeter.UserTokenInfoGeter
}

func (s *sUserInfo) GetUserInfo(ctx context.Context, userToken string) (userInfo *userInfoGeter.UserInfo, err error) {
	if userToken == "" {
		return nil, mpccode.CodeParamInvalid()
	}
	///
	// 用户信息示例
	// "id": 10,
	// "appPubKey": "038c90b87d77f2cc3d26132e1ea26e14646d663e3f43f17180345df3d54b8b5c70",
	// "email": "sunwenhao0421@163.com",
	// "loginType": "tkey-auth0-twitter-cyan",
	// "address": "0xe73E35d8Ecc3972481138D01799ED3934cc57853",
	// "keyHash": "U2FsdGVkX1/O6j9czaWzdjjDo/XPjk1hI8pIoaxSuS52zIxVuStK/nS07ucgiM5si8NjN97rAux3aH7Ld2i5oO8UuL6tpNZmLMG9ZpwVTxvGkCa3H14vTxWNz+yBoWG8",
	// "create_time": 1691118876

	///
	info, err := s.userGeter.GetUserInfo(ctx, userToken)
	if err != nil {
		g.Log().Error(ctx, "GetUserInfo:", "toekn:", userToken, "err:", err)
		return info, mpccode.CodeInternalError()
	}
	return info, err
}

func new() *sUserInfo {
	///
	url := conf.Config.UserTokenUrl
	///
	userGeter := userInfoGeter.NewUserInfoGeter(url, conf.Config.Cache.Duration)
	_, err := userGeter.GetUserInfo(context.Background(), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBQdWJLZXkiOiJhYmNkIiwiaWF0IjoxNjk0NDk5Njg5LCJleHAiOjE3MjYwMzU2ODl9.OsI4nFQoSoegZJbzTQnWBaB1shMjaPinhWZlnntGub4")
	if err != nil {
		panic(err)
	}
	//
	s := &sUserInfo{
		userGeter: userGeter,
	}
	///

	return s
}

func init() {
	service.RegisterUserInfo(new())
}
