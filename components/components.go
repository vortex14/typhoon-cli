package components

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"typhoon-cli/environment"
	"typhoon-cli/interfaces"
	v1_1 "typhoon-cli/migrates/v1.1"
)

type Directory struct {
	Path string
}

type Component struct {
	Path string
	Name string
}

type Project struct {
	Path string
	Name string
	Components []string
	ConfigFile string
	AutoReload bool
	Version string
}

func (p *Project) Migrate()  {

	color.Yellow("Migrate project to %s !", p.GetVersion())

	if p.Version == "v1.1" {
		prMigrates := v1_1.ProjectMigrate{
			Project: p,
			Dir: &interfaces.FileObject{
				Path: "../migrates/v1.1",
			},
		}
		prMigrates.MigrateV11()
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
	return p.Components
}

func (p *Project) GetConfigFile() string {
	return p.ConfigFile
}

func (d *Directory) GetDataFromDirectory(path string) interfaces.MapFileObjects {
	currentData := make(interfaces.MapFileObjects, 0)


	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		typeFile := "file"
		if file.IsDir() {
			typeFile = "dir"
		}

		currentData[file.Name()] = &interfaces.FileObject{

			Type: typeFile,
			Path: file.Name(),
		}

	}



	return currentData
}

func (c *Component) CheckDirectory(required []string, pathComponent string) bool  {
	var status = true

	color.Yellow("checking %s", pathComponent)
	dir := &Directory{
		Path: pathComponent,
	}
	dataDir := dir.GetDataFromDirectory(dir.Path)
	for _, reqFile := range required {
		if _, ok := dataDir[reqFile]; !ok {
			color.Red("%s not exist in %s", reqFile, pathComponent)
			status = false
		} else {
			color.Green("checked %s in %s", reqFile, pathComponent)
		}


	}


	return status
}

func (d *Directory) IsExistDir(path string) bool  {
	var status = false
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		status = true
	}

	return status
}

func (c *Component) CheckComponent() bool {
	var status = false

	pathComponent := fmt.Sprintf("project/%s",c.Name)



	if _, err := os.Stat(pathComponent); !os.IsNotExist(err) {

		if c.Name == "fetcher" {
			required := []string{"executions", "responses", "__init__.py"}

			status = c.CheckDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Fetcher dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if c.Name == "processor" {

			required := []string{"executable", "__init__.py"}
			status = c.CheckDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Processor dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if c.Name == "scheduler" {
			required := []string{"__init__.py"}
			status = c.CheckDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Scheduler dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if c.Name == "donor" {
			required := []string{"__init__.py", "v1", "routes.py"}
			status = c.CheckDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Scheduler dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if c.Name == "result_transporter" {
			required := []string{"__init__.py", "consumers"}
			status = c.CheckDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Scheduler dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}
		}
	} else {
		color.Red("path %s doesn't exist", c.Name)
	}

	fileName := fmt.Sprintf("%s.py", c.Name)
	required := []string{fileName}
	status = c.CheckDirectory(required, ".")
	logVal := fmt.Sprintf("%s.py is %t", c.Name, status)

	if status == true {
		color.Green(logVal)
	} else {
		color.Red(logVal)
	}


	return status
}

func (p *Project) CheckProject() bool {
	var status = true
	var statuses = make(map[string]bool)

	for _, componentName := range p.Components {
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


	return status
}


