package main

import (
	awsConfig "sample/go-gcp/amazon"
	"sample/go-gcp/handlers"
	"sample/go-gcp/inmate"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	awsConfig.Init()
	todoContext := awsConfig.Context
	dynamodb := awsConfig.Dynamo

	inmateSvc := inmate.NewService(todoContext, *dynamodb)
	inmateHandler := handlers.NewInmate(*inmateSvc)

	router := gin.Default()
	router.GET("/inmates", inmateHandler.GetInmates)

	router.POST("/inmate", inmateHandler.PutInmate)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.POST("/inmate/attempt", inmateHandler.Attempt)

	router.GET("/inmate/attempts/:inmateId", inmateHandler.GetInmateAttempts)

	router.Run()
}
