package sampledivider

import (
	"fmt"
	"ghgroups/frame/utils"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

	ghgroupscontext "ghgroups/frame/ghgroups_context"
)

func TestSampleSelfConstructDividerMulti(t *testing.T) {
	constructor := utils.BuildConstructor("")

	sampleSelfConstructDividerMultiNameA := "sample_self_construct_divider_multi_a"
	sampleSelfConstructDividerMultiA := NewSampleSelfConstructDividerMulti(sampleSelfConstructDividerMultiNameA)
	constructor.RegisterDivider(sampleSelfConstructDividerMultiA.Name(), sampleSelfConstructDividerMultiA)

	sampleSelfConstructDividerMultiNameB := "sample_self_construct_divider_multi_b"
	sampleSelfConstructDividerMultiB := NewSampleSelfConstructDividerMulti(sampleSelfConstructDividerMultiNameB)
	constructor.RegisterDivider(sampleSelfConstructDividerMultiB.Name(), sampleSelfConstructDividerMultiB)

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
		t.Run("Input="+sampleSelfConstructDividerMultiNameA, parallelTest)
		t.Run("Input="+sampleSelfConstructDividerMultiNameB, parallelTest)
	})

}
