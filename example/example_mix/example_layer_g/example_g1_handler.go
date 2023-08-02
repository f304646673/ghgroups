package examplelayerb

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleG1Handler struct {
	frame.HandlerBaseInterface
}

func NewExampleG1Handler() *ExampleG1Handler {
	return &ExampleG1Handler{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (e *ExampleG1Handler) Name() string {
	return reflect.TypeOf(*e).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// HandlerBaseInterface
func (e *ExampleG1Handler) Handle(*ghgroupscontext.GhGroupsContext) bool {
	fmt.Printf("run %s\n", e.Name())
	return true
}
