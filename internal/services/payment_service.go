package services

import (
	"mo-da-backend/internal/repositories"
)

type DebtService struct {
	*BaseService
}

func NewDebtService() *DebtService {
	repo := repositories.NewDebtRepo()
	return &DebtService{BaseService: NewBaseService(repo.BaseRepo)}
}

type ReconcileService struct {
	*BaseService
}

func NewReconcileService() *ReconcileService {
	repo := repositories.NewReconcileRepo()
	return &ReconcileService{BaseService: NewBaseService(repo.BaseRepo)}
}
