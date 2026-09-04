package handlers

import (
	"net/http"
	"strconv"
	"time"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/models"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	search := r.URL.Query().Get("search")
	category := r.URL.Query().Get("category")

	query := `
		SELECT id, code, name, COALESCE(category, '') as category, COALESCE(unit, 'm³') as unit,
		       COALESCE(density, 1.5) as density, COALESCE(sale_price, 0) as sale_price,
		       COALESCE(purchase_price, 0) as purchase_price, COALESCE(storage_loc, '') as storage_loc,
		       COALESCE(standard, '') as standard, COALESCE(min_stock, 0) as min_stock,
		       COALESCE(current_stock, 0) as current_stock, COALESCE(status, 'active') as status,
		       COALESCE(notes, '') as notes, created_at, updated_at
		FROM inventory_products
		WHERE 1=1
	`
	var args []interface{}
	idx := 1

	if search != "" {
		query += ` AND (code ILIKE $` + strconv.Itoa(idx) + ` OR name ILIKE $` + strconv.Itoa(idx) + ` OR storage_loc ILIKE $` + strconv.Itoa(idx) + `)`
		args = append(args, "%"+search+"%")
		idx++
	}
	if category != "" && category != "all" {
		query += ` AND category = $` + strconv.Itoa(idx)
		args = append(args, category)
		idx++
	}

	query += ` ORDER BY id ASC`

	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		JSON(w, []interface{}{})
		return
	}
	defer rows.Close()

	var list []models.InventoryProduct
	for rows.Next() {
		var p models.InventoryProduct
		if err := rows.Scan(&p.ID, &p.Code, &p.Name, &p.Category, &p.Unit, &p.Density, &p.SalePrice, &p.PurchasePrice, &p.StorageLoc, &p.Standard, &p.MinStock, &p.CurrentStock, &p.Status, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err == nil {
			list = append(list, p)
		}
	}
	if list == nil {
		list = []models.InventoryProduct{}
	}
	JSON(w, map[string]interface{}{"data": list, "total": len(list)})
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	if idOrCode == "" {
		JSONError(w, http.StatusBadRequest, "Thiếu mã hoặc ID sản phẩm")
		return
	}

	var p models.InventoryProduct
	var err error

	numID, parseErr := strconv.Atoi(idOrCode)
	if parseErr == nil {
		err = database.Pool.QueryRow(ctx, `
			SELECT id, code, name, COALESCE(category, ''), COALESCE(unit, 'm³'),
			       COALESCE(density, 1.5), COALESCE(sale_price, 0), COALESCE(purchase_price, 0),
			       COALESCE(storage_loc, ''), COALESCE(standard, ''), COALESCE(min_stock, 0),
			       COALESCE(current_stock, 0), COALESCE(status, 'active'), COALESCE(notes, ''),
			       created_at, updated_at
			FROM inventory_products
			WHERE id = $1 OR code = $2
		`, numID, idOrCode).Scan(&p.ID, &p.Code, &p.Name, &p.Category, &p.Unit, &p.Density, &p.SalePrice, &p.PurchasePrice, &p.StorageLoc, &p.Standard, &p.MinStock, &p.CurrentStock, &p.Status, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	} else {
		err = database.Pool.QueryRow(ctx, `
			SELECT id, code, name, COALESCE(category, ''), COALESCE(unit, 'm³'),
			       COALESCE(density, 1.5), COALESCE(sale_price, 0), COALESCE(purchase_price, 0),
			       COALESCE(storage_loc, ''), COALESCE(standard, ''), COALESCE(min_stock, 0),
			       COALESCE(current_stock, 0), COALESCE(status, 'active'), COALESCE(notes, ''),
			       created_at, updated_at
			FROM inventory_products
			WHERE code = $1
		`, idOrCode).Scan(&p.ID, &p.Code, &p.Name, &p.Category, &p.Unit, &p.Density, &p.SalePrice, &p.PurchasePrice, &p.StorageLoc, &p.Standard, &p.MinStock, &p.CurrentStock, &p.Status, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	}

	if err != nil {
		JSONError(w, http.StatusNotFound, "Không tìm thấy sản phẩm "+idOrCode)
		return
	}
	JSON(w, p)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var p models.InventoryProduct
	if err := decodeJSON(r, &p); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	if p.Code == "" {
		p.Code = "SP-" + time.Now().Format("0601021504")
	}

	err := database.Pool.QueryRow(ctx, `
		INSERT INTO inventory_products (code, name, category, unit, density, sale_price, purchase_price, storage_loc, standard, min_stock, current_stock, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`, p.Code, p.Name, p.Category, p.Unit, p.Density, p.SalePrice, p.PurchasePrice, p.StorageLoc, p.Standard, p.MinStock, p.CurrentStock, p.Status, p.Notes).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể tạo sản phẩm: "+err.Error())
		return
	}
	JSON(w, p)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	var p models.InventoryProduct
	if err := decodeJSON(r, &p); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	res, err := database.Pool.Exec(ctx, `
		UPDATE inventory_products
		SET name = $1, category = $2, unit = $3, density = $4, sale_price = $5,
		    purchase_price = $6, storage_loc = $7, standard = $8, min_stock = $9,
		    current_stock = $10, status = $11, notes = $12, updated_at = NOW()
		WHERE code = $13 OR id = $14
	`, p.Name, p.Category, p.Unit, p.Density, p.SalePrice, p.PurchasePrice, p.StorageLoc, p.Standard, p.MinStock, p.CurrentStock, p.Status, p.Notes, idOrCode, p.ID)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể cập nhật sản phẩm: "+err.Error())
		return
	}
	if res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy sản phẩm")
		return
	}
	JSON(w, map[string]interface{}{"status": "success", "message": "Cập nhật sản phẩm thành công"})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	res, err := database.Pool.Exec(ctx, `DELETE FROM inventory_products WHERE code = $1 OR id = $2`, idOrCode, numID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể xóa sản phẩm: "+err.Error())
		return
	}
	if res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy sản phẩm")
		return
	}
	JSON(w, map[string]interface{}{"status": "success", "message": "Đã xóa sản phẩm"})
}
