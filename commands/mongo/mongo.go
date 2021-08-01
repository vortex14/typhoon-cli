package mongo

import (
	Context "context"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	typhoon "github.com/vortex14/gotyphoon"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/services"
	"github.com/vortex14/gotyphoon/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
)

var Commands = []*cli.Command{
	{
		Name:  "show",
		Usage: "Show collections of the project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:    "group",
				Aliases: []string{"g"},
				Value:   "main",
				Usage:   "Services group mongo name.",
			},
		},
		Action: func(context *cli.Context) error {
			groupName := context.String("group")
			config := context.String("config")
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}

			projectServices := &services.Services{
				Project: project,
			}

			project.Services = projectServices

			projectServices.LoadMongoServices()

			var tableData [][]string
			header := []string{"№", "Group", "Database", "Collection", "Count"}
			mongoService := project.Services.Collections.Mongo[groupName]
			if mongoService == nil {
				color.Red("Mongo group %s not found for the project", groupName)
				os.Exit(1)
			}
			collections := mongoService.GetCollections()
			for i, collection := range collections {
				query := &interfaces.MongoQuery{
					Timeout:    5,
					Filter:     bson.D{},
					Context:    Context.TODO(),
					Database:   project.GetName(),
					Collection: collection.Name,
					Options:    &options.CountOptions{},
				}

				count := mongoService.GetCountDocuments(query)
				tableData = append(tableData, []string{
					strconv.Itoa(i + 1),
					groupName,
					project.GetName(),
					collection.Name,
					strconv.FormatInt(count, 10),
				})
			}

			u := utils.Utils{}
			u.RenderTableOutput(header, tableData)
			return nil

		},
	},
	{
		Name:  "export",
		Usage: "Export Data from collection of the project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:    "group",
				Aliases: []string{"g"},
				Value:   "main",
				Usage:   "Services group mongo name.",
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "test-output-mongo.json",
				Usage:   "Export file name.",
			},
			&cli.StringFlag{
				Name:    "collection",
				Aliases: []string{"i"},
				Value:   "test",
				Usage:   "Export collection name.",
			},
		},
		Action: func(context *cli.Context) error {
			groupName := context.String("group")
			config := context.String("config")
			pathProject, _ := os.Getwd()
			collectionName := context.String("collection")
			outFile := context.String("file")
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}

			projectServices := &services.Services{
				Project: project,
			}

			project.Services = projectServices

			projectServices.LoadMongoServices()
			mongoService := project.Services.Collections.Mongo[groupName]
			if mongoService == nil {
				color.Red("Mongo group %s not found for the project", groupName)
				os.Exit(1)
			}
			writer, file, count, err := mongoService.Export(project.GetName(), collectionName, outFile)
			defer file.Close()

			if err != nil {
				color.Red("%s", err.Error())
				os.Exit(1)
			}

			err = writer.Flush()
			if err != nil {
				color.Red("%s", err.Error())
				os.Exit(1)
			}

			var tableData [][]string
			header := []string{"Database", "Collection", "Exported", "File"}

			tableData = append(tableData, []string{
				project.GetName(),
				collectionName,
				strconv.FormatInt(count, 10),
				outFile,
			})

			u := utils.Utils{}
			u.RenderTableOutput(header, tableData)

			return nil

		},
	},
	{
		Name:  "import",
		Usage: "Import data of project to collection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:    "group",
				Aliases: []string{"g"},
				Value:   "main",
				Usage:   "Services group mongo name.",
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "test-output-mongo.json",
				Usage:   "Import file name.",
			},
			&cli.StringFlag{
				Name:    "collection",
				Aliases: []string{"i"},
				Value:   "import-test",
				Usage:   "Export collection name.",
			},
		},
		Action: func(context *cli.Context) error {
			groupName := context.String("group")
			config := context.String("config")
			pathProject, _ := os.Getwd()
			collectionName := context.String("collection")
			inputFile := context.String("file")
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}
			projectServices := &services.Services{
				Project: project,
			}
			project.Services = projectServices
			projectServices.LoadMongoServices()
			mongoService := project.Services.Collections.Mongo[groupName]
			if mongoService == nil {
				color.Red("Mongo group %s not found for the project", groupName)
				os.Exit(1)
			}

			err, imported := mongoService.Import(project.GetName(), collectionName, inputFile)
			if err != nil {
				color.Red("%s", err.Error())
				os.Exit(1)
			}

			var tableData [][]string
			header := []string{"Database", "Collection", "Imported", "File"}
			tableData = append(tableData, []string{
				project.GetName(),
				collectionName,
				strconv.FormatUint(imported, 10),
				inputFile,
			})
			u := utils.Utils{}
			u.RenderTableOutput(header, tableData)

			return nil

		},
	},
}
