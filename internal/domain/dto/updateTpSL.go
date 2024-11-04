package dto

type UpdateTpSl struct {
	Id        int64  `json:"id"`
	Symbol    string `json:"symbol"`
	Exchange  string `json:"exchange"`
	TpPercent string `json:"tpPercent"`
	SlPercent string `json:"slPercent"`
	TpPrice   string `json:"tpPrice"`
	SlPrice   string `json:"slPrice"`
}
