package local_exchange

import (
	"Cryptotrage/lib/exchanges/exchange_api"
	"Cryptotrage/lib/market"
	"fmt"
)

type LocalExchange struct {
	exchangeName string
	exchangeAPI exchange_api.ExchangeAPI
	markets []market.Market
}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */
func New( exchangeAPI exchange_api.ExchangeAPI ) LocalExchange {
	exchangeName := exchangeAPI.GetExchangeName()
	markets := make([]market.Market, 0)
	localExchange := LocalExchange{exchangeName, exchangeAPI, markets }
	localExchange.InitMarkets()
	return localExchange
}
func ( le *LocalExchange ) InitMarkets() {
	fmt.Println(">>> Initializing Exchange : " + le.exchangeName)

	keys := le.exchangeAPI.GetMarkets()
	for _, key := range keys {
		mrkt := market.New(key.GetTicker(), key.GetBase(), key.GetQuote(), le.exchangeAPI)
		le.markets = append(le.markets, mrkt)
	}
}

/* ************************************************************ */
/*							METHODS								*/
/* ************************************************************ */


/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */
func ( le LocalExchange ) GetExchangeName() string {
	return le.exchangeName
}
func ( le LocalExchange ) GetMarkets() []market.Market {
	return le.markets
}