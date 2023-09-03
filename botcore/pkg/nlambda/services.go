package nlambda

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/telegram-serverless/botcore/pkg/nlambda/base"
	"github.com/telegram-serverless/botcore/pkg/nlambda/ddb"
	"go.uber.org/zap"
)

var (
	servicesOnce     sync.Once //nolint:gochecknoglobals
	servicesInstance *Services //nolint:gochecknoglobals
)

type Services struct {
	*base.Base
	Ddb ddb.API
}

func Instance() (*Services, error) {
	servicesOnce.Do(func() {
		base, err := base.Instance()
		if err != nil {
			fmt.Printf("unable to create base service: %v\n", err) //nolint:forbidigo
			return
		}

		ddb := ddb.New(base)

		servicesInstance = &Services{
			Base: base,
			Ddb:  ddb,
		}
	})

	if servicesInstance == nil {
		return nil, errors.New("unable to build nlambda services")
	}

	return servicesInstance, nil
}

type svcWrap struct {
	svcDisableInitialLog bool
}

type svcWrapOption func(*svcWrap)

func svcDisableInitialLog(disable bool) svcWrapOption {
	return func(w *svcWrap) {
		w.svcDisableInitialLog = disable
	}
}

func (s *Services) Start(ctx context.Context, event interface{}, opts ...svcWrapOption) error {
	const act = "start nlambda service"

	var (
		logger  *zap.Logger
		err     error
		svcWrap svcWrap
	)

	defer func() {
		if r := recover(); r != nil {
			det := fmt.Sprintf("encountered panic, unable to start service(s) - %v", r)
			err = errors.New(fmt.Sprintf("%s\n%s", det, act))
		}
	}()

	for _, opt := range opts {
		opt(&svcWrap)
	}

	s.Initialise(ctx)

	logger, err = base.InitLogger()
	if err != nil {
		return err
	}

	s.Log = logger.With(s.Lambdazapper.ContextValues(ctx)...)

	if svcWrap.svcDisableInitialLog {
		s.Log.Info("initial event message")
	} else {
		s.Log.Info("initial event message", zap.Any("event", event))
	}

	return err
}

func (s *Services) Shutdown(response interface{}, err error) {
	const (
		resField    = "response"
		shutdownMsg = "shutdown"
		skips       = 2
	)

	shutdownLogger := s.Log.WithOptions(zap.AddCallerSkip(skips))

	if err != nil {
		shutdownLogger.Error(shutdownMsg, zap.Any("err", err), zap.Any(resField, response))
	} else if response != nil {
		shutdownLogger.Info(shutdownMsg, zap.Any(resField, response))
	} else {
		shutdownLogger.Info(shutdownMsg)
	}

	// base close happens here only
	s.Close()

	if r := recover(); r != nil {
		fmt.Printf("shutdown with panic: %#v\n", r) //nolint:forbidigo
		debug.PrintStack()
		panic(r)
	}
}
