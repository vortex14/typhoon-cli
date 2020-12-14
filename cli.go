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
	//red := color.New(color.).PrintfFunc()
	//a: = red("tetst")
	//color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>\n")

	typhoonComponents := []string{"fetcher", "result_transporter", "donor", "processor", "scheduler"}

	app := &cli.App{
		Name: "Typhoon",
		UsageText: `
			typhoon up	--config=config.local.yaml
			typhoon run --component=scheduler
			typhoon run --components=scheduler,fetcher
`		,
		Description: "For running typhoon lite in command line",
		HelpName: "test thelp",
		EnableBashCompletion: true,
		Usage: "cli app",
		Commands: []*cli.Command{
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
