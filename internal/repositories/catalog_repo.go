package repositories

type CatalogMaterialRepo struct {
	*BaseRepo
}

func NewCatalogMaterialRepo() *CatalogMaterialRepo {
	return &CatalogMaterialRepo{BaseRepo: NewBaseRepo("catalog_materials", "id")}
}

type SupplierRepo struct {
	*BaseRepo
}

func NewSupplierRepo() *SupplierRepo {
	return &SupplierRepo{BaseRepo: NewBaseRepo("suppliers", "id")}
}

type CustomerRepo struct {
	*BaseRepo
}

func NewCustomerRepo() *CustomerRepo {
	return &CustomerRepo{BaseRepo: NewBaseRepo("customers", "id")}
}

type TicketTypeRepo struct {
	*BaseRepo
}

func NewTicketTypeRepo() *TicketTypeRepo {
	return &TicketTypeRepo{BaseRepo: NewBaseRepo("ticket_types", "id")}
}

type VehicleCatalogRepo struct {
	*BaseRepo
}

func NewVehicleCatalogRepo() *VehicleCatalogRepo {
	return &VehicleCatalogRepo{BaseRepo: NewBaseRepo("vehicle_catalogs", "id")}
}

type StationRepo struct {
	*BaseRepo
}

func NewStationRepo() *StationRepo {
	return &StationRepo{BaseRepo: NewBaseRepo("stations", "id")}
}

type EquipmentRepo struct {
	*BaseRepo
}

func NewEquipmentRepo() *EquipmentRepo {
	return &EquipmentRepo{BaseRepo: NewBaseRepo("equipment", "id")}
}
