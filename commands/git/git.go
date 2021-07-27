package git

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/vortex14/gotyphoon"
	"github.com/vortex14/gotyphoon/integrations/docker"
	"github.com/vortex14/gotyphoon/integrations/git"
	"os"
)

var Commands = []*cli.Command{
	{
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
	{
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
	{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.local.yaml",
				Usage:   "Project config yaml",
			},
			&cli.StringFlag{
				Name:    "backup",
				Aliases: []string{"b"},
				Value:   "true",
				Usage:   "force git backup",
			},
			&cli.StringFlag{
				Name:    "remote",
				Aliases: []string{"r"},
				Usage:   "git remote",
			},
			&cli.StringFlag{
				Name:    "branch",
				Aliases: []string{"bc"},
				Usage:   "branch",
			},
		},
		Name: "reset",
		Usage: "Synchronize local files by remote branch",
		Action: func(context *cli.Context) error {
			config := context.String("config")
			backup := context.String("backup")
			remote := context.String("remote")
			branch := context.String("branch")

			if len(branch) == 0 || len(remote) == 0 {
				color.Red("Required Arguments: remote, branch")
				os.Exit(1)
			}

			backupFlag := false
			if backup == "true" {
				backupFlag = true
			}

			project := &typhoon.Project{
				ConfigFile: config,
			}
			project.LoadConfig()
			projectGit := git.Git{
				Path: project.GetProjectPath(),
			}
			color.Green("Remove untracked files for %s:", project.Name)
			projectGit.LocalResetLikeRemote(remote, branch, backupFlag)
			fmt.Println("")


			color.Yellow("Run git reset for all %s. backup: %t", project.GetName(), backupFlag)

			return nil
		},
	},
	{
		Name:   "push",
		Usage: "Push Docker resources to registry",
		Subcommands: []*cli.Command{
			{
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
			{
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




