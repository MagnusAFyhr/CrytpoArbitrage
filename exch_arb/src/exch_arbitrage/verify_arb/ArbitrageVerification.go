package verify_arb

import (
	"Cryptotrage/lib/exchanges/exchange_api"
	"Cryptotrage/lib/market"
	"Cryptotrage/lib/order"
	"fmt"
)

type ArbitrageVerification struct {
	valid bool
	arbitrageMeta BasicArbMeta
	executionMeta ArbExecutionMeta
}
type BasicArbMeta struct {
	buyMarket market.Market
	sellMarket market.Market
	whatWePay float64 // INVESTMENT AMOUNT (BASE)
	whatWeGet float64 // ^^^AFTER TRADE FEE (QUOTE)
	whatWeSend float64 // ^^^AFTER WITHDRAWAL FEE (QUOTE)
	whatWeSell float64 // ^^^AFTER TRANSFER FEE (QUOTE - GAS)
	whatWeHave float64 // ^^^AFTER TRADE FEE (BASE)
	whatWeMake float64 // ^^^AFTER WITHDRAWAL FEE (BASE)
}
type ArbExecutionMeta struct {
	arbitrageCost float64

	buyMarketAPI exchange_api.ExchangeAPI
	buyMarketPrice float64
	buyMarketQuantity float64
	depositBuyCoinTicker string
	withdrawalBuyCoinTicker string

	sellMarketAPI exchange_api.ExchangeAPI
	sellMarketPrice float64
	sellMarketQuantity float64
	depositSellCoinTicker string
	withdrawalSellCoinTicker string
}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */
func VerifyArbitrage(buyMarket market.Market, sellMarket market.Market) { // ArbitrageVerification {
	basicArbMeta := CalcArbitrage(buyMarket, sellMarket)
	// arbExecutionMeta := basicArbMeta.GetExecutionMeta()

	//valid := false // CHECK : PROFITABLE, MEETS MIN AMOUNT REQUIREMENTS, VALID EXECUTION META
	//test1 := false
	//test2 := false
	// test3 := false

	// TEST 1 >>> PROFITABILITY TEST (BasicArbMeta)
	roi := basicArbMeta.whatWeMake / basicArbMeta.whatWePay
	if basicArbMeta.whatWeMake > basicArbMeta.whatWePay {
		if roi < 1.02 {
			fmt.Println(" <<< 	NOT ENOUGH OPPORTUNITY	>>> ")
		} else {
			fmt.Println(" <<< 								>>>")
			fmt.Println(" <<<	PROFITABLE ARBITRAGE FOUND	>>> ")
			fmt.Printf(" <<< MARKET        : %s 		>>> \n", buyMarket.GetStandardizedTicker())
			fmt.Printf(" <<< BUY EXCHANGE  : %s 		>>> \n", buyMarket.GetExchangeName())
			fmt.Printf(" <<< SELL EXCHANGE : %s 		>>> \n", sellMarket.GetExchangeName())

			fmt.Printf(" <<< AMOUNT        : %f 		>>> \n", basicArbMeta.whatWeSell)
			fmt.Printf(" <<< COST          : %f 		>>> \n", basicArbMeta.whatWePay)
			fmt.Printf(" <<< REVENUE       : %f 		>>> \n", basicArbMeta.whatWeMake)
			fmt.Printf(" <<< PROFIT        : %f 		>>> \n", basicArbMeta.whatWeMake-basicArbMeta.whatWePay)

			fmt.Println(" <<< 								>>>")
		}
	} else {
		fmt.Println(" <<< 	NO ARBITRAGE OPPORTUNITY	>>> ")
	}

	// TEST 2 >>> MINIMUM AMOUNT REQUIREMENTS TEST (BasicArbMeta)

	// TEST 3 >>> VALID EXECUTION TEST (ArbExecutionMeta)

	//if test1 && test2 && test3 {
	//	valid = true
	//	}

	// return ArbitrageVerification{valid, basicArbMeta,arbExecutionMeta}
}
//func (ba *BasicArbMeta) GetExecutionMeta() ArbExecutionMeta {
//	arbitrageCost := ba.whatWePay // + transferFeeIn
//
//	buyMarketAPI := ba.buyMarket.GetAPI()
//	buyMarketPrice := 0.0
//	buyMarketQuantity := 0.0
//	// depositBuyCoinTicker := ba.buyMarket.GetBaseCoin()
//	// withdrawalBuyCoinTicker := ba.buyMarket.GetQuoteCoin()
//
//	sellMarketAPI := ba.sellMarket.GetAPI()
//	sellMarketPrice := 0.0
//	sellMarketQuantity := 0.0
//	// depositSellCoinTicker := ba.sellMarket.GetQuoteCoin()
//	// withdrawalSellCoinTicker := ba.sellMarket.GetBaseCoin()
//
//	execution := ArbExecutionMeta{ arbitrageCost,
//		buyMarketAPI,buyMarketPrice,buyMarketQuantity,depositBuyCoinTicker, withdrawalBuyCoinTicker,
//	sellMarketAPI,sellMarketPrice,sellMarketQuantity,depositSellCoinTicker,withdrawalSellCoinTicker }
//	return execution
//}

/* ************************************************************ */
/*							METHODS								*/
/* ************************************************************ */
func CalcArbitrage(buyMarket market.Market, sellMarket market.Market) BasicArbMeta {
	buyMarket.LoadOrderBook()
	sellMarket.LoadOrderBook()

	buyOrders := sellMarket.GetOrderBook().GetBuys()
	sellOrders := buyMarket.GetOrderBook().GetSells()

	fmt.Printf("MARKET : %s\n", buyMarket.GetStandardizedTicker())
	whatWePay, whatWeGet, whatWeSend, whatWeSell, whatWeHave, whatWeMake := RecursiveCalcArbitrage(true, buyMarket.GetAPI(), sellMarket.GetAPI(), buyOrders, sellOrders)

	basicArbMeta := BasicArbMeta{ buyMarket,sellMarket,
		whatWePay, whatWeSend, whatWeSell, whatWeGet, whatWeHave, whatWeMake }

	return basicArbMeta
}
func RecursiveCalcArbitrage(first bool, buyerAPI exchange_api.ExchangeAPI, sellerAPI exchange_api.ExchangeAPI, buys []order.Order, sells []order.Order) (payAMT float64, getAMT float64, sendAMT float64, sellAMT float64, haveAMT float64, makeAMT float64) {

	// THIS PROGRAM ASSUMES DEPOSIT FEES ARE 0% & 0.0
	// THIS PROGRAM ASSUMES TRANSFER FEE BETWEEN EXCHANGE A & B is 0.0
	// YOU BUY FROM SELLS & SELL TO BUYS

	// whatWePay -> whatWeGet -> whatWeSend -> whatWeSell -> whatWeHave -> whatWeMake
	// [PAY] -> TRADE FEE -> [GET] -> WITHDRAWAL FEE -> [SEND] -> TRANSFER FEE -> [SELL] -> TRADE FEE -> [HAVE] -> WITHDRAWAL FEE -> [MAKE]

	whatWePay := 0.0  // INVESTMENT AMOUNT // IN PARENT COIN
	whatWeGet := 0.0 // ^^^AFTER TRADE FEE // IN CHILD COIN
	whatWeSend := 0.0 // ^^^AFTER WITHDRAWAL FEE // IN CHILD COIN
	whatWeSell := 0.0 // ^^^AFTER TRANSFER FEE // IN GAS
	whatWeHave := 0.0  // ^^^AFTER TRADE FEE // IN PARENT COIN
	whatWeMake := 0.0 // ^^^AFTER WITHDRAWAL FEE // IN PARENT COIN

	if len(buys) == 0 || len(sells) == 0 {
		return whatWePay, whatWeGet, whatWeSend, whatWeSell, whatWeHave, whatWeMake
	}

	buy := buys[0]
	sell := sells[0]
	buys = append(buys[1:])
	sells = append(sells[1:])

	tradeFeeA := 0.0 // DELETE THIS & REPLACE ALL INSTANCES with apiA.getTradeFee( coinTicker )
	withdrawalFeeA := 0.0 // DELETE THIS & REPLACE ALL INSTANCES with apiA.getWithdrawalFee( coinTicker )
	transferFee := 0.0 // // DELETE THIS & REPLACE ALL INSTANCES with smartContract.getTransferFee( stuff )
	tradeFeeB := 0.0 // DELETE THIS & REPLACE ALL INSTANCES with apiB.getTradeFee( coinTicker )
	withdrawalFeeB := 0.0 // DELETE THIS & REPLACE ALL INSTANCES with apiB.getWithdrawalFee( coinTicker )

	// Adjust price to contain fees : Trade Fee & Withdrawal Fee
	adjSellPrice := (sell.GetVolume()) / (sell.GetQuantity() * (1 - tradeFeeA) * (1 - withdrawalFeeA))
	adjBuyPrice := (buy.GetVolume() * (1 - tradeFeeB) * (1 - withdrawalFeeB)) / (buy.GetQuantity())
	reqSellQuantity := buy.GetQuantity()
	if first {
		reqSellQuantity = ( reqSellQuantity / (1 - withdrawalFeeA) + transferFee) / (1 - tradeFeeA) // AMOUNT OF COINS NEEDED TO FULFILL BUY ORDER (AFTER FEES)
	} else {
		reqSellQuantity = reqSellQuantity / (1 - withdrawalFeeA) / (1 - tradeFeeA) // AMOUNT OF COINS NEEDED TO FULFILL BUY ORDER (AFTER FEES)
	}

	/* //	LOGIC	// */
	if adjSellPrice < adjBuyPrice { /* Profitable */
		// if buy.GetQuantity() > sell.GetQuantity() /* ADJUST FOR FEES */ {
		if reqSellQuantity > sell.GetQuantity() { /* More coins to sell than buy */ /* WILL HAVE INCOMPLETE BUY ORDER */
			whatWePay = sell.GetVolume() // VERIFIED
			whatWeGet = sell.GetQuantity() * (1 - tradeFeeA) // VERIFIED
			whatWeSend = whatWeGet * (1 - withdrawalFeeA) // VERIFIED
			whatWeSell = whatWeSend // VERIFIED
			if first {
				whatWeSell -= transferFee
			}
			// Add Incomplete Buy Order Back To Buy Array
			remBuyQuantity := buy.GetQuantity() - whatWeSell // VERIFIED
			remBuyVolume := buy.GetPrice() * remBuyQuantity // VERIFIED
			remBuy := order.New(buy.GetPrice(), remBuyQuantity, remBuyVolume)
			buys = append([]order.Order{remBuy}, buys...)

			// Calculate Profit
			whatWeHave = (buy.GetVolume() - remBuyVolume) * (1 - tradeFeeB) // VERIFIED
			whatWeMake = whatWeHave * (1 - withdrawalFeeB) // VERIFIED

		//} else if buy.GetQuantity() < sell.GetQuantity() /* ADJUST FOR FEES */ {
		} else if reqSellQuantity < sell.GetQuantity() { /* More coins to buy than sell */ /* WILL HAVE INCOMPLETE SELL ORDER */

			// Add Incomplete Sell Order Back To Sell Array
			remSellQuantity := sell.GetQuantity() - reqSellQuantity
			remSellVolume := sell.GetPrice() * remSellQuantity
			remSell := order.New(sell.GetPrice(), remSellQuantity, remSellVolume)
			sells = append([]order.Order{remSell}, sells...)

			// Calculate Value Variables
			whatWePay = sell.GetVolume() - remSellVolume // sell.GetPrice() * reqSellQuantity
			whatWeGet = reqSellQuantity * (1 - tradeFeeA)
			whatWeSend = whatWeGet * (1 - withdrawalFeeA) // whatWeSell = buy.GetQuantity
			whatWeSell = whatWeSend
			if first {
				whatWeSell -= transferFee
			}
			whatWeHave = buy.GetVolume() * (1 - tradeFeeB)
			whatWeMake = whatWeHave * (1 - withdrawalFeeB)
		} else if reqSellQuantity == sell.GetQuantity() {
			// Calculate Value Variables
			whatWePay = sell.GetVolume()
			whatWeGet = sell.GetQuantity() * (1 - tradeFeeA)
			whatWeSend = whatWeGet * (1 - withdrawalFeeA) // whatWeSell = buy.GetQuantity
			whatWeSell = whatWeSend
			if first {
				whatWeSell -= transferFee
			}
			whatWeHave = buy.GetVolume() * (1 - tradeFeeB)
			whatWeMake = whatWeHave * (1 - withdrawalFeeB)
		}
		// RECURSION UNTIL NOT PROFITABLE
		recWhatWePay, recWhatWeGet, recWhatWeSend, recWhatWeSell, recWhatWeHave, recWhatWeMake := RecursiveCalcArbitrage(false, buyerAPI, sellerAPI, buys, sells)
		whatWePay += recWhatWePay
		whatWeGet += recWhatWeGet
		whatWeSend += recWhatWeSend
		whatWeSell += recWhatWeSell
		whatWeHave += recWhatWeHave
		whatWeMake += recWhatWeMake
	} else { /* Not Profitable */
		fmt.Printf("BUY PRICE : %f\n", buy.GetPrice())
		fmt.Printf("SELL PRICE : %f\n", sell.GetPrice())
		return whatWePay, whatWeGet, whatWeSend, whatWeSell, whatWeHave, whatWeMake
	}
	/* //			  	// */

	return whatWePay, whatWeGet, whatWeSend, whatWeSell, whatWeHave, whatWeMake
}

/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */
