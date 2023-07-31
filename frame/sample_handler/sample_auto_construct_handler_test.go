package samplehandler

import (
	"fmt"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"ghgroups/frame/utils"

	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	filemanager "git-codecommit.us-east-1.amazonaws.com/v1/repos/go-utils.git/file_manager"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromMemory(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	assert.Nil(t, errGetWd)
	testDataPath := path.Join(runPath, "test_data")

	t.Run("Input=invalid/sample_error.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)
		data, err := filemanager.GetFileContent(confPath)
		assert.Nil(t, err)
		sampleAutoConstructHandler := NewSampleAutoConstructHandler()
		err = sampleAutoConstructHandler.LoadConfigFromMemory(data)
		assert.ErrorContains(t, err, "yaml: unmarshal errors")
	})

	t.Run("Input=valid/sample_handler_a.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)
		data, err := filemanager.GetFileContent(confPath)
		assert.Nil(t, err)
		sampleAutoConstructHandler := NewSampleAutoConstructHandler()
		err = sampleAutoConstructHandler.LoadConfigFromMemory(data)
		assert.Nil(t, err)
	})
}

func TestSampleAutoConstructHandler(t *testing.T) {

	runPath, errGetWd := os.Getwd()
	assert.Nil(t, errGetWd)
	testDataPath := path.Join(runPath, "test_data/valid")

	constructor := utils.BuildConstructor("")

	// 因为通过反射+配置文件创建handler，所以需要注册handler的类型
	constructor.Register(reflect.TypeOf(SampleAutoConstructHandler{}))

	err := constructor.ParseHandlerConfFolder(testDataPath)
	assert.Nil(t, err)

	sampleHandlerNameA := "sample_handler_a"
	sampleHandlerNameB := "sample_handler_b"

	parallelTest := func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		handlerBaseInterface, err := constructor.GetHandler(confName)
		assert.Nil(t, err)

		patch := monkey.Patch(fmt.Sprintln, func(a ...interface{}) string {
			s := make([]interface{}, len(a))
			output := ""
			for i, v := range a {
				s[i] = fmt.Sprint(v)
				output = fmt.Sprintf("%s", s[i])
			}
			assert.Equal(t, confName, output)
			return output
		})
		defer patch.Unpatch()
		var context ghgroupscontext.GhGroupsContext
		suc := handlerBaseInterface.Handle(&context)
		assert.True(t, suc)
	}

	t.Run("group", func(t *testing.T) {
		t.Run("Input="+sampleHandlerNameA, parallelTest)
		t.Run("Input="+sampleHandlerNameB, parallelTest)
	})
}
