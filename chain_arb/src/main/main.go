package main

import (
	"Cryptotrage/chain_arb/src/chain_manager"
	"Cryptotrage/chain_arb/src/exchanges/kucoin"
	"fmt"
)

func main() {
	key := "5c2e7a8853b52058333b8ee2"
	secret := "f9a56526-778c-4f8b-abb1-d789cb3c975c"
	kucoinTest := kucoin.New(key, secret)

	baseCoins := make([]string, 0)
	baseCoins = append(baseCoins,"BTC", "ETH")

	chainManager := chain_manager.New(kucoinTest, baseCoins)
	chains := chainManager.GetChains()
	fmt.Println(len(chains))
	for _, chain := range chains {
		chainAsString := chain.GetChainAsString()
		fmt.Println(chainAsString)
	}
	return
}