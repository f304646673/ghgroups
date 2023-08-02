package handlergroup

import (
	"fmt"
	"ghgroups/frame"
	debughelper "ghgroups/frame/debug_helper"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type HandlerGroupConf struct {
	Name     string   `yaml:"name"`
	Handlers []string `yaml:"handlers"`
}

type HandlerGroup struct {
	HandlerGroupInterface
	conf                 *HandlerGroupConf
	handlers             []frame.HandlerBaseInterface
	constructorInterface frame.ConstructorInterface
}

func NewHandlerGroup(constructor frame.ConstructorInterface) *HandlerGroup {
	return &HandlerGroup{
		handlers:             make([]frame.HandlerBaseInterface, 0),
		constructorInterface: constructor,
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// frame.HandlerBaseInterface
func (h *HandlerGroup) Handle(context *ghgroupscontext.GhGroupsContext) bool {
	for _, handler := range h.handlers {
		if handlerWithSaveDuration(handler, handler.Name(), context) {
			continue
		}
		return false
	}
	return true
}

func handlerWithSaveDuration(handlerBaseInterface frame.HandlerBaseInterface, handlerName string, ctx *ghgroupscontext.GhGroupsContext) bool {
	if ctx.ShowDuration {
		defer debughelper.DealDuration(time.Now(), handlerName, ctx)
	}
	return handlerBaseInterface.Handle(ctx)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (h *HandlerGroup) LoadConfigFromFile(confPath string) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	return h.LoadConfigFromMemory(data)
}

func (h *HandlerGroup) LoadConfigFromMemory(configure []byte) error {
	conf := new(HandlerGroupConf)
	err := yaml.Unmarshal([]byte(configure), conf)
	if err != nil {
		return err
	}
	h.conf = conf

	return h.initHandlers()
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// frame.ConcreteInterface
func (h *HandlerGroup) Name() string {
	return h.conf.Name
}

func (h *HandlerGroup) Add(handlderInterface frame.HandlerBaseInterface) error {
	h.handlers = append(h.handlers, handlderInterface)
	return nil
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (h *HandlerGroup) SetConstructorInterface(constructorInterface any) {
	constructorInterfaceNew, ok := constructorInterface.(frame.ConstructorInterface)
	if !ok {
		panic("constructorInterface is not frame.ConstructorInterface")
	}
	h.constructorInterface = constructorInterfaceNew
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (h *HandlerGroup) initHandlers() error {
	for _, handlerName := range h.conf.Handlers {
		if err := h.constructorInterface.CreateConcrete(handlerName); err != nil {
			return err
		}

		if someInterface, err := h.constructorInterface.GetConcrete(handlerName); err != nil {
			return err
		} else {
			if handlerInterface, ok := someInterface.(frame.HandlerBaseInterface); !ok {
				return fmt.Errorf("handler %s is not frame.HandlerBaseInterface", handlerName)
			} else {
				err = h.Add(handlerInterface)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
