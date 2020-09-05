package util

import (
	"chat-test/conf"
	"chat-test/log"
	"runtime"
)

// PrintPanicStack recover并打印堆栈
func PrintPanicStack() {
	if r := recover(); r != nil {
		buf := make([]byte, conf.LenStackBuf)
		l := runtime.Stack(buf, false)
		log.Error("%v: %s", r, string(buf[:l]))
	}
}
