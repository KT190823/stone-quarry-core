package services

import (
	"mo-da-backend/internal/repositories"
)

type InventoryInboundService struct {
	*BaseService
}

func NewInventoryInboundService() *InventoryInboundService {
	repo := repositories.NewInventoryInboundRepo()
	return &InventoryInboundService{BaseService: NewBaseService(repo.BaseRepo)}
}

type InventoryOutboundService struct {
	*BaseService
}

func NewInventoryOutboundService() *InventoryOutboundService {
	repo := repositories.NewInventoryOutboundRepo()
	return &InventoryOutboundService{BaseService: NewBaseService(repo.BaseRepo)}
}

type InventoryStocktakeService struct {
	*BaseService
}

func NewInventoryStocktakeService() *InventoryStocktakeService {
	repo := repositories.NewInventoryStocktakeRepo()
	return &InventoryStocktakeService{BaseService: NewBaseService(repo.BaseRepo)}
}

type InventoryMovementService struct {
	*BaseService
}

func NewInventoryMovementService() *InventoryMovementService {
	repo := repositories.NewInventoryMovementRepo()
	return &InventoryMovementService{BaseService: NewBaseService(repo.BaseRepo)}
}
