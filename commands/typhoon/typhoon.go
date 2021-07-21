package typhoon

import (
	"github.com/urfave/cli/v2"
	"github.com/vortex14/gotyphoon/src"
	"os"
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
			&cli.StringFlag{
				Name:    "source",
				Aliases: []string{"s"},
				Value: "source.json",
				Usage:   "Load source from `FILE`",
			},
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
				Value: "https://200.ru",
				Usage:   "Response url",
			},

		},
		Name: "cache",
		Usage: "Load cache and save to Redis Storage",
		Action: func(context *cli.Context) error {
			//color.Red("Load cache and save to Redis Storage")
			version := context.String("version")
			config := context.String("config")
			url := context.String("url")
			pathProject, _ := os.Getwd()
			sourceFile := context.String("source")
			project := &typhoon.Project{
				Version:    version,
				ConfigFile: config,
				Path:       pathProject,
			}
			project.ImportResponseData(url, sourceFile)
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
				Version:    version,
				ConfigFile: config,
				Path:       pathProject,
			}
			project.TestFunc()
			return nil
		},
	},
}





