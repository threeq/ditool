package stock

import (
	"github.com/threeq/ditool/mathx"
	"math"
)

type StoreStrategy struct {
	channel  *SaleChannel
	minSales int
	maxSales int
}

func (strat *StoreStrategy) ResetLimit(total int) {
	allSales := strat.Sales()
	strat.maxSales = total + allSales
}

func (strat *StoreStrategy) ResetHold(n int) {
	strat.minSales = n
}

func (strat *StoreStrategy) Bind(channel *SaleChannel) {
	strat.channel = channel
}

func (strat *StoreStrategy) Limit() int {
	return strat.maxSales
}

func (strat *StoreStrategy) Hold() (int, int) {
	holdCfg := strat.minSales
	holdRem := mathx.Max(holdCfg-strat.channel.sales, 0)

	for _, saleChannel := range strat.channel.Child() {
		cHoldCfg, cHoldRemaining := saleChannel.strategy.Hold()
		holdCfg += cHoldCfg
		holdRem += cHoldRemaining
	}
	return holdCfg, holdRem
}

func (strat *StoreStrategy) Sales() int {
	ch := strat.channel.sales
	for _, saleChannel := range strat.channel.Child() {
		ch += saleChannel.strategy.Sales()
	}
	return ch
}

func (strat *StoreStrategy) IsShared() bool {
	return strat.maxSales == Shared
}

func (strat *StoreStrategy) CompleteStatus() *SaleChannelStatus {
	selfMaxSales := 0
	selfRemaining := 0
	holdCfg, holdRem := strat.Hold()
	allSales := strat.Sales()
	childSales := allSales - strat.channel.sales
	selfHoldRem := mathx.Max(0, strat.minSales-strat.channel.sales)

	switch strat.maxSales {
	case NotSell:
		return &SaleChannelStatus{
			Max:           0,
			SelfSales:     strat.channel.sales,
			SelfHoldRem:   selfHoldRem,
			SelfHoldCfg:   strat.minSales,
			SelfRemaining: 0,
			ChildHoldRem:  holdRem,
			ChildHoldCfg:  holdCfg - strat.minSales,
			ChildSales:    childSales,
		}
	case Shared:
		ps := strat.channel.Parent().Status()
		selfMaxSales = ps.Max - holdCfg + strat.minSales
		//selfHoldRem * 2 的原因是
		// 1、计算 parent.SelfRemaining 时，已经将自己的 selfHoldRem 扣减一遍
		// 2、计算 holdRem 是，也将自己的s selfHoldRem 扣减一遍
		selfRemaining = ps.SelfRemaining - holdRem + selfHoldRem*2
	default:
		selfMaxSales = strat.maxSales
		selfRemaining = selfMaxSales - allSales - holdRem + selfHoldRem*2
	}

	return &SaleChannelStatus{
		Max:           selfMaxSales,
		SelfSales:     strat.channel.sales,
		SelfHoldRem:   selfHoldRem,
		SelfHoldCfg:   strat.minSales,
		SelfRemaining: selfRemaining,
		ChildHoldRem:  holdRem - selfHoldRem,
		ChildHoldCfg:  holdCfg - strat.minSales,
		ChildSales:    childSales,
	}
}

func NewStoreStrategy(opts ...Option) *StoreStrategy {
	s := &StoreStrategy{minSales: 0, maxSales: Shared}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ---------------------------------------------------------------------------------

type Option func(*StoreStrategy)

func Share() Option {
	return func(strategy *StoreStrategy) {
		strategy.maxSales = Shared
	}
}

func Limit(n int) Option {
	return func(strategy *StoreStrategy) {
		strategy.maxSales = n
	}
}

func Hold(n int) Option {
	return func(strategy *StoreStrategy) {
		strategy.minSales = n
	}
}

// ------------------------------------------------------------

type CombSaleStrategy struct {
	inventory *ComboInventory
	channel   *SaleChannel
}

func (combo *CombSaleStrategy) ResetLimit(total int) {

}

func (combo *CombSaleStrategy) Limit() int {
	return -1
}

func (combo *CombSaleStrategy) Hold() (int, int) {
	holdCfg := 0
	holdRem := 0
	for _, saleChannel := range combo.channel.Child() {
		cHoldCfg, cHoldRem := saleChannel.strategy.Hold()
		holdCfg += cHoldCfg
		holdRem += cHoldRem
	}
	return holdCfg, holdRem
}

func (combo *CombSaleStrategy) Sales() int {
	ch := combo.channel.sales
	for _, saleChannel := range combo.channel.Child() {
		ch += saleChannel.strategy.Sales()
	}
	return ch
}

func (combo *CombSaleStrategy) CompleteStatus() *SaleChannelStatus {
	maxSales, remaining := combo.comboInventoryStatus()

	holdCfg, holdRem := combo.Hold()
	allSales := combo.Sales()

	return &SaleChannelStatus{
		Max:           maxSales,
		SelfSales:     combo.channel.sales,
		SelfHoldCfg:   0,
		SelfRemaining: remaining - holdRem,
		SelfHoldRem:   0,
		ChildHoldRem:  holdRem,
		ChildHoldCfg:  holdCfg,
		ChildSales:    allSales - combo.channel.sales,
	}

}

func (combo *CombSaleStrategy) comboInventoryStatus() (int, int) {

	max := math.MaxInt
	remaining := math.MaxInt

	w := combo.inventory.inventoriesWeight
	for i := range combo.inventory.inventories {
		inventoryChanID := combo.inventory.comboChannelID
		istatus := combo.inventory.inventories[i].Channel(&inventoryChanID).Status()

		tmp := istatus.SelfRemaining / w[i]
		if tmp < remaining {
			remaining = tmp
		}

		tmp = istatus.Max / w[i]
		if tmp < max {
			max = tmp
		}
	}
	return max, remaining
}

func (combo *CombSaleStrategy) Bind(channel *SaleChannel) {
	combo.channel = channel
}

func (combo *CombSaleStrategy) IsShared() bool {
	return false
}

func (combo *CombSaleStrategy) ResetHold(n int) {

}
