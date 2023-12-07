package domain

type Settings struct {
	Id        int64  `json:"id"`
	OrderId   int64  `json:"order_id"`
	TpPercent string `json:"tp_percent"`
	SlPercent string `json:"sl_percent"`
	TpPrice   string `json:"tp_price"`
	SlPrice   string `json:"sl_price"`
	TpType    string `json:"tp_type"`
	SlType    string `json:"sl_type"`
	Ts        string `json:"ts"`
	Date
}

func (s *Settings) IsTpEmpty() bool {
	return s.TpPrice == "" && s.TpPercent == ""
}

func (s *Settings) IsSlEmpty() bool {
	return s.SlPrice == "" && s.SlPercent == ""
}
