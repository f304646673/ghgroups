package examplelayera

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

type ExampleA2Handler struct {
	frame.HandlerBaseInterface
}

func NewExampleA2Handler() *ExampleA2Handler {
	return &ExampleA2Handler{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (e *ExampleA2Handler) Name() string {
	return reflect.TypeOf(*e).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// HandlerBaseInterface
func (e *ExampleA2Handler) Handle(*ghgroupscontext.GhGroupsContext) bool {
	fmt.Printf("run %s\n", e.Name())
	return true
}
