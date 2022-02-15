package dto

type Result struct {
	Status  int64       `json:"status"`  //状态码  -1失败   0成功
	Message string      `json:"message"` //通知消息
	Data    interface{} `json:"data"`
}

func (r *Result) Fail(code int64, msg string) Result {
	if code == 0 {
		code = -1
	}
	return Result{code, msg, nil}
}

func (r *Result) Sucess(msg string) Result {
	return Result{Status: 1, Message: msg}
}

func (r *Result) Sucess2(msg string, data interface{}) Result {
	return Result{Status: 1, Message: msg, Data: data}
}
