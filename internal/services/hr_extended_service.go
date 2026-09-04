package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/repositories"
)

// generic wrapper service exposing the underlying base repo via .Base
type GenericService struct {
	Base *repositories.BaseRepo
}

func newGeneric(table string) *GenericService {
	return &GenericService{Base: repositories.NewBaseRepo(table, "id")}
}

func NewHrWorkHistoryService() *GenericService     { return newGeneric("hr_employee_work_history") }
func NewHrOnboardingService() *GenericService      { return newGeneric("hr_onboarding_checklists") }
func NewHrDecisionService() *GenericService        { return newGeneric("hr_decisions") }
func NewHrEsignDocumentService() *GenericService   { return newGeneric("hr_esign_documents") }
func NewHrEsignSignerService() *GenericService     { return newGeneric("hr_esign_signers") }
func NewHrInsuranceRecordService() *GenericService { return newGeneric("hr_insurance_records") }
func NewHrInsuranceDeclarationService() *GenericService {
	return newGeneric("hr_insurance_declarations")
}
func NewHrRecruitmentRequestService() *GenericService { return newGeneric("hr_recruitment_requests") }
func NewHrCampaignService() *GenericService           { return newGeneric("hr_recruitment_campaigns") }
func NewHrCandidateService() *GenericService          { return newGeneric("hr_candidates") }
func NewHrCandidateApplicationService() *GenericService {
	return newGeneric("hr_candidate_applications")
}
func NewHrCandidateStatusHistoryService() *GenericService {
	return newGeneric("hr_candidate_status_history")
}
func NewHrInterviewScheduleService() *GenericService { return newGeneric("hr_interview_schedules") }
func NewHrTrainingCourseService() *GenericService    { return newGeneric("hr_training_courses") }
func NewHrTrainingSessionService() *GenericService   { return newGeneric("hr_training_sessions") }
func NewHrTrainingParticipantService() *GenericService {
	return newGeneric("hr_training_participants")
}
func NewHrTrainingAttendanceService() *GenericService { return newGeneric("hr_training_attendances") }
func NewHrTrainingEvaluationService() *GenericService { return newGeneric("hr_training_evaluations") }
func NewHrKpiTargetService() *GenericService          { return newGeneric("hr_kpi_targets") }
func NewHrOkrObjectiveService() *GenericService       { return newGeneric("hr_okr_objectives") }
func NewHrOkrKeyResultService() *GenericService       { return newGeneric("hr_okr_key_results") }
func NewHrWorkflowTemplateService() *GenericService   { return newGeneric("hr_workflow_templates") }
func NewHrWorkflowInstanceService() *GenericService   { return newGeneric("hr_workflow_instances") }
func NewHrWorkflowStepService() *GenericService       { return newGeneric("hr_workflow_steps") }
func NewHrShiftAssignmentService() *GenericService    { return newGeneric("hr_shift_assignments") }
func NewHrHolidayService() *GenericService            { return newGeneric("hr_holidays") }
func NewHrAttendanceLogService() *GenericService      { return newGeneric("hr_attendance_logs") }
func NewHrTimesheetDetailService() *GenericService    { return newGeneric("hr_timesheet_details") }
func NewHrPayrollPeriodService() *GenericService      { return newGeneric("hr_payroll_periods") }
func NewHrPayrollRecordService() *GenericService      { return newGeneric("hr_payroll_records") }
func NewHrPayrollFormulaColumnService() *GenericService {
	return newGeneric("hr_payroll_formula_columns")
}
func NewHrSalaryAdvanceService() *GenericService { return newGeneric("hr_salary_advances") }

func CountRows(table, where string, args ...interface{}) (int, error) {
	return repositories.NewBaseRepo(table, "id").CountRows(table, where, args...)
}

// ---------------------- Employee lifecycle ----------------------

type EmployeeStatusService struct{}

func (s *EmployeeStatusService) ChangeStatus(employeeID, toStatus, effectiveDate, reason, changedBy string) (map[string]interface{}, error) {
	ctx := context.Background()

	now := time.Now().Format("02/01/2006 15:04")
	// Update employee status
	if _, err := database.Pool.Exec(ctx,
		`UPDATE hr_employees SET status=$1, updated_at=NOW() WHERE id=$2`, toStatus, employeeID); err != nil {
		return nil, err
	}

	var name string
	database.Pool.QueryRow(ctx, `SELECT name FROM hr_employees WHERE id=$1`, employeeID).Scan(&name)

	// Record work history
	eventType := toStatus
	switch toStatus {
	case "dang_lam_viec":
		eventType = "bat_dau_lam_viec"
	case "nghi_tam_thoi":
		eventType = "nghi_tam_thoi"
	case "nghi_viec":
		eventType = "nghi_viec"
	}
	title := toStatus
	detail := reason
	if effectiveDate != "" {
		detail = "Hiệu lực từ " + effectiveDate + ". " + reason
	}
	database.Pool.Exec(ctx,
		`INSERT INTO hr_employee_work_history (id, employee_id, employee_name, event_type, title, detail, changed_by, changed_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		"WH-"+nowid(), employeeID, name, eventType, title, detail, changedBy, now)

	// Side-effects on resignation
	if toStatus == "nghi_viec" {
		database.Pool.Exec(ctx, `UPDATE hr_contracts SET status='het_hieu_luc' WHERE employee_id=$1 AND status='hiệu lực'`, employeeID)
		database.Pool.Exec(ctx, `UPDATE hr_insurance_records SET status='chua_tham_gia', stop_date=$1 WHERE employee_id=$2`, effectiveDate, employeeID)
		database.Pool.Exec(ctx,
			`INSERT INTO hr_insurance_declarations (id, code, employee_id, employee_name, declaration_type, period, status, notes)
			 VALUES ($1,$2,$3,$4,'giam',$5,'cho_ke_khai',$6) ON CONFLICT (id) DO NOTHING`,
			"INSDEC-"+nowid(), "BAOGIAM-"+nowid(), employeeID, name, currentPeriod(), reason)
	}

	// Return updated employee
	rows, err := database.Pool.Query(ctx, `SELECT row_to_json(t) FROM (SELECT * FROM hr_employees WHERE id=$1) t`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var item map[string]interface{}
	if rows.Next() {
		var data []byte
		rows.Scan(&data)
		jsonUnmarshal(data, &item)
	}
	return item, nil
}

// ---------------------- Onboarding ----------------------

type OnboardingService struct{}

func (s *OnboardingService) GenerateChecklist(employeeID string) error {
	ctx := context.Background()
	var name string
	database.Pool.QueryRow(ctx, `SELECT name FROM hr_employees WHERE id=$1`, employeeID).Scan(&name)

	tasks := []string{
		"Ảnh thẻ / Ảnh cá nhân 3x4",
		"Bản sao Giấy khai sinh",
		"Bản sao Chứng minh nhân dân / Căn cước công dân",
		"Bản sao Bằng cấp / Chứng chỉ",
		"Sổ hộ khẩu / Giấy xác nhận cư trú",
		"Giấy khám sức khỏe nghề nghiệp",
		"Sổ tiết kiệm / Tài khoản ngân hàng nhận lương",
	}
	for i, task := range tasks {
		database.Pool.Exec(ctx,
			`INSERT INTO hr_onboarding_checklists (id, employee_id, employee_name, task_name, is_completed, status)
			 VALUES ($1,$2,$3,$4,FALSE,'cho_bo_sung') ON CONFLICT (id) DO NOTHING`,
			fmt.Sprintf("OBC-%s-%d", nowid(), i), employeeID, name, task)
	}
	return nil
}

func (s *OnboardingService) CompleteItem(itemID string) error {
	ctx := context.Background()
	_, err := database.Pool.Exec(ctx,
		`UPDATE hr_onboarding_checklists SET is_completed=TRUE, status='da_bo_sung', completed_at=NOW() WHERE id=$1`, itemID)
	return err
}

// ---------------------- Decisions ----------------------

type DecisionService struct{}

func (s *DecisionService) Issue(id string) (map[string]interface{}, error) {
	ctx := context.Background()
	var decision map[string]interface{}
	var data []byte
	err := database.Pool.QueryRow(ctx, `SELECT row_to_json(t) FROM (SELECT * FROM hr_decisions WHERE id=$1) t`, id).Scan(&data)
	if err != nil {
		return nil, err
	}
	jsonUnmarshal(data, &decision)

	if _, err := database.Pool.Exec(ctx, `UPDATE hr_decisions SET status='da_ban_hanh', updated_at=NOW() WHERE id=$1`, id); err != nil {
		return nil, err
	}

	empID, _ := decision["employee_id"].(string)
	empName, _ := decision["employee_name"].(string)
	now := time.Now().Format("02/01/2006 15:04")
	decisionType, _ := decision["decision_type"].(string)
	content, _ := decision["content_detail"].(map[string]interface{})

	database.Pool.Exec(ctx,
		`INSERT INTO hr_employee_work_history (id, employee_id, employee_name, event_type, title, detail, changed_by, changed_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		"WH-"+nowid(), empID, empName, decisionType, decisionType, fmt.Sprintf("%v", content), userOrDefault(), now)

	switch decisionType {
	case "tang_luong":
		if newSalary, ok := content["luong_moi"].(float64); ok {
			database.Pool.Exec(ctx,
				`UPDATE hr_employees SET base_salary=$1, updated_at=NOW() WHERE id=$2`, newSalary, empID)
		}
	case "dieu_chuyen":
		if dept, ok := content["phong_ban_moi"].(string); ok {
			database.Pool.Exec(ctx,
				`UPDATE hr_employees SET department=$1, updated_at=NOW() WHERE id=$2`, dept, empID)
		}
	}

	var updated []byte
	database.Pool.QueryRow(ctx, `SELECT row_to_json(t) FROM (SELECT * FROM hr_decisions WHERE id=$1) t`, id).Scan(&updated)
	var out map[string]interface{}
	jsonUnmarshal(updated, &out)
	return out, nil
}

// ---------------------- E-sign ----------------------

type EsignService struct{}

func (s *EsignService) CreateForContract(contractID, contractCode, subject, employeeID, employeeName, hrName string) (map[string]interface{}, error) {
	ctx := context.Background()
	docID := "ESIGN-" + nowid()
	if _, err := database.Pool.Exec(ctx,
		`INSERT INTO hr_esign_documents (id, code, contract_id, contract_code, subject, status)
		 VALUES ($1,$2,$3,$4,$5,'cho_ky')`,
		docID, "ES-"+nowid(), contractID, contractCode, subject); err != nil {
		return nil, err
	}
	// Two signers: employee order 1, HR order 2
	database.Pool.Exec(ctx,
		`INSERT INTO hr_esign_signers (id, document_id, signer_type, employee_id, employee_name, sign_order, sign_status)
		 VALUES ($1,$2,'nhan_vien',$3,$4,1,'cho_ky')`,
		"SG-"+nowid(), docID, employeeID, employeeName)
	database.Pool.Exec(ctx,
		`INSERT INTO hr_esign_signers (id, document_id, signer_type, employee_id, employee_name, sign_order, sign_status)
		 VALUES ($1,$2,'quan_ly_nhan_su',NULL,$3,2,'cho_ky')`,
		"SG-"+nowid(), docID, hrName)

	var data []byte
	database.Pool.QueryRow(ctx, `SELECT row_to_json(t) FROM (SELECT * FROM hr_esign_documents WHERE id=$1) t`, docID).Scan(&data)
	var out map[string]interface{}
	jsonUnmarshal(data, &out)
	return out, nil
}

func (s *EsignService) Sign(signerID, status, note string) error {
	ctx := context.Background()
	now := time.Now().Format("02/01/2006 15:04")
	if _, err := database.Pool.Exec(ctx,
		`UPDATE hr_esign_signers SET sign_status=$1, sign_date=$2, sign_note=$3 WHERE id=$4`, status, now, note, signerID); err != nil {
		return err
	}
	var docID string
	database.Pool.QueryRow(ctx, `SELECT document_id FROM hr_esign_signers WHERE id=$1`, signerID).Scan(&docID)

	// Check all signers resolved
	var pending int
	database.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM hr_esign_signers WHERE document_id=$1 AND sign_status='cho_ky'`, docID).Scan(&pending)
	var rejected int
	database.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM hr_esign_signers WHERE document_id=$1 AND sign_status='tu_choi'`, docID).Scan(&rejected)

	docStatus := "cho_ky"
	if rejected > 0 {
		docStatus = "tu_choi"
	} else if pending == 0 {
		docStatus = "da_ky"
		var contractID string
		database.Pool.QueryRow(ctx, `SELECT contract_id FROM hr_esign_documents WHERE id=$1`, docID).Scan(&contractID)
		database.Pool.Exec(ctx, `UPDATE hr_contracts SET status='dang_hieu_luc' WHERE id=$1`, contractID)
	}
	_, err := database.Pool.Exec(ctx, `UPDATE hr_esign_documents SET status=$1 WHERE id=$2`, docStatus, docID)
	return err
}

func (s *EsignService) SignDocument(docIDOrCode, signerID, signatureData, note string) (map[string]interface{}, error) {
	ctx := context.Background()
	now := time.Now().Format("02/01/2006 15:04")

	var docID string
	err := database.Pool.QueryRow(ctx,
		`SELECT id FROM hr_esign_documents WHERE id=$1 OR code=$1 LIMIT 1`, docIDOrCode).Scan(&docID)
	if err != nil {
		docID = docIDOrCode
	}

	var targetSignerID string
	if signerID != "" && signerID != "admin" {
		_ = database.Pool.QueryRow(ctx,
			`SELECT id FROM hr_esign_signers WHERE document_id=$1 AND (id=$2 OR employee_id=$2) AND sign_status='cho_ky' LIMIT 1`,
			docID, signerID).Scan(&targetSignerID)
	}
	if targetSignerID == "" {
		_ = database.Pool.QueryRow(ctx,
			`SELECT id FROM hr_esign_signers WHERE document_id=$1 AND sign_status='cho_ky' ORDER BY sign_order ASC LIMIT 1`,
			docID).Scan(&targetSignerID)
	}

	if targetSignerID != "" {
		_, _ = database.Pool.Exec(ctx,
			`UPDATE hr_esign_signers SET sign_status='da_ky', sign_date=$1, sign_note=$2 WHERE id=$3`,
			now, note, targetSignerID)
	}

	var pending int
	_ = database.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM hr_esign_signers WHERE document_id=$1 AND sign_status='cho_ky'`, docID).Scan(&pending)
	var rejected int
	_ = database.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM hr_esign_signers WHERE document_id=$1 AND sign_status='tu_choi'`, docID).Scan(&rejected)

	docStatus := "cho_ky"
	var currentDocStatus string
	_ = database.Pool.QueryRow(ctx, `SELECT status FROM hr_esign_documents WHERE id=$1`, docID).Scan(&currentDocStatus)
	if currentDocStatus == "da_chuyen" {
		docStatus = "da_chuyen"
	}
	if rejected > 0 {
		docStatus = "tu_choi"
	} else if pending == 0 {
		docStatus = "da_ky"
		var contractID string
		_ = database.Pool.QueryRow(ctx, `SELECT contract_id FROM hr_esign_documents WHERE id=$1`, docID).Scan(&contractID)
		if contractID != "" {
			_, _ = database.Pool.Exec(ctx, `UPDATE hr_contracts SET status='dang_hieu_luc' WHERE id=$1`, contractID)
		}
	}

	_, _ = database.Pool.Exec(ctx, `UPDATE hr_esign_documents SET status=$1 WHERE id=$2`, docStatus, docID)

	var data []byte
	_ = database.Pool.QueryRow(ctx, `SELECT row_to_json(t) FROM (SELECT * FROM hr_esign_documents WHERE id=$1) t`, docID).Scan(&data)
	var out map[string]interface{}
	jsonUnmarshal(data, &out)
	return out, nil
}

func (s *EsignService) DelegateDocument(docIDOrCode, fromSignerID, toSignerID, reason string) (map[string]interface{}, error) {
	ctx := context.Background()
	now := time.Now().Format("02/01/2006 15:04")

	var docID string
	err := database.Pool.QueryRow(ctx,
		`SELECT id FROM hr_esign_documents WHERE id=$1 OR code=$1 LIMIT 1`, docIDOrCode).Scan(&docID)
	if err != nil {
		docID = docIDOrCode
	}

	var targetSignerID string
	_ = database.Pool.QueryRow(ctx,
		`SELECT id FROM hr_esign_signers WHERE document_id=$1 AND sign_status='cho_ky' ORDER BY sign_order ASC LIMIT 1`,
		docID).Scan(&targetSignerID)

	if targetSignerID != "" {
		_, _ = database.Pool.Exec(ctx,
			`UPDATE hr_esign_signers SET sign_status='da_chuyen', sign_date=$1, sign_note=$2 WHERE id=$3`,
			now, reason, targetSignerID)
	}

	_, _ = database.Pool.Exec(ctx, `UPDATE hr_esign_documents SET status='da_chuyen' WHERE id=$1`, docID)

	var data []byte
	_ = database.Pool.QueryRow(ctx, `SELECT row_to_json(t) FROM (SELECT * FROM hr_esign_documents WHERE id=$1) t`, docID).Scan(&data)
	var out map[string]interface{}
	jsonUnmarshal(data, &out)
	return out, nil
}

// ---------------------- Insurance ----------------------

type InsuranceService struct{}

func (s *InsuranceService) CreateDeclarationForEmployee(employeeID, declType, period string) error {
	ctx := context.Background()
	var name string
	database.Pool.QueryRow(ctx, `SELECT name FROM hr_employees WHERE id=$1`, employeeID).Scan(&name)
	_, err := database.Pool.Exec(ctx,
		`INSERT INTO hr_insurance_declarations (id, code, employee_id, employee_name, declaration_type, period, status, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,'cho_ke_khai','') ON CONFLICT (id) DO NOTHING`,
		"INSDEC-"+nowid(), "BK"+declType+"-"+nowid(), employeeID, name, declType, period)
	return err
}

// ---------------------- Training ----------------------

type TrainingService struct{}

func (s *TrainingService) RecalculateCompletion(courseID string) error {
	ctx := context.Background()
	var threshold float64 = 80
	database.Pool.QueryRow(ctx, `SELECT attendance_threshold FROM hr_training_courses WHERE id=$1`, courseID).Scan(&threshold)

	rows, err := database.Pool.Query(ctx,
		`SELECT id, employee_id, employee_name, attendance_rate FROM hr_training_participants WHERE course_id=$1`, courseID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var pid, eid, ename string
		var rate float64
		rows.Scan(&pid, &eid, &ename, &rate)
		status := "khong_hoan_thanh"
		if rate >= threshold {
			status = "hoan_thanh"
		}
		database.Pool.Exec(ctx, `UPDATE hr_training_participants SET completion_status=$1 WHERE id=$2`, status, pid)
		database.Pool.Exec(ctx,
			`INSERT INTO hr_employee_work_history (id, employee_id, employee_name, event_type, title, detail, changed_at)
			 VALUES ($1,$2,$3,'dao_tao','Hoàn thành khoá đào tạo', $4, NOW())`,
			"WH-"+nowid(), eid, ename, fmt.Sprintf("Tỷ lệ điểm danh %v%%", rate))
	}
	// Update course status
	var incomplete int
	database.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM hr_training_participants WHERE course_id=$1 AND completion_status='khong_hoan_thanh'`, courseID).Scan(&incomplete)
	if incomplete == 0 {
		database.Pool.Exec(ctx, `UPDATE hr_training_courses SET status='hoan_thanh' WHERE id=$1`, courseID)
	} else {
		database.Pool.Exec(ctx, `UPDATE hr_training_courses SET status='khong_hoan_thanh' WHERE id=$1`, courseID)
	}
	return nil
}

// ---------------------- KPI / OKR ----------------------

type KpiService struct{}

func (s *KpiService) RecalculateEmployee(employeeID, period string) (map[string]interface{}, error) {
	ctx := context.Background()
	rows, err := database.Pool.Query(ctx,
		`SELECT id, target_value, actual_value, weight_percent FROM hr_kpi_targets WHERE employee_id=$1 AND period=$2`, employeeID, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	total := 0.0
	for rows.Next() {
		var id string
		var target, actual, weight float64
		rows.Scan(&id, &target, &actual, &weight)
		score := 0.0
		if target > 0 {
			score = (actual / target) * weight
		}
		database.Pool.Exec(ctx, `UPDATE hr_kpi_targets SET score=$1, status='da_ban_hanh' WHERE id=$2`, score, id)
		total += score
	}
	return map[string]interface{}{"employee_id": employeeID, "period": period, "total_score": total}, nil
}

type OkrService struct{}

func (s *OkrService) RecalculateKeyResult(keyResultID string) (float64, error) {
	ctx := context.Background()
	var start, target, current float64
	database.Pool.QueryRow(ctx,
		`SELECT start_value, target_value, current_value FROM hr_okr_key_results WHERE id=$1`, keyResultID).Scan(&start, &target, &current)
	progress := 0.0
	if target != start {
		progress = (current - start) / (target - start) * 100
	}
	database.Pool.Exec(ctx, `UPDATE hr_okr_key_results SET progress_percent=$1 WHERE id=$2`, progress, keyResultID)
	return progress, nil
}

func (s *OkrService) RecalculateObjective(objectiveID string) (float64, error) {
	ctx := context.Background()
	rows, err := database.Pool.Query(ctx,
		`SELECT id, progress_percent FROM hr_okr_key_results WHERE objective_id=$1`, objectiveID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	sum, count := 0.0, 0
	for rows.Next() {
		var id string
		var p float64
		rows.Scan(&id, &p)
		sum += p
		count++
	}
	avg := 0.0
	if count > 0 {
		avg = sum / float64(count)
	}
	database.Pool.Exec(ctx, `UPDATE hr_okr_objectives SET progress_percent=$1 WHERE id=$2`, avg, objectiveID)
	return avg, nil
}

// ---------------------- Workflow engine ----------------------

type WorkflowService struct{}

func (s *WorkflowService) Start(appliesTo, templateID, objectID, objectRef, objectOwner string) (map[string]interface{}, error) {
	ctx := context.Background()
	instID := "WF-" + nowid()
	if _, err := database.Pool.Exec(ctx,
		`INSERT INTO hr_workflow_instances (id, code, template_id, applies_to, object_id, object_ref, status, current_step, started_at)
		 VALUES ($1,$2,$3,$4,$5,$6,'dang_chay',0,NOW())`,
		instID, "WF-"+nowid(), templateID, appliesTo, objectID, objectRef); err != nil {
		return nil, err
	}
	// Load template definition & create steps
	var definition []byte
	database.Pool.QueryRow(ctx, `SELECT definition FROM hr_workflow_templates WHERE id=$1`, templateID).Scan(&definition)
	var def struct {
		Nodes []struct {
			ID           string `json:"id"`
			Type         string `json:"type"`
			ApproverRole string `json:"approver_role"`
		} `json:"nodes"`
	}
	jsonUnmarshal(definition, &def)
	order := 0
	for _, n := range def.Nodes {
		if n.Type != "approval" {
			continue
		}
		order++
		database.Pool.Exec(ctx,
			`INSERT INTO hr_workflow_steps (id, instance_id, step_order, step_type, approver_role, step_status)
			 VALUES ($1,$2,$3,$4,$5,'cho_duyet')`,
			"WFS-"+nowid(), instID, order, "approval", n.ApproverRole)
	}
	var data []byte
	database.Pool.QueryRow(ctx, `SELECT row_to_json(t) FROM (SELECT * FROM hr_workflow_instances WHERE id=$1) t`, instID).Scan(&data)
	var out map[string]interface{}
	jsonUnmarshal(data, &out)
	return out, nil
}

func (s *WorkflowService) ResolveStep(instanceID, decision, comment, approverID, approverName string) (map[string]interface{}, error) {
	ctx := context.Background()
	// Find current pending step
	stepID := ""
	database.Pool.QueryRow(ctx,
		`SELECT id FROM hr_workflow_steps WHERE instance_id=$1 AND step_status='cho_duyet' ORDER BY step_order LIMIT 1`, instanceID).Scan(&stepID)
	if stepID == "" {
		return nil, fmt.Errorf("no pending step")
	}
	status := "da_duyet"
	if decision == "reject" {
		status = "tu_choi"
	}
	database.Pool.Exec(ctx,
		`UPDATE hr_workflow_steps SET step_status=$1, comment=$2, approver_id=$3, approver_name=$4, decided_at=NOW() WHERE id=$5`,
		status, comment, approverID, approverName, stepID)

	// Check remaining pending steps
	var pending int
	database.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM hr_workflow_steps WHERE instance_id=$1 AND step_status='cho_duyet'`, instanceID).Scan(&pending)
	var rejected int
	database.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM hr_workflow_steps WHERE instance_id=$1 AND step_status='tu_choi'`, instanceID).Scan(&rejected)

	if rejected > 0 {
		database.Pool.Exec(ctx, `UPDATE hr_workflow_instances SET status='da_tu_choi', finished_at=NOW() WHERE id=$1`, instanceID)
	} else if pending == 0 {
		database.Pool.Exec(ctx, `UPDATE hr_workflow_instances SET status='da_hoan_thanh', finished_at=NOW() WHERE id=$1`, instanceID)
	}
	var data []byte
	database.Pool.QueryRow(ctx, `SELECT row_to_json(t) FROM (SELECT * FROM hr_workflow_instances WHERE id=$1) t`, instanceID).Scan(&data)
	var out map[string]interface{}
	jsonUnmarshal(data, &out)
	return out, nil
}

// ---------------------- Recruitment ----------------------

type RecruitmentService struct{}

func (s *RecruitmentService) Submit(id string) error {
	ctx := context.Background()
	_, err := database.Pool.Exec(ctx,
		`UPDATE hr_recruitment_requests SET approval_status='cho_duyet', current_step=1, updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (s *RecruitmentService) Approve(id, comment string) error {
	ctx := context.Background()
	_, err := database.Pool.Exec(ctx,
		`UPDATE hr_recruitment_requests SET approval_status='da_duyet', comment=$1, updated_at=NOW() WHERE id=$2`, comment, id)
	return err
}

func (s *RecruitmentService) Reject(id, comment string) error {
	ctx := context.Background()
	if comment == "" {
		return fmt.Errorf("comment required on reject")
	}
	_, err := database.Pool.Exec(ctx,
		`UPDATE hr_recruitment_requests SET approval_status='tu_choi', comment=$1, updated_at=NOW() WHERE id=$2`, comment, id)
	return err
}

// ---------------------- Candidate pipeline ----------------------

type CandidateService struct{}

func (s *CandidateService) ChangeStage(candidateID, toStage, note, changedBy string) error {
	ctx := context.Background()

	// Record history
	var appID, name string
	database.Pool.QueryRow(ctx,
		`SELECT application_id, candidate_id FROM hr_candidate_status_history WHERE candidate_id=$1 ORDER BY created_at DESC LIMIT 1`, candidateID).Scan(&appID, &name)
	var cityCandidateID string
	var fromStage string
	database.Pool.QueryRow(ctx, `SELECT id, pipeline_stage FROM hr_candidates WHERE id=$1`, candidateID).Scan(&cityCandidateID, &fromStage)

	database.Pool.Exec(ctx,
		`INSERT INTO hr_candidate_status_history (id, candidate_id, from_stage, to_stage, changed_by, note, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW())`,
		"CSH-"+nowid(), candidateID, fromStage, toStage, changedBy, note)
	database.Pool.Exec(ctx,
		`UPDATE hr_candidates SET pipeline_stage=$1, note=$2, updated_at=NOW() WHERE id=$3`, toStage, note, candidateID)
	// Sync latest application if any
	database.Pool.Exec(ctx,
		`UPDATE hr_candidate_applications SET pipeline_stage=$1 WHERE candidate_id=$2`, toStage, candidateID)
	return nil
}

func (s *CandidateService) ConvertToEmployee(candidateID string) (string, error) {
	ctx := context.Background()
	var c struct {
		Name string
		Pos  string
		Dept string
		Note string
	}
	err := database.Pool.QueryRow(ctx,
		`SELECT name, position, source, note FROM hr_candidates WHERE id=$1`, candidateID).Scan(&c.Name, &c.Pos, &c.Dept, &c.Note)
	if err != nil {
		return "", err
	}
	empID := "EMP-" + nowid()
	now := time.Now()
	_, err = database.Pool.Exec(ctx,
		`INSERT INTO hr_employees (id, code, name, job_position, department, manager, join_date, contract_type, status, base_salary, created_at)
		 VALUES ($1,$2,$3,$4,$5,'',$6,'Thử việc','cho_nhan_viec',0,NOW())
		 ON CONFLICT (id) DO NOTHING`,
		empID, "NV-"+nowid(), c.Name, c.Pos, c.Dept, now.Format("02/01/2006"))
	if err != nil {
		return "", err
	}
	database.Pool.Exec(ctx,
		`UPDATE hr_candidates SET pipeline_stage='da_tao_nhan_su', note=$1 WHERE id=$2`, c.Note, candidateID)
	database.Pool.Exec(ctx,
		`UPDATE hr_candidate_applications SET pipeline_stage='da_tao_nhan_su' WHERE candidate_id=$1`, candidateID)
	database.Pool.Exec(ctx,
		`INSERT INTO hr_candidate_status_history (id, candidate_id, from_stage, to_stage, changed_by, note, created_at)
		 VALUES ($1,$2,'trung_tuyen','da_tao_nhan_su','Hệ thống','Đã tạo nhân sự ' + $3, NOW())`,
		"CSH-"+nowid(), candidateID, c.Name)
	// Generate onboarding checklist
	if err := (&OnboardingService{}).GenerateChecklist(empID); err != nil {
		return empID, nil
	}
	return empID, nil
}

// ---------------------- Interview ----------------------

type InterviewService struct{}

func (s *InterviewService) SetResult(id, result, note string) error {
	ctx := context.Background()
	_, err := database.Pool.Exec(ctx,
		`UPDATE hr_interview_schedules SET result=$1, note=$2, status='da_co_ket_qua' WHERE id=$3`, result, note, id)
	if err != nil {
		return err
	}
	// Sync candidate pipeline stage based on result
	var candidateID string
	database.Pool.QueryRow(ctx, `SELECT candidate_id FROM hr_interview_schedules WHERE id=$1`, id).Scan(&candidateID)
	if result == "dat" {
		database.Pool.Exec(ctx, `UPDATE hr_candidates SET pipeline_stage='pv_dat' WHERE id=$1`, candidateID)
		database.Pool.Exec(ctx, `UPDATE hr_candidate_applications SET pipeline_stage='pv_dat' WHERE candidate_id=$1`, candidateID)
	} else if result == "khong_dat" {
		database.Pool.Exec(ctx, `UPDATE hr_candidates SET pipeline_stage='khong_trung_tuyen' WHERE id=$1`, candidateID)
		database.Pool.Exec(ctx, `UPDATE hr_candidate_applications SET pipeline_stage='khong_trung_tuyen' WHERE candidate_id=$1`, candidateID)
	}
	return nil
}

// ---------------------- Timesheet & payroll calculation ----------------------

type TimesheetCalcService struct{}

// Calculate aggregates attendance logs + requests + holidays for one period
// (month "MM/YYYY") and writes hr_timesheet_details per employee per day.
func (s *TimesheetCalcService) Calculate(month string) (map[string]interface{}, error) {
	ctx := context.Background()
	// Delete previous details for the period to allow idempotent re-calculation
	details, err := database.Pool.Query(ctx,
		`SELECT id FROM hr_timesheet_details WHERE created_at::text ILIKE $1`, "%")
	if err != nil {
		return nil, err
	}
	details.Close()

	// wipe all detail rows (full recompute for demo period)
	database.Pool.Exec(ctx, `DELETE FROM hr_timesheet_details`)

	type attLog struct {
		EmployeeID string
		Name       string
		Date       string
		Value      float64
		Code       string
		OT         float64
	}
	rows, err := database.Pool.Query(ctx,
		`SELECT employee_id, employee_name, date, work_value, status_code, overtime_hours FROM hr_attendance_logs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []attLog
	for rows.Next() {
		var a attLog
		rows.Scan(&a.EmployeeID, &a.Name, &a.Date, &a.Value, &a.Code, &a.OT)
		logs = append(logs, a)
	}

	written := 0
	for _, a := range logs {
		if _, err := database.Pool.Exec(ctx,
			`INSERT INTO hr_timesheet_details (id, timesheet_id, employee_id, employee_name, date, work_value, status_code, overtime_hours)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (id) DO NOTHING`,
			"TSD-"+nowid(), month, a.EmployeeID, a.Name, a.Date, a.Value, a.Code, a.OT); err != nil {
			continue
		}
		written++
	}
	return map[string]interface{}{"month": month, "details_written": written}, nil
}

type PayrollCalcService struct{}

// Calculate snapshots active employees into payroll_records for a period and
// computes salary fields from their contract allowances and basic salary.
func (s *PayrollCalcService) Calculate(periodID, month string) (map[string]interface{}, error) {
	ctx := context.Background()
	var employees struct {
		ID     string
		Name   string
		Dept   string
		Status string
	}
	rows, err := database.Pool.Query(ctx,
		`SELECT id, name, department, status FROM hr_employees WHERE status='dang_lam_viec'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var empList []struct{ ID, Name, Dept string }
	for rows.Next() {
		rows.Scan(&employees.ID, &employees.Name, &employees.Dept, &employees.Status)
		empList = append(empList, struct{ ID, Name, Dept string }{employees.ID, employees.Name, employees.Dept})
	}
	written := 0
	for _, e := range empList {
		// pull latest contract
		var c struct {
			Base        float64
			Hazard      float64
			Safety      float64
			Meal        float64
			Resp        float64
			Att         float64
			InsBase     float64
			ContractTyp string
		}
		database.Pool.QueryRow(ctx,
			`SELECT base_salary, hazard_allowance, safety_allowance, meal_allowance, responsibility_allowance, attendance_allowance, social_insurance_base, contract_type
			 FROM hr_contracts WHERE employee_id=$1 ORDER BY created_at DESC LIMIT 1`, e.ID).
			Scan(&c.Base, &c.Hazard, &c.Safety, &c.Meal, &c.Resp, &c.Att, &c.InsBase, &c.ContractTyp)

		// aggregate worked days & OT from timesheet details
		var workedDays float64
		var otHours float64
		database.Pool.QueryRow(ctx,
			`SELECT COALESCE(SUM(work_value),0), COALESCE(SUM(overtime_hours),0)
			 FROM hr_timesheet_details WHERE employee_id=$1`, e.ID).Scan(&workedDays, &otHours)

		allowances := c.Hazard + c.Safety + c.Meal + c.Resp + c.Att
		working := c.Base * workedDays / 26.0
		otSalary := otHours * (c.Base / 26.0 / 8.0) * 1.5
		gross := working + otSalary + allowances
		socialIns := c.InsBase * 0.08
		healthIns := c.InsBase * 0.015
		unempIns := c.InsBase * 0.01
		insurance := socialIns + healthIns + unempIns
		tax := gross * 0.05
		net := gross - insurance - tax

		database.Pool.Exec(ctx,
			`INSERT INTO hr_payroll_records
			 (id, period_id, employee_id, employee_name, department, contract_type, bank_name, bank_account, base_salary,
			  actual_days, standard_days, working_salary, piecework_salary, overtime_salary, allowances, gross_salary,
			  social_insurance, health_insurance, unemployment_insurance, total_insurance, personal_income_tax,
			  advance_deduction, total_deductions, net_salary, status)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25)
			 ON CONFLICT (id) DO NOTHING`,
			"PR-"+nowid(), periodID, e.ID, e.Name, e.Dept, c.ContractTyp, "", "", c.Base,
			workedDays, 26.0, working, 0, otSalary, allowances, gross,
			socialIns, healthIns, unempIns, insurance, tax,
			0, insurance+tax, net, "da_tinh")
		written++
	}
	return map[string]interface{}{"month": month, "records_written": written}, nil
}

// ---------------------- helpers ----------------------

func jsonUnmarshal(b []byte, v interface{}) {
	_ = json.Unmarshal(b, v)
}

func userOrDefault() string {
	return "Hệ thống"
}

func nowid() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), time.Now().Nanosecond()%1000)
}

func currentPeriod() string {
	return time.Now().Format("01/2006")
}
