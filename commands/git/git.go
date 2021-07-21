package git

import (
	"github.com/urfave/cli/v2"
	"github.com/vortex14/gotyphoon/src"
	"github.com/vortex14/gotyphoon/src/integrations/docker"
	"github.com/vortex14/gotyphoon/src/integrations/git"
)

var Commands = []*cli.Command{
	&cli.Command{
		Name:   "remove-untracked",
		Usage: "Remove all untracked files",
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
			project.LoadConfig()
			projectGit := git.Git{
				Path: project.GetProjectPath(),
			}
			projectGit.RemoveAllUnTrackingFiles()
			return nil
		},
	},
	&cli.Command{
		Name:   "status",
		Usage: "Status git files",
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
			project.LoadConfig()
			projectGit := git.Git{
				Path:    project.GetProjectPath(),
				Project: project,
			}
			projectGit.RepoStatus()
			return nil
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
}




