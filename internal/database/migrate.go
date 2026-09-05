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
			bi DOUBLE PRECISION DEFAULT 0,
			rfid TEXT,
			tai_trong DOUBLE PRECISION DEFAULT 0,
			unit TEXT DEFAULT 'tấn',
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
			dvt TEXT DEFAULT 'tấn',
			unit TEXT DEFAULT 'tấn',
			density DOUBLE PRECISION DEFAULT 0,
			ty_trong DOUBLE PRECISION DEFAULT 0,
			dinh_muc DOUBLE PRECISION DEFAULT 0,
			price DOUBLE PRECISION DEFAULT 0,
			gia DOUBLE PRECISION DEFAULT 0,
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
			qty DOUBLE PRECISION DEFAULT 0,
			quantity DOUBLE PRECISION DEFAULT 0,
			unit TEXT DEFAULT 'tấn',
			date TEXT,
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
			qty DOUBLE PRECISION DEFAULT 0,
			quantity DOUBLE PRECISION DEFAULT 0,
			unit TEXT DEFAULT 'tấn',
			date TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS inventory_stocktake (
			id SERIAL PRIMARY KEY,
			code TEXT,
			zone TEXT,
			item TEXT,
			volume DOUBLE PRECISION DEFAULT 0,
			survey DOUBLE PRECISION DEFAULT 0,
			erp DOUBLE PRECISION DEFAULT 0,
			book DOUBLE PRECISION DEFAULT 0,
			actual DOUBLE PRECISION DEFAULT 0,
			diff DOUBLE PRECISION DEFAULT 0,
			quantity DOUBLE PRECISION DEFAULT 0,
			unit TEXT DEFAULT 'tấn',
			date TEXT,
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
			qty DOUBLE PRECISION DEFAULT 0,
			quantity DOUBLE PRECISION DEFAULT 0,
			unit TEXT DEFAULT 'tấn',
			date TEXT,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS payments_invoices (
			id SERIAL PRIMARY KEY,
			code TEXT,
			partner TEXT,
			customer TEXT,
			amount TEXT,
			balance TEXT,
			due TEXT,
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
			plan DOUBLE PRECISION DEFAULT 0,
			actual DOUBLE PRECISION DEFAULT 0,
			diff TEXT,
			unit TEXT DEFAULT 'tấn',
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
			tai_trong DOUBLE PRECISION DEFAULT 0,
			unit TEXT DEFAULT 'tấn',
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
			title TEXT,
			recipient TEXT,
			period TEXT,
			date TEXT,
			mined_volume DOUBLE PRECISION DEFAULT 0,
			unit TEXT DEFAULT 'tấn',
			tax_amount DOUBLE PRECISION DEFAULT 0,
			env_fee_amount DOUBLE PRECISION DEFAULT 0,
			status TEXT,
			status_label TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS resource_taxes (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			mineral_type TEXT,
			mined_volume DOUBLE PRECISION DEFAULT 0,
			unit TEXT DEFAULT 'tấn',
			tax_price_per_unit DOUBLE PRECISION DEFAULT 0,
			tax_rate TEXT,
			resource_tax_amount DOUBLE PRECISION DEFAULT 0,
			environmental_fee DOUBLE PRECISION DEFAULT 0,
			total_payable DOUBLE PRECISION DEFAULT 0,
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS production_stages (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			stage_number INTEGER DEFAULT 0,
			stage_name TEXT,
			icon TEXT,
			volume_month DOUBLE PRECISION DEFAULT 0,
			volume_ytd DOUBLE PRECISION DEFAULT 0,
			unit TEXT DEFAULT 'tấn',
			loss_rate TEXT,
			loss_status TEXT,
			measurement_method TEXT,
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
			document_type TEXT,
			partner_name TEXT,
			contract_value NUMERIC DEFAULT 0,
			volume TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`ALTER TABLE hr_esign_documents ADD COLUMN IF NOT EXISTS document_type TEXT`,
		`ALTER TABLE hr_esign_documents ADD COLUMN IF NOT EXISTS partner_name TEXT`,
		`ALTER TABLE hr_esign_documents ADD COLUMN IF NOT EXISTS contract_value NUMERIC DEFAULT 0`,
		`ALTER TABLE hr_esign_documents ADD COLUMN IF NOT EXISTS volume TEXT`,
		`CREATE TABLE IF NOT EXISTS hr_esign_signers (
			id TEXT PRIMARY KEY,
			document_id TEXT REFERENCES hr_esign_documents(id),
			signer_type TEXT,
			employee_id TEXT,
			employee_name TEXT,
			role_title TEXT,
			sign_order INTEGER DEFAULT 1,
			sign_status TEXT,
			sign_date TEXT,
			sign_note TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`ALTER TABLE hr_esign_signers ADD COLUMN IF NOT EXISTS role_title TEXT`,

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

		// Unit columns standardization
		`ALTER TABLE catalog_materials ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE vehicle_catalogs ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE inventory_inbound ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE inventory_outbound ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE inventory_stocktake ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE inventory_movements ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE statutory_reports ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE production_stages ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE reports ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE resource_taxes ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'tấn'`,
		`ALTER TABLE customers ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'VNĐ'`,
		`ALTER TABLE payments_debt ADD COLUMN IF NOT EXISTS unit TEXT DEFAULT 'VNĐ'`,

		// Clean thousand separators before casting
		`UPDATE statutory_reports SET tax_amount = NULLIF(regexp_replace(tax_amount::text, '[^0-9]', '', 'g'), '') WHERE tax_amount::text LIKE '%.%.%'`,
		`UPDATE statutory_reports SET env_fee_amount = NULLIF(regexp_replace(env_fee_amount::text, '[^0-9]', '', 'g'), '') WHERE env_fee_amount::text LIKE '%.%.%'`,
		`UPDATE production_stages SET volume_ytd = NULLIF(regexp_replace(volume_ytd::text, '[^0-9]', '', 'g'), '') WHERE volume_ytd::text LIKE '%.%.%'`,
		`UPDATE resource_taxes SET environmental_fee = NULLIF(regexp_replace(environmental_fee::text, '[^0-9]', '', 'g'), '') WHERE environmental_fee::text LIKE '%.%.%'`,

		// Safe cast to numeric / double precision
		`ALTER TABLE vehicles ALTER COLUMN bi TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(bi::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE vehicles ALTER COLUMN tai_trong TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(tai_trong::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE vehicle_catalogs ALTER COLUMN tai_trong TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(tai_trong::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE catalog_materials ALTER COLUMN density TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(density::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE catalog_materials ALTER COLUMN ty_trong TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(ty_trong::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE catalog_materials ALTER COLUMN dinh_muc TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(dinh_muc::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE catalog_materials ALTER COLUMN price TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(price::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE catalog_materials ALTER COLUMN gia TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(gia::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_inbound ALTER COLUMN qty TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(qty::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_inbound ALTER COLUMN quantity TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(quantity::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_outbound ALTER COLUMN qty TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(qty::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_outbound ALTER COLUMN quantity TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(quantity::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_stocktake ALTER COLUMN book TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(book::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_stocktake ALTER COLUMN actual TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(actual::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_stocktake ALTER COLUMN diff TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(diff::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_stocktake ALTER COLUMN quantity TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(quantity::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_stocktake ALTER COLUMN volume TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(volume::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_stocktake ALTER COLUMN erp TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(erp::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_stocktake ALTER COLUMN survey TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(survey::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_movements ALTER COLUMN qty TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(qty::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE inventory_movements ALTER COLUMN quantity TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(quantity::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE statutory_reports ALTER COLUMN mined_volume TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(mined_volume::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE statutory_reports ALTER COLUMN tax_amount TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(tax_amount::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE statutory_reports ALTER COLUMN env_fee_amount TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(env_fee_amount::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE production_stages ALTER COLUMN volume_month TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(volume_month::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE production_stages ALTER COLUMN volume_ytd TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(volume_ytd::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE resource_taxes ALTER COLUMN environmental_fee TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(environmental_fee::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE reports ALTER COLUMN plan TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(plan::text, '[^0-9.-]', '', 'g'), '')::double precision)`,
		`ALTER TABLE reports ALTER COLUMN actual TYPE DOUBLE PRECISION USING (NULLIF(regexp_replace(actual::text, '[^0-9.-]', '', 'g'), '')::double precision)`,

		// Normalize DVT and Unit to 'tấn'
		`UPDATE catalog_materials SET dvt = 'tấn', unit = 'tấn' WHERE dvt ILIKE '%m3%' OR dvt ILIKE '%m³%' OR dvt IS NULL`,
		`UPDATE mining_plans SET unit = 'tấn' WHERE unit ILIKE '%m3%' OR unit ILIKE '%m³%' OR unit IS NULL`,
		`UPDATE hr_kpi_targets SET unit = 'tấn' WHERE unit ILIKE '%m3%' OR unit ILIKE '%m³%'`,
		`UPDATE hr_okr_key_results SET unit = 'tấn' WHERE unit ILIKE '%m3%' OR unit ILIKE '%m³%'`,

		// ===== New modules from meeting requirements =====

		// Production costs - Chi phí theo sản lượng
		`CREATE TABLE IF NOT EXISTS production_costs (
			id SERIAL PRIMARY KEY,
			cost_type TEXT NOT NULL,
			cost_category TEXT,
			norm_value DOUBLE PRECISION DEFAULT 0,
			norm_unit TEXT,
			actual_value DOUBLE PRECISION DEFAULT 0,
			actual_unit TEXT,
			period TEXT,
			mine_area TEXT,
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Cost norms - Định mức chi phí
		`CREATE TABLE IF NOT EXISTS cost_norms (
			id SERIAL PRIMARY KEY,
			norm_name TEXT NOT NULL,
			norm_type TEXT,
			unit_cost DOUBLE PRECISION DEFAULT 0,
			unit TEXT,
			material_type TEXT,
			vehicle_type TEXT,
			effective_date TEXT,
			status TEXT DEFAULT 'active',
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Vehicle trips - Chuyến xe Camera AI
		`CREATE TABLE IF NOT EXISTS vehicle_trips (
			id SERIAL PRIMARY KEY,
			vehicle_id TEXT,
			license_plate TEXT,
			driver_name TEXT,
			camera_id TEXT,
			direction TEXT,
			check_in_time TIMESTAMPTZ,
			check_out_time TIMESTAMPTZ,
			trip_number INTEGER,
			estimated_quantity DOUBLE PRECISION DEFAULT 0,
			actual_quantity DOUBLE PRECISION DEFAULT 0,
			image_evidence TEXT,
			ai_confidence DOUBLE PRECISION DEFAULT 0,
			gps_start_point TEXT,
			gps_end_point TEXT,
			status TEXT DEFAULT 'completed',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Delivery confirmations - Xác nhận giao hàng 3 nguồn
		`CREATE TABLE IF NOT EXISTS delivery_confirmations (
			id SERIAL PRIMARY KEY,
			order_code TEXT,
			contract_code TEXT,
			customer_name TEXT,
			vehicle_plate TEXT,
			driver_name TEXT,
			product_type TEXT,
			quantity_ordered DOUBLE PRECISION DEFAULT 0,
			quantity_delivered DOUBLE PRECISION DEFAULT 0,
			quantity_confirmed DOUBLE PRECISION DEFAULT 0,
			source_order TEXT,
			source_scale TEXT,
			source_warehouse TEXT,
			confirmation_status TEXT DEFAULT 'pending',
			confirmed_by TEXT,
			confirmed_at TIMESTAMPTZ,
			evidence_photo TEXT,
			evidence_location TEXT,
			evidence_timestamp TIMESTAMPTZ,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Accounting entries - Bút toán kế toán
		`CREATE TABLE IF NOT EXISTS accounting_entries (
			id SERIAL PRIMARY KEY,
			entry_code TEXT UNIQUE,
			entry_date TEXT,
			entry_type TEXT,
			account_code TEXT,
			account_name TEXT,
			description TEXT,
			debit_amount DOUBLE PRECISION DEFAULT 0,
			credit_amount DOUBLE PRECISION DEFAULT 0,
			balance DOUBLE PRECISION DEFAULT 0,
			reference_type TEXT,
			reference_code TEXT,
			tax_rate DOUBLE PRECISION DEFAULT 0,
			tax_amount DOUBLE PRECISION DEFAULT 0,
			period TEXT,
			status TEXT DEFAULT 'posted',
			created_by TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Tax records - Thuế phải nộp
		`CREATE TABLE IF NOT EXISTS tax_records (
			id SERIAL PRIMARY KEY,
			tax_type TEXT NOT NULL,
			tax_code TEXT,
			period TEXT,
			taxable_amount DOUBLE PRECISION DEFAULT 0,
			tax_rate DOUBLE PRECISION DEFAULT 0,
			tax_amount DOUBLE PRECISION DEFAULT 0,
			paid_amount DOUBLE PRECISION DEFAULT 0,
			due_date TEXT,
			paid_date TEXT,
			status TEXT DEFAULT 'pending',
			authority TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Delegations - Ủy quyền
		`CREATE TABLE IF NOT EXISTS delegations (
			id SERIAL PRIMARY KEY,
			delegator_name TEXT,
			delegator_position TEXT,
			delegate_name TEXT,
			delegate_position TEXT,
			permission_type TEXT,
			scope TEXT,
			start_date TEXT,
			end_date TEXT,
			document_url TEXT,
			status TEXT DEFAULT 'active',
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Geofences - Vùng giám sát GPS
		`CREATE TABLE IF NOT EXISTS geofences (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			fence_type TEXT,
			center_lat DOUBLE PRECISION DEFAULT 0,
			center_lng DOUBLE PRECISION DEFAULT 0,
			radius DOUBLE PRECISION DEFAULT 0,
			polygon_points TEXT,
			allowed_vehicles TEXT,
			alert_on_exit BOOLEAN DEFAULT true,
			alert_on_enter BOOLEAN DEFAULT false,
			status TEXT DEFAULT 'active',
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Risk alerts - Cảnh báo rủi ro AI
		`CREATE TABLE IF NOT EXISTS risk_alerts (
			id SERIAL PRIMARY KEY,
			alert_type TEXT NOT NULL,
			severity TEXT DEFAULT 'medium',
			title TEXT,
			description TEXT,
			module_source TEXT,
			entity_type TEXT,
			entity_id TEXT,
			ai_confidence DOUBLE PRECISION DEFAULT 0,
			suggested_action TEXT,
			assigned_to TEXT,
			status TEXT DEFAULT 'open',
			resolved_at TIMESTAMPTZ,
			resolved_by TEXT,
			resolution_notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Authorization & E-signature - Ký số & Ủy quyền
		`CREATE TABLE IF NOT EXISTS authorizations (
			id SERIAL PRIMARY KEY,
			auth_type TEXT,
			document_type TEXT,
			document_code TEXT,
			authorizer_name TEXT,
			authorizer_position TEXT,
			signer_name TEXT,
			signer_position TEXT,
			signature_method TEXT,
			signed_at TIMESTAMPTZ,
			valid_from TEXT,
			valid_to TEXT,
			status TEXT DEFAULT 'pending',
			tx_hash TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Inventory Products - Sản phẩm kho đá / Master Data
		`CREATE TABLE IF NOT EXISTS inventory_products (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			name TEXT NOT NULL,
			category TEXT,
			unit TEXT DEFAULT 'm³',
			density DOUBLE PRECISION DEFAULT 1.5,
			sale_price DOUBLE PRECISION DEFAULT 0,
			purchase_price DOUBLE PRECISION DEFAULT 0,
			storage_loc TEXT,
			standard TEXT,
			min_stock DOUBLE PRECISION DEFAULT 0,
			current_stock DOUBLE PRECISION DEFAULT 0,
			status TEXT DEFAULT 'active',
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Purchase Vouchers - Phiếu mua hàng
		`CREATE TABLE IF NOT EXISTS purchase_vouchers (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			supplier_code TEXT,
			supplier_name TEXT,
			date TEXT,
			warehouse_loc TEXT,
			total_amount DOUBLE PRECISION DEFAULT 0,
			vat_amount DOUBLE PRECISION DEFAULT 0,
			grand_total DOUBLE PRECISION DEFAULT 0,
			payment_status TEXT DEFAULT 'unpaid',
			status TEXT DEFAULT 'completed',
			created_by TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Purchase Voucher Items - Chi tiết dòng sản phẩm mua
		`CREATE TABLE IF NOT EXISTS purchase_voucher_items (
			id SERIAL PRIMARY KEY,
			voucher_id INT,
			voucher_code TEXT,
			product_code TEXT,
			product_name TEXT,
			unit TEXT,
			density DOUBLE PRECISION DEFAULT 1.5,
			unit_price DOUBLE PRECISION DEFAULT 0,
			quantity DOUBLE PRECISION DEFAULT 0,
			total_amount DOUBLE PRECISION DEFAULT 0,
			weight_ton DOUBLE PRECISION DEFAULT 0,
			standard TEXT,
			storage_loc TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Sales Vouchers - Phiếu bán hàng
		`CREATE TABLE IF NOT EXISTS sales_vouchers (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			customer_code TEXT,
			customer_name TEXT,
			date TEXT,
			warehouse_loc TEXT,
			license_plate TEXT,
			ticket_code TEXT,
			total_amount DOUBLE PRECISION DEFAULT 0,
			vat_amount DOUBLE PRECISION DEFAULT 0,
			grand_total DOUBLE PRECISION DEFAULT 0,
			paid_amount DOUBLE PRECISION DEFAULT 0,
			debt_amount DOUBLE PRECISION DEFAULT 0,
			payment_status TEXT DEFAULT 'paid',
			status TEXT DEFAULT 'completed',
			created_by TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Sales Voucher Items - Chi tiết dòng sản phẩm bán
		`CREATE TABLE IF NOT EXISTS sales_voucher_items (
			id SERIAL PRIMARY KEY,
			voucher_id INT,
			voucher_code TEXT,
			product_code TEXT,
			product_name TEXT,
			unit TEXT,
			density DOUBLE PRECISION DEFAULT 1.5,
			unit_price DOUBLE PRECISION DEFAULT 0,
			quantity DOUBLE PRECISION DEFAULT 0,
			total_amount DOUBLE PRECISION DEFAULT 0,
			weight_ton DOUBLE PRECISION DEFAULT 0,
			standard TEXT,
			storage_loc TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Return Vouchers - Phiếu trả hàng
		`CREATE TABLE IF NOT EXISTS return_vouchers (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			return_type TEXT DEFAULT 'sales_return',
			partner_code TEXT,
			partner_name TEXT,
			ref_voucher_code TEXT,
			date TEXT,
			warehouse_loc TEXT,
			reason TEXT,
			total_amount DOUBLE PRECISION DEFAULT 0,
			status TEXT DEFAULT 'approved',
			created_by TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Return Voucher Items - Chi tiết dòng sản phẩm trả
		`CREATE TABLE IF NOT EXISTS return_voucher_items (
			id SERIAL PRIMARY KEY,
			voucher_id INT,
			voucher_code TEXT,
			product_code TEXT,
			product_name TEXT,
			unit TEXT,
			unit_price DOUBLE PRECISION DEFAULT 0,
			quantity DOUBLE PRECISION DEFAULT 0,
			total_amount DOUBLE PRECISION DEFAULT 0,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Receipt Vouchers - Phiếu thu
		`CREATE TABLE IF NOT EXISTS receipt_vouchers (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			partner_code TEXT,
			partner_name TEXT,
			partner_type TEXT DEFAULT 'customer',
			payer_name TEXT,
			reason TEXT,
			ref_code TEXT,
			amount DOUBLE PRECISION DEFAULT 0,
			amount_in_words TEXT,
			payment_method TEXT DEFAULT 'bank_transfer',
			fund_account TEXT,
			status TEXT DEFAULT 'posted',
			date TEXT,
			created_by TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Payment Vouchers - Phiếu chi
		`CREATE TABLE IF NOT EXISTS payment_vouchers (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			partner_code TEXT,
			partner_name TEXT,
			partner_type TEXT DEFAULT 'supplier',
			receiver_name TEXT,
			reason TEXT,
			ref_code TEXT,
			amount DOUBLE PRECISION DEFAULT 0,
			amount_in_words TEXT,
			payment_method TEXT DEFAULT 'bank_transfer',
			fund_account TEXT,
			status TEXT DEFAULT 'posted',
			date TEXT,
			created_by TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Sales Contracts - Hợp đồng thương mại bán đá cam kết sản lượng
		`CREATE TABLE IF NOT EXISTS sales_contracts (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			customer_code TEXT,
			customer_name TEXT,
			project_name TEXT,
			start_date TEXT,
			end_date TEXT,
			total_committed_tons DOUBLE PRECISION DEFAULT 0,
			total_delivered_tons DOUBLE PRECISION DEFAULT 0,
			completion_rate DOUBLE PRECISION DEFAULT 0,
			credit_limit DOUBLE PRECISION DEFAULT 0,
			payment_terms TEXT DEFAULT 'prepaid_wallet',
			status TEXT DEFAULT 'active',
			notes TEXT,
			created_by TEXT,
			items JSONB DEFAULT '[]',
			whitelisted_plates JSONB DEFAULT '[]',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS sales_contract_items (
			id SERIAL PRIMARY KEY,
			contract_id INT,
			contract_code TEXT,
			product_code TEXT,
			product_name TEXT,
			committed_tons DOUBLE PRECISION DEFAULT 0,
			delivered_tons DOUBLE PRECISION DEFAULT 0,
			unit_price DOUBLE PRECISION DEFAULT 0,
			discount_pct DOUBLE PRECISION DEFAULT 0,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Consolidated Delivery Orders - Gom phiếu cân xuất kho chính thức & xuất hóa đơn điện tử
		`CREATE TABLE IF NOT EXISTS consolidated_delivery_orders (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE,
			customer_code TEXT,
			customer_name TEXT,
			contract_code TEXT,
			period_start TEXT,
			period_end TEXT,
			total_trips INT DEFAULT 0,
			total_tons DOUBLE PRECISION DEFAULT 0,
			total_amount DOUBLE PRECISION DEFAULT 0,
			vat_amount DOUBLE PRECISION DEFAULT 0,
			grand_total DOUBLE PRECISION DEFAULT 0,
			einvoice_no TEXT,
			einvoice_status TEXT DEFAULT 'pending',
			einvoice_lookup_url TEXT,
			status TEXT DEFAULT 'draft',
			ticket_codes JSONB DEFAULT '[]',
			created_by TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Blasting Passports extended columns (QCVN 01:2019/BCT)
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS bench_level TEXT`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS hole_diameter_mm DOUBLE PRECISION DEFAULT 105`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS burden_m DOUBLE PRECISION DEFAULT 3.0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS spacing_m DOUBLE PRECISION DEFAULT 3.5`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS row_spacing_m DOUBLE PRECISION DEFAULT 3.0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS stemming_length_m DOUBLE PRECISION DEFAULT 3.2`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS rock_hardness_f TEXT DEFAULT 'f = 10 - 12'`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS firing_method TEXT DEFAULT 'Vi sai phi điện MS'`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS delay_interval_ms INTEGER DEFAULT 25`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS status_code TEXT DEFAULT 'APPROVED'`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS certificate_expiry TEXT`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS approved_by TEXT`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS approved_at TEXT`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS license_number TEXT`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS quota_annual_limit_tons DOUBLE PRECISION DEFAULT 180`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS quota_used_ytd_tons DOUBLE PRECISION DEFAULT 0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS quota_remaining_tons DOUBLE PRECISION DEFAULT 180`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS safety_perimeter_m INTEGER DEFAULT 300`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS police_notified BOOLEAN DEFAULT true`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS commune_notified BOOLEAN DEFAULT true`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS notification_doc_ref TEXT`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS all_guards_confirmed BOOLEAN DEFAULT true`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS evacuation_confirmed BOOLEAN DEFAULT true`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS siren_alerts_completed BOOLEAN DEFAULT true`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS guard_posts JSONB DEFAULT '[]'`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS materials JSONB DEFAULT '[]'`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS same_day_return_closed BOOLEAN DEFAULT true`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS return_completed_time TEXT`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS post_blast_clearance BOOLEAN DEFAULT true`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS smoke_clearing_minutes INTEGER DEFAULT 20`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS misfire_reported BOOLEAN DEFAULT false`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS misfire_count INTEGER DEFAULT 0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS misfire_details TEXT`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS misfire_resolution TEXT`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS explosive_cost DOUBLE PRECISION DEFAULT 0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS drilling_cost DOUBLE PRECISION DEFAULT 0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS blasting_service_fee DOUBLE PRECISION DEFAULT 0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS guard_labor_cost DOUBLE PRECISION DEFAULT 0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS total_blast_cost DOUBLE PRECISION DEFAULT 0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS cost_per_m3_rock DOUBLE PRECISION DEFAULT 0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS cost_per_ton_rock DOUBLE PRECISION DEFAULT 0`,
		`ALTER TABLE blasting_passports ADD COLUMN IF NOT EXISTS notes TEXT`,
		`ALTER TABLE equipment_fuel_logs ADD COLUMN IF NOT EXISTS fuel_dispense_logs JSONB DEFAULT '[]'`,
		`ALTER TABLE equipment_fuel_logs ADD COLUMN IF NOT EXISTS tank_capacity_liters DOUBLE PRECISION DEFAULT 400`,
		`ALTER TABLE equipment_fuel_logs ADD COLUMN IF NOT EXISTS current_fuel_liters DOUBLE PRECISION DEFAULT 280`,
		`ALTER TABLE equipment_fuel_logs ADD COLUMN IF NOT EXISTS last_dispense_at TEXT DEFAULT '07:30 28/10/2026'`,
		`ALTER TABLE equipment_fuel_logs ADD COLUMN IF NOT EXISTS next_maintenance_hours DOUBLE PRECISION DEFAULT 5500`,
		`ALTER TABLE equipment_fuel_logs ADD COLUMN IF NOT EXISTS engine_specs JSONB DEFAULT '{}'`,
	}

	for _, a := range alters {
		stepCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		Pool.Exec(stepCtx, a)
		cancel()
	}

	MigrateQuarryModule()

	fmt.Println("Database migration completed successfully!")
}

