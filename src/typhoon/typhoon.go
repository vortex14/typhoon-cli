package typhoon

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"typhoon-cli/src/integrations/gitlab"
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
