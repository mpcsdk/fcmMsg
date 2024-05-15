package reciver

import (
	"encoding/json"
	"fcmMsg/internal/service"

	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/nats-io/nats.go/jetstream"
)

func (s *sReceiver) transferMsgConsum(msg jetstream.Msg) {
	tx := &entity.ChainTx{}
	json.Unmarshal(msg.Data(), tx)
	service.Fcm().PushByAddr(s.ctx, tx.From, "tansferFrom", tx.To, string(msg.Data()))
	msg.Ack()
}
