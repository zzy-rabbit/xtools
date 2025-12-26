package xtrace

import (
	"context"
	"fmt"
	"github.com/zzy-rabbit/xtools/xlog"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func Trace(ctx context.Context) func(args ...any) {
	log := xlog.GetDefaultLogger(ctx)
	startTime := time.Now()
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	funcName := ""
	if fn != nil {
		temps := strings.Split(fn.Name(), ".")
		for i := len(temps) - 1; i >= 0; i-- {
			if funcName != "" {
				funcName = temps[i] + "." + funcName
			} else {
				funcName = temps[i]
			}
			matched, _ := regexp.MatchString(`^\d+$|^func\d+$`, temps[i])
			if matched {
				continue
			}
			break
		}
	}

	log.Info(ctx, "%s start", funcName)
	return func(args ...interface{}) {
		paramStr := ""
		for i, v := range args {
			// 获取字段值
			fieldName := byte('a' + i)
			fieldValue := fmt.Sprintf("%v", v)
			paramStr += fmt.Sprintf(", %c:%v", fieldName, fieldValue)
		}
		cost := time.Since(startTime)
		log.Info(ctx, "%s end cost:%v, %s\n", funcName, cost, paramStr)
	}
}
