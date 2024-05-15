package reciver

import (
	"context"
	"encoding/json"
	"fcmMsg/internal/conf"
	"fcmMsg/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/mpcsdk/mpcCommon/mq"
	"github.com/nats-io/nats.go/jetstream"
)

type sReceiver struct {
	ctx  context.Context
	nats *mq.NatsServer
	///
	cons_transfer jetstream.Consumer
	cons_mint     jetstream.Consumer
}

func new() *sReceiver {

	ctx := gctx.GetInitCtx()
	///
	nats := mq.New(conf.Config.Nrpc.NatsUrl)

	////transfer msg
	cons_transfer, err := nats.GetConsumer("fcmTransfer", mq.JetStream_SyncChain, mq.JetSub_SyncChainTransfer)
	if err != nil {
		panic(err)
	}
	////mint msg
	cons_mint, err := nats.GetConsumer("fcmMint", mq.JetStream_SyncChain, mq.JetSub_SyncChainTransfer)
	if err != nil {
		panic(err)
	}
	////takerbid msg
	///crossererc20 msg
	////swapcoress msg
	///
	r := g.Redis("aggRiskCtrl")
	_, err = r.Conn(gctx.GetInitCtx())
	if err != nil {
		panic(err)
	}
	///
	s := &sReceiver{
		ctx:           ctx,
		nats:          nats,
		cons_transfer: cons_transfer,
		cons_mint:     cons_mint,
	}

	cons_transfer.Consume(s.transferMsgConsum)
	cons_mint.Consume(func(msg jetstream.Msg) {
		tx := &entity.ChainTx{}
		json.Unmarshal(msg.Data(), tx)

		msg.Ack()
	})
	///
	///

	///
	return s
}

func init() {
	service.RegisterReceiver(new())
}