package models

type HrEmployee struct {
	ID           string   `json:"id"`
	Code         string   `json:"code"`
	Name         string   `json:"name"`
	Avatar       string   `json:"avatar"`
	Gender       string   `json:"gender"`
	DOB          string   `json:"dob"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	IDCard       string   `json:"idCard"`
	Address      string   `json:"address"`
	Department   string   `json:"department"`
	JobPosition  string   `json:"jobPosition"`
	Manager      string   `json:"manager"`
	JoinDate     string   `json:"joinDate"`
	ContractType string   `json:"contractType"`
	WorkLocation string   `json:"workLocation"`
	Status       string   `json:"status"`
	Certificates []string `json:"certificates"`
	BaseSalary   float64  `json:"baseSalary"`
}

type HrContract struct {
	ID                      string  `json:"id"`
	Code                    string  `json:"code"`
	EmployeeID              string  `json:"employeeId"`
	EmployeeName            string  `json:"employeeName"`
	Department              string  `json:"department"`
	ContractType            string  `json:"contractType"`
	StartDate               string  `json:"startDate"`
	EndDate                 string  `json:"endDate"`
	BaseSalary              float64 `json:"baseSalary"`
	HazardAllowance         float64 `json:"hazardAllowance"`
	SafetyAllowance         float64 `json:"safetyAllowance"`
	MealAllowance           float64 `json:"mealAllowance"`
	ResponsibilityAllowance float64 `json:"responsibilityAllowance"`
	AttendanceAllowance     float64 `json:"attendanceAllowance"`
	SocialInsuranceBase     float64 `json:"socialInsuranceBase"`
	Status                  string  `json:"status"`
}

type HrShift struct {
	ID              string  `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	StartTime       string  `json:"startTime"`
	EndTime         string  `json:"endTime"`
	BreakHours      float64 `json:"breakHours"`
	WorkHours       float64 `json:"workHours"`
	NightMultiplier float64 `json:"nightMultiplier"`
	Color           string  `json:"color"`
	Description     string  `json:"description"`
}

type HrShiftSchedule struct {
	ID           string `json:"id"`
	EmployeeID   string `json:"employeeId"`
	EmployeeName string `json:"employeeName"`
	Department   string `json:"department"`
	ShiftCode    string `json:"shiftCode"`
	ShiftName    string `json:"shiftName"`
	Date         string `json:"date"`
	TimeSlot     string `json:"timeSlot"`
	Note         string `json:"note,omitempty"`
	Status       string `json:"status"`
}

type HrAttendance struct {
	ID            string  `json:"id"`
	EmployeeID    string  `json:"employeeId"`
	EmployeeName  string  `json:"employeeName"`
	Department    string  `json:"department"`
	JobPosition   string  `json:"jobPosition"`
	EquipmentType string  `json:"equipmentType"`
	EquipmentCode string  `json:"equipmentCode"`
	VehiclePlate  string  `json:"vehiclePlate"`
	Date          string  `json:"date"`
	ShiftName     string  `json:"shiftName"`
	CheckInTime   string  `json:"checkInTime"`
	CheckOutTime  string  `json:"checkOutTime"`
	WorkedHours   float64 `json:"workedHours"`
	LateMinutes   int     `json:"lateMinutes"`
	EarlyMinutes  int     `json:"earlyMinutes"`
	OTHours       float64 `json:"otHours"`
	Method        string  `json:"method"`
	Location      string  `json:"location"`
	WorkArea      string  `json:"workArea"`
	Status        string  `json:"status"`
	Notes         string  `json:"notes"`
}

type HrLeaveRequest struct {
	ID           string  `json:"id"`
	Code         string  `json:"code"`
	EmployeeID   string  `json:"employeeId"`
	EmployeeName string  `json:"employeeName"`
	Department   string  `json:"department"`
	LeaveType    string  `json:"leaveType"`
	StartDate    string  `json:"startDate"`
	EndDate      string  `json:"endDate"`
	TotalDays    float64 `json:"totalDays"`
	TotalHours   float64 `json:"totalHours"`
	Reason       string  `json:"reason"`
	Stage        string  `json:"stage"`
	Approver     string  `json:"approver"`
	CreatedAt    string  `json:"createdAt"`
}

type HrLeaveAllocation struct {
	EmployeeID     string `json:"employeeId"`
	EmployeeName   string `json:"employeeName"`
	Department     string `json:"department"`
	Year           int    `json:"year"`
	TotalAllocated int    `json:"totalAllocated"`
	UsedDays       int    `json:"usedDays"`
	RemainingDays  int    `json:"remainingDays"`
}

type HrDailyTimesheetEntry struct {
	Day   int     `json:"day"`
	Code  string  `json:"code"`
	Hours float64 `json:"hours"`
	OT    float64 `json:"ot"`
}

type HrMonthlyTimesheet struct {
	ID               string                  `json:"id"`
	EmployeeID       string                  `json:"employeeId"`
	EmployeeName     string                  `json:"employeeName"`
	Department       string                  `json:"department"`
	JobPosition      string                  `json:"jobPosition"`
	Month            string                  `json:"month"`
	StandardDays     int                     `json:"standardDays"`
	ActualDays       int                     `json:"actualDays"`
	PaidLeaveDays    int                     `json:"paidLeaveDays"`
	UnpaidDays       int                     `json:"unpaidDays"`
	SickDays         int                     `json:"sickDays"`
	HolidayDays      int                     `json:"holidayDays"`
	TotalPaidDays    int                     `json:"totalPaidDays"`
	LateCount        int                     `json:"lateCount"`
	LateMinutesTotal int                     `json:"lateMinutesTotal"`
	OTNormalHours    float64                 `json:"otNormalHours"`
	OTNightHours     float64                 `json:"otNightHours"`
	PieceworkTons    float64                 `json:"pieceworkTons,omitempty"`
	PieceworkType    string                  `json:"pieceworkType,omitempty"`
	Status           string                  `json:"status"`
	DailyEntries     []HrDailyTimesheetEntry `json:"dailyEntries"`
}

type HrPayslip struct {
	ID                      string  `json:"id"`
	Code                    string  `json:"code"`
	EmployeeID              string  `json:"employeeId"`
	EmployeeName            string  `json:"employeeName"`
	Department              string  `json:"department"`
	JobPosition             string  `json:"jobPosition"`
	Month                   string  `json:"month"`
	BaseSalary              float64 `json:"baseSalary"`
	StandardDays            int     `json:"standardDays"`
	ActualDays              int     `json:"actualDays"`
	WorkingSalary           float64 `json:"workingSalary"`
	PieceworkTons           float64 `json:"pieceworkTons"`
	PieceworkRate           float64 `json:"pieceworkRate"`
	PieceworkSalary         float64 `json:"pieceworkSalary"`
	OTNormalHours           float64 `json:"otNormalHours"`
	OTNormalSalary          float64 `json:"otNormalSalary"`
	OTNightHours            float64 `json:"otNightHours"`
	OTNightSalary           float64 `json:"otNightSalary"`
	HazardAllowance         float64 `json:"hazardAllowance"`
	SafetyAllowance         float64 `json:"safetyAllowance"`
	MealAllowance           float64 `json:"mealAllowance"`
	ResponsibilityAllowance float64 `json:"responsibilityAllowance"`
	AttendanceAllowance     float64 `json:"attendanceAllowance"`
	TotalAllowances         float64 `json:"totalAllowances"`
	GrossSalary             float64 `json:"grossSalary"`
	SocialInsurance         float64 `json:"socialInsurance"`
	HealthInsurance         float64 `json:"healthInsurance"`
	UnemploymentInsurance   float64 `json:"unemploymentInsurance"`
	TotalInsurance          float64 `json:"totalInsurance"`
	PersonalIncomeTax       float64 `json:"personalIncomeTax"`
	AdvanceDeduction        float64 `json:"advanceDeduction"`
	SafetyPenaltyDeduction  float64 `json:"safetyPenaltyDeduction"`
	TotalDeductions         float64 `json:"totalDeductions"`
	NetSalary               float64 `json:"netSalary"`
	BankName                string  `json:"bankName"`
	BankAccount             string  `json:"bankAccount"`
	Status                  string  `json:"status"`
	PaidDate                string  `json:"paidDate,omitempty"`
}
