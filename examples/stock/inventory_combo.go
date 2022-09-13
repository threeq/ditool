package stock

import (
	"context"
	"github.com/threeq/ditool/lockx"
	"sort"
)

type ComboInventory struct {
	*InventoryMeta
	channels          map[string]*SaleChannel
	comboChannelID    string
	inventories       []Inventory
	inventoriesLock   []Inventory
	inventoriesWeight []int
	root              *SaleChannel
}

type ComboWeight map[int]Inventory

func NewComboInventory(id string, comboChannelID string, comboWeight ComboWeight, lf lockx.LockerFactory) *ComboInventory {
	locker, _ := lf.MutexL2(context.Background(), lockx.Key(id))
	ci := &ComboInventory{
		InventoryMeta: &InventoryMeta{
			Id:     &id,
			locker: locker,
		},
		comboChannelID: comboChannelID,
		channels:       map[string]*SaleChannel{},
	}
	cs := &CombSaleStrategy{inventory: ci}
	sc := NewSaleChannel(cs, nil)
	ci.root = sc

	for w, inventory := range comboWeight {
		ci.inventories = append(ci.inventories, inventory)
		ci.inventoriesWeight = append(ci.inventoriesWeight, w)
	}
	return ci
}

func (ni *ComboInventory) Sync(total int) (int, error) {
	return 0, nil
}

func (ni *ComboInventory) Distribution() Inventory {
	panic("implement me")
}

func (ni *ComboInventory) Status(chanID *string) *SaleChannelStatus {
	if chanID == nil {
		return ni.root.Status()
	}
	return ni.channels[*chanID].Status()
}

func (ni *ComboInventory) SafeDec(chanID *string, number int) (Inventory, error) {
	ni.Lock()
	defer ni.UnLock()
	return ni.Dec(chanID, number)
}

func (ni *ComboInventory) Dec(chanID *string, number int) (Inventory, error) {
	channel := ni.Channel(chanID)
	status := channel.Status()
	if status.SelfRemaining >= number {
		channel.sales += number

		for i, inventory := range ni.inventories {
			_, _ = inventory.Dec(&ni.comboChannelID, number*ni.inventoriesWeight[i])
		}
	} else {
		return nil, ErrRemainingLacking
	}
	return nil, nil
}

func (ni *ComboInventory) SafeInc(chanID *string, number int) (Inventory, error) {
	ni.Lock()
	defer ni.UnLock()
	return ni.Inc(chanID, number)
}

func (ni *ComboInventory) Inc(chanID *string, number int) (Inventory, error) {
	channel := ni.Channel(chanID)
	if channel.sales < number {
		return nil, ErrSalesNumber
	}
	channel.sales -= number
	for i, inventory := range ni.inventories {
		_, _ = inventory.Inc(&ni.comboChannelID, number*ni.inventoriesWeight[i])
	}
	return nil, nil
}

func (ni *ComboInventory) Channel(name *string) *SaleChannel {
	if name == nil {
		return ni.root
	}

	return ni.channels[*name]
}

func (ni *ComboInventory) NewChannel(name *string, strategy SaleStrategy, from *SaleChannel) *SaleChannel {
	if from == nil {
		from = ni.root
	}
	ch := NewSaleChannel(strategy, from)
	ni.channels[*name] = ch
	return ch
}

func (ni *ComboInventory) Lock() error {
	ni.InventoryMeta.Lock()
	ni.inventoriesLock = make([]Inventory, len(ni.inventories))
	copy(ni.inventoriesLock, ni.inventories)
	sort.Slice(ni.inventoriesLock, func(i, j int) bool {
		return *ni.inventoriesLock[i].ID() > *ni.inventoriesLock[j].ID()
	})
	for i := 0; i < len(ni.inventoriesLock); i++ {
		// TODO 失败处理，存在失败主动释放一次之前加过的锁
		ni.inventoriesLock[i].Lock()
	}
	return nil
}

func (ni *ComboInventory) UnLock() error {
	for i := len(ni.inventoriesLock) - 1; i >= 0; i-- {
		ni.inventoriesLock[i].UnLock()
	}
	ni.InventoryMeta.UnLock()
	return nil
}
