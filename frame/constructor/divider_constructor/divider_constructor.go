package dividerconstructor

import (
	"fmt"
	"ghgroups/frame"
	"os"
	"path"

	"reflect"

	folder "git-codecommit.us-east-1.amazonaws.com/v1/repos/go-utils.git/folder"
	"gopkg.in/yaml.v2"
)

type DividerConstructor struct {
	frame.DividerConstructorInterface
	frame.ConstructorSetterInterface
	constructorInterface frame.ConstructorInterface
	dividers             map[string]frame.DividerBaseInterface
	dividersConfPath     map[string]string
}

func NewDividerConstructor(constructorInterface frame.ConstructorInterface) *DividerConstructor {
	return &DividerConstructor{
		constructorInterface: constructorInterface,
		dividers:             make(map[string]frame.DividerBaseInterface),
		dividersConfPath:     make(map[string]string),
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////
// DividerConstructorInterface

func (d *DividerConstructor) CreateDividerWithConfPath(confFilePath string) error {
	dividerName, errGetdividerName := d.getDividerName(confFilePath)
	if errGetdividerName != nil {
		return errGetdividerName
	}
	if path, ok := d.dividersConfPath[dividerName]; ok {
		return fmt.Errorf("divider name %s already exists,path is %s", dividerName, path)
	}
	d.dividersConfPath[dividerName] = confFilePath
	if d.existDivider(dividerName) {
		return fmt.Errorf("divider %s already exists", dividerName)
	}

	dividerInterface, err := d.constructDividerFromName(dividerName)
	if err != nil {
		return err
	}
	err = d.RegisterDivider(dividerName, dividerInterface)
	if err != nil {
		return err
	}
	return nil
}

func (d *DividerConstructor) ParseDividerConfFolder(confFolderPath string) error {
	files, err := folder.OSReadDir(confFolderPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		dividerName, errGetdividerName := d.getDividerName(f)
		if errGetdividerName != nil {
			return errGetdividerName
		}
		if path, ok := d.dividersConfPath[dividerName]; ok {
			return fmt.Errorf("divider name %s already exists,path is %s", dividerName, path)
		}
		d.dividersConfPath[dividerName] = f
	}

	for dividerName := range d.dividersConfPath {
		if d.existDivider(dividerName) {
			return fmt.Errorf("divider %s already exists", dividerName)
		}

		dividerInterface, err := d.constructDividerFromName(dividerName)
		if err != nil {
			return err
		}
		err = d.RegisterDivider(dividerName, dividerInterface)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DividerConstructor) GetDivider(name string) (frame.DividerBaseInterface, error) {
	dividerInterface, exist := d.dividers[name]
	if exist {
		return dividerInterface, nil
	}
	dividerInterface, err := d.constructDividerFromName(name)
	if err != nil {
		return nil, err
	}
	err = d.RegisterDivider(name, dividerInterface)
	if err != nil {
		return nil, err
	}
	return dividerInterface, nil
}

func (d *DividerConstructor) RegisterDivider(name string, divider_interface frame.DividerBaseInterface) error {
	if d.existDivider(name) {
		return fmt.Errorf("divider %s is exist", name)
	}
	d.dividers[name] = divider_interface
	return nil
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DividerConstructor) constructDividerFromName(dividerName string) (frame.DividerBaseInterface, error) {
	dividerConfPath, ok := d.dividersConfPath[dividerName]
	if !ok {
		return nil, fmt.Errorf("%s: divider's configuration is not set", dividerName)
	}
	dividerInterface, err := d.constructDividerFromFile(dividerConfPath)
	if err != nil {
		return nil, err
	}

	if dividerInterface.Name() != dividerName {
		return nil, fmt.Errorf("divider Name (%s) mismatch with configuration file Name (%s) path (%s)", dividerInterface.Name(), dividerName, dividerConfPath)
	}
	return dividerInterface, err
}

func (d *DividerConstructor) constructDividerFromFile(filePath string) (frame.DividerBaseInterface, error) {
	configure, errReadFile := os.ReadFile(filePath)
	if errReadFile != nil {
		return nil, errReadFile
	}
	dividerConf := make(map[any]any)
	errReadFile = yaml.Unmarshal([]byte(configure), &dividerConf)
	if errReadFile != nil {
		return nil, errReadFile
	}
	return d.constructDivider(dividerConf)
}

func (d *DividerConstructor) constructDivider(divider_conf map[any]any) (frame.DividerBaseInterface, error) {
	for k, v := range divider_conf {
		switch k.(type) {
		case string:
			if k == "type" {
				originConf, _ := yaml.Marshal(divider_conf)
				typeName := v.(string)
				concrete, err := d.constructorInterface.Create(typeName, originConf, d.constructorInterface)
				if err != nil {
					return nil, err
				}
				dividerInterface, ok := concrete.(frame.DividerBaseInterface)
				if !ok {
					return nil, fmt.Errorf("convert %s to DividerBaseInterface error", typeName)
				}

				return dividerInterface, nil
			}
		}
	}
	return nil, fmt.Errorf("divider configuration must have type %v", divider_conf)
}

func (d *DividerConstructor) getDividerName(filePath string) (string, error) {
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

func (d *DividerConstructor) existDivider(name string) bool {
	_, exist := d.dividers[name]
	return exist
}

func (d *DividerConstructor) LoadConfigFromFile(confPath string) error {
	return nil
}

func (d *DividerConstructor) LoadConfigFromMemory(configure []byte) error {
	return nil
}

func (d *DividerConstructor) LoadEnvironmentConf() error {
	d.dividers = make(map[string]frame.DividerBaseInterface)
	d.dividersConfPath = make(map[string]string)
	return nil
}

// //////////////////////////////////////////////////////////////////////////////////
func (d *DividerConstructor) Name() string {
	return reflect.TypeOf(*d).Name()
}

func (d *DividerConstructor) SetConstructorInterface(constructorInterface any) {
	constructorInterfaceNew, ok := constructorInterface.(frame.ConstructorInterface)
	if !ok {
		panic("constructorInterface is not frame.ConstructorInterface")
	}
	d.constructorInterface = constructorInterfaceNew
}
