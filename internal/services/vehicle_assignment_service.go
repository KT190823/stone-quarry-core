package services

import (
	"context"
	"fmt"
	"time"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/repositories"
)

type VehicleAssignmentService struct {
	*BaseService
	assignmentRepo *repositories.VehicleAssignmentRepo
}

func NewVehicleAssignmentService() *VehicleAssignmentService {
	repo := repositories.NewVehicleAssignmentRepo()
	return &VehicleAssignmentService{
		BaseService:    NewBaseService(repo.BaseRepo),
		assignmentRepo: repo,
	}
}

func (s *VehicleAssignmentService) GetByID(idOrCode string) (map[string]interface{}, error) {
	return s.assignmentRepo.GetByID(idOrCode)
}

func (s *VehicleAssignmentService) GetActiveByPlate(plate string) (map[string]interface{}, error) {
	return s.assignmentRepo.GetActiveByPlate(plate)
}

func (s *VehicleAssignmentService) CheckIn(ctx context.Context, idOrCode string, startOdo, startFuel float64) (map[string]interface{}, error) {
	nowStr := time.Now().Format("15:04 02/01/2006")
	query := `
		UPDATE vehicle_assignments
		SET status = 'checked_in',
		    checkin_time = $1,
		    start_odo = CASE WHEN $2 > 0 THEN $2 ELSE start_odo END,
		    start_fuel_liters = CASE WHEN $3 > 0 THEN $3 ELSE start_fuel_liters END,
		    updated_at = NOW()
		WHERE CAST(id AS TEXT) = $4 OR code ILIKE $4
		RETURNING id, code, license_plate, driver_name, status, checkin_time
	`
	var id int
	var code, plate, driver, status, checkin string
	err := database.Pool.QueryRow(ctx, query, nowStr, startOdo, startFuel, idOrCode).Scan(&id, &code, &plate, &driver, &status, &checkin)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy lệnh điều xe %s để chấm công nhận xe: %w", idOrCode, err)
	}

	// Update vehicle's current driver
	_, _ = database.Pool.Exec(ctx, `
		UPDATE vehicles
		SET current_driver_name = $1, status = 'Hoạt động trong ca'
		WHERE bs ILIKE $2
	`, driver, plate)

	return map[string]interface{}{
		"id":           id,
		"code":         code,
		"licensePlate": plate,
		"driverName":   driver,
		"status":       status,
		"checkinTime":  checkin,
		"message":      fmt.Sprintf("Tài xế %s đã chấm công nhận xe %s thành công!", driver, plate),
	}, nil
}

func (s *VehicleAssignmentService) CheckOut(ctx context.Context, idOrCode string, endOdo, endFuel float64) (map[string]interface{}, error) {
	nowStr := time.Now().Format("15:04 02/01/2006")
	query := `
		UPDATE vehicle_assignments
		SET status = 'completed',
		    checkout_time = $1,
		    end_odo = CASE WHEN $2 > 0 THEN $2 ELSE end_odo END,
		    end_fuel_liters = CASE WHEN $3 > 0 THEN $3 ELSE end_fuel_liters END,
		    updated_at = NOW()
		WHERE CAST(id AS TEXT) = $4 OR code ILIKE $4
		RETURNING id, code, license_plate, driver_name, status, checkout_time, trips_completed, tons_hauled
	`
	var id, trips int
	var code, plate, driver, status, checkout string
	var tons float64
	err := database.Pool.QueryRow(ctx, query, nowStr, endOdo, endFuel, idOrCode).Scan(&id, &code, &plate, &driver, &status, &checkout, &trips, &tons)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy lệnh điều xe %s để kết thúc ca: %w", idOrCode, err)
	}

	// Reset vehicle current driver
	_, _ = database.Pool.Exec(ctx, `
		UPDATE vehicles
		SET current_driver_name = '', status = 'Sẵn sàng nhận ca'
		WHERE bs ILIKE $1
	`, plate)

	return map[string]interface{}{
		"id":             id,
		"code":           code,
		"licensePlate":   plate,
		"driverName":     driver,
		"status":         status,
		"checkoutTime":   checkout,
		"tripsCompleted": trips,
		"tonsHauled":     tons,
		"message":        fmt.Sprintf("Tài xế %s đã bàn giao xe %s thành công. Hoàn thành %d chuyến, %0.1f Tấn!", driver, plate, trips, tons),
	}, nil
}
