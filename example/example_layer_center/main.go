package main

import (
	"fmt"
	examplelayera "ghgroups/example/example_layer_center/example_layer_a"
	examplelayerb "ghgroups/example/example_layer_center/example_layer_b"
	"ghgroups/frame"
	constructorbuilder "ghgroups/frame/constructor_builder"
	"ghgroups/frame/factory"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
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
	factory.Register(reflect.TypeOf(examplelayera.ExampleA1Handler{}))
	factory.Register(reflect.TypeOf(examplelayera.ExampleA2Handler{}))
	factory.Register(reflect.TypeOf(examplelayerb.ExampleB1Handler{}))
	factory.Register(reflect.TypeOf(examplelayerb.ExampleB2Handler{}))
	factory.Register(reflect.TypeOf(examplelayera.ExampleADivider{}))

	factory.Register(reflect.TypeOf(examplelayerb.ExampleBDivider{}))

	constructor := constructorbuilder.BuildConstructor(factory, concretePath)

	mainProcess := "layer_center"
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
