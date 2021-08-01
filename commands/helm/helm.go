package helm

import (
	"github.com/urfave/cli/v2"
	"github.com/vortex14/gotyphoon"
	"github.com/vortex14/gotyphoon/integrations/helm"
)

var Commands = []*cli.Command{
	{
		Name:  "init",
		Usage: "Create helm resources",
		Subcommands: []*cli.Command{
			{
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
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "config.local.yaml",
						Usage:   "Load configuration from `FILE`",
					},
				},
				Name:  "minikube",
				Usage: "create helm manifest for minikube",
				Action: func(context *cli.Context) error {
					version := context.String("version")
					name := context.String("name")
					config := context.String("config")
					project := &typhoon.Project{
						Version:    version,
						ConfigFile: config,
						Name:       name,
					}
					helmResources := helm.Resources{
						Project: project,
					}
					helmResources.BuildHelmMinikubeResources()
					return nil
				},
			},
		},
	},
	{
		Name:  "remove",
		Usage: "Create helm resources",
		Subcommands: []*cli.Command{
			{
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
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "config.local.yaml",
						Usage:   "Load configuration from `FILE`",
					},
				},
				Name:  "minikube",
				Usage: "remove helm manifest for minikube",
				Action: func(context *cli.Context) error {
					helmResources := helm.Resources{}
					helmResources.RemoveHelmMinikubeManifests()
					return nil
				},
			},
		},
	},
}
