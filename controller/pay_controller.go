package controller

import (
	"Pay-Server/dto"
	"Pay-Server/log"
	"Pay-Server/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
	"sync"
)

type PayController struct{}

func (o *PayController) NewPayController() PayController {
	return PayController{}
}

func (o *PayController) CreateOrder(c *gin.Context) {
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
	}()

	var result dto.CreateOrderResult
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Logger.Println("json参数错误:", err.Error())
		c.JSON(200, result.Fail(-1, "json error"))
		return
	}

	var data dto.CreateOrderBody
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Logger.Println("json反序列化错误:", err.Error())
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	if data.Code == "" {
		c.JSON(200, result.Fail(-2, "sucess"))
		return
	}

	if data.CreateTime == 0 {
		c.JSON(200, result.Fail(-3, "sucess"))
		return
	}

	if data.RandKey == "" {
		c.JSON(200, result.Fail(-5, "sucess"))
		return
	}

	if data.Sign == "" {
		c.JSON(200, result.Fail(-6, "sucess"))
		return
	}

	price := service.GetPayPrice()

	sucess, createResult := service.CreateOrderService(data.Code, data.RandKey, data.Sign, price, data.CreateTime)
	if !sucess {
		c.JSON(200, result.Fail(-7, "sucess"))
		return
	}

	c.JSON(200, createResult)
	return
}

func (o *PayController) SelectOrder(c *gin.Context) {
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
	}()
	var result dto.Result
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Logger.Println("json参数错误:", err.Error())
		c.JSON(200, result.Fail(-1, "json error"))
		return
	}

	var selectOrderData dto.SelectOrderData
	err = json.Unmarshal(body, &selectOrderData)
	if err != nil {
		log.Logger.Println("json反序列化错误:", err.Error())
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	if selectOrderData.OrderId == "" {
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	if selectOrderData.SelectTime == 0 {
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	if selectOrderData.Sign == "" {
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	paySucess, oauthCode := service.SelectOrderService(selectOrderData.OrderId, selectOrderData.Sign, selectOrderData.SelectTime)

	if !paySucess {
		c.JSON(200, result.Fail(-60, "fail"))
		return
	}
	c.JSON(200, result.Sucess(strconv.FormatInt(oauthCode.EndTime, 10)))
}

var (
	CallBackLock sync.Mutex
)

func (o *PayController) PayCallBack(c *gin.Context) {
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
		CallBackLock.Unlock()
	}()
	CallBackLock.Lock()
	var callBackData dto.CallBackData
	err := c.ShouldBind(&callBackData)
	if err != nil {
		log.Logger.Println("收款回调失败:", err)
		return
	}

	if callBackData.Sign == "" {
		log.Logger.Println("收款回调失败2:", err)
		return
	}
	_ = service.PayCallBackService(callBackData)
	service.RecoverPrice(callBackData.Price)
	c.String(200, "success")
}

//这个方法是为了给VMQ返回一个返回值，不然VMQ一直是失败的状态
func (o *PayController) WaitCall(c *gin.Context) {
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
	}()
	var callBackData dto.CallBackData
	err := c.ShouldBind(&callBackData)
	if err != nil {
		log.Logger.Println(err)
	}
	c.String(200, "success")
	return
}

func (o *PayController) DeleteOrder(c *gin.Context) {
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
	}()

	//var result dto.Result
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Logger.Println("json参数错误:", err.Error())
		//c.JSON(200,result.Fail(-1,"json error"))
		return
	}

	var deleteOrder dto.DeleteOrderData
	err = json.Unmarshal(body, &deleteOrder)
	if err != nil {
		log.Logger.Println("json反序类化失败:", err)
		return
	}

	if deleteOrder.OrderId == "" {
		return
	}

	if deleteOrder.DeleteTime == 0 {
		return
	}

	if deleteOrder.Sign == "" {
		log.Logger.Println("来自", c.ClientIP(), "无sign验证的DeleteOrder")
		return
	}

	orderService := service.DeleteOrderService(deleteOrder.OrderId, deleteOrder.Sign, deleteOrder.DeleteTime, deleteOrder.Price)
	if orderService {
		service.RecoverPrice(deleteOrder.Price)
	}
	c.String(200, "success")
}
