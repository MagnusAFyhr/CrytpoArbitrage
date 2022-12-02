package market

import (
	"Cryptotrage/lib/exchanges/exchange_api"
	"Cryptotrage/lib/order"
	"Cryptotrage/lib/orderbook"
)

type Market struct {
	marketTicker string // EX : LTC-ETH, LTC_ETH, LTCETH, LTC/ETH
	standardizedTicker string
	baseCoin string // EX : LTC
	quoteCoin string // EX : ETH
	exchangeName string // EX : Binance
	api exchange_api.ExchangeAPI
	orderBook orderbook.OrderBook
}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */
func New( marketTicker string, baseCoin string, quoteCoin string, api exchange_api.ExchangeAPI) Market {
	standardizedTicker := baseCoin + "-" + quoteCoin
	exchangeName := api.GetExchangeName()

	buys := make([]order.Order, 0)
	sells := make([]order.Order, 0)
	orderBook := orderbook.New(exchangeName, buys, sells)

	// return Market{marketTicker, baseCoin, quoteCoin, exchangeName, api, orderBook }
	return Market{marketTicker, standardizedTicker, baseCoin,quoteCoin, exchangeName, api, orderBook }
}
/* ************************************************************ */
/*							METHODS								*/
/* ************************************************************ */
func (m *Market) LoadOrderBook() {
	m.orderBook = m.api.GetOrderBook(m.marketTicker, 100)
}
/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */
func (m Market) GetMarketTicker() string {
	return m.marketTicker
}
func (m Market) GetStandardizedTicker() string {
	return m.standardizedTicker
}
func (m Market) GetBaseCoin() string {
	return m.baseCoin
}
func (m Market) GetQuoteCoin() string {
	return m.quoteCoin
}
func (m Market) GetExchangeName() string {
	return m.exchangeName
}
func (m Market) GetAPI() exchange_api.ExchangeAPI {
	return m.api
}
func (m Market) GetOrderBook() orderbook.OrderBook {
	return m.orderBook
}