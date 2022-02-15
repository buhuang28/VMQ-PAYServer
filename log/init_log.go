package log

import (
	"Pay-Server/util"
	"log"
	"os"
	"time"
)

func init() {
	InitLog()
}

var (
	Logger *log.Logger
)

func InitLog() {
	exist := false
	exist = util.CheckFileExits(`log`)
	if !exist {
		_ = os.Mkdir("./log", os.ModePerm)
	}
	exist = false

	logFileNmae := `./log/` + time.Now().Format("20060102") + ".log"
	exist = util.CheckFileExits(logFileNmae)

	var f *os.File
	if !exist {
		f, _ = os.Create(logFileNmae)
	} else {
		//如果存在文件则 追加log
		f, _ = os.OpenFile(logFileNmae, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	}
	Logger = log.New(f, "", log.LstdFlags)
	Logger.SetFlags(log.LstdFlags | log.Lshortfile)
}
