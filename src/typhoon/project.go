package typhoon

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"typhoon-cli/src/environment"
	"typhoon-cli/src/interfaces"
	v1_1 "typhoon-cli/src/migrates/v1.1"
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
	SelectedComponent []string
	components        components
	ConfigFile        string
	AutoReload        bool
	Version           string
	BuilderOptions    *interfaces.BuilderOptions
	task              *Task
	EnvSettings       *environment.Settings
	Watcher           fsnotify.Watcher
}


func watchDirTeet(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func (p *Project) WatchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return p.Watcher.Add(path)
	}

	return nil
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


				componentChanged := "processor"

				for _, component := range p.SelectedComponent {
					if strings.Contains(event.Name, component) {
						color.Yellow("reloading %s ... !", component)
						componentChanged = component
						break
					}

				}

				color.Yellow("Reload %s ...", componentChanged)
				color.Yellow("event %+v",event)




				component := p.components.ActiveComponents[componentChanged]


				//p.AddPromise()
				go component.Restart(p)

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


func (p *Project) GetBuilderOptions() *interfaces.BuilderOptions {
	return p.BuilderOptions
}

func (p *Project) Migrate()  {

	color.Yellow("Migrate project to %s !", p.GetVersion())

	if p.Version == "v1.1" {
		prMigrates := v1_1.ProjectMigrate{
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

		component.Start(p)

		p.components.ActiveComponents[componentName] = component

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
	return p.Name
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

	configPath := fmt.Sprintf("%s/%s", p.Path, p.ConfigFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		color.Red("Config %s does not exists in project :%s", p.ConfigFile, configPath )
		os.Exit(1)
	}



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
