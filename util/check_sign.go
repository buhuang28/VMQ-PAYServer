package util

import (
	"crypto/md5"
	"encoding/hex"
)

var (
	SESSIONKEY = "payserver-test"
)

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetSign(params ...string) string {
	content := ""
	for _, v := range params {
		content += v
	}
	return Md5(content)
}

//结果与Pay-Client的Sign一致
func GetClientSign(params ...string) string {
	content := ""
	for _, v := range params {
		content += v
	}
	return Md5(SESSIONKEY + content + SESSIONKEY)
}
