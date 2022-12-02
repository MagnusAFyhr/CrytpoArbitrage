package bittrex

import (
	"Cryptotrage/lib/market/market_key"
	"Cryptotrage/lib/order"
	"Cryptotrage/lib/orderbook"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	BASEURL                 = "https://api.bittrex.com/api/v1.1"
	MARKET_ENDPOINT         = "/public/getmarkets"
	ORDERBOOK_ENDPOINT      = "/public/getorderbook?market="
	DEPOSITADDRESS_ENDPOINT = "/account/getdepositaddress?apikey="
	BALANCE_ENDPOINT        = "/account/getbalance?apikey="
)

type Bittrex struct { // STANDARDIZED
	apiKey    string
	apiSecret string
	client    *http.Client
}

type Pair struct { // STANDARDIZED
	Ticker string `json:"MarketName"`
	Base string `json:"MarketCurrency"`
	Quote string `json:"BaseCurrency"`
}

type Markets struct {
	AssetPairs []Pair `json:"result"`
}

type BalanceMeta struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  struct {
		Currency  string  `json:"Currency"`
		Balance   float64 `json:"Balance"`
		Available float64 `json:"Available"`
		Pending   float64 `json:"Pending"`
	} `json:"result"`
}

type AddressMeta struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  struct {
		Currency string `json:"Currency"`
		Address  string `json:"Address"`
	} `json:"result"`
}

type OrderBookMeta struct {
	Orders OrderBook `json:"result"`
}

type OrderBook struct {
	Buys  []Order `json:"buy"`
	Sells []Order `json:"sell"`
}

type Order struct {
	Price    float64 `json:"rate"`
	Quantity float64 `json:"quantity"`
}

func New(key, secret string) *Bittrex {
	return &Bittrex{key, secret, http.DefaultClient}
}

func (*Bittrex) GetExchangeName() string {
	return "Bittrex"
}

func (b *Bittrex) GetMarkets() []market_key.Key {
	markets := make([]market_key.Key, 0)
	response, err := http.Get(BASEURL + MARKET_ENDPOINT)
	if err != nil {
		log.Printf("The HTTP request failed with error : %s", err.Error())
		return nil
	}
	result := Markets{}
	if err := json.NewDecoder(response.Body).Decode(&result); err == nil {
		for _, pair := range result.AssetPairs {
			mrkt := market_key.New(pair.Ticker, pair.Base, pair.Quote)
			markets = append(markets, mrkt)
		}
	} else {
		log.Printf("Error Decoding Response : %s",err.Error())
		return nil
	}
	return markets
}

func (b *Bittrex) GetOrderBook(market string, limit int) orderbook.OrderBook {
	response, err := http.Get(BASEURL + ORDERBOOK_ENDPOINT + market + "&type=both")
	if err != nil {
		log.Printf("The HTTP request failed with error : %s", err.Error())
		return orderbook.OrderBook{}
	}
	result:=OrderBookMeta{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Printf("Error Decoding Response : %s",err.Error())
		return orderbook.OrderBook{}
	}
	// STANDARDIZING
	if len(result.Orders.Buys) < limit {
		limit = len(result.Orders.Buys)
	}
	if len(result.Orders.Sells) < limit {
		limit = len(result.Orders.Sells)
	}
	buys, sells := make([]order.Order, 0), make([]order.Order, 0)
	for _, buy := range result.Orders.Buys[:limit] {
		stanBuy := order.New(buy.Price, buy.Quantity, buy.Price * buy.Quantity)
		buys = append(buys, stanBuy)
	}
	for _, sell := range result.Orders.Sells[:limit] {
		stanSell := order.New(sell.Price, sell.Quantity, sell.Price * sell.Quantity)
		sells = append(buys, stanSell)
	}
	return orderbook.New(market, buys, sells)
}

func (b *Bittrex) GetDepositAddress(coinTicker string) AddressMeta {
	resp, err := b.doPrivate("GET", BASEURL+DEPOSITADDRESS_ENDPOINT, coinTicker)
	if err != nil {
		log.Printf("The HTTP request failed with error : %s", err.Error())
		return AddressMeta{Success: false, Message: err.Error()}
	}
	defer resp.Body.Close()
	resData, _ := ioutil.ReadAll(resp.Body)
	addressMeta := AddressMeta{}
	if err := json.Unmarshal(resData,&addressMeta); err != nil {
		log.Printf("Error Decoding Response : %s",err.Error())
		return AddressMeta{Success: true, Message: err.Error()}
	}
	return addressMeta

}

func (b *Bittrex) GetBalance(coinTicker string) BalanceMeta {
	resp, err := b.doPrivate("GET", BASEURL+BALANCE_ENDPOINT, coinTicker)
	if err != nil {
		log.Printf("The HTTP request failed with error : %s", err.Error())
		return BalanceMeta{Success: false, Message: err.Error()}
	}
	defer resp.Body.Close()

	resData, _ := ioutil.ReadAll(resp.Body)
	balanceMeta := BalanceMeta{}
	if err := json.Unmarshal(resData,&balanceMeta); err != nil {
		log.Printf("Error Decoding Response : %s",err.Error())
		return BalanceMeta{Success: false, Message: err.Error()}
	}
	return balanceMeta
}

func (b *Bittrex) doPrivate(method, path, ticker string) (*http.Response, error) {
	req, err := http.NewRequest(method, path+b.apiKey+"&currency="+ticker, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	nonce := time.Now().UnixNano()
	q := req.URL.Query()
	q.Set("nonce", fmt.Sprintf("%d", nonce))
	req.URL.RawQuery = q.Encode()
	req.Header.Add("apisign", createSig(b.apiSecret, req))
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func createSig(secret string, req *http.Request) string {
	mac := hmac.New(sha512.New, []byte(secret))
	_, err := mac.Write([]byte(req.URL.String()))
	if err != nil {
		log.Println("Error Creating Signature : ", err)
		return ""
	}
	return hex.EncodeToString(mac.Sum(nil))
}
