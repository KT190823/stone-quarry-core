package services

import (
	"mo-da-backend/internal/repositories"
)

type TicketService struct {
	*BaseService
	repo *repositories.TicketRepo
}

func NewTicketService() *TicketService {
	repo := repositories.NewTicketRepo()
	return &TicketService{
		BaseService: NewBaseService(repo.BaseRepo),
		repo:        repo,
	}
}

func (s *TicketService) GetByID(id string) (map[string]interface{}, error) {
	return s.repo.GetByID(id)
}
