package kraken

import (
	"Cryptotrage/lib/market/market_key"
	"Cryptotrage/lib/order"
	"Cryptotrage/lib/orderbook"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	BASEURL            = "https://api.kraken.com/0"
	MARKET_ENDPOINT    = "/public/AssetPairs"
	ORDERBOOK_ENDPOINT = "/public/Depth?pair="
)

type (
	Kraken struct { // STANDARDIZED
		key    string
		secret string
		client *http.Client
	}
	Markets struct {
		AssetPairs map[string]interface{} `json:"result"`
	}

	Depth struct {
		OrderBooks map[string]interface{} `json:"result"`
	}

	OrderBook struct {
		Buys  []Order
		Sells []Order
	}

	Order struct {
		Price  interface{}
		Volume interface{}
	}

	AddressMeta struct {
		Address      string `json:"address"`
		Expiretm     string `json:"expiretm"`
		New          bool   `json:"new,omitempty"`
		ErrorMessage string `json:"error,omitempty"`
	}
)

func New(key, secret string) *Kraken {
	return &Kraken{key, secret, http.DefaultClient}
}

func (*Kraken) GetExchangeName() string { // STANDARDIZED
	return "Kraken"
}

func (kraken *Kraken) GetMarkets() []market_key.Key { // STANDARDIZED
	var markets []market_key.Key
	response, err := http.Get(BASEURL + MARKET_ENDPOINT)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		return nil
	}

	result := Markets{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error Decoding : " + err.Error())
		return nil
	}
	for _, pair := range result.AssetPairs {
		m := pair.(map[string]interface{})
		if m["wsname"] != nil {
			ticker := m["wsname"].(string)
			s := strings.Split(ticker, "/")
			base, quote := s[0], s[1]
			mrkt := market_key.New(ticker, base, quote)
			markets = append(markets, mrkt)
		}
	}
	return markets
}

func (kraken *Kraken) GetOrderBook(market string, limit int) orderbook.OrderBook {
	sells, buys := make([]order.Order, 0), make([]order.Order, 0)
	response, err := http.Get(BASEURL + ORDERBOOK_ENDPOINT + market + "&count=" + strconv.Itoa(limit))
	if err == nil {
		result := Depth{}
		if err := json.NewDecoder(response.Body).Decode(&result); err == nil {
			for _, value := range result.OrderBooks {
				orderBook := value.(map[string]interface{})
				for _, ask := range orderBook["asks"].([]interface{}) {
					temp := ask.([]interface{})
					price := temp[0].(float64)
					volume := temp[1].(float64)
					quantity := volume / price
					sell := order.New(price, quantity, volume)
					sells = append(sells, sell)
				}
				for _, bid := range orderBook["bids"].([]interface{}) {
					temp := bid.([]interface{})
					price := temp[0].(float64)
					volume := temp[1].(float64)
					quantity := volume / price
					buy := order.New(price, quantity, volume)
					buys = append(buys, buy)
				}
			}
		} else {
			log.Printf("Response Decoding Error : %s",err.Error())
			return orderbook.OrderBook{}
		}
	}else {
		log.Printf("The HTTP request failed with error %s\n", err.Error())
		return  orderbook.OrderBook{}
	}
	return orderbook.New(market, buys, sells)
}

func (kraken *Kraken) DepositMethods(asset string)string {

	response, err := kraken.doPrivate("DepositMethods", url.Values{
		"asset": {asset},
		"nonce": {string(time.Now().UnixNano())},
	})
	if err != nil {
		log.Printf("The HTTP request failed with Error : %s", err.Error())
		return ""
	}
	res, _ := ioutil.ReadAll(response.Body)
	return string(res)
}

func (kraken *Kraken) GetDepositAddress(asset string, method string) AddressMeta {
	response, err := kraken.doPrivate("DepositAddresses", url.Values{
		"asset":  {asset},
		"method": {method},
	})
	if err != nil {
		log.Printf("The HTTP request failed with Error : %s", err.Error())
		return AddressMeta{ErrorMessage: err.Error()}
	}
	res, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(res))
	addressMeta := AddressMeta{}
	if err := json.Unmarshal(res, &addressMeta); err != nil {
		log.Printf("Response Decoding Error  : %s", err.Error())
		return AddressMeta{ErrorMessage:err.Error()}
	}
	return addressMeta

}

func(kraken *Kraken) GetBalance(asset string) string {
	response, err := kraken.doPrivate("Balance", url.Values{})
	if err != nil {
		log.Printf("The HTTP request failed with Error : %s", err.Error())
		return ""
	}
	res, _ := ioutil.ReadAll(response.Body)
	return string(res)
}

func (kraken *Kraken) doPrivate(method string, values url.Values) (*http.Response, error) {
	urlPath := fmt.Sprintf("/private/%s", method)
	reqURL := fmt.Sprintf("%s%s", BASEURL, urlPath)
	secret, _ := base64.StdEncoding.DecodeString(kraken.secret)
	values.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(values.Encode()))
	if err != nil {
		log.Println("Creating Request Error : ", err.Error())
		return nil, err
	}
	req.Header.Add("API-Key", kraken.key)
	req.Header.Add("API-Sign", createSignature(urlPath, values, secret))
	response, err := kraken.client.Do(req)
	if err != nil {
		log.Println("Making Request Error : ", err.Error())
		return nil, err
	}
	return response, nil
}

func createSignature(url string, values url.Values, secret []byte) string {
	shaSum := getSha256([]byte(values.Get("nonce") + values.Encode()))
	macSum := getHMacSha512(append([]byte(url), shaSum...), secret)
	return base64.StdEncoding.EncodeToString(macSum)
}

func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}

func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}
