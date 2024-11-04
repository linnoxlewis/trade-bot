package domain

type Deal struct {
	BaseOrder *Order    //ордер созданный пользователем
	TpOrder   *Order    // тп -ордер, автоматически созданный из базового с тп
	SlOrder   *Order    // сл -ордер, автоматически созданный из базового с ck
	Settings  *Settings //структура с настройками тп сл
}

func (d *Deal) GetBaseOrder() *Order {
	return d.BaseOrder
}

func (d *Deal) GetTpOrder() *Order {
	return d.TpOrder
}

func (d *Deal) GetSlOrder() *Order {
	return d.SlOrder
}

func (d *Deal) GetSettings() *Settings {
	return d.Settings
}

func (d *Deal) SetBaseOrder(order *Order) {
	d.BaseOrder = order
}

func (d *Deal) SetTpOrder(tpOrder *Order) {
	d.TpOrder = tpOrder
}

func (d *Deal) SetSlOrder(slOrder *Order) {
	d.SlOrder = slOrder
}

func (d *Deal) SetSettings(settings *Settings) {
	d.Settings = settings
}
