package gitlab

import (
	"github.com/urfave/cli/v2"
	"github.com/vortex14/gotyphoon/src"
)

var Commands = []*cli.Command{
	&cli.Command{
		Name:   "init",
		Usage: "Create Ci base templates",
		Subcommands: []*cli.Command{
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "typhoon-image",
						Aliases: []string{"v"},
						Value:   "typhoon-s1.ru/typhoon-lite/typhoon:2020.04.13-2",
						Usage:   "Create for available version",
					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "test-kube-project",
						Usage:   "Project name",
					},
				},
				Name: "templates",
				Usage: "create ci templates for k8 cluster",
				Action: func(context *cli.Context) error {
					version := context.String("typhoon-image")
					name := context.String("name")
					project := &typhoon.Project{
						Version: version,
						Name:    name,
					}
					project.BuildCIResources()
					return nil
				},
			},
		},
	},
}



