package orderbook

import "Cryptotrage/lib/order"

type OrderBook struct {
	market string // EX : LTC-ETH
	buys []order.Order
	sells []order.Order
}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */
func New(market string, buys []order.Order, sells []order.Order) OrderBook {
	orders := OrderBook { market, buys, sells }
	return orders
}

/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */
func (book OrderBook) GetMarket() string {
	return book.market
}
func (book OrderBook) GetBuys() []order.Order {
	return book.buys
}
func (book OrderBook) GetSells() []order.Order {
	return book.sells
}