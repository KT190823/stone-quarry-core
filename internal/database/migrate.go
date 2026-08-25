package database

import (
	"context"
	"fmt"
	"log"
	"time"
)

func Migrate() {
	if Pool == nil {
		fmt.Println("⚠️ Skipping DB migration as database pool is nil.")
		return
	}
	ctx := context.Background()
	if err := Pool.Ping(ctx); err != nil {
		fmt.Printf("⚠️ Skipping DB migration: %v\n", err)
		return
	}

	schemas := []string{
		`CREATE TABLE IF NOT EXISTS tickets (
			id TEXT PRIMARY KEY,
			stt INTEGER,
			ben_ban TEXT,
			ben_mua TEXT,
			bien_so TEXT,
			loai_xe TEXT,
			lai_xe TEXT,
			sdt_lai_xe TEXT,
			rfid TEXT,
			loai TEXT,
			stage TEXT,
			stage_label TEXT,
			can_l1 DOUBLE PRECISION DEFAULT 0,
			kl1 TEXT,
			can_l2 TEXT,
			kl2 TEXT,
			kl_hang TEXT,
			kl_tap_chat DOUBLE PRECISION DEFAULT 0,
			kl_tinh_tien TEXT,
			don_gia DOUBLE PRECISION DEFAULT 0,
			thanh_tien TEXT,
			time1 TEXT,
			time2 TEXT,
			date TEXT,
			nguoi_can1 TEXT,
			nguoi_can2 TEXT,
			mat_hang TEXT,
			quy_cach TEXT,
			do_code TEXT,
			tram_can TEXT,
			cong_can TEXT,
			ghi_chu TEXT,
			hoa_don_so TEXT,
			cameras JSONB DEFAULT '{}',
			chatter JSONB DEFAULT '[]',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS vehicles (
			id SERIAL PRIMARY KEY,
			bs TEXT UNIQUE,
			loai TEXT,
			bi TEXT,
			rfid TEXT,
			tai_trong TEXT,
			status TEXT,
			count INTEGER DEFAULT 0,
			han_dang_kiem TEXT,
			date TEXT,
			chu_xe TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS alerts (
			id TEXT PRIMARY KEY,
			title TEXT,
			bs TEXT,
			note TEXT,
			time TEXT,
			date TEXT,
			status TEXT,
			severity TEXT,
			phieu TEXT,
			bi_dang_ky DOUBLE PRECISION DEFAULT 0,
			bi_thuc_te DOUBLE PRECISION DEFAULT 0,
			lech_bi DOUBLE PRECISION DEFAULT 0,
			cam TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS catalog_materials (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			ma TEXT,
			name TEXT,
			ten TEXT,
			dvt TEXT,
			density TEXT,
			ty_trong TEXT,
			dinh_muc TEXT,
			price TEXT,
			gia TEXT,
			kho TEXT,
			standard TEXT,
			status TEXT,
			date TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS suppliers (
			id SERIAL PRIMARY KEY,
			name TEXT,
			ten TEXT,
			tax_code TEXT,
			mst TEXT,
			phone TEXT,
			sdt TEXT,
			items TEXT,
			hang TEXT,
			no TEXT,
			status TEXT,
			rating TEXT,
			date TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			name TEXT,
			ten TEXT,
			tax_code TEXT,
			mst TEXT,
			contact TEXT,
			phone TEXT,
			sdt TEXT,
			group_name TEXT,
			doanh_so TEXT,
			cong_no TEXT,
			han_muc TEXT,
			status TEXT,
			date TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS weigh_bridges (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			name TEXT,
			ten TEXT,
			location TEXT,
			ip TEXT,
			type TEXT,
			capacity TEXT,
			loadcell TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS camera_devices (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			name TEXT,
			ten TEXT,
			location TEXT,
			ip TEXT,
			res TEXT,
			status TEXT,
			type TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS system_users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE,
			name TEXT,
			ten TEXT,
			role TEXT,
			dept TEXT,
			email TEXT,
			last_login TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS hr_employees (
			id TEXT PRIMARY KEY,
			code TEXT,
			name TEXT,
			avatar TEXT,
			gender TEXT,
			dob TEXT,
			phone TEXT,
			email TEXT,
			id_card TEXT,
			address TEXT,
			department TEXT,
			job_position TEXT,
			manager TEXT,
			join_date TEXT,
			contract_type TEXT,
			contracted_hours DOUBLE PRECISION DEFAULT 8.0,
			work_location TEXT,
			status TEXT,
			certificates JSONB DEFAULT '[]',
			base_salary DOUBLE PRECISION DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS hr_contracts (
			id TEXT PRIMARY KEY,
			code TEXT,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			department TEXT,
			contract_type TEXT,
			contracted_hours DOUBLE PRECISION DEFAULT 8.0,
			start_date TEXT,
			end_date TEXT,
			base_salary DOUBLE PRECISION DEFAULT 0,
			hazard_allowance DOUBLE PRECISION DEFAULT 0,
			safety_allowance DOUBLE PRECISION DEFAULT 0,
			meal_allowance DOUBLE PRECISION DEFAULT 0,
			responsibility_allowance DOUBLE PRECISION DEFAULT 0,
			attendance_allowance DOUBLE PRECISION DEFAULT 0,
			social_insurance_base DOUBLE PRECISION DEFAULT 0,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS hr_attendances (
			id TEXT PRIMARY KEY,
			employee_id TEXT,
			employee_name TEXT,
			department TEXT,
			job_position TEXT,
			equipment_type TEXT,
			equipment_code TEXT,
			vehicle_plate TEXT,
			date TEXT,
			check_in_time TEXT,
			check_out_time TEXT,
			contracted_hours DOUBLE PRECISION DEFAULT 8.0,
			worked_hours DOUBLE PRECISION DEFAULT 0,
			ot_hours DOUBLE PRECISION DEFAULT 0,
			method TEXT,
			location TEXT,
			work_area TEXT,
			status TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS hr_leave_requests (
			id TEXT PRIMARY KEY,
			code TEXT,
			employee_id TEXT,
			employee_name TEXT,
			department TEXT,
			leave_type TEXT,
			start_date TEXT,
			end_date TEXT,
			total_days DOUBLE PRECISION DEFAULT 0,
			total_hours DOUBLE PRECISION DEFAULT 0,
			reason TEXT,
			stage TEXT,
			approver TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS hr_leave_allocations (
			id SERIAL PRIMARY KEY,
			employee_id TEXT,
			employee_name TEXT,
			department TEXT,
			year INTEGER,
			total_allocated INTEGER DEFAULT 0,
			used_days INTEGER DEFAULT 0,
			remaining_days INTEGER DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS hr_timesheets (
			id TEXT PRIMARY KEY,
			data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS hr_payslips (
			id TEXT PRIMARY KEY,
			data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS mining_permits (
			id TEXT PRIMARY KEY,
			code TEXT,
			title TEXT,
			mine_name TEXT,
			category TEXT,
			category_label TEXT,
			issuer TEXT,
			license_number TEXT,
			issue_date TEXT,
			expiry_date TEXT,
			capacity TEXT,
			approved_reserve TEXT,
			mined_so_far TEXT,
			mined_percent INTEGER DEFAULT 0,
			depth_level TEXT,
			area TEXT,
			coordinates TEXT,
			status TEXT,
			status_label TEXT,
			days_remaining INTEGER DEFAULT 0,
			files JSONB DEFAULT '[]',
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS mining_plans (
			id TEXT PRIMARY KEY,
			mine TEXT,
			item TEXT,
			annual_target DOUBLE PRECISION DEFAULT 0,
			unit TEXT,
			q1_plan DOUBLE PRECISION DEFAULT 0,
			q1_actual DOUBLE PRECISION DEFAULT 0,
			q2_plan DOUBLE PRECISION DEFAULT 0,
			q2_actual DOUBLE PRECISION DEFAULT 0,
			q3_plan DOUBLE PRECISION DEFAULT 0,
			q3_actual DOUBLE PRECISION DEFAULT 0,
			q4_plan DOUBLE PRECISION DEFAULT 0,
			q4_actual DOUBLE PRECISION DEFAULT 0,
			ytd_actual DOUBLE PRECISION DEFAULT 0,
			completion_rate DOUBLE PRECISION DEFAULT 0,
			status TEXT,
			status_label TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS blasting_passports (
			id TEXT PRIMARY KEY,
			code TEXT,
			mine_name TEXT,
			blast_date TEXT,
			blast_time TEXT,
			location TEXT,
			hole_count INTEGER DEFAULT 0,
			hole_depth_meters DOUBLE PRECISION DEFAULT 0,
			anfo_explosive_kg DOUBLE PRECISION DEFAULT 0,
			emulsion_explosive_kg DOUBLE PRECISION DEFAULT 0,
			detonator_count INTEGER DEFAULT 0,
			designed_rock_volume_m3 DOUBLE PRECISION DEFAULT 0,
			powder_factor_kg_per_m3 DOUBLE PRECISION DEFAULT 0,
			actual_rock_mined_m3 DOUBLE PRECISION DEFAULT 0,
			safety_status TEXT,
			blaster_in_charge TEXT,
			certified_number TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS crusher_plants (
			id SERIAL PRIMARY KEY,
			plant_code TEXT,
			plant_name TEXT,
			capacity_ton_per_hour INTEGER DEFAULT 0,
			input_rock_today_tons INTEGER DEFAULT 0,
			power_consumption_kwh INTEGER DEFAULT 0,
			kwh_per_ton DOUBLE PRECISION DEFAULT 0,
			yield_breakdown JSONB DEFAULT '[]',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS equipment_fuel_logs (
			id TEXT PRIMARY KEY,
			equipment_code TEXT,
			equipment_name TEXT,
			category TEXT,
			operator_name TEXT,
			hours_worked_today DOUBLE PRECISION DEFAULT 0,
			total_hours_meter DOUBLE PRECISION DEFAULT 0,
			fuel_quota_liters_per_hour DOUBLE PRECISION DEFAULT 0,
			actual_fuel_issued_liters DOUBLE PRECISION DEFAULT 0,
			actual_fuel_consumed_liters DOUBLE PRECISION DEFAULT 0,
			fuel_variance_liters DOUBLE PRECISION DEFAULT 0,
			variance_status TEXT,
			location TEXT,
			maintenance_status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS inventory_inbound (
			id SERIAL PRIMARY KEY,
			code TEXT,
			source TEXT,
			loc TEXT,
			item TEXT,
			qty TEXT,
			quantity TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS inventory_outbound (
			id SERIAL PRIMARY KEY,
			code TEXT,
			customer TEXT,
			dest TEXT,
			item TEXT,
			qty TEXT,
			quantity TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS inventory_stocktake (
			id SERIAL PRIMARY KEY,
			code TEXT,
			zone TEXT,
			item TEXT,
			volume TEXT,
			survey TEXT,
			erp TEXT,
			book TEXT,
			actual TEXT,
			diff TEXT,
			quantity TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS inventory_movements (
			id SERIAL PRIMARY KEY,
			code TEXT,
			from_loc TEXT,
			to_loc TEXT,
			item TEXT,
			qty TEXT,
			quantity TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS payments_debt (
			id SERIAL PRIMARY KEY,
			code TEXT,
			partner TEXT,
			customer TEXT,
			limit_val TEXT,
			balance TEXT,
			amount TEXT,
			due TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS payments_reconcile (
			id SERIAL PRIMARY KEY,
			code TEXT,
			customer TEXT,
			scale_amount TEXT,
			bank_amount TEXT,
			diff TEXT,
			status TEXT,
			date TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS payments_pricing (
			id SERIAL PRIMARY KEY,
			code TEXT,
			name TEXT,
			ten TEXT,
			price TEXT,
			gia TEXT,
			vat TEXT,
			unit TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS payments_prepaid (
			id SERIAL PRIMARY KEY,
			code TEXT,
			customer TEXT,
			initial_balance TEXT,
			current_balance TEXT,
			status TEXT,
			date TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS gps_fleet (
			id SERIAL PRIMARY KEY,
			code TEXT,
			data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS fuel_theft_audits (
			id SERIAL PRIMARY KEY,
			code TEXT,
			data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS fuel_norms (
			id SERIAL PRIMARY KEY,
			code TEXT,
			data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS yard_checkinout (
			id SERIAL PRIMARY KEY,
			code TEXT,
			data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS google_drive_files (
			id SERIAL PRIMARY KEY,
			code TEXT,
			data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS google_emails (
			id SERIAL PRIMARY KEY,
			code TEXT,
			data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS google_map_entries (
			id SERIAL PRIMARY KEY,
			code TEXT,
			data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS google_photo_entries (
			id SERIAL PRIMARY KEY,
			code TEXT,
			data JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS print_templates (
			id TEXT PRIMARY KEY,
			code TEXT UNIQUE,
			name TEXT,
			doc_type TEXT,
			size TEXT DEFAULT 'A5',
			orientation TEXT DEFAULT 'portrait',
			description TEXT,
			layout TEXT,
			is_default BOOLEAN DEFAULT FALSE,
			status TEXT DEFAULT 'active',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE,
			name TEXT,
			ten TEXT,
			role TEXT,
			dept TEXT,
			email TEXT,
			phone TEXT,
			sdt TEXT,
			password_hash TEXT,
			last_login TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS user_roles (
			id SERIAL PRIMARY KEY,
			role TEXT,
			"desc" TEXT,
			description TEXT,
			users TEXT,
			level TEXT,
			permission TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS user_logs (
			id SERIAL PRIMARY KEY,
			time TEXT,
			"user" TEXT,
			username TEXT,
			name TEXT,
			action TEXT,
			target TEXT,
			ip TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS reports (
			id TEXT PRIMARY KEY,
			name TEXT,
			item TEXT,
			type TEXT,
			period TEXT,
			plan TEXT,
			actual TEXT,
			diff TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS settings (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			key TEXT,
			name TEXT,
			val TEXT,
			scope TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS ticket_types (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			name TEXT,
			ten TEXT,
			loai TEXT,
			description TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS vehicle_catalogs (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			name TEXT,
			loai TEXT,
			tai_trong TEXT,
			so_truc TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS stations (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			name TEXT,
			ten TEXT,
			location TEXT,
			ip TEXT,
			capacity TEXT,
			type TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS equipment (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			name TEXT,
			ten TEXT,
			loai TEXT,
			capacity TEXT,
			location TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS hr_shifts (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			name TEXT,
			start_time TEXT,
			end_time TEXT,
			hours DOUBLE PRECISION DEFAULT 8.0,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS hr_shift_schedules (
			id SERIAL PRIMARY KEY,
			code TEXT,
			shift_id TEXT,
			employee_id TEXT,
			date TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS statutory_reports (
			id TEXT PRIMARY KEY,
			code TEXT,
			name TEXT,
			item TEXT,
			period TEXT,
			plan TEXT,
			actual TEXT,
			diff TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS resource_taxes (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			material TEXT,
			tax_rate TEXT,
			environmental_fee TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS production_stages (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			name TEXT,
			stage TEXT,
			description TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// ===================== HRM EXTENDED MODULES =====================
		// Employee lifecycle
		`CREATE TABLE IF NOT EXISTS hr_employee_work_history (
			id TEXT PRIMARY KEY,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			event_type TEXT,
			title TEXT,
			detail TEXT,
			changed_by TEXT,
			changed_at TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_onboarding_checklists (
			id TEXT PRIMARY KEY,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			task_name TEXT,
			task_description TEXT,
			attachment_id TEXT,
			is_completed BOOLEAN DEFAULT FALSE,
			completed_at TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_decisions (
			id TEXT PRIMARY KEY,
			code TEXT,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			decision_type TEXT,
			effective_date TEXT,
			content_detail JSONB DEFAULT '{}',
			status TEXT,
			issued_by TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_esign_documents (
			id TEXT PRIMARY KEY,
			code TEXT,
			contract_id TEXT,
			contract_code TEXT,
			subject TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_esign_signers (
			id TEXT PRIMARY KEY,
			document_id TEXT REFERENCES hr_esign_documents(id),
			signer_type TEXT,
			employee_id TEXT,
			employee_name TEXT,
			sign_order INTEGER DEFAULT 1,
			sign_status TEXT,
			sign_date TEXT,
			sign_note TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Insurance
		`CREATE TABLE IF NOT EXISTS hr_insurance_records (
			id TEXT PRIMARY KEY,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			insurance_number TEXT,
			type TEXT,
			participation_date TEXT,
			stop_date TEXT,
			social_insurance_base DOUBLE PRECISION DEFAULT 0,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_insurance_declarations (
			id TEXT PRIMARY KEY,
			code TEXT,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			declaration_type TEXT,
			period TEXT,
			status TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Recruitment
		`CREATE TABLE IF NOT EXISTS hr_recruitment_requests (
			id TEXT PRIMARY KEY,
			code TEXT,
			title TEXT,
			department TEXT,
			position TEXT,
			headcount INTEGER DEFAULT 1,
			reason TEXT,
			requested_by TEXT,
			approval_status TEXT,
			current_step INTEGER DEFAULT 0,
			comment TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_recruitment_campaigns (
			id TEXT PRIMARY KEY,
			code TEXT,
			name TEXT,
			request_id TEXT,
			position TEXT,
			department TEXT,
			open_date TEXT,
			recruit_deadline TEXT,
			target_headcount INTEGER DEFAULT 0,
			hired_count INTEGER DEFAULT 0,
			source TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_candidates (
			id TEXT PRIMARY KEY,
			code TEXT,
			name TEXT,
			phone TEXT,
			email TEXT,
			position TEXT,
			source TEXT,
			cv_text TEXT,
			pipeline_stage TEXT,
			list_status TEXT,
			note TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_candidate_applications (
			id TEXT PRIMARY KEY,
			candidate_id TEXT REFERENCES hr_candidates(id),
			campaign_id TEXT,
			pipeline_stage TEXT,
			applied_at TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_candidate_status_history (
			id TEXT PRIMARY KEY,
			application_id TEXT,
			candidate_id TEXT,
			from_stage TEXT,
			to_stage TEXT,
			changed_by TEXT,
			note TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_interview_schedules (
			id TEXT PRIMARY KEY,
			candidate_id TEXT REFERENCES hr_candidates(id),
			candidate_name TEXT,
			position TEXT,
			interview_date TEXT,
			interview_time TEXT,
			interviewer TEXT,
			round INTEGER DEFAULT 1,
			meeting_type TEXT,
			result TEXT,
			note TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Training
		`CREATE TABLE IF NOT EXISTS hr_training_courses (
			id TEXT PRIMARY KEY,
			code TEXT,
			name TEXT,
			category TEXT,
			provider TEXT,
			start_date TEXT,
			end_date TEXT,
			venue TEXT,
			instructor TEXT,
			duration_hours DOUBLE PRECISION DEFAULT 0,
			attendance_threshold DOUBLE PRECISION DEFAULT 80,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_training_sessions (
			id TEXT PRIMARY KEY,
			course_id TEXT REFERENCES hr_training_courses(id),
			session_name TEXT,
			session_date TEXT,
			start_time TEXT,
			end_time TEXT,
			instructor TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_training_participants (
			id TEXT PRIMARY KEY,
			course_id TEXT REFERENCES hr_training_courses(id),
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			attendance_rate DOUBLE PRECISION DEFAULT 0,
			completion_status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_training_attendances (
			id TEXT PRIMARY KEY,
			session_id TEXT,
			participant_id TEXT,
			employee_id TEXT,
			employee_name TEXT,
			attendance_status TEXT,
			note TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_training_evaluations (
			id TEXT PRIMARY KEY,
			participant_id TEXT,
			course_id TEXT,
			employee_id TEXT,
			employee_name TEXT,
			criteria TEXT,
			rating NUMERIC DEFAULT 0,
			comment TEXT,
			submitted_at TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// KPI / OKR
		`CREATE TABLE IF NOT EXISTS hr_kpi_targets (
			id TEXT PRIMARY KEY,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			period TEXT,
			kpi_code TEXT,
			kpi_name TEXT,
			target_value DOUBLE PRECISION DEFAULT 0,
			actual_value DOUBLE PRECISION DEFAULT 0,
			unit TEXT,
			weight_percent DOUBLE PRECISION DEFAULT 0,
			score DOUBLE PRECISION DEFAULT 0,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_okr_objectives (
			id TEXT PRIMARY KEY,
			code TEXT,
			title TEXT,
			parent_objective_id TEXT,
			owner_text TEXT,
			period TEXT,
			level TEXT,
			progress_percent DOUBLE PRECISION DEFAULT 0,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_okr_key_results (
			id TEXT PRIMARY KEY,
			objective_id TEXT REFERENCES hr_okr_objectives(id),
			key_result TEXT,
			start_value DOUBLE PRECISION DEFAULT 0,
			target_value DOUBLE PRECISION DEFAULT 0,
			current_value DOUBLE PRECISION DEFAULT 0,
			unit TEXT,
			progress_percent DOUBLE PRECISION DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Workflow engine (generic approval)
		`CREATE TABLE IF NOT EXISTS hr_workflow_templates (
			id TEXT PRIMARY KEY,
			code TEXT,
			name TEXT,
			applies_to TEXT,
			definition JSONB DEFAULT '{}',
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_workflow_instances (
			id TEXT PRIMARY KEY,
			code TEXT,
			template_id TEXT,
			applies_to TEXT,
			object_id TEXT,
			object_ref TEXT,
			status TEXT,
			current_step INTEGER DEFAULT 0,
			started_at TEXT,
			finished_at TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_workflow_steps (
			id TEXT PRIMARY KEY,
			instance_id TEXT REFERENCES hr_workflow_instances(id),
			step_order INTEGER DEFAULT 0,
			step_type TEXT,
			approver_role TEXT,
			approver_id TEXT,
			approver_name TEXT,
			step_status TEXT,
			comment TEXT,
			decided_at TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Timesheet & payroll relational (replaces JSONB per-period)
		`CREATE TABLE IF NOT EXISTS hr_shift_assignments (
			id TEXT PRIMARY KEY,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			shift_code TEXT,
			shift_name TEXT,
			date TEXT,
			time_slot TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_holidays (
			id TEXT PRIMARY KEY,
			holiday_name TEXT,
			holiday_date TEXT,
			date_type TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_attendance_logs (
			id TEXT PRIMARY KEY,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			date TEXT,
			check_in_time TEXT,
			check_out_time TEXT,
			gps_lat DOUBLE PRECISION DEFAULT 0,
			gps_lng DOUBLE PRECISION DEFAULT 0,
			device_id TEXT,
			work_value DOUBLE PRECISION DEFAULT 0,
			status_code TEXT,
			overtime_hours DOUBLE PRECISION DEFAULT 0,
			note TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_timesheet_details (
			id TEXT PRIMARY KEY,
			timesheet_id TEXT,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			date TEXT,
			work_value DOUBLE PRECISION DEFAULT 0,
			status_code TEXT,
			overtime_hours DOUBLE PRECISION DEFAULT 0,
			note TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_payroll_periods (
			id TEXT PRIMARY KEY,
			code TEXT,
			month TEXT,
			status TEXT,
			is_locked BOOLEAN DEFAULT FALSE,
			start_date TEXT,
			end_date TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_payroll_records (
			id TEXT PRIMARY KEY,
			period_id TEXT,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			department TEXT,
			contract_type TEXT,
			bank_name TEXT,
			bank_account TEXT,
			base_salary DOUBLE PRECISION DEFAULT 0,
			actual_days DOUBLE PRECISION DEFAULT 0,
			standard_days DOUBLE PRECISION DEFAULT 0,
			working_salary DOUBLE PRECISION DEFAULT 0,
			piecework_salary DOUBLE PRECISION DEFAULT 0,
			overtime_salary DOUBLE PRECISION DEFAULT 0,
			allowances DOUBLE PRECISION DEFAULT 0,
			gross_salary DOUBLE PRECISION DEFAULT 0,
			social_insurance DOUBLE PRECISION DEFAULT 0,
			health_insurance DOUBLE PRECISION DEFAULT 0,
			unemployment_insurance DOUBLE PRECISION DEFAULT 0,
			total_insurance DOUBLE PRECISION DEFAULT 0,
			personal_income_tax DOUBLE PRECISION DEFAULT 0,
			advance_deduction DOUBLE PRECISION DEFAULT 0,
			total_deductions DOUBLE PRECISION DEFAULT 0,
			net_salary DOUBLE PRECISION DEFAULT 0,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_payroll_formula_columns (
			id TEXT PRIMARY KEY,
			code TEXT,
			name TEXT,
			col_type TEXT,
			formula_expression TEXT,
			sort_order INTEGER DEFAULT 0,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS hr_salary_advances (
			id TEXT PRIMARY KEY,
			code TEXT,
			employee_id TEXT REFERENCES hr_employees(id),
			employee_name TEXT,
			amount DOUBLE PRECISION DEFAULT 0,
			request_date TEXT,
			deduction_period TEXT,
			reason TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
	}

	for i, schema := range schemas {
		stepCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		_, err := Pool.Exec(stepCtx, schema)
		cancel()
		if err != nil {
			log.Printf("Migration step %d error: %v", i+1, err)
		}
	}

	// Create view fallbacks for compatibility
	views := []string{
		`CREATE OR REPLACE VIEW driver_gps_fleet AS SELECT * FROM gps_fleet`,
		`CREATE OR REPLACE VIEW fuel_norm_configs AS SELECT * FROM fuel_norms`,
		`CREATE OR REPLACE VIEW yard_check_in_outs AS SELECT * FROM yard_checkinout`,
	}
	for _, v := range views {
		stepCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		Pool.Exec(stepCtx, v)
		cancel()
	}

	// Add optional columns that may be missing on pre-existing tables
	alters := []string{
		`ALTER TABLE hr_attendance_logs ADD COLUMN IF NOT EXISTS overtime_hours DOUBLE PRECISION DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT`,
		`ALTER TABLE statutory_reports ADD COLUMN IF NOT EXISTS title TEXT`,
		`ALTER TABLE statutory_reports ADD COLUMN IF NOT EXISTS recipient TEXT`,
		`ALTER TABLE statutory_reports ADD COLUMN IF NOT EXISTS date TEXT`,
		`ALTER TABLE statutory_reports ADD COLUMN IF NOT EXISTS mined_volume TEXT`,
		`ALTER TABLE statutory_reports ADD COLUMN IF NOT EXISTS tax_amount TEXT`,
		`ALTER TABLE statutory_reports ADD COLUMN IF NOT EXISTS env_fee_amount TEXT`,
		`ALTER TABLE statutory_reports ADD COLUMN IF NOT EXISTS status_label TEXT`,
		`ALTER TABLE resource_taxes ADD COLUMN IF NOT EXISTS mineral_type TEXT`,
		`ALTER TABLE resource_taxes ADD COLUMN IF NOT EXISTS mined_volume NUMERIC DEFAULT 0`,
		`ALTER TABLE resource_taxes ADD COLUMN IF NOT EXISTS tax_price_per_unit NUMERIC DEFAULT 0`,
		`ALTER TABLE resource_taxes ADD COLUMN IF NOT EXISTS resource_tax_amount NUMERIC DEFAULT 0`,
		`ALTER TABLE resource_taxes ADD COLUMN IF NOT EXISTS total_payable NUMERIC DEFAULT 0`,
		`ALTER TABLE production_stages ADD COLUMN IF NOT EXISTS stage_number INTEGER DEFAULT 0`,
		`ALTER TABLE production_stages ADD COLUMN IF NOT EXISTS stage_name TEXT`,
		`ALTER TABLE production_stages ADD COLUMN IF NOT EXISTS icon TEXT`,
		`ALTER TABLE production_stages ADD COLUMN IF NOT EXISTS volume_month TEXT`,
		`ALTER TABLE production_stages ADD COLUMN IF NOT EXISTS volume_ytd TEXT`,
		`ALTER TABLE production_stages ADD COLUMN IF NOT EXISTS loss_rate TEXT`,
		`ALTER TABLE production_stages ADD COLUMN IF NOT EXISTS loss_status TEXT`,
		`ALTER TABLE production_stages ADD COLUMN IF NOT EXISTS measurement_method TEXT`,
		`ALTER TABLE payments_debt ADD COLUMN IF NOT EXISTS debt TEXT`,
		`ALTER TABLE payments_debt ADD COLUMN IF NOT EXISTS partner TEXT`,
		`ALTER TABLE payments_debt ADD COLUMN IF NOT EXISTS limit_amount TEXT`,
		`ALTER TABLE payments_debt ADD COLUMN IF NOT EXISTS limit_val TEXT`,
		`ALTER TABLE payments_debt ADD COLUMN IF NOT EXISTS balance TEXT`,
		`ALTER TABLE payments_debt ADD COLUMN IF NOT EXISTS due TEXT`,
		`ALTER TABLE payments_reconcile ADD COLUMN IF NOT EXISTS day TEXT`,
		`ALTER TABLE payments_reconcile ADD COLUMN IF NOT EXISTS period TEXT`,
		`ALTER TABLE payments_reconcile ADD COLUMN IF NOT EXISTS scale_rev TEXT`,
		`ALTER TABLE payments_reconcile ADD COLUMN IF NOT EXISTS acc_rev TEXT`,
		`ALTER TABLE payments_reconcile ADD COLUMN IF NOT EXISTS scale_revenue TEXT`,
		`ALTER TABLE payments_reconcile ADD COLUMN IF NOT EXISTS erp_revenue TEXT`,
		`ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS items TEXT`,
		`ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS phone TEXT`,
		`ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS sdt TEXT`,
		`ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS no TEXT`,
		`ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS rating TEXT`,
		`ALTER TABLE inventory_inbound ADD COLUMN IF NOT EXISTS date TEXT`,
		`ALTER TABLE inventory_outbound ADD COLUMN IF NOT EXISTS date TEXT`,
		`ALTER TABLE inventory_stocktake ADD COLUMN IF NOT EXISTS date TEXT`,
		`ALTER TABLE inventory_movements ADD COLUMN IF NOT EXISTS date TEXT`,
	}
	for _, a := range alters {
		stepCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		Pool.Exec(stepCtx, a)
		cancel()
	}

	fmt.Println("Database migration completed successfully!")
}
