package examplelayerc

import (
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleCDivider struct {
	frame.DividerBaseInterface
}

func NewExampleCDivider() *ExampleCDivider {
	return &ExampleCDivider{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *ExampleCDivider) Name() string {
	return reflect.TypeOf(*s).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// DividerBaseInterface
func (s *ExampleCDivider) Select(context *ghgroupscontext.GhGroupsContext) string {
	return "ExampleC2Handler"
}
