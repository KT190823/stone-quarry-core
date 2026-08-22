package services

import (
	"mo-da-backend/internal/repositories"
)

type VehicleService struct {
	*BaseService
}

func NewVehicleService() *VehicleService {
	repo := repositories.NewVehicleRepo()
	return &VehicleService{BaseService: NewBaseService(repo.BaseRepo)}
}
