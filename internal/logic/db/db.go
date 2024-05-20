package db

import (
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
	return &sDB{}
}

func init() {
	service.RegisterDB(NewDB())
}
