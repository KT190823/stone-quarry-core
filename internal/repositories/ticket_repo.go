package repositories

import (
	"context"
	"encoding/json"
	"mo-da-backend/internal/database"
)

type TicketRepo struct {
	*BaseRepo
}

func NewTicketRepo() *TicketRepo {
	return &TicketRepo{BaseRepo: NewBaseRepo("tickets", "id")}
}

func (r *TicketRepo) GetByID(idOrPlate string) (map[string]interface{}, error) {
	ctx := context.Background()
	query := "SELECT row_to_json(t) FROM (SELECT * FROM tickets WHERE id = $1 OR bien_so ILIKE $1 ORDER BY created_at DESC LIMIT 1) t"

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
