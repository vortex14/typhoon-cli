package environment

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Settings struct {
	Path string
	Status string
	Projects string

}

type Environment struct {}

func (e *Environment) Load()  {
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
		log.Fatal(err)
	}
	s := bufio.NewScanner(bytes.NewReader(bs))

	for s.Scan() {
		kv := strings.SplitN(s.Text(), "=", 2)
		if strings.Contains(strings.ToLower(kv[0]), "typhoon") {
			os.Setenv(kv[0], kv[1])
		}
	}

}

func (e Environment) Set()  {

}

func (e Environment) Get()  {

}

func (e *Environment) GetSettings() (error error, settings *Settings) {
	e.Load()
	envSetting := &Settings{}
	err := envconfig.Process("typhoon", envSetting)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil, envSetting
}
