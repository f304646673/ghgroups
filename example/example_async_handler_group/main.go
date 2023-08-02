package main

import (
	"fmt"
	"ghgroups/frame"
	"ghgroups/frame/constructor"
	constructorbuilder "ghgroups/frame/constructor_builder"
	"ghgroups/frame/factory"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"os"
	"path"
	"reflect"
)

func main() {
	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(ExampleAHandler{}))
	factory.Register(reflect.TypeOf(ExampleBHandler{}))

	runPath, errGetWd := os.Getwd()
	if errGetWd != nil {
		fmt.Printf("%v", errGetWd)
		return
	}
	concretePath := path.Join(runPath, "conf")
	constructor := constructorbuilder.BuildConstructor(factory, concretePath)
	mainProcess := "async_handler_group_a"

	run(constructor, mainProcess)
}

func run(constructor *constructor.Constructor, mainProcess string) {
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
			// context.ShowDuration = true
			mainHandlerGroup.Handle(context)
		}
	}
}
