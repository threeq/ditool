package retry

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExecuteTemplate_Run(t *testing.T) {

	t.Run("正常调用", func(t *testing.T) {
		runningNum := 0
		Running(func() (any, error) {
			runningNum += 1
			return nil, nil
		})
		assert.Equal(t, 1, runningNum)

	})

	t.Run("默认不重试", func(t *testing.T) {
		runningNum := 0

		Running(func() (any, error) {
			runningNum += 1
			return nil, errors.New("112")
		}, IfError(func(err error) bool {
			return err != nil
		}))
		assert.Equal(t, 1, runningNum)
	})

	t.Run("Error 重试", func(t *testing.T) {
		runningNum := 0

		Running(func() (any, error) {
			runningNum += 1
			return nil, errors.New("112")
		}, IfError(func(err error) bool {
			return err != nil
		}), StopStrategy(StopAfterAttemptStrategy(3)))
		assert.Equal(t, 4, runningNum)
	})

	t.Run("Result 重试", func(t *testing.T) {
		runningNum := 0

		Running(func() (any, error) {
			runningNum += 1
			return "xxxx", errors.New("112")
		}, IfResult(func(res any) bool {
			return "xxxx" == res.(string)
		}), StopStrategy(StopAfterAttemptStrategy(3)),
			WaitStrategy(IncrementingWaitStrategy(100, 20)))
		assert.Equal(t, 4, runningNum)

		Running(func() (any, error) {
			runningNum += 1
			return "bbbb", errors.New("112")
		}, IfResult(func(res any) bool {
			return "xxxx" == res.(string)
		}), StopStrategy(StopAfterAttemptStrategy(3)),
			WaitStrategy(IncrementingWaitStrategy(100, 20)))
		assert.Equal(t, 5, runningNum)
	})

	t.Run("Panic 重试", func(t *testing.T) {
		runningNum := 0

		Running(func() (any, error) {
			runningNum += 1
			panic("xxxx")
		}, IfPanic(func(res any) bool {
			return "xxxx" == res.(string)
		}), StopStrategy(StopAfterAttemptStrategy(3)),
			WaitStrategy(ExponentialWaitStrategy(100, 20)))
		assert.Equal(t, 4, runningNum)

		Running(func() (any, error) {
			runningNum += 1
			panic("12345")
		}, IfPanic(func(res any) bool {
			return "xxxx" == res.(string)
		}), StopStrategy(StopAfterAttemptStrategy(3)),
			WaitStrategy(ExponentialWaitStrategy(100, 2)))
		assert.Equal(t, 5, runningNum)
	})

	t.Run("重试超时-IncrementingWaitStrategy", func(t *testing.T) {
		runningNum := 0

		err := Running(func() (any, error) {
			runningNum += 1
			panic("xxxx")
		}, IfPanic(func(res any) bool {
			return "xxxx" == res.(string)
		}), StopStrategy(StopAfterDelayStrategy(100*time.Millisecond)),
			WaitStrategy(IncrementingWaitStrategy(10*time.Millisecond, 20*time.Millisecond)))
		assert.Equal(t, ErrTimeout, err)
		assert.Equal(t, 6, runningNum)
	})

	t.Run("重试超时-ExponentialWaitStrategy", func(t *testing.T) {
		runningNum := 0

		err := Running(func() (any, error) {
			runningNum += 1
			panic("xxxx")
		}, IfPanic(func(res any) bool {
			return "xxxx" == res.(string)
		}), StopStrategy(StopAfterDelayStrategy(100*time.Millisecond)),
			WaitStrategy(ExponentialWaitStrategy(20*time.Millisecond, 2)))
		assert.Equal(t, ErrTimeout, err)
		assert.Equal(t, 4, runningNum)
	})

	t.Run("重试超时-ExponentialWaitStrategy", func(t *testing.T) {
		runningNum := 0

		err := Running(func() (any, error) {
			runningNum += 1
			panic("xxxx")
		}, IfPanic(func(res any) bool {
			return "xxxx" == res.(string)
		}), StopStrategy(StopAfterDelayStrategy(100*time.Millisecond)),
			WaitStrategy(FibonacciWaitStrategy(10*time.Millisecond)))
		assert.Equal(t, ErrTimeout, err)
		assert.Equal(t, 6, runningNum)
	})
}
