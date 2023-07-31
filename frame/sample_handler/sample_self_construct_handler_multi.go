package samplehandler

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
)

// 这是相对简单的handler，它只用实现HandlerInterface两个接口
// 系统使用名称作为唯一检索键，通过构造不同的对象拥有不同的名字，可以在系统中有多个该名字的handler实例，即一个类型（struct)可以有多个该名字的handler实例

type SampleSelfConstructHandlerMulti struct {
	frame.HandlerBaseInterface
	name string
}

func NewSampleSelfConstructHandlerMulti(name string) *SampleSelfConstructHandlerMulti {
	return &SampleSelfConstructHandlerMulti{
		name: name,
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *SampleSelfConstructHandlerMulti) Name() string {
	return s.name
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// HandlerBaseInterface
func (s *SampleSelfConstructHandlerMulti) Handle(*ghgroupscontext.GhGroupsContext) bool {
	fmt.Sprintln(s.Name())
	return true
}

// ///////////////////////////////////////////////////////////////////////////////////////////
