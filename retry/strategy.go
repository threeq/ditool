package retry

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ------------------------------------------------------------------------------

type Stoped interface {
	Stop() <-chan time.Time
}

type limitStopStrategy int64

func (l *limitStopStrategy) Stop() <-chan time.Time {
	if *l <= 0 {
		return nil
	}
	*l = *l - 1
	return make(chan time.Time)
}

func NeverStopStrategy() Stoped {
	return StopAfterAttemptStrategy(math.MaxInt64)
}

func StopAfterAttemptStrategy(num int64) Stoped {
	st := new(limitStopStrategy)
	*st = limitStopStrategy(num)
	return st
}

type delayStopStrategy time.Duration

func (d *delayStopStrategy) Stop() <-chan time.Time {
	return time.After(time.Duration(*d))
}

func StopAfterDelayStrategy(t time.Duration) Stoped {
	st := new(delayStopStrategy)
	*st = delayStopStrategy(t)
	return st
}

// ------------------------------------------------------------------------------

type Waiting interface {
	Wait() <-chan time.Time
}

type delayWaitStrategy struct {
	waitNum int
	delay   time.Duration
}

type FixedWaitStrategy time.Duration

func (f FixedWaitStrategy) Wait() <-chan time.Time {
	return time.After(time.Duration(f))
}

type randomWS struct {
	min, max time.Duration
}

func (r *randomWS) Wait() <-chan time.Time {
	return time.After(time.Duration(rand.Int63n(int64(r.max-r.max))) + r.min)
}

func RandomWaitStrategy(min, max time.Duration) Waiting {
	return &randomWS{min: min, max: max}
}

type incWaitStrategy struct {
	curr, step time.Duration
}

func (i *incWaitStrategy) Wait() <-chan time.Time {
	i.curr += i.step
	return time.After(i.curr)
}

func IncrementingWaitStrategy(init, step time.Duration) Waiting {
	return &incWaitStrategy{curr: init - step, step: step}
}

type expWaitStrategy struct {
	curr time.Duration
	rate int64
}

func (e *expWaitStrategy) Wait() <-chan time.Time {
	e.curr = e.curr * time.Duration(e.rate)
	return time.After(e.curr)
}

func ExponentialWaitStrategy(init time.Duration, rate int64) Waiting {
	return &expWaitStrategy{curr: init / time.Duration(rate), rate: rate}
}

type fibWaitStrategy struct {
	n1 time.Duration
	n2 time.Duration
}

func FibonacciWaitStrategy(n1 time.Duration) Waiting {
	return &fibWaitStrategy{n1: n1}
}

func (f *fibWaitStrategy) Wait() <-chan time.Time {
	n := f.n1 + f.n2
	f.n2 = f.n1
	f.n1 = n
	return time.After(n)
}
