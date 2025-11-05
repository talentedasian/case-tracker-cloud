package aws

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	Context context.Context
	Config  config.Config
	Dynamo  *dynamodb.Client
)

func Init() {
	Context = context.TODO()
	cfg, err := config.LoadDefaultConfig(Context, config.WithSharedConfigProfile("user-assume-role"), config.WithRegion("ap-southeast-1"))
	if err != nil {
		slog.Error("unable to load SDK config", "error", err.Error())
		panic("Unable to Initialize AWS")
	}
	Config = &cfg

	stsSvc := sts.NewFromConfig(cfg)
	identity, identityErr := stsSvc.GetCallerIdentity(Context, &sts.GetCallerIdentityInput{})

	if identityErr != nil {
		slog.Error("Unable to retrieve current identity", "error", identityErr.Error())
		panic("Unable to Initialize AWS")
	}

	slog.Debug("Current Identity Access", "identity", *identity.Arn)

	cfg.Credentials = stscreds.NewAssumeRoleProvider(stsSvc, "arn:aws:iam::407464631290:role/dynamodb_read_access")

	Config = &cfg
	Dynamo = dynamodb.NewFromConfig(cfg)
}
