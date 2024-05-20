package reciver

import (
	"encoding/json"
	"fcmMsg/internal/service"
	"fmt"
	"math/big"

	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/nats-io/nats.go/jetstream"
)

func (s *sReceiver) transferMsgConsum(msg jetstream.Msg) {
	tx := &entity.ChainTx{}
	json.Unmarshal(msg.Data(), tx)
	////removed , status
	contract := s.contracts[tx.Contract]
	if contract == nil {
		return
	}
	/////
	if contract.ContractKind == "ft" {
		val, _ := big.NewInt(0).SetString(tx.Value, 16)
		val.Div(val, big.NewInt(int64(contract.Decimal)))
		////
		title := `FT接收成功`
		body := fmt.Sprint("您的钱包地址：", tx.To, "已接收", val.Int64(), contract.ContractName, "，请前往交易记录查看详情。")
		service.Fcm().PushByAddr(s.ctx, tx.From, title, body, string(msg.Data()))
		////
		title = `FT发送成功`
		body = fmt.Sprint("您的钱包地址：", tx.To, "已发送", val.Int64(), contract.ContractName, "，请前往交易记录查看详情。")
		service.Fcm().PushByAddr(s.ctx, tx.From, title, body, string(msg.Data()))
		/////
	} else if contract.ContractKind == "nft" {
		title := `NFT接收成功`
		body := fmt.Sprint("您的钱包地址：", tx.To, "已接收", tx.TokenId, contract.ContractName, "，请前往交易记录查看详情。")
		service.Fcm().PushByAddr(s.ctx, tx.From, title, body, string(msg.Data()))
		////
		title = `NFT发送成功`
		body = fmt.Sprint("您的钱包地址：", tx.To, "已发送", tx.TokenId, contract.ContractName, "，请前往交易记录查看详情。")
		service.Fcm().PushByAddr(s.ctx, tx.From, title, body, string(msg.Data()))
		/////
	}

	msg.Ack()
}
