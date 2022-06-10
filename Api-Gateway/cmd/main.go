package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/mahmud3253/Project/Api-Gateway/api"
	"github.com/mahmud3253/Project/Api-Gateway/config"
	"github.com/mahmud3253/Project/Api-Gateway/pkg/logger"
	"github.com/mahmud3253/Project/Api-Gateway/services"
	rds "github.com/mahmud3253/Project/Api-Gateway/storage/redis"
)

func main() {
	//	var redisRepo repo.RedisRepositoryStorage
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api-gateway")

	pool := redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	redisRepo := rds.NewRedisRepo(&pool)

	serviceManager, err := services.NewServiceManager(&cfg)

	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		RedisRepo: redisRepo,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}
}
