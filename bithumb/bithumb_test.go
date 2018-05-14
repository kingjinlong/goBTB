package bithumb

import (
	"log"
	"net/http"
	"testing"

	"github.com/BTB/goBTB/base"
)

// ConnectKey	15ef9b5fae03dc7914d0773cf4d02232
// SecretKey 427f722f23da73cf77b0fa5e2fe604b0
// var accesskey, secretkey string
// var accesskey = "15ef9b5fae03dc7914d0773cf4d02232"
// var secretkey = "427f722f23da73cf77b0fa5e2fe604b0"

var accesskey = ""
var secretkey = ""
var orderID = ""
var bh = New(http.DefaultClient, accesskey, secretkey)

// func TestBithumb_GetTicker(t *testing.T) {
// 	return
// 	ticker, err := bh.GetTicker(gobtb.NewCurrencyPair2("ALL_KAW"))
// 	t.Log("err=>", err)
// 	t.Log("ticker=>", ticker)
// }

func TestBithumb_GetAccount(t *testing.T) {
	return
	log.Printf("测试账户信息")
	back, _ := bh.GetAccount()
	if back != nil {
		for _, views := range back.SubAccounts {
			log.Println(views)
		}
	}
}

func TestBithumb_GetDepth(t *testing.T) {
	return
	// return
	log.Printf("测试获取bithumb合约行情")
	m_index := map[string]gobtb.CurrencyPair{
		// "BTC_KRW": gobtb.BTC_KRW,
		// "LTC_KRW": gobtb.LTC_KRW,
		"EOS_KRW": gobtb.EOS_KRW,
		"ETH_KRW": gobtb.ETH_KRW,
	}
	for _, views := range m_index {
		dep, _ := bh.GetDepth(1, views)
		// t.Log("err=>", err)
		// t.Log("asks=>", dep.AskList)
		// t.Log("bids=>", dep.BidList)
		log.Println(views)
		//盘口
		log.Println(dep.AskList[len(dep.AskList)-1])
		log.Println(dep.BidList[0])
		//log.Printf("Tick[%s],Close[%.2f]", md.Symbol.Code, md.Close)
	}
}

func TestBithumb_placeOrder(t *testing.T) {
	return
	log.Printf("测试:下单")
	// side := "ask"
	amount := "0.1"
	price := "22000"
	var pair gobtb.CurrencyPair
	pair.CurrencyA = gobtb.EOS
	back, _ := bh.LimitSell(amount, price, pair)
	if back != nil {
		log.Println(back)
		log.Println("订单编号为：" + string(back.OrderID))
		log.Println(back.OrderID)
		//cancel

	}

}

func TestBithumb_GetUnfinishOrders(t *testing.T) {
	//	return
	log.Printf("测试:查询报单")
	var pair gobtb.CurrencyPair
	pair.CurrencyA = gobtb.EOS
	back, _ := bh.GetUnfinishOrders(pair)
	if back != nil {
		log.Println(back)
		// var sd string

		for _, views := range back {
			// log.Println(views.OrderID)
			// var side string
			// switch views.Side {
			// case gobtb.SELL:
			// 	side = "SELL"
			// case gobtb.BUY:
			// 	side = "BUY"
			// }
			// side = gobtb.TradeSide.String(views.Side)
			back, err := bh.CancelOrder2(gobtb.TradeSide.String(views.Side), views.OrderID2, views.Currency)
			if back == false {
				log.Printf("测试:失败")
				log.Println(err)
				//cancel
			} else {
				log.Printf("测试:成功")
			}
		}

		//cancel
	}
}

func TestBithumb_CancelOrder2(t *testing.T) {
	return
	log.Printf("测试:撤单")
	orderID := "1526005653075389"
	side := "SELL"
	var pair gobtb.CurrencyPair
	pair.CurrencyA = gobtb.EOS
	back, err := bh.CancelOrder2(side, orderID, pair)
	if back == false {
		log.Printf("测试:失败")
		log.Println(err)
		//cancel
	} else {
		log.Printf("测试:成功")
	}

}
