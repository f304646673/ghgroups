package main

import (
	"fmt"
	"ghgroups/frame"
	"ghgroups/frame/constructor"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"ghgroups/frame/utils"
	"reflect"
)

func main() {
	constructor := utils.BuildConstructor("")
	constructor.Register(reflect.TypeOf(ExampleHandler{}))
	mainProcess := reflect.TypeOf(ExampleHandler{}).Name()
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
