package samplehandler

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"reflect"
)

// 这是最简单的handler，它只用实现HandlerInterface接口
// 因为Name方法限定了名称，且系统使用名称作为唯一检索键，这样整个系统里只能有一个该名字的handler实例，即一个类型（struct)只能有一个该名字的handler实例
// 如果希望一个struct在系统中有多个实例，那么可以使用SampleSelfConstructHandlerMulti

type SampleSelfConstructHandlerSingle struct {
	frame.HandlerBaseInterface
}

func NewSampleSelfConstructHandlerSingle() *SampleSelfConstructHandlerSingle {
	return &SampleSelfConstructHandlerSingle{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *SampleSelfConstructHandlerSingle) Name() string {
	return reflect.TypeOf(*s).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// HandlerBaseInterface
func (s *SampleSelfConstructHandlerSingle) Handle(*ghgroupscontext.GhGroupsContext) bool {
	fmt.Sprintln(s.Name())
	return true
}

// ///////////////////////////////////////////////////////////////////////////////////////////
