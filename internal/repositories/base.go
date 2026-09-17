package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"mo-da-backend/internal/database"
)

// tableColumnsCache caches the physical column names per table so that
// Create/Update can ignore payload keys that are not real columns instead of
// failing the whole request with "column ... does not exist".
var tableColumnsCache sync.Map

type ListParams struct {
	Page       int
	PageSize   int
	Search     string
	Sort       string
	Order      string
	QuarryCode string
}

type BaseRepo struct {
	Table    string
	IDColumn string
}

func NewBaseRepo(table, idColumn string) *BaseRepo {
	return &BaseRepo{Table: table, IDColumn: idColumn}
}

// columnSet returns the set of physical column names of the backing table.
// It returns nil when the metadata lookup fails, in which case callers must
// NOT filter the payload (preserve the previous behaviour).
func (r *BaseRepo) columnSet() map[string]struct{} {
	if cached, ok := tableColumnsCache.Load(r.Table); ok {
		if cols, ok := cached.(map[string]struct{}); ok {
			return cols
		}
	}

	ctx := context.Background()
	rows, err := database.Pool.Query(ctx,
		`SELECT column_name FROM information_schema.columns WHERE table_schema = 'public' AND table_name = $1`,
		r.Table)
	if err != nil {
		return nil
	}
	defer rows.Close()

	cols := map[string]struct{}{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			cols[name] = struct{}{}
		}
	}
	if len(cols) == 0 {
		return nil
	}

	tableColumnsCache.Store(r.Table, cols)
	return cols
}

// NormalizeQuarryFilter maps any quarry code, slug, UUID, or name to canonical code and name
func NormalizeQuarryFilter(input string) (canonicalCode string, locationName string) {
	raw := strings.TrimSpace(input)
	if raw == "" || strings.EqualFold(raw, "all") || strings.EqualFold(raw, "ttc-all") {
		return "", ""
	}
	lower := strings.ToLower(raw)
	upper := strings.ToUpper(raw)

	switch {
	case strings.Contains(upper, "PT") || strings.Contains(lower, "phu-tho") || strings.Contains(lower, "phutho") || strings.Contains(lower, "phu tho") || strings.Contains(lower, "phú thọ") || strings.Contains(lower, "thanh ba") || strings.Contains(upper, "ECD6FC42"):
		return "MO-PT-01", "Phú Thọ"
	case strings.Contains(upper, "TU") || strings.Contains(lower, "tan-uyen") || strings.Contains(lower, "tanuyen") || strings.Contains(lower, "tan uyen") || strings.Contains(lower, "tân uyên") || strings.Contains(lower, "binh duong") || strings.Contains(lower, "bình dương") || strings.Contains(upper, "8D84B1B1"):
		return "MO-TU-02", "Tân Uyên"
	case strings.Contains(upper, "HN") || strings.Contains(lower, "ha-nam") || strings.Contains(lower, "hanam") || strings.Contains(lower, "ha nam") || strings.Contains(lower, "hà nam") || strings.Contains(lower, "kien khe") || strings.Contains(lower, "kiện khê") || strings.Contains(upper, "17B4EA5E"):
		return "MO-HN-03", "Hà Nam"
	case strings.Contains(upper, "BP") || strings.Contains(lower, "binh-phuoc") || strings.Contains(lower, "binhphuoc") || strings.Contains(lower, "binh phuoc") || strings.Contains(lower, "bình phước") || strings.Contains(lower, "chon thanh") || strings.Contains(lower, "chơn thành") || strings.Contains(upper, "4730475F"):
		return "MO-BP-04", "Bình Phước"
	default:
		return raw, ""
	}
}

func (r *BaseRepo) List(params ListParams) ([]map[string]interface{}, int, error) {
	ctx := context.Background()
	offset := (params.Page - 1) * params.PageSize
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 50
	}

	whereConditions := []string{}
	args := []interface{}{}
	argIdx := 1

	if params.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("CAST(row_to_json(%s) AS TEXT) ILIKE $%d", r.Table, argIdx))
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	qCode := strings.TrimSpace(params.QuarryCode)
	if qCode != "" && qCode != "TTC-ALL" && qCode != "ALL" && qCode != "all" {
		cCode, locName := NormalizeQuarryFilter(qCode)
		if locName != "" {
			whereConditions = append(whereConditions, fmt.Sprintf("(CAST(row_to_json(%s) AS TEXT) ILIKE $%d OR CAST(row_to_json(%s) AS TEXT) ILIKE $%d OR CAST(row_to_json(%s) AS TEXT) ILIKE $%d)", r.Table, argIdx, r.Table, argIdx+1, r.Table, argIdx+2))
			args = append(args, "%"+qCode+"%", "%"+cCode+"%", "%"+locName+"%")
			argIdx += 3
		} else {
			whereConditions = append(whereConditions, fmt.Sprintf("CAST(row_to_json(%s) AS TEXT) ILIKE $%d", r.Table, argIdx))
			args = append(args, "%"+qCode+"%")
			argIdx++
		}
	}

	where := ""
	if len(whereConditions) > 0 {
		where = " WHERE " + strings.Join(whereConditions, " AND ")
	}

	orderClause := " ORDER BY created_at DESC"
	if params.Sort != "" {
		dir := "ASC"
		if strings.ToUpper(params.Order) == "DESC" {
			dir = "DESC"
		}
		orderClause = fmt.Sprintf(" ORDER BY %s %s", params.Sort, dir)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s%s", r.Table, where)
	var total int
	err := database.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT row_to_json(t) FROM (SELECT * FROM %s%s%s LIMIT $%d OFFSET $%d) t", r.Table, where, orderClause, argIdx, argIdx+1)
	args = append(args, params.PageSize, offset)

	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			continue
		}
		var item map[string]interface{}
		if err := json.Unmarshal(data, &item); err != nil {
			continue
		}
		results = append(results, item)
	}

	return results, total, nil
}

func (r *BaseRepo) GetByID(id string) (map[string]interface{}, error) {
	ctx := context.Background()
	cols := r.columnSet()
	hasCode := false
	if _, ok := cols["code"]; ok {
		hasCode = true
	}
	var query string
	var args []interface{}
	numID, err := strconv.Atoi(id)
	if err == nil && numID > 0 {
		if hasCode {
			query = fmt.Sprintf("SELECT row_to_json(t) FROM (SELECT * FROM %s WHERE %s = $1 OR code = $2) t", r.Table, r.IDColumn)
			args = []interface{}{numID, id}
		} else {
			query = fmt.Sprintf("SELECT row_to_json(t) FROM (SELECT * FROM %s WHERE %s = $1) t", r.Table, r.IDColumn)
			args = []interface{}{numID}
		}
	} else {
		if hasCode {
			query = fmt.Sprintf("SELECT row_to_json(t) FROM (SELECT * FROM %s WHERE code = $1) t", r.Table)
			args = []interface{}{id}
		} else {
			query = fmt.Sprintf("SELECT row_to_json(t) FROM (SELECT * FROM %s WHERE %s = $1) t", r.Table, r.IDColumn)
			args = []interface{}{id}
		}
	}

	var data []byte
	scanErr := database.Pool.QueryRow(ctx, query, args...).Scan(&data)
	if scanErr != nil {
		return nil, scanErr
	}

	var item map[string]interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *BaseRepo) Create(data map[string]interface{}) (map[string]interface{}, error) {
	ctx := context.Background()

	colNames := []string{}
	vals := []interface{}{}
	placeholders := []string{}
	idx := 1
	colSet := r.columnSet()

	for k, v := range data {
		if k == "id" && v == nil {
			continue
		}
		if k == "created_at" || k == "updated_at" {
			continue
		}
		col := toSnakeCase(k)
		if colSet != nil {
			if _, ok := colSet[col]; !ok {
				continue
			}
		}
		colNames = append(colNames, col)
		vals = append(vals, toJSON(v))
		placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
		idx++
	}

	if len(colNames) == 0 {
		return nil, fmt.Errorf("no valid columns to insert into %s", r.Table)
	}

	query := fmt.Sprintf(
		"WITH ins AS (INSERT INTO %s (%s) VALUES (%s) RETURNING *) SELECT row_to_json(ins) FROM ins",
		r.Table, strings.Join(colNames, ", "), strings.Join(placeholders, ", "),
	)

	var result []byte
	err := database.Pool.QueryRow(ctx, query, vals...).Scan(&result)
	if err != nil {
		return nil, err
	}

	var item map[string]interface{}
	json.Unmarshal(result, &item)
	return item, nil
}

func (r *BaseRepo) Update(id string, data map[string]interface{}) (map[string]interface{}, error) {
	ctx := context.Background()

	sets := []string{}
	vals := []interface{}{}
	idx := 1
	colSet := r.columnSet()

	for k, v := range data {
		if k == r.IDColumn || k == "created_at" || k == "id" {
			continue
		}
		col := toSnakeCase(k)
		if colSet != nil {
			if _, ok := colSet[col]; !ok {
				continue
			}
		}
		sets = append(sets, fmt.Sprintf("%s = $%d", col, idx))
		vals = append(vals, toJSON(v))
		idx++
	}

	if len(sets) == 0 {
		return r.GetByID(id)
	}

	sets = append(sets, fmt.Sprintf("updated_at = NOW()"))
	vals = append(vals, id)

	query := fmt.Sprintf(
		"WITH upd AS (UPDATE %s SET %s WHERE %s = $%d RETURNING *) SELECT row_to_json(upd) FROM upd",
		r.Table, strings.Join(sets, ", "), r.IDColumn, idx,
	)

	var result []byte
	err := database.Pool.QueryRow(ctx, query, vals...).Scan(&result)
	if err != nil {
		return nil, err
	}

	var item map[string]interface{}
	json.Unmarshal(result, &item)
	return item, nil
}

func (r *BaseRepo) Delete(id string) error {
	ctx := context.Background()
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", r.Table, r.IDColumn)
	_, err := database.Pool.Exec(ctx, query, id)
	return err
}

func (r *BaseRepo) ListJSONB(params ListParams) ([]map[string]interface{}, int, error) {
	ctx := context.Background()
	offset := (params.Page - 1) * params.PageSize
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 50
	}

	whereConditions := []string{}
	args := []interface{}{}
	argIdx := 1

	if params.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("data::text ILIKE $%d", argIdx))
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	qCode := strings.TrimSpace(params.QuarryCode)
	if qCode != "" && qCode != "TTC-ALL" && qCode != "ALL" && qCode != "all" {
		cCode, locName := NormalizeQuarryFilter(qCode)
		if locName != "" {
			whereConditions = append(whereConditions, fmt.Sprintf("(data::text ILIKE $%d OR data::text ILIKE $%d OR data::text ILIKE $%d)", argIdx, argIdx+1, argIdx+2))
			args = append(args, "%"+qCode+"%", "%"+cCode+"%", "%"+locName+"%")
			argIdx += 3
		} else {
			whereConditions = append(whereConditions, fmt.Sprintf("data::text ILIKE $%d", argIdx))
			args = append(args, "%"+qCode+"%")
			argIdx++
		}
	}

	where := ""
	if len(whereConditions) > 0 {
		where = " WHERE " + strings.Join(whereConditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s%s", r.Table, where)
	var total int
	database.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)

	query := fmt.Sprintf("SELECT data FROM %s%s LIMIT $%d OFFSET $%d", r.Table, where, argIdx, argIdx+1)
	args = append(args, params.PageSize, offset)

	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			continue
		}
		var item map[string]interface{}
		if err := json.Unmarshal(data, &item); err != nil {
			continue
		}
		results = append(results, item)
	}

	return results, total, nil
}

func toJSON(v interface{}) interface{} {
	switch v.(type) {
	case map[string]interface{}, []interface{}:
		b, _ := json.Marshal(v)
		return string(b)
	default:
		return v
	}
}

func toSnakeCase(s string) string {
	var result []byte
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, byte(c+32))
		} else {
			result = append(result, byte(c))
		}
	}
	return string(result)
}
