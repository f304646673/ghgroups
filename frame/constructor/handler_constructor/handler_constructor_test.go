package handlerconstructor

import (
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"ghgroups/frame/constructor"
	"ghgroups/frame/factory"

	aynchandlergroupconstructor "ghgroups/frame/constructor/async_handler_group_constructor"
	dividerconstructor "ghgroups/frame/constructor/divider_constructor"
	handlergroupconstructor "ghgroups/frame/constructor/handler_group_constructor"
	layercenterconstructor "ghgroups/frame/constructor/layer_center_constructor"
	layerconstructor "ghgroups/frame/constructor/layer_constructor"
	samplehandler "ghgroups/frame/sample_handler"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestGetHandlerName(t *testing.T) {
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
		h := NewHandlerConstructor(nil)
		handlerName, err := h.getHandlerName(confPath)
		assert.Empty(t, handlerName)
		assert.NotNil(t, err)
	})

	t.Run("Input=", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		h := NewHandlerConstructor(nil)
		handlerName, err := h.getHandlerName(confPath)
		assert.Empty(t, handlerName)
		assert.ErrorContains(t, err, "is not a directory")
	})

	t.Run("Input=valid/sample_a.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		h := NewHandlerConstructor(nil)
		handlerName, err := h.getHandlerName(confPath)
		assert.Equal(t, handlerName, "sample_a")
		assert.Nil(t, err)
	})
}

func TestRegisterHandler(t *testing.T) {
	t.Run("Input=valid/sample_a.yaml", func(t *testing.T) {
		t.Parallel()
		h := NewHandlerConstructor(nil)
		handlerName := "handler_name"
		exist := h.existHandler(handlerName)
		assert.False(t, exist)
		err := h.RegisterHandler(handlerName, nil)
		assert.Nil(t, err)
		exist = h.existHandler(handlerName)
		assert.True(t, exist)
		err = h.RegisterHandler(handlerName, nil)
		assert.NotNil(t, err)
	})
}

func TestInit(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(layerconstructor.LayerConstructor{}))
	factory.Register(reflect.TypeOf(dividerconstructor.DividerConstructor{}))
	factory.Register(reflect.TypeOf(HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	constructor := constructor.NewConstructor(factory, "")

	err := constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))
	assert.Nil(t, err)

	t.Run("Input=valid", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		handlerConstructor := NewHandlerConstructor(constructor)
		err := handlerConstructor.ParseHandlerConfFolder(confPath)
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

		handlerConstructor := NewHandlerConstructor(constructor)
		err := handlerConstructor.ParseHandlerConfFolder(confPath)
		assert.ErrorContains(t, err, "no such file or directory")
	})

	t.Run("Input=invalid/handler_name_not_match", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		handlerConstructor := NewHandlerConstructor(constructor)
		err := handlerConstructor.ParseHandlerConfFolder(confPath)
		assert.ErrorContains(t, err, "mismatch with configuration file Name")
	})

	t.Run("Input=invalid/same_handler_name", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		handlerConstructor := NewHandlerConstructor(constructor)
		err := handlerConstructor.ParseHandlerConfFolder(confPath)
		assert.ErrorContains(t, err, "handler name sample_e already exists,path is")
	})

	t.Run("Input=invalid/same_handler_name_with_error_path", func(t *testing.T) {
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

		handlerConstructor := NewHandlerConstructor(constructor)
		handlerConstructor.ParseHandlerConfFolder(confPath)
		assert.FailNow(t, "must panic when handler name is the same")
	})

	t.Run("Input=invalid/no_type", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		handlerConstructor := NewHandlerConstructor(constructor)
		err := handlerConstructor.ParseHandlerConfFolder(confPath)
		assert.ErrorContains(t, err, "handler configuration must have type")
	})
}

func TestConstructHandlerFromFile(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(layerconstructor.LayerConstructor{}))
	factory.Register(reflect.TypeOf(dividerconstructor.DividerConstructor{}))
	factory.Register(reflect.TypeOf(HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	constructor := constructor.NewConstructor(factory, "")

	t.Run("Input=InvalidFilePath", func(t *testing.T) {
		handlerConstructor := NewHandlerConstructor(constructor)
		handlerInterface, err := handlerConstructor.constructHandlerFromFile("")
		assert.Nil(t, handlerInterface)
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

		handlerConstructor := NewHandlerConstructor(constructor)
		handlerInterface, err := handlerConstructor.constructHandlerFromFile(confFilePath)
		assert.Nil(t, handlerInterface)
		assert.ErrorContains(t, err, "cannot unmarshal")
	})
}

func TestConstructHandlerFromName(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(layerconstructor.LayerConstructor{}))
	factory.Register(reflect.TypeOf(dividerconstructor.DividerConstructor{}))
	factory.Register(reflect.TypeOf(HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	constructor := constructor.NewConstructor(factory, "")

	err := constructor.Register(reflect.TypeOf(SampleWithoutHandleHandler{}))
	assert.Nil(t, err)

	t.Run("Input=InvalidFilePath", func(t *testing.T) {
		handlerConstructor := NewHandlerConstructor(constructor)
		handlerInterface, err := handlerConstructor.constructHandlerFromName("")
		assert.Nil(t, handlerInterface)
		assert.ErrorContains(t, err, "handler's configuration is not set")
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

		handlerConstructor := NewHandlerConstructor(constructor)
		handlerInterface, err := handlerConstructor.constructHandlerFromFile(confFilePath)
		assert.Nil(t, handlerInterface)
		assert.ErrorContains(t, err, "cannot unmarshal")
	})

	t.Run("Input=invalid/handler_without_handle", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)
		confFilePath := path.Join(confPath, "sample_j.yaml")

		handlerConstructor := NewHandlerConstructor(constructor)
		handlerInterface, err := handlerConstructor.constructHandlerFromFile(confFilePath)
		assert.Nil(t, handlerInterface)
		assert.ErrorContains(t, err, "convert SampleWithoutHandleHandler to HandlerBaseInterface error")
	})
}

type SampleWithoutHandleHandlerConf struct {
	Name string `yaml:"name"`
}

type SampleWithoutHandleHandler struct {
	// frame.HandlerBaseInterface
	conf SampleWithoutHandleHandlerConf
}

func NewSampleWithoutHandleHandler() *SampleWithoutHandleHandler {
	return &SampleWithoutHandleHandler{}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// DynamicLoadingInterface
func (s *SampleWithoutHandleHandler) LoadConfigFromFile(confPath string) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	return s.LoadConfigFromMemory(data)
}

func (s *SampleWithoutHandleHandler) LoadConfigFromMemory(configure []byte) error {
	sampleHandlerConf := new(SampleWithoutHandleHandlerConf)
	err := yaml.Unmarshal([]byte(configure), sampleHandlerConf)
	if err != nil {
		return err
	}
	s.conf = *sampleHandlerConf
	return nil
}

func (s *SampleWithoutHandleHandler) LoadEnvironmentConf() error {
	return nil
}

func (s *SampleWithoutHandleHandler) Name() string {
	return s.conf.Name
}
