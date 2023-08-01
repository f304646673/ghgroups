package main

import (
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleDivider struct {
	frame.DividerBaseInterface
}

func NewExampleDivider() *ExampleDivider {
	return &ExampleDivider{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *ExampleDivider) Name() string {
	return reflect.TypeOf(*s).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// DividerBaseInterface
func (s *ExampleDivider) Select(context *ghgroupscontext.GhGroupsContext) string {
	return "ExampleBHandler"
}
