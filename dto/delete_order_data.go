package dto

type DeleteOrderData struct {
	OrderId    string  `json:"order_id"`
	DeleteTime int64   `json:"delete_time"`
	Sign       string  `json:"sign"`
	Price      float64 `json:"price"`
}
