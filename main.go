package main

import (
	"context"
	"encoding/json"
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

	router.Run()
}

type Inmate struct {
	ID          uint64    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
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
	var inmate Inmate
	if err := c.BindJSON(&inmate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := map[string]types.AttributeValue{
		"inmate_id":        &types.AttributeValueMemberN{Value: strconv.FormatUint(inmate.ID, 10)},
		"inmate_last_name": &types.AttributeValueMemberS{Value: inmate.LastName},
	}

	slog.Info("Adding inmate", "inmate", item)

	putItemOut, err := svc.PutItem(todoContext, &dynamodb.PutItemInput{
		TableName:              aws.String("case_tracker"),
		Item:                   item,
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		ReturnValues:           types.ReturnValueAllOld,
	})

	if putItemOut != nil {
		jsonByteData, jErr := json.Marshal(&putItemOut.Attributes)

		slog.Info("Total consumred write capacity", "write capacity",
			strconv.FormatFloat(*putItemOut.ConsumedCapacity.CapacityUnits, 'f', 2, 64))

		if jErr != nil {
			slog.Error("Failed to marshal put item output attributes", "error", jErr.Error())
		}

		slog.Info("Successfully added inmate", "inmate", string(jsonByteData))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}
