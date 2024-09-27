package temputil

import (
	"context"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type Activity[A, R any] func(ctx context.Context, arg A) (R, error)
type ActivityNoResult[A any] func(ctx context.Context, arg A) error

// RunSyncActivity executes an activity synchronously, blocking until the result is returned.
func RunSyncActivity[A, R any](ctx workflow.Context, a Activity[A, R], arg A, opts ...workflow.ActivityOptions) (result R, err error) {
	activityCtx := workflow.WithValue(ctx, "", "")

	for _, opt := range opts {
		activityCtx = workflow.WithActivityOptions(activityCtx, opt)
	}

	err = workflow.ExecuteActivity(ctx, a, arg).Get(ctx, &result)
	return result, err
}

// RunSyncActivityNoResult executes an activity synchronously without a result, blocking until the workflow execution ends.
func RunSyncActivityNoResult[A any](ctx workflow.Context, a ActivityNoResult[A], arg A, opts ...workflow.ActivityOptions) (err error) {
	activityCtx := workflow.WithValue(ctx, "", "")

	for _, opt := range opts {
		activityCtx = workflow.WithActivityOptions(activityCtx, opt)
	}

	err = workflow.ExecuteActivity(ctx, a, arg).Get(ctx, nil)
	return err
}

type ActivityFuture[R any] func(ctx workflow.Context) (R, error)
type ActivityFutureNoResult func(ctx workflow.Context) error

// RunAsyncActivity executes an activity asynchronously.
func RunAsyncActivity[A, R any](ctx workflow.Context, a Activity[A, R], arg A, opts ...workflow.ActivityOptions) ActivityFuture[R] {
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

// RunAsyncActivityNoResult executes an activity asynchronously without a result.
func RunAsyncActivityNoResult[A any](ctx workflow.Context, a ActivityNoResult[A], arg A, opts ...workflow.ActivityOptions) ActivityFutureNoResult {
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

type Workflow[A, R any] func(ctx workflow.Context, arg A) (R, error)
type WorkflowNoResult[A any] func(ctx workflow.Context, arg A) error

type WorkflowFuture[R any] func(ctx context.Context) (R, error)
type WorkflowFutureNoResult func(ctx context.Context) error

// RunSyncWorkflow executes a workflow synchronously, blocking until the workflow execution ends and the result is returned.
func RunSyncWorkflow[A, R any](ctx context.Context, tc client.Client, opts client.StartWorkflowOptions, w Workflow[A, R], arg A) (result R, err error) {
	run, err := tc.ExecuteWorkflow(ctx, opts, w, arg)
	if err != nil {
		return result, err
	}

	err = run.Get(ctx, &result)
	return result, err
}

// RunSyncWorkflowNoResult executes a workflow synchronously without a result, blocking until the workflow execution ends.
func RunSyncWorkflowNoResult[A any](ctx context.Context, tc client.Client, opts client.StartWorkflowOptions, w WorkflowNoResult[A], arg A) (err error) {
	run, err := tc.ExecuteWorkflow(ctx, opts, w, arg)
	if err != nil {
		return err
	}

	err = run.Get(ctx, nil)
	return err
}

// RunAsyncWorkflow executes a workflow asynchronously.
func RunAsyncWorkflow[A, R any](ctx context.Context, tc client.Client, opts client.StartWorkflowOptions, w Workflow[A, R], arg A) (WorkflowFuture[R], error) {
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

// RunAsyncWorkflowNoResult executes a workflow asynchronously without a result.
func RunAsyncWorkflowNoResult[A any](ctx context.Context, tc client.Client, opts client.StartWorkflowOptions, w WorkflowNoResult[A], arg A) (WorkflowFutureNoResult, error) {
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

// RunSyncChildWorkflow executes a child workflow synchronously, blocking until the workflow execution ends and the result is returned.
func RunSyncChildWorkflow[A, R any](ctx workflow.Context, opts workflow.ChildWorkflowOptions, w Workflow[A, R], arg A) (result R, err error) {
	run := workflow.ExecuteChildWorkflow(ctx, opts, w, arg)
	err = run.Get(ctx, &result)
	return result, err
}

// RunSyncChildWorkflowNoResult executes a child workflow synchronously without a result, blocking until the workflow execution ends.
func RunSyncChildWorkflowNoResult[A any](ctx workflow.Context, opts workflow.ChildWorkflowOptions, w WorkflowNoResult[A], arg A) (err error) {
	run := workflow.ExecuteChildWorkflow(ctx, opts, w, arg)
	err = run.Get(ctx, nil)
	return err
}

type ChildWorkflowFuture[R any] func(ctx workflow.Context) (R, error)
type ChildWorkflowFutureNoResult func(ctx workflow.Context) error

// RunAsyncChildWorkflow executes a child workflow asynchronously.
func RunAsyncChildWorkflow[A, R any](ctx workflow.Context, opts workflow.ChildWorkflowOptions, w Workflow[A, R], arg A) ChildWorkflowFuture[R] {
	run := workflow.ExecuteChildWorkflow(ctx, opts, w, arg)
	return func(ctx workflow.Context) (R, error) {
		var result R
		err := run.Get(ctx, &result)
		return result, err
	}
}

// RunAsyncChildWorkflowNoResult executes a child workflow asynchronously without a result.
func RunAsyncChildWorkflowNoResult[A any](ctx workflow.Context, opts workflow.ChildWorkflowOptions, w WorkflowNoResult[A], arg A) ChildWorkflowFutureNoResult {
	run := workflow.ExecuteChildWorkflow(ctx, opts, w, arg)
	return func(ctx workflow.Context) error {
		err := run.Get(ctx, nil)
		return err
	}
}
