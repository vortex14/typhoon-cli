package docker
import (
	"github.com/urfave/cli/v2"
	"typhoon-cli/src/typhoon"
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
					project.DockerBuild()
					return nil
				},
			},
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
					project.DockerListContainers()
					return nil
				},
			},
		},
	},
}




