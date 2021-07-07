package interfaces

import (
	"typhoon-cli/src/builders"
	"typhoon-cli/src/environment"
	"typhoon-cli/src/migrates"
	"typhoon-cli/src/typhoon/config"
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
	Data string
	FileSystem
}

type ClusterProject struct {
	Name string
	path string
	Branch string
	Git string
	Remote string
	GitlabId int
}

type GitlabProject struct {
	Name string
	Git string
	Id int
}

type GitlabServer interface {
	GetAllProjectsList() []*GitlabProject
	SyncGitlabProjects()
	Deploy()
	HistoryPipelines()
}

type Cluster interface {
	Add()
	Show()
	Create()
	Deploy()
	SaveConfig()
	GetName() string
	GetConfigName() string
	GetClusterConfigPath() string
	GetProjects() [] *ClusterProject
	GetMeta() map[string] interface{}
	GetEnvSettings() *environment.Settings
}

type ReplaceLabel struct {
	Label string
	Value string
}

type ReplaceLabels []ReplaceLabel


type MapFileObjects map[string]*FileObject
type BuilderOptions builders.BuildOptions

type FileSystem interface {
	GetDataFromDirectory(path string) MapFileObjects
	IsExistDir (path string) bool
}

type TestData interface {
	GetFields()
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


type Service interface {
	TestConnect() bool
}

type Services struct {}

type GrafanaInterface interface {
	ImportGrafanaConfig()
	RemoveGrafanaDashboard()
	CreateBaseGrafanaConfig()
	CreateGrafanaMonitoringTemplates()
}


type DockerInterface interface {
	BuildImage()
	ListContainers()
	ProjectBuild()
	RunComponent(component string) error
}

type HelmInterface interface {
	BuildHelmMinikubeResources()
	RemoveHelmMinikubeManifests()
}

type GitlabInterface interface {
	BuildCIResources()
}

type Project interface {
	Run()
	Close()
	Watch()
	CheckProject()
	GetTag() string
	GetName() string
	GetVersion() string
	GetLogLevel() string
	GetConfigFile() string
	GetProjectPath() string
	builders.ProjectBuilder
	migrates.ProjectMigrate
	GetComponents() []string
	CreateSymbolicLink() error
	GetDockerImageName() string
	GetComponentPort(name string) int
	LoadConfig() *config.Project
	GetBuilderOptions() *BuilderOptions
	GetEnvSettings() *environment.Settings
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
	Stop(project Project)
	//Restart(project *Project)
	goPromise
}





