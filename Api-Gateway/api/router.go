package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/mahmud3253/Project/Api-Gateway/api/handlers/v1"
	"github.com/mahmud3253/Project/Api-Gateway/config"
	"github.com/mahmud3253/Project/Api-Gateway/pkg/logger"
	"github.com/mahmud3253/Project/Api-Gateway/services"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
}

// New ...
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
	})

	api := router.Group("/v1")
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/users/getbyid/:id", handlerV1.GetUserById)
	api.GET("/users/listuser", handlerV1.ListUser)
	api.PUT("/users/update/:id", handlerV1.UpdateUser)
	api.DELETE("/users/delete/:id", handlerV1.DeleteUser)

	return router
}
