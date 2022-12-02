package order

type Order struct {
	price float64 // rate in base / quote
	quantity float64 // amount of quote
	volume float64 // amount of base
}

/* ************************************************************ */
/*							INIT								*/
/* ************************************************************ */
func New(price float64, quantity float64, volume float64) Order {
	order := Order { price, quantity, volume }
	return order
}

/* ************************************************************ */
/*							GETTERS								*/
/* ************************************************************ */
func (o Order) GetPrice() float64 {
	return o.price
}
func (o Order) GetQuantity() float64 {
	return o.quantity
}
func (o Order) GetVolume() float64 {
	return o.volume
}
