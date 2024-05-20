// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IFcm interface {
		FcmToken(address string) string
		SubFcmToken(address string, token string)
		PushByAddr(ctx context.Context, addr string, title string, body string, data string) (string, error)
		Send(ctx context.Context, fcmToken string, title string, body string, data string) (string, error)
	}
)

var (
	localFcm IFcm
)

func Fcm() IFcm {
	if localFcm == nil {
		panic("implement not found for interface IFcm, forgot register?")
	}
	return localFcm
}

func RegisterFcm(i IFcm) {
	localFcm = i
}
