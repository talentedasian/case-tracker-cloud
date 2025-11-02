package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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
	todoContext := context.TODO()
	cfg, err := config.LoadDefaultConfig(todoContext, config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	router := gin.Default()
	router.GET("/inmates", func(c *gin.Context) {
		getInmates(c, todoContext, svc)
	})

	router.POST("/inmate", func(c *gin.Context) {
		putInmate(c, todoContext, svc)
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Run()
}

type Gender uint8

const (
	Female Gender = 0
	Male   Gender = 1
)

type Inmate struct {
	ID          uint64    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Gender      Gender    `json:"gender"`
	MiddleName  string    `json:"middle_name"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	PhoneNumber string    `json:"phone_number"`
}

func getInmates(c *gin.Context, todoContext context.Context, svc *dynamodb.Client) {
	scanOutput, err := svc.Scan(todoContext, &dynamodb.ScanInput{
		TableName: aws.String("case_tracker"),
	})
	if err != nil {
		log.Fatalf("failed to scan table, %v", err)
	}

	c.JSON(http.StatusOK, scanOutput.Items)
}

func putInmate(c *gin.Context, todoContext context.Context, svc *dynamodb.Client) {
	putInmateRequestsCounter.Inc()

	var inmate Inmate
	if err := c.BindJSON(&inmate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := map[string]types.AttributeValue{
		"inmate_id":        &types.AttributeValueMemberN{Value: strconv.FormatUint(inmate.ID, 10)},
		"inmate_last_name": &types.AttributeValueMemberS{Value: inmate.LastName},
		"inmate_gender":    &types.AttributeValueMemberS{Value: inmate.LastName},
	}

	slog.Info("Adding inmate", "inmate", item)

	putItemOut, err := svc.PutItem(todoContext, &dynamodb.PutItemInput{
		TableName:              aws.String("case_tracker"),
		Item:                   item,
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
