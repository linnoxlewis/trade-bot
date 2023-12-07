package domain

type BalanceSymbol struct {
	Symbol   string `json:"symbol"`
	Quantity string `json:"quantity"`
}

func NewBalanceSymbol(symbol, quantity string) BalanceSymbol {
	return BalanceSymbol{
		Symbol:   symbol,
		Quantity: quantity,
	}
}

type Balance []BalanceSymbol
