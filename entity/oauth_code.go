package entity

//对应数据库的
type OauthCode struct {
	ID           uint `gorm:"primary_key"`
	Code         string
	ActivityTime int64
	EndTime      int64
	OkTime       int64
	LastIp       string
	LastKey      string
}
