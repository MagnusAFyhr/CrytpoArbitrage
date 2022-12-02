package kucoin

import (
	"Cryptotrage/lib/order"
	"Cryptotrage/lib/orderbook"
	"fmt"
	"github.com/eeonevision/kucoin-go-master"
)

type V1Kucoin struct {
	exchangeName string
	apiKey string
	apiSecret string
}

func New( apiKey string, apiSecret string ) V1Kucoin {
	exchangeName := "Kucoin"
	client := V1Kucoin{ exchangeName,apiKey, apiSecret }
	return client
}

func (k V1Kucoin) GetExchangeName() string {
	return k.exchangeName
}

func (k V1Kucoin) GetMarkets() []string {
	markets := make([]string, 0)

	/* //		LOGIC		// */
	client := kucoin.New( k.apiKey, k.apiSecret )

	symbols, err1 := client.GetSymbols()
	if err1 != nil {
		fmt.Println(err1)
	}

	for _, symbol := range symbols {
		marketTicker := symbol.CoinType + "-" + symbol.CoinTypePair
		// mrkt := market.New(marketTicker, k.exchangeName, k)
		markets = append(markets, marketTicker)
	}
	/* //		-----		// */

	return markets
}

func (k V1Kucoin) GetOrderBook(market string, limit int) orderbook.OrderBook {
	buys := make([]order.Order, 0)
	sells := make([]order.Order, 0)

	/* //		LOGIC		// */
	client := kucoin.New( k.apiKey, k.apiSecret )

	coinOrders, err1 := client.OrdersBook(market, 0, limit)
	if err1 != nil {
		fmt.Println(err1)
	}

	for _, bid := range coinOrders.BUY {
		price := bid[0]
		quantity := bid[1]
		volume := bid[2]

		buy := order.New(price, quantity, volume)
		buys = append(buys, buy)
	}
	for _, ask := range coinOrders.SELL {
		price := ask[0]
		quantity := ask[1]
		volume := ask[2]

		sell := order.New(price, quantity, volume)
		sells = append(sells, sell)
	}
	/* //		-----		// */

	standardizedOrders := orderbook.New( market, buys, sells )
	return standardizedOrders
}