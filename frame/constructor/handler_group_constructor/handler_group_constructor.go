package handlergroupconstructor

import (
	"fmt"
	"ghgroups/frame"
	"os"
	"path"
	"reflect"

	folder "git-codecommit.us-east-1.amazonaws.com/v1/repos/go-utils.git/folder"

	"gopkg.in/yaml.v2"
)

type HandlerGroupConstructor struct {
	frame.HandlerGroupConstructorInterface
	frame.ConstructorSetterInterface
	constructorInterface frame.ConstructorInterface
	handlers             map[string]frame.HandlerGroupBaseInterface
	handlersConfPath     map[string]string
}

func NewHandlerGroupConstructor(constructorInterface frame.ConstructorInterface) *HandlerGroupConstructor {
	return &HandlerGroupConstructor{
		constructorInterface: constructorInterface,
		handlers:             make(map[string]frame.HandlerGroupBaseInterface),
		handlersConfPath:     make(map[string]string),
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////
// HandlerGroupConstructorInterface

func (h *HandlerGroupConstructor) CreateHandlerGroupWithConfPath(confFilePath string) error {
	handlerName, errGethandlerName := h.getHandlerGroupName(confFilePath)
	if errGethandlerName != nil {
		return errGethandlerName
	}
	if path, ok := h.handlersConfPath[handlerName]; ok {
		return fmt.Errorf("handler name %s already exists,path is %s", handlerName, path)
	}
	h.handlersConfPath[handlerName] = confFilePath
	if h.existHandlerGroup(handlerName) {
		return fmt.Errorf("handler %s already exists", handlerName)
	}

	handlerInterface, err := h.constructHandlerGroupFromName(handlerName)
	if err != nil {
		return err
	}
	err = h.RegisterHandlerGroup(handlerName, handlerInterface)
	if err != nil {
		return err
	}
	return nil
}

func (h *HandlerGroupConstructor) ParseHandlerGroupConfFolder(confFolderPath string) error {
	files, err := folder.OSReadDir(confFolderPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		handlerName, errGethandlerName := h.getHandlerGroupName(f)
		if errGethandlerName != nil {
			return errGethandlerName
		}
		if path, ok := h.handlersConfPath[handlerName]; ok {
			return fmt.Errorf("handler name %s already exists,path is %s", handlerName, path)
		}
		h.handlersConfPath[handlerName] = f
	}

	for handlerName := range h.handlersConfPath {
		if h.existHandlerGroup(handlerName) {
			return fmt.Errorf("handler %s already exists", handlerName)
		}

		handlerInterface, err := h.constructHandlerGroupFromName(handlerName)
		if err != nil {
			return err
		}
		err = h.RegisterHandlerGroup(handlerName, handlerInterface)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *HandlerGroupConstructor) GetHandlerGroup(name string) (frame.HandlerGroupBaseInterface, error) {
	handlerInterface, exist := h.handlers[name]
	if exist {
		return handlerInterface, nil
	}
	handlerInterface, err := h.constructHandlerGroupFromName(name)
	if err != nil {
		return nil, err
	}
	err = h.RegisterHandlerGroup(name, handlerInterface)
	if err != nil {
		return nil, err
	}
	return handlerInterface, nil
}

func (h *HandlerGroupConstructor) RegisterHandlerGroup(name string, handler_interface frame.HandlerGroupBaseInterface) error {
	if h.existHandlerGroup(name) {
		return fmt.Errorf("handler %s is exist", name)
	}
	h.handlers[name] = handler_interface
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *HandlerGroupConstructor) constructHandlerGroupFromName(handlerName string) (frame.HandlerGroupBaseInterface, error) {
	handlerConfPath, ok := h.handlersConfPath[handlerName]
	if !ok {
		return nil, fmt.Errorf("%s: handler's configuration is not set", handlerName)
	}
	handlerInterface, err := h.constructHandlerGroupFromFile(handlerConfPath)
	if err != nil {
		return nil, err
	}

	if handlerInterface.Name() != handlerName {
		return nil, fmt.Errorf("handler Name (%s) mismatch with configuration file Name (%s) path (%s)", handlerInterface.Name(), handlerName, handlerConfPath)
	}
	return handlerInterface, err
}

func (h *HandlerGroupConstructor) constructHandlerGroupFromFile(filePath string) (frame.HandlerGroupBaseInterface, error) {
	configure, errReadFile := os.ReadFile(filePath)
	if errReadFile != nil {
		return nil, errReadFile
	}
	handlerConf := make(map[any]any)
	errReadFile = yaml.Unmarshal([]byte(configure), &handlerConf)
	if errReadFile != nil {
		return nil, errReadFile
	}
	return h.constructHandlerGroup(handlerConf)
}

func (h *HandlerGroupConstructor) constructHandlerGroup(handler_conf map[any]any) (frame.HandlerGroupBaseInterface, error) {
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
				handlerInterface, ok := concrete.(frame.HandlerGroupBaseInterface)
				if !ok {
					return nil, fmt.Errorf("convert %s to HandlerGroupBaseInterface error", typeName)
				}

				return handlerInterface, nil
			}
		}
	}
	return nil, fmt.Errorf("handler configuration must have type %v", handler_conf)
}

func (h *HandlerGroupConstructor) getHandlerGroupName(filePath string) (string, error) {
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

func (h *HandlerGroupConstructor) existHandlerGroup(name string) bool {
	_, exist := h.handlers[name]
	return exist
}

func (h *HandlerGroupConstructor) LoadConfigFromFile(confPath string) error {
	return nil
}

func (h *HandlerGroupConstructor) LoadConfigFromMemory(configure []byte) error {
	return nil
}

func (h *HandlerGroupConstructor) LoadEnvironmentConf() error {
	h.handlers = make(map[string]frame.HandlerGroupBaseInterface)
	h.handlersConfPath = make(map[string]string)
	return nil
}

// //////////////////////////////////////////////////////////////////////////////////
func (h *HandlerGroupConstructor) Name() string {
	return reflect.TypeOf(*h).Name()
}

func (h *HandlerGroupConstructor) SetConstructorInterface(constructorInterface any) {
	constructorInterfaceNew, ok := constructorInterface.(frame.ConstructorInterface)
	if !ok {
		panic("constructorInterface is not frame.ConstructorInterface")
	}
	h.constructorInterface = constructorInterfaceNew
}
