package services

import (
	"mo-da-backend/internal/repositories"
)

type UserService struct {
	*BaseService
}

func NewUserService() *UserService {
	repo := repositories.NewUserRepo()
	return &UserService{BaseService: NewBaseService(repo.BaseRepo)}
}

type UserRoleService struct {
	*BaseService
}

func NewUserRoleService() *UserRoleService {
	repo := repositories.NewUserRoleRepo()
	return &UserRoleService{BaseService: NewBaseService(repo.BaseRepo)}
}

type UserLogService struct {
	*BaseService
}

func NewUserLogService() *UserLogService {
	repo := repositories.NewUserLogRepo()
	return &UserLogService{BaseService: NewBaseService(repo.BaseRepo)}
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
