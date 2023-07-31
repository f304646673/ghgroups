package handlergroup

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"ghgroups/frame/utils"

	samplehandler "ghgroups/frame/sample_handler"

	ghgroupscontext "ghgroups/frame/ghgroups_context"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	handlersConfPath := path.Join(testDataPath, "handlers")
	constructor := utils.BuildConstructor(handlersConfPath)

	err := constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))
	assert.Nil(t, err)
	err = constructor.ParseHandlerConfFolder(handlersConfPath)
	assert.Nil(t, err)

	t.Run("Input=not_exist.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.NoFileExists(t, confPath)
		handlerGroup := NewHandlerGroup(constructor)
		err := handlerGroup.LoadConfigFromFile(confPath)
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
		handlerGroup := NewHandlerGroup(constructor)
		err := handlerGroup.LoadConfigFromFile(confPath)
		assert.Nil(t, err)
	})

	t.Run("Input=no_expected_key.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)
		handlerGroup := NewHandlerGroup(constructor)
		err := handlerGroup.LoadConfigFromFile(confPath)
		assert.ErrorContains(t, err, "did not find expected key")
	})

	t.Run("Input=no_handler.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)
		handlerGroup := NewHandlerGroup(constructor)
		err := handlerGroup.LoadConfigFromFile(confPath)
		fmt.Println(confPath)
		assert.ErrorContains(t, err, "object name  sample_z not found")
	})
}

func TestHandle(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)

	handlersConfPath := path.Join(testDataPath, "handlers")
	constructor := utils.BuildConstructor(handlersConfPath)

	err := constructor.Register(reflect.TypeOf(samplehandler.SampleAutoConstructHandler{}))
	assert.Nil(t, err)

	err = constructor.ParseHandlerConfFolder(handlersConfPath)
	assert.Nil(t, err)

	t.Run("Input=valid.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)

		handlerGroup := NewHandlerGroup(constructor)
		err := handlerGroup.LoadConfigFromFile(confPath)
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
		var ctx ghgroupscontext.GhGroupsContext
		suc := handlerGroup.Handle(&ctx)
		assert.True(t, suc)
		assert.Equal(t, output, " sample_a sample_b")
	})
}

func TestAdd(t *testing.T) {
	constructor := utils.BuildConstructor("")

	handlerGroup := NewHandlerGroup(constructor)
	sampleSelfConstructHandlerSingle := samplehandler.NewSampleSelfConstructHandlerSingle()
	err := handlerGroup.Add(sampleSelfConstructHandlerSingle)
	assert.Nil(t, err)

	called := false
	monkey.PatchInstanceMethod(reflect.TypeOf(sampleSelfConstructHandlerSingle), "Handle", func(*samplehandler.SampleSelfConstructHandlerSingle, *ghgroupscontext.GhGroupsContext) bool {
		called = true
		return true
	})

	assert.True(t, handlerGroup.Handle(nil))
	assert.True(t, called)
}
