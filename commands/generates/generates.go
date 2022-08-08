package generates

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	typhoon "github.com/vortex14/gotyphoon"
	"github.com/vortex14/gotyphoon/elements/models/bar"
	"github.com/vortex14/gotyphoon/elements/models/timer"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)

var Commands = []*cli.Command{
	{
		Name:  "product",
		Usage: "Generate a new product",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: func(context *cli.Context) error {
			var f fake.Product
			err := gofakeit.Struct(&f)
			if err != nil {
				color.Red("%s", err.Error())
				return err
			}
			u := utils.Utils{}
			dump := u.PrintPrettyJson(f)
			color.Green("%s", dump)
			return nil
		},
	},
	{
		Name:  "task",
		Usage: "Generate a new task",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:    "interval",
				Aliases: []string{"i"},
				Value:   "1000",
				Usage:   "Millisecond gen interval",
			},
			&cli.StringFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Value:   "1",
				Usage:   "Time until self-exit in minutes",
			},
			&cli.StringFlag{
				Name:    "component",
				Aliases: []string{"ct"},
				Value:   "Fetcher",
				Usage:   "task for component",
			},
			&cli.StringFlag{
				Name:    "priorities",
				Aliases: []string{"ps"},
				Value:   "1,2,3",
				Usage:   "list of priorities int",
			},
		},
		Action: func(context *cli.Context) error {
			configFile := context.String("config")
			tickS := context.String("interval")
			timeoutS := context.String("timeout")
			timeout, err := strconv.Atoi(timeoutS)
			component := context.String("component")
			priorities := context.String("priorities")

			prioritiesList := strings.Split(priorities, ",")
			if len(prioritiesList) == 0 {
				color.Red("Priorities undefined")
				os.Exit(1)
			}
			u := utils.Utils{}

			prioritiesArr := u.ConvertStringListToIntList(prioritiesList)

			if err != nil {
				color.Red("%s", err.Error())
				return err
			}

			tick, err := strconv.Atoi(tickS)
			if err != nil {
				color.Red("%s", err.Error())
				return err
			} else if tick <= 0 {
				color.Red("--timeout > 0")
				return nil
			}
			project := &typhoon.Project{
				ConfigFile:        configFile,
				SelectedComponent: []string{component},
			}
			project.LoadConfig()

			project.LoadServices(
				interfaces.TyphoonIntegrationsOptions{
					NSQ: interfaces.MessageBrokerOptions{
						Active:          true,
						EnabledConsumer: true,
						EnabledProducer: true,
					},
				},
			)
			if !project.Services.Collections.Nsq.Ping() {
				color.Red("No ping to NSQ")
				os.Exit(1)
			}
			project.Services.Collections.Nsq.Options = interfaces.MessageBrokerOptions{
				EnabledProducer: true,
				EnabledConsumer: false,
			}
			project.Services.RunNSQ()

			count := 0
			bar := bar.Bar{}
			bar.NewOption(0, -1)
			tickerGen := timer.SetInterval(func(args ...interface{}) {

				var f fake.Product
				err = gofakeit.Struct(&f)
				if err != nil {
					color.Red("%s", err.Error())
					return
				}

				fakeTask, _ := fake.CreateFakeTask(interfaces.FakeTaskOptions{
					UserAgent: true,
					Cookies:   true,
					Proxy:     false,
				})

				dump := u.PrintPrettyJson(fakeTask)

				count += 1
				bar.Play(int64(count), "Total generated tasks")
				rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
				priority := prioritiesArr[rand.Intn(len(prioritiesArr))]
				project.Services.Collections.Nsq.PriorityPub(
					priority,
					interfaces.PROCESSOR2PRIORITY,
					dump)
				
			}, tick)

			timer.SetTimeout(func(args ...interface{}) {
				tickerGen.Stop()
			}, timeout*60*1000)

			tickerGen.Await()

			return nil
		},
	},
}
