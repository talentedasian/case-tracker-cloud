package main

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	awsConfig "sample/go-gcp/amazon"
	"sample/go-gcp/handlers"
	"sample/go-gcp/inmate"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	putInmateRequestsCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "inmate_requests_total",
			Help: "Total number of inmate API requests",
		},
	)

	putInmateErrorCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "put_inmate_error_total",
			Help: "Total number of put_inmate errors",
		},
	)

	putInmateWriteCapacityHistgoram = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name: "put_inmate_write_capacity",
			Help: "Measurement of write capacity units consumed during putInmate operations",
		},
	)
)

func main() {
	awsConfig.Init()
	todoContext := awsConfig.Context
	dynamodb := awsConfig.Dynamo

	inmateSvc := inmate.NewService(todoContext, *dynamodb)
	inmateHandler := handlers.NewInmate(*inmateSvc)

	router := gin.Default()
	router.GET("/inmates", inmateHandler.GetInmates)

	router.POST("/inmate", func(c *gin.Context) {
		putInmate(c, todoContext, dynamodb)
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Run()
}

func putInmate(c *gin.Context, todoContext context.Context, svc *dynamodb.Client) {
	putInmateRequestsCounter.Inc()

	var inmate inmate.Inmate
	if err := c.BindJSON(&inmate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inMateItem, avErr := attributevalue.MarshalMap(inmate)

	if avErr != nil {
		slog.Error("Failed to marshal inmate", "error", avErr)
		c.JSON(http.StatusInternalServerError, avErr.Error())
		return
	}

	slog.Info("Adding inmate", "inmate", inMateItem)

	putItemOut, err := svc.PutItem(todoContext, &dynamodb.PutItemInput{
		TableName:              aws.String("case_tracker"),
		Item:                   inMateItem,
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		ReturnValues:           types.ReturnValueAllOld,
	})

	if putItemOut != nil {
		slog.Info("Total consumred write capacity", "write capacity",
			strconv.FormatFloat(*putItemOut.ConsumedCapacity.CapacityUnits, 'f', 2, 64))

		slog.Info("Successfully added inmate", "inmate_id", inmate.ID)

		putInmateWriteCapacityHistgoram.Observe(*putItemOut.ConsumedCapacity.CapacityUnits)
	}

	if err != nil {
		putInmateErrorCounter.Inc()

		c.JSON(http.StatusInternalServerError, err.Error())
	}
}
