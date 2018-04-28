package gateio

import (
	"log"
	"net/http"
	"testing"

	gobtb "github.com/BTB/goBTB/base"
)

var gate = New(http.DefaultClient, "", "")

func TestGate_GetTicker(t *testing.T) {
	return
	ticker, err := gate.GetTicker(gobtb.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}

func TestGate_GetDepth(t *testing.T) {
	// dep, err := gate.GetDepth(1, gobtb.BTC_USDT)
	// log.Println(dep)
	// t.Log("err=>", err)
	// t.Log("asks=>", dep.AskList)
	// t.Log("bids=>", dep.BidList)

	log.Printf("测试获取gateio行情:")
	m_index := map[string]gobtb.CurrencyPair{
		"BTC": gobtb.BTC_USDT,
		//"LTC": gobtb.LTC_USDT,
		//"ETC": gobtb.ETC_ETH,
		//"ETH": gobtb.ETH_ETH,
		//"BCH": gobtb.BCH_ETH,
	}
	for _, views := range m_index {
		dep, err := gate.GetDepth(1, views)
		t.Log("err=>", err)
		log.Println(dep.AskList[len(dep.AskList)-1])
		log.Println(dep.BidList[0])
	}

}
