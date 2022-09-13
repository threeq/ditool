package stock

type FreezeRequest struct {
	id string
	freezeId string
	number int
}


type FreezeFacade struct {
	factory *InventoryFactory
}

func (f *FreezeFacade) Freeze(requests []FreezeRequest) (int, error) {
	panic("implement me")
}

func (f *FreezeFacade) UnFreeze(requests []FreezeRequest) (int, error) {
	panic("implement me")
}

type SkuFacade struct {

}

type InfoFacade interface {
	Log()
	Graph()
}
