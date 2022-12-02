package market_key

type Key struct {
	marketTicker string // EX : LTC-ETH, LTC_ETH, LTCETH, LTC/ETH
	baseCoin string // EX : LTC
	quoteCoin string // EX : ETH
}
func New(ticker, base, quote string) Key {
	return Key{ticker, base, quote }
}
func (m *Key) GetTicker() string {
	return m.marketTicker
}
func (m *Key) GetBase() string {
	return m.baseCoin
}
func (m *Key) GetQuote() string {
	return m.quoteCoin
}