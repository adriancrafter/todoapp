package am

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type (
	Task     func(ctx context.Context) error
	Teardown func(ctx context.Context) error
)

type (
	Supervisor interface {
		Core
		AddTasks(tasks ...Task)
		AddShutdownTasks(tasks ...Teardown)
		Wait() error
		Context() context.Context
		CancelFunc() context.CancelFunc
	}

	supervisor struct {
		*SimpleCore
		tasks    []Task
		teardown []Teardown
		ctx      context.Context
		cancel   context.CancelFunc
	}

	Effect func(opt *opt)

	opt struct {
		parentCtx context.Context
		notify    bool
	}
)

func NewSupervisor(name string, notify bool, opts []Option, effects ...Effect) Supervisor {
	opt := &opt{
		parentCtx: context.Background(),
		notify:    notify,
	}

	for _, apply := range effects {
		apply(opt)
	}

	sv := &supervisor{
		SimpleCore: NewCore(name, opts...),
		tasks:      []Task{},
		teardown:   []Teardown{},
	}

	sv.ctx, sv.cancel = context.WithCancel(opt.parentCtx)
	if opt.notify {
		sv.ctx, sv.cancel = signal.NotifyContext(sv.ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	}

	return sv
}

func (sv *supervisor) AddTasks(tasks ...Task) {
	sv.tasks = append(sv.tasks, tasks...)
}

func (sv *supervisor) AddShutdownTasks(tasks ...Teardown) {
	sv.teardown = append(sv.teardown, tasks...)
}

func (sv *supervisor) Wait() (err error) {
	eg, ctx := errgroup.WithContext(sv.ctx)
	eg.Go(sv.contextDone(ctx))

	for _, t := range sv.tasks {
		task := t
		eg.Go(func() error {
			return task(ctx)
		})
	}

	for _, tt := range sv.teardown {
		teardown := tt
		defer teardown(ctx)
	}

	return eg.Wait()
}

func (sv *supervisor) contextDone(ctx context.Context) func() error {
	return func() error {
		<-ctx.Done()
		sv.cancel()
		return nil
	}
}

func (sv *supervisor) Context() context.Context {
	return sv.ctx
}

func (sv *supervisor) CancelFunc() context.CancelFunc {
	return sv.cancel
}
