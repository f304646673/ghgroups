package constructor

import (
	"fmt"
	"ghgroups/frame"
	"os"
	"reflect"
	"strings"

	concreteconfmanager "ghgroups/frame/concrete_conf_manager"

	"gopkg.in/yaml.v2"
)

type Constructor struct {
	frame.ConstructorInterface

	layerConstructorInterface             frame.LayerConstructorInterface
	dividerConstructorInterface           frame.DividerConstructorInterface
	handlerConstructorInterface           frame.HandlerConstructorInterface
	asyncHandlerGroupConstructorInterface frame.AsyncHandlerGroupConstructorInterface
	handlerGroupConstructorInterface      frame.HandlerGroupConstructorInterface
	layerCenterConstructorInterface       frame.LayerCenterConstructorInterface
	factoryInterface                      frame.FactoryInterface
	concreteConfManager                   *concreteconfmanager.ConcreteConfManager
	deepth                                int
}

func NewConstructor(factory frame.FactoryInterface, confPath string) *Constructor {
	constructor := &Constructor{}

	// 必须使用反射去做，否则有可能会出现循环依赖
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	layerConstructor, err := factory.Create("LayerConstructor", []byte{}, constructor)
	if err != nil {
		panic(err)
	}
	layerConcreteInterface, ok := layerConstructor.(frame.LoadEnvironmentConfInterface)
	if !ok {
		panic("layerConstructor is not frame.LoadEnvironmentConfInterface")
	}
	err = layerConcreteInterface.LoadEnvironmentConf()
	if err != nil {
		panic(err)
	}
	layerConstructorInterface, ok := layerConstructor.(frame.LayerConstructorInterface)
	if !ok {
		panic("layerConstructor is not frame.LayerConstructorInterface")
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////////
	dividerConstructor, err := factory.Create("DividerConstructor", []byte{}, constructor)
	if err != nil {
		panic(err)
	}

	dividerConcreteInterface, ok := dividerConstructor.(frame.LoadEnvironmentConfInterface)
	if !ok {
		panic("dividerConstructor is not frame.LoadEnvironmentConfInterface")
	}

	err = dividerConcreteInterface.LoadEnvironmentConf()
	if err != nil {
		panic(err)
	}

	dividerConstructorInterface, ok := dividerConstructor.(frame.DividerConstructorInterface)
	if !ok {
		panic("dividerConstructor is not frame.DividerConstructorInterface")
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	handlerConstructor, err := factory.Create("HandlerConstructor", []byte{}, constructor)
	if err != nil {
		panic(err)
	}
	handlerConcreteInterface, ok := handlerConstructor.(frame.LoadEnvironmentConfInterface)
	if !ok {
		panic("handlerConstructor is not frame.LoadEnvironmentConfInterface")
	}
	err = handlerConcreteInterface.LoadEnvironmentConf()
	if err != nil {
		panic(err)
	}
	handlerConstructorInterface, ok := handlerConstructor.(frame.HandlerConstructorInterface)
	if !ok {
		panic("handlerConstructor is not frame.HandlerConstructorInterface")
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	layerCenterConstructor, err := factory.Create("LayerCenterConstructor", []byte{}, constructor)
	if err != nil {
		panic(err)
	}
	layerCenterConcreteInterface, ok := layerCenterConstructor.(frame.LoadEnvironmentConfInterface)
	if !ok {
		panic("layerCenterConstructor is not frame.LoadEnvironmentConfInterface")
	}
	err = layerCenterConcreteInterface.LoadEnvironmentConf()
	if err != nil {
		panic(err)
	}
	layerCenterConstructorInterface, ok := layerCenterConstructor.(frame.LayerCenterConstructorInterface)
	if !ok {
		panic("layerCenterConstructor is not frame.LayerCenterConstructorInterface")
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	handlerGroupConstructor, err := factory.Create("HandlerGroupConstructor", []byte{}, constructor)
	if err != nil {
		panic(err)
	}
	handlerGroupConcreteInterface, ok := handlerGroupConstructor.(frame.LoadEnvironmentConfInterface)
	if !ok {
		panic("handlerGroupConstructor is not frame.LoadEnvironmentConfInterface")
	}
	err = handlerGroupConcreteInterface.LoadEnvironmentConf()
	if err != nil {
		panic(err)
	}
	handlerGroupConstructorInterface, ok := handlerGroupConstructor.(frame.HandlerGroupConstructorInterface)
	if !ok {
		panic("handlerGroupConstructor is not frame.HandlerGroupConstructorInterface")
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////、
	asyncHandlerGroupConstructor, err := factory.Create("AsyncHandlerGroupConstructor", []byte{}, constructor)
	if err != nil {
		panic(err)
	}
	asyncHandlerGroupConcreteInterface, ok := asyncHandlerGroupConstructor.(frame.LoadEnvironmentConfInterface)
	if !ok {
		panic("asyncHandlerGroupConstructor is not frame.LoadEnvironmentConfInterface")
	}
	err = asyncHandlerGroupConcreteInterface.LoadEnvironmentConf()
	if err != nil {
		panic(err)
	}
	asyncHandlerGroupConstructorInterface, ok := asyncHandlerGroupConstructor.(frame.AsyncHandlerGroupConstructorInterface)
	if !ok {
		panic("asyncHandlerGroupConstructor is not frame.AsyncHandlerGroupConstructorInterface")
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	concreteConfManager := concreteconfmanager.NewConcreteConfManager()
	if confPath != "" {
		err = concreteConfManager.ParseConfFolder(confPath)
		if err != nil {
			panic("concreteConfManager.ParseConfFolder error")
		}
	}

	constructor.factoryInterface = factory
	constructor.layerConstructorInterface = layerConstructorInterface
	constructor.dividerConstructorInterface = dividerConstructorInterface
	constructor.handlerConstructorInterface = handlerConstructorInterface
	constructor.layerCenterConstructorInterface = layerCenterConstructorInterface
	constructor.handlerGroupConstructorInterface = handlerGroupConstructorInterface
	constructor.concreteConfManager = concreteConfManager
	constructor.asyncHandlerGroupConstructorInterface = asyncHandlerGroupConstructorInterface
	constructor.deepth = -1

	layerConstructorSetterInterface, ok := layerConstructor.(frame.ConstructorSetterInterface)
	if !ok {
		panic("layerConstructor is not frame.ConstructorSetterInterface")
	}
	layerConstructorSetterInterface.SetConstructorInterface(constructor)

	dividerConstructorSetterInterface, ok := dividerConstructor.(frame.ConstructorSetterInterface)
	if !ok {
		panic("dividerConstructor is not frame.ConstructorSetterInterface")
	}
	dividerConstructorSetterInterface.SetConstructorInterface(constructor)

	handlerConstructorSetterInterface, ok := handlerConstructor.(frame.ConstructorSetterInterface)
	if !ok {
		panic("handlerConstructor is not frame.ConstructorSetterInterface")
	}
	handlerConstructorSetterInterface.SetConstructorInterface(constructor)

	return constructor
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// LayerConstructorInterface
func (c *Constructor) GetLayer(name string) (frame.LayerBaseInterface, error) {
	return c.layerConstructorInterface.GetLayer(name)
}

func (c *Constructor) RegisterLayer(name string, layerInterface frame.LayerBaseInterface) error {
	return c.layerConstructorInterface.RegisterLayer(name, layerInterface)
}

func (c *Constructor) ParseLayerConfFolder(confFolderPath string) error {
	return c.layerConstructorInterface.ParseLayerConfFolder(confFolderPath)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// DividerConstructorInterface
func (c *Constructor) GetDivider(name string) (frame.DividerBaseInterface, error) {
	return c.dividerConstructorInterface.GetDivider(name)
}

func (c *Constructor) RegisterDivider(name string, dividerInterface frame.DividerBaseInterface) error {
	return c.dividerConstructorInterface.RegisterDivider(name, dividerInterface)
}

func (c *Constructor) ParseDividerConfFolder(confFolderPath string) error {
	return c.dividerConstructorInterface.ParseDividerConfFolder(confFolderPath)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// HandlerConstructorInterface
func (c *Constructor) GetHandler(name string) (frame.HandlerBaseInterface, error) {
	return c.handlerConstructorInterface.GetHandler(name)
}

func (c *Constructor) RegisterHandler(name string, handlerInterface frame.HandlerBaseInterface) error {
	return c.handlerConstructorInterface.RegisterHandler(name, handlerInterface)
}

func (c *Constructor) ParseHandlerConfFolder(confFolderPath string) error {
	return c.handlerConstructorInterface.ParseHandlerConfFolder(confFolderPath)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// LayerCenterConstructorInterface
func (c *Constructor) GetLayerCenter(name string) (frame.LayerCenterBaseInterface, error) {
	return c.layerCenterConstructorInterface.GetLayerCenter(name)
}

func (c *Constructor) RegisterLayerCenter(name string, layerCenterInterface frame.LayerCenterBaseInterface) error {
	return c.layerCenterConstructorInterface.RegisterLayerCenter(name, layerCenterInterface)
}

func (c *Constructor) ParseLayerCenterConfFolder(confFolderPath string) error {
	return c.layerCenterConstructorInterface.ParseLayerCenterConfFolder(confFolderPath)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// HandlerGroupConstructorInterface
func (c *Constructor) GetHandlerGroup(name string) (frame.HandlerGroupBaseInterface, error) {
	return c.handlerGroupConstructorInterface.GetHandlerGroup(name)
}

func (c *Constructor) RegisterHandlerGroup(name string, handlerGroupInterface frame.HandlerGroupBaseInterface) error {
	return c.handlerGroupConstructorInterface.RegisterHandlerGroup(name, handlerGroupInterface)
}

func (c *Constructor) ParseHandlerGroupConfFolder(confFolderPath string) error {
	return c.handlerGroupConstructorInterface.ParseHandlerGroupConfFolder(confFolderPath)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// AsyncHandlerConstructorInterface
func (c *Constructor) GetAsyncHandlerGroup(name string) (frame.AsyncHandlerGroupBaseInterface, error) {
	return c.asyncHandlerGroupConstructorInterface.GetAsyncHandlerGroup(name)
}

func (c *Constructor) RegisterAsyncHandlerGroup(name string, asyncAsyncHandlerGroupInterface frame.AsyncHandlerGroupBaseInterface) error {
	return c.asyncHandlerGroupConstructorInterface.RegisterAsyncHandlerGroup(name, asyncAsyncHandlerGroupInterface)
}

func (c *Constructor) ParseAsyncHandlerGroupConfFolder(confFolderPath string) error {
	return c.asyncHandlerGroupConstructorInterface.ParseAsyncHandlerGroupConfFolder(confFolderPath)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// FactoryInterface
func (c *Constructor) Register(concreteType reflect.Type) error {
	return c.factoryInterface.Register(concreteType)
}

func (c *Constructor) Create(name string, conf []byte, itf any) (any, error) {
	return c.factoryInterface.Create(name, conf, c)
}

func (c *Constructor) Get(name string) (frame.ConcreteInterface, error) {
	return c.factoryInterface.Get(name)
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type ConstructorType struct {
	Type string `yaml:"type"`
}

const TypeNameHandler = "Handler"
const TypeNameDivider = "Divider"
const TypeNameLayer = "Layer"
const TypeNameHandlerGroup = "HandlerGroup"
const TypeNameLayerCenter = "LayerCenter"
const TypeNameAsyncHandlerGroup = "AsyncHandlerGroup"

func (c *Constructor) createConcreteByTypeName(name string, data []byte) error {
	if strings.HasSuffix(name, TypeNameAsyncHandlerGroup) {
		_, err := c.asyncHandlerGroupConstructorInterface.GetAsyncHandlerGroup(name)
		if err != nil {
			asyncHandlerGroupInterface, err := c.Create(name, data, c)
			if err != nil {
				return err
			}
			if namedInterface, ok := asyncHandlerGroupInterface.(frame.ConcreteInterface); ok {
				return c.asyncHandlerGroupConstructorInterface.RegisterAsyncHandlerGroup(namedInterface.Name(), asyncHandlerGroupInterface.(frame.AsyncHandlerGroupBaseInterface))
			} else {
				return c.asyncHandlerGroupConstructorInterface.RegisterAsyncHandlerGroup(name, asyncHandlerGroupInterface.(frame.AsyncHandlerGroupBaseInterface))
			}
		}
		return nil
	}
	if strings.HasPrefix(name, TypeNameHandlerGroup) {
		_, err := c.handlerGroupConstructorInterface.GetHandlerGroup(name)
		if err != nil {
			handlerGroupInterface, err := c.Create(name, data, c)
			if err != nil {
				return err
			}
			if namedInterface, ok := handlerGroupInterface.(frame.ConcreteInterface); ok {
				return c.handlerGroupConstructorInterface.RegisterHandlerGroup(namedInterface.Name(), handlerGroupInterface.(frame.HandlerGroupBaseInterface))
			} else {
				return c.handlerGroupConstructorInterface.RegisterHandlerGroup(name, handlerGroupInterface.(frame.HandlerGroupBaseInterface))
			}
		}
		return nil
	}
	if strings.HasSuffix(name, TypeNameDivider) {
		_, err := c.dividerConstructorInterface.GetDivider(name)
		if err != nil {
			dividerInterface, err := c.Create(name, data, c)
			if err != nil {
				return err
			}
			if namedInterface, ok := dividerInterface.(frame.ConcreteInterface); ok {
				return c.dividerConstructorInterface.RegisterDivider(namedInterface.Name(), dividerInterface.(frame.DividerBaseInterface))
			} else {
				return c.dividerConstructorInterface.RegisterDivider(name, dividerInterface.(frame.DividerBaseInterface))
			}
		}
		return nil
	}
	if strings.HasSuffix(name, TypeNameHandler) {
		_, err := c.handlerConstructorInterface.GetHandler(name)
		if err != nil {
			handlerInterface, err := c.Create(name, data, c)
			if err != nil {
				return err
			}
			if namedInterface, ok := handlerInterface.(frame.ConcreteInterface); ok {
				return c.handlerConstructorInterface.RegisterHandler(namedInterface.Name(), handlerInterface.(frame.HandlerBaseInterface))
			} else {
				return c.handlerConstructorInterface.RegisterHandler(name, handlerInterface.(frame.HandlerBaseInterface))
			}
		}
		return nil
	}
	if strings.HasSuffix(name, TypeNameLayer) {
		_, err := c.layerConstructorInterface.GetLayer(name)
		if err != nil {
			layerInterface, err := c.Create(name, data, c)
			if err != nil {
				return err
			}
			if namedInterface, ok := layerInterface.(frame.ConcreteInterface); ok {
				return c.layerConstructorInterface.RegisterLayer(namedInterface.Name(), layerInterface.(frame.LayerBaseInterface))
			} else {
				return c.layerConstructorInterface.RegisterLayer(name, layerInterface.(frame.LayerBaseInterface))
			}
		}
		return nil
	}
	if strings.HasSuffix(name, TypeNameLayerCenter) {
		_, err := c.layerCenterConstructorInterface.GetLayerCenter(name)
		if err != nil {
			layerCenterInterface, err := c.Create(name, data, c)
			if err != nil {
				return err
			}
			if namedInterface, ok := layerCenterInterface.(frame.ConcreteInterface); ok {
				return c.layerCenterConstructorInterface.RegisterLayerCenter(namedInterface.Name(), layerCenterInterface.(frame.LayerCenterBaseInterface))
			} else {
				return c.layerCenterConstructorInterface.RegisterLayerCenter(name, layerCenterInterface.(frame.LayerCenterBaseInterface))
			}
		}
		return nil
	}

	return fmt.Errorf("object name  %s not found", name)
}

func (c *Constructor) createConcreteByObjectName(name string) error {
	confPath, err := c.concreteConfManager.GetConfPath(name)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	var constructorType ConstructorType
	if err := yaml.Unmarshal(data, &constructorType); err != nil {
		return err
	}

	switch constructorType.Type {
	case TypeNameHandler:
		_, err := c.handlerConstructorInterface.GetHandler(name)
		if err != nil {
			return c.handlerConstructorInterface.CreateHandlerWithConfPath(confPath)
		}
	case TypeNameDivider:
		_, err := c.dividerConstructorInterface.GetDivider(name)
		if err != nil {
			return c.dividerConstructorInterface.CreateDividerWithConfPath(confPath)
		}
	case TypeNameLayer:
		_, err := c.layerConstructorInterface.GetLayer(name)
		if err != nil {
			return c.layerConstructorInterface.CreateLayerWithConfPath(confPath)
		}
	case TypeNameLayerCenter:
		_, err := c.layerCenterConstructorInterface.GetLayerCenter(name)
		if err != nil {
			return c.layerCenterConstructorInterface.CreateLayerCenterWithConfPath(confPath)
		}
	case TypeNameHandlerGroup:
		_, err := c.handlerGroupConstructorInterface.GetHandlerGroup(name)
		if err != nil {
			return c.handlerGroupConstructorInterface.CreateHandlerGroupWithConfPath(confPath)
		}
	case TypeNameAsyncHandlerGroup:
		_, err := c.asyncHandlerGroupConstructorInterface.GetAsyncHandlerGroup(name)
		if err != nil {
			return c.asyncHandlerGroupConstructorInterface.CreateAsyncHandlerGroupWithConfPath(confPath)
		}
	default:
		err := c.createConcreteByTypeName(constructorType.Type, data)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("class name  %s not found", name)
}

func (c *Constructor) CreateConcrete(name string) error {
	c.deepth++
	for i := 0; i < c.deepth; i++ {
		fmt.Printf("\t")
	}
	fmt.Println(name)

	defer func() {
		c.deepth--
	}()

	_, err := c.GetConcrete(name)
	if err == nil {
		return nil
	}

	if _, err := c.concreteConfManager.GetConfPath(name); err == nil {
		err = c.createConcreteByObjectName(name)
		if err != nil {
			return err
		}
		return nil
	}

	return c.createConcreteByTypeName(name, []byte{})
}

func (c *Constructor) GetConcrete(name string) (any, error) {
	if asyncHandlerBaseInterface, err := c.asyncHandlerGroupConstructorInterface.GetAsyncHandlerGroup(name); err == nil {
		return asyncHandlerBaseInterface, nil
	} else if handlerBaseInterface, err := c.handlerConstructorInterface.GetHandler(name); err == nil {
		return handlerBaseInterface, nil
	} else if dividerBaseInterface, err := c.dividerConstructorInterface.GetDivider(name); err == nil {
		return dividerBaseInterface, nil
	} else if layerBaseInterface, err := c.layerConstructorInterface.GetLayer(name); err == nil {
		return layerBaseInterface, nil
	} else if layerCenterBaseInterface, err := c.layerCenterConstructorInterface.GetLayerCenter(name); err == nil {
		return layerCenterBaseInterface, nil
	} else if handlerGroupBaseInterface, err := c.handlerGroupConstructorInterface.GetHandlerGroup(name); err == nil {
		return handlerGroupBaseInterface, nil
	}

	return nil, fmt.Errorf("concrete %s not found", name)
}
