package controller

import (
	"Pay-Server/dto"
	"Pay-Server/log"
	"Pay-Server/service"
	"Pay-Server/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
	"sync"
	"time"
)

var (
	errorSignIpMap = make(map[string]int64)
	signLock       sync.Mutex
)

type OauthController struct{}

func (o *OauthController) NewOauthController() OauthController {
	return OauthController{}
}

//新增授权码
func (o *OauthController) AddOauthCode(c *gin.Context) {
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
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}
	var data dto.AddCodeBody
	err = json.Unmarshal(body, &data)
	if data.Count == 0 {
		c.JSON(200, result.Fail(-2, "count error"))
		return
	}
	errCount := service.AddOauthCodeService(data.Count)
	content := ""
	if errCount == 0 {
		content = "sucess"
	} else {
		content = strconv.FormatInt(errCount, 10)
	}
	c.JSON(200, result.Sucess(content))
}

//激活接口  -- 返回授权码，修改激活时间
func (o *OauthController) ActiviteOuathCode(c *gin.Context) {
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
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}
	var activityData dto.ActivityData
	err = json.Unmarshal(body, &activityData)
	if err != nil {
		log.Logger.Println("json参数错误:", err.Error())
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}

	if activityData.ActivityTime == 0 {
		log.Logger.Println("来自IP:", c.ClientIP(), "无ActivityTime的访问")
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}

	if activityData.TryQQ == "" {
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}

	if activityData.Sign == "" {
		log.Logger.Println("来自IP:", c.ClientIP(), "无sign的访问")
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}

	sign := util.GetSign(strconv.FormatInt(activityData.ActivityTime, 10), activityData.Invitor, activityData.TryQQ)
	if sign != activityData.Sign {
		log.Logger.Println("来自IP:", c.ClientIP(), "错误sign的访问")
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}

	sucess, oauthCode := service.ActivityCodeService(activityData)
	if sucess {
		c.JSON(200, result.Sucess(oauthCode.Code+"|"+strconv.FormatInt(oauthCode.EndTime, 10)))
	} else {
		c.JSON(200, result.Fail(-100, "sucess"))
	}
	return
}

//定时鉴权接口
func (o *OauthController) CheckOauthCode(c *gin.Context) {
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
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}

	var oauthData dto.OauthBody
	err = json.Unmarshal(body, &oauthData)
	if err != nil {
		log.Logger.Println("json反序列化失败:", err.Error())
		c.JSON(200, result.Fail(-1, "json错误"))
		return
	}

	if oauthData.OauthCode == "" {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthCode的访问")
		c.JSON(200, "sucess")
		return
	}

	if oauthData.OauthTime == 0 {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthTime的访问")
		c.JSON(200, "sucess")
		return
	}

	if oauthData.OauthKey == "" {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthKey的访问")
		c.JSON(200, "sucess")
		return
	}

	if oauthData.OauthSign == "" {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthSign的访问")
		c.JSON(200, "sucess")
		return
	}

	sign := util.GetSign(strconv.FormatInt(oauthData.OauthTime, 10), oauthData.OauthKey, oauthData.OauthCode)
	if sign != oauthData.OauthSign {
		log.Logger.Println("来自IP:", c.ClientIP(), "sign错误")
		signLock.Lock()
		errorSignIpMap[c.ClientIP()] = time.Now().Unix()
		signLock.Unlock()
		c.JSON(200, "sucess")
		return
	}

	sucess := service.CheckOauthCodeService(oauthData.OauthCode, oauthData.OauthKey)
	if sucess {
		c.JSON(200, result.Sucess("sucess"))
	} else {
		c.JSON(200, result.Fail(-105, "sucess"))
	}
	return
}

//更加安全的鉴权接口
func (o *OauthController) CheckOauthCode2(c *gin.Context) {
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
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}

	var oauthData dto.OauthBody
	err = json.Unmarshal(body, &oauthData)
	if err != nil {
		log.Logger.Println("json反序列化失败:", err.Error())
		c.JSON(200, result.Fail(-1, "json错误"))
		return
	}

	if oauthData.OauthCode == "" {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthCode的访问")
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	if oauthData.OauthTime == 0 {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthTime的访问")
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	if oauthData.OauthKey == "" {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthKey的访问")
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	if oauthData.OauthSign == "" {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthSign的访问")
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	sign := util.GetSign(strconv.FormatInt(oauthData.OauthTime, 10), oauthData.OauthKey, oauthData.OauthCode)
	if sign != oauthData.OauthSign {
		log.Logger.Println("来自IP:", c.ClientIP(), "sign错误")
		signLock.Lock()
		errorSignIpMap[c.ClientIP()] = time.Now().Unix()
		signLock.Unlock()
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	sucess := service.CheckOauthCodeService(oauthData.OauthCode, oauthData.OauthKey)
	sucessSign := util.GetSign(oauthData.OauthCode + "1" + strconv.FormatInt(oauthData.OauthTime, 10))
	if sucess {
		c.JSON(200, result.Sucess(sucessSign))
	} else {
		c.JSON(200, result.Fail(-1, sucessSign))
	}
	return
}

//鉴权认证接口  -- 添加key给服务器
func (o *OauthController) AuthOauthCode(c *gin.Context) {
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
	}()

	eIpTs := errorSignIpMap[c.ClientIP()]
	if eIpTs > 0 {
		if time.Now().Unix()-eIpTs < 120 {
			return
		}
	}

	var result dto.Result
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Logger.Println("json参数错误:", err.Error())
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}

	var oauthData dto.OauthBody
	err = json.Unmarshal(body, &oauthData)
	if err != nil {
		log.Logger.Println("json反序列化失败:", err.Error())
		c.JSON(200, result.Fail(-1, "json错误"))
		return
	}

	if oauthData.OauthCode == "" {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthCode的访问")
		c.JSON(200, "sucess")
		return
	}

	if oauthData.OauthTime == 0 {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthTime的访问")
		c.JSON(200, "sucess")
		return
	}

	if oauthData.OauthKey == "" {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthKey的访问")
		c.JSON(200, "sucess")
		return
	}

	if oauthData.OauthSign == "" {
		log.Logger.Println("来自IP:", c.ClientIP(), "无OauthSign的访问")
		c.JSON(200, "sucess")
		return
	}

	sign := util.GetSign(strconv.FormatInt(oauthData.OauthTime, 10), oauthData.OauthKey, oauthData.OauthCode)
	if sign != oauthData.OauthSign {
		log.Logger.Println("来自IP:", c.ClientIP(), "sign错误")
		signLock.Lock()
		errorSignIpMap[c.ClientIP()] = time.Now().Unix()
		signLock.Unlock()
		c.JSON(200, "sucess")
		return
	}

	sucess := service.AuthOauthCodeService(oauthData.OauthCode, oauthData.OauthKey)
	if sucess {
		c.JSON(200, result.Sucess("sucess"))
	} else {
		c.JSON(200, result.Fail(-1, "sucess"))
	}
	return
}

//充钱接口
//func (o *OauthController) ProLongOuathCodeTime(c *gin.Context) {
//	defer func() {
//		e := recover()
//		if e != nil {
//			log.Logger.Println(e)
//		}
//	}()
//	eIpTs := errorSignIpMap[c.ClientIP()]
//	if eIpTs > 0 {
//		if time.Now().Unix() - eIpTs < 60 {
//			return
//		}
//	}
//
//	var result dto.Result
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Logger.Println("json参数错误:",err.Error())
//		c.JSON(200,result.Fail(-1,"参数错误"))
//		return
//	}
//
//	var proLongData entity.ProLongData
//	err = json.Unmarshal(body, &proLongData)
//	if err != nil {
//		log.Logger.Println("json反序列化失败:",err.Error())
//		c.JSON(200,result.Fail(-1,"json错误"))
//		return
//	}
//
//	if proLongData.OauthCode == "" {
//		log.Logger.Println("来自IP:",c.ClientIP(),"无OauthCode的访问")
//		c.JSON(200,"sucess")
//		return
//	}
//
//	if proLongData.OauthAddTime == 0 {
//		log.Logger.Println("来自IP:",c.ClientIP(),"无OauthAddTime的访问")
//		c.JSON(200,"sucess")
//		return
//	}
//
//	if proLongData.OauthTime == 0 {
//		log.Logger.Println("来自IP:",c.ClientIP(),"无OauthTime的访问")
//		c.JSON(200,"sucess")
//		return
//	}
//
//
//	if proLongData.OauthSign == "" {
//		log.Logger.Println("来自IP:",c.ClientIP(),"无sign的访问")
//		signLock.Lock()
//		errorSignIpMap[c.ClientIP()] = time.Now().Unix()
//		signLock.Unlock()
//		c.JSON(200,"sucess")
//		return
//	}
//
//	sign := util.GetSign(strconv.FormatInt(proLongData.OauthTime,10), strconv.FormatInt(proLongData.OauthAddTime,10), proLongData.OauthCode)
//	if sign != proLongData.OauthSign {
//		log.Logger.Println("来自IP:",c.ClientIP(),"sign错误")
//		signLock.Lock()
//		errorSignIpMap[c.ClientIP()] = time.Now().Unix()
//		signLock.Unlock()
//		c.JSON(200,"sucess")
//		return
//	}
//	sucess := service.ProLongService(model.DbLink.Begin(),proLongData.OauthCode)
//	if sucess {
//		c.JSON(200,result.Sucess("sucess"))
//	}else {
//		c.JSON(200,result.Fail(-10,"fail"))
//	}
//	return
//}

//通过授权码查询
func (o *OauthController) RecoverTime(c *gin.Context) {
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
		c.JSON(200, result.Fail(-1, "参数错误"))
		return
	}
	var data dto.RecoverData
	err = json.Unmarshal(body, &data)
	if err != nil {
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	if data.Code == "" {
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	if data.RTime == 0 {
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	if data.Sign == "" {
		c.JSON(200, result.Fail(-1, "sucess"))
		return
	}

	sucess, oauthCode := service.RecoverTimeService(data.Code, data.Sign, data.RTime)
	if sucess {
		c.JSON(200, result.Sucess(data.Code+"|"+strconv.FormatInt(oauthCode.EndTime, 10)))
	} else {
		c.JSON(200, result.Fail(-1, "sucess"))
	}
	return
}
