package main

import (
	"fmt"
	"ghgroups/frame"
	"ghgroups/frame/constructor"
	aynchandlergroupconstructor "ghgroups/frame/constructor/async_handler_group_constructor"
	dividerconstructor "ghgroups/frame/constructor/divider_constructor"
	handlerconstructor "ghgroups/frame/constructor/handler_constructor"
	handlergroupconstructor "ghgroups/frame/constructor/handler_group_constructor"
	layercenterconstructor "ghgroups/frame/constructor/layer_center_constructor"
	layerconstructor "ghgroups/frame/constructor/layer_constructor"
	"ghgroups/frame/factory"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	handlergroup "ghgroups/frame/handler_group"
	"ghgroups/frame/layer"
	layercenter "ghgroups/frame/layer_center"
	"os"
	"path"
	"reflect"
)

func main() {
	runPath, errGetWd := os.Getwd()
	if errGetWd != nil {
		fmt.Printf("%v", errGetWd)
		return
	}
	concretePath := path.Join(runPath, "conf")
	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(handlergroup.HandlerGroup{}))
	factory.Register(reflect.TypeOf(layer.Layer{}))

	factory.Register(reflect.TypeOf(layercenter.LayerCenter{}))
	factory.Register(reflect.TypeOf(layerconstructor.LayerConstructor{}))
	factory.Register(reflect.TypeOf(dividerconstructor.DividerConstructor{}))
	factory.Register(reflect.TypeOf(handlerconstructor.HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(ExampleAHandler{}))
	factory.Register(reflect.TypeOf(ExampleBHandler{}))

	constructor := constructor.NewConstructor(factory, concretePath)

	mainProcess := "handler_group_a"
	if err := constructor.CreateConcrete(mainProcess); err != nil {
		fmt.Printf("%v", err)
	}
	if someInterfaced, err := constructor.GetConcrete(mainProcess); err != nil {
		fmt.Printf("%v", err)
	} else {
		if mainHandlerGroup, ok := someInterfaced.(frame.HandlerBaseInterface); !ok {
			fmt.Printf("mainHandlerGroup %s is not frame.HandlerBaseInterface", mainProcess)
		} else {
			context := ghgroupscontext.NewGhGroupsContext(nil)
			mainHandlerGroup.Handle(context)
		}
	}
}
