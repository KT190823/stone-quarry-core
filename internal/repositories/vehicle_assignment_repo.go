package repositories

import (
	"context"
	"encoding/json"
	"mo-da-backend/internal/database"
)

type VehicleAssignmentRepo struct {
	*BaseRepo
}

func NewVehicleAssignmentRepo() *VehicleAssignmentRepo {
	return &VehicleAssignmentRepo{BaseRepo: NewBaseRepo("vehicle_assignments", "id")}
}

func (r *VehicleAssignmentRepo) GetByID(idOrCode string) (map[string]interface{}, error) {
	ctx := context.Background()
	query := `SELECT row_to_json(t) FROM (
		SELECT * FROM vehicle_assignments WHERE CAST(id AS TEXT) = $1 OR code ILIKE $1 ORDER BY id DESC LIMIT 1
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

func (r *VehicleAssignmentRepo) GetActiveByPlate(plate string) (map[string]interface{}, error) {
	ctx := context.Background()
	query := `SELECT row_to_json(t) FROM (
		SELECT * FROM vehicle_assignments 
		WHERE license_plate ILIKE $1 AND status IN ('checked_in', 'in_progress', 'assigned')
		ORDER BY id DESC LIMIT 1
	) t`

	var data []byte
	err := database.Pool.QueryRow(ctx, query, plate).Scan(&data)
	if err != nil {
		return nil, err
	}

	var item map[string]interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}

	return item, nil
}
