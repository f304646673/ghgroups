package aynchandlergroupconstructor

import (
	"fmt"
	"ghgroups/frame"
	"os"
	"path"
	"reflect"

	folder "git-codecommit.us-east-1.amazonaws.com/v1/repos/go-utils.git/folder"

	"gopkg.in/yaml.v2"
)

type AsyncHandlerGroupConstructor struct {
	frame.AsyncHandlerGroupConstructorInterface
	frame.ConstructorSetterInterface
	constructorInterface frame.ConstructorInterface
	handlers             map[string]frame.AsyncHandlerGroupBaseInterface
	handlersConfPath     map[string]string
}

func NewAsyncHandlerGroupConstructor(constructorInterface frame.ConstructorInterface) *AsyncHandlerGroupConstructor {
	return &AsyncHandlerGroupConstructor{
		constructorInterface: constructorInterface,
		handlers:             make(map[string]frame.AsyncHandlerGroupBaseInterface),
		handlersConfPath:     make(map[string]string),
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////
// AsyncHandlerGroupConstructorInterface

func (a *AsyncHandlerGroupConstructor) CreateAsyncHandlerGroupWithConfPath(confFilePath string) error {
	handlerName, errGethandlerName := a.getAsyncHandlerGroupName(confFilePath)
	if errGethandlerName != nil {
		return errGethandlerName
	}
	if path, ok := a.handlersConfPath[handlerName]; ok {
		return fmt.Errorf("handler name %s already exists,path is %s", handlerName, path)
	}
	a.handlersConfPath[handlerName] = confFilePath
	if a.existAsyncHandlerGroup(handlerName) {
		return fmt.Errorf("handler %s already exists", handlerName)
	}

	handlerInterface, err := a.constructAsyncHandlerGroupFromName(handlerName)
	if err != nil {
		return err
	}
	err = a.RegisterAsyncHandlerGroup(handlerName, handlerInterface)
	if err != nil {
		return err
	}
	return nil
}

func (a *AsyncHandlerGroupConstructor) ParseAsyncHandlerGroupConfFolder(confFolderPath string) error {
	files, err := folder.OSReadDir(confFolderPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		handlerName, errGethandlerName := a.getAsyncHandlerGroupName(f)
		if errGethandlerName != nil {
			return errGethandlerName
		}
		if path, ok := a.handlersConfPath[handlerName]; ok {
			return fmt.Errorf("handler name %s already exists,path is %s", handlerName, path)
		}
		a.handlersConfPath[handlerName] = f
	}

	for handlerName := range a.handlersConfPath {
		if a.existAsyncHandlerGroup(handlerName) {
			return fmt.Errorf("handler %s already exists", handlerName)
		}

		handlerInterface, err := a.constructAsyncHandlerGroupFromName(handlerName)
		if err != nil {
			return err
		}
		err = a.RegisterAsyncHandlerGroup(handlerName, handlerInterface)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AsyncHandlerGroupConstructor) GetAsyncHandlerGroup(name string) (frame.AsyncHandlerGroupBaseInterface, error) {
	handlerInterface, exist := a.handlers[name]
	if exist {
		return handlerInterface, nil
	}
	handlerInterface, err := a.constructAsyncHandlerGroupFromName(name)
	if err != nil {
		return nil, err
	}
	err = a.RegisterAsyncHandlerGroup(name, handlerInterface)
	if err != nil {
		return nil, err
	}
	return handlerInterface, nil
}

func (a *AsyncHandlerGroupConstructor) RegisterAsyncHandlerGroup(name string, handler_interface frame.AsyncHandlerGroupBaseInterface) error {
	if a.existAsyncHandlerGroup(name) {
		return fmt.Errorf("handler %s is exist", name)
	}
	a.handlers[name] = handler_interface
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

func (a *AsyncHandlerGroupConstructor) constructAsyncHandlerGroupFromName(handlerName string) (frame.AsyncHandlerGroupBaseInterface, error) {
	handlerConfPath, ok := a.handlersConfPath[handlerName]
	if !ok {
		return nil, fmt.Errorf("%s: handler's configuration is not set", handlerName)
	}
	handlerInterface, err := a.constructAsyncHandlerGroupFromFile(handlerConfPath)
	if err != nil {
		return nil, err
	}

	if handlerInterface.Name() != handlerName {
		return nil, fmt.Errorf("handler Name (%s) mismatch with configuration file Name (%s) path (%s)", handlerInterface.Name(), handlerName, handlerConfPath)
	}
	return handlerInterface, err
}

func (a *AsyncHandlerGroupConstructor) constructAsyncHandlerGroupFromFile(filePath string) (frame.AsyncHandlerGroupBaseInterface, error) {
	configure, errReadFile := os.ReadFile(filePath)
	if errReadFile != nil {
		return nil, errReadFile
	}
	handlerConf := make(map[any]any)
	errReadFile = yaml.Unmarshal([]byte(configure), &handlerConf)
	if errReadFile != nil {
		return nil, errReadFile
	}
	return a.constructAsyncHandlerGroup(handlerConf)
}

func (a *AsyncHandlerGroupConstructor) constructAsyncHandlerGroup(handler_conf map[any]any) (frame.AsyncHandlerGroupBaseInterface, error) {
	for k, v := range handler_conf {
		switch k.(type) {
		case string:
			if k == "type" {
				originConf, _ := yaml.Marshal(handler_conf)
				typeName := v.(string)
				concrete, err := a.constructorInterface.Create(typeName, originConf, a.constructorInterface)
				if err != nil {
					return nil, err
				}
				handlerInterface, ok := concrete.(frame.AsyncHandlerGroupBaseInterface)
				if !ok {
					return nil, fmt.Errorf("convert %s to AsyncHandlerGroupBaseInterface error", typeName)
				}

				return handlerInterface, nil
			}
		}
	}
	return nil, fmt.Errorf("handler configuration must have type %v", handler_conf)
}

func (a *AsyncHandlerGroupConstructor) getAsyncHandlerGroupName(filePath string) (string, error) {
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

func (a *AsyncHandlerGroupConstructor) existAsyncHandlerGroup(name string) bool {
	_, exist := a.handlers[name]
	return exist
}

func (a *AsyncHandlerGroupConstructor) LoadConfigFromFile(confPath string) error {
	return nil
}

func (a *AsyncHandlerGroupConstructor) LoadConfigFromMemory(configure []byte) error {
	return nil
}

func (a *AsyncHandlerGroupConstructor) LoadEnvironmentConf() error {
	a.handlers = make(map[string]frame.AsyncHandlerGroupBaseInterface)
	a.handlersConfPath = make(map[string]string)
	return nil
}

// //////////////////////////////////////////////////////////////////////////////////
func (a *AsyncHandlerGroupConstructor) Name() string {
	return reflect.TypeOf(*a).Name()
}

func (a *AsyncHandlerGroupConstructor) SetConstructorInterface(constructorInterface any) {
	constructorInterfaceNew, ok := constructorInterface.(frame.ConstructorInterface)
	if !ok {
		panic("constructorInterface is not frame.ConstructorInterface")
	}
	a.constructorInterface = constructorInterfaceNew
}
