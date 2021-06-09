package typhoon

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"typhoon-cli/src/environment"
	"typhoon-cli/src/interfaces"
	v11 "typhoon-cli/src/migrates/v1.1"
	"typhoon-cli/src/utils"
)

type components = struct {
	ActiveComponents 	map[string] *Component
	PathProject 		string
	TyphoonPath			string
	ConfigFile			string
}

type Task struct {
	closed chan struct{}
	wg     sync.WaitGroup
	ticker *time.Ticker
}


type Project struct {
	Path              string
	Name              string
	Tag 			  string
	LogLevel		  string
	SelectedComponent []string
	components        components
	ConfigFile        string
	AutoReload        bool
	Version           string
	BuilderOptions    *interfaces.BuilderOptions
	task              *Task
	EnvSettings       *environment.Settings
	Watcher           fsnotify.Watcher
	Config *ConfigProject
}


func watchDirTeet(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func (p *Project) GetComponentPort(name string) int {
	return p.Config.Config.GetComponentPort(name)
}

func (p *Project) WatchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return p.Watcher.Add(path)
	}

	return nil
}

func (p *Project) CreateProject() {
	color.Yellow("creating project...")
	u := utils.Utils{}
	fileObject := &interfaces.FileObject{
		Path: "../builders/v1.1/project",
	}

	err := u.CopyDir(p.Name, fileObject)


	if err != nil {

		color.Red("Error %s", err)
		os.Exit(0)

	}

	gitIgnore := &interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: ".gitignore",
	}
	errCopyIgnore := u.CopyFile(p.Name + "/.gitignore", gitIgnore)
	if errCopyIgnore != nil {
		color.Red("Error copy %s", err)
	}



	_, confT := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "config.goyaml",

	})
	goTemplate := interfaces.GoTemplate{
		Source: confT,
		ExportPath: p.Name +"/config.local.yaml",
		Data: map[string]string{
			"projectName": p.Name,
			"nsqdAdd": "localhost:4150",
			"redisHost": "localhost",
			"mongoHost": "localhost",
			"redisPort": "6379",
			"debug": "true",
		},
	}

	_= u.GoRunTemplate(&goTemplate)
	goTemplateCompose := interfaces.GoTemplate{
		Source: confT,
		ExportPath: p.Name +"/config.prod.yaml",
		Data: map[string]string{
			"projectName": p.Name,
			"nsqdAdd": "nsqd:4150",
			"redisHost": "redis",
			"redisPort": "6379",
		},
	}

	_= u.GoRunTemplate(&goTemplateCompose)
	//color.Green("Teplate status: %b", status)

	_, dataTDockerLocal := u.GetGoTemplate(&interfaces.FileObject{Path: "../builders/v1.1", Name: "docker-compose.local.goyaml"})

	dataConfig := map[string]string{
		"projectName": p.GetName(),
		"tag": p.GetTag(),
	}

	goTemplateComposeLocal := interfaces.GoTemplate{
		Source: dataTDockerLocal,
		ExportPath: p.Name +"/docker-compose.local.yaml",
		Data: dataConfig,
	}


	u.GoRunTemplate(&goTemplateComposeLocal)
	color.Green("Project %s created !", p.Name)

}

func (p *Project) BuildCIResources() {
	color.Green("Build CI Resources for %s !", p.Name)
	u := utils.Utils{}
	_, confCi := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: ".gitlab-ci.yml",

	})
	goTemplate := interfaces.GoTemplate{
		Source: confCi,
		ExportPath: ".gitlab-ci.yml",
	}

	_= u.GoRunTemplate(&goTemplate)

	_, dockerFile := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "Dockerfile",

	})
	goTemplateDocker := interfaces.GoTemplate{
		Source: dockerFile,
		ExportPath: "Dockerfile",
		Data: map[string]string{
			"TYPHOON_IMAGE": p.Version,
		},
	}

	_= u.GoRunTemplate(&goTemplateDocker)


	_, helmFile := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "helm-review-values.yml",

	})
	goTemplateHelmValues := interfaces.GoTemplate{
		Source: helmFile,
		ExportPath: "helm-review-values.yml",
	}

	_= u.GoRunTemplate(&goTemplateHelmValues)

	_, configFile := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "config-stage.goyaml",

	})
	goTemplateConfig := interfaces.GoTemplate{
		Source: configFile,
		ExportPath: "config.kube-stage.yaml",
		Data: map[string]string{
			"projectName": p.GetName(),
		},
	}

	_= u.GoRunTemplate(&goTemplateConfig)

}

func (p *Project) BuildHelmMinikubeResources()  {
	color.Yellow("build helm minikube resources ...")

	u := utils.Utils{}
	fileObject := &interfaces.FileObject{
		Path: "../builders/v1.1/helm/helm",
	}

	err := u.CopyDirAndReplaceLabel("helm", &interfaces.ReplaceLabel{Label: "{{PROJECT_NAME}}", Value: p.Name}, fileObject)


	if err != nil {

		color.Red("Error %s", err)
		os.Exit(0)

	}

	_, dataTDeployLocal := u.GetGoTemplate(&interfaces.FileObject{Path: "../builders/v1.1/helm", Name: "helm_deploy.gosh"})

	dataConfig := map[string]string{
		"projectName": p.GetName(),
	}

	goTemplateHelmDeployLocal := interfaces.GoTemplate{
		Source: dataTDeployLocal,
		ExportPath: "helm_deploy.sh",
		Data: dataConfig,
	}


	u.GoRunTemplate(&goTemplateHelmDeployLocal)

	_, dataTDumpLocal := u.GetGoTemplate(&interfaces.FileObject{Path: "../builders/v1.1/helm", Name: "helm_dump.gosh"})

	dataDumpConfig := map[string]string{
		"projectName": p.GetName(),
	}

	goTemplateHelmDumpLocal := interfaces.GoTemplate{
		Source: dataTDumpLocal,
		ExportPath: "helm_dump.sh",
		Data: dataDumpConfig,
	}


	u.GoRunTemplate(&goTemplateHelmDumpLocal)


	_, dataTDeleteLocal := u.GetGoTemplate(&interfaces.FileObject{Path: "../builders/v1.1/helm", Name: "helm_delete.gosh"})

	dataDeleteConfig := map[string]string{
		"projectName": p.GetName(),
	}

	goTemplateHelmDeleteLocal := interfaces.GoTemplate{
		Source: dataTDeleteLocal,
		ExportPath: "helm_delete.sh",
		Data: dataDeleteConfig,
	}


	u.GoRunTemplate(&goTemplateHelmDeleteLocal)


	if err := os.Chmod("helm_delete.sh", 0755); err != nil {
		color.Red("%s",err)
	}

	if err := os.Chmod("helm_deploy.sh", 0755); err != nil {
		color.Red("%s",err)
	}

	if err := os.Chmod("helm_dump.sh", 0755); err != nil {
		color.Red("%s",err)
	}

	_, confT := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "config.minikube.goyaml",

	})
	goTemplate := interfaces.GoTemplate{
		Source: confT,
		ExportPath: "config.minikube.yaml",
		Data: map[string]string{
			"projectName": p.Name,
		},
	}

	_= u.GoRunTemplate(&goTemplate)

	color.Green("Generated")


}

func (p *Project) GetEnvSettings() *environment.Settings {
	return p.EnvSettings
}

func (p *Project) AddPromise()  {
	p.task.wg.Add(1)
}
func (p *Project) PromiseDone()  {
	p.task.wg.Done()
}
func (p *Project) WaitPromises()  {
	p.task.wg.Wait()
}
func (p *Project) Run()  {
	p.CheckProject()
	p.task = &Task{
		closed: make(chan struct{}),
		ticker: time.NewTicker(time.Second * 2),
	}
	typhoonDir := &Directory{
		Path: "typhoon",
	}

	if !typhoonDir.IsExistDir("typhoon") {
		_ = p.CreateSymbolicLink()
	}


	color.Magenta("start components")
	p.AddPromise()
	go p.initComponents()
	//
	p.AddPromise()
	go p.task.Run()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go p.Watch()
	//go Watch(&task.wg, typhoonComponent, project.GetConfigFile())
	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		p.AddPromise()
		go p.Close()

	}
	p.task.Stop()

}

func (p *Project) Watch()  {
	color.Green("watch for project ..")
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk("project", watchDirTeet); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:

				if strings.Contains(event.Name, ".pyc") {
					continue
				}

				if strings.Contains(event.String(), "CHMOD") {
					continue
				}

				if strings.Contains(event.Name, ".py~") {
					continue
				}

				if strings.Contains(event.Name, "__pycache__") {
					continue
				}


				componentChanged := ""

				for _, component := range p.SelectedComponent {
					if strings.Contains(event.Name, component) {
						color.Yellow("reloading %s ... !", component)
						componentChanged = component
						break
					}

				}

				if _, ok := p.components.ActiveComponents[componentChanged]; ok {

					color.Yellow("Reload %s ...", componentChanged)
					color.Yellow("event %+v",event)
					component := p.components.ActiveComponents[componentChanged]

					//p.AddPromise()
					go component.Restart(p)


					// "example" is not in the map
				} else {
					color.Yellow("%s isn't running", componentChanged)
				}

				//


				//p.AddPromise()
				//go component.Restart(p)

				//go component.Restart(p)

				//initComponent(wg, tcomponents, componentChanged, configFile)

				// watch for errors
			case err := <-watcher.Errors:
				color.Red("ERROR---->", err)
			}
		}
	}()

	<-done
}

func (p *Project) Close()  {
	defer p.PromiseDone()
	for _, component := range p.components.ActiveComponents {

		if component.Active {
			p.AddPromise()
			go component.Close(p)
		}


	}



}

func (p *Project) Down() {
	p.LoadConfig()
	commandDropProject := fmt.Sprintf("kill -9 $(ps aux | grep \"%s\" | awk '{print $2}')", p.GetName())
	color.Red("Running: %s: ",commandDropProject)
	ctxP, cancelP := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelP()

	if err := exec.CommandContext(ctxP, "bash", "-c", commandDropProject).Run(); err != nil {
		color.Yellow("Project components killed!")
		// This will fail after 100 milliseconds. The 5 second sleep
		// will be interrupted.
	}

	commandDropTyphoon := fmt.Sprintf("kill -9 $(ps aux | grep \"%s\" | awk '{print $2}')", "typhoon")
	ctxT, cancelT := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelT()

	if err := exec.CommandContext(ctxT, "bash", "-c", commandDropTyphoon).Run(); err != nil {

	}
}


func (p *Project) GetBuilderOptions() *interfaces.BuilderOptions {
	return p.BuilderOptions
}

func (p *Project) GetTag() string {
	return p.Tag
}
func (p *Project) Migrate()  {

	color.Yellow("Migrate project to %s !", p.GetVersion())

	if p.Version == "v1.1" {
		prMigrates := v11.ProjectMigrate{
			Project: p,
			Dir: &interfaces.FileObject{
				Path: "../builders/v1.1",
			},
		}
		prMigrates.MigrateV11()
	}
}


func (p *Project) Build()  {
	color.Yellow("builder run... options %+v", p.BuilderOptions)
}


func (p *Project) initComponents()  {
	p.components.ActiveComponents = make(map[string]*Component)
	p.Name = p.Config.Config.ProjectName
	defer p.PromiseDone()

	fmt.Printf(`
												╭━┳━╭━╭━╮╮
												┃┈┈┈┣▅╋▅┫┃
												┃┈┃┈╰━╰━━━━━━╮
												╰┳╯┈┈┈┈┈┈┈┈┈◢▉◣
												╲┃┈┈┈┈┈┈┈┈┈┈▉▉▉
												╲┃┈┈┈┈┈┈┈┈┈┈◥▉◤
												╲┃┈┈┈┈╭━┳━━━━╯
												╲┣━━━━━━┫
	
	
	
`)

	for _, componentName := range p.SelectedComponent {
		component := &Component{
			Name: componentName,
		}

		p.components.ActiveComponents[componentName] = component


		component.Start(p)

	}

}



func (p *Project) GetVersion() string {
	return p.Version
}

func (p *Project) CreateSymbolicLink() error {
	env := &environment.Environment{}
	_, settings := env.GetSettings()

	linkTyphoonPath := fmt.Sprintf("%s/pytyphoon/typhoon", settings.Path)
	color.Yellow("TYPHOON_PATH=%s", settings.Path)
	err := os.Symlink(linkTyphoonPath, "typhoon")

	if err != nil{
		fmt.Printf("err %s",  err)
	}

	return nil
}

func (p *Project) GetName() string {
	projectName := p.Name
	if len(projectName) == 0 {
		projectName = p.Config.Config.ProjectName
	}
	return projectName
}

func (p *Project) GetComponents() []string {
	return p.SelectedComponent
}

func (p *Project) GetConfigFile() string {
	return p.ConfigFile
}

func (p *Project) GetProjectPath() string {
	pathProject, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return pathProject
}
func (p *Project) GetLogLevel() string {
	return p.LogLevel
}

func (p *Project) LoadConfig() {
	configPath := fmt.Sprintf("%s/%s", p.Path, p.ConfigFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		color.Red("Config %s does not exists in project :%s", p.ConfigFile, configPath )
		os.Exit(1)
	}

	var config ConfigProject
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("config.yaml err   #%v ", err)
		os.Exit(1)
	} else {
		err = yaml.Unmarshal(yamlFile, &config.Config)
		if err != nil {
			//log.Fatalf("Unmarshal: %v", err)
			color.Red("Config load error: %s", err )
			os.Exit(1)
		}

	}
	config.ConfigFile = configPath

	p.Config = &config
}

func (p *Project) CheckProject() {
	var status = true
	var statuses = make(map[string]bool)

	p.Path = p.GetProjectPath()

	for _, componentName := range p.SelectedComponent {
		component := &Component{

			Name: componentName,
		}
		color.Yellow("checking: %s...",componentName)

		componentStatus := component.CheckComponent()
		statuses[componentName] = componentStatus
	}

	for componentStatus, statusComponent := range statuses {
		if !statusComponent {
			status = false
		}
		color.Yellow("component %s is: %t", componentStatus, statusComponent)
	}

	p.LoadConfig()




	if status == false {
		color.Red("Project does not exists in the current directory :%s", p.Path )
		os.Exit(1)
	}


	env := &environment.Environment{}
	_, settings := env.GetSettings()

	if len(settings.Path) == 0 || len(settings.Projects) == 0 {
		color.Red("We need set valid environment variables like TYPHOON_PATH and TYPHOON_PROJECTS in %s", env.ProfilePath )
		os.Exit(1)
	}

	p.EnvSettings = settings


}

func (t *Task) Run() {
	t.wg.Done()
	for {
		select {
		case <-t.closed:
			return
		case <-t.ticker.C:
			handle()
		}
	}
}

func (t *Task) Stop() {
	color.Green("Stopping ...")
	close(t.closed)

	t.wg.Wait()
	color.Green("All components are closed")
}

func handle() {
	for i := 0; i < 5; i++ {
		//fmt.Print("#")
		time.Sleep(time.Millisecond * 200)
	}
}
