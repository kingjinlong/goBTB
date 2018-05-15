package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/pflag"
)

// AppConfig AppConfig
type AppConfig struct {
	BITapiKey    string
	BITsecretKey string
}

var configPath string
var appConfig AppConfig

func parseConf(configPath string) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &appConfig); err != nil {
		log.Println(err, data)
		return err
	}
	log.Println(appConfig)
	return nil
}

// AddFlags AddFlags
func AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&configPath, "config", "trade.json", "server port")
	fs.Parse(os.Args[1:])
}

func main() {
	log.Printf("测试获取交易策略")

	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	// "go run main.go --config=../../config/trade.json"
	AddFlags(pflag.CommandLine)
	err := parseConf(configPath)
	if err != nil {
		log.Printf("Get Engine Error: %v\n", err)
		os.Exit(0)
	}
	log.Printf("---")
	log.Println(appConfig)

	//strategy.StrategyHuobiOkex()
	//strategy.Test_USD_LTC_USDT()
	//strategy.QryDepth1()
	// log.Printf("测试合成套利")
	// strategy.StrategyOKex3()
}
