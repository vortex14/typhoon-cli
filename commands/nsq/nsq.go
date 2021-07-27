package nsq

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/vortex14/gotyphoon"
	"github.com/vortex14/gotyphoon/integrations/nsq"
	"github.com/vortex14/gotyphoon/interfaces"
	"os"
	"strconv"
	"strings"
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
		},
		Name: "test",
		Usage: "Test NSQ connection",
		Action: func(context *cli.Context) error {
			config := context.String("config")
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}

			nsqService := nsq.Service{Project: project}
			status := nsqService.Ping()
			color.Yellow(` Test NSQ connection:
			status: %t
			NSQd addresses: %+v
			NSQ LookupD: %+v 
`,
			status,
			project.Config.NsqdNodes,
			project.Config.NsqlookupdIP,
			)
			return nil
		},
	},
	{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"v"},
				Value:   "my producer",
				Usage:   "Producer name",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value: "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:    "task",
				Aliases: []string{"t"},
				Value: "task.json",
				Usage:   "Load task from `FILE`",
			},
			&cli.StringFlag{
				Name:    "message",
				Aliases: []string{"m"},
				Value: "message",
				Usage:   "Load json message from line",
			},
			&cli.StringFlag{
				Name:    "topic",
				Aliases: []string{"i"},
				Value: "test",
				Usage:   "topic name",
			},
			&cli.StringFlag{
				Name:    "channel",
				Aliases: []string{"ch"},
				Value: "tasks",
				Usage:   "Channel name",
			},

		},
		Name: "pub",
		Usage: "Pub message to NSQ",
		Action: func(context *cli.Context) error {
			config := context.String("config")
			task := context.String("task")
			message := context.String("message")
			topic := context.String("topic")
			name := context.String("name")
			channel := context.String("channel")
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}

			nsqService := nsq.Service{Project: project}
			status := nsqService.Ping()
			if !status {
				color.Red("Connection failed to NSQ")
				os.Exit(1)
			}

			nsqService.InitQueue(&interfaces.Queue{
				Channel:    channel,
				Topic:      topic,
				Writable:   true,
			})

			err := nsqService.Pub(name, topic, message)
			if err != nil {
				color.Red("%s", err.Error())
			}

			color.Yellow(` Test NSQ connection:
			status: %t
			NSQd addresses: %+v
			NSQ LookupD: %+v 
			topic: %s
			message: %s
			task: %s
`,
				status,
				project.Config.NsqdNodes,
				project.Config.NsqlookupdIP,
				topic,
				message,
				task,

			)

			nsqService.StopProducers()
			return nil
		},
	},
	{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"v"},
				Value:   "my producer",
				Usage:   "Producer name",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value: "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:    "task",
				Aliases: []string{"t"},
				Value: "task.json",
				Usage:   "Load task from `FILE`",
			},
			&cli.StringFlag{
				Name:    "message",
				Aliases: []string{"m"},
				Value: "message",
				Usage:   "Load json message from line",
			},
			&cli.StringFlag{
				Name:    "topic",
				Aliases: []string{"i"},
				Value: "test",
				Usage:   "topic name",
			},
			&cli.StringFlag{
				Name:    "count",
				Aliases: []string{"ct"},
				Value: "100",
				Usage:   "topic name",
			},

		},
		Name: "batch-pub",
		Usage: "Pub message to NSQ",
		Action: func(context *cli.Context) error {
			config := context.String("config")
			task := context.String("task")
			message := context.String("message")
			topic := context.String("topic")
			name := context.String("name")
			count, err := strconv.Atoi(context.String("count"))
			if err != nil {
				color.Red("%s", err.Error())
				os.Exit(1)
			}
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}

			nsqService := nsq.Service{Project: project}
			status := nsqService.Ping()
			if !status {
				color.Red("Connection failed to NSQ")
				os.Exit(1)
			}
			settings := &interfaces.Queue{}
			settings.SetGroupName(name)
			nsqService.InitQueue(settings)
			iterCount := 0
			for iterCount <= count {
				iterCount += 1
				err := nsqService.Pub(name, topic, message)
				if err != nil {
					color.Red("%s", err.Error())
				}
			}


			color.Yellow(` Test NSQ connection:
			status: %t
			NSQd addresses: %+v
			NSQ LookupD: %+v 
			topic: %s
			message: %s
			task: %s
`,
				status,
				project.Config.NsqdNodes,
				project.Config.NsqlookupdIP,
				topic,
				message,
				task,

			)

			nsqService.StopProducers()
			return nil
		},
	},
	{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value: "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:    "channel",
				Aliases: []string{"ch"},
				Value: "tasks",
				Usage:   "Channel name",
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value: "reader",
				Usage:   "Reader name",
			},
			&cli.StringFlag{
				Name:    "topic",
				Aliases: []string{"i"},
				Value: "test",
				Usage:   "topic name",
			},

		},
		Name: "sub",
		Usage: "Sub messages from NSQ",
		Action: func(context *cli.Context) error {
			config := context.String("config")
			topic := context.String("topic")
			name := context.String("name")
			channel := context.String("channel")
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}

			nsqService := nsq.Service{Project: project}
			status := nsqService.Ping()
			if !status {
				color.Red("Connection failed to NSQ")
				os.Exit(1)
			}

			setting := &interfaces.Queue{Topic: topic, Concurrent: 1, Channel: channel}
			setting.SetGroupName(name)

			consumer := nsqService.InitConsumer(setting)
			var count int
			count = 0
			for msg := range consumer.Messages() {

				color.Yellow("%s", msg.Body)
				msg.Finish()
				count += 1
			}

			color.Yellow(` Test NSQ connection:
			status: %t
			NSQd addresses: %+v
			NSQ LookupD: %+v 
			topic: %s
			read messages: %d	
`,
				status,
				project.Config.NsqdNodes,
				project.Config.NsqlookupdIP,
				topic,
				count,

			)

			nsqService.StopProducers()
			return nil
		},
	},
	{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value: "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:    "channel",
				Aliases: []string{"ch"},
				Value: "tasks",
				Usage:   "Channel name",
			},
			&cli.StringFlag{
				Name:    "topics",
				Aliases: []string{"i"},
				Value: "test,test1,test2",
				Usage:   "topics name with , delimiter",
			},

		},
		Name: "group-sub",
		Usage: "Subscribe messages from a few topics from NSQ. Only Stream",
		Action: func(context *cli.Context) error {
			config := context.String("config")
			topics := context.String("topics")
			topicsArr := strings.Split(topics, ",")
			channel := context.String("channel")
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}

			nsqService := nsq.Service{Project: project}
			status := nsqService.Ping()
			if !status {
				color.Red("Connection failed to NSQ")
				os.Exit(1)
			}

			color.Yellow(`count topics for subscriptions: %d
			topics: %+v
			channel: %s
			`,
			len(topicsArr),
			topicsArr,
			channel,
			)

			for i, topic := range topicsArr {
				setting := &interfaces.Queue{Topic: topic, Channel: channel, Concurrent: 1}
				setting.SetGroupName("reader-"+strconv.Itoa(i))
				nsqService.InitConsumer(setting)
			}

			total := 0
			for yield := range nsqService.Read() {
				total += 1


					color.Green(`
				
				From Topic: %s
				Channel: %s
				Body: %s
				Name: %s
				Total: %d
				`,
						yield.Topic,
						yield.Channel,
						string(yield.Msg.Body),
						yield.Name,
						total,
					)


				yield.Msg.Finish()
			}
			nsqService.StopConsumers()
			return nil
		},
	},
	{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value: "config.local.yaml",
				Usage:   "Load configuration from `FILE`",
			},
			&cli.StringFlag{
				Name:    "channel",
				Aliases: []string{"ch"},
				Value: "tasks",
				Usage:   "Channel name",
			},
			&cli.StringFlag{
				Name:    "topics",
				Aliases: []string{"i"},
				Value: "test,test1,test2",
				Usage:   "topics name with , delimiter",
			},
			&cli.StringFlag{
				Name:    "concurrent",
				Aliases: []string{"ct"},
				Value: "1",
				Usage:   "Concurrent message from Queue",
			},

		},
		Name: "group-batch-sub",
		Usage: "Subscribe messages from a few topics from NSQ. Only Batch",
		Action: func(context *cli.Context) error {
			config := context.String("config")
			topics := context.String("topics")
			concurrentStr := context.String("concurrent")
			concurrent, err := strconv.Atoi(concurrentStr)
			if err != nil {
				color.Red("%s", err.Error())
				os.Exit(1)
			}
			topicsArr := strings.Split(topics, ",")
			channel := context.String("channel")
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}

			nsqService := nsq.Service{Project: project}
			status := nsqService.Ping()
			if !status {
				color.Red("Connection failed to NSQ")
				os.Exit(1)
			}

			color.Yellow(`count topics for subscriptions: %d
			topics: %+v
			channel: %s
			`,
				len(topicsArr),
				topicsArr,
				channel,
			)

			for i, topic := range topicsArr {
				setting := &interfaces.Queue{
					Topic: topic,
					Concurrent: concurrent,
					Channel: channel,
				}
				setting.SetGroupName("reader-"+strconv.Itoa(i))
				nsqService.InitConsumer(setting)
			}


			//nsqService.BatchRead()


			nsqService.StopConsumers()






			//			consumer := nsqService.InitConsumer(name, topic, channel, 1)
			//			var count int
			//
			//
			//			count = 0
			//			for msg := range consumer.Messages() {
			//
			//				color.Yellow("%s", msg.Body)
			//				msg.Finish()
			//				count += 1
			//			}
			//
			//			color.Yellow(` Test NSQ connection:
			//			status: %t
			//			NSQd addresses: %+v
			//			NSQ LookupD: %+v
			//			topic: %s
			//			read messages: %d
			//`,
			//				status,
			//				project.Config.Config.NsqdNodes,
			//				project.Config.Config.NsqlookupdIP,
			//				topic,
			//				count,
			//
			//			)
			//
			//			nsqService.StopProducers()
			return nil
		},
	},

}





