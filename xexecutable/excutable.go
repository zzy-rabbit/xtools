package xexecutable

import (
	"os"
)

type service struct {
	processPath string
}

var instance = &service{}

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	instance.processPath = path
}

// GetProcessAbsPath 获取当前执行文件绝对路径
func GetProcessAbsPath() string {
	return instance.processPath
}
