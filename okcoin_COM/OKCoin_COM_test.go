package okcoin

import (
	"log"
	"net/http"
	"testing"

	gobtb "github.com/BTB/goBTB/base"
	"github.com/stretchr/testify/assert"
)

// 比特币-usd
var okcn = New(http.DefaultClient, "", "")

// 币币交易
var okexSpot = NewOKExSpot(http.DefaultClient, "", "")

func initial() {

}

func TestOKCoinCN_API_GetTicker(t *testing.T) {
	return
	ticker, _ := okcn.GetTicker(gobtb.ETH_USD)
	t.Log(ticker)
	log.Println("res: ETH_USD")
	log.Println(ticker)
}

func TestOKCoinCN_API_GetDepth(t *testing.T) {
	log.Printf("测试获取okcoin行情:几个货币兑美元的价格")
	m_index := map[string]gobtb.CurrencyPair{
		"BTC": gobtb.BTC_USD,
		"LTC": gobtb.LTC_USD,
		"ETC": gobtb.ETC_USD,
		"ETH": gobtb.ETH_USD,
		"BCH": gobtb.BCH_USD,
	}
	for _, views := range m_index {
		dep, err := okcn.GetDepth(1, views)
		t.Log("err=>", err)
		// t.Log("asks=>", dep.AskList)
		// t.Log("bids=>", dep.BidList)
		// log.Println(views)
		//盘口
		log.Println(dep.AskList[len(dep.AskList)-1])
		log.Println(dep.BidList[0])
	}
}

func TestOKExSpot_GetTicker(t *testing.T) {
	return
	ticker, err := okexSpot.GetTicker(gobtb.ETC_BTC)
	assert.Nil(t, err)
	t.Log(ticker)
}

func TestOKExSpot_GetDepth(t *testing.T) {
	return // 尚未测试通过，vpn
	log.Printf("测试获取okcoin行情:TestOKExSpot_GetDepth")
	dep, err := okexSpot.GetDepth(2, gobtb.ETC_BTC)
	assert.Nil(t, err)
	t.Log(dep)
}
