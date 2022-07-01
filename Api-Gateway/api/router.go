package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mahmud3253/Project/Api-Gateway/api/auth"
	"github.com/mahmud3253/Project/Api-Gateway/api/casbin"
	v1 "github.com/mahmud3253/Project/Api-Gateway/api/handlers/v1"
	"github.com/mahmud3253/Project/Api-Gateway/config"
	"github.com/mahmud3253/Project/Api-Gateway/pkg/logger"
	"github.com/mahmud3253/Project/Api-Gateway/services"
	"github.com/mahmud3253/Project/Api-Gateway/storage/repo"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	casbinN "github.com/casbin/casbin/v2"
	_ "github.com/mahmud3253/Project/Api-Gateway/api/docs"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	RedisRepo      repo.RedisRepositoryStorage
	Casbin         *casbinN.Enforcer

}

// New @BasePath /v1
// New ...
// @SecurityDefinitions.apikey BearerAuth
// @Description GetMyProfile
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	jwtHandler := &auth.JwtHandler{
		SigninKey: option.Conf.SigninKey,
		Log:       option.Logger,
	}

	router.Use(casbin.NewJwtRoleStruct(option.Casbin, option.Conf, *jwtHandler))

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Redis:          option.RedisRepo,
	})

	api := router.Group("/v1")
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/users/getbyid/:id", handlerV1.GetUserById)
	api.GET("/users/listuser", handlerV1.ListUser)
	api.PUT("/users/update/:id", handlerV1.UpdateUser)
	api.DELETE("/users/delete/:id", handlerV1.DeleteUser)
	api.POST("/users/register", handlerV1.RegisterUser)
	api.POST("/users/verfication", handlerV1.VerifyUser)
	api.GET("/users/login/:email/:password", handlerV1.Login)
	api.GET("/users/loginbyauth", handlerV1.LoginByAuth)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
