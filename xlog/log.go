package xlog

import (
	"context"
	"fmt"
	"github.com/zzy-rabbit/xtools/xcontext"
	"github.com/zzy-rabbit/xtools/xerror"
	"github.com/zzy-rabbit/xtools/xruntime"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"sync"
)

type service struct {
	config Config
	mutex  sync.RWMutex
	level  int
	skip   int
	logger *log.Logger
	core   *lumberjack.Logger
}

var defaultLogger ILogger

func init() {
	ctx := xcontext.Background()
	logger, err := New(ctx, Config{})
	if err != nil {
		panic("init log fail " + err.Error())
	}
	defaultLogger = logger
}

func GetDefaultLogger(ctx context.Context) ILogger {
	return defaultLogger
}

func New(ctx context.Context, config Config) (ILogger, xerror.IError) {
	s := &service{config: config}
	setDefault(ctx, &s.config)

	if err := os.MkdirAll(s.config.Path, os.ModePerm); err != nil {
		fmt.Printf("mkdir %s fail %v\n", s.config.Path, err)
		return nil, xerror.Extend(xerror.ErrFileOperationFail, err.Error())
	}

	s.core = &lumberjack.Logger{
		Filename:   filepath.Join(s.config.Path, s.config.Name+s.config.Suffix),
		MaxSize:    s.config.MaxSize,
		MaxAge:     7,
		MaxBackups: 0,
		LocalTime:  false,
		Compress:   true,
	}
	s.logger = log.New(s.core, "", log.LstdFlags|log.Lmicroseconds)
	return s, nil
}

func (s *service) SetLevel(ctx context.Context, level int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.level = level
}

func (s *service) SetSkip(ctx context.Context, skip int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.skip = skip
}

func (s *service) Close(ctx context.Context) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_ = s.core.Close()
}

func (s *service) Printf(ctx context.Context, level int, format string, v ...any) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.level <= level {
		funcName, line := xruntime.GetShortFuncName(ctx, s.skip+2)
		s.logger.Printf("<"+funcName+":"+strconv.Itoa(line)+"> ["+levelDescMap[level]+"] trace: "+
			xcontext.GetTrace(ctx)+" "+format+"\n", v...)
	}
}

func (s *service) Debug(ctx context.Context, format string, v ...any) {
	s.Printf(ctx, LevelDebug, format, v...)
}

func (s *service) Info(ctx context.Context, format string, v ...any) {
	s.Printf(ctx, LevelInfo, format, v...)
}

func (s *service) Warn(ctx context.Context, format string, v ...any) {
	s.Printf(ctx, LevelWarn, format, v...)
}

func (s *service) Error(ctx context.Context, format string, v ...any) {
	s.Printf(ctx, LevelError, format, v...)
}

func (s *service) Fatal(ctx context.Context, format string, v ...any) {
	s.Printf(ctx, LevelFatal, format, v...)
	s.Printf(ctx, LevelFatal, "%s", debug.Stack())
	panic("fatal")
}

func (s *service) Stack(ctx context.Context) {
	s.Printf(ctx, LevelStack, "%s", debug.Stack())
}
