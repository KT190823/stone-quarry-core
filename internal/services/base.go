package services

import (
	"mo-da-backend/internal/repositories"
)

type ListParams = repositories.ListParams

type BaseService struct {
	repo *repositories.BaseRepo
}

func NewBaseService(repo *repositories.BaseRepo) *BaseService {
	return &BaseService{repo: repo}
}

func (s *BaseService) List(params repositories.ListParams) ([]map[string]interface{}, int, error) {
	return s.repo.List(params)
}

func (s *BaseService) GetByID(id string) (map[string]interface{}, error) {
	return s.repo.GetByID(id)
}

func (s *BaseService) Create(data map[string]interface{}) (map[string]interface{}, error) {
	return s.repo.Create(data)
}

func (s *BaseService) Update(id string, data map[string]interface{}) (map[string]interface{}, error) {
	return s.repo.Update(id, data)
}

func (s *BaseService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *BaseService) ListJSONB(params repositories.ListParams) ([]map[string]interface{}, int, error) {
	return s.repo.ListJSONB(params)
}
