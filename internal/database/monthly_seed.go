package database

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
)

type quarrySeedSpec struct {
	Code     string
	Name     string
	Location string
	Prefix   string
	Scale    string
	Gate     string
	Ratio    float64
}

var quarrySpecs = []quarrySeedSpec{
	{
		Code:     "MO-PT-01",
		Name:     "Công ty Cổ phần Mỏ Đá TTC - Chi nhánh Phú Thọ",
		Location: "Thanh Ba, Phú Thọ",
		Prefix:   "PT",
		Scale:    "Trạm Cân 01 (100 Tấn)",
		Gate:     "Cổng Cân Số 1",
		Ratio:    1.0,
	},
	{
		Code:     "MO-TU-02",
		Name:     "Công ty Cổ phần Mỏ Đá TTC - Chi nhánh Bình Dương (Tân Uyên)",
		Location: "Tân Uyên, Bình Dương",
		Prefix:   "TU",
		Scale:    "Trạm Cân Tân Uyên 100T",
		Gate:     "Cổng Cân Số 1",
		Ratio:    0.85,
	},
	{
		Code:     "MO-HN-03",
		Name:     "Công ty Cổ phần Mỏ Đá TTC - Chi nhánh Hà Nam (Kiện Khê)",
		Location: "Kiện Khê, Thanh Liêm, Hà Nam",
		Prefix:   "HN",
		Scale:    "Trạm Cân Kiện Khê 120T",
		Gate:     "Cổng Cân Số 2",
		Ratio:    0.90,
	},
	{
		Code:     "MO-BP-04",
		Name:     "Công ty Cổ phần Mỏ Đá TTC - Chi nhánh Bình Phước",
		Location: "Chơn Thành, Bình Phước",
		Prefix:   "BP",
		Scale:    "Trạm Cân Chơn Thành 100T",
		Gate:     "Cổng Cân Số 1",
		Ratio:    0.75,
	},
}

type materialSpec struct {
	Name     string
	Code     string
	Unit     string
	Standard string
	Price    float64
}

var materialCatalog = []materialSpec{
	{"Đá 1x2 Bê tông", "SP-DA-1X2", "tấn", "TCVN 7570:2006", 240000},
	{"Đá Base Cấp Phối Dmax25", "SP-DA-BASE", "tấn", "TCVN 8859:2011", 180000},
	{"Cát Nghiền Nhân Tạo VSI", "SP-CAT-VSI", "tấn", "TCVN 9205:2012", 260000},
	{"Đá 2x4 Xây Dựng", "SP-DA-2X4", "tấn", "TCVN 7570:2006", 235000},
	{"Đá 4x6 Kè Móng", "SP-DA-4X6", "tấn", "TCVN 7570:2006", 220000},
	{"Đá Mi Bụi Đắp Nền", "SP-DA-MIBUI", "tấn", "TCVN 8859:2011", 150000},
}

type customerSpec struct {
	Code string
	Name string
}

var customerCatalog = []customerSpec{
	{"CUST-DEOCA", "Tập đoàn Đèo Cả (Dự án Cao tốc Bắc - Nam)"},
	{"CUST-BQP319", "Tổng Công ty 319 - Bộ Quốc Phòng"},
	{"CUST-TRUONGSON", "Tổng Công Ty XD Trường Sơn"},
	{"CUST-BECAMEX", "Công ty CP Bê Tông Becamex"},
	{"CUST-VICEM", "Công ty CP Xi Măng Vicem Bút Sơn"},
	{"CUST-CC1", "Tổng Công Ty Xây Dựng Số 1 (CC1)"},
	{"CUST-VINACONEX9", "Công ty CP Xây Dựng Vinaconex 9"},
	{"CUST-VIETTRI", "Công ty TNHH Bê Tông Việt Trì"},
}

type fleetSpec struct {
	Plate      string
	Truck      string
	Driver     string
	Phone      string
	Tare       float64
	MaxPayload float64
	BaseNet    float64
	Axles      string
}

var fleetCatalog = []fleetSpec{
	{"19H-056.22", "Xe ben Chenglong 4 chân (8x4)", "Nguyễn Văn Toàn", "0912.888.999", 14.80, 17.90, 17.20, "4 chân (8x4)"},
	{"88H-042.27", "Xe ben HOWO Sinotruk 4 chân (8x4)", "Trần Đình Khang", "0904.112.334", 15.42, 17.50, 17.10, "4 chân (8x4)"},
	{"19C-128.45", "Xe ben Howo 371HP 4 chân (8x4)", "Lê Hữu Thắng", "0983.234.567", 15.20, 17.80, 17.40, "4 chân (8x4)"},
	{"29C-781.90", "Xe ben Shacman 4 chân X3000 (8x4)", "Tô Quốc Huy", "0977.890.123", 15.10, 17.20, 16.90, "4 chân (8x4)"},
	{"29C-345.67", "Xe ben 4 chân Sinotruk (8x4)", "Trần Văn Kiên", "0988.341.992", 15.60, 17.40, 17.00, "4 chân (8x4)"},
	{"36C-789.12", "Xe ben Howo 8x4 371HP", "Lê Văn Cường", "0972.113.445", 15.80, 17.20, 16.80, "4 chân (8x4)"},
	{"19C-098.76", "Xe tải ben 15 tấn 3 chân Hino 500 FL", "Vũ Quốc Đạt", "0978.556.223", 10.80, 14.50, 14.20, "3 chân (6x4)"},
	{"90C-054.67", "Xe ben 15 tấn 3 chân Howo (6x4)", "Nguyễn Hoàng Long", "0981.332.654", 11.20, 15.00, 14.60, "3 chân (6x4)"},
	{"24C-222.11", "Xe tải 12 tấn 2 dí 1 cầu Chenglong", "Trương Văn Nam", "0938.667.129", 9.80, 12.00, 11.80, "3 trục (2 dí)"},
	{"14B-567.89", "Xe tải ben 20 tấn 5 chân Howo A7", "Bùi Đức Thịnh", "0972.441.902", 16.80, 20.00, 19.50, "5 chân (10x4)"},
	{"29H-882.19", "Xe bồn trộn bê tông 12m3 Howo", "Hoàng Minh Đức", "0977.456.123", 16.10, 16.00, 15.80, "3 trục (6x4)"},
	{"61C-445.89", "Xe ben 4 chân Shacman 380HP", "Lê Quốc Bảo", "0903.551.234", 15.30, 17.50, 17.10, "4 chân (8x4)"},
}

// SeedMonthlyQuarryData seeds comprehensive ticket, trip, voucher, cost, fuel, and attendance data
// for all 4 quarries across August and September 2026 with realistic, balanced financial ratios (Gross Margin ~36.5%).
func SeedMonthlyQuarryData() {
	if Pool == nil {
		return
	}
	ctx := context.Background()
	ict := time.FixedZone("ICT", 7*60*60)
	now := time.Now().In(ict)

	fmt.Println("🌱 Seeding Comprehensive Balanced Quarry Data (Aug & Sep 2026, 4 Quarries)...")

	// 1. Thoroughly clean old 2026 data to ensure 100% coherence and prevent double counting
	_, _ = Pool.Exec(ctx, "DELETE FROM tickets WHERE id LIKE 'TK-%-2026%'")
	_, _ = Pool.Exec(ctx, "DELETE FROM vehicle_trips WHERE camera_id = 'CAM-GATE-01' AND check_in_time >= '2026-01-01'")
	_, _ = Pool.Exec(ctx, "DELETE FROM sales_vouchers WHERE code LIKE 'PB-TK-%'")
	_, _ = Pool.Exec(ctx, "DELETE FROM production_costs WHERE created_at >= '2026-01-01'")
	_, _ = Pool.Exec(ctx, "DELETE FROM equipment_fuel_logs WHERE created_at >= '2026-01-01'")
	_, _ = Pool.Exec(ctx, "DELETE FROM hr_attendances WHERE created_at >= '2026-01-01'")
	_, _ = Pool.Exec(ctx, "DELETE FROM inventory_inbound WHERE created_at >= '2026-01-01'")
	_, _ = Pool.Exec(ctx, "DELETE FROM inventory_outbound WHERE created_at >= '2026-01-01'")
	_, _ = Pool.Exec(ctx, "DELETE FROM alerts WHERE created_at >= '2026-01-01'")
	_, _ = Pool.Exec(ctx, "UPDATE tickets SET kl_hang = '30.00', kl_tinh_tien = '30.00', thanh_tien = '7.200.000 đ' WHERE id = 'NA181026-2031'")

	rnd := rand.New(rand.NewSource(20260910))
	batch := &pgx.Batch{}

	// Helper for inserting ticket + vehicle_trip + sales_voucher
	insertTicket := func(b *pgx.Batch, q quarrySeedSpec, tDate time.Time, seq int) (float64, float64) {
		ticketID := fmt.Sprintf("TK-%s-%s-%03d", q.Prefix, tDate.Format("20060102"), seq)
		mat := materialCatalog[(seq+int(tDate.Month())+int(tDate.Day()))%len(materialCatalog)]
		cust := customerCatalog[(seq+int(tDate.Day()))%len(customerCatalog)]
		flt := fleetCatalog[(seq+int(tDate.Month()))%len(fleetCatalog)]

		// Tải trọng hàng thực tế chuẩn 10 - 20 tấn theo thông số xe Trọng Tấn
		netWeight := flt.BaseNet + float64(rnd.Intn(90))/100.0 - 0.40 // Dao động nhẹ quanh tải trọng chuẩn
		tareWeight := flt.Tare
		grossWeight := tareWeight + netWeight
		unitPrice := mat.Price
		totalPrice := netWeight * unitPrice

		time1Str := tDate.Add(-25 * time.Minute).Format("15:04")
		time2Str := tDate.Format("15:04")
		dateStr := tDate.Format("02/01/2006")
		doCode := fmt.Sprintf("DO-%s-%03d", q.Prefix, 100+seq)
		formattedPrice := fmt.Sprintf("%.0f đ", totalPrice)

		b.Queue(`
			INSERT INTO tickets (
				id, stt, ben_ban, ben_mua, bien_so, loai_xe, lai_xe, sdt_lai_xe, rfid,
				loai, stage, stage_label, can_l1, kl1, can_l2, kl2, kl_hang, kl_tap_chat,
				kl_tinh_tien, don_gia, thanh_tien, time1, time2, date, nguoi_can1, nguoi_can2,
				mat_hang, quy_cach, do_code, tram_can, cong_can, ghi_chu, hoa_don_so,
				quarry_code, created_at, updated_at, cameras, chatter
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9,
				$10, $11, $12, $13, $14, $15, $16, $17, $18,
				$19, $20, $21, $22, $23, $24, $25, $26,
				$27, $28, $29, $30, $31, $32, $33,
				$34, $35, $36, $37, $38
			)
			ON CONFLICT (id) DO UPDATE SET
				kl_hang = EXCLUDED.kl_hang,
				don_gia = EXCLUDED.don_gia,
				thanh_tien = EXCLUDED.thanh_tien
		`,
			ticketID, seq, q.Name, cust.Name, flt.Plate, flt.Truck, flt.Driver, flt.Phone, "RFID-"+flt.Plate,
			"Cân bán hàng", "confirmed", "Đã chốt số", tareWeight, fmt.Sprintf("%.2f", tareWeight),
			grossWeight, fmt.Sprintf("%.2f", grossWeight), fmt.Sprintf("%.2f", netWeight), 0.0,
			fmt.Sprintf("%.2f", netWeight), unitPrice, formattedPrice, time1Str, time2Str, dateStr,
			"Nguyễn Văn Dũng", "Lê Văn Cân 02",
			mat.Name, mat.Standard, doCode, q.Scale, q.Gate, "Xuất mỏ đủ tải", "HD-2026-"+ticketID,
			q.Code, tDate, tDate,
			fmt.Sprintf(`{"front":{"camera":"Camera 01 ANPR","time":"%s %s"},"rear":{"camera":"Camera 02 Thùng","time":"%s %s"}}`, dateStr, time1Str, dateStr, time2Str),
			fmt.Sprintf(`[{"author":"Hệ thống Cân AI","content":"Xe %s hoàn tất cân tự động %s tại %s.","time":"%s %s"}]`, flt.Plate, mat.Name, q.Scale, dateStr, time2Str),
		)

		b.Queue(`
			INSERT INTO vehicle_trips (
				vehicle_id, license_plate, driver_name, camera_id, direction,
				check_in_time, check_out_time, trip_number, estimated_quantity, actual_quantity,
				status, created_at, updated_at
			) VALUES (
				$1, $2, $3, 'CAM-GATE-01', 'outbound',
				$4, $5, $6, $7, $8,
				'completed', $9, $10
			)
			ON CONFLICT DO NOTHING
		`,
			flt.Plate, flt.Plate, flt.Driver,
			tDate.Add(-25*time.Minute), tDate, seq,
			netWeight, netWeight,
			tDate, tDate,
		)

		voucherCode := "PB-" + ticketID
		b.Queue(`
			INSERT INTO sales_vouchers (
				code, customer_code, customer_name, date, warehouse_loc, license_plate,
				ticket_code, total_amount, vat_amount, grand_total, paid_amount, debt_amount,
				payment_status, status, created_by, notes, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10, $11, 0,
				'paid', 'completed', 'Nguyễn Văn Dũng', $12, $13, $14
			)
			ON CONFLICT (code) DO UPDATE SET total_amount = EXCLUDED.total_amount
		`,
			voucherCode, cust.Code, cust.Name, tDate.Format("2006-01-02"), q.Location, flt.Plate,
			ticketID, totalPrice, totalPrice*0.1, totalPrice*1.1, totalPrice*1.1,
			"Xuất hàng trạm cân "+q.Scale, tDate, tDate,
		)

		return netWeight, totalPrice
	}

	// Helper to insert production costs strictly calibrated to 63.5% of ticket revenue (Margin ~36.5%)
	insertDailyCost := func(b *pgx.Batch, q quarrySeedSpec, entryDate time.Time, dayRevenue float64) {
		mineArea := fmt.Sprintf("%s (%s)", q.Name, q.Code)
		periodStr := entryDate.Format("2006-01")
		totalCost := dayRevenue * 0.635
		categories := []struct {
			CostType string
			Category string
			Share    float64
			Desc     string
		}{
			{"Sản xuất", "Chi phí sản xuất & nổ mìn", 0.30, "Khoan nổ mìn tầng khai thác và điện cấp liệu trạm nghiền"},
			{"Nhiên liệu", "Nhiên liệu diesel cơ giới", 0.28, "Dầu diesel máy xúc PC450 và xe ben nội bộ"},
			{"Nhân công", "Nhân công & Tiền lương ca", 0.22, "Tiền lương ca vận hành mỏ, trạm cân và cơ giới"},
			{"Khấu hao", "Khấu hao trạm nghiền sàng", 0.10, "Khấu hao thiết bị nghiền sàng và băng tải"},
			{"Vận chuyển", "Vận chuyển nội bộ moong", 0.06, "Cước vận chuyển đá hộc đáy moong lên trạm nghiền"},
			{"Khác", "Bảo hộ ATLĐ & Môi trường", 0.04, "Phun sương dập bụi, quan trắc môi trường và bảo hộ lao động"},
		}

		for _, c := range categories {
			actual := totalCost * c.Share
			norm := actual * 1.02
			b.Queue(`
				INSERT INTO production_costs (
					cost_type, cost_category, norm_value, norm_unit, actual_value, actual_unit,
					period, mine_area, description, created_at
				) VALUES (
					$1, $2, $3, 'VNĐ', $4, 'VNĐ',
					$5, $6, $7, $8
				)
			`, c.CostType, c.Category, norm, actual, periodStr, mineArea, c.Desc, entryDate)
		}
	}

	// Helper for inventory movement (Inbound ~1.02x production, Outbound ~0.98x production)
	insertDailyInventory := func(b *pgx.Batch, q quarrySeedSpec, entryDate time.Time, dayTonnage float64) {
		dateStr := entryDate.Format("02/01/2006")
		inID := fmt.Sprintf("INV-IN-%s-%s", q.Prefix, entryDate.Format("20060102"))
		qtyIn := dayTonnage * 1.02
		b.Queue(`
			INSERT INTO inventory_inbound (code, source, loc, item, qty, quantity, unit, date, status, created_at)
			VALUES ($1, 'Trạm Nghiền Sàng', $2, 'Đá 1x2 Bê tông', $3, $3, 'tấn', $4, 'Hoàn thành', $5)
			ON CONFLICT DO NOTHING
		`, inID, q.Location, qtyIn, dateStr, entryDate)

		outID := fmt.Sprintf("INV-OUT-%s-%s", q.Prefix, entryDate.Format("20060102"))
		qtyOut := dayTonnage * 0.98
		b.Queue(`
			INSERT INTO inventory_outbound (code, customer, dest, item, qty, quantity, unit, date, status, created_at)
			VALUES ($1, 'Khách hàng theo phiếu cân', $2, 'Đá 1x2 Bê tông', $3, $3, 'tấn', $4, 'Hoàn thành', $5)
			ON CONFLICT DO NOTHING
		`, outID, q.Location, qtyOut, dateStr, entryDate)
	}

	// Helper for fuel logs (~2.2 L of diesel per ton of stone)
	insertDailyFuel := func(b *pgx.Batch, q quarrySeedSpec, entryDate time.Time, dayTonnage float64) {
		dateStr := entryDate.Format("02/01/2006")
		totalFuel := dayTonnage * 2.2
		machines := []struct {
			Code, Name, Category string
			Share                float64
			Hours                float64
		}{
			{"MX-01", "Máy Xúc Komatsu PC450", "Máy Xúc", 0.40, 7.5},
			{"DCN-01", "Dây Chuyền Nghiền 01", "Dây Chuyền Nghiền", 0.42, 7.0},
			{"XB-01", "Xe Ben Howo Moong", "Xe Vận Tải Moong", 0.18, 6.5},
		}
		for _, m := range machines {
			id := fmt.Sprintf("FUEL-%s-%s-%s", q.Prefix, entryDate.Format("20060102"), m.Code)
			consumed := totalFuel * m.Share
			issued := consumed * 1.01
			quota := consumed / m.Hours
			b.Queue(`
				INSERT INTO equipment_fuel_logs (
					id, equipment_code, equipment_name, category, operator_name,
					hours_worked_today, total_hours_meter, fuel_quota_liters_per_hour,
					actual_fuel_issued_liters, actual_fuel_consumed_liters, fuel_variance_liters,
					variance_status, location, maintenance_status, tank_capacity_liters,
					current_fuel_liters, last_dispense_at, next_maintenance_hours, quarry_code, created_at
				) VALUES (
					$1, $2, $3, $4, 'Thợ máy vận hành',
					$5, 4500, $6,
					$7, $8, $9,
					'Chuẩn định mức', $10, 'Hoạt động tốt', 400,
					280, $11, 5000, $12, $13
				)
				ON CONFLICT (id) DO UPDATE SET
					actual_fuel_consumed_liters = EXCLUDED.actual_fuel_consumed_liters
			`,
				id, m.Code, m.Name, m.Category,
				m.Hours, quota,
				issued, consumed, issued-consumed,
				q.Location, dateStr+" 17:00", q.Code, entryDate,
			)
		}
	}

	// Helper for attendances (6 workers per quarry)
	insertDailyAttendance := func(b *pgx.Batch, q quarrySeedSpec, entryDate time.Time) {
		dateStr := entryDate.Format("02/01/2006")
		workers := []struct {
			Dept, Role, Name string
		}{
			{"Tổ Vận Hành Trạm Cân", "Nhân viên trạm cân", "Nguyễn Văn Dũng"},
			{"Tổ Vận Hành Trạm Cân", "Thủ kho bãi", "Lê Văn Cân 02"},
			{"Đội Cơ Giới & Vận Tải", "Tài xế xe ben", "Nguyễn Văn Mạnh"},
			{"Xưởng Nghiền Sàng", "Thợ máy nghiền", "Phạm Văn Cường"},
			{"Phòng Kỹ Thuật", "Kỹ sư trắc địa", "Trần Văn Kiên"},
			{"Ban Điều Hành", "Điều phối ca", "Nguyễn Đức Trường"},
		}
		for idx, w := range workers {
			attID := fmt.Sprintf("ATT-%s-%s-%d", q.Prefix, entryDate.Format("20060102"), idx)
			checkInTime := time.Date(entryDate.Year(), entryDate.Month(), entryDate.Day(), 6, 30+idx*5, 0, 0, ict)
			b.Queue(`
				INSERT INTO hr_attendances (
					id, employee_id, employee_name, department, job_position,
					date, check_in_time, check_out_time, contracted_hours, worked_hours,
					status, location, created_at
				) VALUES (
					$1, $2, $3, $4, $5,
					$6, '06:45', '15:15', 8.0, 8.0,
					'Đúng giờ', $7, $8
				)
				ON CONFLICT (id) DO UPDATE SET created_at = EXCLUDED.created_at
			`,
				attID, fmt.Sprintf("EMP-%s-%02d", q.Prefix, idx+1), w.Name,
				w.Dept, w.Role,
				dateStr, q.Location, checkInTime,
			)
		}
	}

	// 2. Loop through August 2026 (Aug 1 to Aug 31)
	for d := 1; d <= 31; d++ {
		for _, q := range quarrySpecs {
			var dayTonnage, dayRevenue float64
			ticketCount := 2
			if q.Ratio >= 0.95 {
				ticketCount = 3
			}

			for i := 1; i <= ticketCount; i++ {
				var hour int
				if i == 1 {
					hour = 7
				} else if i == 2 {
					hour = 9
				} else {
					hour = 14
				}
				min := 15 + ((d*7 + i*13) % 35)
				tDate := time.Date(2026, time.August, d, hour, min, 0, 0, ict)
				tTons, tRev := insertTicket(batch, q, tDate, i)
				dayTonnage += tTons
				dayRevenue += tRev
			}

			costDate := time.Date(2026, time.August, d, 7, 15, 0, 0, ict)
			insertDailyCost(batch, q, costDate, dayRevenue)
			insertDailyInventory(batch, q, costDate, dayTonnage)
			insertDailyFuel(batch, q, costDate, dayTonnage)
			insertDailyAttendance(batch, q, costDate)
		}
	}

	// 3. Loop through September 2026 up to today (now.Day())
	maxDay := now.Day()
	if maxDay < 11 {
		maxDay = 11
	}
	for d := 1; d <= maxDay; d++ {
		for _, q := range quarrySpecs {
			var dayTonnage, dayRevenue float64
			ticketCount := 2
			if q.Ratio >= 0.95 {
				ticketCount = 3
			}
			if d >= 7 && q.Ratio >= 0.85 {
				ticketCount = 3
			}

			for i := 1; i <= ticketCount; i++ {
				var hour int
				if i == 1 {
					hour = 7
				} else if i == 2 {
					hour = 9
				} else {
					hour = 14
				}
				min := 10 + ((d*5 + i*17) % 35)
				tDate := time.Date(2026, time.September, d, hour, min, 0, 0, ict)
				tTons, tRev := insertTicket(batch, q, tDate, i)
				dayTonnage += tTons
				dayRevenue += tRev
			}

			costDate := time.Date(2026, time.September, d, 7, 15, 0, 0, ict)
			insertDailyCost(batch, q, costDate, dayRevenue)
			insertDailyInventory(batch, q, costDate, dayTonnage)
			insertDailyFuel(batch, q, costDate, dayTonnage)
			insertDailyAttendance(batch, q, costDate)
		}
	}

	// 4. Seed Alerts throughout August and September
	alertTemplates := []struct {
		Title, Plate, Severity, Cam, Note string
		Diff                               float64
	}{
		{"Phát hiện xe tải lệch bì kiểm định", "19H-056.22", "critical", "Cam 01 ANPR", "Cân bì lệch +380kg so với đăng kiểm", 380},
		{"Lệch chuẩn định mức nhiên liệu cơ giới", "MX-PC450", "warning", "Cảm biến dầu IoT", "Mức tiêu hao vượt 12% định mức ca sáng", 24},
		{"Khách hàng có công nợ quá hạn hợp đồng", "CUST-319", "warning", "Kế toán mỏ", "Công nợ quá hạn 520 triệu đ cần đối soát", 520000000},
	}

	for i := 1; i <= 18; i++ {
		day := 1 + (i * 31 / 19)
		tmpl := alertTemplates[i%len(alertTemplates)]
		q := quarrySpecs[i%len(quarrySpecs)]
		aDate := time.Date(2026, time.August, day, 10, 30, 0, 0, ict)
		id := fmt.Sprintf("AL-%s-%s-%02d", q.Prefix, aDate.Format("20060102"), i)
		batch.Queue(`
			INSERT INTO alerts (
				id, title, bs, note, time, date, status, severity, phieu,
				bi_dang_ky, bi_thuc_te, lech_bi, cam, created_at
			) VALUES (
				$1, $2, $3, $4, '10:30', $5, 'Đã xác nhận', $6, $7,
				16.0, 16.38, $8, $9, $10
			)
			ON CONFLICT (id) DO UPDATE SET created_at = EXCLUDED.created_at
		`, id, tmpl.Title+" - "+q.Location, tmpl.Plate, tmpl.Note, aDate.Format("02/01/2006"), tmpl.Severity, "TK-"+q.Prefix+"-01", tmpl.Diff, tmpl.Cam, aDate)
	}

	for i := 1; i <= 16; i++ {
		day := 1 + (i * 9 / 17)
		tmpl := alertTemplates[i%len(alertTemplates)]
		q := quarrySpecs[i%len(quarrySpecs)]
		aDate := time.Date(2026, time.September, day, 11, 15, 0, 0, ict)
		id := fmt.Sprintf("AL-%s-%s-%02d", q.Prefix, aDate.Format("20060102"), i)
		batch.Queue(`
			INSERT INTO alerts (
				id, title, bs, note, time, date, status, severity, phieu,
				bi_dang_ky, bi_thuc_te, lech_bi, cam, created_at
			) VALUES (
				$1, $2, $3, $4, '11:15', $5, 'Đã xác nhận', $6, $7,
				16.0, 16.38, $8, $9, $10
			)
			ON CONFLICT (id) DO UPDATE SET created_at = EXCLUDED.created_at
		`, id, tmpl.Title+" - "+q.Location, tmpl.Plate, tmpl.Note, aDate.Format("02/01/2006"), tmpl.Severity, "TK-"+q.Prefix+"-01", tmpl.Diff, tmpl.Cam, aDate)
	}

	todayAlerts := []struct {
		QIdx int
		Tmpl int
		Hour int
	}{
		{0, 0, 8}, // 08:30 AM
		{2, 2, 9}, // 09:30 AM
	}
	for idx, ta := range todayAlerts {
		q := quarrySpecs[ta.QIdx]
		tmpl := alertTemplates[ta.Tmpl]
		aDate := time.Date(2026, time.September, now.Day(), ta.Hour, 30, 0, 0, ict)
		id := fmt.Sprintf("AL-%s-%s-%02d", q.Prefix, aDate.Format("20060102"), idx+1)
		batch.Queue(`
			INSERT INTO alerts (
				id, title, bs, note, time, date, status, severity, phieu,
				bi_dang_ky, bi_thuc_te, lech_bi, cam, created_at
			) VALUES (
				$1, $2, $3, $4, '08:30', $5, 'Chờ xử lý', $6, $7,
				16.0, 16.38, $8, $9, $10
			)
			ON CONFLICT (id) DO UPDATE SET created_at = EXCLUDED.created_at
		`, id, tmpl.Title+" - "+q.Location, tmpl.Plate, tmpl.Note, aDate.Format("02/01/2006"), tmpl.Severity, "TK-"+q.Prefix+"-01", tmpl.Diff, tmpl.Cam, aDate)
	}

	fmt.Printf("📦 Executing monthly seed pipeline (%d operations)...\n", batch.Len())
	br := Pool.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		fmt.Printf("⚠️ Error executing seed batch: %v\n", err)
		return
	}

	fmt.Println("✅ Comprehensive monthly quarry data seed completed successfully!")
}
