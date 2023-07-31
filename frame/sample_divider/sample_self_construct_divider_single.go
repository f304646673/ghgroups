package sampledivider

import (
	"ghgroups/frame"
	"reflect"

	ghgroupscontext "ghgroups/frame/ghgroups_context"
)

// 这是最简单的divider，它只用实现DividerInterface接口
// 因为Name方法限定了名称，且系统使用名称作为唯一检索键，这样整个系统里只能有一个该名字的divider实例，即一个类型（struct)只能有一个该名字的divider实例
// 如果希望一个struct在系统中有多个实例，那么可以使用SampleSelfConstructDividerMulti

type SampleSelfConstructDividerSingle struct {
	frame.DividerBaseInterface
}

func NewSampleSelfConstructDividerSingle() *SampleSelfConstructDividerSingle {
	return &SampleSelfConstructDividerSingle{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *SampleSelfConstructDividerSingle) Name() string {
	return reflect.TypeOf(*s).Name()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// DividerBaseInterface
func (s *SampleSelfConstructDividerSingle) Select(context *ghgroupscontext.GhGroupsContext) string {
	return "SampleSelfConstructHandlerSingle"
}

// ///////////////////////////////////////////////////////////////////////////////////////////
