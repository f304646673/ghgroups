package frame

import (
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type LoadConfigFromFileInterface interface {
	LoadConfigFromFile(filePath string) error
}

type LoadConfigFromMemoryInterface interface {
	LoadConfigFromMemory(configure []byte) error
}

type LoadEnvironmentConfInterface interface {
	LoadEnvironmentConf() error
}

type ConcreteInterface interface {
	Name() string
}

type FactoryInterface interface {
	Register(reflect.Type) error
	Create(string, []byte, any) (any, error)
	Get(string) (ConcreteInterface, error)
}

type HandlerBaseInterface interface {
	ConcreteInterface
	Handle(context *ghgroupscontext.GhGroupsContext) bool
}

type DividerBaseInterface interface {
	ConcreteInterface
	Select(context *ghgroupscontext.GhGroupsContext) string
}

type LayerBaseInterface interface {
	HandlerBaseInterface
}

type LayerWithBuilderInterface interface {
	LayerBaseInterface
	ConstructorSetterInterface
	SetDivider(string, DividerBaseInterface) error
	AddHandler(string, HandlerBaseInterface) error
}

type DividerConstructorInterface interface {
	GetDivider(name string) (DividerBaseInterface, error)
	RegisterDivider(name string, dividerInterface DividerBaseInterface) error
	ParseDividerConfFolder(confFolderPath string) error
	CreateDividerWithConfPath(confFilePath string) error
}

type HandlerConstructorInterface interface {
	GetHandler(name string) (HandlerBaseInterface, error)
	RegisterHandler(name string, handlerInterface HandlerBaseInterface) error
	ParseHandlerConfFolder(confFolderPath string) error
	CreateHandlerWithConfPath(confFilePath string) error
}

type LayerConstructorInterface interface {
	GetLayer(name string) (LayerBaseInterface, error)
	RegisterLayer(name string, layerInterface LayerBaseInterface) error
	ParseLayerConfFolder(confFolderPath string) error
	CreateLayerWithConfPath(confFilePath string) error
}

type ConstructorInterface interface {
	LayerConstructorInterface
	DividerConstructorInterface
	HandlerConstructorInterface
	LayerCenterConstructorInterface
	FactoryInterface
	CreateConcrete(string) error
	GetConcrete(string) (any, error)
}

type ConstructorSetterInterface interface {
	SetConstructorInterface(any)
}

type LayerCenterConstructorInterface interface {
	GetLayerCenter(name string) (LayerCenterBaseInterface, error)
	RegisterLayerCenter(name string, layerCenterBaseInterface LayerCenterBaseInterface) error
	ParseLayerCenterConfFolder(confFolderPath string) error
	CreateLayerCenterWithConfPath(confFilePath string) error
}

type LayerCenterBaseInterface interface {
	HandlerBaseInterface
}

type LayerCenterInterface interface {
	LayerCenterBaseInterface
	Add(LayerWithBuilderInterface)
}

type HandlerGroupConstructorInterface interface {
	GetHandlerGroup(name string) (HandlerGroupBaseInterface, error)
	RegisterHandlerGroup(name string, handlerGroupInterface HandlerGroupBaseInterface) error
	ParseHandlerGroupConfFolder(confFolderPath string) error
	CreateHandlerGroupWithConfPath(confFilePath string) error
}

type HandlerGroupBaseInterface interface {
	HandlerBaseInterface
}

type AsyncHandlerGroupConstructorInterface interface {
	GetAsyncHandlerGroup(name string) (AsyncHandlerGroupBaseInterface, error)
	RegisterAsyncHandlerGroup(name string, handlerGroupInterface AsyncHandlerGroupBaseInterface) error
	ParseAsyncHandlerGroupConfFolder(confFolderPath string) error
	CreateAsyncHandlerGroupWithConfPath(confFilePath string) error
}

type AsyncHandlerGroupBaseInterface interface {
	HandlerBaseInterface
}
