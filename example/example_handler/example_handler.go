package main

import (
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleHandler struct {
	frame.HandlerBaseInterface
}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (e *ExampleHandler) Name() string {
	return reflect.TypeOf(*e).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// HandlerBaseInterface
func (e *ExampleHandler) Handle(*ghgroupscontext.GhGroupsContext) bool {
	return true
}
