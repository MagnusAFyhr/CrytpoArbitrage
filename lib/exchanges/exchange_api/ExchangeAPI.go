package exchange_api

import (
	"Cryptotrage/lib/market/market_key"
	"Cryptotrage/lib/orderbook"
)

type ExchangeAPI interface {
	GetExchangeName() string
	GetMarkets() []market_key.Key
	GetOrderBook( market string, limit int ) orderbook.OrderBook

	// GetDepositAddress( coinTicker string ) AddressMeta
	// GetBalance( coinTicker string ) BalanceMeta

	// GetOrder( orderID string ) OrderMeta
	// GetOpenOrders( market string, orderType string ) []OrderMeta
	// CreateOrder( market string, orderType string, price float64, quantity float64 ) OrderMeta
	// CancelOrder( orderID string ) OrderMeta

	// GetDepositFee( coinTicker string, amount float64 ) float64
	// GetTradeFee( marketTicker string, amount float64 ) bool, float64
	// GetWithdrawFee ( coinTicker string, amount float64 ) float64
}

type AddressMeta struct { // SATISFIED
	success bool // IF CALL WAS SUCCESSFUL
	coinTicker string // CURRENCY (EX : ETH)
	address float64 // ADDRESS OF THIS CURRENCY
	errMessage string // WHY CALL WAS NOT SUCCESSFUL
}

type BalanceMeta struct { // SATISFIED
	success bool // IF CALL WAS SUCCESSFUL
	coinTicker string // CURRENCY (EX : ETH)
	balance float64 // CURRENT BALANCE OF THIS CURRENCY
	errMessage string // WHY CALL WAS NOT SUCCESSFUL
}

type OrderMeta struct { // SATISFIED
	success bool // IF CALL WAS SUCCESSFUL
	orderID string // ID ON THE ORDER
	market string // ASSET PAIR (EX : LTC-ETH)
	orderType string // BUY, SELL
	startAmount float64 // AMOUNT AT BEGINNING OF ORDER
	pendingAmount float64 // PENDING AMOUNT OF ORDER
	rate float64 // PRICE OF BUYING/SELLING
	// opened DateTime // DATETIME WHEN THE ORDER WAS PLACED
	// closed DateTime // DATETIME WHEN THE ORDER WAS FINISHED/CANCELLED
	status int // 0 = ACTIVE, 1 = FULFILLED & CLOSED, 2 = CANCELLED, 3 = PARTIALLY FULFILLED & CANCELLED
	errMessage string // WHY CALL WAS NOT SUCCESSFUL
}
