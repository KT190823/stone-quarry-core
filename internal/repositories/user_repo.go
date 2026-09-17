package repositories

import (
	"context"
	"encoding/json"
	"strconv"

	"mo-da-backend/internal/database"
)

type UserRepo struct {
	*BaseRepo
}

func NewUserRepo() *UserRepo {
	return &UserRepo{BaseRepo: NewBaseRepo("users", "id")}
}

func (r *UserRepo) GetByID(idOrCodeOrUsername string) (map[string]interface{}, error) {
	ctx := context.Background()
	numID, err := strconv.Atoi(idOrCodeOrUsername)
	var query string
	var args []interface{}
	if err == nil && numID > 0 {
		query = "SELECT row_to_json(u) FROM (SELECT * FROM users WHERE id = $1 OR code = $2 OR username = $3) u"
		args = []interface{}{numID, idOrCodeOrUsername, idOrCodeOrUsername}
	} else {
		query = "SELECT row_to_json(u) FROM (SELECT * FROM users WHERE code = $1 OR username = $2) u"
		args = []interface{}{idOrCodeOrUsername, idOrCodeOrUsername}
	}
	var data []byte
	if err := database.Pool.QueryRow(ctx, query, args...).Scan(&data); err != nil {
		return nil, err
	}
	var item map[string]interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}
	return item, nil
}

type UserRoleRepo struct {
	*BaseRepo
}

func NewUserRoleRepo() *UserRoleRepo {
	return &UserRoleRepo{BaseRepo: NewBaseRepo("user_roles", "id")}
}

func (r *UserRoleRepo) GetByID(idOrRole string) (map[string]interface{}, error) {
	ctx := context.Background()
	numID, err := strconv.Atoi(idOrRole)
	var query string
	var args []interface{}
	if err == nil && numID > 0 {
		query = "SELECT row_to_json(ur) FROM (SELECT * FROM user_roles WHERE id = $1 OR role = $2) ur"
		args = []interface{}{numID, idOrRole}
	} else {
		query = "SELECT row_to_json(ur) FROM (SELECT * FROM user_roles WHERE role = $1) ur"
		args = []interface{}{idOrRole}
	}
	var data []byte
	if err := database.Pool.QueryRow(ctx, query, args...).Scan(&data); err != nil {
		return nil, err
	}
	var item map[string]interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}
	return item, nil
}

type UserLogRepo struct {
	*BaseRepo
}

func NewUserLogRepo() *UserLogRepo {
	return &UserLogRepo{BaseRepo: NewBaseRepo("user_logs", "id")}
}

func (r *UserLogRepo) GetByID(id string) (map[string]interface{}, error) {
	ctx := context.Background()
	numID, err := strconv.Atoi(id)
	var query string
	var args []interface{}
	if err == nil && numID > 0 {
		query = "SELECT row_to_json(ul) FROM (SELECT * FROM user_logs WHERE id = $1) ul"
		args = []interface{}{numID}
	} else {
		query = "SELECT row_to_json(ul) FROM (SELECT * FROM user_logs WHERE id::text = $1) ul"
		args = []interface{}{id}
	}
	var data []byte
	if err := database.Pool.QueryRow(ctx, query, args...).Scan(&data); err != nil {
		return nil, err
	}
	var item map[string]interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}
	return item, nil
}

type ReportRepo struct {
	*BaseRepo
}

func NewReportRepo() *ReportRepo {
	return &ReportRepo{BaseRepo: NewBaseRepo("reports", "id")}
}

type SettingRepo struct {
	*BaseRepo
}

func NewSettingRepo() *SettingRepo {
	return &SettingRepo{BaseRepo: NewBaseRepo("settings", "id")}
}
