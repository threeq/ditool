package lockx

import (
	"context"
	"github.com/threeq/ditool/retry"
	"time"
)

// Locker 锁接口
type Locker interface {
	Lock() error
	Unlock() error
}

// RWLocker 读写锁接口
type RWLocker interface {
	Locker
	RLock() error
	RUnlock() error
	RLocker() Locker
}

// LockerFactory 锁创建工厂实现接口
type LockerFactory interface {
	Mutex(context.Context, ...Option) (Locker, error)
	MutexL2(ctx context.Context, options ...Option) (Locker, error)
	RWMutex(context.Context, ...Option) (RWLocker, error)
}

// --------------------------------------------------------------------

// LockerMeta 锁配置元数据
type LockerMeta struct {
	key          string
	ttl          time.Duration
	retryFactory func() *retry.Executor
}

func (m *LockerMeta) retry() *retry.Executor {
	return retry.NewExecutor(retry.IfError(func(err error) bool {
		return err != nil
	}), retry.IfPanic(func(ex any) bool {
		return ex != nil
	}), retry.StopStrategy(retry.StopAfterDelayStrategy(2*m.ttl)),
		retry.WaitStrategy(retry.FibonacciWaitStrategy(1*time.Millisecond)))
}

// Option 锁配置元数据设置
type Option func(*LockerMeta)

// Key 锁 id
func Key(id string) Option {
	return func(meta *LockerMeta) {
		meta.key = id
	}
}

// TTL 锁过期时间
func TTL(ttl time.Duration) Option {
	return func(meta *LockerMeta) {
		meta.ttl = ttl
	}
}

var _factory LockerFactory

func init() {
	_factory = NewLocalLockerFactory()
}

func Init(factory LockerFactory) {
	_factory = factory
}

func Mutex(ctx context.Context, opts ...Option) (Locker, error) {
	return _factory.Mutex(ctx, opts...)
}
func MutexL2(ctx context.Context, opts ...Option) (Locker, error) {
	return _factory.MutexL2(ctx, opts...)
}
func RWMutex(ctx context.Context, opts ...Option) (RWLocker, error) {
	return _factory.RWMutex(ctx, opts...)
}
