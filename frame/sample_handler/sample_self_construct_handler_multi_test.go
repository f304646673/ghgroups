package samplehandler

import (
	"fmt"
	"ghgroups/frame/utils"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

	ghgroupscontext "ghgroups/frame/ghgroups_context"
)

func TestSampleSelfConstructHandlerMulti(t *testing.T) {
	constructor := utils.BuildConstructor("")

	sampleSelfConstructHandlerMultiNameA := "sample_self_construct_handler_multi_a"
	sampleSelfConstructHandlerMultiA := NewSampleSelfConstructHandlerMulti(sampleSelfConstructHandlerMultiNameA)
	constructor.RegisterHandler(sampleSelfConstructHandlerMultiA.Name(), sampleSelfConstructHandlerMultiA)

	sampleSelfConstructHandlerMultiNameB := "sample_self_construct_handler_multi_b"
	sampleSelfConstructHandlerMultiB := NewSampleSelfConstructHandlerMulti(sampleSelfConstructHandlerMultiNameB)
	constructor.RegisterHandler(sampleSelfConstructHandlerMultiB.Name(), sampleSelfConstructHandlerMultiB)

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
		t.Run("Input="+sampleSelfConstructHandlerMultiNameA, parallelTest)
		t.Run("Input="+sampleSelfConstructHandlerMultiNameB, parallelTest)
	})

}
