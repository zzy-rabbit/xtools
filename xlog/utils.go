package xlog

import (
	"context"
	"fmt"
	"github.com/zzy-rabbit/xtools/xexecutable"
	"path/filepath"
)

func setDefault(ctx context.Context, config *Config) {
	if config.Level == 0 {
		config.Level = LevelDebug
	}
	if config.MaxSize == 0 {
		config.MaxSize = 10
	}
	if config.Name == "" {
		config.Name = "default"
	}
	if config.Path == "" {
		config.Path = filepath.Join(xexecutable.GetProcessAbsPath(), "log")
	}
	if config.Suffix == "" {
		config.Suffix = ".log"
	}
	fmt.Printf("logger config %+v\n", config)
}
