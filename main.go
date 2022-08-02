package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"user-management-service/docs"
	"user-management-service/src/configuration"
	"user-management-service/src/controllers"
	"user-management-service/src/helpers"
	"user-management-service/src/middlewares"
	"user-management-service/src/services"
	"user-management-service/src/validators"
)

func main() {
	router := gin.New()
	docs.SwaggerInfo.Title = "user api"
	docs.SwaggerInfo.Description = "User Management Service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = ""

	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "@timestamp",
			log.FieldKeyMsg:  "message",
		},
	})
	logger.SetLevel(log.TraceLevel)

	file, err := os.OpenFile("out.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
	}
	defer file.Close()

	router.Use(middlewares.CORSMiddleware)

	config := configuration.NewConfig()

	connection := helpers.New(config.Database.ConnectionString)

	connection.ConnectDb()

	userValidator := validators.NewUserValidator(logger)

	userService := services.NewUserService(userValidator, logger)

	userController := controllers.NewUserController(userService, logger)

	user := router.Group("/users")
	{
		user.POST("", userController.AddUser)
		user.PATCH("", userController.UpdateUser)

		user.DELETE("/:id", func(context *gin.Context) {
			id := context.Param("id")

			userController.DeleteUser(context, id)
		})

		user.GET("/:id", func(context *gin.Context) {
			id := context.Param("id")

			userController.GetUser(context, id)
		})

		user.GET("", userController.GetAllUser)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
