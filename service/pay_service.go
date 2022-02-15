package service

import (
	"Pay-Server/dao"
	"Pay-Server/db"
	"Pay-Server/dto"
	"Pay-Server/entity"
	"Pay-Server/log"
	"Pay-Server/util"
	"encoding/json"
	"fmt"
	"github.com/valyala/fastjson"
	"strconv"
	"strings"
)

var (
	SESSIONKEY     = "VMQ上的通信密钥，用于保证数据的安全性"
	IP             = "127.0.0.1"
	CREATEORDERAPI = "http://" + IP + ":8080/createOrder"
	CLOSEORDERAPI  = "http://" + IP + ":8080/closeOrder"
)

func CreateOrderService(code, randKey, sign string, price float64, createTime int64) (ok bool, createOrderResult dto.CreateOrderResult) {
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
			ok = false
		}
	}()

	//验证Sign
	getSign := util.GetClientSign(code, randKey, strconv.FormatInt(createTime, 10))

	if sign != getSign {
		return false, createOrderResult
	}

	//发起对VMQ Server端请求创建订单
	var createOrderData dto.CreateOrderData
	createOrderData.PayType = 2
	createOrderData.Price = price

	createOrderData.Param = code
	createOrderData.PayId = strconv.FormatInt(createTime, 10) + code
	createOrderData.Sign = util.GetSign(createOrderData.PayId, createOrderData.Param, strconv.FormatInt(createOrderData.PayType, 10), fmt.Sprint(createOrderData.Price), SESSIONKEY)
	priceString := fmt.Sprint(createOrderData.Price)
	if !strings.Contains(priceString, `.`) {
		priceString = priceString + ".0"
	}
	createOrderData.Sign = util.GetSign(createOrderData.PayId, createOrderData.Param, strconv.FormatInt(createOrderData.PayType, 10), priceString, SESSIONKEY)
	marshal, _ := json.Marshal(createOrderData)
	rMap := make(map[string]string)
	_ = json.Unmarshal(marshal, &rMap)
	rMap["type"] = strconv.FormatInt(createOrderData.PayType, 10)
	rMap["price"] = priceString
	sucess, result := util.GetRequest(CREATEORDERAPI, rMap)
	fmt.Println(string(result))
	if !sucess || result == nil {
		sucess, result = util.GetRequest(CREATEORDERAPI, rMap)
		if !sucess || result == nil {
			log.Logger.Println("HttpPost无响应")
			return false, createOrderResult
		}
	}
	err := json.Unmarshal(result, &createOrderResult)
	if err != nil {
		log.Logger.Println("创建订单json反序列化失败:", err)
		return false, createOrderResult
	}
	if createOrderResult.Code == 1 {
		go func() {
			var payOrder entity.PayOrder
			payOrder.Sign = createOrderData.Sign
			payOrder.Param = createOrderData.Param
			payOrder.PayId = createOrderData.PayId
			payOrder.Price = createOrderData.Price
			payOrder.PayType = createOrderData.PayType
			payOrder.OrderId = createOrderResult.Data.OrderID
			payOrder.ReallyPrice = createOrderResult.Data.ReallyPrice
			er := dao.CreateOrder(db.DbLink, &payOrder)
			if er != nil {
				_ = dao.CreateOrder(db.DbLink, &payOrder)
			}
		}()
	}
	//持久化订单
	return sucess, createOrderResult
}

func SelectOrderService(orderId, sign string, selectTime int64) (sucess bool, oauthCode entity.OauthCode) {
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
			sucess = false
		}
	}()

	getSign := util.GetSign(orderId, strconv.FormatInt(selectTime, 10))
	if getSign != sign {
		return false, entity.OauthCode{}
	}

	ok, order := dao.SelectOrderByOrderId(db.DbLink, orderId)
	selectSucess, oauthCode := dao.SelectOneOauthCodeByCode(order.Param)
	//查得到订单并且成功支付
	if ok && order.PayStatus == 1 && selectSucess {
		return true, oauthCode
	}
	return false, entity.OauthCode{}
}

func PayCallBackService(callBackData dto.CallBackData) bool {
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
	}()

	priceString := fmt.Sprint(callBackData.Price)
	ReallyPriceString := fmt.Sprint(callBackData.ReallyPrice)
	if !strings.Contains(priceString, `.`) {
		priceString = priceString + ".0"
	}

	if !strings.Contains(ReallyPriceString, `.`) {
		ReallyPriceString = ReallyPriceString + ".0"
	}

	getSign := util.GetSign(callBackData.PayId, callBackData.Param, strconv.FormatInt(callBackData.Type, 10),
		priceString, ReallyPriceString, SESSIONKEY)
	if getSign != callBackData.Sign {
		return false
	}

	var payOrder entity.PayOrder
	payOrder.PayStatus = 1
	payOrder.PayId = callBackData.PayId
	payOrder.Param = callBackData.Param

	status := dao.SelectOrderStatus(db.DbLink, &payOrder)
	if status == 1 {
		return true
	}
	if status != 0 {
		log.Logger.Println("查询不到订单:", payOrder.PayId, "授权码:", payOrder.Param)
		return false
	}

	tx := db.DbLink.Begin()
	err := dao.UpdateOrderStatusByOrderId(tx, &payOrder)
	if err != nil {
		log.Logger.Println(err)
		tx.Rollback()
	}

	sucess := ProLongService(tx, callBackData.Param)
	if !sucess {
		tx.Rollback()
	}
	tx.Commit()
	if sucess && err == nil {
		return true
	}
	return false
}

func DeleteOrderService(orderId, sign string, selectTime int64, price float64) bool {
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
	}()

	getSign := util.GetSign(orderId, strconv.FormatInt(selectTime, 10), fmt.Sprint(price))
	if getSign != sign {
		return false
	}

	data := make(map[string]string)
	data["orderId"] = orderId
	data["sign"] = util.GetSign(orderId, SESSIONKEY)

	ok, bytes := util.GetRequest(CLOSEORDERAPI, data)

	if !ok || bytes == nil {
		return false
	}

	values, err := fastjson.ParseBytes(bytes)
	if err != nil {
		return false
	}

	code := values.GetInt64("code")
	return code == 1
}
