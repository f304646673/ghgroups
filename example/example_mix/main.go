package main

import (
	"fmt"
	examplelayera "ghgroups/example/example_mix/example_layer_center_a/example_layer_a"
	examplelayerb "ghgroups/example/example_mix/example_layer_center_a/example_layer_b"

	exampleasynchandlergroupf "ghgroups/example/example_mix/example_async_handler_group_f"
	examplehandlerd "ghgroups/example/example_mix/example_handler_d"
	examplehandlergroupe "ghgroups/example/example_mix/example_handler_group_e"
	examplelayerc "ghgroups/example/example_mix/example_layer_c"

	examplelayerg "ghgroups/example/example_mix/example_layer_g"
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
	factory.Register(reflect.TypeOf(examplelayera.ExampleADivider{}))

	factory.Register(reflect.TypeOf(examplelayerb.ExampleB1Handler{}))
	factory.Register(reflect.TypeOf(examplelayerb.ExampleB2Handler{}))
	factory.Register(reflect.TypeOf(examplelayerb.ExampleBDivider{}))

	factory.Register(reflect.TypeOf(examplelayerc.ExampleC1Handler{}))
	factory.Register(reflect.TypeOf(examplelayerc.ExampleC2Handler{}))
	factory.Register(reflect.TypeOf(examplelayerc.ExampleCDivider{}))

	factory.Register(reflect.TypeOf(examplehandlerd.ExampleDHandler{}))

	factory.Register(reflect.TypeOf(examplehandlergroupe.ExampleE1Handler{}))
	factory.Register(reflect.TypeOf(examplehandlergroupe.ExampleE2Handler{}))

	factory.Register(reflect.TypeOf(exampleasynchandlergroupf.ExampleF1Handler{}))
	factory.Register(reflect.TypeOf(exampleasynchandlergroupf.ExampleF2Handler{}))

	factory.Register(reflect.TypeOf(examplelayerg.ExampleGDivider{}))
	factory.Register(reflect.TypeOf(examplelayerg.ExampleG1Handler{}))
	factory.Register(reflect.TypeOf(examplelayerg.ExampleG2Handler{}))

	constructor := constructorbuilder.BuildConstructor(factory, concretePath)

	mainProcess := "layer_center_main"
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
