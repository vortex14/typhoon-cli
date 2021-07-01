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
	"strings"
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

func (d *DashBoard) getClient(configProject *config.Project) (context.Context, *sdk.Client) {
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

func (d *DashBoard) CreateGrafanaMonitoringTemplates()  {
	d.Project.LoadConfig()
	color.Yellow("Creating Grafana monitoring template ...")
	u := utils.Utils{}

	fileObject := &interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "grafana-template.gojson",
	}
	validProjectName := strings.ReplaceAll(d.Project.GetName(), "-", "_")
	err := u.CopyFileAndReplaceLabel("monitoring-grafana.json",&interfaces.ReplaceLabel{Label: "{{.projectName}}", Value: validProjectName}, fileObject)

	if err != nil {

		color.Red("Error %s", err)
		os.Exit(0)

	}
}

func (d *DashBoard) CreateBaseGrafanaConfig()  {
	color.Yellow("Creating base grafana properties into typhoon project config.yaml")
	configProject := d.Project.LoadConfig()
	configProject.Config.Grafana = config.GrafanaConfig{
		Name: "Typhoon project dashboard",
		Id: "0000000",
		Token: "eyJrIjoiTGZqYUY3NWFsVk92MUc5TFFnTXlkYTg3WFJPME4wQVIiLCJuIjoidHlwaG9vbiIsImlkIjoxfQ==",
		Endpoint: "http://localhost:3000",
	}

	configDumpData, _ := yaml.Marshal(&configProject.Config)

	u := &utils.Utils{}
	err := u.DumpToFile(&interfaces.FileObject{
		Name: d.Project.GetConfigFile(),
		Data: string(configDumpData),
		Path: configProject.ConfigFile,
	})

	if err != nil {
		return
	}

	color.Green("%s updated.", d.Project.GetConfigFile())
}

func (d *DashBoard) CreateGrafanaNSQMonitoringTemplates()  {
	d.Project.LoadConfig()
	color.Yellow("Creating NSQ Grafana monitoring templates ...")
	u := utils.Utils{}

	fileObject := &interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "grafana-nsq-template.gojson",
	}
	err := u.CopyFileAndReplaceLabel("grafana-nsq-monitoring.json",&interfaces.ReplaceLabel{Label: "{{.projectName}}", Value: d.Project.GetName()}, fileObject)

	if err != nil {

		color.Red("Error %s", err)
		os.Exit(0)

	}
}
