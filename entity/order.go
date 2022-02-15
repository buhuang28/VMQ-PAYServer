package entity

//创建订单的数据
type PayOrder struct {
	ID          uint `gorm:"primary_key"`
	PayId       string
	OrderId     string
	Param       string
	Price       float64
	ReallyPrice float64
	PayType     int64
	PayStatus   int64
	Sign        string
}
