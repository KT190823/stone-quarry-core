package repositories

type InventoryInboundRepo struct {
	*BaseRepo
}

func NewInventoryInboundRepo() *InventoryInboundRepo {
	return &InventoryInboundRepo{BaseRepo: NewBaseRepo("inventory_inbound", "id")}
}

type InventoryOutboundRepo struct {
	*BaseRepo
}

func NewInventoryOutboundRepo() *InventoryOutboundRepo {
	return &InventoryOutboundRepo{BaseRepo: NewBaseRepo("inventory_outbound", "id")}
}

type InventoryStocktakeRepo struct {
	*BaseRepo
}

func NewInventoryStocktakeRepo() *InventoryStocktakeRepo {
	return &InventoryStocktakeRepo{BaseRepo: NewBaseRepo("inventory_stocktake", "id")}
}

type InventoryMovementRepo struct {
	*BaseRepo
}

func NewInventoryMovementRepo() *InventoryMovementRepo {
	return &InventoryMovementRepo{BaseRepo: NewBaseRepo("inventory_movements", "id")}
}
