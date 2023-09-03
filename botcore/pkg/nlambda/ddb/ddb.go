package ddb

import (
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	naws "github.com/telegram-serverless/botcore/pkg/nlambda/aws"
	"github.com/telegram-serverless/botcore/pkg/nlambda/base"
	"go.uber.org/zap"
)

//go:generate mockery --name=API --filename=api_mock.go
type API interface {
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, config TransactWriteItemsConfig) (*dynamodb.TransactWriteItemsOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type DDB struct {
	*base.Base
	client   *dynamodb.Client
	loadOnce sync.Once
}

func New(base *base.Base) *DDB {
	return &DDB{Base: base}
}

// New create a single (shared) instance of DDB service.
func (ddb *DDB) Load(ctx context.Context) error {
	const act string = "load ddb service"

	ddb.loadOnce.Do(func() {
		ddb.Log.Info("creating service", zap.String("service", "ddb"))

		cfg, err := naws.LoadConfigWithOpts(ctx)
		if err != nil {
			ddb.Log.Error(err.Error())
			return

		}

		awsv2.AWSV2Instrumentor(&cfg.APIOptions)
		ddb.client = dynamodb.NewFromConfig(cfg)

		ddb.Log.Debug("instance attributes",
			zap.String("name", "ddb"),
			zap.Any("attributes", ddb))
	})

	if ddb.client == nil {
		return errors.New("something went wrong")
	}

	return nil
}

func (ddb *DDB) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if err := ddb.Load(ctx); err != nil {
		return nil, err
	}
	return ddb.client.Query(ctx, params, optFns...)
}

func (ddb *DDB) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if err := ddb.Load(ctx); err != nil {
		return nil, err
	}
	return ddb.client.GetItem(ctx, params, optFns...)
}

func (ddb *DDB) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if err := ddb.Load(ctx); err != nil {
		return nil, err
	}
	return ddb.client.PutItem(ctx, params, optFns...)
}

func (ddb *DDB) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	if err := ddb.Load(ctx); err != nil {
		return nil, err
	}
	return ddb.client.UpdateItem(ctx, params, optFns...)
}

func (ddb *DDB) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	if err := ddb.Load(ctx); err != nil {
		return nil, err
	}
	return ddb.client.DeleteItem(ctx, params, optFns...)
}

type TransactWriteItemsConfig struct {
	AllowConditionalCheckFailed bool
}

func (ddb *DDB) TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, config TransactWriteItemsConfig) (*dynamodb.TransactWriteItemsOutput, error) {
	if err := ddb.Load(ctx); err != nil {
		return nil, err
	}

	out, err := ddb.client.TransactWriteItems(ctx, params)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (ddb *DDB) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	if err := ddb.Load(ctx); err != nil {
		return nil, err
	}
	return ddb.client.Scan(ctx, params, optFns...)
}
