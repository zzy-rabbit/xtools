package xplugin

import (
	"context"
	"github.com/zzy-rabbit/xtools/xerror"
	"reflect"
	"sync"
)

type service struct {
	config    Config
	mutex     sync.RWMutex
	pluginMap map[string]IPlugin
}

var instance = &service{
	pluginMap: make(map[string]IPlugin),
}

func (s *service) Save(ctx context.Context, plugin IPlugin) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.pluginMap[plugin.GetName(ctx)] = plugin
}

func (s *service) Get(ctx context.Context, name string) (IPlugin, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	plugin, ok := s.pluginMap[name]
	return plugin, ok
}

func (s *service) Inject(ctx context.Context, target any) xerror.IError {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return xerror.Extend(xerror.ErrInvalidParam, "target must be a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		// 跳过不可设置字段（未导出）
		if !fieldValue.CanSet() {
			continue
		}

		pluginName := fieldType.Tag.Get("xplugin")
		if pluginName == "" {
			continue
		}

		plugin, ok := s.pluginMap[pluginName]
		if !ok {
			continue
		}

		pluginValue := reflect.ValueOf(plugin)

		// 类型校验
		if !pluginValue.Type().AssignableTo(fieldValue.Type()) {
			continue
		}

		fieldValue.Set(pluginValue)
	}
	return nil
}
