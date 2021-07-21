package services

import (
	"strconv"
	"typhoon-cli/src/integrations/mongo"
	"typhoon-cli/src/integrations/nsq"
	"typhoon-cli/src/integrations/redis"
	"typhoon-cli/src/interfaces"
	"typhoon-cli/src/utils"
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
	projectConfig := s.Project.LoadConfig()
	var tableData [][]string

	redisService := redis.Service{Config: &projectConfig.Config}
	mongoService := mongo.Service{Config: &projectConfig.Config}
	nsqService := nsq.Service{Config: &projectConfig.Config}

	services := Statuses{
		ServiceStatus{Name: "Redis", Status: redisService.TestConnect()},
		ServiceStatus{Name: "Mongo", Status: mongoService.TestConnect()},
		ServiceStatus{Name: "Nsq", Status: nsqService.TestConnect()},
	}

	u := utils.Utils{}

	header := []string{"Name", "Status"}

	for _, service := range services {
		tableData = append(tableData, []string{service.Name, strconv.FormatBool(service.Status)})
	}

	u.RenderTableOutput(header, tableData)
}
