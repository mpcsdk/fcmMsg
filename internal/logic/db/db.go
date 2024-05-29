package db

import (
	"fcmMsg/internal/conf"
	"fcmMsg/internal/service"

	"github.com/mpcsdk/mpcCommon/mpcdao"
)

type sDB struct {
	fmc *mpcdao.Fcm
}

func (s *sDB) Fcm() *mpcdao.Fcm {
	return s.fmc
}

// /
func NewDB() *sDB {
	return &sDB{
		fmc: mpcdao.NewFcm(nil, conf.Config.Cache.Duration),
	}
}

func init() {
	service.RegisterDB(NewDB())
}
