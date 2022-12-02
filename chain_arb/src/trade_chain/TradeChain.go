package trade_chain

import (
	"Cryptotrage/lib/market"
	"bytes"
	"fmt"
)

type Trade struct {
	tradePair string
	startCoin string
	endCoin string
}

type TradeChain struct {
	baseCoin string
	tradeChain []Trade
}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */
func New( baseCoin string ) TradeChain {
	tradeChain := make([]Trade, 0)
	return TradeChain { baseCoin, tradeChain}
}

/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */
func (chain TradeChain) GetBaseCoin() string {
	return chain.baseCoin
}
func (chain TradeChain) GetLastTrade() Trade {
	chainLen := len(chain.tradeChain)
	if chainLen > 0 {
		lastTrade := chain.tradeChain[chainLen-1]
		return lastTrade
	}
	return Trade{"ERROR","ERROR","ERROR"}
}
func (chain TradeChain) GetTradeChain() []Trade {
	return chain.tradeChain
}
func (chain TradeChain) GetChainAsString() string {
	var chainAsString bytes.Buffer
	chainAsString.WriteString("Trade : "+ chain.baseCoin +"\n")
	for _, trade := range chain.tradeChain {
		chainAsString.WriteString(trade.startCoin)
		chainAsString.WriteString(" -> ")
		chainAsString.WriteString(trade.endCoin)
		chainAsString.WriteString("\n")
	}
	return chainAsString.String()
}

/* ************************************************************ */
/*							METHODS								*/
/* ************************************************************ */
func (chain TradeChain) IsValidNewLink( mrkt market.Market) bool {
	currCoin := chain.baseCoin
	if len(chain.tradeChain) > 0 {
		lastTrade := chain.GetLastTrade()
		currCoin = lastTrade.endCoin
	}
	for i := 1; i < len(chain.tradeChain); i++ {
		trade := chain.tradeChain[i]
		if trade.startCoin == mrkt.GetBaseCoin() || trade.startCoin == mrkt.GetQuoteCoin() {
			return false
		}
	}
	if mrkt.GetBaseCoin() == currCoin || mrkt.GetQuoteCoin() == currCoin {
		return true
	}
	return false
}
func (chain *TradeChain) AddNewLink( mrkt market.Market ) {
	currCoin := chain.baseCoin
	if len(chain.tradeChain) > 0 {
		lastTrade := chain.GetLastTrade()
		currCoin = lastTrade.endCoin
	}
	newTrade := Trade{ }
	if mrkt.GetBaseCoin() == currCoin {
		newTrade = Trade{ mrkt.GetMarketTicker(), mrkt.GetBaseCoin(), mrkt.GetQuoteCoin() }
	} else if mrkt.GetQuoteCoin() == currCoin {
		newTrade = Trade{ mrkt.GetMarketTicker(), mrkt.GetQuoteCoin(), mrkt.GetBaseCoin() }
	} else {
		// ERROR : trade wasn't actually valid
		fmt.Println("ERROR : INVALID NEW TRADE")
	}
	chain.tradeChain = append( chain.tradeChain, newTrade)
}
func (chain TradeChain) IsComplete( ) bool {
	if len(chain.tradeChain) < 2 {
		return false
	}
	lastTrade := chain.GetLastTrade()
	if lastTrade.endCoin == chain.baseCoin {
		return true
	}
	return false
}
func (chain TradeChain) IsValidChain( ) bool {
	currCoin := chain.baseCoin
	// CHECK IF TRADE PATH IS LOGICALLY CORRECT
	for _, link := range chain.tradeChain {
		if link.startCoin == currCoin {
			currCoin = link.endCoin
		} else {
			return false
		}
	}
	// CHECK IF WHAT WE START WITH IS WHAT WE END WITH
	if currCoin != chain.baseCoin {
		return false
	}
	// CHECK IF CHAIN SIZE > 2 BECAUSE : ETH-LTC -> LTC-ETH IS NOT GOING TO BE PROFITABLE
	if len(chain.tradeChain) < 3 {
		return false
	}
	// ALL TESTS PASSED
	return true
}

