package tests


import (
	"github.com/urfave/cli/v2"
	"os"
	"typhoon-cli/src/typhoon"
)

var Commands = []*cli.Command{
	&cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Value:   "v1.1",
				Usage:   "Available version",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value: "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},

		},
		Name: "services",
		Usage: "Run test for connection to services from config.yaml",
		Action: func(context *cli.Context) error {
			version := context.String("version")
			config := context.String("config")
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				Version: version,
				ConfigFile: config,
				Path: pathProject,
			}
			project.RunTestServices()
			return nil
		},
	},
	&cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Value:   "v1.1",
				Usage:   "Available version",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value: "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},

		},
		Name: "check",
		Usage: "test check interface",
		Action: func(context *cli.Context) error {
			version := context.String("version")
			config := context.String("config")
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				Version: version,
				ConfigFile: config,
				Path: pathProject,
			}
			project.TestFunc()
			return nil
		},
	},
}



