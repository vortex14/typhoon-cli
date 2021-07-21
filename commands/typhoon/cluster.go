package typhoon

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/vortex14/gotyphoon"
	"github.com/vortex14/gotyphoon/integrations/git"
	"github.com/vortex14/gotyphoon/integrations/gitlab"
	"github.com/vortex14/gotyphoon/integrations/grafana"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
	"os"
	"strconv"
)


var ClusterGitCommands = []*cli.Command{
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
		Name: "status",
		Usage: "Cluster repository status",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			envSettings := cluster.GetEnvSettings()
			projects := cluster.GetProjects()

			for _, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
				}
				project.LoadConfig()
				projectGit := git.Git{
					Path: project.GetProjectPath(),
				}
				_, branchName := project.GetBranch()
				color.Green("Git status of %s", project.GetName())
				color.Green("Project branch: %s", branchName)
				projectGit.RepoStatus()
				fmt.Println("")

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
		Name: "remove-untracked",
		Usage: "Remove all untracked files",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			envSettings := cluster.GetEnvSettings()
			projects := cluster.GetProjects()

			for _, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
				}
				projectGit := git.Git{
					Path: project.GetProjectPath(),
				}
				color.Green("Remove untracked files for %s:", projectCluster.Name)
				projectGit.RemoveAllUnTrackingFiles()
				fmt.Println("")
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
		Name: "remove-ignores",
		Usage: "Remove all ignore files",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			envSettings := cluster.GetEnvSettings()
			projects := cluster.GetProjects()

			for _, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
				}
				projectGit := git.Git{
					Path: project.GetProjectPath(),
				}
				color.Green("Remove untracked files for %s:", projectCluster.Name)
				projectGit.RemovePyCacheFiles()
				fmt.Println("")
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
			&cli.StringFlag{
				Name:    "branch",
				Aliases: []string{"b"},
				Usage:   "Cluster branch name",
			},
		},
		Name: "change-branch",
		Usage: "Change project branch to cluster branch",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			clusterBranch := context.String("branch")
			if len(clusterBranch) == 0 {
				color.Red("Not found arg branch")
				os.Exit(1)
			}
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			projects := cluster.GetProjects()
			for _, projectCluster := range projects {
				projectCluster.Labels.Git.Branch = clusterBranch
			}
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
		Name: "show-ignores",
		Usage: "Show all ignore files",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			envSettings := cluster.GetEnvSettings()
			projects := cluster.GetProjects()

			for _, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
				}
				projectGit := git.Git{
					Path: project.GetProjectPath(),
				}
				color.Green("Remove untracked files for %s:", projectCluster.Name)
				//projectGit.RemovePyCacheFiles()
				projectGit.ShowPyCacheFiles()
				fmt.Println("")
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
			&cli.StringFlag{
				Name:    "backup",
				Aliases: []string{"b"},
				Value:   "true",
				Usage:   "force git backup",
			},
		},
		Name: "reset",
		Usage: "Synchronize local files by remote branch",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			backup := context.String("backup")
			backupFlag := false
			if backup == "true" {
				backupFlag = true
			}
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			color.Yellow("Run git reset for all %s. backup: %t", cluster.Name, backupFlag)
			envSettings := cluster.GetEnvSettings()
			projects := cluster.GetProjects()

			for _, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
				}
				projectGit := git.Git{
					Path: project.GetProjectPath(),
				}
				color.Green("Remove untracked files for %s:", projectCluster.Name)
				projectGit.LocalResetLikeRemote(projectCluster.Labels.Git.Remote, projectCluster.Labels.Git.Branch, backupFlag)
				fmt.Println("")
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
		Name: "push",
		Usage: "Push cluster projects to git",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			envSettings := cluster.GetEnvSettings()
			projects := cluster.GetProjects()

			for _, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
				}
				projectGit := git.Git{
					Path: project.GetProjectPath(),
				}
				color.Green("Git Push %s:", projectCluster.Name)
				err, branch := project.GetBranch()
				if err != nil {
					color.Red("%s", err.Error())
					continue
				}
				projectGit.Push(projectCluster.Labels.Git.Remote, branch)
				fmt.Println("")
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
			&cli.StringFlag{
				Name:    "branch",
				Aliases: []string{"b"},
				Usage:   "new branch name",
			},
			&cli.StringFlag{
				Name:    "message",
				Aliases: []string{"m"},
				Usage:   "Update project",
			},
		},
		Name: "new-branch",
		Usage: "Create a new local branch and add new files",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			newBranchName := context.String("branch")
			message := context.String("message")
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			envSettings := cluster.GetEnvSettings()
			projects := cluster.GetProjects()

			for _, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
				}
				projectGit := git.Git{
					Path: project.GetProjectPath(),
				}
				projectGit.CreateBranchAndCommit(message, newBranchName)
				fmt.Println("")
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
			&cli.StringFlag{
				Name:    "message",
				Aliases: []string{"m"},
				Usage:   "Git message",
			},
		},
		Name: "commit",
		Usage: "Create a new commit and add new files",
		Action: func(context *cli.Context) error {
			clusterName := context.String("name")
			configClusterName := context.String("config")
			message := context.String("message")
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			envSettings := cluster.GetEnvSettings()
			projects := cluster.GetProjects()

			for _, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
				}
				projectGit := git.Git{
					Path: project.GetProjectPath(),
				}
				projectGit.AddAndCommit(message)
				fmt.Println("")
			}

			return nil
		},
	},
}


var ClusterDeployCommands = []*cli.Command{
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
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}

			gitLabServer := gitlab.Server{
				Cluster: &cluster,
			}

			gitLabServer.Deploy()

			return nil
		},
	},
}
var ClusterGrafanaCommands = []*cli.Command{
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
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}

			projects := cluster.GetProjects()
			envSettings := cluster.GetEnvSettings()
			clusterConfig := cluster.LoadConfig(envSettings)
			configDashboard := context.String("grafana-dashboard")
			folderId := clusterConfig.Meta.Grafana.FolderId
			header := []string{"№", "name", "url"}
			var data[][]string
			for i, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
				}
				dashboard := grafana.DashBoard{
					ConfigName: configDashboard,
					Project:    project,
				}
				configImportedDashboard := dashboard.ImportGrafanaConfig(folderId)
				projectCluster.Labels.Grafana = append(projectCluster.Labels.Grafana, configImportedDashboard)
				data = append(data, []string{strconv.Itoa(i +1),  projectCluster.Name,  configImportedDashboard.DashboardUrl})

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
					cluster := typhoon.Cluster{
						Config: configClusterName,
						Name:   clusterName,
					}
					projects := cluster.GetProjects()
					envSettings := cluster.GetEnvSettings()
					for _, projectCluster := range projects {
						project := &typhoon.Project{
							ConfigFile: projectCluster.Config,
							Path:       envSettings.Projects + "/" + projectCluster.Name,
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
					cluster := typhoon.Cluster{
						Config: configClusterName,
						Name:   clusterName,
					}
					projects := cluster.GetProjects()
					envSettings := cluster.GetEnvSettings()
					for _, projectCluster := range projects {
						project := &typhoon.Project{
							ConfigFile: projectCluster.Config,
							Path:       envSettings.Projects + "/" + projectCluster.Name,
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
			cluster := &typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			projects := cluster.GetProjects()
			envSettings := cluster.GetEnvSettings()
			configDashboard := context.String("grafana-dashboard")

			for _, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
				}
				dashboard := grafana.DashBoard{
					ConfigName: configDashboard,
					Project:    project,
					Cluster:    cluster,
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
}

var ClusterDockerCommands = []*cli.Command{
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
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
			}
			projects := cluster.GetProjects()
			envSettings := cluster.GetEnvSettings()

			cluster.Meta.Docker.Image = imageTag

			header := []string{"№", "name"}
			var data[][]string
			for i, projectCluster := range projects {
				project := &typhoon.Project{
					ConfigFile: projectCluster.Config,
					Path:       envSettings.Projects + "/" + projectCluster.Name,
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
				data = append(data, []string{strconv.Itoa(i +1),  projectCluster.Name})

			}

			cluster.SaveConfig()
			u.RenderTableOutput(header, data)
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
			cluster := typhoon.Cluster{
				Name:        name,
				Description: description,
				Config:      configName,
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
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
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
			cluster := typhoon.Cluster{}
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
			cluster := typhoon.Cluster{
				Config: configClusterName,
				Name:   clusterName,
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
		Subcommands: ClusterDeployCommands,
		Name:        "deploy",
		Usage:       "Typhoon deploy the cluster into env",
	},
	&cli.Command{
		Subcommands: ClusterGrafanaCommands,
		Name:        "grafana",
		Usage:       "Integration of Typhoon cluster with Grafana",
	},
	&cli.Command{
		Subcommands: ClusterDockerCommands,
		Name:        "docker",
		Usage:       "Integration of Typhoon cluster with Docker",
	},
	&cli.Command{
		Subcommands: ClusterGitCommands,
		Name:        "git",
		Usage:       "Integration of Typhoon cluster with Git",
	},
}

