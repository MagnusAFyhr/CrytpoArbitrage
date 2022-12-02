package record_arb

import (
	"Cryptotrage/exch_arb/src/exch_arbitrage/execute_arb"
	"Cryptotrage/exch_arb/src/exch_arbitrage/verify_arb"
)

type ArbitrageRecord struct {

}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */
func New(verification verify_arb.ArbitrageVerification, execution execute_arb.ArbitrageExecution) ArbitrageRecord {


	return ArbitrageRecord{ }
}

/* ************************************************************ */
/*							METHODS								*/
/* ************************************************************ */


/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */
func (arb *ArbitrageRecord) GetRecordAsString() string {
	statement := ""

	return statement
}