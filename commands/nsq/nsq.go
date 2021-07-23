package nsq

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/vortex14/gotyphoon"
	"github.com/vortex14/gotyphoon/integrations/nsq"
	"os"
	"strconv"
	"strings"
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
			status := nsqService.TestConnect()
			color.Yellow(` Test NSQ connection:
			status: %t
			NSQd addresses: %+v
			NSQ LookupD: %+v 
`,
			status,
			project.Config.Config.NsqdNodes,
			project.Config.Config.NsqlookupdIP,
			)
			return nil
		},
	},
	&cli.Command{
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

		},
		Name: "pub",
		Usage: "Pub message to NSQ",
		Action: func(context *cli.Context) error {
			config := context.String("config")
			task := context.String("task")
			message := context.String("message")
			topic := context.String("topic")
			name := context.String("name")
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}

			nsqService := nsq.Service{Project: project}
			status := nsqService.TestConnect()
			if !status {
				color.Red("Connection failed to NSQ")
				os.Exit(1)
			}

			nsqService.InitProducer(name)

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
				project.Config.Config.NsqdNodes,
				project.Config.Config.NsqlookupdIP,
				topic,
				message,
				task,

			)

			nsqService.StopProducers()
			return nil
		},
	},
	&cli.Command{
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

		},
		Name: "batch-pub",
		Usage: "Pub message to NSQ",
		Action: func(context *cli.Context) error {
			config := context.String("config")
			task := context.String("task")
			message := context.String("message")
			topic := context.String("topic")
			name := context.String("name")
			pathProject, _ := os.Getwd()
			project := &typhoon.Project{
				ConfigFile: config,
				Path:       pathProject,
			}

			nsqService := nsq.Service{Project: project}
			status := nsqService.TestConnect()
			if !status {
				color.Red("Connection failed to NSQ")
				os.Exit(1)
			}

			nsqService.InitProducer(name)

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
				project.Config.Config.NsqdNodes,
				project.Config.Config.NsqlookupdIP,
				topic,
				message,
				task,

			)

			nsqService.StopProducers()
			return nil
		},
	},
	&cli.Command{
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
			status := nsqService.TestConnect()
			if !status {
				color.Red("Connection failed to NSQ")
				os.Exit(1)
			}

			consumer := nsqService.InitConsumer(name, topic, channel, 1)
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
				project.Config.Config.NsqdNodes,
				project.Config.Config.NsqlookupdIP,
				topic,
				count,

			)

			nsqService.StopProducers()
			return nil
		},
	},
	&cli.Command{
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
			status := nsqService.TestConnect()
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
				nsqService.InitConsumer("reader-"+strconv.Itoa(i), topic, channel, 1)
			}


			nsqService.Read()


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
	&cli.Command{
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
		Name: "group-batch-sub",
		Usage: "Subscribe messages from a few topics from NSQ. Only Batch",
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
			status := nsqService.TestConnect()
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
				nsqService.InitConsumer("reader-"+strconv.Itoa(i), topic, channel, 1)
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





