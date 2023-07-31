package sampledivider

import (
	"ghgroups/frame"

	ghgroupscontext "ghgroups/frame/ghgroups_context"
)

// 这是相对简单的divider，它只用实现DividerInterface两个接口
// 系统使用名称作为唯一检索键，通过构造不同的对象拥有不同的名字，可以在系统中有多个该名字的divider实例，即一个类型（struct)可以有多个该名字的divider实例

type SampleSelfConstructDividerMulti struct {
	frame.DividerBaseInterface
	name string
}

func NewSampleSelfConstructDividerMulti(name string) *SampleSelfConstructDividerMulti {
	return &SampleSelfConstructDividerMulti{
		name: name,
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *SampleSelfConstructDividerMulti) Name() string {
	return s.name
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// DividerBaseInterface
func (s *SampleSelfConstructDividerMulti) Select(context *ghgroupscontext.GhGroupsContext) string {
	return "SampleSelfConstructHandlerSingle"
}

// ///////////////////////////////////////////////////////////////////////////////////////////
