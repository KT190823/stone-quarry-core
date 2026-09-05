package repositories

import (
	"context"
	"encoding/json"

	"mo-da-backend/internal/database"
)

type SalesContractRepo struct {
	*BaseRepo
}

func NewSalesContractRepo() *SalesContractRepo {
	return &SalesContractRepo{BaseRepo: NewBaseRepo("sales_contracts", "id")}
}

func (r *SalesContractRepo) GetByID(idOrCode string) (map[string]interface{}, error) {
	ctx := context.Background()
	query := `SELECT row_to_json(t) FROM (
		SELECT * FROM sales_contracts WHERE id::text = $1 OR code ILIKE $1 ORDER BY created_at DESC LIMIT 1
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

type ConsolidatedDeliveryOrderRepo struct {
	*BaseRepo
}

func NewConsolidatedDeliveryOrderRepo() *ConsolidatedDeliveryOrderRepo {
	return &ConsolidatedDeliveryOrderRepo{BaseRepo: NewBaseRepo("consolidated_delivery_orders", "id")}
}

func (r *ConsolidatedDeliveryOrderRepo) GetByID(idOrCode string) (map[string]interface{}, error) {
	ctx := context.Background()
	query := `SELECT row_to_json(t) FROM (
		SELECT * FROM consolidated_delivery_orders WHERE id::text = $1 OR code ILIKE $1 ORDER BY created_at DESC LIMIT 1
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
