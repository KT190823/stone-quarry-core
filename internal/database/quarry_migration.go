package database

import (
	"context"
	"fmt"
	"time"
)

// MigrateQuarryModule initializes all PostgreSQL tables for Quarry 3D Survey, Volume Calculation, Weighing & Reconciliation
func MigrateQuarryModule() {
	if Pool == nil {
		return
	}
	ctx := context.Background()

	// Try enabling PostGIS if supported
	postgisCtx, cancelPostgis := context.WithTimeout(ctx, 3*time.Second)
	_, _ = Pool.Exec(postgisCtx, "CREATE EXTENSION IF NOT EXISTS postgis;")
	cancelPostgis()

	quarrySchemas := []string{
		`CREATE TABLE IF NOT EXISTS organizations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			code VARCHAR(50) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			tax_code VARCHAR(50),
			address TEXT,
			status VARCHAR(30) DEFAULT 'active',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS quarries (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			code VARCHAR(50) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			address TEXT,
			longitude DECIMAL(12, 8),
			latitude DECIMAL(12, 8),
			total_area NUMERIC(15, 2),
			boundary TEXT,
			status VARCHAR(30) DEFAULT 'active',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS quarry_areas (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			quarry_id UUID NOT NULL REFERENCES quarries(id) ON DELETE CASCADE,
			code VARCHAR(50) NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			boundary TEXT,
			area_m2 NUMERIC(18, 2) DEFAULT 0,
			status VARCHAR(30) DEFAULT 'active',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(quarry_id, code)
		)`,

		`CREATE TABLE IF NOT EXISTS survey_cycles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			quarry_id UUID NOT NULL REFERENCES quarries(id) ON DELETE CASCADE,
			cycle_code VARCHAR(100) NOT NULL,
			survey_date DATE NOT NULL DEFAULT CURRENT_DATE,
			survey_started_at TIMESTAMPTZ,
			survey_completed_at TIMESTAMPTZ,
			previous_cycle_id UUID REFERENCES survey_cycles(id) ON DELETE SET NULL,
			survey_method VARCHAR(50) DEFAULT 'drone',
			status VARCHAR(30) DEFAULT 'completed',
			operator_name VARCHAR(255),
			external_id VARCHAR(255),
			external_source VARCHAR(100),
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(quarry_id, cycle_code)
		)`,

		`CREATE TABLE IF NOT EXISTS survey_area_results (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			survey_cycle_id UUID NOT NULL REFERENCES survey_cycles(id) ON DELETE CASCADE,
			quarry_area_id UUID NOT NULL REFERENCES quarry_areas(id) ON DELETE CASCADE,
			area_m2 NUMERIC(18, 2) DEFAULT 0,
			min_elevation NUMERIC(12, 3) DEFAULT 0,
			max_elevation NUMERIC(12, 3) DEFAULT 0,
			avg_elevation NUMERIC(12, 3) DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(survey_cycle_id, quarry_area_id)
		)`,

		`CREATE TABLE IF NOT EXISTS volume_calculations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			survey_cycle_id UUID NOT NULL REFERENCES survey_cycles(id) ON DELETE CASCADE,
			quarry_area_id UUID REFERENCES quarry_areas(id) ON DELETE SET NULL,
			previous_survey_cycle_id UUID REFERENCES survey_cycles(id) ON DELETE SET NULL,
			calculation_type VARCHAR(50) NOT NULL DEFAULT 'extracted',
			previous_volume_m3 NUMERIC(20, 3) DEFAULT 0,
			current_volume_m3 NUMERIC(20, 3) DEFAULT 0,
			extracted_volume_m3 NUMERIC(20, 3) DEFAULT 0,
			fill_volume_m3 NUMERIC(20, 3) DEFAULT 0,
			net_volume_m3 NUMERIC(20, 3) DEFAULT 0,
			calculation_method VARCHAR(100) DEFAULT 'tin_surface',
			external_id VARCHAR(255),
			calculated_at TIMESTAMPTZ DEFAULT NOW(),
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS surface_models (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			survey_cycle_id UUID NOT NULL REFERENCES survey_cycles(id) ON DELETE CASCADE,
			quarry_area_id UUID REFERENCES quarry_areas(id) ON DELETE SET NULL,
			model_type VARCHAR(50) NOT NULL DEFAULT 'point_cloud',
			file_format VARCHAR(30) DEFAULT 'laz',
			storage_path TEXT NOT NULL,
			file_size BIGINT DEFAULT 0,
			coordinate_system VARCHAR(100) DEFAULT 'EPSG:4326',
			min_x NUMERIC(18, 6) DEFAULT 0,
			min_y NUMERIC(18, 6) DEFAULT 0,
			min_z NUMERIC(18, 6) DEFAULT 0,
			max_x NUMERIC(18, 6) DEFAULT 0,
			max_y NUMERIC(18, 6) DEFAULT 0,
			max_z NUMERIC(18, 6) DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS quarry_materials (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			quarry_id UUID REFERENCES quarries(id) ON DELETE CASCADE,
			code VARCHAR(50) NOT NULL,
			name VARCHAR(255) NOT NULL,
			density_ton_per_m3 NUMERIC(12, 4) DEFAULT 1.55,
			status VARCHAR(30) DEFAULT 'active',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(quarry_id, code)
		)`,

		`CREATE TABLE IF NOT EXISTS quarry_vehicles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			quarry_id UUID REFERENCES quarries(id) ON DELETE CASCADE,
			license_plate VARCHAR(30) NOT NULL,
			vehicle_type VARCHAR(100),
			capacity_ton NUMERIC(12, 3) DEFAULT 0,
			status VARCHAR(30) DEFAULT 'active',
			external_id VARCHAR(255),
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(quarry_id, license_plate)
		)`,

		`CREATE TABLE IF NOT EXISTS weighing_transactions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			quarry_id UUID NOT NULL REFERENCES quarries(id) ON DELETE CASCADE,
			vehicle_id UUID REFERENCES quarry_vehicles(id) ON DELETE SET NULL,
			material_id UUID REFERENCES quarry_materials(id) ON DELETE SET NULL,
			survey_cycle_id UUID REFERENCES survey_cycles(id) ON DELETE SET NULL,
			ticket_number VARCHAR(100),
			weight_in_kg NUMERIC(20, 3) DEFAULT 0,
			weight_out_kg NUMERIC(20, 3) DEFAULT 0,
			net_weight_kg NUMERIC(20, 3) DEFAULT 0,
			weighed_in_at TIMESTAMPTZ,
			weighed_out_at TIMESTAMPTZ,
			status VARCHAR(30) DEFAULT 'completed',
			external_id VARCHAR(255),
			external_source VARCHAR(100),
			raw_data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(external_source, external_id)
		)`,

		`CREATE TABLE IF NOT EXISTS vehicle_trip_images (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			weighing_transaction_id UUID NOT NULL REFERENCES weighing_transactions(id) ON DELETE CASCADE,
			image_type VARCHAR(50) DEFAULT 'scale',
			storage_path TEXT NOT NULL,
			captured_at TIMESTAMPTZ DEFAULT NOW(),
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS production_reconciliations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			quarry_id UUID NOT NULL REFERENCES quarries(id) ON DELETE CASCADE,
			quarry_area_id UUID REFERENCES quarry_areas(id) ON DELETE SET NULL,
			survey_cycle_id UUID NOT NULL REFERENCES survey_cycles(id) ON DELETE CASCADE,
			material_id UUID REFERENCES quarry_materials(id) ON DELETE SET NULL,
			volume_m3 NUMERIC(20, 3) DEFAULT 0,
			density_ton_per_m3 NUMERIC(12, 4) DEFAULT 1.55,
			expected_weight_ton NUMERIC(20, 3) DEFAULT 0,
			actual_weight_ton NUMERIC(20, 3) DEFAULT 0,
			difference_ton NUMERIC(20, 3) DEFAULT 0,
			difference_percent NUMERIC(10, 4) DEFAULT 0,
			status VARCHAR(30) DEFAULT 'matched',
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS quarry_alerts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			quarry_id UUID NOT NULL REFERENCES quarries(id) ON DELETE CASCADE,
			quarry_area_id UUID REFERENCES quarry_areas(id) ON DELETE SET NULL,
			survey_cycle_id UUID REFERENCES survey_cycles(id) ON DELETE SET NULL,
			alert_type VARCHAR(100) NOT NULL,
			severity VARCHAR(20) NOT NULL DEFAULT 'warning',
			title VARCHAR(255),
			message TEXT,
			expected_value NUMERIC(20, 3) DEFAULT 0,
			actual_value NUMERIC(20, 3) DEFAULT 0,
			difference_percent NUMERIC(10, 4) DEFAULT 0,
			is_read BOOLEAN DEFAULT FALSE,
			is_resolved BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			resolved_at TIMESTAMPTZ
		)`,

		`CREATE TABLE IF NOT EXISTS integration_events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			source VARCHAR(100) NOT NULL,
			event_type VARCHAR(100) NOT NULL,
			external_id VARCHAR(255),
			payload JSONB NOT NULL DEFAULT '{}',
			status VARCHAR(30) DEFAULT 'received',
			error_message TEXT,
			received_at TIMESTAMPTZ DEFAULT NOW(),
			processed_at TIMESTAMPTZ
		)`,
	}

	for _, schema := range quarrySchemas {
		stepCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		if _, err := Pool.Exec(stepCtx, schema); err != nil {
			fmt.Printf("⚠️ Quarry schema notice: %v\n", err)
		}
		cancel()
	}

	fmt.Println("⛏️ Quarry Core PostgreSQL tables migrated successfully!")
}
