package stock

import (
	"context"
	"github.com/threeq/ditool/lockx"
)

type NormalInventory struct {
	*InventoryMeta
	channels map[string]*SaleChannel
	root     *SaleChannel
}

func (ni *NormalInventory) NewChannel(name *string, strategy SaleStrategy, from *SaleChannel) *SaleChannel {
	if from == nil {
		from = ni.root
	}

	ch := NewSaleChannel(strategy, from)
	ni.channels[*name] = ch
	return ch
}

func (ni *NormalInventory) Channel(name *string) *SaleChannel {
	if name == nil {
		return ni.root
	}

	return ni.channels[*name]
}

func (ni *NormalInventory) Sync(total int) (int, error) {
	ni.Lock()
	defer ni.UnLock()
	return ni.doSync(total)
}

func (ni *NormalInventory) doSync(total int) (int, error) {
	status := ni.root.Status()
	ni.root.strategy.ResetLimit(total)
	if total < status.ChildHoldRem {
		for _, channel := range ni.channels {
			channel.Strategy().ResetHold(0)
		}
	}
	return total, nil
}

func (ni *NormalInventory) Distribution() Inventory {
	panic("implement me")
}

func (ni *NormalInventory) Status(chanID *string) *SaleChannelStatus {
	if chanID == nil {
		return ni.root.Status()
	}
	return ni.channels[*chanID].Status()
}

func (ni *NormalInventory) SafeDec(chanID *string, number int) (Inventory, error) {
	ni.Lock()
	defer ni.UnLock()
	return ni.Dec(chanID, number)
}

func (ni *NormalInventory) Dec(chanID *string, number int) (Inventory, error) {
	channel := ni.Channel(chanID)
	status := channel.Status()
	if status.SelfRemaining >= number {
		channel.sales += number
	} else {
		return nil, ErrRemainingLacking
	}
	return nil, nil
}

func (ni *NormalInventory) SafeInc(chanID *string, number int) (Inventory, error) {
	ni.Lock()
	defer ni.UnLock()
	return ni.Inc(chanID, number)
}

func (ni *NormalInventory) Inc(chanID *string, number int) (Inventory, error) {
	channel := ni.Channel(chanID)
	if channel.sales < number {
		return nil, ErrSalesNumber
	}
	channel.sales -= number
	return nil, nil
}

func NewNormalInventory(id string, total int, lf lockx.LockerFactory) *NormalInventory {
	s := NewStoreStrategy(Limit(total))
	sc := NewSaleChannel(s, nil)

	locker, _ := lf.MutexL2(context.Background(), lockx.Key(id))
	ni := &NormalInventory{
		InventoryMeta: &InventoryMeta{
			Id:     &id,
			locker: locker,
		},
		channels: map[string]*SaleChannel{},
		root:     sc,
	}
	return ni
}
