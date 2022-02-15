package dao

import (
	"Pay-Server/entity"
	"github.com/jinzhu/gorm"
)

func CreateOrder(link *gorm.DB, order *entity.PayOrder) error {
	return link.Create(order).Error
}

func UpdateOrderStatusByOrderId(link *gorm.DB, order *entity.PayOrder) error {
	update := link.Model(order).Where("param = ? and pay_id = ?", order.Param, order.PayId).Update("pay_status", 1)
	return update.Error
}

func SelectOrderStatus(link *gorm.DB, order *entity.PayOrder) int64 {
	var newOrder entity.PayOrder
	first := link.Where("param = ? and pay_id = ?", order.Param, order.PayId).First(&newOrder)
	if first.Error != nil {
		return -1
	}
	return newOrder.PayStatus
}

func SelectOrderByOrderId(link *gorm.DB, orderId string) (bool, entity.PayOrder) {
	var order entity.PayOrder
	first := link.Where("order_id = ?", orderId).First(&order)
	if first.Error != nil || order.ID < 1 {
		return false, order
	}
	return true, order
}

//查询未支付的订单并且更新
func ProcessTimeOutOrder(link *gorm.DB) []float64 {
	var orders []entity.PayOrder
	link.Raw("select id,price from pay_order where UNIX_TIMESTAMP() - UNIX_TIMESTAMP(create_time) > ? and  pay_status = 0", 300).Scan(&orders)
	var ids []uint
	var prices []float64
	for _, v := range orders {
		ids = append(ids, v.ID)
		prices = append(prices, v.Price)
	}
	if len(ids) > 0 {
		link.Exec("update pay_order set pay_status = -1 where id in (?)", ids)
	}
	return prices
}
