package sampledivider

import (
	"ghgroups/frame"

	ghgroupscontext "ghgroups/frame/ghgroups_context"

	"gopkg.in/yaml.v2"
)

// 自动构建divider，它会自动从配置文件中读取配置，然后根据配置构建divider
// 因为系统使用名称作为唯一检索键，所以自动构建divider在构建过程中，就要被命名，而名称应该来源于配置文件
// 这就要求配置文件中必须有一个名为name的字段，用于指定divider的名称
// 下面例子中confs配置不是必须的，divider的实现者，需要自行解析配置文件，以确保Name方法返回的名称与配置文件中的name字段一致

type SampleAutoConstructDividerConf struct {
	Name   string                              `yaml:"name"`
	Select string                              `yaml:"select"`
	Confs  []SampleAutoConstructDividerEnvConf `yaml:"confs"`
}

type SampleAutoConstructDividerEnvConf struct {
	Env         string                                 `yaml:"env"`
	RegionsConf []SampleAutoConstructDividerRegionConf `yaml:"regions_conf"`
}

type SampleAutoConstructDividerRegionConf struct {
	Region          string `yaml:"region"`
	AwsRegion       string `yaml:"aws_region"`
	AccessKeyId     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	IntKey          int32  `yaml:"int_key"`
}

type SampleAutoConstructDivider struct {
	frame.DividerBaseInterface
	frame.LoadConfigFromMemoryInterface
	conf SampleAutoConstructDividerConf
}

func NewSampleAutoConstructDivider() *SampleAutoConstructDivider {
	return &SampleAutoConstructDivider{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// LoadConfigFromMemoryInterface
func (s *SampleAutoConstructDivider) LoadConfigFromMemory(configure []byte) error {
	sampleDividerConf := new(SampleAutoConstructDividerConf)
	err := yaml.Unmarshal([]byte(configure), sampleDividerConf)
	if err != nil {
		return err
	}
	s.conf = *sampleDividerConf
	return nil
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ConcreteInterface
func (s *SampleAutoConstructDivider) Name() string {
	return s.conf.Name
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// DividerBaseInterface
func (s *SampleAutoConstructDivider) Select(context *ghgroupscontext.GhGroupsContext) string {
	return s.conf.Select
}

// ///////////////////////////////////////////////////////////////////////////////////////////
