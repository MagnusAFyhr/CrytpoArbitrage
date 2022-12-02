package execute_arb

import (
	"Cryptotrage/exch_arb/src/exch_arbitrage/verify_arb"
)

type ArbitrageExecution struct {
	success bool
	date string
	time string
	duration float64

	parentCoin string
	childCoin string
	results ArbExecutionResults
}
type ArbExecutionResults struct {
	// 1. Start Amount
	// * SMART CONTRACT FEE *
	// 2. First Post Deposit Amount
	// * TRADE FEE *
	// 3. First Post Trade Amount
	// * WITHDRAWAL FEE *
	// 4. First Withdrawal Amount
	// * SMART CONTRACT FEE *
	// 5. Second Post Deposit Amount
	// * TRADE FEE *
	// 6. Second Post Trade Amount
	// * WITHDRAWAL FEE *
	// 7. Second Withdrawal Amount
	// * SMART CONTRACT FEE *
	// 8. End Amount
}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */
func ExecuteArbitrage(verification verify_arb.ArbitrageVerification) ArbitrageExecution {

	return ArbitrageExecution{}
}

/* ************************************************************ */
/*							METHODS								*/
/* ************************************************************ */


/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */