package concreteconfmanager

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLayerName(t *testing.T) {
	runPath, errGetWd := os.Getwd()
	testDataPath := path.Join(runPath, "test_data")
	assert.Nil(t, errGetWd)
	t.Run("same_file_name", func(t *testing.T) {
		same_file_name := testDataPath + "/same_file_name"
		concreteConfManager := NewConcreteConfManager()
		err := concreteConfManager.ParseConfFolder(same_file_name)
		assert.NotNil(t, err)
	})

	t.Run("same_file_name_without_suffix", func(t *testing.T) {
		same_file_name_without_suffix := testDataPath + "/same_file_name_without_suffix"
		concreteConfManager := NewConcreteConfManager()
		err := concreteConfManager.ParseConfFolder(same_file_name_without_suffix)
		assert.NotNil(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		valid := testDataPath + "/valid"
		concreteConfManager := NewConcreteConfManager()
		err := concreteConfManager.ParseConfFolder(valid)
		assert.Nil(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		invalid := testDataPath + "/invalid"
		concreteConfManager := NewConcreteConfManager()
		err := concreteConfManager.ParseConfFolder(invalid)
		assert.NotNil(t, err)
	})

}
