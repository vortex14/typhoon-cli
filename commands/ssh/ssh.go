package ssh

import (
	"github.com/urfave/cli/v2"
	"github.com/vortex14/gotyphoon"
	"github.com/vortex14/gotyphoon/integrations/ssh"
	"os"
)

var Commands = []*cli.Command{
	{
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
				Name:    "ip",
				Aliases: []string{"i"},
				Value: "127.0.0.1",
				Usage:   "Remote ip Addresses",
			},
			&cli.StringFlag{
				Name:    "login",
				Aliases: []string{"i"},
				Value: "root",
				Usage:   "Remote login",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Value: "password",
				Usage:   "Remote password",
			},

		},
		Name: "test-connection",
		Usage: "Run test for connection to services from config.yaml",
		Action: func(context *cli.Context) error {
			ssh := ssh.SSH{
				Login: context.String("login"),
				Password: context.String("password"),
				Ip: context.String("ip"),
			}
			ssh.TestConnection()
			return nil
		},
	},
	{
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



