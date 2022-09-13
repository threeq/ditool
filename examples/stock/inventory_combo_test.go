package stock

import (
	"github.com/stretchr/testify/assert"
	"github.com/threeq/ditool/lockx"
	"testing"
)

func TestComboInventory_Remaining(t *testing.T) {
	lf := lockx.NewLocalLockerFactory()
	comboChannelID := "123"
	n1 := NewNormalInventory("inventory_c_normal_1", 1000, lf)
	n1.NewChannel(&comboChannelID, NewStoreStrategy(Share()), nil)

	n2 := NewNormalInventory("inventory_c_normal_2",400, lf)
	n2.NewChannel(&comboChannelID, NewStoreStrategy(Share()), nil)

	c1 := NewComboInventory("inventory_combo_1", comboChannelID, ComboWeight{
		10: n1,
		2:  n2,
	}, lf)

	channelID := "345"
	c1.NewChannel(&channelID, NewStoreStrategy(Share()), nil)

	assert.Equal(t, 100, c1.Status(nil).SelfRemaining)
}

func TestComboInventory_Dec_Inc(t *testing.T) {
	lf := lockx.NewLocalLockerFactory()
	invent1 := NewNormalInventory("inventory_c_normal_3",1000, lf)
	invent2 := NewNormalInventory("inventory_c_normal_4", 1000, lf)
	invent3 := NewNormalInventory("inventory_c_normal_5", 1000, lf)

	comboChannelID := "000"
	invent1.NewChannel(&comboChannelID, NewStoreStrategy(Share()), nil)
	invent2.NewChannel(&comboChannelID, NewStoreStrategy(Share()), nil)
	invent3.NewChannel(&comboChannelID, NewStoreStrategy(Share()), nil)
	combo1 := NewComboInventory("inventory_combo_2", comboChannelID, ComboWeight{
		1: invent1,
		2: invent2,
		3: invent3,
	}, lf)

	assertHelper(t, &SaleChannelStatus{
		Max:           333,
		SelfRemaining: 333,
		SelfSales:     0,
		SelfHoldCfg:   0,
		SelfHoldRem:   0,
		ChildSales:    0,
		ChildHoldRem:  0,
		ChildHoldCfg:  0,
	}, combo1.Status(nil))

	ch1 := "123"
	combo1.NewChannel(&ch1, NewStoreStrategy(Share(), Hold(30)), nil)

	assertHelper(t, &SaleChannelStatus{
		Max:           333,
		SelfRemaining: 303,
		SelfSales:     0,
		SelfHoldCfg:   0,
		SelfHoldRem:   0,
		ChildSales:    0,
		ChildHoldRem:  30,
		ChildHoldCfg:  30,
	}, combo1.Status(nil))
	assertHelper(t, &SaleChannelStatus{
		Max:           333,
		SelfRemaining: 333,
		SelfSales:     0,
		SelfHoldCfg:   30,
		SelfHoldRem:   30,
		ChildSales:    0,
		ChildHoldRem:  0,
		ChildHoldCfg:  0,
	}, combo1.Status(&ch1))

	t.Run("错误-扣减和回退", func(t *testing.T) {
		_, err := combo1.SafeDec(nil, 400)
		assert.Equal(t, "剩余库存不足", err.Error())
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 303,
			SelfSales:     0,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  30,
			ChildHoldCfg:  30,
		}, combo1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 333,
			SelfSales:     0,
			SelfHoldCfg:   30,
			SelfHoldRem:   30,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, combo1.Status(&ch1))
		_, err = combo1.SafeInc(nil, 400)
		assert.Equal(t, "销售数量错误", err.Error())
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 303,
			SelfSales:     0,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  30,
			ChildHoldCfg:  30,
		}, combo1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 333,
			SelfSales:     0,
			SelfHoldCfg:   30,
			SelfHoldRem:   30,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, combo1.Status(&ch1))

	})

	t.Run("扣减-无占用渠道 20", func(t *testing.T) {
		combo1.SafeDec(nil, 20)
		// 验证组合本身
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 283,
			SelfSales:     20,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  30,
			ChildHoldCfg:  30,
		}, combo1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 313,
			SelfSales:     0,
			SelfHoldCfg:   30,
			SelfHoldRem:   30,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, combo1.Status(&ch1))
		//	验证组合中的单品
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 980,
			SelfSales:     20,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&comboChannelID))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 960,
			SelfSales:     40,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent2.Status(&comboChannelID))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 940,
			SelfSales:     60,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent3.Status(&comboChannelID))
	})

	t.Run("扣减-有占用渠道 20", func(t *testing.T) {
		combo1.SafeDec(&ch1, 20)
		// 验证组合本身
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 283,
			SelfSales:     20,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    20,
			ChildHoldRem:  10,
			ChildHoldCfg:  30,
		}, combo1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 293,
			SelfSales:     20,
			SelfHoldCfg:   30,
			SelfHoldRem:   10,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, combo1.Status(&ch1))
		//	验证组合中的单品
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 960,
			SelfSales:     40,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&comboChannelID))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 920,
			SelfSales:     80,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent2.Status(&comboChannelID))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 880,
			SelfSales:     120,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent3.Status(&comboChannelID))
	})

	t.Run("回退-无占用渠道 10", func(t *testing.T) {
		combo1.SafeInc(nil, 10)
		// 验证组合本身
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 293,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    20,
			ChildHoldRem:  10,
			ChildHoldCfg:  30,
		}, combo1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 303,
			SelfSales:     20,
			SelfHoldCfg:   30,
			SelfHoldRem:   10,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, combo1.Status(&ch1))
		//	验证组合中的单品
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 970,
			SelfSales:     30,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&comboChannelID))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 940,
			SelfSales:     60,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent2.Status(&comboChannelID))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 910,
			SelfSales:     90,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent3.Status(&comboChannelID))
	})

	t.Run("回退-有占用渠道 10", func(t *testing.T) {
		combo1.SafeInc(&ch1, 10)
		// 验证组合本身
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 293,
			SelfSales:     10,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    10,
			ChildHoldRem:  20,
			ChildHoldCfg:  30,
		}, combo1.Status(nil))
		assertHelper(t, &SaleChannelStatus{
			Max:           333,
			SelfRemaining: 313,
			SelfSales:     10,
			SelfHoldCfg:   30,
			SelfHoldRem:   20,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, combo1.Status(&ch1))
		//	验证组合中的单品
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 980,
			SelfSales:     20,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent1.Status(&comboChannelID))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 960,
			SelfSales:     40,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent2.Status(&comboChannelID))
		assertHelper(t, &SaleChannelStatus{
			Max:           1000,
			SelfRemaining: 940,
			SelfSales:     60,
			SelfHoldCfg:   0,
			SelfHoldRem:   0,
			ChildSales:    0,
			ChildHoldRem:  0,
			ChildHoldCfg:  0,
		}, invent3.Status(&comboChannelID))
	})
}
