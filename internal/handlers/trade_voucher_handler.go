package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/models"
)

type TradeVoucherHandler struct{}

func NewTradeVoucherHandler() *TradeVoucherHandler {
	return &TradeVoucherHandler{}
}

// ==========================================
// 1. PURCHASE VOUCHERS (PHIẾU MUA HÀNG)
// ==========================================

func (h *TradeVoucherHandler) ListPurchases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	query := `
		SELECT id, code, COALESCE(supplier_code, ''), COALESCE(supplier_name, ''),
		       COALESCE(date, ''), COALESCE(warehouse_loc, ''), COALESCE(total_amount, 0),
		       COALESCE(vat_amount, 0), COALESCE(grand_total, 0), COALESCE(payment_status, 'unpaid'),
		       COALESCE(status, 'completed'), COALESCE(created_by, ''), COALESCE(notes, ''),
		       created_at, updated_at
		FROM purchase_vouchers
		WHERE 1=1
	`
	var args []interface{}
	idx := 1
	if search != "" {
		query += ` AND (code ILIKE $` + strconv.Itoa(idx) + ` OR supplier_name ILIKE $` + strconv.Itoa(idx) + ` OR warehouse_loc ILIKE $` + strconv.Itoa(idx) + `)`
		args = append(args, "%"+search+"%")
		idx++
	}
	if status != "" && status != "all" {
		query += ` AND status = $` + strconv.Itoa(idx)
		args = append(args, status)
		idx++
	}
	query += ` ORDER BY id DESC`

	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		JSON(w, []interface{}{})
		return
	}
	defer rows.Close()

	var list []models.PurchaseVoucher
	for rows.Next() {
		var pv models.PurchaseVoucher
		if err := rows.Scan(&pv.ID, &pv.Code, &pv.SupplierCode, &pv.SupplierName, &pv.Date, &pv.WarehouseLoc, &pv.TotalAmount, &pv.VatAmount, &pv.GrandTotal, &pv.PaymentStatus, &pv.Status, &pv.CreatedBy, &pv.Notes, &pv.CreatedAt, &pv.UpdatedAt); err == nil {
			list = append(list, pv)
		}
	}
	if list == nil {
		list = []models.PurchaseVoucher{}
	}
	JSON(w, map[string]interface{}{"data": list, "total": len(list)})
}

func (h *TradeVoucherHandler) GetPurchase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	var pv models.PurchaseVoucher
	err := database.Pool.QueryRow(ctx, `
		SELECT id, code, COALESCE(supplier_code, ''), COALESCE(supplier_name, ''),
		       COALESCE(date, ''), COALESCE(warehouse_loc, ''), COALESCE(total_amount, 0),
		       COALESCE(vat_amount, 0), COALESCE(grand_total, 0), COALESCE(payment_status, 'unpaid'),
		       COALESCE(status, 'completed'), COALESCE(created_by, ''), COALESCE(notes, ''),
		       created_at, updated_at
		FROM purchase_vouchers
		WHERE code = $1 OR id = $2
	`, idOrCode, numID).Scan(&pv.ID, &pv.Code, &pv.SupplierCode, &pv.SupplierName, &pv.Date, &pv.WarehouseLoc, &pv.TotalAmount, &pv.VatAmount, &pv.GrandTotal, &pv.PaymentStatus, &pv.Status, &pv.CreatedBy, &pv.Notes, &pv.CreatedAt, &pv.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu mua "+idOrCode)
		return
	}

	// Fetch items
	itemRows, itemErr := database.Pool.Query(ctx, `
		SELECT id, COALESCE(voucher_id, 0), COALESCE(voucher_code, ''), COALESCE(product_code, ''), COALESCE(product_name, ''),
		       COALESCE(unit, 'm³'), COALESCE(density, 1.5), COALESCE(unit_price, 0),
		       COALESCE(quantity, 0), COALESCE(total_amount, 0), COALESCE(weight_ton, 0),
		       COALESCE(standard, ''), COALESCE(storage_loc, ''), COALESCE(notes, ''),
		       COALESCE(vat_rate, 10)
		FROM purchase_voucher_items
		WHERE (voucher_id > 0 AND voucher_id = $1) OR voucher_code = $2
		ORDER BY id ASC
	`, pv.ID, pv.Code)

	if itemErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var itm models.PurchaseItem
			if err := itemRows.Scan(&itm.ID, &itm.VoucherID, &itm.VoucherCode, &itm.ProductCode, &itm.ProductName, &itm.Unit, &itm.Density, &itm.UnitPrice, &itm.Quantity, &itm.TotalAmount, &itm.WeightTon, &itm.Standard, &itm.StorageLoc, &itm.Notes, &itm.VatRate); err == nil {
				pv.Items = append(pv.Items, itm)
			}
		}
	}
	if pv.Items == nil {
		pv.Items = []models.PurchaseItem{}
	}

	JSON(w, pv)
}

func (h *TradeVoucherHandler) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var pv models.PurchaseVoucher
	if err := decodeJSON(r, &pv); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	if pv.Code == "" {
		pv.Code = "PM-" + time.Now().Format("20060102-1504")
	}
	if pv.Date == "" {
		pv.Date = time.Now().Format("2006-01-02")
	}
	if pv.Status == "" {
		pv.Status = "completed"
	}

	err := database.Pool.QueryRow(ctx, `
		INSERT INTO purchase_vouchers (code, supplier_code, supplier_name, date, warehouse_loc, total_amount, vat_amount, grand_total, payment_status, status, created_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`, pv.Code, pv.SupplierCode, pv.SupplierName, pv.Date, pv.WarehouseLoc, pv.TotalAmount, pv.VatAmount, pv.GrandTotal, pv.PaymentStatus, pv.Status, pv.CreatedBy, pv.Notes).Scan(&pv.ID, &pv.CreatedAt, &pv.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể tạo phiếu mua: "+err.Error())
		return
	}

	for _, itm := range pv.Items {
		Pool := database.Pool
		vatRate := itm.VatRate
		if vatRate == 0 {
			vatRate = 10
		}
		Pool.Exec(ctx, `
			INSERT INTO purchase_voucher_items (voucher_id, voucher_code, product_code, product_name, unit, density, unit_price, quantity, total_amount, weight_ton, standard, storage_loc, notes, vat_rate)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		`, pv.ID, pv.Code, itm.ProductCode, itm.ProductName, itm.Unit, itm.Density, itm.UnitPrice, itm.Quantity, itm.TotalAmount, itm.WeightTon, itm.Standard, itm.StorageLoc, itm.Notes, vatRate)
	}

	JSON(w, pv)
}

func (h *TradeVoucherHandler) UpdatePurchase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	var pv models.PurchaseVoucher
	if err := decodeJSON(r, &pv); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	numID, _ := strconv.Atoi(idOrCode)
	res, err := database.Pool.Exec(ctx, `
		UPDATE purchase_vouchers
		SET supplier_code = $1, supplier_name = $2, date = $3, warehouse_loc = $4,
		    total_amount = $5, vat_amount = $6, grand_total = $7, payment_status = $8,
		    status = $9, notes = $10, updated_at = NOW()
		WHERE code = $11 OR id = $12
	`, pv.SupplierCode, pv.SupplierName, pv.Date, pv.WarehouseLoc, pv.TotalAmount, pv.VatAmount, pv.GrandTotal, pv.PaymentStatus, pv.Status, pv.Notes, idOrCode, numID)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể cập nhật phiếu mua: "+err.Error())
		return
	}
	if res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu mua")
		return
	}

	// Update items if supplied
	if len(pv.Items) > 0 {
		database.Pool.Exec(ctx, `DELETE FROM purchase_voucher_items WHERE voucher_code = $1 OR voucher_id = $2`, idOrCode, numID)
		for _, itm := range pv.Items {
			vatRate := itm.VatRate
			if vatRate == 0 {
				vatRate = 10
			}
			database.Pool.Exec(ctx, `
				INSERT INTO purchase_voucher_items (voucher_id, voucher_code, product_code, product_name, unit, density, unit_price, quantity, total_amount, weight_ton, standard, storage_loc, notes, vat_rate)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			`, pv.ID, pv.Code, itm.ProductCode, itm.ProductName, itm.Unit, itm.Density, itm.UnitPrice, itm.Quantity, itm.TotalAmount, itm.WeightTon, itm.Standard, itm.StorageLoc, itm.Notes, vatRate)
		}
	}

	JSON(w, map[string]interface{}{"status": "success", "message": "Cập nhật phiếu mua thành công"})
}

func (h *TradeVoucherHandler) DeletePurchase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	database.Pool.Exec(ctx, `DELETE FROM purchase_voucher_items WHERE voucher_code = $1 OR voucher_id = $2`, idOrCode, numID)
	res, err := database.Pool.Exec(ctx, `DELETE FROM purchase_vouchers WHERE code = $1 OR id = $2`, idOrCode, numID)
	if err != nil || res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu mua")
		return
	}
	JSON(w, map[string]interface{}{"status": "success", "message": "Đã xóa phiếu mua"})
}

// ==========================================
// 2. SALES VOUCHERS (PHIẾU BÁN HÀNG)
// ==========================================

func (h *TradeVoucherHandler) ListSales(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	query := `
		SELECT id, code, COALESCE(customer_code, ''), COALESCE(customer_name, ''),
		       COALESCE(date, ''), COALESCE(warehouse_loc, ''), COALESCE(license_plate, ''),
		       COALESCE(ticket_code, ''), COALESCE(total_amount, 0), COALESCE(vat_amount, 0),
		       COALESCE(grand_total, 0), COALESCE(paid_amount, 0), COALESCE(debt_amount, 0),
		       COALESCE(payment_status, 'paid'), COALESCE(status, 'completed'),
		       COALESCE(created_by, ''), COALESCE(notes, ''), created_at, updated_at
		FROM sales_vouchers
		WHERE 1=1
	`
	var args []interface{}
	idx := 1
	if search != "" {
		query += ` AND (code ILIKE $` + strconv.Itoa(idx) + ` OR customer_name ILIKE $` + strconv.Itoa(idx) + ` OR license_plate ILIKE $` + strconv.Itoa(idx) + ` OR ticket_code ILIKE $` + strconv.Itoa(idx) + `)`
		args = append(args, "%"+search+"%")
		idx++
	}
	if status != "" && status != "all" {
		query += ` AND status = $` + strconv.Itoa(idx)
		args = append(args, status)
		idx++
	}
	query += ` ORDER BY id DESC`

	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		JSON(w, []interface{}{})
		return
	}
	defer rows.Close()

	var list []models.SalesVoucher
	for rows.Next() {
		var sv models.SalesVoucher
		if err := rows.Scan(&sv.ID, &sv.Code, &sv.CustomerCode, &sv.CustomerName, &sv.Date, &sv.WarehouseLoc, &sv.LicensePlate, &sv.TicketCode, &sv.TotalAmount, &sv.VatAmount, &sv.GrandTotal, &sv.PaidAmount, &sv.DebtAmount, &sv.PaymentStatus, &sv.Status, &sv.CreatedBy, &sv.Notes, &sv.CreatedAt, &sv.UpdatedAt); err == nil {
			list = append(list, sv)
		}
	}
	if list == nil {
		list = []models.SalesVoucher{}
	}
	JSON(w, map[string]interface{}{"data": list, "total": len(list)})
}

func (h *TradeVoucherHandler) GetSales(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	var sv models.SalesVoucher
	err := database.Pool.QueryRow(ctx, `
		SELECT id, code, COALESCE(customer_code, ''), COALESCE(customer_name, ''),
		       COALESCE(date, ''), COALESCE(warehouse_loc, ''), COALESCE(license_plate, ''),
		       COALESCE(ticket_code, ''), COALESCE(total_amount, 0), COALESCE(vat_amount, 0),
		       COALESCE(grand_total, 0), COALESCE(paid_amount, 0), COALESCE(debt_amount, 0),
		       COALESCE(payment_status, 'paid'), COALESCE(status, 'completed'),
		       COALESCE(created_by, ''), COALESCE(notes, ''), created_at, updated_at
		FROM sales_vouchers
		WHERE code = $1 OR id = $2
	`, idOrCode, numID).Scan(&sv.ID, &sv.Code, &sv.CustomerCode, &sv.CustomerName, &sv.Date, &sv.WarehouseLoc, &sv.LicensePlate, &sv.TicketCode, &sv.TotalAmount, &sv.VatAmount, &sv.GrandTotal, &sv.PaidAmount, &sv.DebtAmount, &sv.PaymentStatus, &sv.Status, &sv.CreatedBy, &sv.Notes, &sv.CreatedAt, &sv.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu bán "+idOrCode)
		return
	}

	// Fetch items
	itemRows, itemErr := database.Pool.Query(ctx, `
		SELECT id, COALESCE(voucher_id, 0), COALESCE(voucher_code, ''), COALESCE(product_code, ''), COALESCE(product_name, ''),
		       COALESCE(unit, 'm³'), COALESCE(density, 1.5), COALESCE(unit_price, 0),
		       COALESCE(quantity, 0), COALESCE(total_amount, 0), COALESCE(weight_ton, 0),
		       COALESCE(standard, ''), COALESCE(storage_loc, ''), COALESCE(notes, ''),
		       COALESCE(vat_rate, 10)
		FROM sales_voucher_items
		WHERE (voucher_id > 0 AND voucher_id = $1) OR voucher_code = $2
		ORDER BY id ASC
	`, sv.ID, sv.Code)

	if itemErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var itm models.SalesItem
			if err := itemRows.Scan(&itm.ID, &itm.VoucherID, &itm.VoucherCode, &itm.ProductCode, &itm.ProductName, &itm.Unit, &itm.Density, &itm.UnitPrice, &itm.Quantity, &itm.TotalAmount, &itm.WeightTon, &itm.Standard, &itm.StorageLoc, &itm.Notes, &itm.VatRate); err == nil {
				sv.Items = append(sv.Items, itm)
			}
		}
	}
	if len(sv.Items) == 0 && (sv.TicketCode != "" || strings.HasPrefix(sv.Code, "PB-TK-")) {
		ticketID := sv.TicketCode
		if ticketID == "" && strings.HasPrefix(sv.Code, "PB-TK-") {
			ticketID = strings.TrimPrefix(sv.Code, "PB-")
		}
		var matHang, quyCach, klStr, tramCan string
		var donGia float64
		tErr := database.Pool.QueryRow(ctx, `
			SELECT COALESCE(mat_hang, ''), COALESCE(quy_cach, ''), COALESCE(kl_tinh_tien::text, '0'),
			       COALESCE(don_gia, 0), COALESCE(tram_can, '')
			FROM tickets WHERE id = $1
		`, ticketID).Scan(&matHang, &quyCach, &klStr, &donGia, &tramCan)
		if tErr == nil && matHang != "" {
			qty, _ := strconv.ParseFloat(strings.TrimSpace(klStr), 64)
			if qty <= 0 && donGia > 0 {
				qty = sv.TotalAmount / donGia
			}
			if qty <= 0 {
				qty = 1
			}
			unitPrice := donGia
			if unitPrice <= 0 && qty > 0 {
				unitPrice = sv.TotalAmount / qty
			}
			itmTotal := sv.TotalAmount
			if itmTotal <= 0 {
				itmTotal = qty * unitPrice
			}
			prodCode := "SP-DA-01"
			switch {
			case strings.Contains(matHang, "1x2"):
				prodCode = "SP-DA-1X2"
			case strings.Contains(matHang, "Base") || strings.Contains(matHang, "BASE"):
				prodCode = "SP-DA-BASE"
			case strings.Contains(matHang, "Cát") || strings.Contains(matHang, "VSI"):
				prodCode = "SP-CAT-VSI"
			case strings.Contains(matHang, "2x4"):
				prodCode = "SP-DA-2X4"
			case strings.Contains(matHang, "4x6"):
				prodCode = "SP-DA-4X6"
			case strings.Contains(matHang, "Mi Bụi") || strings.Contains(matHang, "MIBUI"):
				prodCode = "SP-DA-MIBUI"
			}
			storageLoc := sv.WarehouseLoc
			if storageLoc == "" {
				storageLoc = tramCan
			}
			var prodVatRate float64 = 10
			_ = database.Pool.QueryRow(ctx, `SELECT COALESCE(vat_rate, 10) FROM inventory_products WHERE code = $1 OR name ILIKE $2 LIMIT 1`, prodCode, "%"+matHang+"%").Scan(&prodVatRate)
			item := models.SalesItem{
				VoucherID:   sv.ID,
				VoucherCode: sv.Code,
				ProductCode: prodCode,
				ProductName: matHang,
				Unit:        "tấn",
				Density:     1.5,
				UnitPrice:   unitPrice,
				Quantity:    qty,
				TotalAmount: itmTotal,
				VatRate:     prodVatRate,
				WeightTon:   qty,
				Standard:    quyCach,
				StorageLoc:  storageLoc,
				Notes:       "Xuất kho tự động từ phiếu cân " + ticketID,
			}
			var newItemID int
			_ = database.Pool.QueryRow(ctx, `
				INSERT INTO sales_voucher_items (voucher_id, voucher_code, product_code, product_name, unit, density, unit_price, quantity, total_amount, weight_ton, standard, storage_loc, notes, vat_rate)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
				RETURNING id
			`, sv.ID, sv.Code, item.ProductCode, item.ProductName, item.Unit, item.Density, item.UnitPrice, item.Quantity, item.TotalAmount, item.WeightTon, item.Standard, item.StorageLoc, item.Notes, item.VatRate).Scan(&newItemID)
			item.ID = newItemID
			sv.Items = append(sv.Items, item)
		}
	}
	if sv.Items == nil {
		sv.Items = []models.SalesItem{}
	}

	JSON(w, sv)
}

func (h *TradeVoucherHandler) CreateSales(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var sv models.SalesVoucher
	if err := decodeJSON(r, &sv); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	if sv.Code == "" {
		sv.Code = "PB-" + time.Now().Format("20060102-1504")
	}
	if sv.Date == "" {
		sv.Date = time.Now().Format("2006-01-02")
	}
	if sv.Status == "" {
		sv.Status = "completed"
	}

	err := database.Pool.QueryRow(ctx, `
		INSERT INTO sales_vouchers (code, customer_code, customer_name, date, warehouse_loc, license_plate, ticket_code, total_amount, vat_amount, grand_total, paid_amount, debt_amount, payment_status, status, created_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at
	`, sv.Code, sv.CustomerCode, sv.CustomerName, sv.Date, sv.WarehouseLoc, sv.LicensePlate, sv.TicketCode, sv.TotalAmount, sv.VatAmount, sv.GrandTotal, sv.PaidAmount, sv.DebtAmount, sv.PaymentStatus, sv.Status, sv.CreatedBy, sv.Notes).Scan(&sv.ID, &sv.CreatedAt, &sv.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể tạo phiếu bán: "+err.Error())
		return
	}

	for _, itm := range sv.Items {
		vatRate := itm.VatRate
		if vatRate == 0 {
			vatRate = 10
		}
		database.Pool.Exec(ctx, `
			INSERT INTO sales_voucher_items (voucher_id, voucher_code, product_code, product_name, unit, density, unit_price, quantity, total_amount, weight_ton, standard, storage_loc, notes, vat_rate)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		`, sv.ID, sv.Code, itm.ProductCode, itm.ProductName, itm.Unit, itm.Density, itm.UnitPrice, itm.Quantity, itm.TotalAmount, itm.WeightTon, itm.Standard, itm.StorageLoc, itm.Notes, vatRate)
	}

	JSON(w, sv)
}

func (h *TradeVoucherHandler) UpdateSales(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	var sv models.SalesVoucher
	if err := decodeJSON(r, &sv); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	numID, _ := strconv.Atoi(idOrCode)
	res, err := database.Pool.Exec(ctx, `
		UPDATE sales_vouchers
		SET customer_code = $1, customer_name = $2, date = $3, warehouse_loc = $4,
		    license_plate = $5, ticket_code = $6, total_amount = $7, vat_amount = $8,
		    grand_total = $9, paid_amount = $10, debt_amount = $11, payment_status = $12,
		    status = $13, notes = $14, updated_at = NOW()
		WHERE code = $15 OR id = $16
	`, sv.CustomerCode, sv.CustomerName, sv.Date, sv.WarehouseLoc, sv.LicensePlate, sv.TicketCode, sv.TotalAmount, sv.VatAmount, sv.GrandTotal, sv.PaidAmount, sv.DebtAmount, sv.PaymentStatus, sv.Status, sv.Notes, idOrCode, numID)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể cập nhật phiếu bán: "+err.Error())
		return
	}
	if res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu bán")
		return
	}

	// Update items if supplied
	if len(sv.Items) > 0 {
		database.Pool.Exec(ctx, `DELETE FROM sales_voucher_items WHERE voucher_code = $1 OR voucher_id = $2`, idOrCode, numID)
		for _, itm := range sv.Items {
			vatRate := itm.VatRate
			if vatRate == 0 {
				vatRate = 10
			}
			database.Pool.Exec(ctx, `
				INSERT INTO sales_voucher_items (voucher_id, voucher_code, product_code, product_name, unit, density, unit_price, quantity, total_amount, weight_ton, standard, storage_loc, notes, vat_rate)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			`, sv.ID, sv.Code, itm.ProductCode, itm.ProductName, itm.Unit, itm.Density, itm.UnitPrice, itm.Quantity, itm.TotalAmount, itm.WeightTon, itm.Standard, itm.StorageLoc, itm.Notes, vatRate)
		}
	}

	JSON(w, map[string]interface{}{"status": "success", "message": "Cập nhật phiếu bán thành công"})
}

func (h *TradeVoucherHandler) DeleteSales(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	database.Pool.Exec(ctx, `DELETE FROM sales_voucher_items WHERE voucher_code = $1 OR voucher_id = $2`, idOrCode, numID)
	res, err := database.Pool.Exec(ctx, `DELETE FROM sales_vouchers WHERE code = $1 OR id = $2`, idOrCode, numID)
	if err != nil || res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu bán")
		return
	}
	JSON(w, map[string]interface{}{"status": "success", "message": "Đã xóa phiếu bán"})
}

// ==========================================
// 3. RETURN VOUCHERS (PHIẾU TRẢ HÀNG)
// ==========================================

func (h *TradeVoucherHandler) ListReturns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	search := r.URL.Query().Get("search")
	returnType := r.URL.Query().Get("type")

	query := `
		SELECT id, code, COALESCE(return_type, 'sales_return'), COALESCE(partner_code, ''),
		       COALESCE(partner_name, ''), COALESCE(ref_voucher_code, ''), COALESCE(date, ''),
		       COALESCE(warehouse_loc, ''), COALESCE(reason, ''), COALESCE(total_amount, 0),
		       COALESCE(status, 'approved'), COALESCE(created_by, ''), COALESCE(notes, ''),
		       created_at, updated_at
		FROM return_vouchers
		WHERE 1=1
	`
	var args []interface{}
	idx := 1
	if search != "" {
		query += ` AND (code ILIKE $` + strconv.Itoa(idx) + ` OR partner_name ILIKE $` + strconv.Itoa(idx) + ` OR ref_voucher_code ILIKE $` + strconv.Itoa(idx) + ` OR reason ILIKE $` + strconv.Itoa(idx) + `)`
		args = append(args, "%"+search+"%")
		idx++
	}
	if returnType != "" && returnType != "all" {
		query += ` AND return_type = $` + strconv.Itoa(idx)
		args = append(args, returnType)
		idx++
	}
	query += ` ORDER BY id DESC`

	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		JSON(w, []interface{}{})
		return
	}
	defer rows.Close()

	var list []models.ReturnVoucher
	for rows.Next() {
		var rv models.ReturnVoucher
		if err := rows.Scan(&rv.ID, &rv.Code, &rv.ReturnType, &rv.PartnerCode, &rv.PartnerName, &rv.RefVoucherCode, &rv.Date, &rv.WarehouseLoc, &rv.Reason, &rv.TotalAmount, &rv.Status, &rv.CreatedBy, &rv.Notes, &rv.CreatedAt, &rv.UpdatedAt); err == nil {
			list = append(list, rv)
		}
	}
	if list == nil {
		list = []models.ReturnVoucher{}
	}
	JSON(w, map[string]interface{}{"data": list, "total": len(list)})
}

func (h *TradeVoucherHandler) GetReturn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	var rv models.ReturnVoucher
	err := database.Pool.QueryRow(ctx, `
		SELECT id, code, COALESCE(return_type, 'sales_return'), COALESCE(partner_code, ''),
		       COALESCE(partner_name, ''), COALESCE(ref_voucher_code, ''), COALESCE(date, ''),
		       COALESCE(warehouse_loc, ''), COALESCE(reason, ''), COALESCE(total_amount, 0),
		       COALESCE(status, 'approved'), COALESCE(created_by, ''), COALESCE(notes, ''),
		       created_at, updated_at
		FROM return_vouchers
		WHERE code = $1 OR id = $2
	`, idOrCode, numID).Scan(&rv.ID, &rv.Code, &rv.ReturnType, &rv.PartnerCode, &rv.PartnerName, &rv.RefVoucherCode, &rv.Date, &rv.WarehouseLoc, &rv.Reason, &rv.TotalAmount, &rv.Status, &rv.CreatedBy, &rv.Notes, &rv.CreatedAt, &rv.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu trả "+idOrCode)
		return
	}

	// Fetch items
	itemRows, itemErr := database.Pool.Query(ctx, `
		SELECT id, COALESCE(voucher_id, 0), COALESCE(voucher_code, ''), COALESCE(product_code, ''), COALESCE(product_name, ''),
		       COALESCE(unit, 'm³'), COALESCE(unit_price, 0), COALESCE(quantity, 0),
		       COALESCE(total_amount, 0), COALESCE(notes, ''), COALESCE(vat_rate, 10)
		FROM return_voucher_items
		WHERE (voucher_id > 0 AND voucher_id = $1) OR voucher_code = $2
		ORDER BY id ASC
	`, rv.ID, rv.Code)

	if itemErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var itm models.ReturnItem
			if err := itemRows.Scan(&itm.ID, &itm.VoucherID, &itm.VoucherCode, &itm.ProductCode, &itm.ProductName, &itm.Unit, &itm.UnitPrice, &itm.Quantity, &itm.TotalAmount, &itm.Notes, &itm.VatRate); err == nil {
				rv.Items = append(rv.Items, itm)
			}
		}
	}
	if rv.Items == nil {
		rv.Items = []models.ReturnItem{}
	}

	JSON(w, rv)
}

func (h *TradeVoucherHandler) CreateReturn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var rv models.ReturnVoucher
	if err := decodeJSON(r, &rv); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	if rv.Code == "" {
		rv.Code = "PT-" + time.Now().Format("20060102-1504")
	}
	if rv.Date == "" {
		rv.Date = time.Now().Format("2006-01-02")
	}
	if rv.Status == "" {
		rv.Status = "approved"
	}

	err := database.Pool.QueryRow(ctx, `
		INSERT INTO return_vouchers (code, return_type, partner_code, partner_name, ref_voucher_code, date, warehouse_loc, reason, total_amount, status, created_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`, rv.Code, rv.ReturnType, rv.PartnerCode, rv.PartnerName, rv.RefVoucherCode, rv.Date, rv.WarehouseLoc, rv.Reason, rv.TotalAmount, rv.Status, rv.CreatedBy, rv.Notes).Scan(&rv.ID, &rv.CreatedAt, &rv.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể tạo phiếu trả: "+err.Error())
		return
	}

	for _, itm := range rv.Items {
		vatRate := itm.VatRate
		if vatRate == 0 {
			vatRate = 10
		}
		database.Pool.Exec(ctx, `
			INSERT INTO return_voucher_items (voucher_id, voucher_code, product_code, product_name, unit, unit_price, quantity, total_amount, notes, vat_rate)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, rv.ID, rv.Code, itm.ProductCode, itm.ProductName, itm.Unit, itm.UnitPrice, itm.Quantity, itm.TotalAmount, itm.Notes, vatRate)
	}

	JSON(w, rv)
}

func (h *TradeVoucherHandler) UpdateReturn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	var rv models.ReturnVoucher
	if err := decodeJSON(r, &rv); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	numID, _ := strconv.Atoi(idOrCode)
	res, err := database.Pool.Exec(ctx, `
		UPDATE return_vouchers
		SET return_type = $1, partner_code = $2, partner_name = $3, ref_voucher_code = $4,
		    date = $5, warehouse_loc = $6, reason = $7, total_amount = $8, status = $9,
		    notes = $10, updated_at = NOW()
		WHERE code = $11 OR id = $12
	`, rv.ReturnType, rv.PartnerCode, rv.PartnerName, rv.RefVoucherCode, rv.Date, rv.WarehouseLoc, rv.Reason, rv.TotalAmount, rv.Status, rv.Notes, idOrCode, numID)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể cập nhật phiếu trả: "+err.Error())
		return
	}
	if res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu trả")
		return
	}

	// Update items if supplied
	if len(rv.Items) > 0 {
		database.Pool.Exec(ctx, `DELETE FROM return_voucher_items WHERE voucher_code = $1 OR voucher_id = $2`, idOrCode, numID)
		for _, itm := range rv.Items {
			database.Pool.Exec(ctx, `
				INSERT INTO return_voucher_items (voucher_id, voucher_code, product_code, product_name, unit, unit_price, quantity, total_amount, notes)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, rv.ID, rv.Code, itm.ProductCode, itm.ProductName, itm.Unit, itm.UnitPrice, itm.Quantity, itm.TotalAmount, itm.Notes)
		}
	}

	JSON(w, map[string]interface{}{"status": "success", "message": "Cập nhật phiếu trả thành công"})
}

func (h *TradeVoucherHandler) DeleteReturn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	database.Pool.Exec(ctx, `DELETE FROM return_voucher_items WHERE voucher_code = $1 OR voucher_id = $2`, idOrCode, numID)
	res, err := database.Pool.Exec(ctx, `DELETE FROM return_vouchers WHERE code = $1 OR id = $2`, idOrCode, numID)
	if err != nil || res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu trả")
		return
	}
	JSON(w, map[string]interface{}{"status": "success", "message": "Đã xóa phiếu trả"})
}
