package xplugin

import (
	"context"
	"fmt"
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

func (s *service) Inject(ctx context.Context, target any) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
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
			return fmt.Errorf("plugin %q not found for field %s", pluginName, fieldType.Name)
		}

		pluginValue := reflect.ValueOf(plugin)

		// 类型校验
		if !pluginValue.Type().AssignableTo(fieldValue.Type()) {
			return fmt.Errorf("plugin %s type %s cannot be assigned to field %s (%s)", pluginName, pluginValue.Type(), fieldType.Name, fieldValue.Type())
		}

		fieldValue.Set(pluginValue)
	}
	return nil
}
