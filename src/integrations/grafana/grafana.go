package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/grafana-tools/sdk"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"typhoon-cli/src/interfaces"
	"typhoon-cli/src/typhoon/config"
	"typhoon-cli/src/utils"
)

type DashBoard struct {
	ConfigName string
	Project interfaces.Project
}

type Config struct {
	Name string
	Endpoint string
	token string
}

func (d *DashBoard) getClient(configProject *config.ConfigProject) (context.Context, *sdk.Client) {
	c := sdk.NewClient(configProject.Config.Grafana.Endpoint, configProject.Config.Grafana.Token, sdk.DefaultHTTPClient)
	ctx := context.Background()
	return ctx, c

}

func (d *DashBoard) ImportGrafanaConfig()  {
	configProject := d.Project.LoadConfig()
	ctx, c := d.getClient(configProject)
	rawBoard, errBoard := ioutil.ReadFile(d.ConfigName)
	if errBoard != nil {
		log.Println(errBoard)
		os.Exit(1)
	}
	var board sdk.Board
	if err := json.Unmarshal(rawBoard, &board); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	out, _ := exec.Command("uuidgen").Output()
	board.UID = string(out[:8])
	board.Title = configProject.Config.Grafana.Name
	color.Green("Creating dashboard %s", board.Title)
	params := sdk.SetDashboardParams{
		FolderID:  sdk.DefaultFolderId,
		Overwrite: false,
	}
	configProject.Config.Grafana.Id = board.UID
	_, err := c.SetDashboard(ctx, board, params)
	if err != nil {
		color.Red("Error %s. board: %s", err, board.Title)
		//os.Exit(1)
	}
	color.Yellow("config file name: %s", configProject.ConfigFile)
	color.Green("config %+v", configProject.Config.Grafana)

	configDumpData, err := yaml.Marshal(&configProject.Config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	u := &utils.Utils{}
	err = u.DumpToFile(&interfaces.FileObject{
		Name: d.ConfigName,
		Data: string(configDumpData),
		Path: configProject.ConfigFile,
	})
	if err != nil {
		return
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(configDumpData))
}

func (d *DashBoard) RemoveGrafanaDashboard()  {
	configProject := d.Project.LoadConfig()
	ctx, c := d.getClient(configProject)
	_, err := c.DeleteDashboardByUID(ctx, "5DEF1A09")
	if err != nil {
		color.Red("%+v", err)
		os.Exit(1)
	}
	//_, data := c.GetAllFolders(ctx)
	//color.Red("%+v", data)
}
