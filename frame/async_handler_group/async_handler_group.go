package asynchandlergroup

import (
	"fmt"
	"ghgroups/frame"
	"os"

	debughelper "ghgroups/frame/debug_helper"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"sync"

	"gopkg.in/yaml.v2"
)

type HandlerGroupConf struct {
	Name     string   `yaml:"name"`
	Handlers []string `yaml:"handlers"`
}

type AsyncHandlerGroup struct {
	HandlerGroupInterface
	conf                 *HandlerGroupConf
	handlers             []frame.HandlerBaseInterface
	constructorInterface frame.ConstructorInterface
}

func NewAsyncHandlerGroup(constructor frame.ConstructorInterface) *AsyncHandlerGroup {
	return &AsyncHandlerGroup{
		handlers:             make([]frame.HandlerBaseInterface, 0),
		constructorInterface: constructor,
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// frame.HandlerBaseInterface
func (a *AsyncHandlerGroup) Handle(context *ghgroupscontext.GhGroupsContext) bool {
	wg := sync.WaitGroup{}
	checkChan := make(chan bool, len(a.handlers))
	for _, handler := range a.handlers {
		wg.Add(1)
		go func(handler frame.HandlerBaseInterface) {
			status := debughelper.HandleWithShowDuration(handler, handler.Name(), context)
			checkChan <- status
			wg.Done()
		}(handler)

	}
	wg.Wait()
	for {
		select {
		case status := <-checkChan:
			if !status {
				return false
			}
		default:
			return true
		}
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (a *AsyncHandlerGroup) LoadConfigFromFile(confPath string) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	return a.LoadConfigFromMemory(data)
}

func (a *AsyncHandlerGroup) LoadConfigFromMemory(configure []byte) error {
	conf := new(HandlerGroupConf)
	err := yaml.Unmarshal([]byte(configure), conf)
	if err != nil {
		return err
	}
	a.conf = conf

	return a.initHandlers()
}

func (a *AsyncHandlerGroup) initHandlers() error {
	for _, handlerName := range a.conf.Handlers {
		if err := a.constructorInterface.CreateConcrete(handlerName); err != nil {
			return err
		}

		if someInterface, err := a.constructorInterface.GetConcrete(handlerName); err != nil {
			return err
		} else {
			if handlerInterface, ok := someInterface.(frame.HandlerBaseInterface); !ok {
				return fmt.Errorf("handler %s is not frame.HandlerBaseInterface", handlerName)
			} else {
				err = a.Add(handlerInterface)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// frame.ConcreteInterface
func (a *AsyncHandlerGroup) Name() string {
	return a.conf.Name
}

func (a *AsyncHandlerGroup) Add(handlderInterface frame.HandlerBaseInterface) error {
	a.handlers = append(a.handlers, handlderInterface)
	return nil
}

func (a *AsyncHandlerGroup) SetConstructorInterface(constructorInterface any) {
	constructorInterfaceNew, ok := constructorInterface.(frame.ConstructorInterface)
	if !ok {
		panic("constructorInterface is not frame.ConstructorInterface")
	}
	a.constructorInterface = constructorInterfaceNew
}
