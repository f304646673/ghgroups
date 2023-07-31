package layerconstructor

import (
	"fmt"
	"ghgroups/frame"
	"os"
	"path"
	"reflect"

	folder "git-codecommit.us-east-1.amazonaws.com/v1/repos/go-utils.git/folder"
	"gopkg.in/yaml.v2"
)

type LayerConstructor struct {
	frame.LayerConstructorInterface
	frame.ConstructorSetterInterface
	constructorInterface frame.ConstructorInterface
	layers               map[string]frame.LayerBaseInterface
	layersConfPath       map[string]string
}

func NewLayerConstructor(constructorInterface frame.ConstructorInterface) *LayerConstructor {
	return &LayerConstructor{
		constructorInterface: constructorInterface,
		layers:               make(map[string]frame.LayerBaseInterface),
		layersConfPath:       make(map[string]string),
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////
// LayerConstructorInterface
func (l *LayerConstructor) CreateLayerWithConfPath(confFilePath string) error {
	layerName, errGetlayerName := l.getLayerName(confFilePath)
	if errGetlayerName != nil {
		return errGetlayerName
	}
	if path, ok := l.layersConfPath[layerName]; ok {
		return fmt.Errorf("layer name %s already exists,path is %s", layerName, path)
	}
	l.layersConfPath[layerName] = confFilePath
	if l.existLayer(layerName) {
		return fmt.Errorf("layer %s already exists", layerName)
	}

	layerInterface, err := l.constructLayerFromName(layerName)
	if err != nil {
		return err
	}
	err = l.RegisterLayer(layerName, layerInterface)
	if err != nil {
		return err
	}
	return nil
}

func (l *LayerConstructor) ParseLayerConfFolder(confFolderPath string) error {
	files, err := folder.OSReadDir(confFolderPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		layerName, errGetlayerName := l.getLayerName(f)
		if errGetlayerName != nil {
			return errGetlayerName
		}
		if path, ok := l.layersConfPath[layerName]; ok {
			return fmt.Errorf("layer name %s already exists,path is %s", layerName, path)
		}
		l.layersConfPath[layerName] = f
	}

	for layerName := range l.layersConfPath {
		if l.existLayer(layerName) {
			return fmt.Errorf("layer %s already exists", layerName)
		}

		layerInterface, err := l.constructLayerFromName(layerName)
		if err != nil {
			return err
		}
		err = l.RegisterLayer(layerName, layerInterface)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LayerConstructor) GetLayer(name string) (frame.LayerBaseInterface, error) {
	layerInterface, exist := l.layers[name]
	if exist {
		return layerInterface, nil
	}
	layerInterface, err := l.constructLayerFromName(name)
	if err != nil {
		return nil, err
	}
	err = l.RegisterLayer(name, layerInterface)
	if err != nil {
		return nil, err
	}
	return layerInterface, nil
}

func (l *LayerConstructor) RegisterLayer(name string, layer_interface frame.LayerBaseInterface) error {
	if l.existLayer(name) {
		return fmt.Errorf("layer %s is exist", name)
	}
	l.layers[name] = layer_interface
	return nil
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////

func (l *LayerConstructor) constructLayerFromName(layerName string) (frame.LayerWithBuilderInterface, error) {
	layerConfPath, ok := l.layersConfPath[layerName]
	if !ok {
		return nil, fmt.Errorf("%s: layer's configuration is not set", layerName)
	}
	layerInterface, err := l.constructLayerFromFile(layerConfPath)
	if err != nil {
		return nil, err
	}

	if layerInterface.Name() != layerName {
		return nil, fmt.Errorf("layer Name (%s) mismatch with configuration file Name (%s) path (%s)", layerInterface.Name(), layerName, layerConfPath)
	}
	return layerInterface, err
}

func (l *LayerConstructor) constructLayerFromFile(filePath string) (frame.LayerWithBuilderInterface, error) {
	configure, errReadFile := os.ReadFile(filePath)
	if errReadFile != nil {
		return nil, errReadFile
	}
	layerConf := make(map[any]any)
	errReadFile = yaml.Unmarshal([]byte(configure), &layerConf)
	if errReadFile != nil {
		return nil, errReadFile
	}
	return l.constructLayer(layerConf)
}

func (l *LayerConstructor) constructLayer(layer_conf map[any]any) (frame.LayerWithBuilderInterface, error) {
	for k, v := range layer_conf {
		switch k.(type) {
		case string:
			if k == "type" {
				originConf, _ := yaml.Marshal(layer_conf)
				typeName := v.(string)

				layer, err := l.constructorInterface.Create(typeName, originConf, l.constructorInterface)
				layerInterface := layer.(frame.LayerWithBuilderInterface)
				// layerInterface.SetConstructorInterface(l.constructorInterface)
				// layerInterface.LoadEnvironmentConf(env, region, l.constructorInterface)
				if err != nil {
					return nil, err
				}

				return layerInterface, nil
			}
		}
	}
	return nil, fmt.Errorf("layer configuration must have type %v", layer_conf)
}

func (l *LayerConstructor) getLayerName(filePath string) (string, error) {
	s, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	if s.IsDir() {
		return "", fmt.Errorf("file path %s is not a directory", filePath)
	}

	fileNameWithExt := path.Base(filePath)
	ext := path.Ext(filePath)
	fileName := fileNameWithExt[0 : len(fileNameWithExt)-len(ext)]
	return fileName, nil
}

func (l *LayerConstructor) existLayer(name string) bool {
	_, exist := l.layers[name]
	return exist
}

func (l *LayerConstructor) SetConstructorInterface(constructorInterface any) {
	constructorInterfaceNew, ok := constructorInterface.(frame.ConstructorInterface)
	if !ok {
		panic("constructorInterface is not frame.ConstructorInterface")
	}
	l.constructorInterface = constructorInterfaceNew
}

// //////////////////////////////////////////////////////////////////////////////////
func (l *LayerConstructor) LoadConfigFromFile(confPath string) error {
	return nil
}

func (l *LayerConstructor) LoadConfigFromMemory(configure []byte) error {
	return nil
}

func (l *LayerConstructor) LoadEnvironmentConf() error {
	l.layers = make(map[string]frame.LayerBaseInterface)
	l.layersConfPath = make(map[string]string)
	return nil
}

// //////////////////////////////////////////////////////////////////////////////////
func (l *LayerConstructor) Name() string {
	return reflect.TypeOf(*l).Name()
}

// //////////////////////////////////////////////////////////////////////////////////
