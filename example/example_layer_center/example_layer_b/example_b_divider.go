package examplelayerb

import (
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleBDivider struct {
	frame.DividerBaseInterface
}

func NewExampleBDivider() *ExampleBDivider {
	return &ExampleBDivider{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *ExampleBDivider) Name() string {
	return reflect.TypeOf(*s).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// DividerBaseInterface
func (s *ExampleBDivider) Select(context *ghgroupscontext.GhGroupsContext) string {
	return "ExampleB1Handler"
}
