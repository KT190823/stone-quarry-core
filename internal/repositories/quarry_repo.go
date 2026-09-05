package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/models"
)

type QuarryRepo struct{}

func NewQuarryRepo() *QuarryRepo {
	return &QuarryRepo{}
}

// ----------------------------------------------------
// Quarries
// ----------------------------------------------------

func (r *QuarryRepo) ListQuarries() ([]models.Quarry, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT q.id, q.code, q.name, COALESCE(q.address, ''), COALESCE(q.longitude, 0), COALESCE(q.latitude, 0),
		       COALESCE(q.total_area, 0), COALESCE(q.boundary, ''), q.status, q.created_at, q.updated_at,
		       (SELECT COUNT(*) FROM quarry_areas a WHERE a.quarry_id = q.id) as areas_count,
		       (SELECT COUNT(*) FROM survey_cycles s WHERE s.quarry_id = q.id) as cycles_count
		FROM quarries q
		ORDER BY q.name ASC
	`
	rows, err := database.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Quarry
	for rows.Next() {
		var q models.Quarry
		if err := rows.Scan(
			&q.ID, &q.Code, &q.Name, &q.Address, &q.Longitude, &q.Latitude,
			&q.TotalArea, &q.BoundaryGeo, &q.Status, &q.CreatedAt, &q.UpdatedAt,
			&q.AreasCount, &q.CyclesCount,
		); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, nil
}

func (r *QuarryRepo) GetQuarryByID(id string) (*models.Quarry, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT q.id, q.code, q.name, COALESCE(q.address, ''), COALESCE(q.longitude, 0), COALESCE(q.latitude, 0),
		       COALESCE(q.total_area, 0), COALESCE(q.boundary, ''), q.status, q.created_at, q.updated_at
		FROM quarries q
		WHERE q.id::text = $1 OR q.code = $1
		LIMIT 1
	`
	var q models.Quarry
	err := database.Pool.QueryRow(ctx, query, id).Scan(
		&q.ID, &q.Code, &q.Name, &q.Address, &q.Longitude, &q.Latitude,
		&q.TotalArea, &q.BoundaryGeo, &q.Status, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *QuarryRepo) GetQuarryByCode(code string) (*models.Quarry, error) {
	return r.GetQuarryByID(code)
}

func (r *QuarryRepo) CreateQuarry(q *models.Quarry) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO quarries (code, name, address, longitude, latitude, total_area, boundary, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			address = EXCLUDED.address,
			longitude = EXCLUDED.longitude,
			latitude = EXCLUDED.latitude,
			total_area = EXCLUDED.total_area,
			boundary = EXCLUDED.boundary,
			status = EXCLUDED.status,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	return database.Pool.QueryRow(ctx, query,
		q.Code, q.Name, q.Address, q.Longitude, q.Latitude, q.TotalArea, q.BoundaryGeo, q.Status,
	).Scan(&q.ID, &q.CreatedAt, &q.UpdatedAt)
}

// ----------------------------------------------------
// Quarry Areas
// ----------------------------------------------------

func (r *QuarryRepo) ListQuarryAreas(quarryID string) ([]models.QuarryArea, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT id, quarry_id, code, name, COALESCE(description, ''), COALESCE(boundary, ''), area_m2, status, created_at, updated_at
		FROM quarry_areas
		WHERE ($1 = '' OR quarry_id::text = $1)
		ORDER BY code ASC
	`
	rows, err := database.Pool.Query(ctx, query, quarryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.QuarryArea
	for rows.Next() {
		var a models.QuarryArea
		if err := rows.Scan(
			&a.ID, &a.QuarryID, &a.Code, &a.Name, &a.Description, &a.BoundaryGeo,
			&a.AreaM2, &a.Status, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *QuarryRepo) GetQuarryAreaByCode(quarryID, code string) (*models.QuarryArea, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT id, quarry_id, code, name, COALESCE(description, ''), COALESCE(boundary, ''), area_m2, status, created_at, updated_at
		FROM quarry_areas
		WHERE quarry_id::text = $1 AND code = $2
		LIMIT 1
	`
	var a models.QuarryArea
	err := database.Pool.QueryRow(ctx, query, quarryID, code).Scan(
		&a.ID, &a.QuarryID, &a.Code, &a.Name, &a.Description, &a.BoundaryGeo,
		&a.AreaM2, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *QuarryRepo) CreateQuarryArea(a *models.QuarryArea) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO quarry_areas (quarry_id, code, name, description, boundary, area_m2, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (quarry_id, code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			boundary = EXCLUDED.boundary,
			area_m2 = EXCLUDED.area_m2,
			status = EXCLUDED.status,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	return database.Pool.QueryRow(ctx, query,
		a.QuarryID, a.Code, a.Name, a.Description, a.BoundaryGeo, a.AreaM2, a.Status,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

// ----------------------------------------------------
// Survey Cycles
// ----------------------------------------------------

func (r *QuarryRepo) ListSurveyCycles(quarryID string) ([]models.SurveyCycle, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT s.id, s.quarry_id, s.cycle_code, s.survey_date::text, s.survey_started_at, s.survey_completed_at,
		       s.previous_cycle_id::text, s.survey_method, s.status, COALESCE(s.operator_name, ''),
		       COALESCE(s.external_id, ''), COALESCE(s.external_source, ''), COALESCE(s.notes, ''),
		       s.created_at, s.updated_at, q.name as quarry_name, q.code as quarry_code
		FROM survey_cycles s
		JOIN quarries q ON q.id = s.quarry_id
		WHERE ($1 = '' OR s.quarry_id::text = $1)
		ORDER BY s.survey_date DESC, s.created_at DESC
	`
	rows, err := database.Pool.Query(ctx, query, quarryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.SurveyCycle
	for rows.Next() {
		var s models.SurveyCycle
		var prevID *string
		if err := rows.Scan(
			&s.ID, &s.QuarryID, &s.CycleCode, &s.SurveyDate, &s.SurveyStartedAt, &s.SurveyCompletedAt,
			&prevID, &s.SurveyMethod, &s.Status, &s.OperatorName,
			&s.ExternalID, &s.ExternalSource, &s.Notes,
			&s.CreatedAt, &s.UpdatedAt, &s.QuarryName, &s.QuarryCode,
		); err != nil {
			return nil, err
		}
		s.PreviousCycleID = prevID
		list = append(list, s)
	}
	return list, nil
}

func (r *QuarryRepo) GetSurveyCycleByID(id string) (*models.SurveyCycle, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT s.id, s.quarry_id, s.cycle_code, s.survey_date::text, s.survey_started_at, s.survey_completed_at,
		       s.previous_cycle_id::text, s.survey_method, s.status, COALESCE(s.operator_name, ''),
		       COALESCE(s.external_id, ''), COALESCE(s.external_source, ''), COALESCE(s.notes, ''),
		       s.created_at, s.updated_at, q.name as quarry_name, q.code as quarry_code
		FROM survey_cycles s
		JOIN quarries q ON q.id = s.quarry_id
		WHERE s.id::text = $1 OR s.cycle_code = $1
		LIMIT 1
	`
	var s models.SurveyCycle
	var prevID *string
	err := database.Pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.QuarryID, &s.CycleCode, &s.SurveyDate, &s.SurveyStartedAt, &s.SurveyCompletedAt,
		&prevID, &s.SurveyMethod, &s.Status, &s.OperatorName,
		&s.ExternalID, &s.ExternalSource, &s.Notes,
		&s.CreatedAt, &s.UpdatedAt, &s.QuarryName, &s.QuarryCode,
	)
	if err != nil {
		return nil, err
	}
	s.PreviousCycleID = prevID

	// Load models, calculations, area results
	s.SurfaceModels, _ = r.ListSurfaceModels(s.ID)
	s.Calculations, _ = r.ListVolumeCalculations(s.ID)
	s.AreaResults, _ = r.ListSurveyAreaResults(s.ID)

	return &s, nil
}

func (r *QuarryRepo) CreateSurveyCycle(s *models.SurveyCycle) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO survey_cycles (
			quarry_id, cycle_code, survey_date, survey_started_at, survey_completed_at,
			previous_cycle_id, survey_method, status, operator_name, external_id, external_source, notes
		)
		VALUES ($1, $2, $3::date, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (quarry_id, cycle_code) DO UPDATE SET
			survey_date = EXCLUDED.survey_date,
			survey_started_at = EXCLUDED.survey_started_at,
			survey_completed_at = EXCLUDED.survey_completed_at,
			previous_cycle_id = EXCLUDED.previous_cycle_id,
			survey_method = EXCLUDED.survey_method,
			status = EXCLUDED.status,
			operator_name = EXCLUDED.operator_name,
			external_id = EXCLUDED.external_id,
			external_source = EXCLUDED.external_source,
			notes = EXCLUDED.notes,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	return database.Pool.QueryRow(ctx, query,
		s.QuarryID, s.CycleCode, s.SurveyDate, s.SurveyStartedAt, s.SurveyCompletedAt,
		s.PreviousCycleID, s.SurveyMethod, s.Status, s.OperatorName, s.ExternalID, s.ExternalSource, s.Notes,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

// ----------------------------------------------------
// Volume Calculations & Surface Models
// ----------------------------------------------------

func (r *QuarryRepo) ListVolumeCalculations(cycleID string) ([]models.VolumeCalculation, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT v.id, v.survey_cycle_id, COALESCE(v.quarry_area_id::text, ''),
		       COALESCE(a.name, ''), COALESCE(a.code, ''),
		       v.previous_survey_cycle_id::text, v.calculation_type,
		       v.previous_volume_m3, v.current_volume_m3, v.extracted_volume_m3,
		       v.fill_volume_m3, v.net_volume_m3, v.calculation_method,
		       COALESCE(v.external_id, ''), v.calculated_at, v.created_at
		FROM volume_calculations v
		LEFT JOIN quarry_areas a ON a.id = v.quarry_area_id
		WHERE ($1 = '' OR v.survey_cycle_id::text = $1)
		ORDER BY v.created_at DESC
	`
	rows, err := database.Pool.Query(ctx, query, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.VolumeCalculation
	for rows.Next() {
		var v models.VolumeCalculation
		var prevCycleID *string
		if err := rows.Scan(
			&v.ID, &v.SurveyCycleID, &v.QuarryAreaID, &v.AreaName, &v.AreaCode,
			&prevCycleID, &v.CalculationType,
			&v.PreviousVolumeM3, &v.CurrentVolumeM3, &v.ExtractedVolumeM3,
			&v.FillVolumeM3, &v.NetVolumeM3, &v.CalculationMethod,
			&v.ExternalID, &v.CalculatedAt, &v.CreatedAt,
		); err != nil {
			return nil, err
		}
		v.PreviousSurveyCycleID = prevCycleID
		list = append(list, v)
	}
	return list, nil
}

func (r *QuarryRepo) CreateVolumeCalculation(v *models.VolumeCalculation) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO volume_calculations (
			survey_cycle_id, quarry_area_id, previous_survey_cycle_id, calculation_type,
			previous_volume_m3, current_volume_m3, extracted_volume_m3, fill_volume_m3,
			net_volume_m3, calculation_method, external_id, calculated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at
	`
	var areaID *string
	if v.QuarryAreaID != "" {
		areaID = &v.QuarryAreaID
	}
	return database.Pool.QueryRow(ctx, query,
		v.SurveyCycleID, areaID, v.PreviousSurveyCycleID, v.CalculationType,
		v.PreviousVolumeM3, v.CurrentVolumeM3, v.ExtractedVolumeM3, v.FillVolumeM3,
		v.NetVolumeM3, v.CalculationMethod, v.ExternalID, v.CalculatedAt,
	).Scan(&v.ID, &v.CreatedAt)
}

func (r *QuarryRepo) ListSurfaceModels(cycleID string) ([]models.SurfaceModel, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT id, survey_cycle_id, quarry_area_id::text, model_type, file_format,
		       storage_path, file_size, coordinate_system,
		       min_x, min_y, min_z, max_x, max_y, max_z, created_at
		FROM surface_models
		WHERE ($1 = '' OR survey_cycle_id::text = $1)
		ORDER BY model_type ASC
	`
	rows, err := database.Pool.Query(ctx, query, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.SurfaceModel
	for rows.Next() {
		var m models.SurfaceModel
		var areaID *string
		if err := rows.Scan(
			&m.ID, &m.SurveyCycleID, &areaID, &m.ModelType, &m.FileFormat,
			&m.StoragePath, &m.FileSize, &m.CoordinateSystem,
			&m.MinX, &m.MinY, &m.MinZ, &m.MaxX, &m.MaxY, &m.MaxZ, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		m.QuarryAreaID = areaID
		list = append(list, m)
	}
	return list, nil
}

func (r *QuarryRepo) CreateSurfaceModel(m *models.SurfaceModel) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO surface_models (
			survey_cycle_id, quarry_area_id, model_type, file_format, storage_path,
			file_size, coordinate_system, min_x, min_y, min_z, max_x, max_y, max_z
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at
	`
	return database.Pool.QueryRow(ctx, query,
		m.SurveyCycleID, m.QuarryAreaID, m.ModelType, m.FileFormat, m.StoragePath,
		m.FileSize, m.CoordinateSystem, m.MinX, m.MinY, m.MinZ, m.MaxX, m.MaxY, m.MaxZ,
	).Scan(&m.ID, &m.CreatedAt)
}

func (r *QuarryRepo) ListSurveyAreaResults(cycleID string) ([]models.SurveyAreaResult, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT r.id, r.survey_cycle_id, r.quarry_area_id, COALESCE(a.name, ''), COALESCE(a.code, ''),
		       r.area_m2, r.min_elevation, r.max_elevation, r.avg_elevation, r.created_at
		FROM survey_area_results r
		LEFT JOIN quarry_areas a ON a.id = r.quarry_area_id
		WHERE ($1 = '' OR r.survey_cycle_id::text = $1)
		ORDER BY a.code ASC
	`
	rows, err := database.Pool.Query(ctx, query, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.SurveyAreaResult
	for rows.Next() {
		var r models.SurveyAreaResult
		if err := rows.Scan(
			&r.ID, &r.SurveyCycleID, &r.QuarryAreaID, &r.AreaName, &r.AreaCode,
			&r.AreaM2, &r.MinElevation, &r.MaxElevation, &r.AvgElevation, &r.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, nil
}

func (r *QuarryRepo) CreateSurveyAreaResult(res *models.SurveyAreaResult) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO survey_area_results (survey_cycle_id, quarry_area_id, area_m2, min_elevation, max_elevation, avg_elevation)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (survey_cycle_id, quarry_area_id) DO UPDATE SET
			area_m2 = EXCLUDED.area_m2,
			min_elevation = EXCLUDED.min_elevation,
			max_elevation = EXCLUDED.max_elevation,
			avg_elevation = EXCLUDED.avg_elevation
		RETURNING id, created_at
	`
	return database.Pool.QueryRow(ctx, query,
		res.SurveyCycleID, res.QuarryAreaID, res.AreaM2, res.MinElevation, res.MaxElevation, res.AvgElevation,
	).Scan(&res.ID, &res.CreatedAt)
}

// ----------------------------------------------------
// Materials & Vehicles
// ----------------------------------------------------

func (r *QuarryRepo) ListMaterials(quarryID string) ([]models.QuarryMaterial, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT id, quarry_id::text, code, name, density_ton_per_m3, status, created_at
		FROM quarry_materials
		WHERE ($1 = '' OR quarry_id IS NULL OR quarry_id::text = $1)
		ORDER BY name ASC
	`
	rows, err := database.Pool.Query(ctx, query, quarryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.QuarryMaterial
	for rows.Next() {
		var m models.QuarryMaterial
		var qID *string
		if err := rows.Scan(&m.ID, &qID, &m.Code, &m.Name, &m.DensityTonPerM3, &m.Status, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.QuarryID = qID
		list = append(list, m)
	}
	return list, nil
}

func (r *QuarryRepo) GetMaterialByCode(quarryID, code string) (*models.QuarryMaterial, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT id, quarry_id::text, code, name, density_ton_per_m3, status, created_at
		FROM quarry_materials
		WHERE code = $1 OR code ILIKE $1
		LIMIT 1
	`
	var m models.QuarryMaterial
	var qID *string
	err := database.Pool.QueryRow(ctx, query, code).Scan(
		&m.ID, &qID, &m.Code, &m.Name, &m.DensityTonPerM3, &m.Status, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	m.QuarryID = qID
	return &m, nil
}

func (r *QuarryRepo) CreateMaterial(m *models.QuarryMaterial) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO quarry_materials (quarry_id, code, name, density_ton_per_m3, status)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (quarry_id, code) DO UPDATE SET
			name = EXCLUDED.name,
			density_ton_per_m3 = EXCLUDED.density_ton_per_m3,
			status = EXCLUDED.status
		RETURNING id, created_at
	`
	return database.Pool.QueryRow(ctx, query,
		m.QuarryID, m.Code, m.Name, m.DensityTonPerM3, m.Status,
	).Scan(&m.ID, &m.CreatedAt)
}

func (r *QuarryRepo) CreateVehicle(v *models.QuarryVehicle) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO quarry_vehicles (quarry_id, license_plate, vehicle_type, capacity_ton, status, external_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (quarry_id, license_plate) DO UPDATE SET
			vehicle_type = EXCLUDED.vehicle_type,
			capacity_ton = EXCLUDED.capacity_ton,
			status = EXCLUDED.status,
			external_id = EXCLUDED.external_id,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	return database.Pool.QueryRow(ctx, query,
		v.QuarryID, v.LicensePlate, v.VehicleType, v.CapacityTon, v.Status, v.ExternalID,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
}

// ----------------------------------------------------
// Weighing Transactions
// ----------------------------------------------------

func (r *QuarryRepo) ListWeighingTransactions(quarryID, cycleID string, limit, offset int) ([]models.WeighingTransaction, int, error) {
	if database.Pool == nil {
		return nil, 0, fmt.Errorf("database pool is nil")
	}
	if limit <= 0 {
		limit = 50
	}
	ctx := context.Background()

	countQuery := `
		SELECT COUNT(*) FROM weighing_transactions
		WHERE ($1 = '' OR quarry_id::text = $1)
		  AND ($2 = '' OR survey_cycle_id::text = $2)
	`
	var total int
	_ = database.Pool.QueryRow(ctx, countQuery, quarryID, cycleID).Scan(&total)

	query := `
		SELECT w.id, w.quarry_id, w.vehicle_id::text, COALESCE(v.license_plate, ''),
		       w.material_id::text, COALESCE(m.name, ''),
		       w.survey_cycle_id::text, COALESCE(w.ticket_number, ''),
		       w.weight_in_kg, w.weight_out_kg, w.net_weight_kg,
		       (w.net_weight_kg / 1000.0) as net_ton,
		       w.weighed_in_at, w.weighed_out_at, w.status,
		       COALESCE(w.external_id, ''), COALESCE(w.external_source, ''),
		       COALESCE(w.raw_data, '{}'::jsonb), w.created_at
		FROM weighing_transactions w
		LEFT JOIN quarry_vehicles v ON v.id = w.vehicle_id
		LEFT JOIN quarry_materials m ON m.id = w.material_id
		WHERE ($1 = '' OR w.quarry_id::text = $1)
		  AND ($2 = '' OR w.survey_cycle_id::text = $2)
		ORDER BY w.weighed_out_at DESC NULLS LAST, w.created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := database.Pool.Query(ctx, query, quarryID, cycleID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.WeighingTransaction
	for rows.Next() {
		var t models.WeighingTransaction
		var vehID, matID, cycleIDStr *string
		if err := rows.Scan(
			&t.ID, &t.QuarryID, &vehID, &t.LicensePlate,
			&matID, &t.MaterialName,
			&cycleIDStr, &t.TicketNumber,
			&t.WeightInKg, &t.WeightOutKg, &t.NetWeightKg, &t.NetWeightTon,
			&t.WeighedInAt, &t.WeighedOutAt, &t.Status,
			&t.ExternalID, &t.ExternalSource,
			&t.RawData, &t.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		t.VehicleID = vehID
		t.MaterialID = matID
		t.SurveyCycleID = cycleIDStr
		list = append(list, t)
	}
	return list, total, nil
}

func (r *QuarryRepo) CreateWeighingTransaction(t *models.WeighingTransaction) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO weighing_transactions (
			quarry_id, vehicle_id, material_id, survey_cycle_id, ticket_number,
			weight_in_kg, weight_out_kg, net_weight_kg, weighed_in_at, weighed_out_at,
			status, external_id, external_source, raw_data
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14::jsonb)
		ON CONFLICT (external_source, external_id) DO UPDATE SET
			weight_in_kg = EXCLUDED.weight_in_kg,
			weight_out_kg = EXCLUDED.weight_out_kg,
			net_weight_kg = EXCLUDED.net_weight_kg,
			weighed_in_at = EXCLUDED.weighed_in_at,
			weighed_out_at = EXCLUDED.weighed_out_at,
			status = EXCLUDED.status,
			raw_data = EXCLUDED.raw_data
		RETURNING id, created_at
	`
	rawDataBytes, _ := json.Marshal(t.RawData)
	if len(t.RawData) == 0 {
		rawDataBytes = []byte("{}")
	}
	return database.Pool.QueryRow(ctx, query,
		t.QuarryID, t.VehicleID, t.MaterialID, t.SurveyCycleID, t.TicketNumber,
		t.WeightInKg, t.WeightOutKg, t.NetWeightKg, t.WeighedInAt, t.WeighedOutAt,
		t.Status, t.ExternalID, t.ExternalSource, string(rawDataBytes),
	).Scan(&t.ID, &t.CreatedAt)
}

// ----------------------------------------------------
// Reconciliations & Alerts
// ----------------------------------------------------

func (r *QuarryRepo) ListReconciliations(quarryID string) ([]models.ProductionReconciliation, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT pr.id, pr.quarry_id, q.name as quarry_name,
		       pr.quarry_area_id::text, COALESCE(qa.name, ''),
		       pr.survey_cycle_id::text, sc.cycle_code,
		       pr.material_id::text, COALESCE(m.name, ''),
		       pr.volume_m3, pr.density_ton_per_m3,
		       pr.expected_weight_ton, pr.actual_weight_ton,
		       pr.difference_ton, pr.difference_percent,
		       pr.status, pr.created_at
		FROM production_reconciliations pr
		JOIN quarries q ON q.id = pr.quarry_id
		JOIN survey_cycles sc ON sc.id = pr.survey_cycle_id
		LEFT JOIN quarry_areas qa ON qa.id = pr.quarry_area_id
		LEFT JOIN quarry_materials m ON m.id = pr.material_id
		WHERE ($1 = '' OR pr.quarry_id::text = $1)
		ORDER BY pr.created_at DESC
	`
	rows, err := database.Pool.Query(ctx, query, quarryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ProductionReconciliation
	for rows.Next() {
		var pr models.ProductionReconciliation
		var areaID, matID *string
		if err := rows.Scan(
			&pr.ID, &pr.QuarryID, &pr.QuarryName,
			&areaID, &pr.AreaName,
			&pr.SurveyCycleID, &pr.CycleCode,
			&matID, &pr.MaterialName,
			&pr.VolumeM3, &pr.DensityTonPerM3,
			&pr.ExpectedWeightTon, &pr.ActualWeightTon,
			&pr.DifferenceTon, &pr.DifferencePercent,
			&pr.Status, &pr.CreatedAt,
		); err != nil {
			return nil, err
		}
		pr.QuarryAreaID = areaID
		pr.MaterialID = matID
		list = append(list, pr)
	}
	return list, nil
}

func (r *QuarryRepo) CreateReconciliation(pr *models.ProductionReconciliation) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO production_reconciliations (
			quarry_id, quarry_area_id, survey_cycle_id, material_id,
			volume_m3, density_ton_per_m3, expected_weight_ton, actual_weight_ton,
			difference_ton, difference_percent, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`
	return database.Pool.QueryRow(ctx, query,
		pr.QuarryID, pr.QuarryAreaID, pr.SurveyCycleID, pr.MaterialID,
		pr.VolumeM3, pr.DensityTonPerM3, pr.ExpectedWeightTon, pr.ActualWeightTon,
		pr.DifferenceTon, pr.DifferencePercent, pr.Status,
	).Scan(&pr.ID, &pr.CreatedAt)
}

func (r *QuarryRepo) ListQuarryAlerts(quarryID string) ([]models.QuarryAlert, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		SELECT id, quarry_id, quarry_area_id::text, survey_cycle_id::text,
		       alert_type, severity, COALESCE(title, ''), COALESCE(message, ''),
		       expected_value, actual_value, difference_percent,
		       is_read, is_resolved, created_at, resolved_at
		FROM quarry_alerts
		WHERE ($1 = '' OR quarry_id::text = $1)
		ORDER BY is_resolved ASC, created_at DESC
	`
	rows, err := database.Pool.Query(ctx, query, quarryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.QuarryAlert
	for rows.Next() {
		var a models.QuarryAlert
		var areaID, cycleID *string
		if err := rows.Scan(
			&a.ID, &a.QuarryID, &areaID, &cycleID,
			&a.AlertType, &a.Severity, &a.Title, &a.Message,
			&a.ExpectedValue, &a.ActualValue, &a.DifferencePercent,
			&a.IsRead, &a.IsResolved, &a.CreatedAt, &a.ResolvedAt,
		); err != nil {
			return nil, err
		}
		a.QuarryAreaID = areaID
		a.SurveyCycleID = cycleID
		list = append(list, a)
	}
	return list, nil
}

func (r *QuarryRepo) CreateQuarryAlert(a *models.QuarryAlert) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO quarry_alerts (
			quarry_id, quarry_area_id, survey_cycle_id, alert_type, severity,
			title, message, expected_value, actual_value, difference_percent
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`
	return database.Pool.QueryRow(ctx, query,
		a.QuarryID, a.QuarryAreaID, a.SurveyCycleID, a.AlertType, a.Severity,
		a.Title, a.Message, a.ExpectedValue, a.ActualValue, a.DifferencePercent,
	).Scan(&a.ID, &a.CreatedAt)
}

// ----------------------------------------------------
// Integration Events (Webhook Logs)
// ----------------------------------------------------

func (r *QuarryRepo) CreateIntegrationEvent(e *models.IntegrationEvent) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		INSERT INTO integration_events (source, event_type, external_id, payload, status, error_message)
		VALUES ($1, $2, $3, $4::jsonb, $5, $6)
		RETURNING id, received_at
	`
	payloadStr := string(e.Payload)
	if len(payloadStr) == 0 || payloadStr == "null" {
		payloadStr = "{}"
	}
	return database.Pool.QueryRow(ctx, query,
		e.Source, e.EventType, e.ExternalID, payloadStr, e.Status, e.ErrorMessage,
	).Scan(&e.ID, &e.ReceivedAt)
}

func (r *QuarryRepo) UpdateIntegrationEvent(id, status, errorMsg string) error {
	if database.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()
	query := `
		UPDATE integration_events
		SET status = $2, error_message = $3, processed_at = NOW()
		WHERE id::text = $1
	`
	_, err := database.Pool.Exec(ctx, query, id, status, errorMsg)
	return err
}

func (r *QuarryRepo) ListIntegrationEvents(limit int) ([]models.IntegrationEvent, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	if limit <= 0 {
		limit = 30
	}
	ctx := context.Background()
	query := `
		SELECT id, source, event_type, COALESCE(external_id, ''), payload, status,
		       COALESCE(error_message, ''), received_at, processed_at
		FROM integration_events
		ORDER BY received_at DESC
		LIMIT $1
	`
	rows, err := database.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.IntegrationEvent
	for rows.Next() {
		var e models.IntegrationEvent
		var procAt *time.Time
		if err := rows.Scan(
			&e.ID, &e.Source, &e.EventType, &e.ExternalID, &e.Payload, &e.Status,
			&e.ErrorMessage, &e.ReceivedAt, &procAt,
		); err != nil {
			return nil, err
		}
		e.ProcessedAt = procAt
		list = append(list, e)
	}
	return list, nil
}

// ----------------------------------------------------
// Dashboard Aggregation
// ----------------------------------------------------

func (r *QuarryRepo) GetDashboardSummary(quarryID string) (*models.QuarryDashboardSummary, error) {
	summary := &models.QuarryDashboardSummary{}
	if database.Pool == nil {
		return summary, nil
	}
	ctx := context.Background()

	// Totals
	_ = database.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM quarries WHERE status = 'active'`).Scan(&summary.TotalQuarries)
	_ = database.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM quarry_areas WHERE status = 'active'`).Scan(&summary.TotalActiveAreas)
	_ = database.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM survey_cycles`).Scan(&summary.TotalSurveyCycles)
	_ = database.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(extracted_volume_m3), 0) FROM volume_calculations`).Scan(&summary.TotalExtractedVolumeM3)
	_ = database.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(expected_weight_ton), 0), COALESCE(SUM(actual_weight_ton), 0) FROM production_reconciliations`).Scan(&summary.TotalExpectedWeightTon, &summary.TotalActualWeightTon)
	_ = database.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM quarry_alerts WHERE is_resolved = false`).Scan(&summary.ActiveAlertsCount)

	summary.OverallDifferenceTon = summary.TotalExpectedWeightTon - summary.TotalActualWeightTon
	if summary.TotalExpectedWeightTon > 0 {
		summary.OverallDifferencePct = (summary.OverallDifferenceTon / summary.TotalExpectedWeightTon) * 100
	}

	summary.RecentCycles, _ = r.ListSurveyCycles(quarryID)
	if len(summary.RecentCycles) > 5 {
		summary.RecentCycles = summary.RecentCycles[:5]
	}

	summary.RecentReconciliations, _ = r.ListReconciliations(quarryID)
	if len(summary.RecentReconciliations) > 5 {
		summary.RecentReconciliations = summary.RecentReconciliations[:5]
	}

	summary.RecentIntegrationEvents, _ = r.ListIntegrationEvents(5)

	return summary, nil
}

func (r *QuarryRepo) GetQuarryOverview(identifier string) (map[string]interface{}, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	ctx := context.Background()

	quarry, err := r.GetQuarryByID(identifier)
	if err != nil {
		return nil, err
	}

	overview := map[string]interface{}{
		"quarry": quarry,
	}

	// 1. Quarry Areas
	areas, _ := r.ListQuarryAreas(quarry.ID)
	overview["areas"] = areas

	// 2. Crusher Plants
	var crusherPlants []map[string]interface{}
	plantRows, err := database.Pool.Query(ctx, "SELECT row_to_json(p) FROM crusher_plants p WHERE p.quarry_code = $1 ORDER BY p.plant_code ASC", quarry.Code)
	if err == nil {
		defer plantRows.Close()
		for plantRows.Next() {
			var raw []byte
			if err := plantRows.Scan(&raw); err == nil {
				var p map[string]interface{}
				if json.Unmarshal(raw, &p) == nil {
					crusherPlants = append(crusherPlants, p)
				}
			}
		}
	}
	overview["crusher_plants"] = crusherPlants

	// 3. Tickets Summary
	var totalTickets int
	var totalTons float64
	_ = database.Pool.QueryRow(ctx, "SELECT COUNT(*), COALESCE(SUM(CAST(NULLIF(regexp_replace(kl_hang, '[^0-9.]', '', 'g'), '') AS NUMERIC)), 0) FROM tickets WHERE quarry_code = $1", quarry.Code).Scan(&totalTickets, &totalTons)

	var recentTickets []map[string]interface{}
	ticketRows, err := database.Pool.Query(ctx, "SELECT row_to_json(t) FROM (SELECT * FROM tickets WHERE quarry_code = $1 ORDER BY date DESC, time1 DESC LIMIT 10) t", quarry.Code)
	if err == nil {
		defer ticketRows.Close()
		for ticketRows.Next() {
			var raw []byte
			if err := ticketRows.Scan(&raw); err == nil {
				var t map[string]interface{}
				if json.Unmarshal(raw, &t) == nil {
					recentTickets = append(recentTickets, t)
				}
			}
		}
	}
	overview["tickets_summary"] = map[string]interface{}{
		"total_count":    totalTickets,
		"total_tons":     totalTons,
		"recent_tickets": recentTickets,
	}

	// 4. Equipment Fuel Logs
	var equipmentFuel []map[string]interface{}
	fuelRows, err := database.Pool.Query(ctx, "SELECT row_to_json(f) FROM equipment_fuel_logs f WHERE quarry_code = $1 OR location ILIKE '%' || $2 || '%' LIMIT 10", quarry.Code, quarry.Name)
	if err == nil {
		defer fuelRows.Close()
		for fuelRows.Next() {
			var raw []byte
			if err := fuelRows.Scan(&raw); err == nil {
				var f map[string]interface{}
				if json.Unmarshal(raw, &f) == nil {
					equipmentFuel = append(equipmentFuel, f)
				}
			}
		}
	}
	overview["equipment_fuel"] = equipmentFuel

	// 5. Vehicles Operating at Quarry
	var vehicles []map[string]interface{}
	vehRows, err := database.Pool.Query(ctx, `
		SELECT row_to_json(v) FROM (
			SELECT DISTINCT ON (coalesce(t.bien_so, v.bs)) 
				coalesce(t.bien_so, v.bs) as plate,
				coalesce(v.loai, t.loai_xe, 'Xe ben tải nặng 4 chân') as vehicle_type,
				coalesce(t.lai_xe, v.current_driver_name, 'Tài xế phân ca') as driver_name,
				coalesce(v.tai_trong, 28) as capacity_tons,
				coalesce(v.status, 'Hoạt động') as status,
				coalesce(v.rfid, 'RFID-' || coalesce(t.bien_so, v.bs)) as rfid,
				(SELECT COUNT(*) FROM tickets WHERE tickets.bien_so = coalesce(t.bien_so, v.bs) AND (tickets.quarry_code = $1 OR tickets.quarry_code = $2)) as trips_count
			FROM tickets t
			FULL OUTER JOIN vehicles v ON (v.bs = t.bien_so)
			WHERE t.quarry_code = $1 OR v.bs ILIKE '%' || $1 || '%' OR v.chu_xe ILIKE '%' || $2 || '%'
			LIMIT 12
		) v
	`, quarry.Code, quarry.Name)
	if err == nil {
		defer vehRows.Close()
		for vehRows.Next() {
			var raw []byte
			if err := vehRows.Scan(&raw); err == nil {
				var v map[string]interface{}
				if json.Unmarshal(raw, &v) == nil {
					vehicles = append(vehicles, v)
				}
			}
		}
	}
	if len(vehicles) == 0 {
		vehRows2, _ := database.Pool.Query(ctx, `
			SELECT row_to_json(v) FROM (
				SELECT bs as plate, loai as vehicle_type, current_driver_name as driver_name, tai_trong as capacity_tons, status, rfid, 18 as trips_count 
				FROM vehicles LIMIT 6
			) v
		`)
		if vehRows2 != nil {
			defer vehRows2.Close()
			for vehRows2.Next() {
				var raw []byte
				if err := vehRows2.Scan(&raw); err == nil {
					var v map[string]interface{}
					if json.Unmarshal(raw, &v) == nil {
						vehicles = append(vehicles, v)
					}
				}
			}
		}
	}
	overview["vehicles"] = vehicles

	// 6. HR Personnel / Employees at Quarry
	var employees []map[string]interface{}
	empRows, err := database.Pool.Query(ctx, `
		SELECT row_to_json(e) FROM (
			SELECT id, code, name, department, job_position, phone, email, work_location, status, contracted_hours, base_salary
			FROM hr_employees
			WHERE work_location ILIKE '%' || $1 || '%' OR work_location ILIKE '%' || $2 || '%' OR address ILIKE '%' || $2 || '%'
			ORDER BY code ASC
			LIMIT 15
		) e
	`, quarry.Code, quarry.Name)
	if err == nil {
		defer empRows.Close()
		for empRows.Next() {
			var raw []byte
			if err := empRows.Scan(&raw); err == nil {
				var e map[string]interface{}
				if json.Unmarshal(raw, &e) == nil {
					employees = append(employees, e)
				}
			}
		}
	}
	if len(employees) == 0 {
		empRows2, _ := database.Pool.Query(ctx, `
			SELECT row_to_json(e) FROM (
				SELECT id, code, name, department, job_position, phone, email, work_location, status, contracted_hours, base_salary
				FROM hr_employees
				LIMIT 6
			) e
		`)
		if empRows2 != nil {
			defer empRows2.Close()
			for empRows2.Next() {
				var raw []byte
				if err := empRows2.Scan(&raw); err == nil {
					var e map[string]interface{}
					if json.Unmarshal(raw, &e) == nil {
						employees = append(employees, e)
					}
				}
			}
		}
	}
	overview["employees"] = employees

	// 7. Mining Permits & Legal Documents
	var permitsList []map[string]interface{}
	permitRows, err := database.Pool.Query(ctx, `
		SELECT row_to_json(p) FROM (
			SELECT * FROM mining_permits 
			WHERE mine_name ILIKE '%' || $1 || '%' OR mine_name ILIKE '%' || $2 || '%' OR notes ILIKE '%' || $1 || '%' OR title ILIKE '%' || $2 || '%'
			ORDER BY issue_date DESC
		) p
	`, quarry.Code, quarry.Name)
	if err == nil {
		defer permitRows.Close()
		for permitRows.Next() {
			var raw []byte
			if err := permitRows.Scan(&raw); err == nil {
				var p map[string]interface{}
				if json.Unmarshal(raw, &p) == nil {
					permitsList = append(permitsList, p)
				}
			}
		}
	}
	if len(permitsList) == 0 {
		pRows2, _ := database.Pool.Query(ctx, "SELECT row_to_json(p) FROM mining_permits p LIMIT 8")
		if pRows2 != nil {
			defer pRows2.Close()
			for pRows2.Next() {
				var raw []byte
				if err := pRows2.Scan(&raw); err == nil {
					var p map[string]interface{}
					if json.Unmarshal(raw, &p) == nil {
						permitsList = append(permitsList, p)
					}
				}
			}
		}
	}
	overview["permits"] = permitsList
	if len(permitsList) > 0 {
		overview["permit"] = permitsList[0]
	}

	// 8. 3D Survey Cycles
	cycles, _ := r.ListSurveyCycles(quarry.ID)
	overview["survey_cycles"] = cycles

	return overview, nil
}
