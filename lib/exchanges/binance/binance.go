package binance

import (
	"Cryptotrage/lib/market/market_key"
	"Cryptotrage/lib/order"
	"Cryptotrage/lib/orderbook"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	BASEURL                 = "https://api.binance.com"
	MARKET_ENDPOINT         = "/api/v1/exchangeInfo"
	ORDERBOOK_ENDPOINT      = "/api/v1/depth?symbol="
	DEPOSITADDRESS_ENDPOINT = "/wapi/v3/depositAddress.html"
	BALANCE_ENDPOINT        = "/api/v3/account"
)

type Binance struct { // STANDARDIZED
	apiKey    string
	apiSecret string
	client    *http.Client
}

type Pair struct { // STANDARDIZED
	Ticker string	`json:"symbol"`
	Base string		`json:"baseAsset"`
	Quote string	`json:"quoteAsset"`
}

type Markets struct {
	AssetPairs []Pair `json:"symbols"`
}

type BalanceMeta struct {
	Asset  string  `json:"asset,omitempty"`
	Free   float64 `json:"free,omitempty"`
	Locked float64 `json:"locked,omitempty"`
}

type Wallets struct {
	CanWithdraw bool          `json:"canWithdraw,omitempty"`
	CanDeposit  bool          `json:"canDeposit,omitempty"`
	CanTrade    bool          `json:"canTrade,omitempty"`
	Balances    []BalanceMeta `json:"balances,omitempty"`
	ErrMessage  string        `json:"errMessage,omitempty"`
}

type AddressMeta struct {
	Success    bool   `json:"success,omitempty"`
	CoinTicker string `json:"asset,omitempty"`
	Address    string `json:"address,omitempty"`
	ErrMessage string `json:"errMessage,omitempty"`
}

type OrderBook struct {
	Buys  []Order
	Sells []Order
}

type Order struct {
	Price    interface{} `json:"price"`
	Quantity interface{} `json:"quantity"`
	Volume   interface{} `json:"volume"`
}

func New(key, secret string) *Binance {
	return &Binance{key, secret, &http.Client{}}
}

func (*Binance) GetExchangeName() string {
	return "Binance"
}

func (binance *Binance) GetMarkets() []market_key.Key {
	var markets []market_key.Key
	response, err := http.Get(BASEURL + MARKET_ENDPOINT)
	if err != nil {
		log.Println("HTTP request failed with error %s\n", err)
		return nil
	}
	result:=Markets{}
	if err := json.NewDecoder(response.Body).Decode(&result); err == nil {
		for _, pair := range result.AssetPairs {
			mrkt := market_key.New(pair.Ticker, pair.Base, pair.Quote)
			markets = append(markets, mrkt)
		}
	} else {
		log.Println("Error Decoding")
		return nil
	}
	return markets
}

func (binance *Binance) GetOrderBook(market string, limit int) orderbook.OrderBook {
	sells, buys := make([]order.Order, 0), make([]order.Order, 0)
	response, err := http.Get(BASEURL + ORDERBOOK_ENDPOINT + market)
	if err == nil {
		result := struct {
			Bids [][]interface{} `json:"bids"`
			Asks [][]interface{} `json:"asks"`
		}{}
		bidCounter, askCounter := 0, 0
		if err := json.NewDecoder(response.Body).Decode(&result); err == nil {
			for _, ask := range result.Asks {
				if askCounter <= limit {
					if price, err := strconv.ParseFloat(ask[0].(string), 64); err == nil {
						if quantity, err := strconv.ParseFloat(ask[1].(string), 64); err == nil {
							sells = append(sells, order.New(price, quantity, price * quantity))
						}
					}
				} else {
					break
				}
				askCounter += 1
			}
			for _, bid := range result.Bids {
				if bidCounter <= limit {
					if price, _ := strconv.ParseFloat(bid[0].(string), 64); err == nil {
						if quantity, _ := strconv.ParseFloat(bid[1].(string), 64); err == nil {
							buys = append(buys, order.New(price, quantity, price * quantity))
						}
					}
				} else {
					break
				}
				bidCounter += 1
			}
		} else {
			log.Println("Decoding Error")
			return orderbook.OrderBook{}
		}
	}
	return orderbook.New(market, buys, sells)
}

func (binance *Binance) GetDepositAddress(coinTicker string) AddressMeta {
	params := map[string]string{
		"status":     "true",
		"asset":      coinTicker,
		"recvWindow": strconv.FormatInt(int64(5*time.Second)/int64(time.Millisecond), 10),
		"timestamp":  strconv.FormatInt(time.Now().Unix()*1000, 10),
	}
	response, err := binance.doPrivate("GET", BASEURL+DEPOSITADDRESS_ENDPOINT, params)
	if err != nil {
		log.Printf("The HTTP request failed with error : %s", err)
		return AddressMeta{Success: false, CoinTicker: coinTicker, ErrMessage: err.Error()}
	}
	resData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return AddressMeta{Success: true, CoinTicker: coinTicker, ErrMessage: err.Error()}
	}
	log.Println(string(resData))
	addressMeta:=AddressMeta{}
	if err:=json.Unmarshal(resData,&addressMeta);err!=nil{
		log.Printf("Error Decoding Response : %s",err)
		return AddressMeta{Success:true,CoinTicker:coinTicker,ErrMessage:err.Error()}
	}
	//addressMeta.CoinTicker=coinTicker
	return addressMeta
}

func (binance *Binance) GetBalance(coinTicker string) Wallets {
	params := map[string]string{
		"recvWindow": strconv.FormatInt(int64(5*time.Second)/int64(time.Millisecond), 10),
		"timestamp":  strconv.FormatInt(time.Now().Unix()*1000, 10),
	}
	response, err := binance.doPrivate("GET", BASEURL+BALANCE_ENDPOINT, params)
	if err != nil {
		log.Printf("The HTTP request failed with error : %s", err)
		return Wallets{ErrMessage: err.Error()}
	}

	resData, _ := ioutil.ReadAll(response.Body)
	wallets := Wallets{}
	if err = json.Unmarshal(resData, &wallets); err != nil {
		log.Printf("Error Decoding Response : %s",err)
		return Wallets{ErrMessage: err.Error()}
	}
	return wallets
}

func (binance *Binance) doPrivate(method, path string, queries map[string]string) (*http.Response, error) {
	params := url.Values{}
	for key, value := range queries {
		params.Set(key, value)
	}
	requestUrl := path + "?" + params.Encode()
	requestUrl += fmt.Sprintf("&signature=%s", createSig([]byte(params.Encode()), []byte(binance.apiSecret)))
	req, _ := http.NewRequest(method, requestUrl, nil)
	req.Header.Set("X-MBX-APIKEY", binance.apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func createSig(input, key []byte) string {
	hmac := hmac.New(sha256.New, key)
	hmac.Write(input)
	return hex.EncodeToString(hmac.Sum(nil))
}
