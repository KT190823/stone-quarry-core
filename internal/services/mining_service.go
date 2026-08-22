package services

import (
	"mo-da-backend/internal/repositories"
)

type MiningPermitService struct {
	*BaseService
}

func NewMiningPermitService() *MiningPermitService {
	repo := repositories.NewMiningPermitRepo()
	return &MiningPermitService{BaseService: NewBaseService(repo.BaseRepo)}
}

type MiningPlanService struct {
	*BaseService
}

func NewMiningPlanService() *MiningPlanService {
	repo := repositories.NewMiningPlanRepo()
	return &MiningPlanService{BaseService: NewBaseService(repo.BaseRepo)}
}

type BlastingPassportService struct {
	*BaseService
}

func NewBlastingPassportService() *BlastingPassportService {
	repo := repositories.NewBlastingPassportRepo()
	return &BlastingPassportService{BaseService: NewBaseService(repo.BaseRepo)}
}

type CrusherPlantService struct {
	*BaseService
}

func NewCrusherPlantService() *CrusherPlantService {
	repo := repositories.NewCrusherPlantRepo()
	return &CrusherPlantService{BaseService: NewBaseService(repo.BaseRepo)}
}

type EquipmentFuelLogService struct {
	*BaseService
}

func NewEquipmentFuelLogService() *EquipmentFuelLogService {
	repo := repositories.NewEquipmentFuelLogRepo()
	return &EquipmentFuelLogService{BaseService: NewBaseService(repo.BaseRepo)}
}

type StatutoryReportService struct {
	*BaseService
}

func NewStatutoryReportService() *StatutoryReportService {
	repo := repositories.NewStatutoryReportRepo()
	return &StatutoryReportService{BaseService: NewBaseService(repo.BaseRepo)}
}

type ResourceTaxService struct {
	*BaseService
}

func NewResourceTaxService() *ResourceTaxService {
	repo := repositories.NewResourceTaxRepo()
	return &ResourceTaxService{BaseService: NewBaseService(repo.BaseRepo)}
}

type ProductionStageService struct {
	*BaseService
}

func NewProductionStageService() *ProductionStageService {
	repo := repositories.NewProductionStageRepo()
	return &ProductionStageService{BaseService: NewBaseService(repo.BaseRepo)}
}
