package repositories

type HrEmployeeRepo struct {
	*BaseRepo
}

func NewHrEmployeeRepo() *HrEmployeeRepo {
	return &HrEmployeeRepo{BaseRepo: NewBaseRepo("hr_employees", "id")}
}

type HrContractRepo struct {
	*BaseRepo
}

func NewHrContractRepo() *HrContractRepo {
	return &HrContractRepo{BaseRepo: NewBaseRepo("hr_contracts", "id")}
}

type HrShiftRepo struct {
	*BaseRepo
}

func NewHrShiftRepo() *HrShiftRepo {
	return &HrShiftRepo{BaseRepo: NewBaseRepo("hr_shifts", "id")}
}

type HrShiftScheduleRepo struct {
	*BaseRepo
}

func NewHrShiftScheduleRepo() *HrShiftScheduleRepo {
	return &HrShiftScheduleRepo{BaseRepo: NewBaseRepo("hr_shift_schedules", "id")}
}

type HrAttendanceRepo struct {
	*BaseRepo
}

func NewHrAttendanceRepo() *HrAttendanceRepo {
	return &HrAttendanceRepo{BaseRepo: NewBaseRepo("hr_attendances", "id")}
}

type HrLeaveRequestRepo struct {
	*BaseRepo
}

func NewHrLeaveRequestRepo() *HrLeaveRequestRepo {
	return &HrLeaveRequestRepo{BaseRepo: NewBaseRepo("hr_leave_requests", "id")}
}

type HrLeaveAllocationRepo struct {
	*BaseRepo
}

func NewHrLeaveAllocationRepo() *HrLeaveAllocationRepo {
	return &HrLeaveAllocationRepo{BaseRepo: NewBaseRepo("hr_leave_allocations", "id")}
}

type HrPayslipRepo struct {
	*BaseRepo
}

func NewHrPayslipRepo() *HrPayslipRepo {
	return &HrPayslipRepo{BaseRepo: NewBaseRepo("hr_payslips", "id")}
}

type HrTimesheetRepo struct {
	*BaseRepo
}

func NewHrTimesheetRepo() *HrTimesheetRepo {
	return &HrTimesheetRepo{BaseRepo: NewBaseRepo("hr_timesheets", "id")}
}
