package stock

import (
	"github.com/stretchr/testify/assert"
	"github.com/threeq/ditool/lockx"
	"testing"
)

func assertHelper(t *testing.T, expected, actual *SaleChannelStatus) {
	assert.Equal(t, expected.Max, actual.Max, "总数")
	assert.Equal(t, expected.SelfRemaining, actual.SelfRemaining, "渠道剩余数")
	assert.Equal(t, expected.SelfSales, actual.SelfSales, "渠道售卖数")
	assert.Equal(t, expected.SelfHoldCfg, actual.SelfHoldCfg, "渠道占用配置数")
	assert.Equal(t, expected.SelfHoldRem, actual.SelfHoldRem, "渠道占用剩余数")
	assert.Equal(t, expected.ChildSales, actual.ChildSales, "子渠道售卖总数")
	assert.Equal(t, expected.ChildHoldRem, actual.ChildHoldRem, "子渠道占用剩余总数")
	assert.Equal(t, expected.ChildHoldCfg, actual.ChildHoldCfg, "子渠道占用配置数")
}

func TestNormalInventory_NewChannel(t *testing.T) {
	lf := lockx.NewLocalLockerFactory()
	ch1 := "123"
	ch2 := "456"
	ch3 := "789"
	invent1 := NewNormalInventory("TestNormalInventory_NewChannel", 1000, lf)
	invent1.NewChannel(&ch1, NewStoreStrategy(Share(), Hold(10)), nil)
	invent1.NewChannel(&ch2, NewStoreStrategy(Share(), Hold(20)), nil)
	invent1.NewChannel(&ch3, NewStoreStrategy(Share()), nil)

	status := invent1.Status(nil)
	assert.Equal(t, 970, status.SelfRemaining)
	status = invent1.Status(&ch1)
	assert.Equal(t, 980, status.SelfRemaining)
	status = invent1.Status(&ch2)
	assert.Equal(t, 990, status.SelfRemaining)
	status = invent1.Status(&ch3)
	assert.Equal(t, 970, status.SelfRemaining)
}

func TestNormalInventory_Dec_Inc(t *testing.T) {
	lf := lockx.NewLocalLockerFactory()
	ch1 := "123"
	ch2 := "456"
	ch3 := "789"
	invent1 := NewNormalInventory("TestNormalInventory_Dec_Inc", 1000, lf)

	initDec(t, invent1, ch1, ch2, ch3)

	t.Run("回退-根渠道: 2", func(t *testing.T) {
		invent1.SafeInc(nil, 2)
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 947,
			SelfSales:     8,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    40,
			ChildHoldCfg:  30,
			ChildHoldRem:  5,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 952,
			SelfSales:     5,
			SelfHoldCfg:   10,
			SelfHoldRem:   5,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 947,
			SelfSales:     25,
			SelfHoldCfg:   20,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 947,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})

	t.Run("回退-无最小占用渠道：ch3 : 3", func(t *testing.T) {
		invent1.SafeInc(&ch3, 3)
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 950,
			SelfSales:     8,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    37,
			ChildHoldCfg:  30,
			ChildHoldRem:  5,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 955,
			SelfSales:     5,
			SelfHoldCfg:   10,
			SelfHoldRem:   5,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 950,
			SelfSales:     25,
			SelfHoldCfg:   20,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 950,
			SelfSales:     7,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})

	t.Run("回退-有最小占用渠道：ch1 最小占用有: 3", func(t *testing.T) {
		invent1.SafeInc(&ch1, 3)
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 950,
			SelfSales:     8,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    34,
			ChildHoldCfg:  30,
			ChildHoldRem:  8,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 958,
			SelfSales:     2,
			SelfHoldCfg:   10,
			SelfHoldRem:   8,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 950,
			SelfSales:     25,
			SelfHoldCfg:   20,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 950,
			SelfSales:     7,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})

	t.Run("回退-有最小占用渠道：ch2 最小占用之上: 3", func(t *testing.T) {
		invent1.SafeInc(&ch2, 3)
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 953,
			SelfSales:     8,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    31,
			ChildHoldCfg:  30,
			ChildHoldRem:  8,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 961,
			SelfSales:     2,
			SelfHoldCfg:   10,
			SelfHoldRem:   8,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 953,
			SelfSales:     22,
			SelfHoldCfg:   20,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 953,
			SelfSales:     7,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})

	t.Run("回退-有最小占用渠道：ch2 最小占用以下: 10", func(t *testing.T) {
		invent1.SafeInc(&ch2, 10)
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 955,
			SelfSales:     8,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    21,
			ChildHoldCfg:  30,
			ChildHoldRem:  16,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 963,
			SelfSales:     2,
			SelfHoldCfg:   10,
			SelfHoldRem:   8,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 963,
			SelfSales:     12,
			SelfHoldCfg:   20,
			SelfHoldRem:   8,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 955,
			SelfSales:     7,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})

	t.Run("错误-回归过量：ch2： 50", func(t *testing.T) {
		_, err := invent1.SafeInc(&ch2, 50)
		assert.Equal(t, "销售数量错误", err.Error())
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 955,
			SelfSales:     8,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    21,
			ChildHoldCfg:  30,
			ChildHoldRem:  16,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 963,
			SelfSales:     2,
			SelfHoldCfg:   10,
			SelfHoldRem:   8,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 963,
			SelfSales:     12,
			SelfHoldCfg:   20,
			SelfHoldRem:   8,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 955,
			SelfSales:     7,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})

	t.Run("错误-购买过量：ch3： 1000", func(t *testing.T) {
		_, err := invent1.SafeDec(&ch3, 1000)
		assert.Equal(t, "剩余库存不足", err.Error())
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 955,
			SelfSales:     8,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    21,
			ChildHoldCfg:  30,
			ChildHoldRem:  16,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 963,
			SelfSales:     2,
			SelfHoldCfg:   10,
			SelfHoldRem:   8,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 963,
			SelfSales:     12,
			SelfHoldCfg:   20,
			SelfHoldRem:   8,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 955,
			SelfSales:     7,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})
}

func initDec(t *testing.T, invent1 *NormalInventory, ch1 string, ch2 string, ch3 string) {
	t.Run("初始化库存渠道", func(t *testing.T) {
		invent1.NewChannel(&ch1, NewStoreStrategy(Share(), Hold(10)), nil)
		invent1.NewChannel(&ch2, NewStoreStrategy(Share(), Hold(20)), nil)
		invent1.NewChannel(&ch3, NewStoreStrategy(Share()), nil)

		status := invent1.Status(nil)
		assert.Equal(t, 970, status.SelfRemaining)
		status = invent1.Status(&ch1)
		assert.Equal(t, 980, status.SelfRemaining)
		status = invent1.Status(&ch2)
		assert.Equal(t, 990, status.SelfRemaining)
		status = invent1.Status(&ch3)
		assert.Equal(t, 970, status.SelfRemaining)
	})

	t.Run("扣减-跟渠道: 10", func(t *testing.T) {
		// 跟渠道售卖
		invent1.SafeDec(nil, 10)
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 960,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  30,
			ChildHoldCfg:  30,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 970,
			SelfSales:     0,
			SelfHoldCfg:   10,
			SelfHoldRem:   10,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 980,
			SelfSales:     0,
			SelfHoldCfg:   20,
			SelfHoldRem:   20,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 960,
			SelfSales:     0,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})

	t.Run("扣减-有最小占用渠道：ch1 售卖数小于占用数: 5", func(t *testing.T) {
		invent1.SafeDec(&ch1, 5)
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 960,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    5,
			ChildHoldCfg:  30,
			ChildHoldRem:  25,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 965,
			SelfSales:     5,
			SelfHoldCfg:   10,
			SelfHoldRem:   5,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 980,
			SelfSales:     0,
			SelfHoldCfg:   20,
			SelfHoldRem:   20,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 960,
			SelfSales:     0,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})

	t.Run("扣减-无最小占用渠道：ch3 售卖: 10", func(t *testing.T) {
		invent1.SafeDec(&ch3, 10)
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 950,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    15,
			ChildHoldCfg:  30,
			ChildHoldRem:  25,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 955,
			SelfSales:     5,
			SelfHoldCfg:   10,
			SelfHoldRem:   5,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 970,
			SelfSales:     0,
			SelfHoldCfg:   20,
			SelfHoldRem:   20,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 950,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})

	t.Run("扣减-有最小占用渠道：ch2 售卖数大于占用数: 25", func(t *testing.T) {
		invent1.SafeDec(&ch2, 25)
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 945,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    40,
			ChildHoldCfg:  30,
			ChildHoldRem:  5,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 950,
			SelfSales:     5,
			SelfHoldCfg:   10,
			SelfHoldRem:   5,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 945,
			SelfSales:     25,
			SelfHoldCfg:   20,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 945,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})
}

func TestNormalInventory_Sync(t *testing.T) {
	lf := lockx.NewLocalLockerFactory()
	ch1 := "123"
	ch2 := "456"
	ch3 := "789"
	invent1 := NewNormalInventory("TestNormalInventory_Dec_Inc", 1000, lf)

	initDec(t, invent1, ch1, ch2, ch3)

	t.Run("同步-同步数量大于占用 500", func(t *testing.T) {
		invent1.Sync(500)
		assertHelper(t, &SaleChannelStatus{
			Max:           550,
			SelfRemaining: 495,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    40,
			ChildHoldCfg:  30,
			ChildHoldRem:  5,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           550,
			SelfRemaining: 500,
			SelfSales:     5,
			SelfHoldCfg:   10,
			SelfHoldRem:   5,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           550,
			SelfRemaining: 495,
			SelfSales:     25,
			SelfHoldCfg:   20,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           550,
			SelfRemaining: 495,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})

	t.Run("同步-同步数量小于占用 5", func(t *testing.T) {
		invent1.Sync(3)
		assertHelper(t, &SaleChannelStatus{
			Max:           53,
			SelfRemaining: 3,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    40,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           53,
			SelfRemaining: 3,
			SelfSales:     5,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch1))
		assertHelper(t, &SaleChannelStatus{
			Max:           53,
			SelfRemaining: 3,
			SelfSales:     25,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch2))
		assertHelper(t, &SaleChannelStatus{
			Max:           53,
			SelfRemaining: 3,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&ch3))
	})
}
