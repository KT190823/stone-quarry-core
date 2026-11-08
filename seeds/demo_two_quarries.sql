-- Dữ liệu giả lập cho môi trường demo riêng, không tự chạy khi khởi động.
-- Yêu cầu schema từ Migrate và MigrateQuarryModule đã được tạo đầy đủ.
-- Chạy thủ công trong một phiên PostgreSQL; không chạy đồng thời nhiều phiên seed.
-- Không xóa hoặc cập nhật dữ liệu cũ. Mã DEMO2 được dành riêng cho bộ này.
-- Mỗi bảng được bao phủ tạo 5 bản ghi/mỏ hoặc 5 bản ghi chung như ghi bên dưới.
-- Nếu mã DEMO2 đã bị sử dụng cho dữ liệu khác, các dòng trùng được giữ nguyên;
-- cần kiểm tra xung đột trước khi dùng trong môi trường demo thực tế.
-- Đường dẫn mô hình/ảnh chỉ là metadata minh họa, không chứa tệp thật.
-- Các bảng không có quan hệ mỏ đáng tin cậy được tạo dữ liệu chung.
-- integration_events không có quan hệ mỏ nên cũng không tạo dữ liệu giả.
-- Seed cũ đã có hai mỏ, vùng, cân và dữ liệu rải rác; số lượng ở đây là
-- số bản ghi demo mới dự kiến, không phải tổng số sau khi cộng dữ liệu cũ.
BEGIN;

INSERT INTO quarries (code, name, address, longitude, latitude, total_area, status)
VALUES ('MO-PT-01', 'Công ty CP TTC - Mỏ đá Phú Thọ', 'Thanh Ba, Phú Thọ', 105.1845, 21.4523, 450000, 'active'),
       ('MO-TU-02', 'Mỏ đá Tân Uyên (Mỏ TTC 02)', 'Tân Uyên, Bình Dương', 106.8241, 11.0825, 320000, 'active')
ON CONFLICT (code) DO NOTHING;

-- Bộ chỉ số chung: mã mỏ và 5 dòng ổn định cho mỗi mỏ.
CREATE TEMP TABLE demo2_rows ON COMMIT DROP AS
SELECT q.id AS quarry_id, q.code AS quarry_code, q.name AS quarry_name,
       CASE q.code WHEN 'MO-PT-01' THEN 'PT' ELSE 'TU' END AS site, n
FROM quarries q CROSS JOIN generate_series(1, 5) AS g(n)
WHERE q.code IN ('MO-PT-01', 'MO-TU-02');

INSERT INTO quarry_areas (id, quarry_id, code, name, area_m2, status)
SELECT md5('demo2:area:' || d.site || ':' || d.n)::uuid, d.quarry_id, 'DEMO2-AREA-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'Khu demo ' || d.n, 10000 + d.n * 1000, 'active'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO quarry_materials (id, quarry_id, code, name, density_ton_per_m3, status)
SELECT md5('demo2:material:' || d.site || ':' || d.n)::uuid, d.quarry_id, 'DEMO2-MAT-' || d.site || '-' || lpad(d.n::text, 2, '0'), (ARRAY['Đá 1x2','Đá 4x6','Đá mi','Đá base','Đá hộc'])[d.n], 1.55, 'active'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO quarry_vehicles (id, quarry_id, license_plate, vehicle_type, capacity_ton, status)
SELECT md5('demo2:vehicle:' || d.site || ':' || d.n)::uuid, d.quarry_id, 'D2-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'Xe ben demo', 30, 'active'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO survey_cycles (id, quarry_id, cycle_code, survey_date, survey_method, status, operator_name)
SELECT md5('demo2:cycle:' || d.site || ':' || d.n)::uuid, d.quarry_id, 'DEMO2-SURVEY-' || d.site || '-' || lpad(d.n::text, 2, '0'), DATE '2026-01-01' + (d.n * 30), 'drone', 'completed', 'Tổ khảo sát demo'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO survey_area_results (survey_cycle_id, quarry_area_id, area_m2, min_elevation, max_elevation, avg_elevation)
SELECT c.id, a.id, a.area_m2, 50, 100, 75
FROM demo2_rows d JOIN quarry_areas a ON a.quarry_id = d.quarry_id AND a.code = 'DEMO2-AREA-' || d.site || '-' || lpad(d.n::text, 2, '0') JOIN survey_cycles c ON c.quarry_id = d.quarry_id AND c.cycle_code = 'DEMO2-SURVEY-' || d.site || '-' || lpad(d.n::text, 2, '0')
ON CONFLICT DO NOTHING;

INSERT INTO volume_calculations (id, survey_cycle_id, quarry_area_id, calculation_type, current_volume_m3, extracted_volume_m3, net_volume_m3, external_id)
SELECT md5('demo2:volume:' || d.site || ':' || d.n)::uuid, c.id, a.id, 'extracted', 100000 - d.n * 20, 20, 20, 'DEMO2-VOL-' || d.site || '-' || lpad(d.n::text, 2, '0')
FROM demo2_rows d JOIN quarry_areas a ON a.quarry_id = d.quarry_id AND a.code = 'DEMO2-AREA-' || d.site || '-' || lpad(d.n::text, 2, '0') JOIN survey_cycles c ON c.quarry_id = d.quarry_id AND c.cycle_code = 'DEMO2-SURVEY-' || d.site || '-' || lpad(d.n::text, 2, '0')
ON CONFLICT DO NOTHING;

INSERT INTO surface_models (id, survey_cycle_id, quarry_area_id, model_type, file_format, storage_path)
SELECT md5('demo2:surface:' || d.site || ':' || d.n)::uuid, c.id, a.id, 'point_cloud', 'laz', 'demo://two-quarries/' || d.site || '/' || d.n || '.laz'
FROM demo2_rows d JOIN quarry_areas a ON a.quarry_id = d.quarry_id AND a.code = 'DEMO2-AREA-' || d.site || '-' || lpad(d.n::text, 2, '0') JOIN survey_cycles c ON c.quarry_id = d.quarry_id AND c.cycle_code = 'DEMO2-SURVEY-' || d.site || '-' || lpad(d.n::text, 2, '0')
ON CONFLICT DO NOTHING;

INSERT INTO weighing_transactions (id, quarry_id, vehicle_id, material_id, survey_cycle_id, ticket_number, weight_in_kg, weight_out_kg, net_weight_kg, weighed_in_at, weighed_out_at, status, external_source, external_id)
SELECT md5('demo2:weigh:' || d.site || ':' || d.n)::uuid, d.quarry_id, v.id, m.id, c.id, 'DEMO2-TICKET-' || d.site || '-' || lpad(d.n::text, 2, '0'), 40000, 15000, 25000, c.survey_date + time '09:00', c.survey_date + time '09:15', 'completed', 'demo_two_quarries', 'DEMO2-WEIGH-' || d.site || '-' || lpad(d.n::text, 2, '0')
FROM demo2_rows d JOIN quarry_vehicles v ON v.quarry_id = d.quarry_id AND v.license_plate = 'D2-' || d.site || '-' || lpad(d.n::text, 2, '0') JOIN quarry_materials m ON m.quarry_id = d.quarry_id AND m.code = 'DEMO2-MAT-' || d.site || '-' || lpad(d.n::text, 2, '0') JOIN survey_cycles c ON c.quarry_id = d.quarry_id AND c.cycle_code = 'DEMO2-SURVEY-' || d.site || '-' || lpad(d.n::text, 2, '0')
ON CONFLICT DO NOTHING;

INSERT INTO vehicle_trip_images (id, weighing_transaction_id, image_type, storage_path)
SELECT md5('demo2:image:' || d.site || ':' || d.n)::uuid, w.id, 'scale', 'demo://two-quarries/' || d.site || '/weigh-' || d.n || '.jpg'
FROM demo2_rows d JOIN weighing_transactions w ON w.external_source = 'demo_two_quarries' AND w.external_id = 'DEMO2-WEIGH-' || d.site || '-' || lpad(d.n::text, 2, '0') AND w.quarry_id = d.quarry_id
ON CONFLICT DO NOTHING;

INSERT INTO production_reconciliations (id, quarry_id, quarry_area_id, survey_cycle_id, material_id, volume_m3, density_ton_per_m3, expected_weight_ton, actual_weight_ton, difference_ton, difference_percent, status)
SELECT md5('demo2:recon:' || d.site || ':' || d.n)::uuid, d.quarry_id, a.id, c.id, m.id, 20, 1.55, 31, 25, 6, 19.3548, 'warning'
FROM demo2_rows d JOIN quarry_areas a ON a.quarry_id = d.quarry_id AND a.code = 'DEMO2-AREA-' || d.site || '-' || lpad(d.n::text, 2, '0') JOIN survey_cycles c ON c.quarry_id = d.quarry_id AND c.cycle_code = 'DEMO2-SURVEY-' || d.site || '-' || lpad(d.n::text, 2, '0') JOIN quarry_materials m ON m.quarry_id = d.quarry_id AND m.code = 'DEMO2-MAT-' || d.site || '-' || lpad(d.n::text, 2, '0')
ON CONFLICT DO NOTHING;

INSERT INTO quarry_alerts (id, quarry_id, quarry_area_id, survey_cycle_id, alert_type, severity, title, message, expected_value, actual_value, difference_percent)
SELECT md5('demo2:alert:' || d.site || ':' || d.n)::uuid, d.quarry_id, a.id, c.id, 'volume_weight_discrepancy', 'warning', 'Đối soát demo ' || d.n, 'Kiểm tra chênh lệch sản lượng demo', 31, 25, 19.3548
FROM demo2_rows d JOIN quarry_areas a ON a.quarry_id = d.quarry_id AND a.code = 'DEMO2-AREA-' || d.site || '-' || lpad(d.n::text, 2, '0') JOIN survey_cycles c ON c.quarry_id = d.quarry_id AND c.cycle_code = 'DEMO2-SURVEY-' || d.site || '-' || lpad(d.n::text, 2, '0')
ON CONFLICT DO NOTHING;

INSERT INTO vehicles (bs, loai, tai_trong, status, quarry_code)
SELECT 'D2-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'Xe ben demo', 30, 'Hoạt động', d.quarry_code
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO tickets (id, quarry_code, ben_ban, ben_mua, bien_so, loai, stage, stage_label, can_l1, can_l2, kl_hang, kl_tinh_tien, don_gia, thanh_tien, date, mat_hang)
SELECT 'DEMO2-TICKET-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, d.quarry_name, 'Khách hàng demo ' || d.site || ' ' || d.n, 'D2-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'Cân bán hàng', 'confirmed', 'Đã chốt số', 15, '40', '25', '25', 200000, '5000000', '28/10/2026', (ARRAY['Đá 1x2','Đá 4x6','Đá mi','Đá base','Đá hộc'])[d.n]
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO vehicle_assignments (code, quarry_code, date, license_plate, vehicle_type, driver_name, status)
SELECT 'DEMO2-ASSIGN-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, '28/10/2026', 'D2-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'Xe ben demo', 'Lái xe demo ' || d.n, 'assigned'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO equipment_fuel_logs (id, quarry_code, equipment_code, equipment_name, hours_worked_today, fuel_quota_liters_per_hour, actual_fuel_consumed_liters, variance_status)
SELECT 'DEMO2-FUEL-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, 'DEMO2-EQUIP-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'Máy xúc demo ' || d.n, 8, 20, 160, 'normal'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO blasting_passports (id, code, quarry_code, mine_name, blast_date, location, hole_count, hole_depth_meters, designed_rock_volume_m3, safety_status)
SELECT 'DEMO2-BLAST-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'DEMO2-BLAST-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, d.quarry_name, '28/10/2026', 'Khu demo ' || d.n, 10, 8, 1000, 'Đạt'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO alerts (id, quarry_code, title, note, date, status, severity, phieu)
SELECT 'DEMO2-ALERT-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, 'Cảnh báo demo ' || d.n, 'Kiểm tra phiếu cân', '28/10/2026', 'open', 'warning', 'DEMO2-TICKET-' || d.site || '-' || lpad(d.n::text, 2, '0')
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO reports (id, quarry_code, quarry_name, name, item, type, period, plan, actual, unit, status)
SELECT 'DEMO2-REPORT-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, d.quarry_name, 'Báo cáo demo ' || d.n, 'Đá xây dựng', 'production', '2026-10', 1000, 950, 'tấn', 'completed'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO production_stages (code, quarry_code, quarry_name, stage_number, stage_name, volume_month, volume_ytd, unit, status)
SELECT 'DEMO2-STAGE-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, d.quarry_name, d.n, 'Công đoạn demo ' || d.n, 1000, 10000, 'tấn', 'active'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO resource_taxes (code, quarry_code, quarry_name, mineral_type, mined_volume, unit, tax_price_per_unit, tax_rate, resource_tax_amount, environmental_fee, total_payable, status)
SELECT 'DEMO2-TAX-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, d.quarry_name, 'Đá xây dựng', 1000, 'tấn', 100000, '10%', 10000000, 1000000, 11000000, 'draft'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO statutory_reports (id, code, quarry_code, quarry_name, name, title, period, mined_volume, unit, status)
SELECT 'DEMO2-STAT-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'DEMO2-STAT-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, d.quarry_name, 'Báo cáo khai thác demo', 'Báo cáo kỳ ' || d.n, '2026-0' || d.n, 1000, 'tấn', 'draft'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

INSERT INTO sales_contracts (code, quarry_code, customer_name, project_name, start_date, end_date, total_committed_tons, total_delivered_tons, completion_rate, status)
SELECT 'DEMO2-CONTRACT-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, 'Khách hàng demo ' || d.site || ' ' || d.n, 'Công trình demo ' || d.n, '01/01/2026', '31/12/2026', 1000, 25, 2.5, 'active'
FROM demo2_rows d
ON CONFLICT DO NOTHING;

-- Customers và suppliers có thể chưa có quarry_code; tên demo chứa mã mỏ.
INSERT INTO customers (name, ten, status)
SELECT 'Khách hàng demo ' || d.site || ' ' || d.n, 'Khách hàng demo ' || d.site || ' ' || d.n, 'active'
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM customers existing WHERE existing.name = 'Khách hàng demo ' || d.site || ' ' || d.n)
ON CONFLICT DO NOTHING;

INSERT INTO suppliers (name, ten, status)
SELECT 'Nhà cung cấp demo ' || d.site || ' ' || d.n, 'Nhà cung cấp demo ' || d.site || ' ' || d.n, 'active'
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM suppliers existing WHERE existing.name = 'Nhà cung cấp demo ' || d.site || ' ' || d.n)
ON CONFLICT DO NOTHING;

INSERT INTO crusher_plants (plant_code, plant_name, quarry_code, capacity_ton_per_hour)
SELECT 'DEMO2-CRUSH-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'Trạm nghiền demo ' || d.n, d.quarry_code, 150
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM crusher_plants existing WHERE existing.plant_code = 'DEMO2-CRUSH-' || d.site || '-' || lpad(d.n::text, 2, '0'))
ON CONFLICT DO NOTHING;

INSERT INTO inventory_inbound (code, quarry_code, quarry_name, source, item, qty, quantity, unit, date, status)
SELECT 'DEMO2-IN-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, d.quarry_name, 'Khai trường', 'Đá 1x2', 100, 100, 'tấn', '28/10/2026', 'completed'
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM inventory_inbound existing WHERE existing.code = 'DEMO2-IN-' || d.site || '-' || lpad(d.n::text, 2, '0'))
ON CONFLICT DO NOTHING;

INSERT INTO inventory_outbound (code, quarry_code, quarry_name, customer, item, qty, quantity, unit, date, status)
SELECT 'DEMO2-OUT-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, d.quarry_name, 'Khách hàng demo ' || d.site || ' ' || d.n, 'Đá 1x2', 25, 25, 'tấn', '28/10/2026', 'completed'
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM inventory_outbound existing WHERE existing.code = 'DEMO2-OUT-' || d.site || '-' || lpad(d.n::text, 2, '0'))
ON CONFLICT DO NOTHING;

INSERT INTO inventory_stocktake (code, quarry_code, quarry_name, zone, item, volume, survey, erp, book, actual, diff, quantity, unit, date, status)
SELECT 'DEMO2-STOCK-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, d.quarry_name, 'Khu demo ' || d.n, 'Đá 1x2', 75, 75, 75, 75, 75, 0, 75, 'tấn', '28/10/2026', 'completed'
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM inventory_stocktake existing WHERE existing.code = 'DEMO2-STOCK-' || d.site || '-' || lpad(d.n::text, 2, '0'))
ON CONFLICT DO NOTHING;

INSERT INTO inventory_movements (code, quarry_code, quarry_name, from_loc, to_loc, item, qty, quantity, unit, date, status)
SELECT 'DEMO2-MOVE-' || d.site || '-' || lpad(d.n::text, 2, '0'), d.quarry_code, d.quarry_name, 'Khai trường', 'Bãi demo ' || d.n, 'Đá 1x2', 100, 100, 'tấn', '28/10/2026', 'completed'
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM inventory_movements existing WHERE existing.code = 'DEMO2-MOVE-' || d.site || '-' || lpad(d.n::text, 2, '0'))
ON CONFLICT DO NOTHING;

-- Các bảng JSON dùng quarry_code bên trong data theo seed hiện tại.
INSERT INTO gps_fleet (code, data)
SELECT 'DEMO2-GPS-' || d.site || '-' || lpad(d.n::text, 2, '0'), jsonb_build_object('id', 'DEMO2-GPS-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'code', 'DEMO2-GPS-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'quarry_code', d.quarry_code, 'quarry_name', d.quarry_name, 'plate', 'D2-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'driverName', 'Lái xe demo ' || d.n, 'status', 'moving', 'lat', CASE d.site WHEN 'PT' THEN 21.4523 ELSE 11.0825 END, 'lng', CASE d.site WHEN 'PT' THEN 105.1845 ELSE 106.8241 END, 'speed', 30)
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM gps_fleet existing WHERE existing.code = 'DEMO2-GPS-' || d.site || '-' || lpad(d.n::text, 2, '0'))
ON CONFLICT DO NOTHING;

INSERT INTO fuel_theft_audits (code, data)
SELECT 'DEMO2-AUDIT-' || d.site || '-' || lpad(d.n::text, 2, '0'), jsonb_build_object('id', 'DEMO2-AUDIT-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'code', 'DEMO2-AUDIT-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'quarry_code', d.quarry_code, 'quarry_name', d.quarry_name, 'plate', 'D2-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'driverName', 'Lái xe demo ' || d.n, 'shiftDate', '28/10/2026', 'actualFuelConsumedLiters', 40, 'theoreticalNormLiters', 40, 'fuelVarianceLiters', 0, 'riskLevel', 'SAVING')
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM fuel_theft_audits existing WHERE existing.code = 'DEMO2-AUDIT-' || d.site || '-' || lpad(d.n::text, 2, '0'))
ON CONFLICT DO NOTHING;

INSERT INTO fuel_norms (code, data)
SELECT 'DEMO2-NORM-' || d.site || '-' || lpad(d.n::text, 2, '0'), jsonb_build_object('id', 'DEMO2-NORM-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'code', 'DEMO2-NORM-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'quarry_code', d.quarry_code, 'quarry_name', d.quarry_name, 'plate', 'D2-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'driverName', 'Lái xe demo ' || d.n, 'vehicleType', 'Xe ben demo', 'normLoadedLitersPer100Km', 40)
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM fuel_norms existing WHERE existing.code = 'DEMO2-NORM-' || d.site || '-' || lpad(d.n::text, 2, '0'))
ON CONFLICT DO NOTHING;

INSERT INTO yard_checkinout (code, data)
SELECT 'DEMO2-YARD-' || d.site || '-' || lpad(d.n::text, 2, '0'), jsonb_build_object('id', 'DEMO2-YARD-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'code', 'DEMO2-YARD-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'quarry_code', d.quarry_code, 'quarry_name', d.quarry_name, 'plate', 'D2-' || d.site || '-' || lpad(d.n::text, 2, '0'), 'driverName', 'Lái xe demo ' || d.n, 'actionType', 'CHECK_OUT', 'timestamp', '07:05 28/10/2026', 'isWithin20mGeofence', true)
FROM demo2_rows d
WHERE NOT EXISTS (SELECT 1 FROM yard_checkinout existing WHERE existing.code = 'DEMO2-YARD-' || d.site || '-' || lpad(d.n::text, 2, '0'))
ON CONFLICT DO NOTHING;


-- Kế hoạch theo mỏ: tổng bốn quý và thực hiện lũy kế khớp chỉ tiêu.
INSERT INTO mining_plans (id, quarry_code, mine, item, annual_target, unit, q1_plan, q1_actual, q2_plan, q2_actual, q3_plan, q3_actual, q4_plan, q4_actual, ytd_actual, completion_rate, status, status_label)
SELECT 'DEMO2-PLAN-' || site || '-' || lpad(n::text, 2, '0'), quarry_code, quarry_name,
       (ARRAY['Đá 1x2','Đá 4x6','Đá mi','Đá base','Đá hộc'])[n],
       1000, 'tấn', 250, 250, 250, 250, 250, 250, 250, 200, 950, 95, 'active', 'Đang thực hiện'
FROM demo2_rows ON CONFLICT (id) DO NOTHING;

-- Danh mục và hồ sơ nhân sự chung, không có khóa mỏ trong schema.
INSERT INTO catalog_materials (code, ma, name, ten, unit, dvt, density, ty_trong, price, gia, status)
SELECT 'DEMO2-PROD-' || n, 'DEMO2-PROD-' || n, (ARRAY['Đá 1x2','Đá 4x6','Đá mi','Đá base','Đá hộc'])[n],
       (ARRAY['Đá 1x2','Đá 4x6','Đá mi','Đá base','Đá hộc'])[n], 'tấn', 'tấn', 1.55, 1.55, 200000, 200000, 'active'
FROM generate_series(1,5) n ON CONFLICT (code) DO NOTHING;

INSERT INTO inventory_products (code, name, category, unit, density, sale_price, purchase_price, current_stock, status)
SELECT 'DEMO2-PROD-' || n, (ARRAY['Đá 1x2','Đá 4x6','Đá mi','Đá base','Đá hộc'])[n], 'Đá xây dựng', 'tấn', 1.55, 200000, 150000, 75, 'active'
FROM generate_series(1,5) n ON CONFLICT (code) DO NOTHING;

INSERT INTO hr_employees (id, code, name, department, job_position, join_date, contract_type, contracted_hours, status, certificates, base_salary)
SELECT 'DEMO2-EMP-' || n, 'DEMO2-EMP-' || n, 'Nhân viên demo ' || n, 'Vận hành', 'Nhân viên', '01/01/2026', 'Hợp đồng 1 - 3 năm', 8, 'Đang làm việc', '[]'::jsonb, 12000000
FROM generate_series(1,5) n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_contracts (id, code, employee_id, employee_name, department, contract_type, contracted_hours, start_date, end_date, base_salary, social_insurance_base, status)
SELECT 'DEMO2-HRC-' || n, 'DEMO2-HRC-' || n, e.id, e.name, e.department, e.contract_type, 8, '01/01/2026', '31/12/2026', 12000000, 12000000, 'Hiệu lực'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id = 'DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_attendances (id, employee_id, employee_name, department, job_position, date, check_in_time, check_out_time, contracted_hours, worked_hours, ot_hours, method, status)
SELECT 'DEMO2-ATT-' || n, e.id, e.name, e.department, e.job_position, '28/10/2026', '08:00', '16:00', 8, 8, 0, 'Thủ công', 'Đủ giờ công'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id = 'DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_leave_allocations (employee_id, employee_name, department, year, total_allocated, used_days, remaining_days)
SELECT e.id, e.name, e.department, 2026, 12, 0, 12 FROM generate_series(1,5) n
JOIN hr_employees e ON e.id = 'DEMO2-EMP-' || n
WHERE NOT EXISTS (SELECT 1 FROM hr_leave_allocations a WHERE a.employee_id = e.id AND a.year = 2026);

-- Giấy phép chỉ lưu tên mỏ dạng văn bản: không suy ra khóa mỏ.
INSERT INTO mining_permits (id, code, title, mine_name, category, category_label, issuer, license_number, issue_date, expiry_date, status, status_label, files, notes)
SELECT 'DEMO2-PERMIT-' || n, 'DEMO2-PERMIT-' || n, 'Giấy phép demo ' || n, 'Mỏ demo', 'mining', 'Khai thác', 'Cơ quan demo', 'DEMO2-GP-' || n, '01/01/2026', '31/12/2030', 'active', 'Còn hiệu lực', '[]'::jsonb, 'Dữ liệu minh họa'
FROM generate_series(1,5) n ON CONFLICT (id) DO NOTHING;

-- Chi tiết hợp đồng theo cha có quarry_code; mã chi tiết ổn định theo hợp đồng/sản phẩm.
INSERT INTO sales_contract_items (contract_id, contract_code, product_code, product_name, committed_tons, delivered_tons, unit_price, discount_pct)
SELECT c.id, c.code, 'DEMO2-PROD-' || d.n, p.name, 1000, 25, 200000, 0
FROM demo2_rows d JOIN sales_contracts c ON c.code = 'DEMO2-CONTRACT-' || d.site || '-' || lpad(d.n::text, 2, '0') AND c.quarry_code = d.quarry_code
JOIN inventory_products p ON p.code = 'DEMO2-PROD-' || d.n
WHERE NOT EXISTS (SELECT 1 FROM sales_contract_items i WHERE i.contract_code = c.code AND i.product_code = p.code);

INSERT INTO consolidated_delivery_orders (code, customer_name, contract_code, period_start, period_end, total_trips, total_tons, total_amount, vat_amount, grand_total, einvoice_status, status, ticket_codes)
SELECT 'DEMO2-DELIVERY-' || d.site || '-' || lpad(d.n::text, 2, '0'), c.customer_name, c.code, '01/10/2026', '31/10/2026', 1, 25, 5000000, 500000, 5500000, 'pending', 'draft',
       jsonb_build_array('DEMO2-TICKET-' || d.site || '-' || lpad(d.n::text, 2, '0'))
FROM demo2_rows d JOIN sales_contracts c ON c.code = 'DEMO2-CONTRACT-' || d.site || '-' || lpad(d.n::text, 2, '0') AND c.quarry_code = d.quarry_code
ON CONFLICT (code) DO NOTHING;


-- Các kỳ lương và bản ghi minh họa có tham chiếu nhân viên; không chứa tài khoản ngân hàng.
INSERT INTO hr_payroll_periods (id, code, month, status, is_locked, start_date, end_date)
SELECT 'DEMO2-PERIOD-' || n, 'DEMO2-PERIOD-' || n, '2026-0' || n, 'draft', false, '01/0' || n || '/2026', '28/0' || n || '/2026'
FROM generate_series(1,5) n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_payroll_records (id, period_id, employee_id, employee_name, department, contract_type, base_salary, actual_days, standard_days, working_salary, gross_salary, total_insurance, total_deductions, net_salary, status)
SELECT 'DEMO2-PAYROLL-' || n, 'DEMO2-PERIOD-' || n, e.id, e.name, e.department, e.contract_type, 12000000, 26, 26, 12000000, 12000000, 0, 0, 12000000, 'draft'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id = 'DEMO2-EMP-' || n
JOIN hr_payroll_periods p ON p.id = 'DEMO2-PERIOD-' || n ON CONFLICT (id) DO NOTHING;

-- Thanh toán là bảng chung, khóa code không unique nên chỉ chạy một phiên seed.
INSERT INTO payments_pricing (code, name, ten, price, gia, vat, unit, status)
SELECT 'DEMO2-PRICE-' || n, p.name, p.name, '200000', '200000', '10%', 'tấn', 'active'
FROM generate_series(1,5) n JOIN inventory_products p ON p.code = 'DEMO2-PROD-' || n
WHERE NOT EXISTS (SELECT 1 FROM payments_pricing x WHERE x.code = 'DEMO2-PRICE-' || n);

INSERT INTO payments_invoices (code, partner, customer, amount, balance, due, status)
SELECT 'DEMO2-INV-' || n, 'Khách hàng chung demo ' || n, 'Khách hàng chung demo ' || n, '5500000', '5500000', '30/11/2026', 'unpaid'
FROM generate_series(1,5) n WHERE NOT EXISTS (SELECT 1 FROM payments_invoices x WHERE x.code = 'DEMO2-INV-' || n);

INSERT INTO payments_debt (code, partner, customer, limit_val, balance, amount, due, status)
SELECT 'DEMO2-DEBT-' || n, 'Khách hàng chung demo ' || n, 'Khách hàng chung demo ' || n, '10000000', '5500000', '5500000', '30/11/2026', 'unpaid'
FROM generate_series(1,5) n WHERE NOT EXISTS (SELECT 1 FROM payments_debt x WHERE x.code = 'DEMO2-DEBT-' || n);

INSERT INTO payments_reconcile (code, customer, scale_amount, bank_amount, diff, status, date)
SELECT 'DEMO2-RECON-' || n, 'Khách hàng chung demo ' || n, '5500000', '5500000', '0', 'matched', '28/10/2026'
FROM generate_series(1,5) n WHERE NOT EXISTS (SELECT 1 FROM payments_reconcile x WHERE x.code = 'DEMO2-RECON-' || n);

INSERT INTO payments_prepaid (code, customer, initial_balance, current_balance, status, date)
SELECT 'DEMO2-PREPAID-' || n, 'Khách hàng chung demo ' || n, '10000000', '4500000', 'active', '28/10/2026'
FROM generate_series(1,5) n WHERE NOT EXISTS (SELECT 1 FROM payments_prepaid x WHERE x.code = 'DEMO2-PREPAID-' || n);


-- Hồ sơ nhân sự chung: chuỗi nhân viên, hợp đồng, nghỉ phép và đào tạo.
INSERT INTO hr_leave_requests (id, code, employee_id, employee_name, department, leave_type, start_date, end_date, total_days, total_hours, reason, stage)
SELECT 'DEMO2-LEAVE-' || n, 'DEMO2-LEAVE-' || n, e.id, e.name, e.department, 'Phép năm', '01/11/2026', '01/11/2026', 1, 8, 'Nghỉ phép demo', 'Chờ duyệt'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_employee_work_history (id, employee_id, employee_name, event_type, title, detail, changed_by, changed_at)
SELECT 'DEMO2-HISTORY-' || n, e.id, e.name, 'onboarding', 'Tiếp nhận demo', 'Hồ sơ minh họa', 'Demo', '01/01/2026'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_onboarding_checklists (id, employee_id, employee_name, task_name, task_description, is_completed, status)
SELECT 'DEMO2-ONBOARD-' || n, e.id, e.name, 'Hướng dẫn an toàn', 'Nội dung minh họa', false, 'pending'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_insurance_records (id, employee_id, employee_name, insurance_number, type, participation_date, social_insurance_base, status)
SELECT 'DEMO2-INS-' || n, e.id, e.name, 'DEMO2-INS-' || n, 'BHXH', '01/01/2026', 12000000, 'active'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_insurance_declarations (id, code, employee_id, employee_name, declaration_type, period, status, notes)
SELECT 'DEMO2-DECL-' || n, 'DEMO2-DECL-' || n, e.id, e.name, 'Tăng mới', '2026-01', 'draft', 'Minh họa'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_training_courses (id, code, name, category, provider, start_date, end_date, duration_hours, attendance_threshold, status)
SELECT 'DEMO2-COURSE-' || n, 'DEMO2-COURSE-' || n, 'Đào tạo an toàn demo ' || n, 'An toàn', 'Đơn vị demo', '01/10/2026', '01/10/2026', 8, 80, 'completed'
FROM generate_series(1,5) n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_training_sessions (id, course_id, session_name, session_date, start_time, end_time, instructor, status)
SELECT 'DEMO2-SESSION-' || n, c.id, 'Buổi đào tạo demo', '01/10/2026', '08:00', '16:00', 'Giảng viên demo', 'completed'
FROM generate_series(1,5) n JOIN hr_training_courses c ON c.id='DEMO2-COURSE-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_training_participants (id, course_id, employee_id, employee_name, attendance_rate, completion_status)
SELECT 'DEMO2-PARTICIPANT-' || n, c.id, e.id, e.name, 100, 'completed'
FROM generate_series(1,5) n JOIN hr_training_courses c ON c.id='DEMO2-COURSE-' || n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_training_attendances (id, session_id, participant_id, employee_id, employee_name, attendance_status)
SELECT 'DEMO2-TRAIN-ATT-' || n, s.id, p.id, e.id, e.name, 'present'
FROM generate_series(1,5) n JOIN hr_training_sessions s ON s.id='DEMO2-SESSION-' || n JOIN hr_training_participants p ON p.id='DEMO2-PARTICIPANT-' || n JOIN hr_employees e ON e.id=p.employee_id ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_training_evaluations (id, participant_id, course_id, employee_id, employee_name, criteria, rating, comment, submitted_at)
SELECT 'DEMO2-EVAL-' || n, p.id, p.course_id, e.id, e.name, 'Hoàn thành nội dung', 4, 'Đánh giá demo', '01/10/2026'
FROM generate_series(1,5) n JOIN hr_training_participants p ON p.id='DEMO2-PARTICIPANT-' || n JOIN hr_employees e ON e.id=p.employee_id ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_attendance_logs (id, employee_id, employee_name, date, check_in_time, check_out_time, work_value, status_code, overtime_hours, note)
SELECT 'DEMO2-LOG-' || n, e.id, e.name, '28/10/2026', '08:00', '16:00', 1, 'present', 0, 'Chấm công demo'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_timesheets (id, data)
SELECT 'DEMO2-TS-' || n, jsonb_build_array(jsonb_build_object('employeeId',e.id,'employeeCode',e.code,'employeeName',e.name,'department',e.department,'jobPosition',e.job_position,'contractedHoursPerDay',8,'standardDays',26,'actualDays',26,'totalContractHours',208,'totalWorkedHours',208,'paidLeaveDays',0,'unpaidLeaveDays',0,'otHours',0,'pieceworkTons',0,'qualifiedDays',26))
FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_timesheet_details (id, timesheet_id, employee_id, employee_name, date, work_value, status_code, overtime_hours, note)
SELECT 'DEMO2-DETAIL-' || n, 'DEMO2-TS-' || n, e.id, e.name, '28/10/2026', 1, 'present', 0, 'Công demo'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

INSERT INTO hr_salary_advances (id, code, employee_id, employee_name, amount, request_date, deduction_period, reason, status)
SELECT 'DEMO2-ADV-' || n, 'DEMO2-ADV-' || n, e.id, e.name, 0, '28/10/2026', '2026-10', 'Minh họa', 'draft'
FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;

-- Xác nhận giao hàng giữ liên hệ tới hợp đồng và phiếu cân cùng mỏ.
INSERT INTO delivery_confirmations (order_code, contract_code, customer_name, vehicle_plate, driver_name, product_type, quantity_ordered, quantity_delivered, quantity_confirmed, source_order, source_scale, source_warehouse, confirmation_status, notes)
SELECT o.code, c.code, c.customer_name, 'D2-' || d.site || '-' || lpad(d.n::text,2,'0'), 'Lái xe demo ' || d.n, p.name, 25, 25, 25, o.code, 'DEMO2-TICKET-' || d.site || '-' || lpad(d.n::text,2,'0'), 'DEMO2-OUT-' || d.site || '-' || lpad(d.n::text,2,'0'), 'pending', 'Dữ liệu demo'
FROM demo2_rows d JOIN sales_contracts c ON c.code='DEMO2-CONTRACT-' || d.site || '-' || lpad(d.n::text,2,'0') AND c.quarry_code=d.quarry_code JOIN consolidated_delivery_orders o ON o.code='DEMO2-DELIVERY-' || d.site || '-' || lpad(d.n::text,2,'0') JOIN inventory_products p ON p.code='DEMO2-PROD-' || d.n
WHERE NOT EXISTS (SELECT 1 FROM delivery_confirmations x WHERE x.order_code=o.code);

INSERT INTO accounting_entries (entry_code, entry_date, entry_type, account_code, account_name, description, debit_amount, credit_amount, balance, reference_type, reference_code, tax_rate, tax_amount, period, status)
SELECT 'DEMO2-ENTRY-' || n, '28/10/2026', 'sale', '131', 'Phải thu khách hàng', 'Bút toán demo', 5500000, 0, 5500000, 'invoice', 'DEMO2-INV-' || n, 10, 500000, '2026-10', 'draft'
FROM generate_series(1,5) n ON CONFLICT (entry_code) DO NOTHING;


INSERT INTO hr_recruitment_requests (id, code, title, department, position, headcount, reason, requested_by, approval_status)
SELECT 'DEMO2-REQ-' || n, 'DEMO2-REQ-' || n, 'Tuyển dụng demo ' || n, 'Vận hành', 'Nhân viên', 1, 'Minh họa', 'Demo', 'draft'
FROM generate_series(1,5) n ON CONFLICT (id) DO NOTHING;
INSERT INTO hr_recruitment_campaigns (id, code, name, request_id, position, department, open_date, recruit_deadline, target_headcount, hired_count, source, status)
SELECT 'DEMO2-CAMP-' || n, 'DEMO2-CAMP-' || n, 'Chiến dịch demo ' || n, r.id, 'Nhân viên', 'Vận hành', '01/10/2026', '30/11/2026', 1, 0, 'Demo', 'draft'
FROM generate_series(1,5) n JOIN hr_recruitment_requests r ON r.id='DEMO2-REQ-' || n ON CONFLICT (id) DO NOTHING;
INSERT INTO hr_candidates (id, code, name, position, source, pipeline_stage, list_status, note)
SELECT 'DEMO2-CAND-' || n, 'DEMO2-CAND-' || n, 'Ứng viên demo ' || n, 'Nhân viên', 'Demo', 'applied', 'active', 'Dữ liệu minh họa'
FROM generate_series(1,5) n ON CONFLICT (id) DO NOTHING;
INSERT INTO hr_candidate_applications (id, candidate_id, campaign_id, pipeline_stage, applied_at)
SELECT 'DEMO2-APPLICATION-' || n, c.id, r.id, 'applied', '01/10/2026' FROM generate_series(1,5) n JOIN hr_candidates c ON c.id='DEMO2-CAND-' || n JOIN hr_recruitment_campaigns r ON r.id='DEMO2-CAMP-' || n ON CONFLICT (id) DO NOTHING;
INSERT INTO hr_candidate_status_history (id, application_id, candidate_id, from_stage, to_stage, changed_by, note)
SELECT 'DEMO2-CAND-HIST-' || n, a.id, a.candidate_id, '', 'applied', 'Demo', 'Minh họa' FROM generate_series(1,5) n JOIN hr_candidate_applications a ON a.id='DEMO2-APPLICATION-' || n ON CONFLICT (id) DO NOTHING;
INSERT INTO hr_interview_schedules (id, candidate_id, candidate_name, position, interview_date, interview_time, interviewer, round, meeting_type, status)
SELECT 'DEMO2-INTERVIEW-' || n, c.id, c.name, c.position, '15/11/2026', '09:00', 'Nhân sự demo', 1, 'Trực tiếp', 'draft' FROM generate_series(1,5) n JOIN hr_candidates c ON c.id='DEMO2-CAND-' || n ON CONFLICT (id) DO NOTHING;
INSERT INTO hr_kpi_targets (id, employee_id, employee_name, period, kpi_code, kpi_name, target_value, actual_value, unit, weight_percent, score, status)
SELECT 'DEMO2-KPI-' || n,e.id,e.name,'2026-10','DEMO2-KPI-' || n,'Khối lượng demo',100,95,'tấn',100,95,'draft' FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;
INSERT INTO hr_okr_objectives (id, code, title, owner_text, period, level, progress_percent, status)
SELECT 'DEMO2-OBJ-' || n,'DEMO2-OBJ-' || n,'Mục tiêu demo ' || n,e.name,'2026-Q4','individual',95,'draft' FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;
INSERT INTO hr_okr_key_results (id, objective_id, key_result, start_value, target_value, current_value, unit, progress_percent)
SELECT 'DEMO2-KR-' || n,o.id,'Sản lượng demo',0,100,95,'tấn',95 FROM generate_series(1,5) n JOIN hr_okr_objectives o ON o.id='DEMO2-OBJ-' || n ON CONFLICT (id) DO NOTHING;
INSERT INTO hr_holidays (id, holiday_name, holiday_date, date_type, status)
SELECT 'DEMO2-HOLIDAY-' || n,'Ngày nghỉ demo ' || n,'0' || n || '/01/2026','holiday','draft' FROM generate_series(1,5) n ON CONFLICT (id) DO NOTHING;


-- Chứng từ mua là dữ liệu chung vì bảng không có khóa mỏ.
INSERT INTO purchase_vouchers (code, supplier_name, date, warehouse_loc, total_amount, vat_amount, grand_total, payment_status, status, notes)
SELECT 'DEMO2-PURCHASE-' || n, 'Nhà cung cấp chung demo ' || n, '28/10/2026', 'Kho demo', 3750000, 375000, 4125000, 'unpaid', 'draft', 'Dữ liệu demo'
FROM generate_series(1,5) n ON CONFLICT (code) DO NOTHING;
INSERT INTO purchase_voucher_items (voucher_id, voucher_code, product_code, product_name, unit, density, unit_price, quantity, total_amount, vat_rate, weight_ton)
SELECT v.id,v.code,p.code,p.name,'tấn',1.55,150000,25,3750000,10,25 FROM generate_series(1,5) n JOIN purchase_vouchers v ON v.code='DEMO2-PURCHASE-' || n JOIN inventory_products p ON p.code='DEMO2-PROD-' || n
WHERE NOT EXISTS (SELECT 1 FROM purchase_voucher_items i WHERE i.voucher_code=v.code AND i.product_code=p.code);
-- sales_vouchers chưa seed: một số DB có cột quarry_code ngoài migration hiện tại;
-- bỏ mặc định có thể gán nhầm phiếu Tân Uyên vào Phú Thọ.
INSERT INTO receipt_vouchers (code,partner_name,partner_type,payer_name,reason,ref_code,amount,payment_method,status,date,notes)
SELECT 'DEMO2-RECEIPT-' || n,'Khách hàng chung demo ' || n,'customer','Người nộp demo','Thu demo','DEMO2-INV-' || n,5500000,'cash','draft','28/10/2026','Minh họa' FROM generate_series(1,5) n ON CONFLICT (code) DO NOTHING;
INSERT INTO payment_vouchers (code,partner_name,partner_type,receiver_name,reason,ref_code,amount,payment_method,status,date,notes)
SELECT 'DEMO2-PAYMENT-' || n,v.supplier_name,'supplier','Người nhận demo','Chi demo',v.code,4125000,'cash','draft','28/10/2026','Minh họa' FROM generate_series(1,5) n JOIN purchase_vouchers v ON v.code='DEMO2-PURCHASE-' || n ON CONFLICT (code) DO NOTHING;


-- Danh mục thiết bị và lịch ca dùng chung, mã DEMO2 để không trùng dữ liệu khác.
INSERT INTO equipment (code,name,ten,loai,capacity,location,status)
SELECT 'DEMO2-EQUIPMENT-' || n,'Thiết bị demo ' || n,'Thiết bị demo ' || n,'Máy xúc','150 tấn/giờ','Kho demo','active' FROM generate_series(1,5) n ON CONFLICT (code) DO NOTHING;
INSERT INTO vehicle_catalogs (code,name,loai,tai_trong,unit,status)
SELECT 'DEMO2-VEHCAT-' || n,'Xe demo ' || n,'Xe ben',30,'tấn','active' FROM generate_series(1,5) n ON CONFLICT (code) DO NOTHING;
INSERT INTO ticket_types (code,name,ten,loai,description,status)
SELECT 'DEMO2-TYPE-' || n,'Loại phiếu demo ' || n,'Loại phiếu demo ' || n,'Cân bán hàng','Minh họa','active' FROM generate_series(1,5) n ON CONFLICT (code) DO NOTHING;
INSERT INTO stations (code,name,ten,location,capacity,type,status)
SELECT 'DEMO2-STATION-' || n,'Trạm demo ' || n,'Trạm demo ' || n,'Địa điểm demo','80 tấn','Trạm cân','inactive' FROM generate_series(1,5) n ON CONFLICT (code) DO NOTHING;
INSERT INTO hr_shifts (code,name,start_time,end_time,hours,status)
SELECT 'DEMO2-SHIFT-' || n,'Ca demo ' || n,'08:00','16:00',8,'Hoạt động' FROM generate_series(1,5) n ON CONFLICT (code) DO NOTHING;
INSERT INTO hr_shift_assignments (id,employee_id,employee_name,shift_code,shift_name,date,time_slot,status)
SELECT 'DEMO2-SHIFT-ASSIGN-' || n,e.id,e.name,s.code,s.name,'28/10/2026','08:00 - 16:00','draft' FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n JOIN hr_shifts s ON s.code='DEMO2-SHIFT-' || n ON CONFLICT (id) DO NOTHING;
INSERT INTO hr_shift_schedules (code,shift_id,employee_id,date,status)
SELECT 'DEMO2-SCHEDULE-' || n,s.code,e.id,'28/10/2026','draft' FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n JOIN hr_shifts s ON s.code='DEMO2-SHIFT-' || n WHERE NOT EXISTS (SELECT 1 FROM hr_shift_schedules x WHERE x.code='DEMO2-SCHEDULE-' || n);
INSERT INTO hr_decisions (id,code,employee_id,employee_name,decision_type,effective_date,content_detail,status,issued_by,notes)
SELECT 'DEMO2-DECISION-' || n,'DEMO2-DECISION-' || n,e.id,e.name,'Tiếp nhận','01/01/2026','{}'::jsonb,'draft','Demo','Minh họa' FROM generate_series(1,5) n JOIN hr_employees e ON e.id='DEMO2-EMP-' || n ON CONFLICT (id) DO NOTHING;
INSERT INTO cost_norms (norm_name,norm_type,unit_cost,unit,material_type,effective_date,status,description)
SELECT 'DEMO2-NORM-COST-' || n,'Sản xuất',10000,'tấn',p.name,'01/01/2026','inactive','Minh họa' FROM generate_series(1,5) n JOIN inventory_products p ON p.code='DEMO2-PROD-' || n WHERE NOT EXISTS (SELECT 1 FROM cost_norms x WHERE x.norm_name='DEMO2-NORM-COST-' || n);
INSERT INTO tax_records (tax_type,tax_code,period,taxable_amount,tax_rate,tax_amount,paid_amount,due_date,status,authority,notes)
SELECT 'Thuế tài nguyên','DEMO2-TAXREC-' || n,'2026-10',100000000,10,10000000,0,'30/11/2026','draft','Cơ quan demo','Minh họa' FROM generate_series(1,5) n WHERE NOT EXISTS (SELECT 1 FROM tax_records x WHERE x.tax_code='DEMO2-TAXREC-' || n);

COMMIT;
