package main

import (
	"Cryptotrage/exch_arb/src/exch_arbitrage/verify_arb"
	"Cryptotrage/exch_arb/src/global_market"
	"Cryptotrage/lib/exchanges/binance"
	"Cryptotrage/lib/exchanges/bittrex"
	"Cryptotrage/lib/exchanges/kraken"
	"Cryptotrage/lib/exchanges/poloniex"
)

var (
	bnc = binance.New(
		"CrgsAcy931vS9DS2VXehkzCqF4eGweD68yHt3pg42dq05WAVIkZ78aO9VJMW8EC7",
		"QSUsaY0T6qPE3ktMyocKFpovNU7UvrYfvh4Ft5TdHboN8HmBo5wijeNHPJv85qs7")
	btrx = bittrex.New(
		"e9de5d0c7dd64e7fa078df6a63c0e79a",
		"9e6d2e8919b742f4a5470c50558c20d8")

	krkn = kraken.New(
		"iGGsepCklkeWSOxKSqmTckma1FEM0/Vo6Rwrj4mO3BTLRhjVyFIXcxSe",
		"vTSvz9APV9iM+tIT0rI3llKah2vGmBbGDGZRbs4JD7YHwzooawmxL9PAO0lZt1Ps7ZZB/ZCUxJYrgmZiMojaLA==")

	plnx = poloniex.New(
		"3W82N24D-LDSB95E7-7PR128OB-EYZVZ7NS",
		"e25c5e397e2f2a6f621e40c950ef9ce3c54b7948eb0cc4e165938e774c63c105db3920122a41a8daaa81ab9a19ed16fb83257e95ff95aff462c60fe7b0216735")
)
func main() {
	// key := "5c2e7a8853b52058333b8ee2"
	// secret := "f9a56526-778c-4f8b-abb1-d789cb3c975c"
	// kucoinTest := kucoin.New(key, secret)


	return
}

// BRAIN FUNCTION
func RunBrain(globalMarkets []global_market.GlobalMarket) {
	for _, globalMarket := range globalMarkets {
		probeGlobalMarket (globalMarket)
	}
}


func probeGlobalMarket(globalMarket global_market.GlobalMarket) {
	for _, a := range globalMarket.GetLocalMarkets() {
		for _, b := range globalMarket.GetLocalMarkets() {
			if a.GetExchangeName() != b.GetExchangeName() {
				verify_arb.VerifyArbitrage(a, b)
			}
		}
	}
}
