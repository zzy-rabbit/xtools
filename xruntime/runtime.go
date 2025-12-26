package xruntime

import (
	"context"
	"regexp"
	"runtime"
	"strings"
)

func GetShortFuncName(ctx context.Context, skip int) (string, int) {
	pc, _, line, _ := runtime.Caller(skip + 1)
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
	return funcName, line
}
