package util

import (
	"github.com/go-ashe/gulu/logx"
	"runtime/debug"
)

// GoSafe 启动一个安全的 goroutine，捕获 panic
func GoSafe(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logx.Error("[GoSafe] Recovered panic", "error", string(debug.Stack()))
			}
		}()
		fn()
	}()
}
