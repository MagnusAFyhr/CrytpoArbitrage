package poloniex

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	BASEURL            = "https://poloniex.com/" // Bittrex API endpoint
	MARKET_ENDPOINT    = "public?command=returnTicker"
	ORDERBOOK_ENDPOINT = "public?command=returnOrderBook&currencyPair="
	PRIVATE            = "tradingApi"
)

type (
	Poloniex struct { // ExchangeName = Kucoin, Binance, Bittrex, Yobit, etc...
		apiKey    string
		apiSecret string
		client    *http.Client
	}

	OrderBook struct {
		Buys   []Order
		Sells  []Order
	}

	Order struct {
		Price    interface{} `json:"price"`
		Quantity interface{} `json:"quantity"`
		Volume   interface{} `json:"volume"`
	}

	BalanceMeta struct {
		success    bool
		coinTicker string
		balance    float64
		errMessage string
	}

	AddressMeta struct {
		success bool
		coinTicker string
		address string
		errMessage string
	}
)

func New(key, secret string) *Poloniex {
	return &Poloniex{key, secret, &http.Client{}}

}

func (*Poloniex) GetExchangeName() string {
	return "Poloniex"
}

func (*Poloniex) GetMarkets() []string {
	res, err := http.Get(BASEURL + MARKET_ENDPOINT)
	if err != nil {
		log.Println("Error in connecing: ", err)
		return nil
	}
	resp, err := ioutil.ReadAll(res.Body)
	var bodyDataMap map[string]interface{}
	if err := json.Unmarshal(resp, &bodyDataMap); err != nil {
		log.Println("Error in unmarshaling: ", err)
		return nil
	}
	markets := make([]string, 0)
	for mkt, _ := range bodyDataMap {
		markets = append(markets, mkt)
	}
	return markets
}
func (*Poloniex) GetOrderBook(market string, limit int) *OrderBook {
	response, err := http.Get(BASEURL + ORDERBOOK_ENDPOINT + market + "&depth=" + strconv.Itoa(limit))
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		return nil
	}
	resData, _ := ioutil.ReadAll(response.Body)
	var orderbook map[string]interface{}
	if err := json.Unmarshal(resData, &orderbook); err != nil {
		log.Println("Error in Decoding : ", err)
		return nil
	}
	sells, bids := make([]Order, 0), make([]Order, 0)
	for _, ask := range orderbook["asks"].([]interface{}) {
		order := ask.([]interface{})
		sells = append(sells, Order{Price: order[0], Quantity: order[1]})
	}
	for _, bid := range orderbook["bids"].([]interface{}) {
		order := bid.([]interface{})
		bids = append(bids, Order{Price: order[0], Quantity: order[1]})
	}
	return &OrderBook{Buys: bids, Sells: sells}
}

func (poloniex *Poloniex) GetDepositAddress(coinTicker string) AddressMeta{
	address:=""
	response, err := poloniex.doPrivate("POST", url.Values{
		"command": {"returnDepositAddresses"},
		"nonce":   {strconv.FormatInt(time.Now().Unix()*1000, 10)},
	})
	if err != nil {
		log.Printf("The HTTP request failed with Error : %s", err.Error())
		return AddressMeta{success:false,coinTicker:coinTicker,errMessage:err.Error()}
	}

	defer response.Body.Close()
	resData, _ := ioutil.ReadAll(response.Body)
	log.Println(string(resData))
	arr := strings.Split(string(resData), ",")
	for _, value := range arr {
		if strings.Split(value, ":")[0] == "\""+coinTicker+"\"" {
			address= strings.Split(value, ":")[1][1:len(strings.Split(value, ":")[1])-1]
		}
	}
	return AddressMeta{success:true,coinTicker:coinTicker,address:address}
}

func (poloniex *Poloniex) GetBalance(coinTicker string) BalanceMeta {
	balance := float64(0)
	response, err := poloniex.doPrivate("POST", url.Values{
		"command": {"returnBalances"},
		"nonce":   {strconv.FormatInt(time.Now().Unix()*1000, 10)},
	})
	if err != nil {
		log.Printf("The HTTP request failed with Error : %s", err.Error())
		return BalanceMeta{success:false,coinTicker: coinTicker,errMessage: err.Error()}
	}
	defer response.Body.Close()
	resData, _ := ioutil.ReadAll(response.Body)
	log.Println(string(resData))
	arr := strings.Split(string(resData), ",")
	for _, value := range arr {
		if strings.Split(value, ":")[0] == "\""+coinTicker+"\"" {
			balance, _ = strconv.ParseFloat(strings.Split(value, ":")[1][1:len(strings.Split(value, ":")[1])-1], 64)
		}
	}

	return BalanceMeta{success:true,coinTicker:coinTicker,balance:balance}
}

func (poloniex *Poloniex) doPrivate(method string, values url.Values) (*http.Response, error) {
	reqURL := BASEURL + PRIVATE + "?" + values.Encode()
	req, err := http.NewRequest(method, reqURL, bytes.NewBuffer([]byte(values.Encode())))
	if err != nil {
		log.Println("Creating Request Error : ", err.Error())
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Key", poloniex.apiKey)
	req.Header.Set("Sign", createSignature([]byte(values.Encode()), []byte(poloniex.apiSecret)))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Making Request Error : ", err.Error())
		return nil, err
	}
	return response, nil
}

func createSignature(message, secret []byte) string {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return hex.EncodeToString(mac.Sum(nil))
}
