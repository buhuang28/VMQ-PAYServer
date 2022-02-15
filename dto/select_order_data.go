package dto

type SelectOrderData struct {
	OrderId    string `json:"order_id"`
	SelectTime int64  `json:"select_time"`
	Sign       string `json:"sign"`
}
