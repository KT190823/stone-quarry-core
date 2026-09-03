package services

import (
	"mo-da-backend/internal/repositories"
)

type VehicleService struct {
	*BaseService
	repo *repositories.VehicleRepo
}

func NewVehicleService() *VehicleService {
	repo := repositories.NewVehicleRepo()
	return &VehicleService{
		BaseService: NewBaseService(repo.BaseRepo),
		repo:        repo,
	}
}

func (s *VehicleService) GetByID(id string) (map[string]interface{}, error) {
	return s.repo.GetByID(id)
}
