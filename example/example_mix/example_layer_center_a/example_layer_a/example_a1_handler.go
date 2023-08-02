package examplelayera

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleA1Handler struct {
	frame.HandlerBaseInterface
}

func NewExampleA1Handler() *ExampleA1Handler {
	return &ExampleA1Handler{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (e *ExampleA1Handler) Name() string {
	return reflect.TypeOf(*e).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// HandlerBaseInterface
func (e *ExampleA1Handler) Handle(*ghgroupscontext.GhGroupsContext) bool {
	fmt.Printf("run %s\n", e.Name())
	return true
}
