package logic

import (
	"context"
	"fcmMsg/internal/service"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/mpcsdk/mpcCommon/mpccode"
	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/mpcsdk/mpcCommon/rand"
	"google.golang.org/api/option"
)

type scanData struct {
	fcmToken string
	addr     string
}
type sFcm struct {
	ctx           context.Context
	client        *messaging.Client
	addrLock      sync.RWMutex
	addrFcmTokens map[string]string
	scanChan      chan scanData
}

func (s *sFcm) FcmToken(address string) string {
	s.addrLock.RLock()
	defer s.addrLock.RUnlock()
	return s.addrFcmTokens[address]
}
func (s *sFcm) SubFcmToken(ctx context.Context, userId, address string, fcmToken string, token string) error {
	///
	fcm, _ := service.DB().Fcm().QueryFcmToken(ctx, &mpcdao.QueryFcmToken{
		Address: address,
	})
	// if err != nil {
	// 	g.Log().Warning(ctx, "SubFcmToken:", "userId:", userId, "address:", address, "fcmToken:", fcmToken, "token:", token, "err:", err)
	// 	return mpccode.CodeInternalError()
	// }
	////
	if fcm == nil {
		err := service.DB().Fcm().InsertFcmToken(ctx, &entity.FcmToken{
			UserId:      userId,
			Token:       token,
			FcmToken:    fcmToken,
			Address:     address,
			CreatedTime: gtime.Now(),
			UpdatedTime: gtime.New(),
		})
		if err != nil {
			g.Log().Warning(ctx, "SubFcmToken:", "userId:", userId, "address:", address, "fcmToken:", fcmToken, "token:", token, "err:", err)
			return mpccode.CodeInternalError()
		}
	} else {
		if fcm.FcmToken == fcmToken && fcm.Address == address {
		} else {
			fcm.FcmToken = fcmToken
			fcm.UpdatedTime = gtime.Now()
			err := service.DB().Fcm().UpdateFcmToken(ctx, fcm.Address, fcm)
			if err != nil {
				g.Log().Warning(ctx, "SubFcmToken:", "userId:", userId, "address:", address, "fcmToken:", fcmToken, "token:", token, "err:", err)
			}
		}
	}
	//
	{
		s.addrLock.Lock()
		s.addrFcmTokens[address] = fcmToken
		s.addrLock.Unlock()
	}
	///scanOffline msg
	s.scanChan <- scanData{
		fcmToken: fcmToken,
		addr:     address,
	}

	return nil
}

func (s *sFcm) scanOfflineMsg() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case scanData := <-s.scanChan:
			s.scanOffline(s.ctx, scanData.addr, scanData.fcmToken)
		default:
		}
	}
}
func (s *sFcm) scanOffline(ctx context.Context, address string, token string) {
	var lastPos mpcdao.PosFcmOffline = nil
	for {
		datas, pos, err := service.DB().Fcm().QueryFcmOfflineMsg(s.ctx, lastPos, 10)
		if err != nil {
			g.Log().Warning(ctx, "scanOffline err:", err)
			return
		}
		////
		ids := []string{}
		for _, data := range datas {
			_, err := s.PushByAddr(ctx, data.Address, data.Title, data.Body, data.Data)
			if err == nil {
				ids = append(ids, data.Id)
			}
		}
		service.DB().Fcm().DeleteOfflineMsgs(ctx, ids)
		if len(datas) < 10 {
			break
		}
		lastPos = pos
	}
}
func (s *sFcm) PushByAddr(ctx context.Context, addr string, title string, body string, data string) (string, error) {
	fcmToken := s.FcmToken(addr)
	g.Log().Debug(ctx, "PushByAddr:", "addr:", addr, "token:", fcmToken, "title:", title, "body:", body, "data:", data)
	if fcmToken == "" {
		g.Log().Info(ctx, "PushByAddr UnSub add:", addr)
		return "", nil
	}
	////
	response, err := s.Send(ctx, fcmToken, title, body, data)
	if err != nil {
		if messaging.IsUnregistered(err) {
			id := rand.GenNewSid()
			err := service.DB().Fcm().InsertFcmOfflineMsg(ctx, &entity.FcmOfflineMsg{
				Id:          id,
				FmcToken:    fcmToken,
				Title:       title,
				Body:        body,
				Data:        data,
				UserId:      "",
				Address:     addr,
				CreatedTime: gtime.Now(),
			})
			if err != nil {
				g.Log().Warning(ctx, "FCM InsertFcmOfflineMsg ", "addr:", addr, "token:", fcmToken, "err:", err)
				return "", mpccode.CodeParamInvalid()
			}
		} else {
			service.DB().Fcm().InsertPushErr(ctx, &entity.PushErr{
				FmcToken: fcmToken,
				Err:      err.Error(),
				Title:    title,
				Body:     body,
				Data:     data,
			})
			g.Log().Warning(ctx, "FCM push ", "addr:", addr, "token:", fcmToken, "err:", err)
			return "", mpccode.CodeInternalError()
		}
	}
	g.Log().Info(ctx, "FCM push ", "addr:", addr, "token:", fcmToken, "response:", response)
	return response, nil
}
func (s *sFcm) Send(ctx context.Context, fcmToken string, title string, body string, data string) (string, error) {
	response, err := s.client.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				ClickAction: "com.mmsdk.opengame",
			},
		},
		Data: map[string]string{
			"data": data,
		},
		Token: fcmToken,
	})
	if err != nil {
		if messaging.IsUnregistered(err) {
			service.DB().Fcm().InsertFcmOfflineMsg(ctx, &entity.FcmOfflineMsg{
				FmcToken: fcmToken,
				Title:    title,
				Body:     body,
				Data:     data,
				// UserId: "",
				// Address: "",
			})
			g.Log().Info(ctx, "FCM IsUnregistered", "fcmToken:", fcmToken)
		} else {
			g.Log().Warning(ctx, "FCM push ", "fcmToken:", fcmToken, "err:", err)
		}
	}
	g.Log().Info(ctx, "response:", response)
	return response, err
}

// /////
func New() *sFcm {
	ctx := gctx.GetInitCtx()
	opt := option.WithCredentialsFile("token.json")
	config := &firebase.Config{ProjectID: "tantalum-f449b"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		panic(err)
	}
	client, err := app.Messaging(context.Background())
	if err != nil {
		panic(err)
	}
	/////
	addrFcmTokens := map[string]string{}
	var lastPos mpcdao.PosFcmToken = nil
	for {
		fcms, pos, err := service.DB().Fcm().QueryFcmTokenAll(ctx, lastPos, 10)
		if err != nil {
			panic(err)
		}

		for _, fcm := range fcms {
			addrFcmTokens[fcm.Address] = fcm.FcmToken
		}
		if len(fcms) < 10 {
			break
		}
		lastPos = pos

	}

	////
	return &sFcm{
		ctx:           ctx,
		client:        client,
		addrFcmTokens: addrFcmTokens,
		scanChan:      make(chan scanData, 1000),
	}
}

func init() {
	service.RegisterFcm(New())
}
