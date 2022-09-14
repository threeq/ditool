package retry

import (
	"errors"
	"time"
)

var (
	ErrTimeout = errors.New("retry timeout")
)

type Executor struct {
	wait    Waiting
	stop    Stoped
	trigger []Trigger
}

func (et *Executor) Run(task func() (any, error)) error {
	running := func() (res any, err error, ex any) {
		defer func() {
			if ex = recover(); ex != nil {

			}
		}()
		res, err = task()
		return
	}

	res, err, ex := running()
	//  停止检查
	for {
		if et.triggerCheck(res, err, ex) {
			s := et.stop.Stop()
			if s == nil {
				return ErrTimeout
			}
			select {
			case <-s:
				return ErrTimeout
			case <-et.wait.Wait():
				res, err, ex = running()
			}
		} else {
			// 已经成功，不需重试
			return nil
		}
	}

}

func (et *Executor) RunAsync(task func() (any, error)) {
	// TODO implement me
	panic("implement me")
}

func (et *Executor) triggerCheck(res any, err error, ex any) bool {
	for _, t := range et.trigger {
		if t.Check(res, err, ex) {
			return true
		}
	}
	return false
}

func NewExecutor(options ...Option) *Executor {
	et := &Executor{
		stop: StopAfterAttemptStrategy(0),
		wait: FixedWaitStrategy(time.Duration(1)),
	}
	for _, option := range options {
		option(et)
	}
	return et
}

type Option func(template *Executor)

func IfResult(checker func(res any) bool) Option {
	return func(template *Executor) {
		template.trigger = append(template.trigger, &checkResult{
			checker: checker,
		})
	}
}

func IfError(checker func(err error) bool) Option {
	return func(b *Executor) {
		b.trigger = append(b.trigger, &checkError{
			checker: checker,
		})
	}
}

func IfPanic(checker func(ex any) bool) Option {
	return func(b *Executor) {
		b.trigger = append(b.trigger, &checkPanic{
			checker: checker,
		})
	}
}

func WaitStrategy(wait Waiting) Option {
	return func(b *Executor) {
		b.wait = wait
	}
}

func StopStrategy(stop Stoped) Option {
	return func(b *Executor) {
		b.stop = stop
	}
}

// ----------------------------------------------------------------------

type Trigger interface {
	Check(res any, err error, ex any) bool
}

type checkResult struct {
	checker func(any) bool
}

func (receiver *checkResult) Check(res any, err error, ex any) bool {
	return receiver.checker(res)
}

type checkError struct {
	checker func(error) bool
}

func (receiver *checkError) Check(res any, err error, ex any) bool {
	return receiver.checker(err)
}

type checkPanic struct {
	checker func(any) bool
}

func (receiver *checkPanic) Check(res any, err error, ex any) bool {
	return receiver.checker(ex)
}

//-------------------------------------------------------------------

func Running(task func() (any, error), options ...Option) error {
	executor := NewExecutor(options...)
	return executor.Run(task)
}

func RunningAsync(task func() (any, error), options ...Option) {
	executor := NewExecutor(options...)
	executor.RunAsync(task)
}
