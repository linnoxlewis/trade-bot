package domain

import "sync"

type Order struct {
	Id          int64  `json:"id"`
	ExecOrderId int64  `json:"execOrderId"`
	UserId      int64  `json:"user_id"`
	Symbol      string `json:"symbol"`
	Side        string `json:"side"`
	OrderType   string `json:"orderType"`
	Quantity    string `json:"quantity"`
	Price       string `json:"price"`
	TimeInForce string `json:"timeInForce"`
	StopPrice   string `json:"stopPrice"`
	IcebergQty  string `json:"icebergQty"`
	Exchange    string `json:"exchange"`
	Status      string `json:"status"`
	TpSl        string `json:"tpSl"`
	Inwork      bool   `json:"-"`
	sync.RWMutex
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
