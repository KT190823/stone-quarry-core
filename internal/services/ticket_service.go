package services

import (
	"mo-da-backend/internal/repositories"
)

type TicketService struct {
	*BaseService
}

func NewTicketService() *TicketService {
	repo := repositories.NewTicketRepo()
	return &TicketService{BaseService: NewBaseService(repo.BaseRepo)}
}
