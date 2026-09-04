package handlers

import (
	"net/http"
	"strconv"
	"time"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/models"
)

type PaymentVoucherHandler struct{}

func NewPaymentVoucherHandler() *PaymentVoucherHandler {
	return &PaymentVoucherHandler{}
}

// ==========================================
// 1. RECEIPTS (PHIẾU THU)
// ==========================================

func (h *PaymentVoucherHandler) ListReceipts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	query := `
		SELECT id, code, COALESCE(partner_code, ''), COALESCE(partner_name, ''),
		       COALESCE(partner_type, 'customer'), COALESCE(payer_name, ''),
		       COALESCE(reason, ''), COALESCE(ref_code, ''), COALESCE(amount, 0),
		       COALESCE(amount_in_words, ''), COALESCE(payment_method, 'bank_transfer'),
		       COALESCE(fund_account, ''), COALESCE(status, 'posted'), COALESCE(date, ''),
		       COALESCE(created_by, ''), COALESCE(notes, ''), created_at, updated_at
		FROM receipt_vouchers
		WHERE 1=1
	`
	var args []interface{}
	idx := 1
	if search != "" {
		query += ` AND (code ILIKE $` + strconv.Itoa(idx) + ` OR partner_name ILIKE $` + strconv.Itoa(idx) + ` OR payer_name ILIKE $` + strconv.Itoa(idx) + ` OR ref_code ILIKE $` + strconv.Itoa(idx) + `)`
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

	var list []models.ReceiptVoucher
	for rows.Next() {
		var rc models.ReceiptVoucher
		if err := rows.Scan(&rc.ID, &rc.Code, &rc.PartnerCode, &rc.PartnerName, &rc.PartnerType, &rc.PayerName, &rc.Reason, &rc.RefCode, &rc.Amount, &rc.AmountInWords, &rc.PaymentMethod, &rc.FundAccount, &rc.Status, &rc.Date, &rc.CreatedBy, &rc.Notes, &rc.CreatedAt, &rc.UpdatedAt); err == nil {
			list = append(list, rc)
		}
	}
	if list == nil {
		list = []models.ReceiptVoucher{}
	}
	JSON(w, map[string]interface{}{"data": list, "total": len(list)})
}

func (h *PaymentVoucherHandler) GetReceipt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	var rc models.ReceiptVoucher
	err := database.Pool.QueryRow(ctx, `
		SELECT id, code, COALESCE(partner_code, ''), COALESCE(partner_name, ''),
		       COALESCE(partner_type, 'customer'), COALESCE(payer_name, ''),
		       COALESCE(reason, ''), COALESCE(ref_code, ''), COALESCE(amount, 0),
		       COALESCE(amount_in_words, ''), COALESCE(payment_method, 'bank_transfer'),
		       COALESCE(fund_account, ''), COALESCE(status, 'posted'), COALESCE(date, ''),
		       COALESCE(created_by, ''), COALESCE(notes, ''), created_at, updated_at
		FROM receipt_vouchers
		WHERE code = $1 OR id = $2
	`, idOrCode, numID).Scan(&rc.ID, &rc.Code, &rc.PartnerCode, &rc.PartnerName, &rc.PartnerType, &rc.PayerName, &rc.Reason, &rc.RefCode, &rc.Amount, &rc.AmountInWords, &rc.PaymentMethod, &rc.FundAccount, &rc.Status, &rc.Date, &rc.CreatedBy, &rc.Notes, &rc.CreatedAt, &rc.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu thu "+idOrCode)
		return
	}
	JSON(w, rc)
}

func (h *PaymentVoucherHandler) CreateReceipt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var rc models.ReceiptVoucher
	if err := decodeJSON(r, &rc); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	if rc.Code == "" {
		rc.Code = "PT-THU-" + time.Now().Format("20060102-1504")
	}
	if rc.Date == "" {
		rc.Date = time.Now().Format("2006-01-02")
	}
	if rc.Status == "" {
		rc.Status = "posted"
	}

	err := database.Pool.QueryRow(ctx, `
		INSERT INTO receipt_vouchers (code, partner_code, partner_name, partner_type, payer_name, reason, ref_code, amount, amount_in_words, payment_method, fund_account, status, date, created_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at
	`, rc.Code, rc.PartnerCode, rc.PartnerName, rc.PartnerType, rc.PayerName, rc.Reason, rc.RefCode, rc.Amount, rc.AmountInWords, rc.PaymentMethod, rc.FundAccount, rc.Status, rc.Date, rc.CreatedBy, rc.Notes).Scan(&rc.ID, &rc.CreatedAt, &rc.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể tạo phiếu thu: "+err.Error())
		return
	}
	JSON(w, rc)
}

func (h *PaymentVoucherHandler) UpdateReceipt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	var rc models.ReceiptVoucher
	if err := decodeJSON(r, &rc); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	numID, _ := strconv.Atoi(idOrCode)
	res, err := database.Pool.Exec(ctx, `
		UPDATE receipt_vouchers
		SET partner_code = $1, partner_name = $2, partner_type = $3, payer_name = $4,
		    reason = $5, ref_code = $6, amount = $7, amount_in_words = $8,
		    payment_method = $9, fund_account = $10, status = $11, date = $12,
		    notes = $13, updated_at = NOW()
		WHERE code = $14 OR id = $15
	`, rc.PartnerCode, rc.PartnerName, rc.PartnerType, rc.PayerName, rc.Reason, rc.RefCode, rc.Amount, rc.AmountInWords, rc.PaymentMethod, rc.FundAccount, rc.Status, rc.Date, rc.Notes, idOrCode, numID)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể cập nhật phiếu thu: "+err.Error())
		return
	}
	if res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu thu")
		return
	}
	JSON(w, map[string]interface{}{"status": "success", "message": "Cập nhật phiếu thu thành công"})
}

func (h *PaymentVoucherHandler) DeleteReceipt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	res, err := database.Pool.Exec(ctx, `DELETE FROM receipt_vouchers WHERE code = $1 OR id = $2`, idOrCode, numID)
	if err != nil || res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu thu")
		return
	}
	JSON(w, map[string]interface{}{"status": "success", "message": "Đã xóa phiếu thu"})
}

// ==========================================
// 2. PAYMENT VOUCHERS (PHIẾU CHI)
// ==========================================

func (h *PaymentVoucherHandler) ListPaymentVouchers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	query := `
		SELECT id, code, COALESCE(partner_code, ''), COALESCE(partner_name, ''),
		       COALESCE(partner_type, 'supplier'), COALESCE(receiver_name, ''),
		       COALESCE(reason, ''), COALESCE(ref_code, ''), COALESCE(amount, 0),
		       COALESCE(amount_in_words, ''), COALESCE(payment_method, 'bank_transfer'),
		       COALESCE(fund_account, ''), COALESCE(status, 'posted'), COALESCE(date, ''),
		       COALESCE(created_by, ''), COALESCE(notes, ''), created_at, updated_at
		FROM payment_vouchers
		WHERE 1=1
	`
	var args []interface{}
	idx := 1
	if search != "" {
		query += ` AND (code ILIKE $` + strconv.Itoa(idx) + ` OR partner_name ILIKE $` + strconv.Itoa(idx) + ` OR receiver_name ILIKE $` + strconv.Itoa(idx) + ` OR ref_code ILIKE $` + strconv.Itoa(idx) + `)`
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

	var list []models.PaymentVoucher
	for rows.Next() {
		var pv models.PaymentVoucher
		if err := rows.Scan(&pv.ID, &pv.Code, &pv.PartnerCode, &pv.PartnerName, &pv.PartnerType, &pv.ReceiverName, &pv.Reason, &pv.RefCode, &pv.Amount, &pv.AmountInWords, &pv.PaymentMethod, &pv.FundAccount, &pv.Status, &pv.Date, &pv.CreatedBy, &pv.Notes, &pv.CreatedAt, &pv.UpdatedAt); err == nil {
			list = append(list, pv)
		}
	}
	if list == nil {
		list = []models.PaymentVoucher{}
	}
	JSON(w, map[string]interface{}{"data": list, "total": len(list)})
}

func (h *PaymentVoucherHandler) GetPaymentVoucher(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	var pv models.PaymentVoucher
	err := database.Pool.QueryRow(ctx, `
		SELECT id, code, COALESCE(partner_code, ''), COALESCE(partner_name, ''),
		       COALESCE(partner_type, 'supplier'), COALESCE(receiver_name, ''),
		       COALESCE(reason, ''), COALESCE(ref_code, ''), COALESCE(amount, 0),
		       COALESCE(amount_in_words, ''), COALESCE(payment_method, 'bank_transfer'),
		       COALESCE(fund_account, ''), COALESCE(status, 'posted'), COALESCE(date, ''),
		       COALESCE(created_by, ''), COALESCE(notes, ''), created_at, updated_at
		FROM payment_vouchers
		WHERE code = $1 OR id = $2
	`, idOrCode, numID).Scan(&pv.ID, &pv.Code, &pv.PartnerCode, &pv.PartnerName, &pv.PartnerType, &pv.ReceiverName, &pv.Reason, &pv.RefCode, &pv.Amount, &pv.AmountInWords, &pv.PaymentMethod, &pv.FundAccount, &pv.Status, &pv.Date, &pv.CreatedBy, &pv.Notes, &pv.CreatedAt, &pv.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu chi "+idOrCode)
		return
	}
	JSON(w, pv)
}

func (h *PaymentVoucherHandler) CreatePaymentVoucher(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var pv models.PaymentVoucher
	if err := decodeJSON(r, &pv); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	if pv.Code == "" {
		pv.Code = "PC-CHI-" + time.Now().Format("20060102-1504")
	}
	if pv.Date == "" {
		pv.Date = time.Now().Format("2006-01-02")
	}
	if pv.Status == "" {
		pv.Status = "posted"
	}

	err := database.Pool.QueryRow(ctx, `
		INSERT INTO payment_vouchers (code, partner_code, partner_name, partner_type, receiver_name, reason, ref_code, amount, amount_in_words, payment_method, fund_account, status, date, created_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at
	`, pv.Code, pv.PartnerCode, pv.PartnerName, pv.PartnerType, pv.ReceiverName, pv.Reason, pv.RefCode, pv.Amount, pv.AmountInWords, pv.PaymentMethod, pv.FundAccount, pv.Status, pv.Date, pv.CreatedBy, pv.Notes).Scan(&pv.ID, &pv.CreatedAt, &pv.UpdatedAt)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể tạo phiếu chi: "+err.Error())
		return
	}
	JSON(w, pv)
}

func (h *PaymentVoucherHandler) UpdatePaymentVoucher(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	var pv models.PaymentVoucher
	if err := decodeJSON(r, &pv); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	numID, _ := strconv.Atoi(idOrCode)
	res, err := database.Pool.Exec(ctx, `
		UPDATE payment_vouchers
		SET partner_code = $1, partner_name = $2, partner_type = $3, receiver_name = $4,
		    reason = $5, ref_code = $6, amount = $7, amount_in_words = $8,
		    payment_method = $9, fund_account = $10, status = $11, date = $12,
		    notes = $13, updated_at = NOW()
		WHERE code = $14 OR id = $15
	`, pv.PartnerCode, pv.PartnerName, pv.PartnerType, pv.ReceiverName, pv.Reason, pv.RefCode, pv.Amount, pv.AmountInWords, pv.PaymentMethod, pv.FundAccount, pv.Status, pv.Date, pv.Notes, idOrCode, numID)

	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Không thể cập nhật phiếu chi: "+err.Error())
		return
	}
	if res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu chi")
		return
	}
	JSON(w, map[string]interface{}{"status": "success", "message": "Cập nhật phiếu chi thành công"})
}

func (h *PaymentVoucherHandler) DeletePaymentVoucher(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idOrCode := r.PathValue("id")
	numID, _ := strconv.Atoi(idOrCode)

	res, err := database.Pool.Exec(ctx, `DELETE FROM payment_vouchers WHERE code = $1 OR id = $2`, idOrCode, numID)
	if err != nil || res.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy phiếu chi")
		return
	}
	JSON(w, map[string]interface{}{"status": "success", "message": "Đã xóa phiếu chi"})
}
