package global_exchange

import (
	"Cryptotrage/exch_arb/src/exch_arbitrage/verify_arb"
	"Cryptotrage/exch_arb/src/global_market"
	"Cryptotrage/lib/exchanges/exchange_api"
	"Cryptotrage/lib/local_exchange"
)

type GlobalExchange struct {
	localExchanges []local_exchange.LocalExchange
	globalMarkets []global_market.GlobalMarket
	arbitragedMarkets []global_market.GlobalMarket
}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */

func New( supportedExchanges []exchange_api.ExchangeAPI ) GlobalExchange {
	localExchanges := InitLocalExchanges( supportedExchanges )
	globalMarkets := InitGlobalMarkets( localExchanges )
	arbitragedMarkets := make([]global_market.GlobalMarket, 0)
	ge := GlobalExchange{ localExchanges, globalMarkets, arbitragedMarkets }
	return ge
}
func InitLocalExchanges( supportedExchanges []exchange_api.ExchangeAPI) []local_exchange.LocalExchange {
	localExchanges :=  make( []local_exchange.LocalExchange, 0 )
	for _, exch := range supportedExchanges {
		localExchanges = append( localExchanges, local_exchange.New(exch) )
	}
	return localExchanges
}
func InitGlobalMarkets( localExchanges []local_exchange.LocalExchange ) []global_market.GlobalMarket {
	globalMarkets := make( []global_market.GlobalMarket, 0)
	for _, localExchange := range localExchanges {
		for _, mrkt := range localExchange.GetMarkets() {
			found := false
			for pos, globalMarket := range globalMarkets {
				if mrkt.GetStandardizedTicker() == globalMarket.GetGlobalMarketTicker() {
					found = true
					globalMarket.AddLocalMarket( mrkt )
					globalMarkets[pos] = globalMarket
					break
				}
			}
			if found == false && (mrkt.GetBaseCoin() != "TUSD" || mrkt.GetBaseCoin() != "USDT") && mrkt.GetQuoteCoin() == "ETH" {
				newGlobalMarket := global_market.New( mrkt.GetStandardizedTicker() )
				newGlobalMarket.AddLocalMarket( mrkt )
				globalMarkets = append(globalMarkets, newGlobalMarket)
			}
		}
	}
	adjGlobalMarkets := make( []global_market.GlobalMarket, 0)
	for _, globalMarket := range globalMarkets {
		if len( globalMarket.GetLocalMarkets() ) > 1 {
			adjGlobalMarkets = append( adjGlobalMarkets, globalMarket )
		}
	}
	return adjGlobalMarkets
}

/* ************************************************************ */
/*							METHODS								*/
/* ************************************************************ */
func (ge *GlobalExchange) FindArbitrages( ) { /* Loops through global markets to find potential arbitrage opportunities */
	for _, globalMarket := range ge.globalMarkets {
		for _, buyMarket := range globalMarket.GetLocalMarkets() {
			for _, sellMarket := range globalMarket.GetLocalMarkets() {
				if buyMarket.GetExchangeName() != sellMarket.GetExchangeName() {
					verify_arb.VerifyArbitrage(buyMarket, sellMarket)
				}
			}
		}
	}
}

//func (ge *GlobalExchange) AddArbitrageMarket( market string ) { /* Tries to remove the specified market from arbitraged markets */
//	for _, globalMarket := range ge.globalMarkets {
//		if globalMarket.GetGlobalMarketTicker() == market {
//			ge.arbitragedMarkets = append(ge.arbitragedMarkets, globalMarket)
//			fmt.Printf("%s : Market Found & Added!\n", market)
//		} else {
//			fmt.Printf("%s : Market Not Found & Could Not Be Added!\n", market)
//		}
//	}
//}
//func (ge *GlobalExchange) RemoveArbitrageMarket( market string ) { /* Tries to remove the specified market from arbitraged markets */
//	for i, globalMarket := range ge.globalMarkets { // Cycles through the global markets
//		if globalMarket.GetGlobalMarketTicker() == market { // If the global market shares the market name, remove this market, else do nothing
//			ge.arbitragedMarkets = append(ge.globalMarkets[:i], ge.globalMarkets[i+1:]...)
//			fmt.Printf("%s : Market Found & Removed!\n", market)
//		} else {
//			fmt.Printf("%s : Market Not Found & Could Not Be Removed!\n", market)
//		}
//	}
//}

/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */

func (ge GlobalExchange) GetLocalExchanges() []local_exchange.LocalExchange { /* Returns the local exchanges */
	return ge.localExchanges
}
	func (ge GlobalExchange) GetGlobalMarkets() []global_market.GlobalMarket { /* Returns the global markets */
	return ge.globalMarkets
}
func (ge GlobalExchange) GetArbitragedMarkets() []global_market.GlobalMarket { /* Returns the arbitraged markets */
	return ge.arbitragedMarkets
}