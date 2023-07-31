package samplehandler

import (
	"fmt"
	"ghgroups/frame"
	ghgroupscontext "ghgroups/frame/ghgroups_context"

	"gopkg.in/yaml.v2"
)

// 自动构建handler，它会自动从配置文件中读取配置，然后根据配置构建handler
// 因为系统使用名称作为唯一检索键，所以自动构建handler在构建过程中，就要被命名，而名称应该来源于配置文件
// 这就要求配置文件中必须有一个名为name的字段，用于指定handler的名称
// 下面例子中confs配置不是必须的，handler的实现者，需要自行解析配置文件，以确保Name方法返回的名称与配置文件中的name字段一致

type SampleAutoConstructHandlerConf struct {
	Name  string                              `yaml:"name"`
	Confs []SampleAutoConstructHandlerEnvConf `yaml:"confs"`
}

type SampleAutoConstructHandlerEnvConf struct {
	Env         string                                 `yaml:"env"`
	RegionsConf []SampleAutoConstructHandlerRegionConf `yaml:"regions_conf"`
}

type SampleAutoConstructHandlerRegionConf struct {
	Region          string `yaml:"region"`
	AwsRegion       string `yaml:"aws_region"`
	AccessKeyId     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	IntKey          int32  `yaml:"int_key"`
}

type SampleAutoConstructHandler struct {
	frame.HandlerBaseInterface
	frame.LoadConfigFromMemoryInterface
	conf SampleAutoConstructHandlerConf
}

func NewSampleAutoConstructHandler() *SampleAutoConstructHandler {
	return &SampleAutoConstructHandler{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// LoadConfigFromMemoryInterface
func (s *SampleAutoConstructHandler) LoadConfigFromMemory(configure []byte) error {
	sampleHandlerConf := new(SampleAutoConstructHandlerConf)
	err := yaml.Unmarshal([]byte(configure), sampleHandlerConf)
	if err != nil {
		return err
	}
	s.conf = *sampleHandlerConf
	return nil
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *SampleAutoConstructHandler) Name() string {
	return s.conf.Name
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// HandlerBaseInterface
func (s *SampleAutoConstructHandler) Handle(*ghgroupscontext.GhGroupsContext) bool {
	fmt.Sprintln(s.conf.Name)
	return true
}

// ///////////////////////////////////////////////////////////////////////////////////////////
