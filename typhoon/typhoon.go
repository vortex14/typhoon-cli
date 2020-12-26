package typhoon

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/go-cmd/cmd"
	"github.com/go-logfmt/logfmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"typhoon-cli/components"
	"typhoon-cli/interfaces"
	"typhoon-cli/utils"
)


type kv struct {
	k, v []byte
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

type Worker struct {
	Command 	  string
	Args    	  []string
	Cmd			  *cmd.Cmd
	Status 		  chan bool
	Name 		  string
	ComponentPath string
	ProjectPath	  string
	Active bool
}

type Components = struct {
	Components 			[]string
	ActiveComponents 	map[string] *Worker
	PathProject 		string
	TyphoonPath			string
	ConfigFile			string
}



func (w *Worker) Run(typhoonPath string) {
	color.Green("Run component %s", w.Name)
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	envCmd := cmd.NewCmdOptions(cmdOptions, w.Command, w.Args...)
	w.Cmd = envCmd
	w.Status = make(chan bool, 1)
	w.Status <- true
	projectEnv := fmt.Sprintf("PYTHONPATH=%s:%s", typhoonPath, w.ProjectPath)
	//color.Red("project path %s; projectEnv: %s", typhoonPath, projectEnv)
	newEnv := append(os.Environ(), projectEnv)
	envCmd.Env = newEnv
}



func CreateTransporterManifest(version string) error {

	//if "v1.1" == version {
	//	//dir := &components.Directory{Path: "project/result_transporter/consumers"}
	//	//dataDir := dir.GetDataFromDirectory("project/result_transporter/consumers")
	//	v1_1.CreateTransporterManifest()
	//	//for _, v := range dataDir {
	//	//	color.Red("k %s, v %s", v.Path, v.Type)
	//	//}
	//} else {
	//	color.Red("Version not found")
	//	return nil
	//}



	return nil
}

var watcher *fsnotify.Watcher

func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func WatchTest()  {
	color.Green("watch for project ..")
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk("project", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done

}


func Watch(wg *sync.WaitGroup, tcomponents *Components, configFile string)  {
	color.Green("watch for project ..")
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()



	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk("project", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:

				if strings.Contains(event.Name, ".pyc") || strings.Contains(event.String(), "CHMOD") {
					continue
				}

				componentChanged := "processor"

				for _, component := range tcomponents.Components {
					if strings.Contains(event.Name, component) {
						color.Yellow("reloading %s ... !", component)
						componentChanged = component
						break
					}

				}





				component := tcomponents.ActiveComponents[componentChanged]
				wg.Add(1)
				closeComponent(wg, component)


				initComponent(wg, tcomponents, componentChanged, configFile)

				// watch for errors
			case err := <-watcher.Errors:
				color.Red("ERROR---->", err)
			}
		}
	}()

	<-done

}


func Migrate(project interfaces.TyphoonProject) error {
	project.Migrate()
	return nil
}



func CreateSymbolicLink() error {
	project := &components.Project{}

	err := project.CreateSymbolicLink()

	return err
}

func ParseLogData(fileObject *interfaces.FileObject) error {
	u := utils.Utils{}
	err := u.ParseLog(fileObject)


	if err != nil {

		color.Red("Error %s", err)
		os.Exit(0)

	}

	return nil


}


func (w *Worker) Logging(wg *sync.WaitGroup) {

	Info := color.New(color.FgWhite, color.BgBlack, color.Bold).SprintFunc()

	defer wg.Done()
	for w.Cmd.Stdout != nil || w.Cmd.Stderr != nil || w.Status != nil {
		select {
		case line, open := <-w.Cmd.Stdout:
			if !open {
				continue
			}
			color.Cyan(line)
			fmt.Printf(`%s Logs ...
`, Info(w.Name))
			logDataMap := logfmt.NewDecoder(strings.NewReader(line))

			for logDataMap.ScanRecord() {
				for logDataMap.ScanKeyval() {



					if w.Name == "processor" {
						color.Yellow("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if w.Name == "result_transporter" {
						color.Green("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if w.Name == "fetcher" {
						color.Blue("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if w.Name == "donor" {
						color.HiBlackString("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if w.Name == "scheduler" {
						color.Cyan("%s = %s", logDataMap.Key(), logDataMap.Value())
					}



				}

			}
			if logDataMap.Err() != nil {
				color.Red("Invalid Log format. Don't use = . Broken line: %s",line)
				//panic(d.Err())
				continue
			}
			fmt.Printf(`
------------
`)




		case line, open := <-w.Cmd.Stderr:
			if !open {
				continue
			}
			errLog := ""
			io.Copy(os.Stderr, bytes.NewBufferString(errLog))
			//errLog = fmt.Sprintf("Component: %s; %s , %s", w.Name, errLog, line)
			//color.Red(errLog)

			err := w.Cmd.Stop()

			if err != nil {
				color.Red(" %s error: %s",w.Name, line)
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

		case status, ok := <-w.Status:
			if ok != true || status == false {

				err := w.Cmd.Stop()

				if err != nil {
					color.Red("Component: %s ,Err: %s",w.Name, err)
				}
				close(w.Status)

				//color.Red("Return from Logging. Component: %s", w.Name)
				status := w.Cmd.Status()
				errKill := syscall.Kill(-status.PID, syscall.SIGKILL)
				if errKill != nil {
					color.Red("Error kill :%s, component: %s", errKill, w.Name)
				}
				return
			}

		}

	}



}

func initComponent(wg *sync.WaitGroup, tcomponents *Components, component string, configFile string)  {
	pathExecute := fmt.Sprintf("%s.py", component)
	configArg := fmt.Sprintf("--config=%s", configFile)
	typhoonComponent := &Worker{Command: "python3.8", Args: []string{pathExecute, configArg}}
	typhoonComponent.Name = component
	typhoonComponent.ComponentPath = fmt.Sprintf("%s/project/%s", tcomponents.PathProject, component )

	typhoonComponent.ProjectPath = tcomponents.PathProject

	color.Red("ProjectPath: %s. file execute : %s", typhoonComponent.ProjectPath, pathExecute)
	typhoonComponent.Run(tcomponents.TyphoonPath)

	typhoonComponent.Cmd.Start()
	typhoonComponent.Cmd.Status()
	wg.Add(1)
	go typhoonComponent.Logging(wg)
	typhoonComponent.Active = true

	tcomponents.ActiveComponents[component] = typhoonComponent
}

func initComponents(wg *sync.WaitGroup, tcomponents *Components, configFile string)  {
	tcomponents.ActiveComponents = make(map[string]*Worker)
	defer wg.Done()

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

	//var wg sync.WaitGroup
	for _, component := range tcomponents.Components {
		initComponent(wg, tcomponents, component, configFile)

	}

}

func closeComponent(wg *sync.WaitGroup, component *Worker) {
	defer wg.Done()

	component.Status <- false
	component.Active = false


	color.Red("component %s was be closed", component.Name)
	//time.Sleep(time.Second * 2)

}

func closeComponents(wg *sync.WaitGroup, tcomponents *Components) {
	defer wg.Done()
	for _, component := range tcomponents.ActiveComponents {

		if component.Active {
			wg.Add(1)
			go closeComponent(wg, component)
		}


	}
}

type Task struct {
	closed chan struct{}
	wg     sync.WaitGroup
	ticker *time.Ticker
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

func Check(project interfaces.TyphoonProject)  {

	pathProject, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	isProject := project.CheckProject()
	if isProject == false {
		color.Red("Project does not exists in the current directory :%s", pathProject )
		os.Exit(1)
	}
}

func Run(project interfaces.TyphoonProject)  {


	task := &Task{
		closed: make(chan struct{}),
		ticker: time.NewTicker(time.Second * 2),
	}

	pathProject, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	var TyphoonPath string


	//configPath := fmt.Sprintf("%s/project/%s", pathProject, configFile)
	//if _, err := os.Stat(configPath); os.IsNotExist(err) {
	//	color.Red("Config %s does not exists in project :%s", configFile, configPath )
	//	os.Exit(1)
	//}



	isProject := project.CheckProject()
	if isProject == false {
		color.Red("Project does not exists in the current directory :%s", pathProject )
		os.Exit(1)
	}

	typhoonDir := &components.Directory{
		Path: "typhoon",
	}

	if typhoonDir.IsExistDir("typhoon") {
		TyphoonPath = "typhoon"
		goto toComponentInit
	}

	//Check TYPHOON_PATH
	for _, s := range os.Environ() {
		kv := strings.SplitN(s, "=", 2) // unpacks "key=value"
		if kv[0] == "TYPHOON_PATH" {
			TyphoonPath = kv[1]
		}
	}
	if len(TyphoonPath) == 0 {
		color.Red("Not found TYPHOON_PATH")
		os.Exit(1)
	}

	toComponentInit:
	//return

	var typhoonComponent = &Components{
		Components:  project.GetComponents(),
		PathProject: pathProject,
		TyphoonPath: TyphoonPath,
	}
	color.Magenta("start components")
	task.wg.Add(1)
	go initComponents(&task.wg, typhoonComponent, project.GetConfigFile())



	task.wg.Add(1)
	go task.Run()
	//
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go Watch(&task.wg, typhoonComponent, project.GetConfigFile())
	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		task.wg.Add(1)
		go closeComponents(&task.wg, typhoonComponent)

	}
	task.Stop()



}

func main() {


}
