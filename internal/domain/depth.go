package domain

type Depth struct {
	LastUpdateID int64
	Bids         []PriceLevel
	Asks         []PriceLevel
}

type PriceLevel struct {
	Price    string
	Quantity string
}
