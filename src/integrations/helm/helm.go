package helm


import (
	"github.com/urfave/cli/v2"
	"typhoon-cli/src/typhoon"
)

var Commands = []*cli.Command{
	&cli.Command{
		Name:   "init",
		Usage: "Create helm resources",
		Subcommands: []*cli.Command{
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "version",
						Aliases: []string{"v"},
						Value:   "v1.1",
						Usage:   "Create for available version",
					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "test-kube-project",
						Usage:   "Project name",
					},
				},
				Name: "minikube",
				Usage: "create helm manifest for minikube",
				Action: func(context *cli.Context) error {
					version := context.String("version")
					name := context.String("name")
					project := &typhoon.Project{
						Version: version,
						Name: name,
					}
					project.BuildHelmMinikubeResources()
					return nil
				},
			},
		},

	},
}




