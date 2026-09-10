package database

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
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
	Plate  string
	Truck  string
	Driver string
	Phone  string
	Tare   float64
}

var fleetCatalog = []fleetSpec{
	{"19H-056.22", "Xe ben Chenglong 4 chân", "Nguyễn Văn Mạnh", "0982.145.882", 14.80},
	{"88H-042.27", "Xe ben Sinotruk Howo 4 chân", "Trần Đình Trọng", "0984.112.334", 16.20},
	{"29C-345.67", "Xe ben 4 chân Sinotruk", "Trần Văn Kiên", "0988.341.992", 16.42},
	{"29H-882.19", "Xe bồn trộn 12m3", "Lê Văn Tuấn", "0912.445.667", 16.23},
	{"60C-312.78", "Xe đầu kéo mooc ben Howo", "Trương Văn Nam", "0938.667.129", 17.80},
	{"61C-445.89", "Xe ben 4 chân Shacman", "Lê Quốc Bảo", "0903.551.234", 16.30},
	{"90C-128.45", "Xe đầu kéo mooc ben", "Đinh Văn Toàn", "0915.223.789", 17.10},
	{"90C-054.67", "Xe bồn trộn bê tông", "Nguyễn Hoàng Long", "0981.332.654", 15.90},
	{"93C-114.28", "Xe tải ben 4 chân Howo", "Vũ Quốc Đạt", "0978.556.223", 16.50},
	{"93C-098.52", "Xe đầu kéo mooc ben", "Bùi Đức Thịnh", "0972.441.902", 17.20},
}

// SeedMonthlyQuarryData seeds comprehensive ticket, trip, voucher, cost, fuel, and attendance data
// for all 4 quarries across all months of 2026 (Jan to Sep 2026), including today (10/09/2026) and yesterday (09/09/2026).
func SeedMonthlyQuarryData() {
	if Pool == nil {
		return
	}
	ctx := context.Background()

	var checkCount int
	_ = Pool.QueryRow(ctx, "SELECT COUNT(*) FROM tickets WHERE id LIKE 'TK-%-2026%'").Scan(&checkCount)
	if checkCount >= 180 {
		fmt.Printf("✅ Monthly quarry data already exists (%d tickets). Refreshing today & yesterday...\n", checkCount)
	} else {
		fmt.Println("🌱 Seeding Comprehensive Monthly Quarry Tickets (Jan - Sep 2026, 4 Quarries)...")
	}

	ict := time.FixedZone("ICT", 7*60*60)
	now := time.Now().In(ict) // 2026-09-10

	var wg sync.WaitGroup

	// Concurrently seed each quarry's tickets & related vouchers
	for _, qSpec := range quarrySpecs {
		wg.Add(1)
		go func(q quarrySeedSpec) {
			defer wg.Done()
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

			if checkCount < 180 {
				// Months 1 to 8 (Jan to Aug 2026): 5-6 tickets per month
				for m := 1; m <= 8; m++ {
					ticketCount := 5 + (m % 3)
					for i := 1; i <= ticketCount; i++ {
						day := 4 + (i * 4)
						hour := 7 + (i % 9)
						minute := 10 + ((i * 7) % 45)
						tDate := time.Date(2026, time.Month(m), day, hour, minute, 0, 0, ict)
						insertSingleQuarryTicket(ctx, q, tDate, i, rnd)
					}
				}

				// Month 9 (September 2026)
				// Early Sept: Sept 2, 4
				for i := 1; i <= 2; i++ {
					tDate := time.Date(2026, time.September, i*2, 8+i, 20, 0, 0, ict)
					insertSingleQuarryTicket(ctx, q, tDate, 10+i, rnd)
				}

				// This week: Sept 7, 8
				insertSingleQuarryTicket(ctx, q, time.Date(2026, time.September, 7, 9, 30, 0, 0, ict), 21, rnd)
				insertSingleQuarryTicket(ctx, q, time.Date(2026, time.September, 7, 14, 15, 0, 0, ict), 22, rnd)
				insertSingleQuarryTicket(ctx, q, time.Date(2026, time.September, 8, 9, 45, 0, 0, ict), 25, rnd)
				insertSingleQuarryTicket(ctx, q, time.Date(2026, time.September, 8, 15, 00, 0, 0, ict), 26, rnd)
			}

			// Yesterday (Sept 9): 5 tickets throughout the day
			yesterdayHours := []int{7, 9, 11, 14, 16}
			for i, h := range yesterdayHours {
				yDate := time.Date(2026, time.September, 9, h, 15+i*8, 0, 0, ict)
				insertSingleQuarryTicket(ctx, q, yDate, 30+i, rnd)
			}

			// Today (Sept 10): 6 tickets throughout morning/afternoon
			todayHours := []int{7, 8, 9, 10, 11, 12}
			for i, h := range todayHours {
				tDate := time.Date(2026, time.September, 10, h, 10+i*7, 0, 0, ict)
				insertSingleQuarryTicket(ctx, q, tDate, 40+i, rnd)
			}
		}(qSpec)
	}

	wg.Wait()

	// Seed monthly production costs
	seedProductionCostsForAllQuarries(ctx, now)

	// Seed equipment fuel logs
	seedFuelLogsForAllQuarries(ctx, now)

	// Seed attendances
	seedAttendancesForAllQuarries(ctx, now)

	// Seed inventory movements
	seedInventoryMovementsForAllQuarries(ctx, now)

	// Seed alerts
	seedAlertsForAllQuarries(ctx, now)

	fmt.Println("✅ Comprehensive monthly quarry data seed completed successfully!")
}

func insertSingleQuarryTicket(ctx context.Context, q quarrySeedSpec, ticketDate time.Time, seq int, rnd *rand.Rand) {
	ticketID := fmt.Sprintf("TK-%s-%s-%03d", q.Prefix, ticketDate.Format("20060102"), seq)

	mat := materialCatalog[(seq+int(ticketDate.Month()))%len(materialCatalog)]
	cust := customerCatalog[(seq+int(ticketDate.Day()))%len(customerCatalog)]
	flt := fleetCatalog[(seq)%len(fleetCatalog)]

	netWeight := 28.50 + float64(rnd.Intn(550))/100.0
	tareWeight := flt.Tare
	grossWeight := tareWeight + netWeight
	unitPrice := mat.Price
	totalPrice := netWeight * unitPrice

	time1Str := ticketDate.Add(-25 * time.Minute).Format("15:04")
	time2Str := ticketDate.Format("15:04")
	dateStr := ticketDate.Format("02/01/2006")
	doCode := fmt.Sprintf("DO-%s-%03d", q.Prefix, 100+seq)
	formattedPrice := fmt.Sprintf("%.0f đ", totalPrice)

	// 1. Insert into tickets
	_, err := Pool.Exec(ctx, `
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
			created_at = EXCLUDED.created_at,
			date = EXCLUDED.date,
			kl_hang = EXCLUDED.kl_hang,
			don_gia = EXCLUDED.don_gia,
			thanh_tien = EXCLUDED.thanh_tien,
			quarry_code = EXCLUDED.quarry_code
	`,
		ticketID, seq, q.Name, cust.Name, flt.Plate, flt.Truck, flt.Driver, flt.Phone, "RFID-"+flt.Plate,
		"Cân bán hàng", "confirmed", "Đã chốt số", tareWeight, fmt.Sprintf("%.2f", tareWeight),
		grossWeight, fmt.Sprintf("%.2f", grossWeight), fmt.Sprintf("%.2f", netWeight), 0.0,
		fmt.Sprintf("%.2f", netWeight), unitPrice, formattedPrice, time1Str, time2Str, dateStr,
		"Nguyễn Văn Dũng", "Lê Văn Cân 02",
		mat.Name, mat.Standard, doCode, q.Scale, q.Gate, "Xuất mỏ đủ tải", "HD-2026-"+ticketID,
		q.Code, ticketDate, ticketDate,
		fmt.Sprintf(`{"front":{"camera":"Camera 01 ANPR","time":"%s %s"},"rear":{"camera":"Camera 02 Thùng","time":"%s %s"}}`, dateStr, time1Str, dateStr, time2Str),
		fmt.Sprintf(`[{"author":"Hệ thống Cân AI","content":"Xe %s hoàn tất cân tự động %s tại %s.","time":"%s %s"}]`, flt.Plate, mat.Name, q.Scale, dateStr, time2Str),
	)
	if err != nil {
		return
	}

	// 2. Insert corresponding vehicle_trip
	Pool.Exec(ctx, `
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
		ticketDate.Add(-25*time.Minute), ticketDate, seq,
		netWeight, netWeight,
		ticketDate, ticketDate,
	)

	// 3. Insert corresponding sales_voucher
	voucherCode := "PB-" + ticketID
	Pool.Exec(ctx, `
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
		voucherCode, cust.Code, cust.Name, ticketDate.Format("2006-01-02"), q.Location, flt.Plate,
		ticketID, totalPrice, totalPrice*0.1, totalPrice*1.1, totalPrice*1.1,
		"Xuất hàng trạm cân "+q.Scale, ticketDate, ticketDate,
	)
}

func seedProductionCostsForAllQuarries(ctx context.Context, now time.Time) {
	var count int
	_ = Pool.QueryRow(ctx, "SELECT COUNT(*) FROM production_costs WHERE mine_area LIKE '%MO-PT-01%' AND period = '2026-09' AND created_at >= '2026-09-07'").Scan(&count)
	if count >= 10 {
		return
	}

	// Clean out any old/sparse 2026 production cost rows so we have a completely consistent, rich timeline
	_, _ = Pool.Exec(ctx, "DELETE FROM production_costs WHERE period LIKE '2026-%' OR mine_area IN ('Moong Tầng 3 (+45m)', 'Toàn mỏ', 'Tuyến nội bộ', 'Trạm nghiền 01', 'Văn phòng & An toàn')")

	categories := []struct {
		CostType    string
		Category    string
		MonthlyNorm float64
		DailyShare  float64
		Desc        string
	}{
		{"Sản xuất", "Chi phí sản xuất & nổ mìn", 230000000, 0.28, "Chi phí nổ mìn, phụ kiện nổ và vận hành nghiền sàng"},
		{"Nhiên liệu", "Nhiên liệu diesel cơ giới", 185000000, 0.32, "Dầu diesel máy xúc PC450, xúc lật và xe tải moong"},
		{"Nhân công", "Nhân công & Tiền lương ca", 155000000, 0.20, "Tiền lương ca mỏ, thợ máy, tài xế và trạm cân"},
		{"Khấu hao", "Khấu hao trạm nghiền sàng", 75000000, 0.10, "Khấu hao tài sản thiết bị dây chuyền nghiền sàng"},
		{"Vận chuyển", "Vận chuyển nội bộ moong", 95000000, 0.06, "Cước vận chuyển nội bộ từ đáy moong lên trạm nghiền"},
		{"Khác", "Bảo hộ ATLĐ & Môi trường", 32000000, 0.04, "Bảo hộ lao động, quan trắc bụi mỏ và phun sương dập bụi"},
	}

	batch := &pgx.Batch{}

	// 1. Seed historical months 1 to 8 (January to August 2026)
	for _, q := range quarrySpecs {
		mineArea := fmt.Sprintf("%s (%s)", q.Name, q.Code)
		for m := 1; m <= 8; m++ {
			periodStr := fmt.Sprintf("2026-%02d", m)
			// Spread across 4 weeks of the month (days 7, 14, 21, 28)
			for w := 1; w <= 4; w++ {
				day := w * 7
				entryDate := time.Date(2026, time.Month(m), day, 11, 30, 0, 0, now.Location())
				for _, c := range categories {
					norm := (c.MonthlyNorm * q.Ratio) / 4.0
					actual := norm * (0.95 + float64((m*7+w+len(q.Code))%12)/100.0)

					batch.Queue(`
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
		}
	}

	// 2. Seed Month 9 (September 2026) - Current month
	periodSep := "2026-09"
	for _, q := range quarrySpecs {
		mineArea := fmt.Sprintf("%s (%s)", q.Name, q.Code)

		// 2A. Week 1 (Tuần trước: 01/09 - 06/09/2026) - daily operating costs
		for day := 1; day <= 5; day++ {
			entryDate := time.Date(2026, time.September, day, 15, 30, 0, 0, now.Location())
			for _, c := range categories {
				norm := (c.MonthlyNorm * q.Ratio) / 26.0
				actual := norm * (0.94 + float64((day*3+len(q.Code))%13)/100.0)

				batch.Queue(`
					INSERT INTO production_costs (
						cost_type, cost_category, norm_value, norm_unit, actual_value, actual_unit,
						period, mine_area, description, created_at
					) VALUES (
						$1, $2, $3, 'VNĐ', $4, 'VNĐ',
						$5, $6, $7, $8
					)
				`, c.CostType, c.Category, norm, actual, periodSep, mineArea, c.Desc, entryDate)
			}
		}

		// 2B. Week 2 (Tuần này):
		// - 07/09 (Thứ 2)
		// - 08/09 (Thứ 3)
		// - 09/09 (Thứ 4 - Hôm qua)
		weekdays := []struct {
			Day  int
			Hour int
		}{
			{7, 16},
			{8, 15},
			{9, 14},
		}
		for _, wd := range weekdays {
			entryDate := time.Date(2026, time.September, wd.Day, wd.Hour, 0, 0, 0, now.Location())
			for _, c := range categories {
				norm := (c.MonthlyNorm * q.Ratio) / 26.0
				actual := norm * (0.97 + float64((wd.Day*5+len(q.Code))%11)/100.0)

				batch.Queue(`
					INSERT INTO production_costs (
						cost_type, cost_category, norm_value, norm_unit, actual_value, actual_unit,
						period, mine_area, description, created_at
					) VALUES (
						$1, $2, $3, 'VNĐ', $4, 'VNĐ',
						$5, $6, $7, $8
					)
				`, c.CostType, c.Category, norm, actual, periodSep, mineArea, c.Desc, entryDate)
			}
		}

		// 2C. Today (10/09/2026 - Thứ 5) - hourly shift cost entries
		todayShifts := []struct {
			Hour     int
			Minute   int
			CostType string
			Category string
			BaseVal  float64
			Desc     string
		}{
			{8, 0, "Sản xuất", "Chi phí sản xuất & nổ mìn", 24500000, "Ca 1: Khoan nổ mìn tầng 3 và điện cấp liệu trạm nghiền"},
			{10, 0, "Nhiên liệu", "Nhiên liệu diesel cơ giới", 21800000, "Cấp phát dầu diesel ca sáng máy xúc PC450 & xe tải moong"},
			{12, 0, "Nhân công", "Nhân công & Tiền lương ca", 16500000, "Tiền công thợ vận hành, phụ cấp ca độc hại và bữa ăn ca"},
			{14, 0, "Vận chuyển", "Vận chuyển nội bộ moong", 11200000, "Cước vận chuyển đá hộc nội bộ moong lên bãi sơ chế"},
			{14, 30, "Khấu hao", "Khấu hao trạm nghiền sàng", 8500000, "Khấu hao máy nghiền côn và sàng rung 3 tầng"},
			{15, 0, "Khác", "Bảo hộ ATLĐ & Môi trường", 3800000, "Phun sương dập bụi tuyến đường mỏ và bảo hộ lao động"},
		}

		for _, ts := range todayShifts {
			entryDate := time.Date(2026, time.September, 10, ts.Hour, ts.Minute, 0, 0, now.Location())
			norm := ts.BaseVal * q.Ratio
			actual := norm * (0.98 + float64((ts.Hour+len(q.Code))%9)/100.0)

			batch.Queue(`
				INSERT INTO production_costs (
					cost_type, cost_category, norm_value, norm_unit, actual_value, actual_unit,
					period, mine_area, description, created_at
				) VALUES (
					$1, $2, $3, 'VNĐ', $4, 'VNĐ',
					$5, $6, $7, $8
				)
			`, ts.CostType, ts.Category, norm, actual, periodSep, mineArea, ts.Desc, entryDate)
		}
	}

	br := Pool.SendBatch(ctx, batch)
	_ = br.Close()
}

func seedFuelLogsForAllQuarries(ctx context.Context, now time.Time) {
	var count int
	_ = Pool.QueryRow(ctx, "SELECT COUNT(*) FROM equipment_fuel_logs WHERE id LIKE 'FUEL-%-2026%'").Scan(&count)
	if count >= 20 {
		return
	}
	machines := []struct {
		Code, Name, Category string
		Quota                float64
	}{
		{"MX-01", "Máy Xúc Komatsu PC450", "Máy Xúc", 32.5},
		{"MX-02", "Máy Xúc Hyundai R380", "Máy Xúc", 28.0},
		{"DCN-01", "Dây Chuyền Nghiền 01", "Dây Chuyền Nghiền", 45.0},
		{"XB-01", "Xe Ben Howo 88H-042.27", "Xe Vận Tải Moong", 18.5},
		{"XB-02", "Xe Ben Chenglong 19H-056.22", "Xe Vận Tải Moong", 17.5},
	}

	for _, q := range quarrySpecs {
		for d := 1; d <= 10; d++ {
			logDate := time.Date(2026, time.September, d, 17, 30, 0, 0, now.Location())
			dateStr := logDate.Format("02/01/2006")

			for idx, m := range machines {
				id := fmt.Sprintf("FUEL-%s-%s-%s", q.Prefix, logDate.Format("20060102"), m.Code)
				hours := 7.0 + float64((idx+d)%3)
				consumed := hours * m.Quota * (0.95 + float64((d+idx)%10)/100.0)
				issued := consumed + float64((idx*3)%10) - 4.0

				Pool.Exec(ctx, `
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
						created_at = EXCLUDED.created_at,
						actual_fuel_consumed_liters = EXCLUDED.actual_fuel_consumed_liters
				`,
					id, m.Code, m.Name, m.Category,
					hours, m.Quota,
					issued, consumed, issued-consumed,
					q.Location, dateStr+" 17:30", q.Code, logDate,
				)
			}
		}
	}
}

func seedAttendancesForAllQuarries(ctx context.Context, now time.Time) {
	var count int
	_ = Pool.QueryRow(ctx, "SELECT COUNT(*) FROM hr_attendances WHERE id LIKE 'ATT-%-2026%'").Scan(&count)
	if count >= 30 {
		return
	}
	depts := []struct {
		Dept, Role string
		Workers    []string
	}{
		{"Tổ Vận Hành Trạm Cân", "Nhân viên trạm cân", []string{"Nguyễn Văn Dũng", "Lê Văn Cân 02"}},
		{"Đội Cơ Giới & Vận Tải Mỏ", "Tài xế xe ben", []string{"Nguyễn Văn Mạnh", "Trần Đình Trọng", "Lê Hữu Thắng"}},
		{"Xưởng Nghiền Sàng Đá", "Thợ máy nghiền sàng", []string{"Phạm Văn Cường", "Đỗ Văn Long"}},
		{"Phòng Kỹ Thuật & An Toàn", "Kỹ sư trắc địa", []string{"Trần Văn Kiên", "Phạm Hoàng Nam"}},
		{"Ban Điều Hành & Kế Toán", "Điều phối & Kế toán", []string{"Nguyễn Đức Trường", "Nguyễn Thị Thủy"}},
	}

	for _, q := range quarrySpecs {
		// Past 7 days
		for _, dayOffset := range []int{0, 1, 2, 3, 4, 5, 6} {
			attDate := now.AddDate(0, 0, -dayOffset)
			dateStr := attDate.Format("02/01/2006")

			for dIdx, d := range depts {
				for wIdx, worker := range d.Workers {
					attID := fmt.Sprintf("ATT-%s-%s-%d-%d", q.Prefix, attDate.Format("20060102"), dIdx, wIdx)
					checkInTime := time.Date(attDate.Year(), attDate.Month(), attDate.Day(), 6, 30+wIdx*10, 0, 0, now.Location())

					Pool.Exec(ctx, `
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
						attID, fmt.Sprintf("EMP-%s-%d%d", q.Prefix, dIdx, wIdx), worker,
						d.Dept, d.Role,
						dateStr, q.Location, checkInTime,
					)
				}
			}
		}
	}
}

func seedInventoryMovementsForAllQuarries(ctx context.Context, now time.Time) {
	var count int
	_ = Pool.QueryRow(ctx, "SELECT COUNT(*) FROM inventory_inbound WHERE code LIKE 'INV-IN-%'").Scan(&count)
	if count >= 20 {
		return
	}
	for _, q := range quarrySpecs {
		for d := 1; d <= 10; d++ {
			invDate := time.Date(2026, time.September, d, 14, 0, 0, 0, now.Location())
			dateStr := invDate.Format("02/01/2006")

			inID := fmt.Sprintf("INV-IN-%s-%02d", q.Prefix, d)
			qtyIn := 120.0 + float64(d*15)
			Pool.Exec(ctx, `
				INSERT INTO inventory_inbound (code, source, loc, item, qty, quantity, unit, date, status, created_at)
				VALUES ($1, 'Moong Khai Thác', $2, 'Đá 1x2 Bê tông', $3, $3, 'tấn', $4, 'Hoàn thành', $5)
				ON CONFLICT DO NOTHING
			`, inID, q.Location, qtyIn, dateStr, invDate)

			outID := fmt.Sprintf("INV-OUT-%s-%02d", q.Prefix, d)
			qtyOut := 110.0 + float64(d*14)
			Pool.Exec(ctx, `
				INSERT INTO inventory_outbound (code, customer, dest, item, qty, quantity, unit, date, status, created_at)
				VALUES ($1, 'Khách hàng dự án', $2, 'Đá 1x2 Bê tông', $3, $3, 'tấn', $4, 'Hoàn thành', $5)
				ON CONFLICT DO NOTHING
			`, outID, q.Location, qtyOut, dateStr, invDate)
		}
	}
}

func seedAlertsForAllQuarries(ctx context.Context, now time.Time) {
	alerts := []struct {
		Title, Plate, Severity, Cam, Note string
		Diff                               float64
	}{
		{"Phát hiện xe tải lệch bì kiểm định", "19H-056.22", "critical", "Cam 01 ANPR", "Cân bì lệch +380kg so với đăng kiểm", 380},
		{"Lệch chuẩn định mức nhiên liệu cơ giới", "MX-PC450", "warning", "Cảm biến dầu IoT", "Mức tiêu hao vượt 18% định mức ca sáng", 24},
		{"Cảnh báo công nợ khách hàng vượt ngưỡng", "CUST-319", "warning", "Hệ thống Kế toán", "Công nợ vượt hạn mức hợp đồng DO-319", 520000000},
	}

	for _, q := range quarrySpecs {
		for i, a := range alerts {
			id := fmt.Sprintf("AL-%s-20260910-%02d", q.Prefix, i+1)
			todayAlertTime := now.Add(-time.Duration(i*2) * time.Hour)

			Pool.Exec(ctx, `
				INSERT INTO alerts (
					id, title, bs, note, time, date, status, severity, phieu,
					bi_dang_ky, bi_thuc_te, lech_bi, cam, created_at
				) VALUES (
					$1, $2, $3, $4, '10:30', '10/09/2026', 'Chờ xử lý', $5, $6,
					16.0, 16.38, $7, $8, $9
				)
				ON CONFLICT (id) DO UPDATE SET created_at = EXCLUDED.created_at
			`,
				id, a.Title+" - "+q.Location, a.Plate, a.Note, a.Severity, "TK-"+q.Prefix+"-01",
				a.Diff, a.Cam, todayAlertTime,
			)

			idY := fmt.Sprintf("AL-%s-20260909-%02d", q.Prefix, i+1)
			yesterdayAlertTime := now.AddDate(0, 0, -1).Add(-time.Duration(i*2) * time.Hour)
			Pool.Exec(ctx, `
				INSERT INTO alerts (
					id, title, bs, note, time, date, status, severity, phieu,
					bi_dang_ky, bi_thuc_te, lech_bi, cam, created_at
				) VALUES (
					$1, $2, $3, $4, '14:20', '09/09/2026', 'Đã xác nhận', $5, $6,
					16.0, 16.35, $7, $8, $9
				)
				ON CONFLICT (id) DO UPDATE SET created_at = EXCLUDED.created_at
			`,
				idY, a.Title+" - "+q.Location, a.Plate, a.Note, a.Severity, "TK-"+q.Prefix+"-02",
				a.Diff, a.Cam, yesterdayAlertTime,
			)
		}
	}
}
