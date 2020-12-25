package typhoon

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-cmd/cmd"
	"github.com/go-logfmt/logfmt"
	"github.com/kelseyhightower/envconfig"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
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
	w.Status = make(chan bool, 1)
	w.Status <- true
	projectEnv := fmt.Sprintf("PYTHONPATH=%s:%s", typhoonPath, w.ProjectPath)
	//color.Red("project path %s; projectEnv: %s", typhoonPath, projectEnv)
	newEnv := append(os.Environ(), projectEnv)
	envCmd.Env = newEnv
}

type Settings struct {
	Path string
	Status string

}

func LoadEnv2()  {
	findCmd := cmd.NewCmd("ls")
	statusChan := findCmd.Start() // non-blocking
	status := findCmd.Status()
	<- statusChan
	fmt.Printf("output---->: %s \n", status.Stdout)


}

func LoadEnv()  {

	loadStatus := false
	var pathProfile string
	pathHome := os.Getenv("HOME")

	pathsProfiles := []string{
		fmt.Sprintf("%s/.bash_profile", pathHome),
		fmt.Sprintf("%s/.bashrc", pathHome),
		fmt.Sprintf("%s/.bashprofile", pathHome),
		fmt.Sprintf("%s/.bash_rc", pathHome),
	}

	for _, _pathProfile := range pathsProfiles {
		fmt.Sprintf("path: %s", _pathProfile)

		if _, err := os.Stat(_pathProfile); !os.IsNotExist(err) {
			pathProfile = _pathProfile
			loadStatus = true
		}
	}

	if !loadStatus {
		color.Red("Not found bash profile !" )
		os.Exit(1)
	}

	color.Yellow("bash profile path: : %s", pathProfile)

	cmdSource := exec.Command("bash", "-c", "source " + pathProfile + "; env")

	bs, err := cmdSource.CombinedOutput()
	if err != nil {
		log.Fatalln(err)
	}
	s := bufio.NewScanner(bytes.NewReader(bs))

	for s.Scan() {
		kv := strings.SplitN(s.Text(), "=", 2)
		if strings.Contains(strings.ToLower(kv[0]), "typhoon") {
			os.Setenv(kv[0], kv[1])
		}
	}

}

func ReadEnv() error {


	LoadEnv()


	var settings Settings
	err := envconfig.Process("typhoon", &settings)
	if err != nil {
		log.Fatal(err.Error())
	}
	format := "Typhoon path: %v\nProjects path: %d\n"
	_, err = fmt.Printf(format, settings.Path, settings.Status)
	if err != nil {
		log.Fatal(err.Error())
	}
	color.Green("Typhoon read Env")
	//for _, e := range os.Environ() {
	//
	//	pair := strings.SplitN(e, "=", 2)
	//	fmt.Printf("%s: %s\n", pair[0], pair[1])
	//}
	return nil
}

func ParseLogData(logFileName string) error {
	currentPath, _ := os.Getwd()
	logPath := fmt.Sprintf("%s/%s", currentPath, logFileName)
	dat, err := ioutil.ReadFile(logPath)

	color.Red("Log file path: %s", logPath)
	if err != nil {

		color.Red("Log file not found")
		os.Exit(0)


	}

	logDataMap := logfmt.NewDecoder(strings.NewReader(string(dat)))
	for logDataMap.ScanRecord() {
		for logDataMap.ScanKeyval() {
			color.Yellow("%s = %s", logDataMap.Key(), logDataMap.Value())
		}
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

func initComponents(wg *sync.WaitGroup, tcomponents *TyphoonComponents, configFile string)  {
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

		tcomponents.ActiveComponents[component] = typhoonComponent

	}

	//color.Yellow("%+f", tcomponents.ActiveComponents)
	//wg.Wait()
}

func closeComponent(wg *sync.WaitGroup, component *Worker) {
	defer wg.Done()
	component.Status <- false
	color.Red("component %s was be closed", component.Name)
	//time.Sleep(time.Second * 2)

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

func Check(typhoonComponents[]string)  {
	pathProject, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	isProject := checkProject(typhoonComponents)
	if isProject == false {
		color.Red("Project does not exists in the current directory :%s", pathProject )
		os.Exit(1)
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

	if isExistDir("typhoon") {
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
