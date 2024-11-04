package domain

type Symbol string

type SymbolList []Symbol

type Symbols struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}
