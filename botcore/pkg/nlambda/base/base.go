package base

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/dougEfresh/lambdazap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/**
 * Setup Services and Init for Lambdas without Ledger Interactions
 */
var (
	baseOnce     sync.Once //nolint:gochecknoglobals
	baseInstance *Base     //nolint:gochecknoglobals
)

// Service describes the services that will be available to all lambda handlers.
type Base struct {
	Context      context.Context //nolint:containedctx
	Lambdazapper *lambdazap.LambdaLogContext
	baseLog      *zap.Logger
	Log          *zap.Logger
}

func (b *Base) Close() {
	b.Log.Debug("about to close base service")
	defer b.Log.Sync() //nolint:errcheck
}

func (b *Base) Initialise(ctx context.Context) {
	b.Context = ctx
	b.Log = b.baseLog.With(b.Lambdazapper.ContextValues(ctx)...)
}

func (b *Base) AddFieldsToLogger(fields ...zapcore.Field) {
	b.Log = b.Log.With(fields...)
}

func Instance() (*Base, error) {
	baseOnce.Do(func() {
		var logger *zap.Logger
		var err error

		if logger, err = InitLogger(); err != nil {
			fmt.Printf("logger init failed: %s\n", err) //nolint:forbidigo
			return
		}

		baseInstance = &Base{
			baseLog:      logger,
			Lambdazapper: lambdazap.New().With(lambdazap.AwsRequestID, lambdazap.FunctionName, lambdazap.InvokeFunctionArn),
		}
	})

	if baseInstance == nil {
		return nil, errors.New("unable to build base service for lambda")
	}

	return baseInstance, nil
}

func InitLogger() (*zap.Logger, error) {
	lcfg := zap.NewProductionConfig()

	logger, err := lcfg.Build()
	if err != nil {
		return nil, errors.New("initialise logger")
	}

	return logger, nil
}
