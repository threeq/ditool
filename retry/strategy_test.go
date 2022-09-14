package retry

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIncrementingWaitStrategy(t *testing.T) {
	ws := IncrementingWaitStrategy(100, 3).(*incWaitStrategy)
	assert.Equal(t, time.Duration(97), ws.curr)
	assert.Equal(t, time.Duration(3), ws.step)
	ws.Wait()
	assert.Equal(t, time.Duration(100), ws.curr)
	ws.Wait()
	assert.Equal(t, time.Duration(103), ws.curr)
	assert.Equal(t, time.Duration(3), ws.step)
	ws.Wait()
	assert.Equal(t, time.Duration(106), ws.curr)
}

func TestExponentialWaitStrategy(t *testing.T) {
	ws := ExponentialWaitStrategy(100, 2).(*expWaitStrategy)
	assert.Equal(t, time.Duration(50), ws.curr)
	assert.Equal(t, int64(2), ws.rate)
	ws.Wait()
	assert.Equal(t, time.Duration(100), ws.curr)
	ws.Wait()
	assert.Equal(t, time.Duration(200), ws.curr)
	ws.Wait()
	assert.Equal(t, time.Duration(400), ws.curr)
}

func TestFibonacciWaitStrategy(t *testing.T) {
	ws := FibonacciWaitStrategy(100).(*fibWaitStrategy)
	assert.Equal(t, time.Duration(100), ws.n1)
	assert.Equal(t, time.Duration(0), ws.n2)

	ws.Wait()
	assert.Equal(t, time.Duration(100), ws.n1)
	assert.Equal(t, time.Duration(100), ws.n2)

	ws.Wait()
	assert.Equal(t, time.Duration(200), ws.n1)
	assert.Equal(t, time.Duration(100), ws.n2)

	ws.Wait()
	assert.Equal(t, time.Duration(300), ws.n1)
	assert.Equal(t, time.Duration(200), ws.n2)

	ws.Wait()
	assert.Equal(t, time.Duration(500), ws.n1)
	assert.Equal(t, time.Duration(300), ws.n2)

	ws.Wait()
	assert.Equal(t, time.Duration(800), ws.n1)
	assert.Equal(t, time.Duration(500), ws.n2)

	ws.Wait()
	assert.Equal(t, time.Duration(1300), ws.n1)
	assert.Equal(t, time.Duration(800), ws.n2)
}
