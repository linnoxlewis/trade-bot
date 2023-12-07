package domain

type Depth struct {
	LastUpdateID int64        `json:"lastUpdateId"`
	Bids         []PriceLevel `json:"bids"`
	Asks         []PriceLevel `json:"asks"`
}

type PriceLevel struct {
	Price    string
	Quantity string
}
