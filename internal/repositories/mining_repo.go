package repositories

import (
	"context"
	"encoding/json"

	"mo-da-backend/internal/database"
)

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

func (r *BlastingPassportRepo) GetByID(idOrCode string) (map[string]interface{}, error) {
	ctx := context.Background()
	query := "SELECT row_to_json(t) FROM (SELECT * FROM blasting_passports WHERE id = $1 OR code ILIKE $1 ORDER BY created_at DESC LIMIT 1) t"

	var data []byte
	err := database.Pool.QueryRow(ctx, query, idOrCode).Scan(&data)
	if err != nil {
		return nil, err
	}

	var item map[string]interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}

	return item, nil
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

func (r *EquipmentFuelLogRepo) GetByID(idOrCode string) (map[string]interface{}, error) {
	ctx := context.Background()
	query := `SELECT row_to_json(t) FROM (
		SELECT * FROM equipment_fuel_logs WHERE id = $1 OR equipment_code ILIKE $1 ORDER BY created_at DESC LIMIT 1
	) t`

	var data []byte
	err := database.Pool.QueryRow(ctx, query, idOrCode).Scan(&data)
	if err != nil {
		return nil, err
	}

	var item map[string]interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}

	return item, nil
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
