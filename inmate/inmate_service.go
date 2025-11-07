package inmate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
