package main

import (
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
	"strings"
	"typhoon-cli/typhoon"
)

func main() {

	typhoonComponents := []string{"fetcher", "result_transporter", "donor", "processor", "scheduler"}

	app := &cli.App{
		Name: "Typhoon",
		UsageText: `
			typhoon up	--config=config.local.yaml
			typhoon run --component=scheduler
			typhoon run --components=scheduler,fetcher


			typhoon logging --file=test-log.log
`		,
		Description: "For running typhoon lite in command line",
		HelpName: "test thelp",
		EnableBashCompletion: true,
		Usage: "cli app",
		Commands: []*cli.Command{
			{
				Name: "logging",
				Usage: "Check typhoon parse log",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Value: "log.log",
						Usage:   "Load log from `FILE` for test parsing",
					},
				},
				Action: func(c *cli.Context) error {

					configFile := c.String("file")
					typhoon.ParseLogData(configFile)
					return nil
				},
			},
			{
				Name: "bashrc",
				Usage: "Read from ~/.bashrc Typhoon variables",
				Action: func(context *cli.Context) error {
					log.Printf("Read from bashrc")
					typhoon.ReadEnv()
					return nil
				},
			},
			{
				Name: "debug",
				Usage: "Run cli Ui for debugging",
				Action: func(context *cli.Context) error {
					typhoon.RunUI()
					return nil
				},


			},
			{
				Name:    "up",
				Aliases: []string{"u"},
				Usage:   "run typhoon project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value: "config.local.yaml",
						Usage:   "Load configuration from `FILE`",
					},
				},
				Action:  func(c *cli.Context) error {
					configFile := c.String("config")
					typhoon.Run(typhoonComponents, configFile)
					return nil
				},
			},
			{
				Name: "check",
				Aliases: []string{"rc"},
				Usage:   "Check health component of dir",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "component",
						Aliases: []string{"ct"},
						Value:   "processor",
						Usage:   "check component of dir",

					},
					&cli.StringFlag{
						Name:    "components",
						Aliases: []string{"cs"},
						Usage:   "check a few component of dir ",
					},
				},
				Action:  func(c *cli.Context) error {
					componentName := c.String("component")
					componentsName := c.String("components")

					if len(componentsName) > 0 {
						componentsArr := strings.Split(componentsName, ",")

						for _, name := range componentsArr {
							_, found := typhoon.Find(typhoonComponents, name)
							//color.Yellow("%t %s", found, name)
							if !found {
								color.Red("component %s isn't valid", name)
								os.Exit(1)
							}
						}

						typhoon.Check(componentsArr)

					} else {
						color.Yellow("run: %s", componentName)
						typhoon.Check([]string{componentName})
					}


					return nil
				},

			},
			{
				Name:    "run",
				Aliases: []string{"rc"},
				Usage:   "Run single component",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "component",
						Aliases: []string{"ct"},
						Value:   "processor",
						Usage:   "Run one component",

					},
					&cli.StringFlag{
						Name:    "components",
						Aliases: []string{"cs"},
						Usage:   "Run a few component",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value: "config.local.yaml",
						Usage:   "Load configuration from `FILE`",
					},
				},
				Action:  func(c *cli.Context) error {
					configFile := c.String("config")
					componentName := c.String("component")
					componentsName := c.String("components")

					if len(componentsName) > 0 {
						componentsArr := strings.Split(componentsName, ",")

						for _, name := range componentsArr {
							_, found := typhoon.Find(typhoonComponents, name)
							color.Yellow("%s %s", found, name)
							if !found {
								color.Red("component %s isn't valid", name)
								os.Exit(1)
							}
						}

						typhoon.Run(componentsArr, configFile)

					} else {
						color.Yellow("run: %s , config: %s", componentName, configFile)
						typhoon.Run([]string{componentName}, configFile)
					}


					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
