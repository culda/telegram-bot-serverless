package nlambda

import "context"

type wrap[E, R any] struct {
	OnError           func(E, R, error) (R, error)
	DisableInitialLog bool
}

type wrapOption[E, R any] func(*wrap[E, R])

type (
	Handler[E, R any] func(context.Context, E) (R, error)
)

//nolint:nakedret
func LambdaWrap[E, R any](svc *Services, f Handler[E, R], opts ...wrapOption[E, R]) Handler[E, R] {
	return func(ctx context.Context, e E) (res R, err error) {
		var wrap wrap[E, R]

		for _, opt := range opts {
			opt(&wrap)
		}

		defer func() {
			// The eror returned from LambdaWrap is used as shutdown error code.
			// The error from `nerr` is still used in returning LambdaWrap.
			shutdownErr := err

			// OnError is allowed to change the response and error returned from
			// LambdaWrap.
			if err != nil && wrap.OnError != nil {
				res, err = wrap.OnError(e, res, err)
			}

			svc.Shutdown(res, shutdownErr)
		}()

		if err = svc.Start(ctx, e, svcDisableInitialLog(wrap.DisableInitialLog)); err != nil {
			return
		}

		res, err = f(ctx, e)
		return
	}
}
