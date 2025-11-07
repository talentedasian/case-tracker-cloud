package inmate

import (
	"context"
	"fmt"
	"log"
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

func (svc *InmateService) GetInmates() []Inmate {
	inmatesScanned, err := svc.dynamodb.Scan(svc.context, &dynamodb.ScanInput{
		TableName: aws.String("case_tracker"),
	})

	if err != nil {
		log.Fatalf("failed to scan table, %v", err)
	}

	opts := func(o *attributevalue.DecoderOptions) {
		o.UseEncodingUnmarshalers = true
	}

	var inmates []Inmate

	if err := attributevalue.UnmarshalListOfMapsWithOptions(inmatesScanned.Items, &inmates, opts); err != nil {
		slog.Error("Failed to unmarshal inmates", "error", err)
		log.Fatalf("Could not unmarshal scan response")
	}

	return inmates
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
