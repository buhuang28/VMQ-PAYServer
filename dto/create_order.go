package dto

//接收创建订单的数据
type CreateOrderBody struct {
	Code       string `json:"code"`
	CreateTime int64  `json:"create_time"`
	Sign       string `json:"sign"`
	RandKey    string `json:"rand_key"`
}

//对VMQ Server端请求创建订单
type CreateOrderData struct {
	PayId   string  `json:"payId"`
	PayType int64   `json:"type"`
	Price   float64 `json:"price"`
	Sign    string  `json:"sign"`
	Param   string  `json:"param"`
}

//创建订单返回的数据
type CreateOrderResult struct {
	Code int64     `json:"code"`
	Msg  string    `json:"msg"`
	Data OrderData `json:"data"`
}

type OrderData struct {
	PayID       string  `json:"payId"`
	OrderID     string  `json:"orderId"`
	PayType     int64   `json:"payType"`
	Price       float64 `json:"price"`
	ReallyPrice float64 `json:"reallyPrice"`
	PayURL      string  `json:"payUrl"`
	IsAuto      int64   `json:"isAuto"`
	State       int64   `json:"state"`
	TimeOut     int64   `json:"timeOut"`
	Date        int64   `json:"date"`
}

func (c *CreateOrderResult) Fail(code int64, msg string) CreateOrderResult {
	return CreateOrderResult{Code: code, Msg: msg}
}
