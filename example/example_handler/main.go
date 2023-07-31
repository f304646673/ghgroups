package main

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"ghgroups/frame/utils"
	"reflect"
)

func main() {
	constructor := utils.BuildConstructor("")
	constructor.Register(reflect.TypeOf(ExampleHandler{}))

	context := ghgroupscontext.NewGhGroupsContext(nil)
	mainProcess := reflect.TypeOf(ExampleHandler{}).Name()
	if err := constructor.CreateConcrete(mainProcess); err != nil {
		fmt.Printf("%v", err)
	}
	if someInterfaced, err := constructor.GetConcrete(mainProcess); err != nil {
		fmt.Printf("%v", err)
	} else {
		if mainHandlerGroup, ok := someInterfaced.(frame.HandlerBaseInterface); !ok {
			fmt.Printf("mainHandlerGroup %s is not frame.HandlerBaseInterface", mainProcess)
		} else {
			mainHandlerGroup.Handle(context)
		}
	}
}
