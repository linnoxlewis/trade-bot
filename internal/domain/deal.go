package domain

type Deal struct {
	BaseOrder *Order
	TpOrder   *Order
	SlOrder   *Order
	Settings  *Settings
}
