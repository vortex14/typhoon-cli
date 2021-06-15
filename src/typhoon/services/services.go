package services

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"typhoon-cli/src/integrations/mongo"
	"typhoon-cli/src/integrations/nsq"
	"typhoon-cli/src/integrations/redis"
	"typhoon-cli/src/interfaces"
)

type Services struct {
	Project interfaces.Project
}
type ServiceStatus struct {
	Name string
	Status bool
}
type Statuses []ServiceStatus

func (s *Services) RunTestServices() {
	services := Statuses{}

	projectConfig := s.Project.LoadConfig()

	redisService := redis.ServiceRedis{Config: &projectConfig.Config}
	redisStatus := redisService.TestConnect()
	services = append(services, ServiceStatus{Name: "Redis", Status: redisStatus})

	mongoService := mongo.ServiceMongo{Config: &projectConfig.Config}
	mongoStatus := mongoService.TestConnect()
	services = append(services, ServiceStatus{Name: "Mongo", Status: mongoStatus})


	nsqService := nsq.ServiceNSQ{Config: &projectConfig.Config}
	nsqStatus := nsqService.TestConnect()
	services = append(services, ServiceStatus{Name: "Nsq", Status: nsqStatus})

	//color.Yellow("Redis status %t", redisStatus)
	//color.Yellow("Mongo service: %t", mongoStatus)
	//color.Yellow("NSQ service: %t", nsqStatus)

	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)
	//lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
	cols, rows := 2, 4

	table.SetCell(0, 0,
		tview.NewTableCell("Service").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).SetMaxWidth(10))
	table.SetCell(0, 1,
		tview.NewTableCell("Status").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).SetMaxWidth(40))

	for r := 1; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			colorError := tcell.ColorRed
			service := services[r-1]
			var name string
			if c == 0 {
				name = service.Name
			} else {
				if !service.Status {
					color = colorError
				}

				name = fmt.Sprintf("%t", service.Status)
			}

			table.SetCell(r, c, tview.NewTableCell(name).SetTextColor(color))
		}
	}
	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}


}
