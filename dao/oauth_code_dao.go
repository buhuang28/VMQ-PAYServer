package dao

import (
	"Pay-Server/db"
	"Pay-Server/entity"
	"github.com/jinzhu/gorm"
	"log"
)

func AddOauthCode(link *gorm.DB, oauthCode *entity.OauthCode) error {
	create := link.Create(oauthCode)
	return create.Error
}

func SelectOauthCodeByInviteKey(inviteKey string) (bool, entity.OauthCode) {
	var user entity.OauthCode
	first := db.DbLink.Where("invite_key = ?", inviteKey).First(&user)
	//查得到数据
	if first == nil && user.ID != 0 {
		return true, user
	}
	return true, user
}

func SelectOneOauthCodeByCode(code string) (bool, entity.OauthCode) {
	var user entity.OauthCode
	first := db.DbLink.Where("code = ?", code).First(&user)
	if first.Error != nil {
		log.Logger.Println(first.Error, code)
		return false, user
	}
	if user.ID == 0 {
		return false, user
	}
	return true, user
}

func UpdateOneOauthCode(link *gorm.DB, oauthCode *entity.OauthCode) error {
	save := link.Model(oauthCode).Where("code = ?", oauthCode.Code).
		Updates(entity.OauthCode{OkTime: oauthCode.OkTime, EndTime: oauthCode.EndTime, LastIp: oauthCode.LastIp})
	return save.Error
}

func UpdateLastKeyByCode(link *gorm.DB, code, key string) error {
	save := link.Model(&entity.OauthCode{}).Where("code = ?", code).Update("last_key", key)
	return save.Error
}

func UpdateOauthCodeActivityInfoByCode(link *gorm.DB, oauthCode *entity.OauthCode) error {
	//save := link.Model(oauthCode).Where("code",oauthCode.Code).Updates("invitor","activity_time","end_time","ok_time")
	save := link.Model(oauthCode).Where("code = ?", oauthCode.Code).Updates(
		entity.OauthCode{ActivityTime: oauthCode.ActivityTime,
			EndTime: oauthCode.EndTime, OkTime: oauthCode.OkTime})
	return save.Error
}

//func UpdateOauthCodeActivityInfoByCode(link *gorm.DB,oauthCode *entity.OauthCode) error {
//	//save := link.entity(oauthCode).Where("code",oauthCode.Code).Updates("invitor","activity_time","end_time","ok_time")
//	save := link.Model(oauthCode).Where("code = ?",oauthCode.Code).Updates(
//		entity.OauthCode{Invitor: oauthCode.Invitor,ActivityTime: oauthCode.ActivityTime,
//			EndTime: oauthCode.EndTime,OkTime: oauthCode.OkTime})
//	return save.Error
//}

func SelectOneOauthCode() (bool, entity.OauthCode) {
	var user entity.OauthCode
	first := db.DbLink.Where("activity_time = ?", 0).First(&user)
	if first.Error != nil {
		log.Logger.Println(first.Error)
		return false, user
	}
	return true, user
}

func SelectOauthCodeList() []entity.OauthCode {
	var users []entity.OauthCode
	find := db.DbLink.Find(&users)
	if find.Error != nil {
		log.Logger.Println(find.Error)
		return nil
	}
	return users
}
