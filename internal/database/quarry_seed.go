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

	// Check if already seeded
	var count int
	_ = Pool.QueryRow(ctx, "SELECT COUNT(*) FROM quarries").Scan(&count)
	if count > 0 {
		return
	}

	fmt.Println("🌱 Seeding Quarry Core Module data...")

	// 1. Organization
	var orgID string
	_ = Pool.QueryRow(ctx, `
		INSERT INTO organizations (code, name, tax_code, address, status)
		VALUES ('TTC-GROUP', 'TẬP ĐOÀN KHAI THÁC & CHẾ BIẾN KHOÁNG SẢN TTC', '0108992345', 'Tòa nhà TTC Tower, Hà Nội', 'active')
		ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`).Scan(&orgID)

	// 2. Quarries
	var qPT, qTU, qHN string
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

	if qPT == "" {
		return
	}

	// 3. Quarry Areas for MO-PT-01
	var area1, area2, area3 string
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

	// 13. Integration Events Log
	_, _ = Pool.Exec(ctx, `
		INSERT INTO integration_events (source, event_type, external_id, payload, status, error_message, received_at, processed_at)
		VALUES
		('MineSurveySoftware', 'survey_completed', 'SURVEY-20260825-003', '{"quarryCode":"MO-PT-01","cycleCode":"CYCLE-2026-Q3-03","surveyDate":"2026-08-25","operatorName":"KS. Trần Hải Đăng","extractedVolumeM3":43500}', 'completed', NULL, '2026-08-25 15:50:00+07', '2026-08-25 15:50:02+07'),
		('SmartWeighScale', 'scale_ticket', 'TICKET-EXT-1002', '{"ticketNumber":"PC-PT260825-02","licensePlate":"19C-128.45","netWeightKg":35800,"materialCode":"DA-1X2"}', 'completed', NULL, '2026-08-25 10:23:45+07', '2026-08-25 10:23:46+07')
	`)

	fmt.Println("🎉 Quarry Core Module seeded successfully!")
}
