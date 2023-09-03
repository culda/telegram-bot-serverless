package aws

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Additional error codes to retry on.
var extraErrorCodes = []string{ //nolint:gochecknoglobals
	(*ddbtypes.TransactionConflictException)(nil).ErrorCode(),
}

const backOffDelay = time.Second * 1

func LoadConfigWithOpts(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	// append retryer to any optFns an individual config loader may have
	optFns = append(optFns, config.WithRetryer(func() aws.Retryer {
		r := retry.AddWithMaxBackoffDelay(retry.NewStandard(), backOffDelay)

		r = retry.AddWithMaxAttempts(r, retry.DefaultMaxAttempts)
		r = retry.AddWithErrorCodes(r, extraErrorCodes...)

		return r
	}))

	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return aws.Config{}, errors.New("load default config")
	}

	return cfg, nil
}
