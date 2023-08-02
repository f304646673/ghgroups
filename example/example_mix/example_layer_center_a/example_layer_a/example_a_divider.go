package examplelayera

import (
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleADivider struct {
	frame.DividerBaseInterface
}

func NewExampleADivider() *ExampleADivider {
	return &ExampleADivider{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *ExampleADivider) Name() string {
	return reflect.TypeOf(*s).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// DividerBaseInterface
func (s *ExampleADivider) Select(context *ghgroupscontext.GhGroupsContext) string {
	return "ExampleA2Handler"
}
