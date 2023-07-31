package sampledivider

import (
	"fmt"
	"ghgroups/frame/utils"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

	ghgroupscontext "ghgroups/frame/ghgroups_context"
)

func TestSampleSelfConstructDividerSingle(t *testing.T) {
	constructor := utils.BuildConstructor("")

	sampleSelfConstructDividerSingle := NewSampleSelfConstructDividerSingle()
	constructor.RegisterDivider(sampleSelfConstructDividerSingle.Name(), sampleSelfConstructDividerSingle)

	name := "SampleSelfConstructDividerSingle"
	dividerBaseInterface, err := constructor.GetDivider(name)
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
	handlerName := dividerBaseInterface.Select(&context)
	assert.Equal(t, handlerName, "SampleSelfConstructHandlerSingle")
}
