package bithumb

import (
	"log"
	"net/http"
	"testing"

	"github.com/BTB/goBTB/base"
)

var bh = New(http.DefaultClient, "", "")

func TestBithumb_GetTicker(t *testing.T) {
	return
	ticker, err := bh.GetTicker(gobtb.NewCurrencyPair2("ALL_KAW"))
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}

func TestBithumb_GetDepth(t *testing.T) {
	log.Printf("测试获取bithumb合约行情")
	m_index := map[string]gobtb.CurrencyPair{
		"BTC": gobtb.BTC_KRW,
		"LTC": gobtb.LTC_KRW,
		"EOS": gobtb.EOS_KRW,
		"ETH": gobtb.ETH_KRW,
	}
	for _, views := range m_index {
		dep, err := bh.GetDepth(1, views)
		t.Log("err=>", err)
		t.Log("asks=>", dep.AskList)
		t.Log("bids=>", dep.BidList)
		log.Println(views)
		//盘口
		log.Println(dep.AskList[len(dep.AskList)-1])
		log.Println(dep.BidList[0])
		//log.Printf("Tick[%s],Close[%.2f]", md.Symbol.Code, md.Close)
	}

	return
	dep, err := bh.GetDepth(1, gobtb.EOS_BTC)
	t.Log("err=>", err)
	t.Log("asks=>", dep.AskList)
	t.Log("bids=>", dep.BidList)
	log.Println("res: EOS_KRW")

	log.Println(dep.AskList[len(dep.AskList)-1])
	log.Println(dep.BidList[0])
	return
	dep, err = bh.GetDepth(1, gobtb.EOS_KRW)
	t.Log("err=>", err)
	t.Log("asks=>", dep.AskList)
	t.Log("bids=>", dep.BidList)
	log.Println("res: EOS_KRW")
	log.Println(dep)

	dep, err = bh.GetDepth(1, gobtb.LTC_KRW)
	t.Log("err=>", err)
	t.Log("asks=>", dep.AskList)
	t.Log("bids=>", dep.BidList)
	log.Println("res: LTC_KRW")
	log.Println(dep)
}
