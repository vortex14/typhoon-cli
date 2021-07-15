package main

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
	"strings"
	"typhoon-cli/src/dirs/docker"
	"typhoon-cli/src/dirs/gitlab"
	"typhoon-cli/src/dirs/grafana"
	"typhoon-cli/src/dirs/helm"
	"typhoon-cli/src/dirs/ssh"
	"typhoon-cli/src/dirs/tests"
	"typhoon-cli/src/environment"
	"typhoon-cli/src/interfaces"
	"typhoon-cli/src/typhoon"
	"typhoon-cli/src/utils"
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

					logFile := c.String("file")
					fileObject := interfaces.FileObject{
						Path: logFile,
					}
					_ = typhoon.ParseLogData(&fileObject)
					return nil
				},
			},
			{
				Name: "init",
				Usage: "create symbolic link to typhoon",
				Action: func(context *cli.Context) error {
					color.Green("create symbolic link to typhoon ")
					project := typhoon.Project{}
					err := project.CreateSymbolicLink()
					return err
				},

			},
			{
				Name: "bashrc",
				Usage: "Read from ~/.bashrc Typhoon variables",
				Action: func(context *cli.Context) error {
					log.Printf("Read from bashrc")
					envSetting := environment.Environment{}
					_, env := envSetting.GetSettings()

					color.Green("TYPHOON_PATH: %s \nTYPHOON_PROJECTS: %s\n", env.Path, env.Projects)
					//log.Printf("%+f", env)
					return nil
				},
			},
			{
				Name: "migrate",
				Usage: "Migrate typhoon project to new version",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "new",
						Aliases: []string{"n"},
						Value:   "v1.1",
						Usage:   "migrate project to v1.1",

					},
					&cli.StringFlag{
						Name:    "tag",
						Aliases: []string{"t"},
						Value:   "v1.0.0",
						Usage:   "docker images tag",

					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"p"},
						Usage:   "Project name ",
						Required: true,
					},
				},
				Action: func(context *cli.Context) error {
					project := typhoon.Project{
						Version: context.String("new"),
						Name: context.String("name"),
						Tag: context.String("tag"),
					}
					project.Migrate()
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
					&cli.StringFlag{
						Name:    "level",
						Aliases: []string{"l"},
						Value: "DEBUG",
						Usage:   "LOG LEVEL",
					},
					&cli.StringFlag{
						Name:    "reload",
						Aliases: []string{"r"},
						Value: "true",
						Usage:   "Auto reloading project",
					},
				},
				Action:  func(c *cli.Context) error {
					configFile := c.String("config")
					reloadF := c.String("reload")
					logLevel := c.String("level")
					var reload bool
					if reloadF == "true" {
						reload = true
					} else {
						reload = false
					}

					pathProject, err := os.Getwd()
					if err != nil {
						log.Println(err)
					}

					project := &typhoon.Project{
						SelectedComponent: typhoonComponents,
						ConfigFile: configFile,
						AutoReload: reload,
						Path: pathProject,
						LogLevel: logLevel,
					}
					project.Run()
					return nil
				},
			},
			{
				Name:    "down",
				Aliases: []string{"u"},
				Usage:   "down typhoon project",
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
					pathProject, err := os.Getwd()
					if err != nil {
						log.Println(err)
					}

					project := &typhoon.Project{
						SelectedComponent: typhoonComponents,
						ConfigFile: configFile,
						Path: pathProject,
					}
					project.Down()
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
					utils := utils.Utils{}
					if len(componentsName) > 0 {
						componentsArr := strings.Split(componentsName, ",")

						for _, name := range componentsArr {
							_, found := utils.CheckSlice(typhoonComponents, name)
							//color.Yellow("%t %s", found, name)
							if !found {
								color.Red("component %s isn't valid", name)
								os.Exit(1)
							}
						}

						project := &typhoon.Project{
							SelectedComponent: componentsArr,
						}

						project.CheckProject()

					} else {
						color.Yellow("run: %s", componentName)
						project := &typhoon.Project{
							SelectedComponent: []string{componentName},
						}

						project.CheckProject()
					}


					return nil
				},

			},
			{
				Name: "watch",
				Usage: "Watch for changing in typhoon project",
				Action: func(context *cli.Context) error {
					typhoon.WatchTest()
					return nil
				},
			},
			{
				Name: "create",
				Usage: "Create a new Typhoon project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "new",
						Aliases: []string{"n"},
						Value:   "v1.1",
						Usage:   "Project skeleton template version.",
					},
					&cli.StringFlag{
						Name:    "tag",
						Aliases: []string{"t"},
						Value:   "v1.0.0",
						Usage:   "docker images tag",
					},
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"p"},
						Usage:    "Project name ",
						Required: true,
					},
				},
				Action: func(context *cli.Context) error {

					project := typhoon.Project{
						Version: context.String("new"),
						Name: context.String("name"),
						Tag: context.String("tag"),
					}
					project.CreateProject()
					return nil
				},
			},
			{
				Name:        "ci-templates",
				Usage:       "CI Templates",
				Subcommands: gitlab.Commands,
			},
			{
				Name: "tests",
				Usage: "Typhoon tests",
				Subcommands: tests.Commands,
			},
			{
				Name: "ssh",
				Usage: "SSH connection",
				Subcommands: ssh.Commands,
			},
			{
				Name: "import",
				Usage: "Typhoon import",
				Subcommands: typhoon.Commands,
			},
			{
				Name: "cluster",
				Usage: "Typhoon cluster",
				Subcommands: typhoon.ClusterCommands,
			},
			{
				Name: "grafana",
				Usage: "Grafana integration",
				Subcommands: grafana.Commands,
			},
			{
				Name: "docker",
				Usage: "Docker integration",
				Subcommands: docker.Commands,
			},
			{
				Name:        "helm",
				Usage:       "Helm integration",
				Subcommands: helm.Commands,
			},
			{

				Name: "transporter",
				Usage: "Manage of transporter component",
				Subcommands: []*cli.Command{
					&cli.Command{
						Name:   "create",
						Usage: "Create resource for component",
						Subcommands: []*cli.Command{
							&cli.Command{
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:    "version",
										Aliases: []string{"v"},
										Value:   "v1.1",
										Usage:   "Create for available version",
									},
								},
								Name: "manifest",
								Usage: "generate transporter yaml manifest",
								Action: func(context *cli.Context) error {
									version := context.String("version")
									project := &typhoon.Project{
										Version: version,
										BuilderOptions: &interfaces.BuilderOptions{
											Component: "transporter",
											Type: "manifest",
										},
									}
									project.Build()
									return nil
								},
							},
						},

					},
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
						Name:    "level",
						Aliases: []string{"l"},
						Value: "DEBUG",
						Usage:   "LOG LEVEL",
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
					&cli.StringFlag{
						Name:    "reload",
						Aliases: []string{"r"},
						Value: "true",
						Usage:   "Auto reloading project",
					},
				},
				Action:  func(c *cli.Context) error {
					configFile := c.String("config")
					componentName := c.String("component")
					componentsName := c.String("components")
					logLevel := c.String("level")
					utils := utils.Utils{}
					//reload := c.String("reload")

					if len(componentsName) > 0 {
						componentsArr := strings.Split(componentsName, ",")

						for _, name := range componentsArr {
							_, found := utils.CheckSlice(typhoonComponents, name)
							color.Yellow("%s %s", found, name)
							if !found {
								color.Red("component %s isn't valid", name)
								os.Exit(1)
							}
						}

						project := &typhoon.Project{
							SelectedComponent: componentsArr,
							AutoReload: true,
							ConfigFile: configFile,
							LogLevel: logLevel,
						}




						project.Run()

					} else {
						color.Yellow("run: %s , config: %s", componentName, configFile)
						project := &typhoon.Project{
							SelectedComponent: []string{componentName},
							AutoReload: true,
							ConfigFile: configFile,
						}

						project.Run()
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
