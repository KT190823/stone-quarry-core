package services

import (
	"mo-da-backend/internal/repositories"
)

type AlertService struct {
	*BaseService
}

func NewAlertService() *AlertService {
	repo := repositories.NewAlertRepo()
	return &AlertService{BaseService: NewBaseService(repo.BaseRepo)}
}
