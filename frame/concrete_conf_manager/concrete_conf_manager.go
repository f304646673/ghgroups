package concreteconfmanager

import (
	"fmt"

	"git-codecommit.us-east-1.amazonaws.com/v1/repos/go-utils.git/folder"
)

type ConcreteConfManager struct {
	confNameFilePathMap map[string]string
}

func NewConcreteConfManager() *ConcreteConfManager {
	return &ConcreteConfManager{
		confNameFilePathMap: make(map[string]string),
	}
}

func (c *ConcreteConfManager) ParseConfFolder(confFolderPath string) error {
	files, err := folder.GetAllFilesPath(confFolderPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		confName, errGetConfName := c.getConfName(f)
		if errGetConfName != nil {
			return errGetConfName
		}
		if path, ok := c.confNameFilePathMap[confName]; ok {
			return fmt.Errorf("conf name %s already exists,path is %s", confName, path)
		}
		c.confNameFilePathMap[confName] = f
	}

	return nil
}

func (c *ConcreteConfManager) GetConfPath(confName string) (string, error) {
	path, exist := c.confNameFilePathMap[confName]
	if !exist {
		return "", fmt.Errorf("conf %s not exist", confName)
	}
	return path, nil
}

func (c *ConcreteConfManager) getConfName(path string) (string, error) {
	return folder.GetFileNameWithoutSuffix(path), nil
}
