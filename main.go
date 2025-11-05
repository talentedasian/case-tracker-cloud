package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"

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
	cfg, err := config.LoadDefaultConfig(todoContext, config.WithSharedConfigProfile("user-assume-role"), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	stsSvc := sts.NewFromConfig(cfg)
	identity, identityErr := stsSvc.GetCallerIdentity(todoContext, &sts.GetCallerIdentityInput{})

	if identityErr != nil {
		log.Fatalf("Unable to retrieve current identity: reason %s", identityErr.Error())
	}

	slog.Info("Current Identity Access", "identity", *identity.Arn)

	cfg.Credentials = stscreds.NewAssumeRoleProvider(stsSvc, "arn:aws:iam::407464631290:role/dynamodb_read_access")

	dynamodb := dynamodb.NewFromConfig(cfg)
	router := gin.Default()
	router.GET("/inmates", func(c *gin.Context) {
		getInmates(c, todoContext, dynamodb)
	})

	router.POST("/inmate", func(c *gin.Context) {
		putInmate(c, todoContext, dynamodb)
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
	ID       uint64 `dynamodbav:"inmate_id" json:"id"`
	LastName string `dynamodbav:"inmate_last_name" json:"last_name"`
	Gender   Gender `dynamodbav:"inmate_gender" json:"gender"`
}

func NewDecoderOptions() attributevalue.DecoderOptions {
	return attributevalue.DecoderOptions{
		UseEncodingUnmarshalers: true,
	}
}

func getInmates(c *gin.Context, todoContext context.Context, svc *dynamodb.Client) {
	scanOutput, err := svc.Scan(todoContext, &dynamodb.ScanInput{
		TableName: aws.String("case_tracker"),
	})
	if err != nil {
		log.Fatalf("failed to scan table, %v", err)
	}

	opts := func(o *attributevalue.DecoderOptions) {
		o.UseEncodingUnmarshalers = true
	}

	var inmates []Inmate

	if err := attributevalue.UnmarshalListOfMapsWithOptions(scanOutput.Items, &inmates, opts); err != nil {
		slog.Error("Failed to unmarshal inmates", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal inmates"})
		return
	}

	c.JSON(http.StatusOK, inmates)
}

func (g *Gender) UnmarshalText(text []byte) error {
	s := strings.ToLower(string(text))
	switch s {
	case "male", "1":
		*g = Male
	case "female", "0":
		*g = Female
	default:
		return fmt.Errorf("invalid gender: %s", text)
	}
	return nil
}

func putInmate(c *gin.Context, todoContext context.Context, svc *dynamodb.Client) {
	putInmateRequestsCounter.Inc()

	var inmate Inmate
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
