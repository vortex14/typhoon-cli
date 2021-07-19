package typhoon

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strconv"
	"typhoon-cli/src/integrations/gitlab"
	"typhoon-cli/src/integrations/grafana"
	"typhoon-cli/src/interfaces"
	"typhoon-cli/src/utils"
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
			project := &Project{
				Version: version,
				ConfigFile: config,
				Path: pathProject,
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
			project := &Project{
				Version: version,
				ConfigFile: config,
				Path: pathProject,
			}
			project.TestFunc()
			return nil
		},
	},
}



var ClusterCommands = []*cli.Command{
	&cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "test",
				Usage:   "Cluster name",
			},
			&cli.StringFlag{
				Name:    "description",
				Aliases: []string{"d"},
				Value:   "Test cluster",
				Usage:   "Cluster description",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "cluster.local.yaml",
				Usage:   "Cluster config yaml",
			},
		},
		Name: "create",
		Usage: "Create a new Typhoon cluster",
		Action: func(context *cli.Context) error {
			name := context.String("name")
			configName := context.String("config")
			description := context.String("description")
			cluster := Cluster{
				Name: name,
				Description: description,
				Config: configName,
			}
			cluster.Create()
			return nil
		},
	},
	&cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "test",
				Usage:   "Cluster name",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "cluster.local.yaml",
				Usage:   "Cluster config yaml",
			},
		},
		Name: "add",
		Usage: "Add the project to Typhoon cluster",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			cluster := Cluster{
				Config: configClusterName,
				Name: clusterName,
			}
			cluster.Add()
			return nil
		},
	},
	&cli.Command{
		Flags: []cli.Flag{},
		Name: "show",
		Usage: "Show all Typhoon clusters",
		Action: func(context *cli.Context) error {
			cluster := Cluster{}
			cluster.Show()
			return nil
		},
	},
	&cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "test",
				Usage:   "Cluster name",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "cluster.local.yaml",
				Usage:   "Cluster config yaml",
			},
		},
		Name: "sync-gitlab",
		Usage: "Typhoon sync cluster with gitlab projects",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			cluster := Cluster{
				Config: configClusterName,
				Name: clusterName,
			}

			gitLabServer := gitlab.Server{
				Cluster: &cluster,
			}

			gitLabServer.SyncGitlabProjects()

			return nil
		},
	},
	&cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "test",
				Usage:   "Cluster name",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "cluster.local.yaml",
				Usage:   "Cluster config yaml",
			},
		},
		Subcommands: []*cli.Command{
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "test",
						Usage:   "Cluster name",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "cluster.local.yaml",
						Usage:   "Cluster config yaml",
					},
				},
				Name: "gitlab",
				Usage: "Deploy Typhoon cluster to Gitlab",
				Action: func(context *cli.Context) error {
					clusterName := context.String("name")
					configClusterName := context.String("config")
					cluster := Cluster{
						Config: configClusterName,
						Name: clusterName,
					}

					gitLabServer := gitlab.Server{
						Cluster: &cluster,
					}

					gitLabServer.Deploy()

					return nil
				},
			},
		},
		Name: "deploy",
		Usage: "Typhoon deploy the cluster into env",
	},
	&cli.Command{
		Subcommands: []*cli.Command{
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "test",
						Usage:   "Cluster name",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "cluster.local.yaml",
						Usage:   "Cluster config yaml",
					},
					&cli.StringFlag{
						Name:    "grafana-dashboard",
						Aliases: []string{"g"},
						Value: "monitoring-grafana.json",
						Usage:   "Load configuration from `FILE`",
					},

				},
				Name: "import",
				Usage: "import cluster of projects template to grafana api",
				Action: func(context *cli.Context) error {
					clusterName := context.String("name")
					configClusterName := context.String("config")
					cluster := Cluster{
						Config: configClusterName,
						Name: clusterName,
					}

					projects := cluster.GetProjects()
					envSettings := cluster.GetEnvSettings()
					clusterConfig := cluster.LoadConfig(envSettings)
					configDashboard := context.String("grafana-dashboard")
					folderId := clusterConfig.Meta.Grafana.FolderId
					header := []string{"№", "name", "url"}
					var data[][]string
					for i, projectCluster := range projects {
						project := &Project{
							ConfigFile: projectCluster.Config,
							Path: envSettings.Projects + "/" + projectCluster.Name,
						}
						dashboard := grafana.DashBoard{
							ConfigName: configDashboard,
							Project: project,
						}
						configImportedDashboard := dashboard.ImportGrafanaConfig(folderId)
						projectCluster.Labels.Grafana = append(projectCluster.Labels.Grafana, configImportedDashboard)
						data = append(data, []string{strconv.Itoa(i+1),  projectCluster.Name,  configImportedDashboard.DashboardUrl})

					}
					u := utils.Utils{}
					u.RenderTableOutput(header, data)
					cluster.Projects = projects

					cluster.SaveConfig()
					return nil
				},
			},
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "test",
						Usage:   "Cluster name",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "cluster.local.yaml",
						Usage:   "Cluster config yaml",
					},
				},
				Name: "init",
				Subcommands: []*cli.Command{
					&cli.Command{
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Value:   "test",
								Usage:   "Cluster name",
							},
							&cli.StringFlag{
								Name:    "config",
								Aliases: []string{"c"},
								Value:   "cluster.local.yaml",
								Usage:   "Cluster config yaml",
							},
						},
						Name: "monitoring",
						Usage: "create base template for monitoring of project",
						Action: func(context *cli.Context) error {
							clusterName := context.String("name")
							configClusterName := context.String("config")
							cluster := Cluster{
								Config: configClusterName,
								Name: clusterName,
							}
							projects := cluster.GetProjects()
							envSettings := cluster.GetEnvSettings()
							for _, projectCluster := range projects {
								project := &Project{
									ConfigFile: projectCluster.Config,
									Path: envSettings.Projects + "/" + projectCluster.Name,
								}
								dashboard := grafana.DashBoard{
									Project: project,
								}
								dashboard.CreateGrafanaMonitoringTemplates()
							}
							return nil
						},
					},
					&cli.Command{
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Value:   "test",
								Usage:   "Cluster name",
							},
							&cli.StringFlag{
								Name:    "config",
								Aliases: []string{"c"},
								Value:   "cluster.local.yaml",
								Usage:   "Cluster config yaml",
							},
						},
						Name: "nsq-monitoring",
						Usage: "create base template for nsq monitoring of project",
						Action: func(context *cli.Context) error {
							clusterName := context.String("name")
							configClusterName := context.String("config")
							cluster := Cluster{
								Config: configClusterName,
								Name: clusterName,
							}
							projects := cluster.GetProjects()
							envSettings := cluster.GetEnvSettings()
							for _, projectCluster := range projects {
								project := &Project{
									ConfigFile: projectCluster.Config,
									Path: envSettings.Projects + "/" + projectCluster.Name,
								}
								dashboard := grafana.DashBoard{
									Project: project,
								}
								dashboard.CreateGrafanaNSQMonitoringTemplates()
							}
							return nil
						},
					},
				},
				Usage: "Create grafana base template monitoring for each project of the cluster",
			},
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "test",
						Usage:   "Cluster name",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "cluster.local.yaml",
						Usage:   "Cluster config yaml",
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
					clusterName := context.String("name")
					configClusterName := context.String("config")
					cluster := &Cluster{
						Config: configClusterName,
						Name: clusterName,
					}
					projects := cluster.GetProjects()
					envSettings := cluster.GetEnvSettings()
					configDashboard := context.String("grafana-dashboard")

					for _, projectCluster := range projects {
						project := &Project{
							ConfigFile: projectCluster.Config,
							Path: envSettings.Projects + "/" + projectCluster.Name,
						}
						dashboard := grafana.DashBoard{
							ConfigName: configDashboard,
							Project: project,
							Cluster: cluster,
						}
						err, imported := dashboard.RemoveGrafanaDashboard()
						if err != nil {
							color.Red("%s", err.Error())
							os.Exit(1)
						}

						for di, dashboardRecord := range projectCluster.Labels.Grafana {
							if dashboardRecord.Id == imported.Id {
								projectCluster.Labels.Grafana = append(projectCluster.Labels.Grafana[:di], projectCluster.Labels.Grafana[di+1:]...)
							}
						}
						//color.Green("%+v", imported)

					}
					cluster.SaveConfig()
					return nil
				},
			},
		},
		Name: "grafana",
		Usage: "Integration of Typhoon cluster with Grafana",
	},
	&cli.Command{
		Subcommands: []*cli.Command{
			&cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "test",
						Usage:   "Cluster name",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "cluster.local.yaml",
						Usage:   "Cluster config yaml",
					},
					&cli.StringFlag{
						Name:    "docker-file",
						Aliases: []string{"d"},
						Value: "Dockerfile",
						Usage:   "Load Dockerfile template from `FILE`",
					},
					&cli.StringFlag{
						Name:    "tag",
						Aliases: []string{"t"},
						Value: "typhoon-lite:latest",
						Usage:   "Typhoon Lite Tag",
					},

				},
				Name: "generate",
				Usage: "Generate own a dockerfile for all cluster",
				Action: func(context *cli.Context) error {
					clusterName := context.String("name")
					imageTag := context.String("tag")
					configClusterName := context.String("config")
					dockerFileTemplate := context.String("docker-file")
					u := utils.Utils{}
					cluster := Cluster{
						Config: configClusterName,
						Name: clusterName,
					}
					projects := cluster.GetProjects()
					envSettings := cluster.GetEnvSettings()
					//clusterConfig := cluster.LoadConfig(envSettings)

					header := []string{"№", "name"}
					var data[][]string
					for i, projectCluster := range projects {
						project := &Project{
							ConfigFile: projectCluster.Config,
							Path: envSettings.Projects + "/" + projectCluster.Name,
						}
						project.LoadConfig()

						path := project.GetProjectPath()
						//color.Green("%s", path)
						exportPath := fmt.Sprintf("%s/Dockerfile", path)

						fileObject := &interfaces.FileObject{
							Path: ".",
							Name: dockerFileTemplate,
						}

						lables := []interfaces.ReplaceLabel{
							interfaces.ReplaceLabel{Label: "{{.Image}}", Value: imageTag},
							interfaces.ReplaceLabel{Label: "{{.Config}}", Value: project.GetConfigFile()},
						}

						err := u.CopyFileAndReplaceLabelsFromHost(exportPath, lables, fileObject)

						if err != nil {

							color.Red("Error: %s", err)
							os.Exit(0)

						}
						data = append(data, []string{strconv.Itoa(i+1),  projectCluster.Name})

					}

					u.RenderTableOutput(header, data)
					return nil
				},
			},
		},
		Name: "docker",
		Usage: "Integration of Typhoon cluster with Docker",
	},
}





var watcher *fsnotify.Watcher

func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func CreateProject() {

}

func WatchTest()  {
	color.Green("watch for project ..")
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories

	if err := filepath.Walk("project", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done

}



func ParseLogData(fileObject *interfaces.FileObject) error {
	u := utils.Utils{}
	err := u.ParseLog(fileObject)


	if err != nil {

		color.Red("Error %s", err)
		os.Exit(0)

	}

	return nil


}
