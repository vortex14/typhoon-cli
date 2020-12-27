package typhoon

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-cmd/cmd"
	"github.com/go-logfmt/logfmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"syscall"
	"typhoon-cli/src/interfaces"
)

type Directory struct {
	Path string
}

type Worker struct {
	Command 	  string
	Args    	  []string
	Cmd			  *cmd.Cmd
	Status 		  chan bool
}

type Component struct {
	Path string
	Name string
	Active bool
	Worker *Worker
	Promise sync.WaitGroup
}

func (c *Component) AddPromise()  {
	c.Promise.Add(1)
}

func (c *Component) PromiseDone()  {
	c.Promise.Done()
}

func (c *Component) WaitPromises()  {
	c.Promise.Wait()
}

func (c *Component) Start(project interfaces.Project)  {
	color.Yellow("init %s", c.Name)

	pathExecute := fmt.Sprintf("%s.py", c.Name)
	configArg := fmt.Sprintf("--config=%s", project.GetConfigFile())
	c.Worker = &Worker{Command: "python3.8", Args: []string{pathExecute, configArg}}
	//c.Path = fmt.Sprintf("%s/project/%s", project.GetProjectPath(), c.Name )
	c.Worker.Run(project)

	color.Yellow("path Ex: %s, arg: %s, path: %s", pathExecute, configArg, c.Path)

	c.Worker.Cmd.Start()
	c.Worker.Cmd.Status()
	c.Active = true
	c.AddPromise()
	go c.Logging()
}


func (c *Component) Close(project interfaces.Project)  {
	defer project.PromiseDone()
	c.Stop()
}

func (c *Component) Stop()  {
	status := c.Worker.Cmd.Status()
	color.Green("%s status.PID %s", status.PID, c.Name)
	if !IsClosed(c.Worker.Status){
		c.Worker.Status <- false
	}
	c.Active = false


	errKill := syscall.Kill(-status.PID, syscall.SIGKILL)
	if errKill != nil {
		color.Red("Error kill :%s, component: %s", errKill, c.Name)
	} else {
		color.Green("%s killed", c.Name)
	}


	color.Red("component %s was be closed", c.Name)

}

func (c Component) Restart(project *Project)  {
	color.Red("Restart component %s ...", c.Name)
	c.Stop()
	c.Start(project)

	project.components.ActiveComponents[c.Name] = &c
}

func (c *Component) GetName() string {
	return c.Name
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

func (w *Worker) Run(project interfaces.Project) {
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	envCmd := cmd.NewCmdOptions(cmdOptions, w.Command, w.Args...)
	w.Cmd = envCmd
	w.Status = make(chan bool, 1)
	w.Status <- true
	projectEnv := fmt.Sprintf("PYTHONPATH=%s:%s", project.GetEnvSettings(), project.GetProjectPath())
	newEnv := append(os.Environ(), projectEnv)
	envCmd.Env = newEnv
}

func IsClosed(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func (c *Component) Logging()  {
	Info := color.New(color.FgWhite, color.BgBlack, color.Bold).SprintFunc()
	for c.Worker.Cmd.Stdout != nil || c.Worker.Cmd.Stderr != nil || c.Worker.Status != nil {
		select {
		case line, open := <-c.Worker.Cmd.Stdout:
			if !open {
				continue
			}



			color.Cyan(line)
			fmt.Printf(`%s Logs ...
`, Info(c.Name))
			logDataMap := logfmt.NewDecoder(strings.NewReader(line))

			for logDataMap.ScanRecord() {
				for logDataMap.ScanKeyval() {



					if c.Name == "processor" {
						color.Yellow("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if c.Name == "result_transporter" {
						color.Green("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if c.Name == "fetcher" {
						color.Blue("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if c.Name == "donor" {
						color.HiBlackString("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if c.Name == "scheduler" {
						color.Cyan("%s = %s", logDataMap.Key(), logDataMap.Value())
					}



				}

			}
			if logDataMap.Err() != nil {
				//color.Red("Invalid Log format. Don't use = . Broken line: %s",line)
				//panic(d.Err())
				continue
			}
			fmt.Printf(`
------------
`)
		case line, open := <-c.Worker.Cmd.Stderr:
			if !open {
				continue
			}
			errLog := ""
			io.Copy(os.Stderr, bytes.NewBufferString(errLog))
			//errLog = fmt.Sprintf("Component: %s; %s , %s", w.Name, errLog, line)
			//color.Red(errLog)

			err := c.Worker.Cmd.Stop()

			if err != nil {
				color.Red(" %s error: %s",c.Name, line)
				//fmt.Fprintln(os.Stderr, line)
			}
			//close(w.Status)

			//color.Red("Return from Logging. Component: %s", w.Name)
			//status := w.Cmd.Status()
			//errKill := syscall.Kill(-status.PID, syscall.SIGKILL)
			//if errKill != nil {
			//	color.Red("Error kill :%s, component: %s", errKill, w.Name)
			//}
			continue
		case status, ok := <-c.Worker.Status:
			if ok != true || status == false {

				err := c.Worker.Cmd.Stop()

				if err != nil {
					color.Red("Component: %s ,Err: %s",c.Name, err)
				}



				if !IsClosed(c.Worker.Status) {
					close(c.Worker.Status)
				}

				//c.Promise.Done()

				color.Blue("promise done ... %s", c.Name)

				return
			}

		}

	}
}