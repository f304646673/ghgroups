package layercenterconstructor

import (
	"fmt"
	"ghgroups/frame"
	"os"
	"path"
	"reflect"

	folder "git-codecommit.us-east-1.amazonaws.com/v1/repos/go-utils.git/folder"

	"gopkg.in/yaml.v2"
)

type LayerCenterConstructor struct {
	frame.LayerCenterConstructorInterface
	frame.ConstructorSetterInterface
	constructorInterface frame.ConstructorInterface
	handlers             map[string]frame.LayerCenterBaseInterface
	handlersConfPath     map[string]string
}

func NewLayerCenterConstructor(constructorInterface frame.ConstructorInterface) *LayerCenterConstructor {
	return &LayerCenterConstructor{
		constructorInterface: constructorInterface,
		handlers:             make(map[string]frame.LayerCenterBaseInterface),
		handlersConfPath:     make(map[string]string),
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////
// LayerCenterConstructorInterface

func (h *LayerCenterConstructor) CreateLayerCenterWithConfPath(confFilePath string) error {
	handlerName, errGethandlerName := h.getLayerCenterName(confFilePath)
	if errGethandlerName != nil {
		return errGethandlerName
	}
	if path, ok := h.handlersConfPath[handlerName]; ok {
		return fmt.Errorf("handler name %s already exists,path is %s", handlerName, path)
	}
	h.handlersConfPath[handlerName] = confFilePath
	if h.existLayerCenter(handlerName) {
		return fmt.Errorf("handler %s already exists", handlerName)
	}

	handlerInterface, err := h.constructLayerCenterFromName(handlerName)
	if err != nil {
		return err
	}
	err = h.RegisterLayerCenter(handlerName, handlerInterface)
	if err != nil {
		return err
	}
	return nil
}

func (h *LayerCenterConstructor) ParseLayerCenterConfFolder(confFolderPath string) error {
	files, err := folder.OSReadDir(confFolderPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		handlerName, errGethandlerName := h.getLayerCenterName(f)
		if errGethandlerName != nil {
			return errGethandlerName
		}
		if path, ok := h.handlersConfPath[handlerName]; ok {
			return fmt.Errorf("handler name %s already exists,path is %s", handlerName, path)
		}
		h.handlersConfPath[handlerName] = f
	}

	for handlerName := range h.handlersConfPath {
		if h.existLayerCenter(handlerName) {
			return fmt.Errorf("handler %s already exists", handlerName)
		}

		handlerInterface, err := h.constructLayerCenterFromName(handlerName)
		if err != nil {
			return err
		}
		err = h.RegisterLayerCenter(handlerName, handlerInterface)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *LayerCenterConstructor) GetLayerCenter(name string) (frame.LayerCenterBaseInterface, error) {
	handlerInterface, exist := h.handlers[name]
	if exist {
		return handlerInterface, nil
	}
	handlerInterface, err := h.constructLayerCenterFromName(name)
	if err != nil {
		return nil, err
	}
	err = h.RegisterLayerCenter(name, handlerInterface)
	if err != nil {
		return nil, err
	}
	return handlerInterface, nil
}

func (h *LayerCenterConstructor) RegisterLayerCenter(name string, handler_interface frame.LayerCenterBaseInterface) error {
	if h.existLayerCenter(name) {
		return fmt.Errorf("handler %s is exist", name)
	}
	h.handlers[name] = handler_interface
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *LayerCenterConstructor) constructLayerCenterFromName(handlerName string) (frame.LayerCenterBaseInterface, error) {
	handlerConfPath, ok := h.handlersConfPath[handlerName]
	if !ok {
		return nil, fmt.Errorf("%s: handler's configuration is not set", handlerName)
	}
	handlerInterface, err := h.constructLayerCenterFromFile(handlerConfPath)
	if err != nil {
		return nil, err
	}

	if handlerInterface.Name() != handlerName {
		return nil, fmt.Errorf("handler Name (%s) mismatch with configuration file Name (%s) path (%s)", handlerInterface.Name(), handlerName, handlerConfPath)
	}
	return handlerInterface, err
}

func (h *LayerCenterConstructor) constructLayerCenterFromFile(filePath string) (frame.LayerCenterBaseInterface, error) {
	configure, errReadFile := os.ReadFile(filePath)
	if errReadFile != nil {
		return nil, errReadFile
	}
	handlerConf := make(map[any]any)
	errReadFile = yaml.Unmarshal([]byte(configure), &handlerConf)
	if errReadFile != nil {
		return nil, errReadFile
	}
	return h.constructLayerCenter(handlerConf)
}

func (h *LayerCenterConstructor) constructLayerCenter(handler_conf map[any]any) (frame.LayerCenterBaseInterface, error) {
	for k, v := range handler_conf {
		switch k.(type) {
		case string:
			if k == "type" {
				originConf, _ := yaml.Marshal(handler_conf)
				typeName := v.(string)
				concrete, err := h.constructorInterface.Create(typeName, originConf, h.constructorInterface)
				if err != nil {
					return nil, err
				}
				handlerInterface, ok := concrete.(frame.LayerCenterBaseInterface)
				if !ok {
					return nil, fmt.Errorf("convert %s to LayerCenterBaseInterface error", typeName)
				}

				return handlerInterface, nil
			}
		}
	}
	return nil, fmt.Errorf("handler configuration must have type %v", handler_conf)
}

func (h *LayerCenterConstructor) getLayerCenterName(filePath string) (string, error) {
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

func (h *LayerCenterConstructor) existLayerCenter(name string) bool {
	_, exist := h.handlers[name]
	return exist
}

func (h *LayerCenterConstructor) LoadConfigFromFile(confPath string) error {
	return nil
}

func (h *LayerCenterConstructor) LoadConfigFromMemory(configure []byte) error {
	return nil
}

func (h *LayerCenterConstructor) LoadEnvironmentConf() error {
	h.handlers = make(map[string]frame.LayerCenterBaseInterface)
	h.handlersConfPath = make(map[string]string)
	return nil
}

// //////////////////////////////////////////////////////////////////////////////////
func (h *LayerCenterConstructor) Name() string {
	return reflect.TypeOf(*h).Name()
}

func (h *LayerCenterConstructor) SetConstructorInterface(constructorInterface any) {
	constructorInterfaceNew, ok := constructorInterface.(frame.ConstructorInterface)
	if !ok {
		panic("constructorInterface is not frame.ConstructorInterface")
	}
	h.constructorInterface = constructorInterfaceNew
}
