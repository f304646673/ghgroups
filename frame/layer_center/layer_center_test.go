package layercenter

import (
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"ghgroups/frame/layer"
	"ghgroups/frame/utils"
	"reflect"

	sampledivider "ghgroups/frame/sample_divider"
	samplehandler "ghgroups/frame/sample_handler"

	"os"
	"path"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	constructor := utils.BuildConstructor(testDataPath)
	constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))
	constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))
	constructor.Register(reflect.TypeOf(layer.Layer{}))

	layerCenter := NewLayerCenter(constructor)

	t.Run("Input=not_exist.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)
		err := layerCenter.LoadConfigFromFile(confPath)
		assert.ErrorContains(t, err, "no such file or directory")
	})

	t.Run("Input=valid.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)
		err := layerCenter.LoadConfigFromFile(confPath)
		assert.Nil(t, err)
	})

	t.Run("Input=invalid.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)
		err := layerCenter.LoadConfigFromFile(confPath)
		assert.ErrorContains(t, err, "yaml: unmarshal errors")
	})
}

func TestInit(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	constructor := utils.BuildConstructor(testDataPath)

	constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))
	constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))
	constructor.Register(reflect.TypeOf(layer.Layer{}))

	handlersConfPath := path.Join(testDataPath, "handlers")
	err := constructor.ParseHandlerConfFolder(handlersConfPath)
	assert.Nil(t, err)

	dividerConfPath := path.Join(testDataPath, "divider")
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

		layerCenter := NewLayerCenter(constructor)

		err = layerCenter.LoadConfigFromFile(confPath)
		assert.Nil(t, err)
	})
}

func TestAdd(t *testing.T) {
	constructor := utils.BuildConstructor("")

	layerCenter := NewLayerCenter(constructor)
	testLayer := layer.NewLayer("test_layer", constructor)
	layerCenter.Add(testLayer)

	called := false
	monkey.PatchInstanceMethod(reflect.TypeOf(testLayer), "Handle", func(*layer.Layer, *ghgroupscontext.GhGroupsContext) bool {
		called = true
		return true
	})

	ctx := ghgroupscontext.NewGhGroupsContext(nil)
	assert.True(t, layerCenter.Handle(ctx))
	assert.True(t, called)
}
