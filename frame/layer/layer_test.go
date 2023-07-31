package layer

import (
	"fmt"
	dividerconstructor "ghgroups/frame/constructor/divider_constructor"
	handlerconstructor "ghgroups/frame/constructor/handler_constructor"
	layerconstructor "ghgroups/frame/constructor/layer_constructor"
	"ghgroups/frame/factory"
	"ghgroups/frame/utils"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	ghgroupscontext "ghgroups/frame/ghgroups_context"
	sampledivider "ghgroups/frame/sample_divider"
	samplehandler "ghgroups/frame/sample_handler"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(layerconstructor.LayerConstructor{}))
	factory.Register(reflect.TypeOf(dividerconstructor.DividerConstructor{}))
	factory.Register(reflect.TypeOf(handlerconstructor.HandlerConstructor{}))
	constructor := utils.BuildConstructor("")

	t.Run("Input=not_exist.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)
		layer := NewLayer("", constructor)
		err := layer.LoadConfigFromFile(confPath)
		assert.ErrorContains(t, err, "no such file or directory")
	})

	t.Run("Input=sample_error.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)
		layer := NewLayer("", constructor)
		err := layer.LoadConfigFromFile(confPath)
		assert.ErrorContains(t, err, "yaml: unmarshal errors")
	})

	t.Run("Input=valid.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)

		constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))
		handlersConfPath := path.Join(testDataPath, "handlers")
		err := constructor.ParseHandlerConfFolder(handlersConfPath)
		assert.Nil(t, err)

		constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))
		dividerConfPath := path.Join(testDataPath, "divider")
		err = constructor.ParseDividerConfFolder(dividerConfPath)
		assert.Nil(t, err)

		layer := NewLayer("test_layer", constructor)
		err = layer.LoadConfigFromFile(confPath)
		assert.Nil(t, err)
	})
}

func TestAutoConstruct(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(layerconstructor.LayerConstructor{}))
	factory.Register(reflect.TypeOf(dividerconstructor.DividerConstructor{}))
	factory.Register(reflect.TypeOf(handlerconstructor.HandlerConstructor{}))
	constructor := utils.BuildConstructor("")

	constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))
	constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))

	handlersConfPath := path.Join(testDataPath, "handlers")
	dividerConfPath := path.Join(testDataPath, "divider")
	err := constructor.ParseHandlerConfFolder(handlersConfPath)
	assert.Nil(t, err)
	err = constructor.ParseDividerConfFolder(dividerConfPath)
	assert.Nil(t, err)

	t.Run("Input=valid.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)

		layer := NewLayer("", constructor)
		layer.SetConstructorInterface(constructor)
		err := layer.LoadConfigFromFile(confPath)
		assert.Nil(t, err)

		var output string
		patch := monkey.Patch(fmt.Sprintln, func(a ...interface{}) string {
			s := make([]interface{}, len(a))
			for i, v := range a {
				s[i] = fmt.Sprint(v)
				output = fmt.Sprintf("%s %s", output, s[i])
			}
			return output
		})
		defer patch.Unpatch()

		var context ghgroupscontext.GhGroupsContext
		suc := layer.Handle(&context)
		assert.True(t, suc)
		assert.Equal(t, output, " handler_sample_b")
	})
}

func TestSelfConstruct(t *testing.T) {
	constructor := utils.BuildConstructor("")

	// 不是自动构建，不需要注册
	// constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))
	// constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))

	sampleSelfConstructHandlerSingle := samplehandler.NewSampleSelfConstructHandlerSingle()
	sampleSelfConstructDividerSingle := sampledivider.NewSampleSelfConstructDividerSingle()

	// 注册手动构建的handler和divider
	// constructor.RegisterHandler(sampleSelfConstructHandlerSingle.Name(), sampleSelfConstructHandlerSingle)
	// constructor.RegisterDivider(sampleSelfConstructDividerSingle.Name(), sampleSelfConstructDividerSingle)

	// 手工构建layer
	testLayer := NewLayer("test_layer", constructor)
	// 手动设置layer的divider和handler
	testLayer.SetDivider(sampleSelfConstructDividerSingle.Name(), sampleSelfConstructDividerSingle)
	testLayer.AddHandler(sampleSelfConstructHandlerSingle.Name(), sampleSelfConstructHandlerSingle)

	patch := monkey.Patch(fmt.Sprintln, func(a ...interface{}) string {
		s := make([]interface{}, len(a))
		output := ""
		for i, v := range a {
			s[i] = fmt.Sprint(v)
			output = fmt.Sprintf("%s", s[i])
		}
		assert.Equal(t, "SampleSelfConstructHandlerSingle", output)
		return output
	})
	defer patch.Unpatch()

	var context ghgroupscontext.GhGroupsContext
	suc := testLayer.Handle(&context)
	assert.True(t, suc)
}

func TestSelfAutoConstruct(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	constructor := utils.BuildConstructor("")

	// 不是自动构建，不需要注册
	// constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))
	// constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))

	// 注册手动构建的handler和divider
	sampleSelfConstructHandlerSingle := samplehandler.NewSampleSelfConstructHandlerSingle()
	constructor.RegisterHandler(sampleSelfConstructHandlerSingle.Name(), sampleSelfConstructHandlerSingle)

	sampleSelfConstructDividerSingle := sampledivider.NewSampleSelfConstructDividerSingle()
	constructor.RegisterDivider(sampleSelfConstructDividerSingle.Name(), sampleSelfConstructDividerSingle)

	// 自动构建Layer，就需要注册
	constructor.Register(reflect.TypeOf(Layer{}))
	constructor.ParseLayerConfFolder(path.Join(testDataPath, "layer"))
	layerBaseInterface, err := constructor.GetLayer("layer_sample_a")
	assert.Nil(t, err)

	patch := monkey.Patch(fmt.Sprintln, func(a ...interface{}) string {
		s := make([]interface{}, len(a))
		output := ""
		for i, v := range a {
			s[i] = fmt.Sprint(v)
			output = fmt.Sprintf("%s", s[i])
		}
		assert.Equal(t, "SampleSelfConstructHandlerSingle", output)
		return output
	})
	defer patch.Unpatch()

	var context ghgroupscontext.GhGroupsContext
	suc := layerBaseInterface.Handle(&context)
	assert.True(t, suc)
}
