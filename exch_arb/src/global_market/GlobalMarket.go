package global_market

import (
	"Cryptotrage/lib/market"
)

type GlobalMarket struct {
	globalMarketTicker string
	localMarkets []market.Market
}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */
func New(globalMarketTicker string) GlobalMarket {
	localMarkets := make([]market.Market, 0)
	gm := GlobalMarket{globalMarketTicker, localMarkets}
	return gm
}

/* ************************************************************ */
/*							METHODS								*/
/* ************************************************************ */
func (gm *GlobalMarket) AddLocalMarket(lm market.Market) {//[]market.Market {
	gm.localMarkets = append( gm.localMarkets, lm )
	//return gm.localMarkets
}

/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */
func (gm GlobalMarket) GetGlobalMarketTicker() string {
	return gm.globalMarketTicker
}
func (gm GlobalMarket) GetLocalMarkets() []market.Market {
	return gm.localMarkets
}
