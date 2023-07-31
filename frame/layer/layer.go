package layer

import (
	"fmt"
	"ghgroups/frame"
	"os"

	ghgroupscontext "ghgroups/frame/ghgroups_context"

	"gopkg.in/yaml.v2"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type LayerConf struct {
	Name     string   `yaml:"name"`
	Divider  string   `yaml:"divider"`
	Handlers []string `yaml:"handlers"`
}

type Layer struct {
	frame.LayerWithBuilderInterface
	conf                 LayerConf
	divider              frame.DividerBaseInterface
	handlers             map[string]frame.HandlerBaseInterface
	constructorInterface frame.ConstructorInterface
}

func NewLayer(name string, constructorInterface frame.ConstructorInterface) *Layer {
	return &Layer{
		constructorInterface: constructorInterface,
		conf:                 LayerConf{Name: name},
		handlers:             make(map[string]frame.HandlerBaseInterface),
	}
}

func (l *Layer) SetDivider(name string, t frame.DividerBaseInterface) error {
	if l.divider == nil {
		l.divider = t
	} else {
		return fmt.Errorf("divider %s already exists", name)
	}
	return nil
}

func (l *Layer) AddHandler(name string, h frame.HandlerBaseInterface) error {
	if _, ok := l.handlers[name]; ok {
		return fmt.Errorf("handler %s already exists", name)
	}
	l.handlers[name] = h
	return nil
}

func (l *Layer) initDivider(dividerName string) error {
	if err := l.constructorInterface.CreateConcrete(dividerName); err != nil {
		return err
	}

	if someInterface, err := l.constructorInterface.GetConcrete(dividerName); err != nil {
		return err
	} else {
		if dividerInterface, ok := someInterface.(frame.DividerBaseInterface); !ok {
			return fmt.Errorf("handler %s is not frame.DividerBaseInterface", dividerName)
		} else {
			err = l.SetDivider(dividerName, dividerInterface)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *Layer) initHandlers(handlersName []string) error {
	for _, handlerName := range handlersName {

		if err := l.constructorInterface.CreateConcrete(handlerName); err != nil {
			return err
		}

		if someInterface, err := l.constructorInterface.GetConcrete(handlerName); err != nil {
			return err
		} else {
			if handlerInterface, ok := someInterface.(frame.HandlerBaseInterface); !ok {
				return fmt.Errorf("handler %s is not frame.HandlerBaseInterface", handlerName)
			} else {
				err = l.AddHandler(handlerName, handlerInterface)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (l *Layer) init() error {
	if l.handlers == nil {
		l.handlers = make(map[string]frame.HandlerBaseInterface)
	}
	err := l.initDivider(l.conf.Divider)
	if err != nil {
		return err
	}

	err = l.initHandlers(l.conf.Handlers)
	if err != nil {
		return err
	}
	return nil
}

func (l *Layer) SetConstructorInterface(constructorInterface any) {
	constructorInterfaceNew, ok := constructorInterface.(frame.ConstructorInterface)
	if !ok {
		panic("constructorInterface is not frame.ConstructorInterface")
	}
	l.constructorInterface = constructorInterfaceNew
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (l *Layer) Name() string {
	return l.conf.Name
}

func (l *Layer) Handle(ctx *ghgroupscontext.GhGroupsContext) bool {
	layerName := l.divider.Select(ctx)
	if handler, ok := l.handlers[layerName]; !ok {
		return false
	} else {
		return handlerWithSaveDuration(handler, layerName, ctx)
	}
}

func handlerWithSaveDuration(handlerBaseInterface frame.HandlerBaseInterface, layeName string, ctx *ghgroupscontext.GhGroupsContext) bool {
	// if ctx.DebugFlag {
	// 	defer debughelper.SaveDurationToContext(time.Now(), layeName, ctx)
	// }
	return handlerBaseInterface.Handle(ctx)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (l *Layer) LoadConfigFromFile(confPath string) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}
	return l.LoadConfigFromMemory(data)
}

func (l *Layer) LoadConfigFromMemory(configure []byte) error {
	var layerConf LayerConf
	err := yaml.Unmarshal(configure, &layerConf)
	if err != nil {
		return err
	}
	l.conf = layerConf
	return l.init()
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
