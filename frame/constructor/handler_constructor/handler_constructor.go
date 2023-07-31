package handlerconstructor

import (
	"fmt"
	"ghgroups/frame"
	"os"
	"path"

	"reflect"

	folder "git-codecommit.us-east-1.amazonaws.com/v1/repos/go-utils.git/folder"
	"gopkg.in/yaml.v2"
)

type HandlerConstructor struct {
	frame.HandlerConstructorInterface
	frame.ConstructorSetterInterface
	constructorInterface frame.ConstructorInterface
	handlers             map[string]frame.HandlerBaseInterface
	handlersConfPath     map[string]string
}

func NewHandlerConstructor(constructorInterface frame.ConstructorInterface) *HandlerConstructor {
	return &HandlerConstructor{
		constructorInterface: constructorInterface,
		handlers:             make(map[string]frame.HandlerBaseInterface),
		handlersConfPath:     make(map[string]string),
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////
// HandlerConstructorInterface

func (h *HandlerConstructor) CreateHandlerWithConfPath(confFilePath string) error {
	handlerName, errGethandlerName := h.getHandlerName(confFilePath)
	if errGethandlerName != nil {
		return errGethandlerName
	}
	if path, ok := h.handlersConfPath[handlerName]; ok {
		return fmt.Errorf("handler name %s already exists,path is %s", handlerName, path)
	}
	h.handlersConfPath[handlerName] = confFilePath
	if h.existHandler(handlerName) {
		return fmt.Errorf("handler %s already exists", handlerName)
	}

	handlerInterface, err := h.constructHandlerFromName(handlerName)
	if err != nil {
		return err
	}
	err = h.RegisterHandler(handlerName, handlerInterface)
	if err != nil {
		return err
	}
	return nil
}

func (h *HandlerConstructor) ParseHandlerConfFolder(confFolderPath string) error {
	files, err := folder.OSReadDir(confFolderPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		handlerName, errGethandlerName := h.getHandlerName(f)
		if errGethandlerName != nil {
			return errGethandlerName
		}
		if path, ok := h.handlersConfPath[handlerName]; ok {
			return fmt.Errorf("handler name %s already exists,path is %s", handlerName, path)
		}
		h.handlersConfPath[handlerName] = f
	}

	for handlerName := range h.handlersConfPath {
		if h.existHandler(handlerName) {
			return fmt.Errorf("handler %s already exists", handlerName)
		}

		handlerInterface, err := h.constructHandlerFromName(handlerName)
		if err != nil {
			return err
		}
		err = h.RegisterHandler(handlerName, handlerInterface)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *HandlerConstructor) GetHandler(name string) (frame.HandlerBaseInterface, error) {
	handlerInterface, exist := h.handlers[name]
	if exist {
		return handlerInterface, nil
	}
	handlerInterface, err := h.constructHandlerFromName(name)
	if err != nil {
		return nil, err
	}
	err = h.RegisterHandler(name, handlerInterface)
	if err != nil {
		return nil, err
	}
	return handlerInterface, nil
}

func (h *HandlerConstructor) RegisterHandler(name string, handler_interface frame.HandlerBaseInterface) error {
	if h.existHandler(name) {
		return fmt.Errorf("handler %s is exist", name)
	}
	h.handlers[name] = handler_interface
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *HandlerConstructor) constructHandlerFromName(handlerName string) (frame.HandlerBaseInterface, error) {
	handlerConfPath, ok := h.handlersConfPath[handlerName]
	if !ok {
		return nil, fmt.Errorf("%s: handler's configuration is not set", handlerName)
	}
	handlerInterface, err := h.constructHandlerFromFile(handlerConfPath)
	if err != nil {
		return nil, err
	}

	if handlerInterface.Name() != handlerName {
		return nil, fmt.Errorf("handler Name (%s) mismatch with configuration file Name (%s) path (%s)", handlerInterface.Name(), handlerName, handlerConfPath)
	}
	return handlerInterface, err
}

func (h *HandlerConstructor) constructHandlerFromFile(filePath string) (frame.HandlerBaseInterface, error) {
	configure, errReadFile := os.ReadFile(filePath)
	if errReadFile != nil {
		return nil, errReadFile
	}
	handlerConf := make(map[any]any)
	errReadFile = yaml.Unmarshal([]byte(configure), &handlerConf)
	if errReadFile != nil {
		return nil, errReadFile
	}
	return h.constructHandler(handlerConf)
}

func (h *HandlerConstructor) constructHandler(handler_conf map[any]any) (frame.HandlerBaseInterface, error) {
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
				handlerInterface, ok := concrete.(frame.HandlerBaseInterface)
				if !ok {
					return nil, fmt.Errorf("convert %s to HandlerBaseInterface error", typeName)
				}

				return handlerInterface, nil
			}
		}
	}
	return nil, fmt.Errorf("handler configuration must have type %v", handler_conf)
}

func (h *HandlerConstructor) getHandlerName(filePath string) (string, error) {
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

func (h *HandlerConstructor) existHandler(name string) bool {
	_, exist := h.handlers[name]
	return exist
}

func (h *HandlerConstructor) LoadConfigFromFile(confPath string) error {
	return nil
}

func (h *HandlerConstructor) LoadConfigFromMemory(configure []byte) error {
	return nil
}

func (h *HandlerConstructor) LoadEnvironmentConf() error {
	h.handlers = make(map[string]frame.HandlerBaseInterface)
	h.handlersConfPath = make(map[string]string)
	return nil
}

// //////////////////////////////////////////////////////////////////////////////////
func (h *HandlerConstructor) Name() string {
	return reflect.TypeOf(*h).Name()
}

func (h *HandlerConstructor) SetConstructorInterface(constructorInterface any) {
	constructorInterfaceNew, ok := constructorInterface.(frame.ConstructorInterface)
	if !ok {
		panic("constructorInterface is not frame.ConstructorInterface")
	}
	h.constructorInterface = constructorInterfaceNew
}
