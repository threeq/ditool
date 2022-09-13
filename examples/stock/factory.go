package stock

type InventoryFactory struct {
	repos Repository
}

func (receiver *InventoryFactory) Builder(skuID string) Inventory {
	panic("implement me")
}
