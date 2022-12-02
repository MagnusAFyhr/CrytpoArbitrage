package chain_manager

import (
	"Cryptotrage/chain_arb/src/trade_chain"
	"Cryptotrage/lib/exchanges/exchange_api"
	"Cryptotrage/lib/local_exchange"
	"Cryptotrage/lib/market"
)

type ChainManager struct {
	localExchange local_exchange.LocalExchange
	arbChains []trade_chain.TradeChain
}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */
func New( exchangeAPI exchange_api.ExchangeAPI, baseCoins []string) ChainManager {
	localExchange := local_exchange.New( exchangeAPI )
	arbChains := make([]trade_chain.TradeChain, 0)
	chainManager := ChainManager{localExchange,arbChains }
	chainManager.InitChains( baseCoins, localExchange.GetMarkets() )
	return chainManager
}
func (cm *ChainManager) InitChains( baseCoins []string, openMarkets []market.Market ) {
	arbChains := make([]trade_chain.TradeChain, 0)
	for _, baseCoin := range baseCoins {
		openMarketsCopy := make([]market.Market, len(openMarkets))
		copy(openMarketsCopy, openMarkets)
		arbChain := trade_chain.New( baseCoin )
		// CALL CRAFT CHAINS
		newArbChains := cm.CraftChains( arbChain, openMarketsCopy)
		// ADD NEW GROUP OF CHAINS TO arbChains
		arbChains = append(arbChains, newArbChains...)
	}
	cm.arbChains = cm.PublishChains( arbChains )
}

/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */
func (cm *ChainManager) GetExchange() local_exchange.LocalExchange {
	return cm.localExchange
}
func (cm *ChainManager) GetChains() []trade_chain.TradeChain {
	return cm.arbChains
}

/* ************************************************************ */
/*							METHODS								*/
/* ************************************************************ */
func (cm *ChainManager) CraftChains(chain trade_chain.TradeChain, openMarkets []market.Market) []trade_chain.TradeChain {
	chains := make([]trade_chain.TradeChain, 0)
	if chain.IsComplete() {
		return append(chains, chain)
	}
	for pos, openMarket := range openMarkets {
		// NEED TO FIND A VALID NEXT TRADE
		if chain.IsValidNewLink( openMarket ) { // NEXT TRADE FOUND
			newArbChain := chain
			newArbChain.AddNewLink( openMarket )
			// REMOVE THE MARKET WE ARE USING FROM MARKETS ARRAY
			openMarketsCopy := openMarkets
			openMarketsCopy[pos] = openMarketsCopy[len(openMarketsCopy)-1]
			openMarketsCopy[len(openMarketsCopy)-1] = market.Market{}
			openMarketsCopy = openMarketsCopy[:len(openMarketsCopy)-1]
			// CALL CRAFT CHAINS
			newArbChains := cm.CraftChains( newArbChain, openMarketsCopy)
			// ADD NEW GROUP OF CHAINS TO arbChains
			chains = append(chains, newArbChains...)
		}
	}
	return chains
}
func (cm *ChainManager) PublishChains(arbChains []trade_chain.TradeChain) []trade_chain.TradeChain { // VERIFIES ALL CHAINS
	verifiedChains := make([]trade_chain.TradeChain, 0)
	for _, arbChain := range arbChains {
		if arbChain.IsValidChain() {
			verifiedChains = append(verifiedChains, arbChain)
		}
	}
	return verifiedChains
}
