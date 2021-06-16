package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/fatih/color"
	"io"
	"strings"
	"time"
	"typhoon-cli/src/environment"
	"typhoon-cli/src/interfaces"
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

func (d *Docker) imageBuild(client *client.Client) error {
	envSetting := environment.Environment{}
	_, env := envSetting.GetSettings()

	color.Green("PatH: %s", env.Path)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	tar, err := archive.TarWithOptions(fmt.Sprintf("%s/", env.Path), &archive.TarOptions{
		ExcludePatterns: []string{
			"extensions/tests/*",
			".git",
			"chrome",
		}},
	)
	if err != nil {
		return err
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{dockerRegistryUserID + "typhoon-lite"},
		Remove:     true,

	}
	res, err := client.ImageBuild(ctx, tar, opts)
	if err != nil {
		color.Red("Error!")
		return err
	}

	defer res.Body.Close()

	d.print(res.Body)

	//if err != nil {
	//	return err
	//}
	return nil
}

func (d *Docker) BuildImage()  {
	color.Yellow("Typhoon docker build ...")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return
	}


	err = d.imageBuild(cli)
	if err != nil {
		fmt.Println(err.Error())
		return
	}


	//opts := types.ImageBuildOptions{
	//	Dockerfile:  "Dockerfile",
	//	Tags:        []string{dockerRegistryUserID + "/node-hello"},
	//	Remove:      true,
	//}
	//res, err := dockerClient.ImageBuild(ctx, tar, opts)
	//if err != nil {
	//	return err
	//}
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