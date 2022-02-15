package log

import (
	"bytes"
	"fmt"
	"runtime"
)

func PrintStackTrace(err interface{}) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v\n", err)
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}
	Logger.Println(buf.String())
}
