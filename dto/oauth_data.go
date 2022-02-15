package dto

//对应接口接收的数据结构体
type OauthBody struct {
	OauthCode string `json:"oauth_code"` //授权码
	OauthTime int64  `json:"oauth_time"` //时间戳
	OauthSign string `json:"oauth_sign"` //sign签名
	OauthKey  string `json:"oauth_key"`  //随机字符串
}
