package services

import (
	"mo-da-backend/internal/repositories"
)

type CatalogMaterialService struct {
	*BaseService
}

func NewCatalogMaterialService() *CatalogMaterialService {
	repo := repositories.NewCatalogMaterialRepo()
	return &CatalogMaterialService{BaseService: NewBaseService(repo.BaseRepo)}
}

type SupplierService struct {
	*BaseService
}

func NewSupplierService() *SupplierService {
	repo := repositories.NewSupplierRepo()
	return &SupplierService{BaseService: NewBaseService(repo.BaseRepo)}
}

type CustomerService struct {
	*BaseService
}

func NewCustomerService() *CustomerService {
	repo := repositories.NewCustomerRepo()
	return &CustomerService{BaseService: NewBaseService(repo.BaseRepo)}
}

type TicketTypeService struct {
	*BaseService
}

func NewTicketTypeService() *TicketTypeService {
	repo := repositories.NewTicketTypeRepo()
	return &TicketTypeService{BaseService: NewBaseService(repo.BaseRepo)}
}

type VehicleCatalogService struct {
	*BaseService
}

func NewVehicleCatalogService() *VehicleCatalogService {
	repo := repositories.NewVehicleCatalogRepo()
	return &VehicleCatalogService{BaseService: NewBaseService(repo.BaseRepo)}
}

type StationService struct {
	*BaseService
}

func NewStationService() *StationService {
	repo := repositories.NewStationRepo()
	return &StationService{BaseService: NewBaseService(repo.BaseRepo)}
}

type EquipmentService struct {
	*BaseService
}

func NewEquipmentService() *EquipmentService {
	repo := repositories.NewEquipmentRepo()
	return &EquipmentService{BaseService: NewBaseService(repo.BaseRepo)}
}
