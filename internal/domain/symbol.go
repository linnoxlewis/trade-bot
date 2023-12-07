package domain

type Symbol string

type SymbolList []Symbol

type Symbols struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}

func (sl SymbolList) IsEmpty() bool {
	return len(sl) == 0
}
