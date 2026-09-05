package database

import (
	"context"
	"fmt"
)

// SeedQuarryModule populates initial realistic demo data for Quarry Survey, 3D Models, Scale & Reconciliations
func SeedQuarryModule() {
	if Pool == nil {
		return
	}
	ctx := context.Background()

	fmt.Println("🌱 Checking & Seeding Quarry Core Module data...")

	// 1. Organization
	var orgID string
	_ = Pool.QueryRow(ctx, `
		INSERT INTO organizations (code, name, tax_code, address, status)
		VALUES ('TTC-GROUP', 'TẬP ĐOÀN KHAI THÁC & CHẾ BIẾN KHOÁNG SẢN TTC', '0108992345', 'Tòa nhà TTC Tower, Hà Nội', 'active')
		ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`).Scan(&orgID)

	// 2. Quarries
	var qPT, qTU, qHN, qBP string
	_ = Pool.QueryRow(ctx, `
		INSERT INTO quarries (code, name, address, longitude, latitude, total_area, boundary, status)
		VALUES ('MO-PT-01', 'Công ty CP TTC - Mỏ đá Phú Thọ', 'Thanh Ba, Phú Thọ', 105.184500, 21.452300, 450000.0, '{"type":"Polygon","coordinates":[[[105.18,21.45],[105.19,21.45],[105.19,21.46],[105.18,21.46],[105.18,21.45]]]}', 'active')
		ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`).Scan(&qPT)

	_ = Pool.QueryRow(ctx, `
		INSERT INTO quarries (code, name, address, longitude, latitude, total_area, boundary, status)
		VALUES ('MO-TU-02', 'Mỏ đá Tân Uyên (Mỏ TTC 02)', 'Tân Uyên, Bình Dương', 106.824100, 11.082500, 320000.0, '{"type":"Polygon","coordinates":[[[106.82,11.08],[106.83,11.08],[106.83,11.09],[106.82,11.09],[106.82,11.08]]]}', 'active')
		ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`).Scan(&qTU)

	_ = Pool.QueryRow(ctx, `
		INSERT INTO quarries (code, name, address, longitude, latitude, total_area, boundary, status)
		VALUES ('MO-HN-03', 'Mỏ đá Hà Nam (Mỏ TTC 03)', 'Kiện Khê, Thanh Liêm, Hà Nam', 105.912400, 20.485100, 580000.0, '{"type":"Polygon","coordinates":[[[105.91,20.48],[105.92,20.48],[105.92,20.49],[105.91,20.49],[105.91,20.48]]]}', 'active')
		ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`).Scan(&qHN)

	_ = Pool.QueryRow(ctx, `
		INSERT INTO quarries (code, name, address, longitude, latitude, total_area, boundary, status)
		VALUES ('MO-BP-04', 'Mỏ đá Bình Phước (Mỏ TTC 04)', 'Chơn Thành, Bình Phước', 106.684500, 11.452300, 410000.0, '{"type":"Polygon","coordinates":[[[106.68,11.45],[106.69,11.45],[106.69,11.46],[106.68,11.46],[106.68,11.45]]]}', 'active')
		ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`).Scan(&qBP)

	// Fallback to fetch IDs if scan was empty
	if qPT == "" {
		_ = Pool.QueryRow(ctx, "SELECT id FROM quarries WHERE code = 'MO-PT-01'").Scan(&qPT)
	}
	if qTU == "" {
		_ = Pool.QueryRow(ctx, "SELECT id FROM quarries WHERE code = 'MO-TU-02'").Scan(&qTU)
	}
	if qHN == "" {
		_ = Pool.QueryRow(ctx, "SELECT id FROM quarries WHERE code = 'MO-HN-03'").Scan(&qHN)
	}
	if qBP == "" {
		_ = Pool.QueryRow(ctx, "SELECT id FROM quarries WHERE code = 'MO-BP-04'").Scan(&qBP)
	}

	// 3. Quarry Areas for MO-PT-01 (Phú Thọ)
	var area1, area2, area3 string
	if qPT != "" {
		_ = Pool.QueryRow(ctx, `
			INSERT INTO quarry_areas (quarry_id, code, name, description, area_m2, status)
			VALUES ($1, 'KV-01', 'Moong Khai Thác Tầng 1 (Cao độ +120m đến +100m)', 'Khu vực nổ mìn bóc đất tầng trên', 120000.0, 'active')
			ON CONFLICT (quarry_id, code) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, qPT).Scan(&area1)

		_ = Pool.QueryRow(ctx, `
			INSERT INTO quarry_areas (quarry_id, code, name, description, area_m2, status)
			VALUES ($1, 'KV-02', 'Moong Khai Thác Tầng 2 (Cao độ +100m đến +80m)', 'Gương đá gốc chất lượng cao bê tông', 150000.0, 'active')
			ON CONFLICT (quarry_id, code) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, qPT).Scan(&area2)

		_ = Pool.QueryRow(ctx, `
			INSERT INTO quarry_areas (quarry_id, code, name, description, area_m2, status)
			VALUES ($1, 'KV-03', 'Bãi Chứa Thành Phẩm & Trạm Nghiền Sàng Số 1', 'Bãi tập kết Đá 1x2, 4x6 và Cát VSI', 85000.0, 'active')
			ON CONFLICT (quarry_id, code) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, qPT).Scan(&area3)
	}

	// Quarry Areas for MO-HN-03 (Hà Nam)
	if qHN != "" {
		_, _ = Pool.Exec(ctx, `
			INSERT INTO quarry_areas (quarry_id, code, name, description, area_m2, status)
			VALUES 
			($1, 'HN-KV-01', 'Moong Kiện Khê Tầng -25m', 'Gương khai thác đá vôi chất lượng cao cho xi măng', 180000.0, 'active'),
			($1, 'HN-KV-02', 'Khu Vực Bãi Nghiền Sàng Cát VSI Kiện Khê', 'Trạm nghiền cát nhân tạo mác cao', 95000.0, 'active')
			ON CONFLICT (quarry_id, code) DO UPDATE SET name = EXCLUDED.name
		`, qHN)
	}

	// Quarry Areas for MO-TU-02 (Tân Uyên)
	if qTU != "" {
		_, _ = Pool.Exec(ctx, `
			INSERT INTO quarry_areas (quarry_id, code, name, description, area_m2, status)
			VALUES 
			($1, 'TU-KV-01', 'Moong Tân Uyên Vỉa 1', 'Vỉa cuội sỏi đồi và đá cấp phối', 140000.0, 'active'),
			($1, 'TU-KV-02', 'Bãi Tuyển Rửa Sỏi Cát Tân Uyên', 'Hồ lắng và bãi thành phẩm', 75000.0, 'active')
			ON CONFLICT (quarry_id, code) DO UPDATE SET name = EXCLUDED.name
		`, qTU)
	}

	// Quarry Areas for MO-BP-04 (Bình Phước)
	if qBP != "" {
		_, _ = Pool.Exec(ctx, `
			INSERT INTO quarry_areas (quarry_id, code, name, description, area_m2, status)
			VALUES 
			($1, 'BP-KV-01', 'Khai Trường Đá Bazan Chơn Thành Tầng 1', 'Vỉa đá bazan cường độ cao', 160000.0, 'active'),
			($1, 'BP-KV-02', 'Bãi Nghiền Base & Thảm Nhựa', 'Trạm nghiền đá thảm bê tông nhựa nóng', 80000.0, 'active')
			ON CONFLICT (quarry_id, code) DO UPDATE SET name = EXCLUDED.name
		`, qBP)
	}

	// 4. Materials
	var mat1, mat2, mat3 string
	_ = Pool.QueryRow(ctx, `
		INSERT INTO quarry_materials (quarry_id, code, name, density_ton_per_m3, status)
		VALUES ($1, 'DA-1X2', 'Đá 1x2 bê tông tiêu chuẩn', 1.55, 'active')
		ON CONFLICT (quarry_id, code) DO UPDATE SET density_ton_per_m3 = EXCLUDED.density_ton_per_m3
		RETURNING id
	`, qPT).Scan(&mat1)

	_ = Pool.QueryRow(ctx, `
		INSERT INTO quarry_materials (quarry_id, code, name, density_ton_per_m3, status)
		VALUES ($1, 'DA-4X6', 'Đá 4x6 móng công trình', 1.58, 'active')
		ON CONFLICT (quarry_id, code) DO UPDATE SET density_ton_per_m3 = EXCLUDED.density_ton_per_m3
		RETURNING id
	`, qPT).Scan(&mat2)

	_ = Pool.QueryRow(ctx, `
		INSERT INTO quarry_materials (quarry_id, code, name, density_ton_per_m3, status)
		VALUES ($1, 'DA-MI', 'Đá mi sàng 0-5mm', 1.52, 'active')
		ON CONFLICT (quarry_id, code) DO UPDATE SET density_ton_per_m3 = EXCLUDED.density_ton_per_m3
		RETURNING id
	`, qPT).Scan(&mat3)

	// 5. Vehicles
	var v1, v2 string
	_ = Pool.QueryRow(ctx, `
		INSERT INTO quarry_vehicles (quarry_id, license_plate, vehicle_type, capacity_ton, status, external_id)
		VALUES ($1, '19H-056.22', 'Xe ben Howo 4 chân', 30.0, 'active', 'VEH-001')
		ON CONFLICT (quarry_id, license_plate) DO UPDATE SET vehicle_type = EXCLUDED.vehicle_type
		RETURNING id
	`, qPT).Scan(&v1)

	_ = Pool.QueryRow(ctx, `
		INSERT INTO quarry_vehicles (quarry_id, license_plate, vehicle_type, capacity_ton, status, external_id)
		VALUES ($1, '19C-128.45', 'Xe đầu kéo mạ kẽm', 35.0, 'active', 'VEH-002')
		ON CONFLICT (quarry_id, license_plate) DO UPDATE SET vehicle_type = EXCLUDED.vehicle_type
		RETURNING id
	`, qPT).Scan(&v2)

	// 6. Survey Cycles
	var cycle1, cycle2, cycle3 string
	_ = Pool.QueryRow(ctx, `
		INSERT INTO survey_cycles (quarry_id, cycle_code, survey_date, survey_started_at, survey_completed_at, survey_method, status, operator_name, external_id, external_source, notes)
		VALUES ($1, 'CYCLE-2026-Q1-01', '2026-03-31', '2026-03-31 08:00:00+07', '2026-03-31 16:30:00+07', 'drone', 'completed', 'KS. Nguyễn Tuấn Anh (TTC Survey)', 'SURVEY-20260331-001', 'MineSurveySoftware', 'Quét mốc cơ sở ban đầu đầu năm 2026')
		ON CONFLICT (quarry_id, cycle_code) DO UPDATE SET status = EXCLUDED.status
		RETURNING id
	`, qPT).Scan(&cycle1)

	_ = Pool.QueryRow(ctx, `
		INSERT INTO survey_cycles (quarry_id, cycle_code, survey_date, survey_started_at, survey_completed_at, previous_cycle_id, survey_method, status, operator_name, external_id, external_source, notes)
		VALUES ($1, 'CYCLE-2026-Q2-02', '2026-06-30', '2026-06-30 08:30:00+07', '2026-06-30 17:00:00+07', $2, 'drone', 'completed', 'KS. Nguyễn Tuấn Anh (TTC Survey)', 'SURVEY-20260630-002', 'MineSurveySoftware', 'Khảo sát tiến độ bóc tầng 1 và khai thác tầng 2')
		ON CONFLICT (quarry_id, cycle_code) DO UPDATE SET status = EXCLUDED.status
		RETURNING id
	`, qPT, cycle1).Scan(&cycle2)

	_ = Pool.QueryRow(ctx, `
		INSERT INTO survey_cycles (quarry_id, cycle_code, survey_date, survey_started_at, survey_completed_at, previous_cycle_id, survey_method, status, operator_name, external_id, external_source, notes)
		VALUES ($1, 'CYCLE-2026-Q3-03', '2026-08-25', '2026-08-25 07:30:00+07', '2026-08-25 15:45:00+07', $2, 'drone', 'completed', 'KS. Trần Hải Đăng (Drone RTK Team)', 'SURVEY-20260825-003', 'MineSurveySoftware', 'Khảo sát định kỳ tháng 8 trước mùa mưa, quét đám mây điểm 3D')
		ON CONFLICT (quarry_id, cycle_code) DO UPDATE SET status = EXCLUDED.status
		RETURNING id
	`, qPT, cycle2).Scan(&cycle3)

	// 7. Volume Calculations
	_, _ = Pool.Exec(ctx, `
		INSERT INTO volume_calculations (survey_cycle_id, quarry_area_id, previous_survey_cycle_id, calculation_type, previous_volume_m3, current_volume_m3, extracted_volume_m3, fill_volume_m3, net_volume_m3, calculation_method, external_id, calculated_at)
		VALUES
		($1, $4, NULL, 'remaining', 0, 500000.0, 0, 0, 500000.0, 'tin_surface', 'VOL-001', '2026-03-31 17:00:00+07'),
		($2, $4, $1, 'extracted', 500000.0, 465000.0, 35000.0, 0, 35000.0, 'tin_surface_diff', 'VOL-002', '2026-06-30 17:30:00+07'),
		($3, $4, $2, 'extracted', 465000.0, 421500.0, 43500.0, 0, 43500.0, 'tin_surface_diff', 'VOL-003', '2026-08-25 16:30:00+07')
	`, cycle1, cycle2, cycle3, area1)

	// 8. Surface Models (3D metadata)
	_, _ = Pool.Exec(ctx, `
		INSERT INTO surface_models (survey_cycle_id, quarry_area_id, model_type, file_format, storage_path, file_size, coordinate_system, min_x, min_y, min_z, max_x, max_y, max_z)
		VALUES
		($1, $2, 'point_cloud', 'laz', 'r2://ttc-quarries/pt-01/cycle-03/pointcloud.laz', 184500000, 'VN2000-UTM48N', 582310.45, 2341200.80, 68.50, 582980.12, 2341850.30, 142.00),
		($1, $2, 'mesh', 'glb', 'r2://ttc-quarries/pt-01/cycle-03/mesh_surface.glb', 52400000, 'VN2000-UTM48N', 582310.45, 2341200.80, 68.50, 582980.12, 2341850.30, 142.00),
		($1, $2, 'dem', 'tif', 'r2://ttc-quarries/pt-01/cycle-03/elevation_dem.tif', 34200000, 'VN2000-UTM48N', 582310.45, 2341200.80, 68.50, 582980.12, 2341850.30, 142.00),
		($1, $2, 'orthophoto', 'tif', 'r2://ttc-quarries/pt-01/cycle-03/ortho_mosaic.tif', 98500000, 'VN2000-UTM48N', 582310.45, 2341200.80, 68.50, 582980.12, 2341850.30, 142.00)
	`, cycle3, area1)

	// 9. Survey Area Elevation Results
	_, _ = Pool.Exec(ctx, `
		INSERT INTO survey_area_results (survey_cycle_id, quarry_area_id, area_m2, min_elevation, max_elevation, avg_elevation)
		VALUES
		($1, $2, 120000.0, 78.2, 142.0, 112.5),
		($1, $3, 150000.0, 68.5, 102.0, 84.3)
		ON CONFLICT (survey_cycle_id, quarry_area_id) DO NOTHING
	`, cycle3, area1, area2)

	// 10. Weighing Transactions
	var t1ID string
	_ = Pool.QueryRow(ctx, `
		INSERT INTO weighing_transactions (quarry_id, vehicle_id, material_id, survey_cycle_id, ticket_number, weight_in_kg, weight_out_kg, net_weight_kg, weighed_in_at, weighed_out_at, status, external_id, external_source, raw_data)
		VALUES ($1, $2, $3, $4, 'PC-PT260825-01', 46520, 12200, 34320, '2026-08-25 09:12:00+07', '2026-08-25 09:15:30+07', 'completed', 'TICKET-EXT-1001', 'SmartWeighScale', '{"camera_front":"cam_01_f.jpg","camera_rear":"cam_01_r.jpg","rfid":"RFID-19H05622","ai_confidence":0.985,"gate_id":"GATE-01"}')
		ON CONFLICT (external_source, external_id) DO UPDATE SET net_weight_kg = EXCLUDED.net_weight_kg
		RETURNING id
	`, qPT, v1, mat1, cycle3).Scan(&t1ID)

	_ = Pool.QueryRow(ctx, `
		INSERT INTO weighing_transactions (quarry_id, vehicle_id, material_id, survey_cycle_id, ticket_number, weight_in_kg, weight_out_kg, net_weight_kg, weighed_in_at, weighed_out_at, status, external_id, external_source, raw_data)
		VALUES ($1, $2, $3, $4, 'PC-PT260825-02', 48900, 13100, 35800, '2026-08-25 10:20:00+07', '2026-08-25 10:23:45+07', 'completed', 'TICKET-EXT-1002', 'SmartWeighScale', '{"camera_front":"cam_02_f.jpg","camera_rear":"cam_02_r.jpg","rfid":"RFID-19C12845","ai_confidence":0.992,"gate_id":"GATE-01"}')
		ON CONFLICT (external_source, external_id) DO UPDATE SET net_weight_kg = EXCLUDED.net_weight_kg
		RETURNING id
	`, qPT, v2, mat1, cycle3).Scan(&t1ID)

	// 11. Production Reconciliations
	// Cycle 2: 35,000 m3 * 1.55 = 54,250 ton vs Actual 53,820 ton (diff 430 ton, 0.79% -> matched)
	_, _ = Pool.Exec(ctx, `
		INSERT INTO production_reconciliations (quarry_id, quarry_area_id, survey_cycle_id, material_id, volume_m3, density_ton_per_m3, expected_weight_ton, actual_weight_ton, difference_ton, difference_percent, status)
		VALUES ($1, $2, $3, $4, 35000.0, 1.55, 54250.0, 53820.0, 430.0, 0.79, 'matched')
	`, qPT, area1, cycle2, mat1)

	// Cycle 3: 43,500 m3 * 1.55 = 67,425 ton vs Actual 65,120 ton (diff 2,305 ton, 3.42% -> warning)
	_, _ = Pool.Exec(ctx, `
		INSERT INTO production_reconciliations (quarry_id, quarry_area_id, survey_cycle_id, material_id, volume_m3, density_ton_per_m3, expected_weight_ton, actual_weight_ton, difference_ton, difference_percent, status)
		VALUES ($1, $2, $3, $4, 43500.0, 1.55, 67425.0, 65120.0, 2305.0, 3.42, 'warning')
	`, qPT, area1, cycle3, mat1)

	// 12. Quarry Alerts
	_, _ = Pool.Exec(ctx, `
		INSERT INTO quarry_alerts (quarry_id, quarry_area_id, survey_cycle_id, alert_type, severity, title, message, expected_value, actual_value, difference_percent, is_read, is_resolved)
		VALUES
		($1, $2, $3, 'volume_weight_discrepancy', 'warning', 'Sai lệch sản lượng chu kỳ Tháng 8: 3.42%', 'Thể tích đo 3D: 43,500 m³ (tương đương 67,425 tấn) so với tổng cân xuất bãi: 65,120 tấn. Chênh lệch 2,305 tấn cần đối chiếu tồn bãi ngoài trời.', 67425.0, 65120.0, 3.42, false, false)
	`, qPT, area1, cycle3)

	// 14. Seed Multi-Quarry Employees
	_, _ = Pool.Exec(ctx, `
		INSERT INTO hr_employees (id, code, name, avatar, gender, dob, phone, email, id_card, address, department, job_position, manager, join_date, contract_type, contracted_hours, work_location, status, certificates, base_salary)
		VALUES
		('EMP-HN-01', 'NV-HN-01', 'Đỗ Quốc Huy', 'H', 'Nam', '10/05/1983', '0912.776.321', 'huydq@ttcgroup.vn', '035083009988', 'TP. Phủ Lý, Tỉnh Hà Nam', 'Ban Giám Đốc', 'Quản Đốc Điều Hành Mỏ Hà Nam', 'HĐQT Tập Đoàn TTC', '01/01/2021', 'Chính thức vô thời hạn', 8.0, 'Văn phòng Điều Hành Mỏ Hà Nam (MO-HN-03)', 'Đang làm việc', '["Chứng chỉ Giám đốc Mỏ Lộ Thiên", "Chứng chỉ An toàn VLNCN"]'::jsonb, 35000000),
		('EMP-HN-02', 'NV-HN-02', 'Bùi Văn Long', 'L', 'Nam', '18/11/1989', '0988.441.223', 'longbv@ttcgroup.vn', '035089004411', 'Kiện Khê, Thanh Liêm, Hà Nam', 'Tổ Vận Hành Trạm Cân', 'Trưởng Trạm Cân Cổng 1 Hà Nam', 'Đỗ Quốc Huy', '15/06/2021', 'Chính thức vô thời hạn', 8.0, 'Trạm Cân Cổng 1 - Mỏ Hà Nam (MO-HN-03)', 'Đang làm việc', '["Chứng chỉ Cân Điện Tử 120T", "Chứng chỉ AI ANPR"]'::jsonb, 17000000),
		('EMP-HN-03', 'NV-HN-03', 'Vũ Đức Mạnh', 'M', 'Nam', '25/08/1991', '0977.123.890', 'manhvd@ttcgroup.vn', '035091005522', 'Thanh Liêm, Hà Nam', 'Đội Cơ Giới & Vận Tải Mỏ', 'Thợ Vận Hành Máy Xúc Moong', 'Đỗ Quốc Huy', '10/08/2021', 'Chính thức vô thời hạn', 8.0, 'Moong Khai Thác Kiện Khê (MO-HN-03)', 'Đang làm việc', '["Chứng chỉ Thợ Lái Máy Xúc Komatsu PC450"]'::jsonb, 19500000),
		('EMP-TU-01', 'NV-TU-01', 'Trương Văn Nam', 'N', 'Nam', '03/09/1986', '0938.667.129', 'namtv@ttcgroup.vn', '075086001234', 'Tân Uyên, Bình Dương', 'Đội Cơ Giới & Vận Tải Mỏ', 'Tài Xế Xe Đầu Kéo Mooc Ben 60C-312.78', 'Nguyễn Đức Trường', '01/03/2022', 'Khoán sản lượng', 8.0, 'Đội Xe Mỏ Tân Uyên (MO-TU-02)', 'Đang làm việc', '["Giấy phép Lái xe FC", "Chứng chỉ An toàn Cơ giới Mỏ"]'::jsonb, 15000000),
		('EMP-TU-02', 'NV-TU-02', 'Võ Văn Thọ', 'T', 'Nam', '14/02/1987', '0903.882.145', 'thovv@ttcgroup.vn', '075087005678', 'Thủ Dầu Một, Bình Dương', 'Tổ Vận Hành Trạm Cân', 'Trưởng Trạm Cân Tân Uyên 100T', 'Nguyễn Đức Trường', '10/05/2022', 'Chính thức vô thời hạn', 8.0, 'Trạm Cân Tân Uyên (MO-TU-02)', 'Đang làm việc', '["Chứng chỉ Kiểm định Cân Điện Tử"]'::jsonb, 17500000),
		('EMP-TU-03', 'NV-TU-03', 'Lê Văn Cân', 'C', 'Nam', '20/07/1990', '0918.552.331', 'canlv@ttcgroup.vn', '075090003344', 'Bến Cát, Bình Dương', 'Xưởng Nghiền Sàng Đá', 'Quản Lý Kho Bãi & Nghiền Sàng Đá Tân Uyên', 'Nguyễn Đức Trường', '15/09/2022', 'Chính thức vô thời hạn', 8.0, 'Xưởng Nghiền Mỏ Tân Uyên (MO-TU-02)', 'Đang làm việc', '["Bằng Kỹ thuật Chế tạo máy", "Chứng chỉ Vận hành Trạm nghiền"]'::jsonb, 18000000),
		('EMP-BP-01', 'NV-BP-01', 'Nguyễn Thanh Sơn', 'S', 'Nam', '12/10/1984', '0979.445.889', 'sonnt@ttcgroup.vn', '070084001122', 'Chơn Thành, Bình Phước', 'Ban Giám Đốc', 'Chỉ Huy Khai Trường Mỏ Bình Phước', 'HĐQT Tập Đoàn TTC', '01/06/2023', 'Chính thức vô thời hạn', 8.0, 'Văn phòng Mỏ Bình Phước (MO-BP-04)', 'Đang làm việc', '["Chứng chỉ Giám đốc Mỏ Lộ Thiên", "Thẻ ATLĐ Nhóm 1"]'::jsonb, 32000000),
		('EMP-BP-02', 'NV-BP-02', 'Phạm Minh Tuấn', 'T', 'Nam', '08/04/1992', '0968.112.556', 'tuanpm@ttcgroup.vn', '070092003344', 'Đồng Xoài, Bình Phước', 'Đội Cơ Giới & Vận Tải Mỏ', 'Tài Xế Xe Ben 93C-114.52', 'Nguyễn Thanh Sơn', '15/07/2023', 'Khoán sản lượng', 8.0, 'Khai trường Mỏ Bình Phước (MO-BP-04)', 'Đang làm việc', '["Giấy phép Lái xe Hạng C/FC", "Thẻ An toàn Vận tải"]'::jsonb, 14000000)
		ON CONFLICT (id) DO UPDATE SET work_location = EXCLUDED.work_location, base_salary = EXCLUDED.base_salary
	`)

	// 15. Seed Multi-Quarry Mining Permits
	_, _ = Pool.Exec(ctx, `
		INSERT INTO mining_permits (id, code, title, mine_name, category, category_label, issuer, license_number, issue_date, expiry_date, capacity, approved_reserve, mined_so_far, mined_percent, depth_level, area, coordinates, status, status_label, days_remaining, files, notes)
		VALUES
		('PERMIT-HN-01', 'GP-HN-2021-BTNMT', 'Giấy phép khai thác đá vôi Mỏ Kiện Khê - Hà Nam (MO-HN-03)', 'Mỏ đá Hà Nam (Mỏ TTC 03)', 'mining_license', 'Giấy phép khai thác mỏ', 'Bộ Tài nguyên và Môi trường', 'Số 142/GP-BTNMT', '15/04/2021', '15/04/2046', '1.500.000 Tấn/năm', '28.000.000 Tấn', '6.850.000 Tấn', 24, 'Mức cao +60m đến -25m', '58.0 Hecta', '20.4851°N, 105.9124°E', 'valid', 'Còn hiệu lực', 7160, '[{"name":"Quyet_Dinh_GP_Kien_Khe_142.pdf","size":"5.1 MB","type":"pdf","url":"#"}]'::jsonb, 'Khai thác đá vôi làm VLXD thông thường và nguyên liệu xi măng'),
		('PERMIT-HN-02', 'DTM-HN-2021-BTNMT', 'Phê duyệt Báo cáo ĐTM & Bảo vệ môi trường Mỏ Kiện Khê (MO-HN-03)', 'Mỏ đá Hà Nam (Mỏ TTC 03)', 'environmental', 'Giấy phép môi trường & ĐTM', 'Bộ Tài nguyên và Môi trường', 'Số 89/QĐ-BTNMT', '20/05/2021', '15/04/2046', 'Hồ lắng xử lý nước rửa đá 600 m³/ngày', 'Ký quỹ BVMT 6.5 Tỷ đ', 'Đã ký quỹ 100%', 100, 'Khu phụ trợ moong', '58.0 Hecta', '20.4851°N, 105.9124°E', 'valid', 'Còn hiệu lực', 7160, '[]'::jsonb, 'Định kỳ quan trắc môi trường 4 đợt/năm'),
		('PERMIT-TU-01', 'GP-TU-2020-UBND', 'Giấy phép khai thác đá xây dựng Mỏ Tân Uyên (MO-TU-02)', 'Mỏ đá Tân Uyên (Mỏ TTC 02)', 'mining_license', 'Giấy phép khai thác mỏ', 'UBND Tỉnh Bình Dương', 'Số 78/GP-UBND', '10/08/2020', '10/08/2040', '900.000 Tấn/năm', '16.500.000 Tấn', '4.120.000 Tấn', 25, 'Mức cao +45m đến -10m', '32.0 Hecta', '11.0825°N, 106.8241°E', 'valid', 'Còn hiệu lực', 5088, '[{"name":"Giay_Phep_Tan_Uyen_78.pdf","size":"4.6 MB","type":"pdf","url":"#"}]'::jsonb, 'Cung ứng vật liệu đá xây dựng cho Vành Đai 3 TP.HCM'),
		('PERMIT-BP-01', 'GP-BP-2022-BTNMT', 'Giấy phép khai thác đá xây dựng Mỏ Chơn Thành - Bình Phước (MO-BP-04)', 'Mỏ đá Bình Phước (Mỏ TTC 04)', 'mining_license', 'Giấy phép khai thác mỏ', 'Bộ Tài nguyên và Môi trường', 'Số 88/GP-BTNMT', '05/11/2022', '05/11/2052', '1.200.000 Tấn/năm', '22.000.000 Tấn', '1.250.000 Tấn', 6, 'Mức cao +80m đến +20m', '41.0 Hecta', '11.4523°N, 106.6845°E', 'valid', 'Còn hiệu lực', 9550, '[{"name":"GP_Chon_Thanh_88.pdf","size":"4.9 MB","type":"pdf","url":"#"}]'::jsonb, 'Khai thác đá granite xây dựng và đá base cấp phối')
		ON CONFLICT (id) DO NOTHING
	`)

	// 16. Seed Multi-Quarry Tickets
	_, _ = Pool.Exec(ctx, `
		INSERT INTO tickets (id, ben_ban, ben_mua, bien_so, loai_xe, lai_xe, sdt_lai_xe, loai, stage, stage_label, can_l1, kl1, can_l2, kl2, kl_hang, kl_tinh_tien, don_gia, thanh_tien, time1, time2, date, nguoi_can1, mat_hang, quy_cach, do_code, tram_can, cong_can, ghi_chu, quarry_code)
		VALUES
		('HN281026-001', 'Công ty Cổ phần Mỏ Đá TTC - Chi nhánh Hà Nam', 'Tổng Công Ty XD Trường Sơn', '88H-042.27', 'Xe ben 4 chân Sinotruk', 'Trần Đình Trọng', '0984.112.334', 'Cân bán hàng', 'confirmed', 'Đã chốt số', 16.50, '16.50', 51.20, '51.20', '34.70', '34.70', 230000, '7.981.000 đ', '08:30', '09:05', '28/10/2026', 'Bùi Văn Long', 'Đá 1x2 Bê Tông Kiện Khê', 'TCVN 7570:2006', 'DO-HN-011', 'Trạm Cân Kiện Khê 120T', 'Cổng 1', 'Xuất mỏ Hà Nam dự án Cao tốc', 'MO-HN-03'),
		('HN281026-002', 'Công ty Cổ phần Mỏ Đá TTC - Chi nhánh Hà Nam', 'Công ty Bê Tông Xuân Mai', '29C-781.90', 'Xe đầu kéo Howo mooc ben', 'Vũ Đức Mạnh', '0977.123.890', 'Cân bán hàng', 'confirmed', 'Đã chốt số', 17.20, '17.20', 54.80, '54.80', '37.60', '37.60', 185000, '6.956.000 đ', '09:15', '09:50', '28/10/2026', 'Bùi Văn Long', 'Đá Base Cấp Phối Dmax25', 'TCVN 8859:2011', 'DO-HN-012', 'Trạm Cân Kiện Khê 120T', 'Cổng 1', 'Đá base xuất mỏ Hà Nam', 'MO-HN-03'),
		('TU281026-001', 'Công ty Cổ phần Mỏ Đá TTC - Chi nhánh Bình Dương', 'Tập đoàn Đèo Cả (Vành Đai 3)', '60C-312.78', 'Xe đầu kéo mooc ben', 'Trương Văn Nam', '0938.667.129', 'Cân bán hàng', 'confirmed', 'Đã chốt số', 17.80, '17.80', 52.30, '52.30', '34.50', '34.50', 220000, '7.590.000 đ', '10:30', '11:00', '28/10/2026', 'Võ Văn Thọ', 'Đá Base Cấp Phối Vành Đai 3', 'TCVN 8859:2011', 'DO-TU-002', 'Trạm Cân Tân Uyên 100T', 'Cổng Cân Số 1', 'Dự án đường Vành Đai 3 TP.HCM', 'MO-TU-02'),
		('BP281026-001', 'Công ty Cổ phần Mỏ Đá TTC - Chi nhánh Bình Phước', 'Tổng Công Ty XD Miền Đông', '93C-114.52', 'Xe ben 4 chân Dongfeng', 'Phạm Minh Tuấn', '0968.112.556', 'Cân bán hàng', 'confirmed', 'Đã chốt số', 16.10, '16.10', 48.90, '48.90', '32.80', '32.80', 215000, '7.052.000 đ', '08:00', '08:35', '28/10/2026', 'Nguyễn Thanh Sơn', 'Đá 1x2 Granite Chơn Thành', 'TCVN 7570:2006', 'DO-BP-001', 'Trạm Cân Chơn Thành 100T', 'Cổng 1', 'Xuất mỏ Chơn Thành', 'MO-BP-04')
		ON CONFLICT (id) DO NOTHING
	`)

	// 17. Seed Vehicles for other quarries
	_, _ = Pool.Exec(ctx, `
		INSERT INTO vehicles (bs, loai, chu_xe, tai_trong, count, rfid, han_dang_kiem, status, date, bi, unit, ownership_type, current_driver_name, quarry_code)
		VALUES
		('88H-042.27', 'Xe ben Howo 4 chân 371HP', 'Công ty Cổ phần Mỏ Đá TTC', 30, 84, 'RFID-88H-042', '15/10/2027', 'Hoạt động', '28/10/2026', 15.42, 'tấn', 'company', 'Trần Đình Trọng', 'MO-HN-03'),
		('29C-781.90', 'Xe đầu kéo Howo mooc ben', 'Công ty Cổ phần Mỏ Đá TTC', 35, 62, 'RFID-29C-781', '20/12/2027', 'Hoạt động', '28/10/2026', 17.20, 'tấn', 'company', 'Vũ Đức Mạnh', 'MO-HN-03'),
		('60C-312.78', 'Xe đầu kéo mooc ben 3 trục', 'Công ty Cổ phần Mỏ Đá TTC', 34, 76, 'RFID-60C-312', '10/08/2027', 'Hoạt động', '28/10/2026', 17.80, 'tấn', 'company', 'Trương Văn Nam', 'MO-TU-02'),
		('93C-114.52', 'Xe tải ben Dongfeng 4 chân', 'Công ty Cổ phần Mỏ Đá TTC', 28, 48, 'RFID-93C-114', '05/06/2027', 'Hoạt động', '28/10/2026', 16.10, 'tấn', 'company', 'Phạm Minh Tuấn', 'MO-BP-04')
		ON CONFLICT (bs) DO UPDATE SET quarry_code = EXCLUDED.quarry_code, current_driver_name = EXCLUDED.current_driver_name;
	`)

	fmt.Println("🎉 Quarry Core Module seeded successfully!")
}
