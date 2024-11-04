package dto

type UpdateOrder struct {
	OrderId     int64   `json:"orderId"`
	Exchange    string  `json:"exchange"`
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	OrderType   string  `json:"orderType"`
	Quantity    string  `json:"quantity"`
	Price       string  `json:"price"`
	TimeInForce string  `json:"timeInForce"`
	StopPrice   string  `json:"stopPrice"`
	StopPercent string  `json:"stopPercent"`
	Tp          float32 `json:"tp"`
	Sl          float32 `json:"sl"`
	Ts          float32 `json:"ts"`
}

func NewUpdateOrder(
	orderId int64,
	symbol,
	side,
	orderType,
	quantity,
	price,
	timeInForce,
	stopPrice,
	stopPercent string,
	tp, sl, ts float32,
) *UpdateOrder {
	return &UpdateOrder{
		OrderId:     orderId,
		Symbol:      symbol,
		Side:        side,
		OrderType:   orderType,
		Quantity:    quantity,
		Price:       price,
		TimeInForce: timeInForce,
		StopPrice:   stopPrice,
		StopPercent: stopPercent,
		Tp:          tp,
		Sl:          sl,
		Ts:          ts,
	}
}
