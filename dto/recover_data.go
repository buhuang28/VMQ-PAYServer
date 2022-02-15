package dto

type RecoverData struct {
	RTime int64  `json:"r_time"`
	Code  string `json:"code"`
	Sign  string `json:"sign"`
}
