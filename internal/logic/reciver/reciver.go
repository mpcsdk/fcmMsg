package reciver

import (
	"context"
	"encoding/json"
	"fcmMsg/internal/conf"
	"fcmMsg/internal/service"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/mpcsdk/mpcCommon/mpcdao"
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
	////db
	contracts map[string]*entity.Contractabi
	chains    map[int64]*entity.Chaincfg
	///
}

func new() *sReceiver {

	ctx := gctx.GetInitCtx()

	nats := mq.New(conf.Config.Nrpc.NatsUrl)

	///
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
	time.Sleep(10 * time.Second)
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
		contracts:     map[string]*entity.Contractabi{},
		chains:        map[int64]*entity.Chaincfg{},
	}
	/////chaincfg
	chaincfgDb := mpcdao.NewChainCfg()
	chaincfgs, err := chaincfgDb.AllCfg(s.ctx)
	if err != nil {
		panic(err)
	}
	for _, c := range chaincfgs {
		s.chains[c.ChainId] = c
	}
	/////
	// contracts, err := chaincfg.GetContractAbiBriefs(s.ctx, 0, "")
	// if err != nil {
	// 	panic(err)
	// }
	// for _, c := range contracts {
	// 	s.contracts[c.ContractAddress] = c
	// }
	///
	go func() {
		time.Sleep(20 * time.Second)
		ruledb := mpcdao.NewRiskCtrlRule(nil, 0)
		contracts, err := ruledb.GetContractAbiBriefs(s.ctx, 0, "")
		if err != nil {
			panic(err)
		}
		for _, c := range contracts {
			s.contracts[c.ContractAddress] = c
		}
		////
		s.cons_transfer.Consume(s.transferMsgConsum)
		s.cons_mint.Consume(func(msg jetstream.Msg) {
			tx := &entity.ChainTx{}
			json.Unmarshal(msg.Data(), tx)

			msg.Ack()
		})
	}()
	///
	return s
}
func (s *sReceiver) StartSub() {

}

func init() {
	service.RegisterReceiver(new())
}
