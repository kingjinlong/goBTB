package huobi

import (
	"log"
	"net/http"
	"testing"

	gobtb "github.com/BTB/goBTB/base"
	"github.com/stretchr/testify/assert"
)

// var hb = New(http.DefaultClient, "", "")
var hbpro = NewHuobiPro(http.DefaultClient, "", "", "")

func TestHuoBi_GetDepth(t *testing.T) {
	log.Println("test huobi")
	//dep, err := hb.GetDepth(2, gobtb.BTC_CNY)
	dep, err := hbpro.GetDepth(2, gobtb.BTC_USDT)
	assert.Nil(t, err)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
	log.Println(dep.AskList[len(dep.AskList)-1])
	log.Println(dep.BidList[0])
}

// func TestHuoBi_GetKlineRecords(t *testing.T) {
// 	return
// 	klines, err := hb.GetKlineRecords(gobtb.BTC_CNY, gobtb.KLINE_PERIOD_4H, 1, -1)
// 	assert.Nil(t, err)
// 	t.Log(klines)
// }
