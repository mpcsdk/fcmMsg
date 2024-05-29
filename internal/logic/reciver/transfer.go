package reciver

import (
	"encoding/json"
	"fcmMsg/internal/service"
	"fmt"
	"math"
	"math/big"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/mpcsdk/mpcCommon/mpcdao/model/entity"
	"github.com/nats-io/nats.go/jetstream"
)

func addrabbr(addr string) string {
	return addr[0:5] + "..." + addr[len(addr)-4:]

}
func (s *sReceiver) transferMsgConsum(msg jetstream.Msg) {
	tx := &entity.ChainTx{}
	json.Unmarshal(msg.Data(), tx)
	g.Log().Debug(s.ctx, "transferMsgConsum:", tx)
	////removed , status
	// if tx.Status == 0 {
	// 	g.Log().Debug(s.ctx, "transferMsgConsum status:", tx)
	// 	msg.Ack()
	// 	return
	// }
	// addr := ""
	// fcmToken := ""
	// title := ""
	// body := ""
	// data := ""
	if tx.Kind == "external" {
		///
		fbalance := big.NewFloat(0)
		fbalance.SetString(tx.Value)
		fval := fbalance.Quo(fbalance, big.NewFloat(math.Pow10(18)))
		// val, _ := big.NewInt(0).SetString(tx.Value, 10)
		// val.Sub()
		// val = val.Div(val, big.NewInt(int64(18)))
		////
		chain := s.chains[tx.ChainId]
		coin := ""
		if chain != nil {
			coin = chain.Coin
		}

		title := coin + `接收成功`
		to := addrabbr(tx.To)
		body := fmt.Sprint("您的钱包地址：", to, "已接收 ", fval.String(), coin, "，请前往交易记录查看详情。")
		service.Fcm().PushByAddr(s.ctx, tx.To, title, body, string(msg.Data()))
		////
		title = coin + `发送成功`
		from := addrabbr(tx.From)
		body = fmt.Sprint("您的钱包地址：", from, "已发送 ", fval.String(), coin, "，请前往交易记录查看详情。")
		service.Fcm().PushByAddr(s.ctx, tx.From, title, body, string(msg.Data()))
		msg.Ack()
		return
	} else {
		contract := s.contracts[tx.Contract]
		if contract == nil {
			g.Log().Debug(s.ctx, "transferMsgConsum no contract:", tx)
			msg.Ack()
			return
		}
		/////
		if contract.ContractKind == "ft" {
			fbalance := big.NewFloat(0)
			fbalance.SetString(tx.Value)
			fval := fbalance.Quo(fbalance, big.NewFloat(math.Pow10(contract.Decimal)))
			////
			title := contract.ContractName + `接收成功`
			to := addrabbr(tx.To)
			body := fmt.Sprint("您的钱包地址：", to, "已接收 ", fval.String(), contract.ContractName, "，请前往交易记录查看详情。")
			service.Fcm().PushByAddr(s.ctx, tx.To, title, body, string(msg.Data()))
			////
			title = contract.ContractName + `发送成功`
			from := addrabbr(tx.From)
			body = fmt.Sprint("您的钱包地址：", from, "已发送 ", fval.String(), contract.ContractName, "，请前往交易记录查看详情。")
			service.Fcm().PushByAddr(s.ctx, tx.From, title, body, string(msg.Data()))
			/////
		} else if contract.ContractKind == "nft" {
			title := `NFT接收成功`
			id := big.NewInt(0)
			err := id.UnmarshalText([]byte(tx.TokenId))
			fmt.Println(err)
			fmt.Println(id.Text(10))
			// tokenId := id.String()
			to := addrabbr(tx.To)
			body := fmt.Sprint("您的钱包地址：", to, "已接收 ", contract.ContractName, id.Text(10), "，请前往交易记录查看详情。")
			service.Fcm().PushByAddr(s.ctx, tx.To, title, body, string(msg.Data()))
			////
			title = `NFT发送成功`
			from := addrabbr(tx.From)
			body = fmt.Sprint("您的钱包地址：", from, "已发送 ", contract.ContractName, id.Text(10), "，请前往交易记录查看详情。")
			service.Fcm().PushByAddr(s.ctx, tx.From, title, body, string(msg.Data()))
			/////
		} else {
			g.Log().Notice(s.ctx, "transferMsgConsum no contract kind:", tx)
		}
		msg.Ack()
	}
}
