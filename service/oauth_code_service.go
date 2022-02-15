package service

import (
	"Pay-Server/dao"
	"Pay-Server/db"
	"Pay-Server/dto"
	"Pay-Server/entity"
	"Pay-Server/log"
	"Pay-Server/util"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	CODELENGTH   = 25
	INVITELENGTH = 6
	KEYLENGTH    = 25

	DAYTIME   int64 = 1 * 60 * 60 * 24
	WEEKTIME  int64 = DAYTIME * 7
	MONTHTIME int64 = DAYTIME * 31
	YEARTIME        = 365 * DAYTIME

	proLongLock  sync.Mutex
	inviteLock   sync.Mutex
	activityLock sync.Mutex

	oauthKeyMap  = make(map[string]string)
	oauthKeyLock sync.Mutex
)

func AddOauthCodeService(Count int64) (errorCount int64) {
	errorCount = 0
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
			errorCount = -1
		}
	}()

	var i int64
	for i = 0; i < Count; i++ {
		var oauthCode entity.OauthCode
		code := util.RandString(CODELENGTH)
		is := strconv.FormatInt(time.Now().Unix(), 10)
		code = util.Md5(is + code)
		code = strings.ToUpper(code)
		oauthCode.Code = code

		//inviteKey := ""
		//for {
		//	inviteKey = util.RandString(INVITELENGTH)
		//	ok,_ := dao.SelectOauthCodeByInviteKey(inviteKey)
		//	if ok {
		//		break
		//	}
		//}
		//oauthCode.InviteKey = inviteKey

		err := dao.AddOauthCode(db.DbLink, &oauthCode)
		if err != nil {
			errorCount++
			log.Logger.Println("create_error:", err)
			continue
		}

	}
	return errorCount
}

//充钱服务
func ProLongService(tx *gorm.DB, code string) (sucess bool) {
	//fmt.Println(code,"开始充值")
	//tx := entity.DbLink.Begin()
	proLongLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
			tx.Rollback()
			sucess = false
		}
		proLongLock.Unlock()
	}()
	//查询该授权码对应的实体
	sucess, oauthCode := dao.SelectOneOauthCodeByCode(code)
	if !sucess || oauthCode.ID < 1 {
		return sucess
	}
	//加时
	nowTime := time.Now().Unix()
	if oauthCode.EndTime < nowTime {
		oauthCode.EndTime = nowTime
	}

	oauthCode.EndTime = oauthCode.EndTime + YEARTIME
	if nowTime < 1630434461 {
		oauthCode.EndTime = oauthCode.EndTime + MONTHTIME + MONTHTIME
	}
	//查出授权码邀请人
	oauthCode.OkTime = YEARTIME
	err := dao.UpdateOneOauthCode(tx, &oauthCode)
	if err != nil {
		tx.Rollback()
		//var failCode entity.FailCode
		//failCode.Code = code
		//failCode.FailTime = nowTime
		//failCode.ProTime = YEARTIME
		//err = dao.AddFailCode(entity.DbLink, &failCode)
		//if err != nil {
		//	log.Logger.Println("授权码用户:",code,"在",nowTime,"延时一年失败")
		//}
		return sucess
	}
	tx.Commit()
	sucess = true

	//go func() {
	//	inviteLock.Lock()
	//	defer func() {
	//		er := recover()
	//		if er != nil {
	//			log.PrintStackTrace(er)
	//		}
	//		inviteLock.Unlock()
	//	}()
	//
	//	//if oauthCode.Invitor == "" {
	//	//	return
	//	//}
	//	//ok, invitor := dao.SelectOauthCodeByInviteKey(oauthCode.Invitor)
	//	//if ok && invitor.Code != "" {
	//	//	if invitor.EndTime < time.Now().Unix() {
	//	//		invitor.EndTime = time.Now().Unix()
	//	//	}
	//	//	invitor.EndTime = invitor.EndTime + MONTHTIME * 3
	//	//	err = dao.UpdateOneOauthCode(entity.DbLink, &invitor)
	//	//	if err != nil {
	//	//		var invitorFail entity.FailCode
	//	//		invitorFail.ProTime = MONTHTIME
	//	//		invitorFail.Code = invitor.Code
	//	//		invitorFail.FailTime = time.Now().Unix()
	//	//		err = dao.AddFailCode(entity.DbLink, &invitorFail)
	//	//		if err != nil {
	//	//			log.Logger.Println("授权码用户:",code,"在",nowTime,"延时一个月失败")
	//	//		}
	//	//	}
	//	//}
	//}()
	return sucess
}

func ActivityCodeService(activityData dto.ActivityData) (bool, entity.OauthCode) {
	activityLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
		activityLock.Unlock()
	}()

	ok, oauthCode := dao.SelectOneOauthCode()
	if !ok {
		time.Sleep(time.Second)
		ok, oauthCode = dao.SelectOneOauthCode()
		if !ok {
			return false, oauthCode
		}
	}

	oauthCode.ActivityTime = activityData.ActivityTime
	//oauthCode.Invitor = activityData.Invitor
	if activityData.Invitor != "" {
		oauthCode.EndTime = activityData.ActivityTime + WEEKTIME
		oauthCode.OkTime = WEEKTIME
	} else {
		oauthCode.EndTime = activityData.ActivityTime + 3*DAYTIME
		oauthCode.OkTime = 3 * DAYTIME
	}
	err := dao.UpdateOauthCodeActivityInfoByCode(db.DbLink, &oauthCode)
	return err == nil, oauthCode
}

func AuthOauthCodeService(code, key string) (sucess bool) {
	sucess = true
	oauthKeyLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
			sucess = false
		}
		oauthKeyLock.Unlock()
	}()
	byCode, oauthCode := dao.SelectOneOauthCodeByCode(code)
	if !byCode || oauthCode.EndTime < time.Now().Unix() {
		return false
	}
	oauthKeyMap[code] = key
	go func() {
		err := dao.UpdateLastKeyByCode(db.DbLink, code, key)
		if err != nil {
			log.Logger.Println(err)
		}
	}()
	return sucess
}

func CheckOauthCodeService(code, key string) (sucess bool) {
	oauthKeyLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
		oauthKeyLock.Unlock()
	}()

	if oauthKeyMap[code] != "" && oauthKeyMap[code] == key {
		return true
	}

	if oauthKeyMap[code] == "" {
		success, oauthCode := dao.SelectOneOauthCodeByCode(code)
		if !success {
			return false
		}
		if time.Now().Unix() > oauthCode.EndTime {
			return false
		}
		oauthKeyMap[code] = key
		return true
	}
	return false
}

func RecoverTimeService(code, sign string, t int64) (bool, entity.OauthCode) {
	getSign := util.GetSign(code, strconv.FormatInt(t, 10))
	if getSign != sign {
		return false, entity.OauthCode{}
	}
	sucess, oauthCode := dao.SelectOneOauthCodeByCode(code)
	return sucess, oauthCode
}
