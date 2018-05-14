/*
 *
 * @brief XCoin API-call sample script (for Go)
 *
 * @author btckorea
 * @date 2017-04-14
 *
 * @note
 * Make sure current system time is correct.
 * If current system time is not correct, API request will not be processed normally.
 *
 * rdate -s time.nist.gov (if necessary)
 *
 */

package main

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"strconv"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/hex"
	"encoding/base64"
	"encoding/json"
	"crypto/hmac"
	"crypto/sha512"
)


/*
	@fn microsectime
	@brief Return current microseconds.
	@return Returns the current time measured in the number of microseconds(Unix timestamp + microseconds) since the Unix Epoch (January 1 1970 00:00:00 GMT).
*/

func microsectime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond);
}


/*
	@fn hash_hmac
	@brief Generate a keyed hash value using the HMAC method.
	@param hmac_key: Shared secret key used for generating the HMAC variant of the message digest.
	@param hmac_data: Message to be hashed.
	@return Returns a string containing the calculated message digest as base64 encoded hexits.
*/

func hash_hmac(hmac_key string, hmac_data string) (hash_hmac_str string) {
	hmh := hmac.New(sha512.New, []byte(hmac_key));
	hmh.Write([]byte(hmac_data));

	hex_data := hex.EncodeToString(hmh.Sum(nil));
	hash_hmac_bytes := []byte(hex_data);
	hmh.Reset();

	hash_hmac_str = base64.StdEncoding.EncodeToString(hash_hmac_bytes);

	return (hash_hmac_str);
}


// A global variable for use only with the xcoinApiInit() function.

var g_api_key = "";
var g_api_secret = "";


/*
	@fn xcoinApiInit
	@brief Initializes Bithumb API Key and API Secret to use in xcoinApiCall () function.
	@param api_key: Bithumb API Key.
	@param api_secret: Bithumb API Secret.
	@return
*/

func xcoinApiInit(api_key string, api_secret string) {
	g_api_key = api_key;
	g_api_secret = api_secret;
}


/*
	@fn xcoinApiCall
	@brief Call the API with the endpoint URI and POST parameters.
	@param endpoint: Endpoint API URI.
	@param params: POST Parameter.
	@return Returns JSON result as a string from the API server.
*/

func xcoinApiCall(endpoint string, params string) (resp_data_str string) {
	var api_url = "https://api.bithumb.com";

	e_endpoint := url.QueryEscape(endpoint);
	params += "&endpoint=" + e_endpoint;

	hmac_key := g_api_secret;

	// Api-Nonce information generation.
	nonce_int64 := microsectime();
	api_nonce := fmt.Sprint(nonce_int64);

	// Api-Sign information generation.
	hmac_data := endpoint + string(0) + params + string(0) + api_nonce;
	hash_hmac_str := hash_hmac(hmac_key, hmac_data);
	api_sign := hash_hmac_str;

	// Connects to Bithumb API server and returns JSON result value.
	client := &http.Client{};
	http_req, _ := http.NewRequest("POST", api_url + endpoint, bytes.NewBufferString(params)); // URL-encoded payload

	http_req.Header.Add("Api-Key", g_api_key);
	http_req.Header.Add("Api-Sign", api_sign);
	http_req.Header.Add("Api-Nonce", api_nonce);
	http_req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	content_length_str := strconv.Itoa(len(params));
	http_req.Header.Add("Content-Length", content_length_str);

	resp, err := client.Do(http_req);
	if err != nil {
		return ("");
	}

	resp_data, err := ioutil.ReadAll(resp.Body);
	resp_data_str = string(resp_data);

	return (resp_data_str);
}


// /public/ticker structure

type ticker_rec struct {
	Opening_price float64 `json:"opening_price,string"`
	Closing_price float64 `json:"closing_price,string"`
	Sell_price float64 `json:"sell_price,string"`
	Buy_price float64 `json:"buy_price,string"`
}


// Ticker JSON structure

type ticker_json_rec struct {
	Status string `json:"status"`
	Data ticker_rec `json:"data"`
}


// /info/account structure

type account_rec struct {
	Created int64 `json:"created,string"`
	Account_id string `json:"account_id"`
	Trade_fee float64 `json:"trade_fee,string"`
	Balance float64 `json:"balance,string"`
}


// Account JSON structure

type account_json_rec struct {
	Status string `json:"status"`
	Data account_rec `json:"data"`
}


func main() {
	var api_key = "api connect key";
	var api_secret = "api secret key";

	var params = "order_currency=BTC&payment_currency=KRW";
	var resp_data_str = "";
	var resp_data_bytes = []byte("");

	var ticker_json_rec_info ticker_json_rec;
	// var account_json_rec_info account_json_rec;


	xcoinApiInit(api_key, api_secret);

	//
	// Public API
	//
	// /public/ticker
	// /public/recent_ticker
	// /public/orderbook
	// /public/recent_transactions

	fmt.Println("Bithumb Public API URI('/public/ticker') Request...");

	resp_data_str = xcoinApiCall("/public/ticker", params);
	fmt.Printf("%s\n", resp_data_str);

	resp_data_bytes = []byte(resp_data_str);

	json.Unmarshal(resp_data_bytes, &ticker_json_rec_info);

	fmt.Printf("- Status Code: %s\n", ticker_json_rec_info.Status);
	fmt.Printf("- Opening Price: %.8f\n", ticker_json_rec_info.Data.Opening_price);
	fmt.Printf("- Closing Price: %.8f\n", ticker_json_rec_info.Data.Closing_price);
	fmt.Printf("- Sell Price: %.8f\n", ticker_json_rec_info.Data.Sell_price);
	fmt.Printf("- Buy Price: %.8f\n", ticker_json_rec_info.Data.Buy_price);
	fmt.Printf("\n\n");


	//
	// private api
	//
	// endpoint => parameters
	// /info/current
	// /info/account
	// /info/balance
	// /info/wallet_address

	/* fmt.Println("Bithumb Private API URI('/info/account') Request...");

	resp_data_str = xcoinApiCall("/info/account", params);
	fmt.Printf("%s\n", resp_data_str);

	resp_data_bytes = []byte(resp_data_str);

	json.Unmarshal(resp_data_bytes, &account_json_rec_info);

	fmt.Printf("- Status Code: %s\n", account_json_rec_info.Status);
	fmt.Printf("- Created: %d\n", account_json_rec_info.Data.Created);
	fmt.Printf("- Account ID: %s\n", account_json_rec_info.Data.Account_id);
	fmt.Printf("- Trade Fee: %.4f\n", account_json_rec_info.Data.Trade_fee);
	fmt.Printf("- Balance: %.8f\n", account_json_rec_info.Data.Balance); */

	os.Exit(0);
}
