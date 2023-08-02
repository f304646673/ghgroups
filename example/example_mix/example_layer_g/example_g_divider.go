package examplelayerb

import (
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleGDivider struct {
	frame.DividerBaseInterface
}

func NewExampleGDivider() *ExampleGDivider {
	return &ExampleGDivider{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *ExampleGDivider) Name() string {
	return reflect.TypeOf(*s).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// DividerBaseInterface
func (s *ExampleGDivider) Select(context *ghgroupscontext.GhGroupsContext) string {
	return "ExampleG1Handler"
}
