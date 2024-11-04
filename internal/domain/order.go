package domain

import "sync"

type Order struct {
	Id          int64
	ExecOrderId int64
	UserId      int64
	Symbol      string
	Side        string
	OrderType   string
	Quantity    string
	Price       string
	TimeInForce string
	StopPrice   string
	IcebergQty  string
	Exchange    string
	Status      string
	TpSl        string
	Inwork      bool
	sync.RWMutex
	Settings
}

type OrderList struct {
	Orders []*Order
	sync.Mutex
}

func NewOrderList() *OrderList {
	return &OrderList{
		Orders: make([]*Order, 0),
	}
}

func (o *OrderList) Add(order *Order) {
	o.Lock()
	defer o.Unlock()
	o.Orders = append(o.Orders, order)
}
