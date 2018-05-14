package main

import (
	"log"

	"github.com/BTB/goBTB/strategy"
)

func main() {
	log.Printf("测试获取合约行情")
	strategy.QryDepthCross()
	//strategy.Test_USD_LTC_USDT()
	//strategy.QryDepth1()
	// log.Printf("测试合成套利")
	// strategy.StrategyOKex3()
}
