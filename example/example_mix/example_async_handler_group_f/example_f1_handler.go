package exampleasynchandlergroupf

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
	"time"
)

type ExampleF1Handler struct {
	frame.HandlerBaseInterface
}

func NewExampleF1Handler() *ExampleF1Handler {
	return &ExampleF1Handler{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (e *ExampleF1Handler) Name() string {
	return reflect.TypeOf(*e).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// HandlerBaseInterface
func (e *ExampleF1Handler) Handle(*ghgroupscontext.GhGroupsContext) bool {
	fmt.Printf("run %s\n", e.Name())
	time.Sleep(1 * time.Second)
	return true
}
