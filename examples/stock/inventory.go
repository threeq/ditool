package stock

import (
	"context"
	"errors"
	"github.com/threeq/ditool/lockx"
)

const (
	Shared  = -1
	NotSell = 0
)

var (
	ErrRemainingLacking = errors.New("剩余库存不足")
	ErrSalesNumber      = errors.New("销售数量错误")
)

// Inventory 库存
type Inventory interface {
	ID() *string
	Lock() error
	UnLock() error
	//Sync 同步总库存
	//@return 返回总库存
	Sync(total int) (int, error)

	//Distribution 库存分布
	//@return 返回库存分布
	Distribution() Inventory
	//Status 渠道剩余可售卖数量
	//@Param freeId 冻结库存剩余数量
	//              nil stockId 全部剩余数量
	//				""  非冻结部分库存剩余数量
	Status(chanID *string) *SaleChannelStatus
	//SafeDec 扣减库存
	SafeDec(chanID *string, number int) (Inventory, error)
	//SafeInc 增加库存
	SafeInc(chanID *string, number int) (Inventory, error)
	Dec(s *string, i int) (Inventory, error)
	Inc(s *string, i int) (Inventory, error)
	//Channel 返回销售渠道
	Channel(*string) *SaleChannel
	//NewChannel 新建渠道
	NewChannel(name *string, strategy SaleStrategy, from *SaleChannel) *SaleChannel
}

// SaleStrategy 销售策略
type SaleStrategy interface {
	IsShared() bool
	Limit() int
	Hold() (int, int)
	Sales() int
	CompleteStatus() *SaleChannelStatus
	Bind(*SaleChannel)
	ResetLimit(total int)
	ResetHold(n int)
}

// SaleChannel 销售渠道
type SaleChannel struct {
	strategy SaleStrategy
	p        *SaleChannel
	child    []*SaleChannel
	sales    int
}

func NewSaleChannel(strategy SaleStrategy, p *SaleChannel) *SaleChannel {
	c := &SaleChannel{strategy: strategy, p: p}
	strategy.Bind(c)
	if p != nil {
		p.child = append(p.child, c)
	}
	return c
}

type SaleChannelStatus struct {
	Max           int
	SelfRemaining int
	SelfSales     int
	SelfHoldCfg   int
	SelfHoldRem   int
	ChildSales    int
	ChildHoldCfg  int
	ChildHoldRem  int
}

func (s *SaleChannel) Strategy() SaleStrategy {
	return s.strategy
}

func (s *SaleChannel) Status() *SaleChannelStatus {
	return s.strategy.CompleteStatus()
}

func (s *SaleChannel) Parent() *SaleChannel {
	return s.p
}

func (s *SaleChannel) Child() []*SaleChannel {
	return s.child
}

type InventoryMeta struct {
	Id     *string
	locker lockx.Locker
}

func NewInventoryMeta(ID *string, lf lockx.LockerFactory) *InventoryMeta {
	locker, err := lf.Mutex(context.Background(), lockx.Key(*ID))
	if err != nil {
		panic(err)
	}
	return &InventoryMeta{Id: ID, locker: locker}

}

func (im *InventoryMeta) Lock() error {
	return im.locker.Lock()
}

func (im *InventoryMeta) UnLock() error {
	return im.locker.Unlock()
}

func (im *InventoryMeta) ID() *string {
	return im.Id
}
