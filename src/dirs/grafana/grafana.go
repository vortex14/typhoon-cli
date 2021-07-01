package grafana

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os"
	"typhoon-cli/src/integrations/grafana"
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
				&cli.StringFlag{
					Name:    "grafana-dashboard",
					Aliases: []string{"g"},
					Value: "monitoring-grafana.json",
					Usage:   "Load configuration from `FILE`",
				},

			},
			Name: "import",
			Usage: "import project template to grafana api",
			Action: func(context *cli.Context) error {
				version := context.String("version")
				config := context.String("config")
				configDashboard := context.String("grafana-dashboard")
				pathProject, _ := os.Getwd()
				project := &typhoon.Project{
					Version: version,
					ConfigFile: config,
					Path: pathProject,
				}

				dashboard := grafana.DashBoard{
					ConfigName: configDashboard,
					Project: project,
				}
				dashboard.ImportGrafanaConfig()

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
				&cli.StringFlag{
					Name:    "grafana-dashboard",
					Aliases: []string{"g"},
					Value: "monitoring-grafana.json",
					Usage:   "Load configuration from `FILE`",
				},

			},
			Name: "remove",
			Usage: "Remove project dashboard from grafana api",
			Action: func(context *cli.Context) error {
				version := context.String("version")
				config := context.String("config")
				configDashboard := context.String("grafana-dashboard")
				pathProject, _ := os.Getwd()
				project := &typhoon.Project{
					Version: version,
					ConfigFile: config,
					Path: pathProject,
				}
				color.Red("Remove Grafana Dashboard")
				dashboard := grafana.DashBoard{
					ConfigName: configDashboard,
					Project: project,
				}
				dashboard.RemoveGrafanaDashboard()
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
			Name: "config",
			Usage: "Add to your config Grafana properties",
			Action: func(context *cli.Context) error {
				version := context.String("version")
				config := context.String("config")
				pathProject, _ := os.Getwd()
				project := &typhoon.Project{
					Version: version,
					ConfigFile: config,
					Path: pathProject,
				}
				dashboard := grafana.DashBoard{Project: project}
				dashboard.CreateBaseGrafanaConfig()
				return nil
			},
		},
		&cli.Command{
			Name:   "init",
			Usage: "Create grafana base template for project monitoring",
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
							Name:    "config",
							Aliases: []string{"c"},
							Value: "config.local.yaml",
							Usage:   "Load configuration from `FILE`",
						},
					},
					Name: "monitoring",
					Usage: "create base template for monitoring of project",
					Action: func(context *cli.Context) error {
						version := context.String("version")
						config := context.String("config")
						pathProject, _ := os.Getwd()
						project := &typhoon.Project{
							Version: version,
							ConfigFile: config,
							Path: pathProject,
						}

						dashboard := grafana.DashBoard{
							Project: project,
						}
						dashboard.CreateGrafanaMonitoringTemplates()
						return nil
					},
				},
				&cli.Command{
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "version",
							Aliases: []string{"v"},
							Value:   "v1.1",
							Usage:   "Create for available version",
						},
						&cli.StringFlag{
							Name:    "config",
							Aliases: []string{"c"},
							Value: "config.local.yaml",
							Usage:   "Load configuration from `FILE`",
						},
					},
					Name: "nsq-monitoring",
					Usage: "create base template for nsq monitoring of project",
					Action: func(context *cli.Context) error {
						version := context.String("version")
						config := context.String("config")
						pathProject, _ := os.Getwd()
						project := &typhoon.Project{
							Version: version,
							ConfigFile: config,
							Path: pathProject,
						}

						dashboard := grafana.DashBoard{
							Project: project,
						}
						dashboard.CreateGrafanaNSQMonitoringTemplates()
						return nil
					},
				},
			},

		},
}



