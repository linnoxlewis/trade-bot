package dto

type TgOrder struct {
	Id        int64  `json:"id"`
	Command   string `json:"command"`
	Exchange  string `json:"exchange"`
	Symbol    string `json:"ccy"`
	OrderType string `json:"type"`
	Side      string `json:"side"`
	Quantity  string `json:"qty"`
	Price     string `json:"price"`
	TpPercent string `json:"tp_percent"`
	SlPercent string `json:"sl_percent"`
	TpPrice   string `json:"tp_price"`
	SlPrice   string `json:"sl_price"`
	TpType    string `json:"tp_type"`
	SlType    string `json:"sl_type"`
	Ts        string `json:"ts"`

	TimeInForce string `json:"tif"`
	StopPercent string `json:"stopPercent"`
	StopPrice   string `json:"stopPrice"`
	IcebergQty  string `json:"icebergQty"`
}
