package interfaces

import (
	"typhoon-cli/environment"
	"typhoon-cli/migrates"
)

type GoTemplate struct {
	Source string
	ExportPath string
	Data interface{}
}

type FileObject struct {
	Type string
	Path string
	Name string
	FileSystem
}


type MapFileObjects map[string]*FileObject

type FileSystem interface {
	GetDataFromDirectory(path string) MapFileObjects
	IsExistDir (path string) bool
}


type Environment interface {
	Load()
	Set()
	Get()
	GetSettings() (error, *environment.Settings)
}


type TyphoonProject interface {
	CheckProject() bool
	GetName() string
	GetComponents() []string
	GetConfigFile() string
	CreateSymbolicLink() error
	GetVersion() string
	migrates.ProjectMigrate
}

type Utils interface {
	GoRunTemplate(goTemplate *GoTemplate) bool
	ParseLog(object *FileObject) error
	GetGoTemplate(object *FileObject) (error, string)
}

type TyphoonComponent interface {
	CheckDirectory(required []string, pathComponent string) bool
	CheckComponent(component string) bool
}





