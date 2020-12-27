package interfaces

import (
	"typhoon-cli/src/builders"
	"typhoon-cli/src/environment"
	"typhoon-cli/src/migrates"
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
type BuilderOptions builders.BuildOptions

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

type goPromise interface {
	AddPromise()
	PromiseDone()
	WaitPromises()
}

type Project interface {
	Run()
	Close()
	Watch()
	CheckProject()
	GetName() string
	GetProjectPath() string
	GetComponents() []string
	GetConfigFile() string
	CreateSymbolicLink() error
	GetVersion() string
	GetEnvSettings() *environment.Settings
	GetBuilderOptions() *BuilderOptions
	migrates.ProjectMigrate
	builders.ProjectBuilder
	goPromise
}

type Utils interface {
	GoRunTemplate(goTemplate *GoTemplate) bool
	ParseLog(object *FileObject) error
	GetGoTemplate(object *FileObject) (error, string)
}

type Component interface {
	CheckDirectory(required []string, pathComponent string) bool
	GetName() string
	//CheckComponent(component string) bool
	Start(project Project)
	Close(project Project)
	Stop()
	Restart(project Project)
	goPromise
}





