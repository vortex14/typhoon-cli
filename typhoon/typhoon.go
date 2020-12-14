package typhoon

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-cmd/cmd"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)



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
}

type TyphoonComponents = struct {
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
	w.Cmd.Dir = w.ComponentPath
	w.Status = make(chan bool, 1)
	w.Status <- true
	//color.Red()
	projectEnv := fmt.Sprintf("PYTHONPATH=%s:%s", typhoonPath, w.ProjectPath)
	//color.Red(projectEnv)
	//projectEnv := "PYTHONPATH=/Users/dmitrijviharev/Desktop/projects/typhoon_projects/typhoon-local/typhoon-lite-log/pytyphoon:/Users/dmitrijviharev/Desktop/projects/typhoon_projects/typhoon-kube-project"
	newEnv := append(os.Environ(), projectEnv)
	envCmd.Env = newEnv
	//color.Green("Setup %s", w.Name)
}

func (w *Worker) Logging(wg *sync.WaitGroup) {

	Info := color.New(color.FgWhite, color.BgBlack, color.Bold).SprintFunc()
	//levelInfo := color.New(color.FgWhite, color.BgHiCyan, color.Bold).SprintFunc()
	//pathInfo := color.New(color.FgWhite, color.BgBlack, color.Bold).SprintFunc()
	//fmt.Printf("This %s rocks!\n", timeInfo("package"))
	//logColorsMap := make(map[string] interface{})

	defer wg.Done()
	for w.Cmd.Stdout != nil || w.Cmd.Stderr != nil || w.Status != nil {
		select {
		case line, open := <-w.Cmd.Stdout:
			if !open {
				continue
			}
			logData := string(line)
			logArr := strings.Split(logData, " ")


			if len(logArr) > 1 {

				logDataMap := make(map[string]string)
				for _, logRaw := range logArr {
					logDetail := strings.Split(logRaw, "=")
					if len(logDetail) != 2 {
						continue
					}

					logDataMap[logDetail[0]] = logDetail[1]

				}

				//logFormat := `
				//
				//event_time=%s level=%s
				//
				//`
				//logValues := make([]interface{}, 0)
				//
				//
				fmt.Printf(`%s Logs ...
`, Info(w.Name))
				for key, value := range logDataMap {

					if w.Name == "processor" {
						color.Yellow("%s = %s", key, value)
					} else if w.Name == "result_transporter" {
						color.Green("%s = %s", key, value)
					} else if w.Name == "fetcher" {
						color.Magenta("%s = %s", key, value)
					} else if w.Name == "donor" {
						color.Cyan("%s = %s", key, value)
					} else if w.Name == "scheduler" {
						color.Blue("%s = %s", key, value)
					}

					//if key == "event_time" {
					//	logValues = append(logValues, timeInfo(value))
					//	continue
					//}
					//if key  == "level" {
					//	logValues = append(logValues, levelInfo(value))
					//	continue
					//}

				}
				fmt.Printf(`
------------
`)
				//
				//fmt.Printf(logFormat, logValues...)

			}



		case line, open := <-w.Cmd.Stderr:
			if !open {
				continue
			}
			errLog := ""
			io.Copy(os.Stderr, bytes.NewBufferString(errLog))
			errLog = fmt.Sprintf("Component: %s; %s , %s", w.Name, errLog, line)
			color.Red(errLog)
			//fmt.Fprintln(os.Stderr, line)
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

func initComponents(wg *sync.WaitGroup, tcomponents *TyphoonComponents, configFile string)  {
	tcomponents.ActiveComponents = make(map[string]*Worker)
	defer wg.Done()

	//var wg sync.WaitGroup
	for _, component := range tcomponents.Components {
		pathExecute := fmt.Sprintf("run.py")
		configArg := fmt.Sprintf("--config=%s", configFile)
		typhoonComponent := &Worker{Command: "python3.8", Args: []string{pathExecute, configArg}}
		typhoonComponent.Name = component
		typhoonComponent.ComponentPath = fmt.Sprintf("%s/project/%s", tcomponents.PathProject, component )
		typhoonComponent.ProjectPath = tcomponents.PathProject
		typhoonComponent.Run(tcomponents.TyphoonPath)

		typhoonComponent.Cmd.Start()
		typhoonComponent.Cmd.Status()
		wg.Add(1)
		go typhoonComponent.Logging(wg)

		tcomponents.ActiveComponents[component] = typhoonComponent

	}

	//color.Yellow("%+f", tcomponents.ActiveComponents)
	//wg.Wait()
}

func closeComponent(wg *sync.WaitGroup, component *Worker) {
	defer wg.Done()
	component.Status <- false
	color.Red("component %s was be closed", component.Name)


}

func closeComponents(wg *sync.WaitGroup, tcomponents *TyphoonComponents) {
	defer wg.Done()
	for _, component := range tcomponents.ActiveComponents {
		wg.Add(1)
		go closeComponent(wg, component)

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


func Run(typhoonComponents[]string, configFile string)  {


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

	isProject := checkProject(typhoonComponents)
	if isProject == false {
		color.Red("Project does not exists in the current directory :%s", pathProject )
		os.Exit(1)
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


	//return

	var typhoonComponent = &TyphoonComponents{
		Components:  typhoonComponents,
		PathProject: pathProject,
		TyphoonPath: TyphoonPath,
	}
	color.Magenta("start components")
	task.wg.Add(1)
	go initComponents(&task.wg, typhoonComponent, configFile)



	task.wg.Add(1)
	go task.Run()
	//
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)


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
