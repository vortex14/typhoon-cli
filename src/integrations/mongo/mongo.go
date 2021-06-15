package mongo

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
	"typhoon-cli/src/typhoon/config"
)

type ServiceMongo struct {
	Config *config.Config
}

func (m *ServiceMongo) connect(service *config.ServiceMongo) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	status := false
	connectionString := fmt.Sprintf("mongodb://%s:%d", service.Details.Host, service.Details.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	errPing := client.Ping(ctx, readpref.Primary())

	if errPing != nil || err != nil {
		color.Red("%s", errPing)
	} else {
		status = true
	}
	return status
}

func (m *ServiceMongo) TestConnect() bool {
	projectConfig := m.Config
	status := false
	for _, service := range projectConfig.Services.Mongo.Debug {
		status = m.connect(&service)
	}
	return status
}