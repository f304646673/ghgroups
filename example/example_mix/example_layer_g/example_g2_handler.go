package examplelayerb

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleG2Handler struct {
	frame.HandlerBaseInterface
}

func NewExampleG2Handler() *ExampleG2Handler {
	return &ExampleG2Handler{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (e *ExampleG2Handler) Name() string {
	return reflect.TypeOf(*e).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// HandlerBaseInterface
func (e *ExampleG2Handler) Handle(*ghgroupscontext.GhGroupsContext) bool {
	fmt.Printf("run %s\n", e.Name())
	return true
}
