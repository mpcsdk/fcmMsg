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
	if tx.Status == 0 {
		g.Log().Warning(s.ctx, "transferMsgConsum status:", tx)
		msg.Ack()
		return
	}
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
		val := fval.Text('f', -1)
		// val, _ := big.NewInt(0).SetString(tx.Value, 10)
		// val.Sub()
		// val = val.Div(val, big.NewInt(int64(18)))
		////
		chain := s.chains[tx.ChainId]
		coin := ""
		if chain != nil {
			coin = chain.Coin
		}

		title := coin + ` Receive Success`
		to := addrabbr(tx.To)
		body := fmt.Sprint("Your wallet address: ", to, " has received ", val, coin, ". Please check your transaction history for details.")
		service.Fcm().PushByAddr(s.ctx, tx.To, title, body, string(msg.Data()))
		////
		title = coin + ` Send Success`
		from := addrabbr(tx.From)
		body = fmt.Sprint("Your wallet address: ", from, " has sent ", val, coin, ". Please check your transaction history for details.")
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
			val := fval.Text('f', -1)
			////
			title := contract.ContractName + ` Receive Success`
			to := addrabbr(tx.To)
			body := fmt.Sprint("Your wallet address: ", to, " has received ", val, contract.ContractName, ". Please check your transaction history for details.")
			service.Fcm().PushByAddr(s.ctx, tx.To, title, body, string(msg.Data()))
			////
			title = contract.ContractName + ` Send Success`
			from := addrabbr(tx.From)
			body = fmt.Sprint("Your wallet address: ", from, " has sent ", val, contract.ContractName, ". Please check your transaction history for details.")
			service.Fcm().PushByAddr(s.ctx, tx.From, title, body, string(msg.Data()))
			/////
		} else if contract.ContractKind == "nft" {
			title := `NFT Receive Success`
			id := big.NewInt(0)
			err := id.UnmarshalText([]byte(tx.TokenId))
			fmt.Println(err)
			tokenId := id.Text(10)
			to := addrabbr(tx.To)
			body := fmt.Sprint("Your wallet address: ", to, " has successfully received ", contract.ContractName, ":", tokenId, ". Please check your transaction history for details.")
			service.Fcm().PushByAddr(s.ctx, tx.To, title, body, string(msg.Data()))
			////
			title = `NFT Send Success`
			from := addrabbr(tx.From)
			body = fmt.Sprint("Your wallet address: ", from, " has successfully sent ", contract.ContractName, ":", tokenId, ". Please check your transaction history for details.")
			service.Fcm().PushByAddr(s.ctx, tx.From, title, body, string(msg.Data()))
			/////
		} else {
			g.Log().Notice(s.ctx, "transferMsgConsum no contract kind:", tx)
		}
		msg.Ack()
	}
}
