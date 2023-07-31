package dividerconstructor

import (
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"ghgroups/frame/constructor"
	aynchandlergroupconstructor "ghgroups/frame/constructor/async_handler_group_constructor"
	handlerconstructor "ghgroups/frame/constructor/handler_constructor"
	handlergroupconstructor "ghgroups/frame/constructor/handler_group_constructor"
	layercenterconstructor "ghgroups/frame/constructor/layer_center_constructor"
	layerconstructor "ghgroups/frame/constructor/layer_constructor"
	"ghgroups/frame/factory"
	sampledivider "ghgroups/frame/sample_divider"

	ghgroupscontext "ghgroups/frame/ghgroups_context"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestGetDividerName(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	t.Run("Input=.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		h := NewDividerConstructor(nil)
		dividerName, err := h.getDividerName(confPath)
		assert.Empty(t, dividerName)
		assert.NotNil(t, err)
	})

	t.Run("Input=", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		h := NewDividerConstructor(nil)
		dividerName, err := h.getDividerName(confPath)
		assert.Empty(t, dividerName)
		assert.ErrorContains(t, err, "is not a directory")
	})

	t.Run("Input=valid/sample_a.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		h := NewDividerConstructor(nil)
		dividerName, err := h.getDividerName(confPath)
		assert.Equal(t, dividerName, "sample_a")
		assert.Nil(t, err)
	})
}

func TestRegisterDivider(t *testing.T) {
	t.Run("Input=valid/sample_a.yaml", func(t *testing.T) {
		t.Parallel()
		h := NewDividerConstructor(nil)
		dividerName := "divider_name"
		exist := h.existDivider(dividerName)
		assert.False(t, exist)
		err := h.RegisterDivider(dividerName, nil)
		assert.Nil(t, err)
		exist = h.existDivider(dividerName)
		assert.True(t, exist)
		err = h.RegisterDivider(dividerName, nil)
		assert.NotNil(t, err)
	})
}

func TestInit(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(layerconstructor.LayerConstructor{}))
	factory.Register(reflect.TypeOf(DividerConstructor{}))
	factory.Register(reflect.TypeOf(handlerconstructor.HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	constructor := constructor.NewConstructor(factory, "")
	constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))

	t.Run("Input=valid", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		dividerConstructor := NewDividerConstructor(constructor)
		err := dividerConstructor.ParseDividerConfFolder(confPath)
		assert.Nil(t, err)
	})

	t.Run("Input=invalid/no_such_path", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		dividerConstructor := NewDividerConstructor(constructor)
		err := dividerConstructor.ParseDividerConfFolder(confPath)
		assert.ErrorContains(t, err, "no such file or directory")
	})

	t.Run("Input=invalid/divider_name_not_match", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		dividerConstructor := NewDividerConstructor(constructor)
		err := dividerConstructor.ParseDividerConfFolder(confPath)
		assert.ErrorContains(t, err, "mismatch with configuration file Name")
	})

	t.Run("Input=invalid/same_divider_name", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		dividerConstructor := NewDividerConstructor(constructor)
		err := dividerConstructor.ParseDividerConfFolder(confPath)
		assert.ErrorContains(t, err, "divider name sample_e already exists,path is")
	})

	t.Run("Input=invalid/same_divider_name_with_error_path", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		defer func() {
			r := recover()
			assert.Contains(t, r, "concrete's Name  sample_f has already been used")
		}()

		dividerConstructor := NewDividerConstructor(constructor)
		dividerConstructor.ParseDividerConfFolder(confPath)
		assert.FailNow(t, "must panic when divider name is the same")
	})

	t.Run("Input=invalid/no_type", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		dividerConstructor := NewDividerConstructor(constructor)
		err := dividerConstructor.ParseDividerConfFolder(confPath)
		assert.ErrorContains(t, err, "divider configuration must have type")
	})
}

func TestConstructDividerFromFile(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(layerconstructor.LayerConstructor{}))
	factory.Register(reflect.TypeOf(DividerConstructor{}))
	factory.Register(reflect.TypeOf(handlerconstructor.HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	constructor := constructor.NewConstructor(factory, "")
	constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))

	t.Run("Input=InvalidFilePath", func(t *testing.T) {
		dividerConstructor := NewDividerConstructor(constructor)
		dividerInterface, err := dividerConstructor.constructDividerFromFile("")
		assert.Nil(t, dividerInterface)
		assert.ErrorContains(t, err, "no such file or directory")
	})

	t.Run("Input=invalid/error_unmarshal", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)
		confFilePath := path.Join(confPath, "sample_h.yaml")

		dividerConstructor := NewDividerConstructor(constructor)
		dividerConstructor.ParseDividerConfFolder(confPath)
		dividerInterface, err := dividerConstructor.constructDividerFromFile(confFilePath)
		assert.Nil(t, dividerInterface)
		assert.ErrorContains(t, err, "cannot unmarshal")
	})
}

func TestConstructDividerFromName(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(layerconstructor.LayerConstructor{}))
	factory.Register(reflect.TypeOf(DividerConstructor{}))
	factory.Register(reflect.TypeOf(handlerconstructor.HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	constructor := constructor.NewConstructor(factory, "")

	constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))
	constructor.Register(reflect.TypeOf(SampleWithoutSelectDivider{}))

	t.Run("Input=InvalidFilePath", func(t *testing.T) {
		dividerConstructor := NewDividerConstructor(constructor)
		dividerInterface, err := dividerConstructor.constructDividerFromName("")
		assert.Nil(t, dividerInterface)
		assert.ErrorContains(t, err, "divider's configuration is not set")
	})

	t.Run("Input=invalid/error_unmarshal", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)
		confFilePath := path.Join(confPath, "sample_h.yaml")

		dividerConstructor := NewDividerConstructor(constructor)
		dividerConstructor.ParseDividerConfFolder(confPath)
		dividerInterface, err := dividerConstructor.constructDividerFromFile(confFilePath)
		assert.Nil(t, dividerInterface)
		assert.ErrorContains(t, err, "cannot unmarshal")
	})

	t.Run("Input=invalid/divider_without_handle", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		dividerConstructor := NewDividerConstructor(constructor)
		dividerConstructor.LoadEnvironmentConf()
		err := dividerConstructor.ParseDividerConfFolder(confPath)
		assert.Nil(t, err)
	})
}

type SampleWithoutSelectDividerConf struct {
	Name string `yaml:"name"`
}

type SampleWithoutSelectDivider struct {
	// frame.DividerBaseInterface
	conf SampleWithoutSelectDividerConf
}

func NewSampleWithoutSelectDivider() *SampleWithoutSelectDivider {
	return &SampleWithoutSelectDivider{}
}

func (s *SampleWithoutSelectDivider) Select(context *ghgroupscontext.GhGroupsContext) string {
	return ""
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// DynamicLoadingInterface
func (s *SampleWithoutSelectDivider) LoadConfigFromFile(confPath string) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	return s.LoadConfigFromMemory(data)
}

func (s *SampleWithoutSelectDivider) LoadConfigFromMemory(configure []byte) error {
	sampleDividerConf := new(SampleWithoutSelectDividerConf)
	err := yaml.Unmarshal([]byte(configure), sampleDividerConf)
	if err != nil {
		return err
	}
	s.conf = *sampleDividerConf
	return nil
}

func (s *SampleWithoutSelectDivider) LoadEnvironmentConf() error {
	return nil
}

func (s *SampleWithoutSelectDivider) Name() string {
	return s.conf.Name
}
