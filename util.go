package temputil

import (
	"context"

	"go.temporal.io/sdk/client"
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

type wf[A, R any] func(ctx workflow.Context, arg A) (R, error)
type wfNoResult[A any] func(ctx workflow.Context, arg A) error

type workflowFuture[R any] func(ctx context.Context) (R, error)
type workflowFutureNoResult func(ctx context.Context) error

func RunSyncWorkflow[A, R any](ctx context.Context, tc client.Client, opts client.StartWorkflowOptions, w wf[A, R], arg A) (result R, err error) {
	run, err := tc.ExecuteWorkflow(ctx, opts, w, arg)
	if err != nil {
		return result, err
	}

	err = run.Get(ctx, &result)
	return result, err
}

func RunSyncWorkflowNoResult[A any](ctx context.Context, tc client.Client, opts client.StartWorkflowOptions, w wfNoResult[A], arg A) (err error) {
	run, err := tc.ExecuteWorkflow(ctx, opts, w, arg)
	if err != nil {
		return err
	}

	err = run.Get(ctx, nil)
	return err
}

func RunAsyncWorkflow[A, R any](ctx context.Context, tc client.Client, opts client.StartWorkflowOptions, w wf[A, R], arg A) (workflowFuture[R], error) {
	var result R

	run, err := tc.ExecuteWorkflow(ctx, opts, w, arg)
	if err != nil {
		return func(ctx context.Context) (R, error) {
			return result, err
		}, err
	}

	return func(ctx context.Context) (R, error) {
		err := run.Get(ctx, &result)
		return result, err
	}, nil
}

func RunAsyncWorkflowNoResult[A any](ctx context.Context, tc client.Client, opts client.StartWorkflowOptions, w wfNoResult[A], arg A) (workflowFutureNoResult, error) {
	run, err := tc.ExecuteWorkflow(ctx, opts, w, arg)
	if err != nil {
		return func(ctx context.Context) error {
			return nil
		}, err
	}

	return func(ctx context.Context) error {
		err := run.Get(ctx, nil)
		return err
	}, nil
}

func RunSyncChildWorkflow[A, R any](ctx workflow.Context, opts workflow.ChildWorkflowOptions, w wf[A, R], arg A) (result R, err error) {
	run := workflow.ExecuteChildWorkflow(ctx, opts, w, arg)
	err = run.Get(ctx, &result)
	return result, err
}

func RunSyncChildWorkflowNoResult[A any](ctx workflow.Context, opts workflow.ChildWorkflowOptions, w wfNoResult[A], arg A) (err error) {
	run := workflow.ExecuteChildWorkflow(ctx, opts, w, arg)
	err = run.Get(ctx, nil)
	return err
}

type childWorkflowFuture[R any] func(ctx workflow.Context) (R, error)
type childWorkflowFutureNoResult func(ctx workflow.Context) error

func RunAsyncChildWorkflow[A, R any](ctx workflow.Context, opts workflow.ChildWorkflowOptions, w wf[A, R], arg A) childWorkflowFuture[R] {
	run := workflow.ExecuteChildWorkflow(ctx, opts, w, arg)
	return func(ctx workflow.Context) (R, error) {
		var result R
		err := run.Get(ctx, &result)
		return result, err
	}
}

func RunAsyncChildWorkflowNoResult[A any](ctx workflow.Context, opts workflow.ChildWorkflowOptions, w wfNoResult[A], arg A) childWorkflowFutureNoResult {
	run := workflow.ExecuteChildWorkflow(ctx, opts, w, arg)
	return func(ctx workflow.Context) error {
		err := run.Get(ctx, nil)
		return err
	}
}
