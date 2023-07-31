package sampledivider

import (
	"fmt"
	"ghgroups/frame/utils"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"bou.ke/monkey"

	filemanager "git-codecommit.us-east-1.amazonaws.com/v1/repos/go-utils.git/file_manager"

	ghgroupscontext "ghgroups/frame/ghgroups_context"

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
		sampleAutoConstructDivider := NewSampleAutoConstructDivider()
		err = sampleAutoConstructDivider.LoadConfigFromMemory(data)
		assert.ErrorContains(t, err, "yaml: unmarshal errors")
	})

	t.Run("Input=valid/sample_divider_a.yaml", func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		confPath := path.Join(testDataPath, confName)
		assert.FileExists(t, confPath)
		data, err := filemanager.GetFileContent(confPath)
		assert.Nil(t, err)
		sampleAutoConstructDivider := NewSampleAutoConstructDivider()
		err = sampleAutoConstructDivider.LoadConfigFromMemory(data)
		assert.Nil(t, err)
	})
}

func TestSampleAutoConstructDivider(t *testing.T) {

	runPath, errGetWd := os.Getwd()
	assert.Nil(t, errGetWd)
	testDataPath := path.Join(runPath, "test_data/valid")

	constructor := utils.BuildConstructor("")

	// 因为通过反射+配置文件创建divider，所以需要注册divider的类型
	constructor.Register(reflect.TypeOf(SampleAutoConstructDivider{}))

	err := constructor.ParseDividerConfFolder(testDataPath)
	assert.Nil(t, err)

	sampleDividerNameA := "sample_divider_a"
	sampleDividerNameB := "sample_divider_b"

	parallelTest := func(t *testing.T) {
		t.Parallel()
		testName := t.Name()
		comma := strings.Index(testName, "=")
		assert.Greater(t, comma, 0)
		confName := testName[comma+1:]
		dividerBaseInterface, err := constructor.GetDivider(confName)
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
		handlerName := dividerBaseInterface.Select(&context)
		assert.Equal(t, handlerName, "SampleSelfConstructHandlerSingle")
	}

	t.Run("group", func(t *testing.T) {
		t.Run("Input="+sampleDividerNameA, parallelTest)
		t.Run("Input="+sampleDividerNameB, parallelTest)
	})
}
