package services

import (
	"mo-da-backend/internal/repositories"
)

type UserService struct {
	*BaseService
	repo *repositories.UserRepo
}

func NewUserService() *UserService {
	repo := repositories.NewUserRepo()
	return &UserService{BaseService: NewBaseService(repo.BaseRepo), repo: repo}
}

func (s *UserService) GetByID(id string) (map[string]interface{}, error) {
	return s.repo.GetByID(id)
}

type UserRoleService struct {
	*BaseService
	repo *repositories.UserRoleRepo
}

func NewUserRoleService() *UserRoleService {
	repo := repositories.NewUserRoleRepo()
	return &UserRoleService{BaseService: NewBaseService(repo.BaseRepo), repo: repo}
}

func (s *UserRoleService) GetByID(id string) (map[string]interface{}, error) {
	return s.repo.GetByID(id)
}

type UserLogService struct {
	*BaseService
	repo *repositories.UserLogRepo
}

func NewUserLogService() *UserLogService {
	repo := repositories.NewUserLogRepo()
	return &UserLogService{BaseService: NewBaseService(repo.BaseRepo), repo: repo}
}

func (s *UserLogService) GetByID(id string) (map[string]interface{}, error) {
	return s.repo.GetByID(id)
}

type ReportService struct {
	*BaseService
}

func NewReportService() *ReportService {
	repo := repositories.NewReportRepo()
	return &ReportService{BaseService: NewBaseService(repo.BaseRepo)}
}

type SettingService struct {
	*BaseService
}

func NewSettingService() *SettingService {
	repo := repositories.NewSettingRepo()
	return &SettingService{BaseService: NewBaseService(repo.BaseRepo)}
}
