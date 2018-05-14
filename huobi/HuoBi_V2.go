package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	. "github.com/BTB/goBTB/base"
)

type HuoBi_V2 struct {
	httpClient *http.Client
	accountId,
	baseUrl,
	accessKey,
	secretKey string
}

type response struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Errmsg  string          `json:"err-msg"`
	Errcode string          `json:"err-code"`
}

func NewV2(httpClient *http.Client, accessKey, secretKey, clientId string) *HuoBi_V2 {
	return &HuoBi_V2{httpClient, clientId, "https://be.huobi.com", accessKey, secretKey}
}

func (hbV2 *HuoBi_V2) GetAccountId() (string, error) {
	path := "/v1/account/accounts"
	params := &url.Values{}
	hbV2.buildPostForm("GET", path, params)

	log.Println(hbV2.baseUrl + path + "?" + params.Encode())

	respmap, err := HttpGet(hbV2.httpClient, hbV2.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return "", err
	}
	//log.Println(respmap)
	if respmap["status"].(string) != "ok" {
		return "", errors.New(respmap["err-code"].(string))
	}

	data := respmap["data"].([]interface{})
	accountIdMap := data[0].(map[string]interface{})
	hbV2.accountId = fmt.Sprintf("%.f", accountIdMap["id"].(float64))

	//log.Println(respmap)
	return hbV2.accountId, nil
}

func (hbV2 *HuoBi_V2) GetAccount() (*Account, error) {
	path := fmt.Sprintf("/v1/account/accounts/%s/balance", hbV2.accountId)
	params := &url.Values{}
	params.Set("accountId-id", hbV2.accountId)
	hbV2.buildPostForm("GET", path, params)

	urlStr := hbV2.baseUrl + path + "?" + params.Encode()
	//println(urlStr)
	respmap, err := HttpGet(hbV2.httpClient, urlStr)

	if err != nil {
		return nil, err
	}

	//log.Println(respmap)

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})
	if datamap["state"].(string) != "working" {
		return nil, errors.New(datamap["state"].(string))
	}

	list := datamap["list"].([]interface{})
	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount, 3)
	acc.Exchange = hbV2.GetExchangeName()

	var (
		cnySubAcc  SubAccount
		bccSubAcc  SubAccount
		etcSubAcc  SubAccount
		ethSubAcc  SubAccount
		btcSubAcc  SubAccount
		ltcSubAcc  SubAccount
		usdtSubAcc SubAccount
	)

	for _, v := range list {
		balancemap := v.(map[string]interface{})
		currency := balancemap["currency"].(string)
		typeStr := balancemap["type"].(string)
		balance := ToFloat64(balancemap["balance"])
		switch currency {
		case "cny":
			cnySubAcc.Currency = CNY
			if typeStr == "trade" {
				cnySubAcc.Amount = balance
			} else {
				cnySubAcc.ForzenAmount = balance
			}
		case "bcc":
			bccSubAcc.Currency = BCC
			if typeStr == "trade" {
				bccSubAcc.Amount = balance
			} else {
				bccSubAcc.ForzenAmount = balance
			}
		case "etc":
			etcSubAcc.Currency = ETC
			if typeStr == "trade" {
				etcSubAcc.Amount = balance
			} else {
				etcSubAcc.ForzenAmount = balance
			}
		case "eth":
			ethSubAcc.Currency = ETH
			if typeStr == "trade" {
				ethSubAcc.Amount = balance
			} else {
				ethSubAcc.ForzenAmount = balance
			}
		case "btc":
			btcSubAcc.Currency = BTC
			if typeStr == "trade" {
				btcSubAcc.Amount = balance
			} else {
				btcSubAcc.ForzenAmount = balance
			}
		case "ltc":
			ltcSubAcc.Currency = LTC
			if typeStr == "trade" {
				ltcSubAcc.Amount = balance
			} else {
				ltcSubAcc.ForzenAmount = balance
			}
		case "usdt":
			usdtSubAcc.Currency = USDT
			if typeStr == "trade" {
				usdtSubAcc.Amount = balance
			} else {
				usdtSubAcc.ForzenAmount = balance
			}
		}
	}

	acc.SubAccounts[CNY] = cnySubAcc
	acc.SubAccounts[BCC] = bccSubAcc
	acc.SubAccounts[ETC] = etcSubAcc
	acc.SubAccounts[ETH] = ethSubAcc
	acc.SubAccounts[BTC] = btcSubAcc
	acc.SubAccounts[USDT] = usdtSubAcc
	acc.SubAccounts[LTC] = ltcSubAcc

	return acc, nil
}

func (hbV2 *HuoBi_V2) placeOrder(amount, price string, pair CurrencyPair, orderType string) (string, error) {
	path := "/v1/order/orders/place"
	params := url.Values{}
	params.Set("account-id", hbV2.accountId)
	params.Set("amount", amount)
	params.Set("symbol", strings.ToLower(pair.ToSymbol("")))
	params.Set("type", orderType)

	switch orderType {
	case "buy-limit", "sell-limit":
		params.Set("price", price)
	}

	hbV2.buildPostForm("POST", path, &params)

	resp, err := HttpPostForm3(hbV2.httpClient, hbV2.baseUrl+path+"?"+params.Encode(), hbV2.toJson(params),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	if err != nil {
		return "", err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return "", err
	}

	if respmap["status"].(string) != "ok" {
		return "", errors.New(respmap["err-code"].(string))
	}

	return respmap["data"].(string), nil
}

func (hbV2 *HuoBi_V2) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbV2.placeOrder(amount, price, currency, "buy-limit")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     BUY}, nil
}

func (hbV2 *HuoBi_V2) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbV2.placeOrder(amount, price, currency, "sell-limit")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     SELL}, nil
}

func (hbV2 *HuoBi_V2) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbV2.placeOrder(amount, price, currency, "buy-market")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     BUY_MARKET}, nil
}

func (hbV2 *HuoBi_V2) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbV2.placeOrder(amount, price, currency, "sell-market")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     SELL_MARKET}, nil
}

func (hbV2 *HuoBi_V2) parseOrder(ordmap map[string]interface{}) Order {
	ord := Order{
		OrderID:    ToInt(ordmap["id"]),
		Amount:     ToFloat64(ordmap["amount"]),
		Price:      ToFloat64(ordmap["price"]),
		DealAmount: ToFloat64(ordmap["field-amount"]),
		Fee:        ToFloat64(ordmap["field-fees"]),
		OrderTime:  ToInt(ordmap["created-at"]),
	}

	state := ordmap["state"].(string)
	switch state {
	case "submitted":
		ord.Status = ORDER_UNFINISH
	case "filled":
		ord.Status = ORDER_FINISH
	case "partial-filled":
		ord.Status = ORDER_PART_FINISH
	case "canceled", "partial-canceled":
		ord.Status = ORDER_CANCEL
	default:
		ord.Status = ORDER_UNFINISH
	}

	if ord.DealAmount > 0.0 {
		ord.AvgPrice = ToFloat64(ordmap["field-cash-amount"]) / ord.DealAmount
	}

	typeS := ordmap["type"].(string)
	switch typeS {
	case "buy-limit":
		ord.Side = BUY
	case "buy-market":
		ord.Side = BUY_MARKET
	case "sell-limit":
		ord.Side = SELL
	case "sell-market":
		ord.Side = SELL_MARKET
	}
	return ord
}

func (hbV2 *HuoBi_V2) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	path := "/v1/order/orders/" + orderId
	params := url.Values{}
	hbV2.buildPostForm("GET", path, &params)
	respmap, err := HttpGet(hbV2.httpClient, hbV2.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})
	order := hbV2.parseOrder(datamap)
	order.Currency = currency
	//log.Println(respmap)
	return &order, nil
}

func (hbV2 *HuoBi_V2) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	path := "/v1/order/orders"
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.ToSymbol("")))
	params.Set("states", "submitted")
	hbV2.buildPostForm("GET", path, &params)
	respmap, err := HttpGet(hbV2.httpClient, fmt.Sprintf("%s%s?%s", hbV2.baseUrl, path, params.Encode()))
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].([]interface{})
	var orders []Order
	for _, v := range datamap {
		ordmap := v.(map[string]interface{})
		ord := hbV2.parseOrder(ordmap)
		ord.Currency = currency
		orders = append(orders, ord)
	}

	//resp, err := HttpPostForm3(hbV2.httpClient, hbV2.baseUrl+path+"?"+params.Encode(), hbV2.toJson(params),
	//	map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	//log.Println(respmap)
	return orders, nil
}

func (hbV2 *HuoBi_V2) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	path := fmt.Sprintf("/v1/order/orders/%s/submitcancel", orderId)
	params := url.Values{}
	hbV2.buildPostForm("POST", path, &params)
	resp, err := HttpPostForm3(hbV2.httpClient, hbV2.baseUrl+path+"?"+params.Encode(), hbV2.toJson(params),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	if err != nil {
		return false, err
	}

	var respmap map[string]interface{}
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	if respmap["status"].(string) != "ok" {
		return false, errors.New(string(resp))
	}

	return true, nil
}

func (hbV2 *HuoBi_V2) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}

func (hbV2 *HuoBi_V2) GetExchangeName() string {
	return "huobi.com"
}

func (hbV2 *HuoBi_V2) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	url := hbV2.baseUrl + "/market/detail/merged?symbol=" + strings.ToLower(currencyPair.ToSymbol(""))
	respmap, err := HttpGet(hbV2.httpClient, url)
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) == "error" {
		return nil, errors.New(respmap["err-msg"].(string))
	}

	tickmap, ok := respmap["tick"].(map[string]interface{})
	if !ok {
		return nil, errors.New("tick assert error")
	}

	ticker := new(Ticker)
	ticker.Vol = ToFloat64(tickmap["amount"])
	ticker.Low = ToFloat64(tickmap["low"])
	ticker.High = ToFloat64(tickmap["high"])
	ticker.Buy = ToFloat64((tickmap["bid"].([]interface{}))[0])
	ticker.Sell = ToFloat64((tickmap["ask"].([]interface{}))[0])
	ticker.Last = ToFloat64(tickmap["close"])
	ticker.Date = ToUint64(respmap["ts"])

	return ticker, nil
}

// GetDepth GetDepth
func (hbV2 *HuoBi_V2) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	//url := hbV2.baseUrl + "/market/depth?symbol=%s&type=step1"
	sd := fmt.Sprintf(hbV2.baseUrl+"/market/depth?symbol=%s&type=step1", strings.ToLower(currency.ToSymbol("")))
	//log.Printf(url)
	log.Printf(sd)
	respmap, err := HttpGet(hbV2.httpClient, sd)
	log.Printf("1")
	//respmap, err := HttpGet(hbV2.httpClient, fmt.Sprintf(url, strings.ToLower(currency.ToSymbol(""))))
	if err != nil {
		log.Printf("3")
		return nil, err
	}
	log.Printf("2")
	if "ok" != respmap["status"].(string) {
		return nil, errors.New(respmap["err-msg"].(string))
	}

	tick, _ := respmap["tick"].(map[string]interface{})
	bids, _ := tick["bids"].([]interface{})
	asks, _ := tick["asks"].([]interface{})

	depth := new(Depth)
	_size := size
	for _, r := range asks {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.AskList = append(depth.AskList, dr)

		_size--
		if _size == 0 {
			break
		}
	}

	_size = size
	for _, r := range bids {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.BidList = append(depth.BidList, dr)

		_size--
		if _size == 0 {
			break
		}
	}

	return depth, nil
}

func (hbV2 *HuoBi_V2) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (hbV2 *HuoBi_V2) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (hbV2 *HuoBi_V2) buildPostForm(reqMethod, path string, postForm *url.Values) error {
	postForm.Set("AccessKeyId", hbV2.accessKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))
	domain := strings.Replace(hbV2.baseUrl, "https://", "", len(hbV2.baseUrl))
	payload := fmt.Sprintf("%s\n%s\n%s\n%s", reqMethod, domain, path, postForm.Encode())
	sign, _ := GetParamHmacSHA256Base64Sign(hbV2.secretKey, payload)
	postForm.Set("Signature", sign)
	return nil
}

func (hbV2 *HuoBi_V2) toJson(params url.Values) string {
	parammap := make(map[string]string)
	for k, v := range params {
		parammap[k] = v[0]
	}
	jsonData, _ := json.Marshal(parammap)
	return string(jsonData)
}
