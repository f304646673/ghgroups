package factory

import (
	"fmt"
	"ghgroups/frame"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	type TestSameName struct {
	}

	t.Run("Input=RegisterSuc", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.Nil(t, r)
		}()

		factory := NewFactory()
		factory.Register(reflect.TypeOf(TestSameName{}))
	})

	t.Run("Input=RegisterPanic", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.Equal(t, r, "TestSameName is already registered.Please modify name.")
		}()

		factory := NewFactory()
		factory.Register(reflect.TypeOf(TestSameName{}))
		factory.Register(reflect.TypeOf(TestSameName{}))
	})
}

type TestConcreteWithInterface struct {
	frame.ConcreteInterface
}

func (t *TestConcreteWithInterface) LoadConfigFromFile(confPath string) error {
	return nil
}

func (t *TestConcreteWithInterface) LoadConfigFromMemory(configure []byte) error {
	return nil
}

func (t *TestConcreteWithInterface) LoadEnvironmentConf() error {
	return nil
}

func (t *TestConcreteWithInterface) Name() string {
	return "TestConcreteWithInterface"
}

func TestCreate(t *testing.T) {
	type TestConcrete struct {
	}

	t.Run("Input=CreateNotExist", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.Equal(t, r, "concrete name not found for TestConcrete")
		}()

		factory := NewFactory()
		factory.Create(reflect.TypeOf(TestConcrete{}).Name(), nil, nil)
	})

	t.Run("Input=CreateNotImplementsInterface", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.Equal(t, r, "concrete TestConcrete conver to ConcreteInterface error")
		}()

		factory := NewFactory()
		factory.Register(reflect.TypeOf(TestConcrete{}))
		factory.Create(reflect.TypeOf(TestConcrete{}).Name(), nil, nil)
	})

	t.Run("Input=CreateInitError", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.Equal(t, r, "concrete LoadConfigFromMemory error for TestConcreteWithInterface .error: ")
		}()

		factory := NewFactory()
		factory.Register(reflect.TypeOf(TestConcreteWithInterface{}))
		testConcreteWithInterface := TestConcreteWithInterface{}

		monkey.PatchInstanceMethod(reflect.TypeOf(&testConcreteWithInterface), "LoadConfigFromMemory", func(*TestConcreteWithInterface, []byte) error {
			return fmt.Errorf("")
		})
		defer monkey.UnpatchAll()

		factory.Create(reflect.TypeOf(testConcreteWithInterface).Name(), nil, nil)
	})

	t.Run("Input=CreateInitError", func(t *testing.T) {
		factory := NewFactory()
		factory.Register(reflect.TypeOf(TestConcreteWithInterface{}))
		testConcreteWithInterface := TestConcreteWithInterface{}
		i, e := factory.Create(reflect.TypeOf(testConcreteWithInterface).Name(), nil, nil)
		assert.NotNil(t, i)
		assert.Nil(t, e)
	})
}
