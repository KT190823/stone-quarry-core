package handlers

import (
	"net/http"

	"mo-da-backend/internal/services"
)

type HrExtendedHandler struct {
	workHistSvc      *services.BaseService
	onboardingSvc    *services.BaseService
	decisionSvc      *services.BaseService
	esignDocSvc      *services.BaseService
	esignSignerSvc   *services.BaseService
	insRecSvc        *services.BaseService
	insDeclSvc       *services.BaseService
	recReqSvc        *services.BaseService
	campaignSvc      *services.BaseService
	candidateSvc     *services.BaseService
	candAppSvc       *services.BaseService
	candHistSvc      *services.BaseService
	interviewSvc     *services.BaseService
	courseSvc        *services.BaseService
	sessionSvc       *services.BaseService
	participantSvc   *services.BaseService
	trainingAttSvc   *services.BaseService
	evaluationSvc    *services.BaseService
	kpiSvc           *services.BaseService
	okrObjSvc        *services.BaseService
	okrKrSvc         *services.BaseService
	wfTemplateSvc    *services.BaseService
	wfInstanceSvc    *services.BaseService
	wfStepSvc        *services.BaseService
	shiftAssignSvc   *services.BaseService
	holidaySvc       *services.BaseService
	attLogSvc        *services.BaseService
	tsDetailSvc      *services.BaseService
	payrollPeriodSvc *services.BaseService
	payrollRecordSvc *services.BaseService
	formulaColSvc    *services.BaseService
	salaryAdvanceSvc *services.BaseService

	empStatusSvc     *services.EmployeeStatusService
	onboardingSrv    *services.OnboardingService
	decisionSrv      *services.DecisionService
	esignSrv         *services.EsignService
	insuranceSrv     *services.InsuranceService
	trainingSrv      *services.TrainingService
	kpiSrv           *services.KpiService
	okrSrv           *services.OkrService
	workflowSrv      *services.WorkflowService
	recruitSrv       *services.RecruitmentService
	candidateSrvSvc  *services.CandidateService
	interviewSrvSvc  *services.InterviewService
	timesheetCalcSvc *services.TimesheetCalcService
	payrollCalcSvc   *services.PayrollCalcService
}

func NewHrExtendedHandler() *HrExtendedHandler {
	return &HrExtendedHandler{
		workHistSvc:      services.NewBaseService(services.NewHrWorkHistoryService().Base),
		onboardingSvc:    services.NewBaseService(services.NewHrOnboardingService().Base),
		decisionSvc:      services.NewBaseService(services.NewHrDecisionService().Base),
		esignDocSvc:      services.NewBaseService(services.NewHrEsignDocumentService().Base),
		esignSignerSvc:   services.NewBaseService(services.NewHrEsignSignerService().Base),
		insRecSvc:        services.NewBaseService(services.NewHrInsuranceRecordService().Base),
		insDeclSvc:       services.NewBaseService(services.NewHrInsuranceDeclarationService().Base),
		recReqSvc:        services.NewBaseService(services.NewHrRecruitmentRequestService().Base),
		campaignSvc:      services.NewBaseService(services.NewHrCampaignService().Base),
		candidateSvc:     services.NewBaseService(services.NewHrCandidateService().Base),
		candAppSvc:       services.NewBaseService(services.NewHrCandidateApplicationService().Base),
		candHistSvc:      services.NewBaseService(services.NewHrCandidateStatusHistoryService().Base),
		interviewSvc:     services.NewBaseService(services.NewHrInterviewScheduleService().Base),
		courseSvc:        services.NewBaseService(services.NewHrTrainingCourseService().Base),
		sessionSvc:       services.NewBaseService(services.NewHrTrainingSessionService().Base),
		participantSvc:   services.NewBaseService(services.NewHrTrainingParticipantService().Base),
		trainingAttSvc:   services.NewBaseService(services.NewHrTrainingAttendanceService().Base),
		evaluationSvc:    services.NewBaseService(services.NewHrTrainingEvaluationService().Base),
		kpiSvc:           services.NewBaseService(services.NewHrKpiTargetService().Base),
		okrObjSvc:        services.NewBaseService(services.NewHrOkrObjectiveService().Base),
		okrKrSvc:         services.NewBaseService(services.NewHrOkrKeyResultService().Base),
		wfTemplateSvc:    services.NewBaseService(services.NewHrWorkflowTemplateService().Base),
		wfInstanceSvc:    services.NewBaseService(services.NewHrWorkflowInstanceService().Base),
		wfStepSvc:        services.NewBaseService(services.NewHrWorkflowStepService().Base),
		shiftAssignSvc:   services.NewBaseService(services.NewHrShiftAssignmentService().Base),
		holidaySvc:       services.NewBaseService(services.NewHrHolidayService().Base),
		attLogSvc:        services.NewBaseService(services.NewHrAttendanceLogService().Base),
		tsDetailSvc:      services.NewBaseService(services.NewHrTimesheetDetailService().Base),
		payrollPeriodSvc: services.NewBaseService(services.NewHrPayrollPeriodService().Base),
		payrollRecordSvc: services.NewBaseService(services.NewHrPayrollRecordService().Base),
		formulaColSvc:    services.NewBaseService(services.NewHrPayrollFormulaColumnService().Base),
		salaryAdvanceSvc: services.NewBaseService(services.NewHrSalaryAdvanceService().Base),

		empStatusSvc:     &services.EmployeeStatusService{},
		onboardingSrv:    &services.OnboardingService{},
		decisionSrv:      &services.DecisionService{},
		esignSrv:         &services.EsignService{},
		insuranceSrv:     &services.InsuranceService{},
		trainingSrv:      &services.TrainingService{},
		kpiSrv:           &services.KpiService{},
		okrSrv:           &services.OkrService{},
		workflowSrv:      &services.WorkflowService{},
		recruitSrv:       &services.RecruitmentService{},
		candidateSrvSvc:  &services.CandidateService{},
		interviewSrvSvc:  &services.InterviewService{},
		timesheetCalcSvc: &services.TimesheetCalcService{},
		payrollCalcSvc:   &services.PayrollCalcService{},
	}
}

// generic CRUD helpers
func (h *HrExtendedHandler) list(svc *services.BaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := parseListParams(r)
		results, total, err := svc.List(params)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		JSON(w, map[string]interface{}{"data": results, "total": total})
	}
}

func (h *HrExtendedHandler) get(svc *services.BaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		result, err := svc.GetByID(id)
		if err != nil {
			http.Error(w, "not found", 404)
			return
		}
		JSON(w, result)
	}
}

func (h *HrExtendedHandler) create(svc *services.BaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readJSON(r)
		if err != nil {
			http.Error(w, "invalid JSON", 400)
			return
		}
		result, err := svc.Create(data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(201)
		JSON(w, result)
	}
}

func (h *HrExtendedHandler) update(svc *services.BaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		data, err := readJSON(r)
		if err != nil {
			http.Error(w, "invalid JSON", 400)
			return
		}
		result, err := svc.Update(id, data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		JSON(w, result)
	}
}

func (h *HrExtendedHandler) delete(svc *services.BaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := svc.Delete(id); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		JSON(w, map[string]string{"status": "deleted"})
	}
}

func (h *HrExtendedHandler) crud(svc *services.BaseService) (http.HandlerFunc, http.HandlerFunc, http.HandlerFunc, http.HandlerFunc) {
	return h.list(svc), h.get(svc), h.create(svc), h.update(svc)
}
func (h *HrExtendedHandler) del(svc *services.BaseService) http.HandlerFunc { return h.delete(svc) }

// Register mounts all HRM-extended routes onto the given mux.
func (h *HrExtendedHandler) Register(mux *http.ServeMux) {
	reg := func(pattern string, fn http.HandlerFunc) { mux.HandleFunc(pattern, fn) }
	L, G := h.list, h.get
	Cl, U := h.create, h.update
	D := h.del

	reg("GET /api/hr-ext/dashboard", h.Dashboard)

	// Recruitment
	reg("GET /api/hr-ext/recruitment-requests", L(h.recReqSvc))
	reg("GET /api/hr-ext/recruitment-requests/{id}", G(h.recReqSvc))
	reg("POST /api/hr-ext/recruitment-requests", Cl(h.recReqSvc))
	reg("PUT /api/hr-ext/recruitment-requests/{id}", U(h.recReqSvc))
	reg("DELETE /api/hr-ext/recruitment-requests/{id}", D(h.recReqSvc))
	reg("PATCH /api/hr-ext/recruitment-requests/{id}/approve", h.RecruitmentApprove)
	reg("PATCH /api/hr-ext/recruitment-requests/{id}/reject", h.RecruitmentReject)
	reg("PATCH /api/hr-ext/recruitment-requests/{id}/submit", h.RecruitmentSubmit)

	reg("GET /api/hr-ext/campaigns", L(h.campaignSvc))
	reg("GET /api/hr-ext/campaigns/{id}", G(h.campaignSvc))
	reg("POST /api/hr-ext/campaigns", Cl(h.campaignSvc))
	reg("PUT /api/hr-ext/campaigns/{id}", U(h.campaignSvc))
	reg("DELETE /api/hr-ext/campaigns/{id}", D(h.campaignSvc))

	reg("GET /api/hr-ext/candidates", L(h.candidateSvc))
	reg("GET /api/hr-ext/candidates/{id}", G(h.candidateSvc))
	reg("POST /api/hr-ext/candidates", Cl(h.candidateSvc))
	reg("PUT /api/hr-ext/candidates/{id}", U(h.candidateSvc))
	reg("DELETE /api/hr-ext/candidates/{id}", D(h.candidateSvc))
	reg("PATCH /api/hr-ext/candidates/{id}/stage", h.CandidateStageChange)
	reg("POST /api/hr-ext/candidates/{id}/convert-to-employee", h.CandidateConvert)

	reg("GET /api/hr-ext/candidate-applications", L(h.candAppSvc))
	reg("POST /api/hr-ext/candidate-applications", Cl(h.candAppSvc))

	reg("GET /api/hr-ext/interview-schedules", L(h.interviewSvc))
	reg("GET /api/hr-ext/interview-schedules/{id}", G(h.interviewSvc))
	reg("POST /api/hr-ext/interview-schedules", Cl(h.interviewSvc))
	reg("PUT /api/hr-ext/interview-schedules/{id}", U(h.interviewSvc))
	reg("DELETE /api/hr-ext/interview-schedules/{id}", D(h.interviewSvc))
	reg("PATCH /api/hr-ext/interview-schedules/{id}/result", h.InterviewResult)

	// Employee lifecycle
	reg("PATCH /api/hr-ext/employees/{id}/status", h.EmployeeStatusChange)
	reg("GET /api/hr-ext/work-history", L(h.workHistSvc))
	reg("GET /api/hr-ext/onboarding-checklist", L(h.onboardingSvc))
	reg("POST /api/hr-ext/employees/{id}/onboarding-checklist", h.OnboardingGenerate)
	reg("PATCH /api/hr-ext/onboarding-checklist/{id}", h.OnboardingComplete)

	// Decisions & E-sign
	reg("GET /api/hr-ext/decisions", L(h.decisionSvc))
	reg("GET /api/hr-ext/decisions/{id}", G(h.decisionSvc))
	reg("POST /api/hr-ext/decisions", Cl(h.decisionSvc))
	reg("PUT /api/hr-ext/decisions/{id}", U(h.decisionSvc))
	reg("DELETE /api/hr-ext/decisions/{id}", D(h.decisionSvc))
	reg("POST /api/hr-ext/decisions/{id}/issue", h.DecisionIssue)

	reg("GET /api/hr-ext/esign-documents", L(h.esignDocSvc))
	reg("GET /api/hr-ext/esign-documents/{id}", G(h.esignDocSvc))
	reg("POST /api/hr-ext/esign-documents", h.EsignCreate)
	reg("GET /api/hr-ext/esign-signers", L(h.esignSignerSvc))
	reg("PATCH /api/hr-ext/esign-signers/{id}", h.EsignSign)

	// Insurance
	reg("GET /api/hr-ext/insurance-records", L(h.insRecSvc))
	reg("GET /api/hr-ext/insurance-records/{id}", G(h.insRecSvc))
	reg("POST /api/hr-ext/insurance-records", Cl(h.insRecSvc))
	reg("PUT /api/hr-ext/insurance-records/{id}", U(h.insRecSvc))
	reg("GET /api/hr-ext/insurance-declarations", L(h.insDeclSvc))
	reg("POST /api/hr-ext/insurance-declarations", Cl(h.insDeclSvc))
	reg("POST /api/hr-ext/insurance-declarations/generate", h.InsuranceDeclaration)

	// Training
	reg("GET /api/hr-ext/training-courses", L(h.courseSvc))
	reg("GET /api/hr-ext/training-courses/{id}", G(h.courseSvc))
	reg("POST /api/hr-ext/training-courses", Cl(h.courseSvc))
	reg("PUT /api/hr-ext/training-courses/{id}", U(h.courseSvc))
	reg("POST /api/hr-ext/training-courses/{id}/recalculate", h.TrainingRecalculate)
	reg("GET /api/hr-ext/training-sessions", L(h.sessionSvc))
	reg("POST /api/hr-ext/training-sessions", Cl(h.sessionSvc))
	reg("GET /api/hr-ext/training-participants", L(h.participantSvc))
	reg("POST /api/hr-ext/training-participants", Cl(h.participantSvc))
	reg("GET /api/hr-ext/training-attendances", L(h.trainingAttSvc))
	reg("POST /api/hr-ext/training-attendances", Cl(h.trainingAttSvc))
	reg("GET /api/hr-ext/training-evaluations", L(h.evaluationSvc))
	reg("POST /api/hr-ext/training-evaluations", Cl(h.evaluationSvc))

	// KPI / OKR
	reg("GET /api/hr-ext/kpi-targets", L(h.kpiSvc))
	reg("POST /api/hr-ext/kpi-targets", Cl(h.kpiSvc))
	reg("PUT /api/hr-ext/kpi-targets/{id}", U(h.kpiSvc))
	reg("POST /api/hr-ext/kpi-targets/{id}/recalculate", h.KpiRecalculate)
	reg("GET /api/hr-ext/okr-objectives", L(h.okrObjSvc))
	reg("POST /api/hr-ext/okr-objectives", Cl(h.okrObjSvc))
	reg("POST /api/hr-ext/okr-objectives/{id}/recalculate", h.OkrRecalculateObjective)
	reg("GET /api/hr-ext/okr-key-results", L(h.okrKrSvc))
	reg("POST /api/hr-ext/okr-key-results", Cl(h.okrKrSvc))
	reg("PATCH /api/hr-ext/okr-key-results/{id}/recalculate", h.OkrRecalculateKr)

	// Workflow engine
	reg("GET /api/hr-ext/workflow-templates", L(h.wfTemplateSvc))
	reg("POST /api/hr-ext/workflow-templates", Cl(h.wfTemplateSvc))
	reg("GET /api/hr-ext/workflow-instances", L(h.wfInstanceSvc))
	reg("GET /api/hr-ext/workflow-instances/{id}", G(h.wfInstanceSvc))
	reg("POST /api/hr-ext/workflow-instances", h.WorkflowStart)
	reg("PATCH /api/hr-ext/workflow-instances/{id}/resolve", h.WorkflowResolve)
	reg("GET /api/hr-ext/workflow-steps", L(h.wfStepSvc))

	// Timesheet & payroll relational
	reg("GET /api/hr-ext/shift-assignments", L(h.shiftAssignSvc))
	reg("POST /api/hr-ext/shift-assignments", Cl(h.shiftAssignSvc))
	reg("GET /api/hr-ext/holidays", L(h.holidaySvc))
	reg("POST /api/hr-ext/holidays", Cl(h.holidaySvc))
	reg("GET /api/hr-ext/attendance-logs", L(h.attLogSvc))
	reg("POST /api/hr-ext/attendance-logs", Cl(h.attLogSvc))
	reg("GET /api/hr-ext/timesheet-details", L(h.tsDetailSvc))
	reg("POST /api/hr-ext/timesheet-details", Cl(h.tsDetailSvc))
	reg("POST /api/hr-ext/timesheets/calculate", h.TimesheetCalculate)
	reg("GET /api/hr-ext/payroll-periods", L(h.payrollPeriodSvc))
	reg("POST /api/hr-ext/payroll-periods", Cl(h.payrollPeriodSvc))
	reg("GET /api/hr-ext/payroll-records", L(h.payrollRecordSvc))
	reg("POST /api/hr-ext/payroll-records", Cl(h.payrollRecordSvc))
	reg("POST /api/hr-ext/payroll-periods/calculate", h.PayrollCalculate)
	reg("GET /api/hr-ext/payroll-formula-columns", L(h.formulaColSvc))
	reg("POST /api/hr-ext/payroll-formula-columns", Cl(h.formulaColSvc))
	reg("GET /api/hr-ext/salary-advances", L(h.salaryAdvanceSvc))
	reg("POST /api/hr-ext/salary-advances", Cl(h.salaryAdvanceSvc))
	reg("PUT /api/hr-ext/salary-advances/{id}", U(h.salaryAdvanceSvc))
}

// ---------------------- Employee lifecycle ----------------------

func (h *HrExtendedHandler) EmployeeStatusChange(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	toStatus, _ := data["to_status"].(string)
	effective, _ := data["effective_date"].(string)
	reason, _ := data["reason"].(string)
	changedBy, _ := data["changed_by"].(string)
	result, err := h.empStatusSvc.ChangeStatus(id, toStatus, effective, reason, changedBy)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrExtendedHandler) OnboardingGenerate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.onboardingSrv.GenerateChecklist(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "generated"})
}

func (h *HrExtendedHandler) OnboardingComplete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.onboardingSrv.CompleteItem(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "completed"})
}

func (h *HrExtendedHandler) DecisionIssue(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.decisionSrv.Issue(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrExtendedHandler) EsignCreate(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	contractID, _ := data["contract_id"].(string)
	contractCode, _ := data["contract_code"].(string)
	subject, _ := data["subject"].(string)
	employeeID, _ := data["employee_id"].(string)
	employeeName, _ := data["employee_name"].(string)
	hrName, _ := data["hr_name"].(string)
	result, err := h.esignSrv.CreateForContract(contractID, contractCode, subject, employeeID, employeeName, hrName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrExtendedHandler) EsignSign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	status, _ := data["status"].(string)
	note, _ := data["note"].(string)
	if err := h.esignSrv.Sign(id, status, note); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "signed"})
}

func (h *HrExtendedHandler) InsuranceDeclaration(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	employeeID, _ := data["employee_id"].(string)
	declType, _ := data["declaration_type"].(string)
	period, _ := data["period"].(string)
	if err := h.insuranceSrv.CreateDeclarationForEmployee(employeeID, declType, period); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "created"})
}

func (h *HrExtendedHandler) TrainingRecalculate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.trainingSrv.RecalculateCompletion(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "recalculated"})
}

func (h *HrExtendedHandler) KpiRecalculate(w http.ResponseWriter, r *http.Request) {
	employeeID := r.PathValue("id")
	period := r.URL.Query().Get("period")
	result, err := h.kpiSrv.RecalculateEmployee(employeeID, period)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrExtendedHandler) OkrRecalculateKr(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	progress, err := h.okrSrv.RecalculateKeyResult(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"key_result_id": id, "progress_percent": progress})
}

func (h *HrExtendedHandler) OkrRecalculateObjective(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	progress, err := h.okrSrv.RecalculateObjective(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"objective_id": id, "progress_percent": progress})
}

func (h *HrExtendedHandler) WorkflowStart(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	appliesTo, _ := data["applies_to"].(string)
	templateID, _ := data["template_id"].(string)
	objectID, _ := data["object_id"].(string)
	objectRef, _ := data["object_ref"].(string)
	result, err := h.workflowSrv.Start(appliesTo, templateID, objectID, objectRef, "")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrExtendedHandler) WorkflowResolve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	decision, _ := data["decision"].(string)
	comment, _ := data["comment"].(string)
	approverID, _ := data["approver_id"].(string)
	approverName, _ := data["approver_name"].(string)
	result, err := h.workflowSrv.ResolveStep(id, decision, comment, approverID, approverName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrExtendedHandler) RecruitmentSubmit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.recruitSrv.Submit(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "submitted"})
}

func (h *HrExtendedHandler) RecruitmentApprove(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, _ := readJSON(r)
	comment, _ := data["comment"].(string)
	if err := h.recruitSrv.Approve(id, comment); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "approved"})
}

func (h *HrExtendedHandler) RecruitmentReject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, _ := readJSON(r)
	comment, _ := data["comment"].(string)
	if err := h.recruitSrv.Reject(id, comment); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	JSON(w, map[string]string{"status": "rejected"})
}

func (h *HrExtendedHandler) CandidateStageChange(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	toStage, _ := data["to_stage"].(string)
	note, _ := data["note"].(string)
	changedBy, _ := data["changed_by"].(string)
	if err := h.candidateSrvSvc.ChangeStage(id, toStage, note, changedBy); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "stage_changed"})
}

func (h *HrExtendedHandler) CandidateConvert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	empID, err := h.candidateSrvSvc.ConvertToEmployee(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"status": "converted", "employee_id": empID})
}

func (h *HrExtendedHandler) InterviewResult(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, _ := data["result"].(string)
	note, _ := data["note"].(string)
	if err := h.interviewSrvSvc.SetResult(id, result, note); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "result_saved"})
}

func (h *HrExtendedHandler) TimesheetCalculate(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	if month == "" {
		month = "10/2026"
	}
	result, err := h.timesheetCalcSvc.Calculate(month)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrExtendedHandler) PayrollCalculate(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	periodID, _ := data["period_id"].(string)
	month, _ := data["month"].(string)
	if periodID == "" {
		periodID = "PR-1026"
	}
	if month == "" {
		month = "10/2026"
	}
	result, err := h.payrollCalcSvc.Calculate(periodID, month)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrExtendedHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	cnt := func(table, where string, args ...interface{}) int {
		n, _ := services.CountRows(table, where, args...)
		return n
	}
	JSON(w, map[string]interface{}{
		"employees":           cnt("hr_employees", ""),
		"activeEmployees":     cnt("hr_employees", "status='dang_lam_viec'"),
		"contracts":           cnt("hr_contracts", ""),
		"leaveRequests":       cnt("hr_leave_requests", ""),
		"pendingLeaves":       cnt("hr_leave_requests", "stage LIKE '%Chờ%'"),
		"recruitmentRequests": cnt("hr_recruitment_requests", ""),
		"approvedRecRequests": cnt("hr_recruitment_requests", "approval_status='da_duyet'"),
		"campaignsActive":     cnt("hr_recruitment_campaigns", "status='dang_tuyen'"),
		"candidates":          cnt("hr_candidates", ""),
		"openCandidates":      cnt("hr_candidates", "pipeline_stage NOT IN ('da_tao_nhan_su','khong_trung_tuyen','blacklist')"),
		"decisions":           cnt("hr_decisions", ""),
		"pendingDecisions":    cnt("hr_decisions", "status='cho_duyet'"),
		"insuranceDecls":      cnt("hr_insurance_declarations", ""),
		"trainingCourses":     cnt("hr_training_courses", ""),
		"kpiTargets":          cnt("hr_kpi_targets", ""),
		"salaryAdvances":      cnt("hr_salary_advances", ""),
		"payrollPeriods":      cnt("hr_payroll_periods", ""),
	})
}
