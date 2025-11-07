package inmate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	putInmateWriteCapacityHistgoram = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name: "put_inmate_write_capacity",
			Help: "Measurement of write capacity units consumed during putInmate operations",
		},
	)

	putInmateErrorCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "put_inmate_error_total",
			Help: "Total number of put_inmate errors",
		},
	)
)

type InmateService struct {
	context  context.Context
	dynamodb dynamodb.Client
}

func NewService(ctx context.Context, dynamodbClient dynamodb.Client) *InmateService {
	return &InmateService{ctx, dynamodbClient}
}

func (svc *InmateService) GetInmates() ([]Inmate, error) {
	var inmates []Inmate

	inmatesScanned, err := svc.dynamodb.Scan(svc.context, &dynamodb.ScanInput{
		TableName: aws.String("case_tracker"),
	})

	if err != nil {
		slog.Error("failed to scan table", "err", err.Error())

		return inmates, errors.New("scanning of table failed")
	}

	opts := func(o *attributevalue.DecoderOptions) {
		o.UseEncodingUnmarshalers = true
	}

	if err := attributevalue.UnmarshalListOfMapsWithOptions(inmatesScanned.Items, &inmates, opts); err != nil {
		slog.Error("Failed to unmarshal inmates", "error", err)

		return inmates, errors.New("failed to parse dynamodb response")
	}

	return inmates, nil
}

func (svc *InmateService) PutInmate(inmate Inmate) error {
	inmateItem, avErr := attributevalue.MarshalMap(inmate)

	if avErr != nil {
		slog.Error("Failed to marshal inmate", "error", avErr)
		return errors.New("failed to parse inmate into dynamodb request")
	}

	slog.Info("Adding inmate", "inmate", inmateItem)

	putItemOut, err := svc.dynamodb.PutItem(svc.context, &dynamodb.PutItemInput{
		TableName:              aws.String("case_tracker"),
		Item:                   inmateItem,
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
		slog.Error("failed to put item", "err", err.Error(), "item", inmateItem)

		return errors.New("failed to write to dynamodb item")
	}

	return nil
}

func NewDecoderOptions() attributevalue.DecoderOptions {
	return attributevalue.DecoderOptions{
		UseEncodingUnmarshalers: true,
	}
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
