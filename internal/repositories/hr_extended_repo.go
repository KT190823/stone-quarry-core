package repositories

import (
	"context"
	"fmt"
	"strings"

	"mo-da-backend/internal/database"
)

// Generic repository helper for HR extended tables.
func newRepo(table string) *BaseRepo {
	return NewBaseRepo(table, "id")
}

type HrWorkHistoryRepo struct {
	*BaseRepo
}

func NewHrWorkHistoryRepo() *HrWorkHistoryRepo {
	return &HrWorkHistoryRepo{BaseRepo: newRepo("hr_employee_work_history")}
}

type HrOnboardingRepo struct {
	*BaseRepo
}

func NewHrOnboardingRepo() *HrOnboardingRepo {
	return &HrOnboardingRepo{BaseRepo: newRepo("hr_onboarding_checklists")}
}

type HrDecisionRepo struct {
	*BaseRepo
}

func NewHrDecisionRepo() *HrDecisionRepo {
	return &HrDecisionRepo{BaseRepo: newRepo("hr_decisions")}
}

type HrEsignDocumentRepo struct {
	*BaseRepo
}

func NewHrEsignDocumentRepo() *HrEsignDocumentRepo {
	return &HrEsignDocumentRepo{BaseRepo: newRepo("hr_esign_documents")}
}

type HrEsignSignerRepo struct {
	*BaseRepo
}

func NewHrEsignSignerRepo() *HrEsignSignerRepo {
	return &HrEsignSignerRepo{BaseRepo: newRepo("hr_esign_signers")}
}

type HrInsuranceRecordRepo struct {
	*BaseRepo
}

func NewHrInsuranceRecordRepo() *HrInsuranceRecordRepo {
	return &HrInsuranceRecordRepo{BaseRepo: newRepo("hr_insurance_records")}
}

type HrInsuranceDeclarationRepo struct {
	*BaseRepo
}

func NewHrInsuranceDeclarationRepo() *HrInsuranceDeclarationRepo {
	return &HrInsuranceDeclarationRepo{BaseRepo: newRepo("hr_insurance_declarations")}
}

type HrRecruitmentRequestRepo struct {
	*BaseRepo
}

func NewHrRecruitmentRequestRepo() *HrRecruitmentRequestRepo {
	return &HrRecruitmentRequestRepo{BaseRepo: newRepo("hr_recruitment_requests")}
}

type HrCampaignRepo struct {
	*BaseRepo
}

func NewHrCampaignRepo() *HrCampaignRepo {
	return &HrCampaignRepo{BaseRepo: newRepo("hr_recruitment_campaigns")}
}

type HrCandidateRepo struct {
	*BaseRepo
}

func NewHrCandidateRepo() *HrCandidateRepo {
	return &HrCandidateRepo{BaseRepo: newRepo("hr_candidates")}
}

type HrCandidateApplicationRepo struct {
	*BaseRepo
}

func NewHrCandidateApplicationRepo() *HrCandidateApplicationRepo {
	return &HrCandidateApplicationRepo{BaseRepo: newRepo("hr_candidate_applications")}
}

type HrCandidateStatusHistoryRepo struct {
	*BaseRepo
}

func NewHrCandidateStatusHistoryRepo() *HrCandidateStatusHistoryRepo {
	return &HrCandidateStatusHistoryRepo{BaseRepo: newRepo("hr_candidate_status_history")}
}

type HrInterviewScheduleRepo struct {
	*BaseRepo
}

func NewHrInterviewScheduleRepo() *HrInterviewScheduleRepo {
	return &HrInterviewScheduleRepo{BaseRepo: newRepo("hr_interview_schedules")}
}

type HrTrainingCourseRepo struct {
	*BaseRepo
}

func NewHrTrainingCourseRepo() *HrTrainingCourseRepo {
	return &HrTrainingCourseRepo{BaseRepo: newRepo("hr_training_courses")}
}

type HrTrainingSessionRepo struct {
	*BaseRepo
}

func NewHrTrainingSessionRepo() *HrTrainingSessionRepo {
	return &HrTrainingSessionRepo{BaseRepo: newRepo("hr_training_sessions")}
}

type HrTrainingParticipantRepo struct {
	*BaseRepo
}

func NewHrTrainingParticipantRepo() *HrTrainingParticipantRepo {
	return &HrTrainingParticipantRepo{BaseRepo: newRepo("hr_training_participants")}
}

type HrTrainingAttendanceRepo struct {
	*BaseRepo
}

func NewHrTrainingAttendanceRepo() *HrTrainingAttendanceRepo {
	return &HrTrainingAttendanceRepo{BaseRepo: newRepo("hr_training_attendances")}
}

type HrTrainingEvaluationRepo struct {
	*BaseRepo
}

func NewHrTrainingEvaluationRepo() *HrTrainingEvaluationRepo {
	return &HrTrainingEvaluationRepo{BaseRepo: newRepo("hr_training_evaluations")}
}

type HrKpiTargetRepo struct {
	*BaseRepo
}

func NewHrKpiTargetRepo() *HrKpiTargetRepo {
	return &HrKpiTargetRepo{BaseRepo: newRepo("hr_kpi_targets")}
}

type HrOkrObjectiveRepo struct {
	*BaseRepo
}

func NewHrOkrObjectiveRepo() *HrOkrObjectiveRepo {
	return &HrOkrObjectiveRepo{BaseRepo: newRepo("hr_okr_objectives")}
}

type HrOkrKeyResultRepo struct {
	*BaseRepo
}

func NewHrOkrKeyResultRepo() *HrOkrKeyResultRepo {
	return &HrOkrKeyResultRepo{BaseRepo: newRepo("hr_okr_key_results")}
}

type HrWorkflowTemplateRepo struct {
	*BaseRepo
}

func NewHrWorkflowTemplateRepo() *HrWorkflowTemplateRepo {
	return &HrWorkflowTemplateRepo{BaseRepo: newRepo("hr_workflow_templates")}
}

type HrWorkflowInstanceRepo struct {
	*BaseRepo
}

func NewHrWorkflowInstanceRepo() *HrWorkflowInstanceRepo {
	return &HrWorkflowInstanceRepo{BaseRepo: newRepo("hr_workflow_instances")}
}

type HrWorkflowStepRepo struct {
	*BaseRepo
}

func NewHrWorkflowStepRepo() *HrWorkflowStepRepo {
	return &HrWorkflowStepRepo{BaseRepo: newRepo("hr_workflow_steps")}
}

type HrShiftAssignmentRepo struct {
	*BaseRepo
}

func NewHrShiftAssignmentRepo() *HrShiftAssignmentRepo {
	return &HrShiftAssignmentRepo{BaseRepo: newRepo("hr_shift_assignments")}
}

type HrHolidayRepo struct {
	*BaseRepo
}

func NewHrHolidayRepo() *HrHolidayRepo {
	return &HrHolidayRepo{BaseRepo: newRepo("hr_holidays")}
}

type HrAttendanceLogRepo struct {
	*BaseRepo
}

func NewHrAttendanceLogRepo() *HrAttendanceLogRepo {
	return &HrAttendanceLogRepo{BaseRepo: newRepo("hr_attendance_logs")}
}

type HrTimesheetDetailRepo struct {
	*BaseRepo
}

func NewHrTimesheetDetailRepo() *HrTimesheetDetailRepo {
	return &HrTimesheetDetailRepo{BaseRepo: newRepo("hr_timesheet_details")}
}

type HrPayrollPeriodRepo struct {
	*BaseRepo
}

func NewHrPayrollPeriodRepo() *HrPayrollPeriodRepo {
	return &HrPayrollPeriodRepo{BaseRepo: newRepo("hr_payroll_periods")}
}

type HrPayrollRecordRepo struct {
	*BaseRepo
}

func NewHrPayrollRecordRepo() *HrPayrollRecordRepo {
	return &HrPayrollRecordRepo{BaseRepo: newRepo("hr_payroll_records")}
}

type HrPayrollFormulaColumnRepo struct {
	*BaseRepo
}

func NewHrPayrollFormulaColumnRepo() *HrPayrollFormulaColumnRepo {
	return &HrPayrollFormulaColumnRepo{BaseRepo: newRepo("hr_payroll_formula_columns")}
}

type HrSalaryAdvanceRepo struct {
	*BaseRepo
}

func NewHrSalaryAdvanceRepo() *HrSalaryAdvanceRepo {
	return &HrSalaryAdvanceRepo{BaseRepo: newRepo("hr_salary_advances")}
}

// Rebuild HR dashboard aggregates.
func (r *BaseRepo) CountRows(table, where string, args ...interface{}) (int, error) {
	ctx := context.Background()
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, where)
	if strings.TrimSpace(where) == "" {
		query = fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	}
	var n int
	err := database.Pool.QueryRow(ctx, query, args...).Scan(&n)
	return n, err
}
