package main

import (
	"Cryptotrage/exch_arb/src/global_exchange"
	"Cryptotrage/lib/exchanges/binance"
	"Cryptotrage/lib/exchanges/bittrex"
	"Cryptotrage/lib/exchanges/exchange_api"
	"Cryptotrage/lib/exchanges/kraken"
	"Cryptotrage/lib/exchanges/poloniex"
	"encoding/json"
	"fmt"
)

var (
	bnc = binance.New(
		"CrgsAcy931vS9DS2VXehkzCqF4eGweD68yHt3pg42dq05WAVIkZ78aO9VJMW8EC7",
		"QSUsaY0T6qPE3ktMyocKFpovNU7UvrYfvh4Ft5TdHboN8HmBo5wijeNHPJv85qs7")
	btrx = bittrex.New(
		"e9de5d0c7dd64e7fa078df6a63c0e79a",
		"9e6d2e8919b742f4a5470c50558c20d8")

	krkn = kraken.New(
		"iGGsepCklkeWSOxKSqmTckma1FEM0/Vo6Rwrj4mO3BTLRhjVyFIXcxSe",
		"vTSvz9APV9iM+tIT0rI3llKah2vGmBbGDGZRbs4JD7YHwzooawmxL9PAO0lZt1Ps7ZZB/ZCUxJYrgmZiMojaLA==")

	plnx = poloniex.New(
		"3W82N24D-LDSB95E7-7PR128OB-EYZVZ7NS",
		"e25c5e397e2f2a6f621e40c950ef9ce3c54b7948eb0cc4e165938e774c63c105db3920122a41a8daaa81ab9a19ed16fb83257e95ff95aff462c60fe7b0216735")
)

func main() {
	//book := bnc.GetOrderBook("ETHBTC", 3)
	//for _, order := range book.Buys {
	//	fmt.Println(order.Volume)
	//}
	TestArbitrage()

	// BinanceMarkets()
	//BinanceOrderBook("ETHBTC",100)
	//BinanceGetBalance("ETH")
	//BinanceGetDepositAddress("ETH")

	//BittrexMarkets()
	//BittrexOrderBook("ETH-TUSD",3)
	//BittrexGetBalance("ETH")
	//BittrexGetDepositAddress("ETH")

	//krkn.GetMarkets()
	//KrakenMarkets()
	//KrakenOrderBook("ADACAD",4)
	//KrakenGetDepositMethods("BTC")
	//KrakenGetDepositAddress("BTC","corressponding method from above function")
	//KrakenGetBalance("BTC")


	//PoloniexMarkets()
	//PoloniexOrderBook("BTC_ETH",3)
	//PoloniexGetDepositAddress("ETH")
	//PoloniexGetBalance("ETH")




}

func TestArbitrage() {
	suppExchanges := make([]exchange_api.ExchangeAPI, 0)
	suppExchanges = append(suppExchanges, btrx)
	suppExchanges = append(suppExchanges, bnc)

	ge := global_exchange.New(suppExchanges)
	ge.FindArbitrages()
	//gms := ge.GetGlobalMarkets()
	//for _, gm := range gms {
	//	lms := gm.GetLocalMarkets()
	//	verify_arb.CalcArbitrage(lms[0], lms[1])
	//}
}

func BinanceMarkets() {
	markets := bnc.GetMarkets()
	for _, market := range markets {
		fmt.Println(market)
	}

}

func BinanceOrderBook(market string, limit int) {
	orderBook := bnc.GetOrderBook(market, limit)
	fmt.Println("BUYS")
	for _, buy := range orderBook.GetBuys() {
		fmt.Println(buy.GetPrice())
		fmt.Println(buy.GetQuantity())
		fmt.Println(buy.GetVolume())
	}
	fmt.Println("SELLS")
	for _, buy := range orderBook.GetBuys() {
		fmt.Println(buy.GetPrice())
		fmt.Println(buy.GetQuantity())
		fmt.Println(buy.GetVolume())
	}
}

func BinanceGetDepositAddress(coin string) {
	addressMeta := bnc.GetDepositAddress(coin)
	jsonFormat, _ := json.Marshal(addressMeta)
	fmt.Println(string(jsonFormat))
}

func BinanceGetBalance(coin string) {
	balanceMeta := bnc.GetBalance(coin)
	jsonFormat, _ := json.Marshal(balanceMeta)
	fmt.Println(string(jsonFormat))

}





func BittrexMarkets() {
	markets := btrx.GetMarkets()
	for _, market := range markets {
		fmt.Println(market)
	}
}

func BittrexOrderBook(market string, limit int) {
	orderBook := btrx.GetOrderBook(market, limit)
	fmt.Println("BUYS")
	for _, buy := range orderBook.GetBuys() {
		fmt.Println(buy.GetPrice())
		fmt.Println(buy.GetQuantity())
		fmt.Println(buy.GetVolume())
	}
	fmt.Println("SELLS")
	for _, buy := range orderBook.GetBuys() {
		fmt.Println(buy.GetPrice())
		fmt.Println(buy.GetQuantity())
		fmt.Println(buy.GetVolume())
	}
}

func BittrexGetDepositAddress(coin string) {
	addressMeta := btrx.GetDepositAddress(coin)
	jsonFormat, _ := json.Marshal(addressMeta)
	fmt.Println(string(jsonFormat))
}

func BittrexGetBalance(coin string) {
	balanceMeta := btrx.GetBalance(coin)
	jsonFormat, _ := json.Marshal(balanceMeta)
	fmt.Println(string(jsonFormat))
}




func KrakenMarkets() {
	markets := krkn.GetMarkets()
	for _, market := range markets {
		fmt.Println(market.GetTicker())
		fmt.Println(market.GetBase())
		fmt.Println(market.GetQuote())
		fmt.Println()
	}
}

func KrakenOrderBook(market string, limit int) {
	orderBook := krkn.GetOrderBook(market, limit)
	for _, order := range orderBook.GetSells() {
		ordr, _ := json.Marshal(order)
		fmt.Println("Sell => ",string(ordr))
	}
	for _, order := range orderBook.GetBuys() {
		ordr, _ := json.Marshal(order)
		fmt.Println("Buy => ",string(ordr))
	}
}

func KrakenGetDepositMethods(coin string)  {
	fmt.Println(krkn.DepositMethods(coin))
}

func KrakenGetDepositAddress(coin ,method string) {
	addressMeta := krkn.GetDepositAddress(coin,method)
	jsonFormat, _ := json.Marshal(addressMeta)
	fmt.Println(string(jsonFormat))
}

func KrakenGetBalance(coin string) {
	balanceMeta := krkn.GetBalance(coin)
	jsonFormat, _ := json.Marshal(balanceMeta)
	fmt.Println(string(jsonFormat))
}




func PoloniexMarkets() {
	markets := plnx.GetMarkets()
	for _, market := range markets {
		fmt.Println(market)
	}
}

func PoloniexOrderBook(market string, limit int) {
	orderBook := plnx.GetOrderBook(market, limit)
	for _, order := range orderBook.Sells {

		ordr, _ := json.Marshal(order)
		fmt.Println("Sell => ",string(ordr))
	}
	for _, order := range orderBook.Buys {
		ordr, _ := json.Marshal(order)
		fmt.Println("Buy => ",string(ordr))
	}
}

func PoloniexGetDepositAddress(coin string) {
	addressMeta := plnx.GetDepositAddress(coin)
	jsonFormat, _ := json.Marshal(addressMeta)
	fmt.Println(string(jsonFormat))
}

func PoloniexGetBalance(coin string) {
	balanceMeta := plnx.GetBalance(coin)
	jsonFormat, _ := json.Marshal(balanceMeta)
	fmt.Println(string(jsonFormat))
}

