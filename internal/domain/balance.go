package domain

type BalanceSymbol struct {
	Symbol   string
	Quantity string
}

func NewBalanceSymbol(symbol, quantity string) BalanceSymbol {
	return BalanceSymbol{
		Symbol:   symbol,
		Quantity: quantity,
	}
}

type Balances []BalanceSymbol
