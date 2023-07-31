package samplehandler

import (
	"fmt"
	"ghgroups/frame/utils"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

	ghgroupscontext "ghgroups/frame/ghgroups_context"
)

func TestSampleSelfConstructHandlerSingle(t *testing.T) {
	constructor := utils.BuildConstructor("")

	sampleSelfConstructHandlerSingle := NewSampleSelfConstructHandlerSingle()
	constructor.RegisterHandler(sampleSelfConstructHandlerSingle.Name(), sampleSelfConstructHandlerSingle)

	name := sampleSelfConstructHandlerSingle.Name()
	handlerBaseInterface, err := constructor.GetHandler(name)
	assert.Nil(t, err)

	patch := monkey.Patch(fmt.Sprintln, func(a ...interface{}) string {
		s := make([]interface{}, len(a))
		output := ""
		for i, v := range a {
			s[i] = fmt.Sprint(v)
			output = fmt.Sprintf("%s", s[i])
		}
		assert.Equal(t, name, output)
		return output
	})
	defer patch.Unpatch()
	var context ghgroupscontext.GhGroupsContext
	suc := handlerBaseInterface.Handle(&context)
	assert.True(t, suc)
}
