package temputil

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

type activity[A, R any] func(ctx context.Context, arg A) (R, error)
type activityNoResult[A any] func(ctx context.Context, arg A) error

func RunSyncActivity[A, R any](ctx workflow.Context, a activity[A, R], arg A, opts ...workflow.ActivityOptions) (result R, err error) {
	activityCtx := workflow.WithValue(ctx, "", "")

	for _, opt := range opts {
		activityCtx = workflow.WithActivityOptions(activityCtx, opt)
	}

	err = workflow.ExecuteActivity(ctx, a, arg).Get(ctx, &result)
	return result, err
}

func RunSyncActivityNoResult[A any](ctx workflow.Context, a activityNoResult[A], arg A, opts ...workflow.ActivityOptions) (err error) {
	activityCtx := workflow.WithValue(ctx, "", "")

	for _, opt := range opts {
		activityCtx = workflow.WithActivityOptions(activityCtx, opt)
	}

	err = workflow.ExecuteActivity(ctx, a, arg).Get(ctx, nil)
	return err
}

type activityFuture[R any] func(ctx workflow.Context) (R, error)
type activityFutureNoResult func(ctx workflow.Context) error

func RunAsyncActivity[A, R any](ctx workflow.Context, a activity[A, R], arg A, opts ...workflow.ActivityOptions) activityFuture[R] {
	activityCtx := workflow.WithValue(ctx, "", "")

	for _, opt := range opts {
		activityCtx = workflow.WithActivityOptions(activityCtx, opt)
	}

	future := workflow.ExecuteActivity(ctx, a, arg)

	return func(ctx workflow.Context) (result R, err error) {
		err = future.Get(ctx, &result)
		return result, err
	}
}

func RunAsyncActivityNoResult[A any](ctx workflow.Context, a activityNoResult[A], arg A, opts ...workflow.ActivityOptions) activityFutureNoResult {
	activityCtx := workflow.WithValue(ctx, "", "")

	for _, opt := range opts {
		activityCtx = workflow.WithActivityOptions(activityCtx, opt)
	}

	future := workflow.ExecuteActivity(ctx, a, arg)

	return func(ctx workflow.Context) (err error) {
		err = future.Get(ctx, nil)
		return err
	}
}
