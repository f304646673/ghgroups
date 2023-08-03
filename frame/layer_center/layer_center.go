package layercenter

import (
	"fmt"
	"ghgroups/frame"
	debughelper "ghgroups/frame/debug_helper"
	ghgroupscontext "ghgroups/frame/ghgroups_context"
	"os"

	"gopkg.in/yaml.v2"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type LayerCenterConf struct {
	Type   string   `yaml:"type"`
	Name   string   `yaml:"name"`
	Layers []string `yaml:"layers"`
}

type LayerCenter struct {
	frame.LayerCenterInterface
	frame.ConstructorSetterInterface
	constructorInterface frame.ConstructorInterface
	conf                 LayerCenterConf
	layers               []frame.LayerBaseInterface
}

func NewLayerCenter(constructorInterface frame.ConstructorInterface) *LayerCenter {
	return &LayerCenter{
		layers:               make([]frame.LayerBaseInterface, 0),
		constructorInterface: constructorInterface,
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// LayerCenterInterface

func (l *LayerCenter) init() error {
	for _, layerName := range l.conf.Layers {
		if err := l.constructorInterface.CreateConcrete(layerName); err != nil {
			return err
		}

		if someInterface, err := l.constructorInterface.GetConcrete(layerName); err != nil {
			return err
		} else {
			if layerBaseInterface, ok := someInterface.(frame.LayerBaseInterface); !ok {
				return fmt.Errorf("layer %s is not frame.LayerBaseInterface", layerName)
			} else {
				l.layers = append(l.layers, layerBaseInterface)
			}
		}
	}

	return nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (l *LayerCenter) Name() string {
	return l.conf.Name
}

func (l *LayerCenter) Add(layerInterface frame.LayerWithBuilderInterface) {
	l.layers = append(l.layers, layerInterface)
}

func (l *LayerCenter) Handle(ctx *ghgroupscontext.GhGroupsContext) bool {
	for _, layer := range l.layers {
		if debughelper.HandleWithShowDuration(layer, layer.Name(), ctx) {
			continue
		} else {
			return false
		}
	}
	return true
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (l *LayerCenter) LoadConfigFromFile(confPath string) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	return l.LoadConfigFromMemory(data)
}

func (l *LayerCenter) LoadConfigFromMemory(configure []byte) error {
	var layerCenterConf LayerCenterConf
	err := yaml.Unmarshal(configure, &layerCenterConf)
	if err != nil {
		return err
	}
	l.conf = layerCenterConf
	return l.init()
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (l *LayerCenter) SetConstructorInterface(constructorInterface any) {
	constructorInterfaceNew, ok := constructorInterface.(frame.ConstructorInterface)
	if !ok {
		panic("constructorInterface is not frame.ConstructorInterface")
	}
	l.constructorInterface = constructorInterfaceNew
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
