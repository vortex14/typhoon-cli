package grafana

import (
	"context"
	"encoding/json"
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
	DashBoardUrl string
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

	folderUID := configProject.Config.Grafana.FolderId
	var FolderId int
	if len(folderUID) > 0 {
		data, _ := c.GetFolderByUID(ctx, folderUID)
		if data.ID == 0 {
			color.Red("Folder not found. UUID %s", folderUID)
			os.Exit(1)
			return
		}

		FolderId = data.ID
	} else {
		FolderId = sdk.DefaultFolderId
	}


	params := sdk.SetDashboardParams{
		FolderID:  FolderId,
		Overwrite: false,
	}
	configProject.Config.Grafana.Id = board.UID
	configProject.Config.Grafana.DashboardUrl = configProject.Config.Grafana.Endpoint + "d/" + configProject.Config.Grafana.Id
	_, err := c.SetDashboard(ctx, board, params)
	if err != nil {
		color.Red("Error %s. board: %s", err, board.Title)
		os.Exit(1)
	}


	//color.Yellow("config file name: %s", configProject.ConfigFile)
	//color.Green("config %+v", configProject.Config.Grafana)

	configDumpData, err := yaml.Marshal(&configProject.Config)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	u := &utils.Utils{}
	err = u.DumpToFile(&interfaces.FileObject{
		Name: d.ConfigName,
		Data: string(configDumpData),
		Path: configProject.ConfigFile,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	color.Green("%s created !", configProject.Config.Grafana.Name)
	//fmt.Printf("--- m dump:\n%s\n\n", string(configDumpData))
}

func (d *DashBoard) RemoveGrafanaDashboard()  {
	configProject := d.Project.LoadConfig()
	ctx, c := d.getClient(configProject)
	_, err := c.DeleteDashboardByUID(ctx, configProject.Config.Grafana.Id)
	if err != nil {
		color.Red("%+v", err)
		os.Exit(1)
	}
	color.Green("%s was be removed.", configProject.Config.Grafana.Name)
	configProject.Config.Grafana.Id = ""
	configDumpData, err := yaml.Marshal(&configProject.Config)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	configProject.Config.Grafana.DashboardUrl = ""
	u := &utils.Utils{}
	err = u.DumpToFile(&interfaces.FileObject{
		Name: d.ConfigName,
		Data: string(configDumpData),
		Path: configProject.ConfigFile,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	//_, data := c.GetAllFolders(ctx)
	//color.Red("%+v", data)
}
