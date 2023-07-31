package layerconstructor

import (
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"ghgroups/frame/constructor"
	dividerconstructor "ghgroups/frame/constructor/divider_constructor"
	handlerconstructor "ghgroups/frame/constructor/handler_constructor"
	handlergroupconstructor "ghgroups/frame/constructor/handler_group_constructor"
	layercenterconstructor "ghgroups/frame/constructor/layer_center_constructor"
	"ghgroups/frame/factory"
	"ghgroups/frame/layer"

	aynchandlergroupconstructor "ghgroups/frame/constructor/async_handler_group_constructor"

	sampledivider "ghgroups/frame/sample_divider"

	samplehandler "ghgroups/frame/sample_handler"

	"github.com/stretchr/testify/assert"
)

func TestGetLayerName(t *testing.T) {
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
		h := NewLayerConstructor(nil)
		layerName, err := h.getLayerName(confPath)
		assert.Empty(t, layerName)
		assert.NotNil(t, err)
	})

	t.Run("Input=", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		h := NewLayerConstructor(nil)
		layerName, err := h.getLayerName(confPath)
		assert.Empty(t, layerName)
		assert.ErrorContains(t, err, "is not a directory")
	})

	t.Run("Input=valid/layer_sample_a.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		h := NewLayerConstructor(nil)
		layerName, err := h.getLayerName(confPath)
		assert.Equal(t, layerName, "layer_sample_a")
		assert.Nil(t, err)
	})
}

func TestRegisterLayer(t *testing.T) {
	t.Run("Input=valid/sample_a.yaml", func(t *testing.T) {
		t.Parallel()
		layerConstructor := NewLayerConstructor(nil)
		layerName := "layer_name"
		exist := layerConstructor.existLayer(layerName)
		assert.False(t, exist)
		err := layerConstructor.RegisterLayer(layerName, nil)
		assert.Nil(t, err)
		exist = layerConstructor.existLayer(layerName)
		assert.True(t, exist)
		err = layerConstructor.RegisterLayer(layerName, nil)
		assert.NotNil(t, err)
	})
}

func TestInit(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(LayerConstructor{}))
	factory.Register(reflect.TypeOf(dividerconstructor.DividerConstructor{}))
	factory.Register(reflect.TypeOf(handlerconstructor.HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	constructor := constructor.NewConstructor(factory, "")
	constructor.Register(reflect.TypeOf(layer.Layer{}))

	constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))
	constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))

	handlersConfPath := path.Join(testDataPath, "handlers")
	err := constructor.ParseHandlerConfFolder(handlersConfPath)
	assert.Nil(t, err)

	dividerConfPath := path.Join(testDataPath, "divider")
	err = constructor.ParseDividerConfFolder(dividerConfPath)
	assert.Nil(t, err)

	t.Run("Input=valid", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		layerConstructor := NewLayerConstructor(constructor)

		err := layerConstructor.ParseLayerConfFolder(confPath)
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

		layerConstructor := NewLayerConstructor(constructor)
		err := layerConstructor.ParseLayerConfFolder(confPath)
		assert.ErrorContains(t, err, "no such file or directory")
	})

	t.Run("Input=invalid/layer_name_not_match", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		layerConstructor := NewLayerConstructor(constructor)
		err := layerConstructor.ParseLayerConfFolder(confPath)
		assert.ErrorContains(t, err, "mismatch with configuration file Name")
	})

	t.Run("Input=invalid/same_layer_name", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		layerConstructor := NewLayerConstructor(constructor)
		err := layerConstructor.ParseLayerConfFolder(confPath)
		assert.ErrorContains(t, err, "layer name sample_e already exists,path is")
	})

	t.Run("Input=invalid/same_layer_name_with_error_path", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)

		layerConstructor := NewLayerConstructor(constructor)
		err := layerConstructor.ParseLayerConfFolder(confPath)
		assert.ErrorContains(t, err, "layer name sample_f already exists")
	})
}

func TestConstructLayerFromFile(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(LayerConstructor{}))
	factory.Register(reflect.TypeOf(dividerconstructor.DividerConstructor{}))
	factory.Register(reflect.TypeOf(handlerconstructor.HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	constructor := constructor.NewConstructor(factory, "")

	constructor.Register(reflect.TypeOf(layer.Layer{}))
	constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))
	constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))

	handlersConfPath := path.Join(testDataPath, "handlers")
	err := constructor.ParseHandlerConfFolder(handlersConfPath)
	assert.Nil(t, err)

	dividerConfPath := path.Join(testDataPath, "divider")
	err = constructor.ParseDividerConfFolder(dividerConfPath)
	assert.Nil(t, err)

	t.Run("Input=InvalidFilePath", func(t *testing.T) {
		layerConstructor := NewLayerConstructor(constructor)
		layerInterface, err := layerConstructor.constructLayerFromFile("")
		assert.Nil(t, layerInterface)
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

		layerConstructor := NewLayerConstructor(constructor)
		layerInterface, err := layerConstructor.constructLayerFromFile(confFilePath)
		assert.Nil(t, layerInterface)
		assert.ErrorContains(t, err, "cannot unmarshal")
	})
}

func TestConstructLayerFromName(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	factory := factory.NewFactory()
	factory.Register(reflect.TypeOf(LayerConstructor{}))
	factory.Register(reflect.TypeOf(dividerconstructor.DividerConstructor{}))
	factory.Register(reflect.TypeOf(handlerconstructor.HandlerConstructor{}))
	factory.Register(reflect.TypeOf(layercenterconstructor.LayerCenterConstructor{}))
	factory.Register(reflect.TypeOf(handlergroupconstructor.HandlerGroupConstructor{}))
	factory.Register(reflect.TypeOf(aynchandlergroupconstructor.AsyncHandlerGroupConstructor{}))
	constructor := constructor.NewConstructor(factory, "")

	constructor.Register(reflect.TypeOf(layer.Layer{}))
	constructor.Register(reflect.TypeOf(sampledivider.SampleAutoConstructDivider{}))
	constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))

	handlersConfPath := path.Join(testDataPath, "handlers")
	err := constructor.ParseHandlerConfFolder(handlersConfPath)
	assert.Nil(t, err)

	dividerConfPath := path.Join(testDataPath, "divider")
	err = constructor.ParseDividerConfFolder(dividerConfPath)
	assert.Nil(t, err)

	t.Run("Input=InvalidFilePath", func(t *testing.T) {
		layerConstructor := NewLayerConstructor(constructor)
		layerInterface, err := layerConstructor.constructLayerFromName("")
		assert.Nil(t, layerInterface)
		assert.ErrorContains(t, err, "layer's configuration is not set")
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

		layerConstructor := NewLayerConstructor(constructor)
		layerInterface, err := layerConstructor.constructLayerFromFile(confFilePath)
		assert.Nil(t, layerInterface)
		assert.ErrorContains(t, err, "cannot unmarshal")
	})
}
