package xplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zzy-rabbit/xtools/xerror"
)

const TagName = "xplugin"

type PluginConfig struct {
	PluginName string `json:"plugin_name"`
	InitParam  string `json:"init_param"`
	RunParam   string `json:"run_param"`
	StopParam  string `json:"stop_param"`
}

type Config struct {
	Plugins []PluginConfig `json:"plugins"`
}

type IPlugin interface {
	GetName(ctx context.Context) string
	Init(ctx context.Context, initParam string) error
	Run(ctx context.Context, runParam string) error
	Stop(ctx context.Context, stopParam string) error
}

func ParseConfig(ctx context.Context, content []byte) xerror.IError {
	err := json.Unmarshal(content, &instance.config)
	if err != nil {
		return xerror.Extend(xerror.ErrInvalidParam, err.Error())
	}
	return nil
}

func Inject(ctx context.Context, obj any) xerror.IError {
	return instance.Inject(ctx, obj)
}

func Register(ctx context.Context, plugin IPlugin) xerror.IError {
	err := instance.Inject(ctx, plugin)
	if err != nil {
		return err
	}
	instance.Save(ctx, plugin)
	fmt.Println("register plugin:", plugin.GetName(ctx))
	return nil
}

func Init(ctx context.Context) xerror.IError {
	for _, obj := range instance.config.Plugins {
		plugin, ok := instance.Get(ctx, obj.PluginName)
		if !ok {
			continue
		}
		err := plugin.Init(ctx, obj.InitParam)
		if err != nil {
			continue
		}
		fmt.Println("init plugin:", plugin.GetName(ctx))
	}
	return nil
}

func Get(ctx context.Context, name string) (IPlugin, bool) {
	return instance.Get(ctx, name)
}

func Run(ctx context.Context) xerror.IError {
	for _, obj := range instance.config.Plugins {
		plugin, ok := instance.Get(ctx, obj.PluginName)
		if !ok {
			continue
		}
		err := plugin.Run(ctx, obj.RunParam)
		if err != nil {
			continue
		}
		fmt.Println("run plugin:", plugin.GetName(ctx))
	}
	return nil
}

func Stop(ctx context.Context) xerror.IError {
	for _, obj := range instance.config.Plugins {
		plugin, ok := instance.Get(ctx, obj.PluginName)
		if !ok {
			continue
		}
		err := plugin.Stop(ctx, obj.StopParam)
		if err != nil {
			continue
		}
		fmt.Println("stop plugin:", plugin.GetName(ctx))
	}
	return nil
}
