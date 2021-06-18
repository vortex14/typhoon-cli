package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/fatih/color"
	"io"
	"os"
	"strings"
	"time"
	"typhoon-cli/src/environment"
	"typhoon-cli/src/interfaces"
	"typhoon-cli/src/utils"
)

type Docker struct {
	Project interfaces.Project
}

var dockerRegistryUserID = ""

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type Log struct {
	Stream string `json:"stream"`
}

func (d *Docker) print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		var log Log
		json.Unmarshal(scanner.Bytes(), &log)
		fmt.Println(strings.ReplaceAll(log.Stream, "\n", ""))
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (d *Docker) build(workDir string, options *archive.TarOptions, opts types.ImageBuildOptions) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}



	color.Green("PatH: %s", workDir)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	tar, err := archive.TarWithOptions(fmt.Sprintf("%s/", workDir), options)
	if err != nil {
		return err
	}

	res, err := cli.ImageBuild(ctx, tar, opts)
	if err != nil {
		color.Red("Error!")
		return err
	}

	defer res.Body.Close()
	_ = d.print(res.Body)

	return nil
}

func (d *Docker) BuildImage()  {
	color.Yellow("Typhoon docker build ...")

	options := &archive.TarOptions{
		ExcludePatterns: []string{
			"extensions/tests/*",
			".git",
			"chrome",
	}}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{dockerRegistryUserID + "typhoon-lite"},
		Remove:     true,

	}

	envSetting := environment.Environment{}
	_, env := envSetting.GetSettings()

	err := d.build(env.Path, options, opts)
	if err != nil {
		color.Red(err.Error())
		return
	}

}

func (d *Docker) ListContainers()  {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

func (d *Docker) RunComponent(component string) error {
	d.Project.LoadConfig()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	if err != nil {
		//fmt.Println(err.Error())
		//color.Red("Error: %s", err)
		return err
	}

	containerConfig := &container.Config{
		Image: fmt.Sprintf("typhoon-lite-%s", d.Project.GetName()),
		Cmd: []string{"python", "donor.py --config=config.local.yaml --level=DEBUG"},
	}


	_, err = cli.ContainerCreate(ctx, containerConfig, nil, nil, nil, "typhoon")

	if err != nil {
		color.Red("ContainerCreate: %s", err)
		return err
	}

	return nil
}

func (d *Docker) ProjectBuild()  {
	u := utils.Utils{}

	_, dockerFile := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "ProjectDockerfile",

	})
	goTemplateDocker := interfaces.GoTemplate{
		Source: dockerFile,
		ExportPath: "Dockerfile",
		Data: map[string]string{
			"TYPHOON_IMAGE": d.Project.GetDockerImageName(),
		},
	}

	err := u.GoRunTemplate(&goTemplateDocker)
	if !err {
		color.Red("creation Dockerfile was fail")
		os.Exit(1)
	}

	color.Green("Dockerfile created!")

	color.Yellow("Typhoon project docker build ...")

	options := &archive.TarOptions{}
	projectConfig := d.Project.LoadConfig()
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{"typhoon-lite-" + projectConfig.Config.ProjectName},
		Remove:     true,

	}

	_ = d.build(d.Project.GetProjectPath(), options, opts)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}





}

func (d *Docker) RemoveResources()  {
	u := utils.Utils{}
	u.RemoveFiles([]string{
		"Dockerfile",
	})

	color.Green("Removed")
}