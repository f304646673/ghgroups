package main

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleBHandler struct {
	frame.HandlerBaseInterface
}

func NewExampleBHandler() *ExampleBHandler {
	return &ExampleBHandler{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (e *ExampleBHandler) Name() string {
	return reflect.TypeOf(*e).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// HandlerBaseInterface
func (e *ExampleBHandler) Handle(*ghgroupscontext.GhGroupsContext) bool {
	fmt.Printf("run %s\n", e.Name())
	return true
}
