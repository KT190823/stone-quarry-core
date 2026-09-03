package repositories

import (
	"context"
	"encoding/json"
	"mo-da-backend/internal/database"
)

type VehicleRepo struct {
	*BaseRepo
}

func NewVehicleRepo() *VehicleRepo {
	return &VehicleRepo{BaseRepo: NewBaseRepo("vehicles", "id")}
}

func (r *VehicleRepo) GetByID(idOrPlate string) (map[string]interface{}, error) {
	ctx := context.Background()
	query := "SELECT row_to_json(t) FROM (SELECT * FROM vehicles WHERE CAST(id AS TEXT) = $1 OR bs ILIKE $1 ORDER BY id ASC LIMIT 1) t"

	var data []byte
	err := database.Pool.QueryRow(ctx, query, idOrPlate).Scan(&data)
	if err != nil {
		return nil, err
	}

	var item map[string]interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}

	return item, nil
}
