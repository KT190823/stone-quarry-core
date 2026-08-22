package repositories

type MiningPermitRepo struct {
	*BaseRepo
}

func NewMiningPermitRepo() *MiningPermitRepo {
	return &MiningPermitRepo{BaseRepo: NewBaseRepo("mining_permits", "id")}
}

type MiningPlanRepo struct {
	*BaseRepo
}

func NewMiningPlanRepo() *MiningPlanRepo {
	return &MiningPlanRepo{BaseRepo: NewBaseRepo("mining_plans", "id")}
}

type BlastingPassportRepo struct {
	*BaseRepo
}

func NewBlastingPassportRepo() *BlastingPassportRepo {
	return &BlastingPassportRepo{BaseRepo: NewBaseRepo("blasting_passports", "id")}
}

type CrusherPlantRepo struct {
	*BaseRepo
}

func NewCrusherPlantRepo() *CrusherPlantRepo {
	return &CrusherPlantRepo{BaseRepo: NewBaseRepo("crusher_plants", "id")}
}

type EquipmentFuelLogRepo struct {
	*BaseRepo
}

func NewEquipmentFuelLogRepo() *EquipmentFuelLogRepo {
	return &EquipmentFuelLogRepo{BaseRepo: NewBaseRepo("equipment_fuel_logs", "id")}
}

type StatutoryReportRepo struct {
	*BaseRepo
}

func NewStatutoryReportRepo() *StatutoryReportRepo {
	return &StatutoryReportRepo{BaseRepo: NewBaseRepo("statutory_reports", "id")}
}

type ResourceTaxRepo struct {
	*BaseRepo
}

func NewResourceTaxRepo() *ResourceTaxRepo {
	return &ResourceTaxRepo{BaseRepo: NewBaseRepo("resource_taxes", "id")}
}

type ProductionStageRepo struct {
	*BaseRepo
}

func NewProductionStageRepo() *ProductionStageRepo {
	return &ProductionStageRepo{BaseRepo: NewBaseRepo("production_stages", "id")}
}
