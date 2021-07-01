package redis

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
	"typhoon-cli/src/typhoon/config"
)

type ServiceRedis struct {
	Config *config.Config
	connection *redis.Client
}

func (r *ServiceRedis) connect(service *config.ServiceRedis) bool {
	var ctx = context.Background()
	status := false
	redisString := fmt.Sprintf("%s:%d", service.Details.Host, service.Details.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisString, // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	r.connection = rdb
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		color.Red("%s", err)
		//os.Exit(1)
	} else {
		status = true
	}
	return status
}

func (r *ServiceRedis) TestConnect() bool {
	//color.Yellow("Run test connection to redis")
	projectConfig := r.Config
	status := false
	for _, service := range projectConfig.Services.Redis.Debug {
		status = r.connect(&service)
		//color.Green("Redis.Debug.%s: %t", service.Name, status)
	}

	return status
}

func (r *ServiceRedis) Set(key string, value interface{})  {
	var ctx = context.Background()
	_ = r.connection.Set(ctx, key, value, 0)
	//color.Yellow("%s", output)
}