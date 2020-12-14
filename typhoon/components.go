package typhoon

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
)


type FileObject struct {
	Type string
	Path string
}

func checkDirectory(required []string, pathComponent string) bool {
	var status = true

	color.Yellow("checking %s", pathComponent)

	dataDir := GetDataFromDirectory(pathComponent)
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

func checkComponent(component string) bool {
	var status = false

	pathComponent := fmt.Sprintf("project/%s",component)
	if _, err := os.Stat(pathComponent); !os.IsNotExist(err) {

		if component == "fetcher" {
			required := []string{"run.py", "executions", "responses", "__init__.py"}
			status = checkDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Fetcher dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if component == "processor" {

			required := []string{"run.py", "executable", "__init__.py"}
			status = checkDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Processor dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if component == "scheduler" {
			required := []string{"run.py", "__init__.py"}
			status = checkDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Scheduler dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if component == "donor" {
			required := []string{"run.py", "__init__.py", "v1", "routes.py"}
			status = checkDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Scheduler dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if component == "result_transporter" {
			required := []string{"run.py", "__init__.py", "consumers"}
			status = checkDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Scheduler dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}
		}
	} else {
		color.Red("path %s doesn't exist", component)
	}


	return status
}

func checkProject(components []string) bool {
	var status = true
	var statuses = make(map[string]bool)

	for _, component := range components {
		color.Yellow("checking: %s...",component)
		componentStatus := checkComponent(component)
		statuses[component] = componentStatus
	}

	for componentStatus, statusComponent := range statuses {
		if !statusComponent {
			status = false
		}
		color.Yellow("component %s is: %t", componentStatus, statusComponent)
	}


	return status
}



func GetDataFromDirectory(path string) map[string]*FileObject {
	currentData := make(map[string]*FileObject, 0)


	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		typeFile := "file"
		if file.IsDir() {
			typeFile = "dir"
		}

		currentData[file.Name()] = &FileObject{

			Type: typeFile,
			Path: file.Name(),
		}

	}



	return currentData
}
