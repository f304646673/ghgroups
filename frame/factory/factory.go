package factory

import (
	"fmt"
	"ghgroups/frame"
	"reflect"
)

type Factory struct {
	frame.FactoryInterface
	concretesType map[string]reflect.Type
	concretes     map[string]frame.ConcreteInterface
}

func NewFactory() *Factory {

	return &Factory{
		concretesType: make(map[string]reflect.Type),
		concretes:     make(map[string]frame.ConcreteInterface),
	}
}

func (f *Factory) Create(concreteTypeName string, configure []byte, constructorInterface any) (concrete any, err error) {
	concreteType, ok := f.concretesType[concreteTypeName]
	if !ok {
		err = fmt.Errorf("concrete name not found for %s", concreteTypeName)
		panic(err.Error())
	}

	concrete = reflect.New(concreteType).Interface()
	concreteInterface, ok := concrete.(frame.ConcreteInterface)
	if !ok || concreteInterface == nil {
		err = fmt.Errorf("concrete %s conver to ConcreteInterface error", concreteTypeName)
		panic(err.Error())
	}

	constructorSetterInterface, ok := concrete.(frame.ConstructorSetterInterface)
	if ok && constructorSetterInterface != nil && constructorInterface != nil {
		constructorSetterInterface.SetConstructorInterface(constructorInterface)
	}

	loadConfigFromMemoryInterface, ok := concrete.(frame.LoadConfigFromMemoryInterface)
	if ok && loadConfigFromMemoryInterface != nil {
		err := loadConfigFromMemoryInterface.LoadConfigFromMemory(configure)
		if err != nil {
			err = fmt.Errorf("concrete LoadConfigFromMemory error for %s .error: %v", concreteTypeName, err)
			panic(err.Error())
		}
	}

	loadEnvironmentConfInterface, ok := concrete.(frame.LoadEnvironmentConfInterface)
	if ok && loadEnvironmentConfInterface != nil {
		err := loadEnvironmentConfInterface.LoadEnvironmentConf()
		if err != nil {
			err = fmt.Errorf("concrete LoadEnvironmentConf error for %s .error: %v", concreteTypeName, err)
			panic(err.Error())
		}
	}

	concreteName := concreteInterface.Name()
	if concreteName == "" {
		err = fmt.Errorf("concrete's Name is empty for %s", concreteTypeName)
		panic(err.Error())
	}
	if _, ok := f.concretes[concreteName]; !ok {
		f.concretes[concreteName] = concreteInterface
	} else {
		err = fmt.Errorf("concrete's Name  %s has already been used, type is %s", concreteName, concreteTypeName)
		panic(err.Error())
	}
	return
}

func (f *Factory) Register(concreteType reflect.Type) (err error) {
	concreteName := concreteType.Name()
	if _, ok := f.concretesType[concreteName]; !ok {
		f.concretesType[concreteName] = concreteType
	} else {
		err = fmt.Errorf(concreteName + " is already registered.Please modify name.")
		panic(err.Error())
	}
	return nil
}

func (f *Factory) Get(concreteName string) (frame.ConcreteInterface, error) {
	if concrete, ok := f.concretes[concreteName]; ok {
		return concrete, nil
	} else {
		return nil, nil
	}
}
