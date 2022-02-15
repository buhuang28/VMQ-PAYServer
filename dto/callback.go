package dto

type CallBackData struct {
	PayId       string  `form:"payId"`
	Param       string  `form:"param"`
	Type        int64   `form:"type"`
	Price       float64 `form:"price"`
	ReallyPrice float64 `form:"reallyPrice"`
	Sign        string  `form:"sign"`
}
