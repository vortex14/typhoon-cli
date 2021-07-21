package docker

import (
	"github.com/urfave/cli/v2"
	"github.com/vortex14/gotyphoon/src"
	"github.com/vortex14/gotyphoon/src/integrations/docker"
)

var Commands = []*cli.Command{
	&cli.Command{
		Name:   "build",
		Usage: "Create new Typhoon build",
		Subcommands: []*cli.Command{
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "typhoon-lite:latest",
						Usage:   "Pass image name",
					},
				},
				Name: "image",
				Usage: "Create new image",
				Action: func(context *cli.Context) error {
					imageName := context.String("name")
					project := &typhoon.Project{
						DockerImageName: imageName,
					}
					projectDocker := docker.Docker{Project: project}
					projectDocker.BuildImage()
					return nil
				},
			},
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "typhoon-lite:latest",
						Usage:   "Pass image name",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value: "config.local.yaml",
						Usage:   "Load configuration from `FILE`",
					},
				},
				Name: "project",
				Usage: "Build image for project",
				Action: func(context *cli.Context) error {
					imageName := context.String("name")
					configFile := context.String("config")
					project := &typhoon.Project{
						DockerImageName: imageName,
						ConfigFile:      configFile,
					}
					projectDocker := docker.Docker{
						Project: project,
					}
					projectDocker.BuildImage()
					projectDocker.PushImage()
					return nil
				},
			},
		},
	},
	&cli.Command{
		Name:   "push",
		Usage: "Push Docker resources to registry",
		Subcommands: []*cli.Command{
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "typhoon-lite:latest",
						Usage:   "Pass image name",
					},
					&cli.StringFlag{
						Name:    "latest-date",
						Aliases: []string{"l"},
						Value:   "true",
						Usage:   "Pass image name",
					},
				},
				Name: "image",
				Usage: "Push Typhoon docker image to Registry",
				Action: func(context *cli.Context) error {
					imageName := context.String("name")
					LatestDate := context.String("latest-date")
					project := &typhoon.Project{
						DockerImageName: imageName,
					}
					projectDocker := docker.Docker{Project: project, LatestTag: LatestDate}
					projectDocker.PushImage()
					return nil
				},
			},
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "typhoon-lite:latest",
						Usage:   "Pass image name",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value: "config.local.yaml",
						Usage:   "Load configuration from `FILE`",
					},
				},
				Name: "project",
				Usage: "Build image for project",
				Action: func(context *cli.Context) error {
					imageName := context.String("name")
					configFile := context.String("config")
					project := &typhoon.Project{
						DockerImageName: imageName,
						ConfigFile:      configFile,
					}
					projectDocker := docker.Docker{
						Project: project,
					}
					projectDocker.ProjectBuild()
					return nil
				},
			},
		},
	},
	&cli.Command{
		Name:   "build-push",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "typhoon-lite:latest",
				Usage:   "Pass image name",
			},
			&cli.StringFlag{
				Name:    "latest-date",
				Aliases: []string{"l"},
				Value:   "true",
				Usage:   "Pass image name",
			},
		},
		Usage: "Build and push Docker resources to registry",
		Action: func(context *cli.Context) error {
			imageName := context.String("name")
			LatestDate := context.String("latest-date")
			project := &typhoon.Project{
				DockerImageName: imageName,
			}
			projectDocker := docker.Docker{Project: project, LatestTag: LatestDate}
			projectDocker.BuildImage()
			projectDocker.PushImage()
			return nil
		},
	},
	&cli.Command{
		Name:   "list",
		Usage: "Create new Typhoon build",
		Subcommands: []*cli.Command{
			&cli.Command{
				Name: "containers",
				Usage: "Show list containers",
				Action: func(context *cli.Context) error {
					project := &typhoon.Project{}
					projectDocker := docker.Docker{Project: project}
					projectDocker.ListContainers()
					return nil
				},
			},
		},
	},
	&cli.Command{
		Name:   "run",
		Usage: "Create new Typhoon build",
		Subcommands: []*cli.Command{
			&cli.Command{
				Name: "component",
				Usage: "Run Typhoon component in docker container",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value: "config.local.yaml",
						Usage:   "Load configuration from `FILE`",
					},
				},
				Action: func(context *cli.Context) error {
					configFile := context.String("config")
					project := &typhoon.Project{
						ConfigFile: configFile,
					}
					projectDocker := docker.Docker{
						Project: project,
					}

					err := projectDocker.RunComponent("test")
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	},
	&cli.Command{
		Name:   "remove",
		Usage: "Remove Docker resources",
		Subcommands: []*cli.Command{
			&cli.Command{
				Name: "files",
				Usage: "Run Typhoon component in docker container",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value: "config.local.yaml",
						Usage:   "Load configuration from `FILE`",
					},
				},
				Action: func(context *cli.Context) error {
					configFile := context.String("config")
					project := &typhoon.Project{
						ConfigFile: configFile,
					}
					projectDocker := docker.Docker{
						Project: project,
					}

					projectDocker.RemoveResources()
					return nil
				},
			},
		},
	},
}




