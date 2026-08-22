package handlers

import (
	"net/http"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/services"
)

type HrHandler struct {
	employeeSvc     *services.HrEmployeeService
	contractSvc     *services.HrContractService
	shiftSvc        *services.HrShiftService
	shiftSchedSvc   *services.HrShiftScheduleService
	attendanceSvc   *services.HrAttendanceService
	leaveReqSvc     *services.HrLeaveRequestService
	leaveAllocSvc   *services.HrLeaveAllocationService
	payslipSvc      *services.HrPayslipService
	timesheetSvc    *services.HrTimesheetService
}

func NewHrHandler() *HrHandler {
	return &HrHandler{
		employeeSvc:   services.NewHrEmployeeService(),
		contractSvc:   services.NewHrContractService(),
		shiftSvc:      services.NewHrShiftService(),
		shiftSchedSvc: services.NewHrShiftScheduleService(),
		attendanceSvc: services.NewHrAttendanceService(),
		leaveReqSvc:   services.NewHrLeaveRequestService(),
		leaveAllocSvc: services.NewHrLeaveAllocationService(),
		payslipSvc:    services.NewHrPayslipService(),
		timesheetSvc:  services.NewHrTimesheetService(),
	}
}

func (h *HrHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	var totalEmployees, totalContracts, totalShifts int
	database.Pool.QueryRow(r.Context(), "SELECT COUNT(*) FROM hr_employees").Scan(&totalEmployees)
	database.Pool.QueryRow(r.Context(), "SELECT COUNT(*) FROM hr_contracts").Scan(&totalContracts)
	database.Pool.QueryRow(r.Context(), "SELECT COUNT(*) FROM hr_shifts").Scan(&totalShifts)

	JSON(w, map[string]interface{}{
		"totalEmployees": totalEmployees,
		"totalContracts": totalContracts,
		"totalShifts":    totalShifts,
	})
}

func (h *HrHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.employeeSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *HrHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.employeeSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.employeeSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *HrHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.employeeSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.employeeSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *HrHandler) ListContracts(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.contractSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *HrHandler) GetContract(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.contractSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) CreateContract(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.contractSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *HrHandler) UpdateContract(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.contractSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) DeleteContract(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.contractSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *HrHandler) ListShifts(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.shiftSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *HrHandler) GetShift(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.shiftSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) CreateShift(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.shiftSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *HrHandler) UpdateShift(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.shiftSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) DeleteShift(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.shiftSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *HrHandler) ListShiftSchedules(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.shiftSchedSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *HrHandler) GetShiftSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.shiftSchedSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) CreateShiftSchedule(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.shiftSchedSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *HrHandler) UpdateShiftSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.shiftSchedSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) DeleteShiftSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.shiftSchedSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *HrHandler) ListAttendances(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.attendanceSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *HrHandler) GetAttendance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.attendanceSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) CreateAttendance(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.attendanceSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *HrHandler) UpdateAttendance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.attendanceSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) DeleteAttendance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.attendanceSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *HrHandler) ListLeaves(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.leaveReqSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *HrHandler) GetLeave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.leaveReqSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) CreateLeave(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.leaveReqSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *HrHandler) UpdateLeave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.leaveReqSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) DeleteLeave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.leaveReqSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *HrHandler) ListLeaveAllocations(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.leaveAllocSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *HrHandler) GetLeaveAllocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.leaveAllocSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) CreateLeaveAllocation(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.leaveAllocSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *HrHandler) UpdateLeaveAllocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.leaveAllocSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *HrHandler) DeleteLeaveAllocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.leaveAllocSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *HrHandler) ListTimesheets(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.timesheetSvc.ListJSONB(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *HrHandler) ListPayslips(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.payslipSvc.ListJSONB(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}
