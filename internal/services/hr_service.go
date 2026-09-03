package services

import (
	"mo-da-backend/internal/repositories"
)

type HrEmployeeService struct {
	*BaseService
}

func NewHrEmployeeService() *HrEmployeeService {
	repo := repositories.NewHrEmployeeRepo()
	return &HrEmployeeService{BaseService: NewBaseService(repo.BaseRepo)}
}

type HrContractService struct {
	*BaseService
	repo *repositories.HrContractRepo
}

func NewHrContractService() *HrContractService {
	repo := repositories.NewHrContractRepo()
	return &HrContractService{
		BaseService: NewBaseService(repo.BaseRepo),
		repo:        repo,
	}
}

func (s *HrContractService) GetByID(idOrCode string) (map[string]interface{}, error) {
	return s.repo.GetByID(idOrCode)
}

type HrShiftService struct {
	*BaseService
}

func NewHrShiftService() *HrShiftService {
	repo := repositories.NewHrShiftRepo()
	return &HrShiftService{BaseService: NewBaseService(repo.BaseRepo)}
}

type HrShiftScheduleService struct {
	*BaseService
}

func NewHrShiftScheduleService() *HrShiftScheduleService {
	repo := repositories.NewHrShiftScheduleRepo()
	return &HrShiftScheduleService{BaseService: NewBaseService(repo.BaseRepo)}
}

type HrAttendanceService struct {
	*BaseService
}

func NewHrAttendanceService() *HrAttendanceService {
	repo := repositories.NewHrAttendanceRepo()
	return &HrAttendanceService{BaseService: NewBaseService(repo.BaseRepo)}
}

type HrLeaveRequestService struct {
	*BaseService
}

func NewHrLeaveRequestService() *HrLeaveRequestService {
	repo := repositories.NewHrLeaveRequestRepo()
	return &HrLeaveRequestService{BaseService: NewBaseService(repo.BaseRepo)}
}

type HrLeaveAllocationService struct {
	*BaseService
}

func NewHrLeaveAllocationService() *HrLeaveAllocationService {
	repo := repositories.NewHrLeaveAllocationRepo()
	return &HrLeaveAllocationService{BaseService: NewBaseService(repo.BaseRepo)}
}

type HrPayslipService struct {
	*BaseService
}

func NewHrPayslipService() *HrPayslipService {
	repo := repositories.NewHrPayslipRepo()
	return &HrPayslipService{BaseService: NewBaseService(repo.BaseRepo)}
}

type HrTimesheetService struct {
	*BaseService
}

func NewHrTimesheetService() *HrTimesheetService {
	repo := repositories.NewHrTimesheetRepo()
	return &HrTimesheetService{BaseService: NewBaseService(repo.BaseRepo)}
}
