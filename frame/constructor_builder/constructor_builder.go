package constructorbuilder

import (
	"ghgroups/frame/constructor"
	aynchandlergroupconstructor "ghgroups/frame/constructor/async_handler_group_constructor"
	dividerconstructor "ghgroups/frame/constructor/divider_constructor"
	handlerconstructor "ghgroups/frame/constructor/handler_constructor"
	handlergroupconstructor "ghgroups/frame/constructor/handler_group_constructor"
	layercenterconstructor "ghgroups/frame/constructor/layer_center_constructor"
	layerconstructor "ghgroups/frame/constructor/layer_constructor"
	"ghgroups/frame/factory"
	"ghgroups/frame/layer"
	layercenter "ghgroups/frame/layer_center"

	asynchandlergroup "ghgroups/frame/async_handler_group"
	handlergroup "ghgroups/frame/handler_group"
	"reflect"
)

func BuildConstructor(factory *factory.Factory, concretePath string) *constructor.Constructor {
	factory.Register(reflect.TypeOf(asynchandlergroup.AsyncHandlerGroup{}))
	factory.Register(reflect.TypeOf(handlergroup.HandlerGroup{}))
	factory.Register(reflect.TypeOf(layer.Layer{}))
	factory.Register(reflect.TypeOf(layercenter.LayerCenter{}))
	factory.Register(reflect.TypeOf(layerconstructor.LayerConstructor{}))
	factory.Register(reflect.TypeOf(dividerconstructor.DividerConstructor{}))
	factory.Register(reflect.TypeOf(handlerconstructor.HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	return constructor.NewConstructor(factory, concretePath)
}
