package logic

import (
	"context"
	"fcmMsg/internal/service"
	"log"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/mpcsdk/mpcCommon/mpcdao"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"google.golang.org/api/option"
)

type sFcm struct {
	ctx           context.Context
	client        *messaging.Client
	addrLock      sync.RWMutex
	addrFcmTokens map[string]string
}

func (s *sFcm) FcmToken(address string) string {
	s.addrLock.RLock()
	defer s.addrLock.RUnlock()
	return s.addrFcmTokens[address]
}
func (s *sFcm) SubFcmToken(address string, token string) {
	s.addrLock.Lock()
	defer s.addrLock.Unlock()
	s.addrFcmTokens[address] = token
}

func (s *sFcm) PushByAddr(ctx context.Context, addr string, title string, body string, data string) (string, error) {
	fcmToken := s.FcmToken(addr)
	if fcmToken == "" {
		return "", nil
	}
	////
	response, err := s.Send(ctx, fcmToken, title, body, data)
	if err != nil {
		service.DB().Fcm().InsertPushErr(ctx, &entity.PushErr{
			FmcToken: fcmToken,
			Err:      err.Error(),
			Title:    title,
			Body:     body,
			Data:     data,
		})
		g.Log().Warning(ctx, "FCM push ", "addr:", addr, "token:", fcmToken, "err:", err)
		return "", err
	}
	g.Log().Info(ctx, "FCM push ", "addr:", addr, "token:", fcmToken, "response:", response)
	return response, nil
}
func (s *sFcm) Send(ctx context.Context, token string, title string, body string, data string) (string, error) {
	response, err := s.client.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{
			"data": data,
		},
		Token: token,
	})
	if messaging.IsRegistrationTokenNotRegistered(err) {
		log.Fatalf("error initializing app: %v\n", err)
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
		fcms, lastPos, err := service.DB().Fcm().QueryFcmTokenAll(ctx, lastPos, 10)
		if err != nil {
			panic(err)
		}

		for _, fcm := range fcms {
			addrFcmTokens[fcm.Address] = fcm.FcmToken
		}
		if len(fcms) < 10 {
			break
		}
		lastPos = lastPos
	}

	////
	return &sFcm{
		ctx:           ctx,
		client:        client,
		addrFcmTokens: addrFcmTokens,
	}
}

func init() {
	service.RegisterFcm(New())
}
