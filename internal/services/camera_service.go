package services

import (
	"mo-da-backend/internal/repositories"
)

type CameraDeviceService struct {
	*BaseService
}

func NewCameraDeviceService() *CameraDeviceService {
	repo := repositories.NewCameraDeviceRepo()
	return &CameraDeviceService{BaseService: NewBaseService(repo.BaseRepo)}
}

type GoogleDriveService struct {
	*BaseService
}

func NewGoogleDriveService() *GoogleDriveService {
	repo := repositories.NewGoogleDriveRepo()
	return &GoogleDriveService{BaseService: NewBaseService(repo.BaseRepo)}
}

type GoogleEmailService struct {
	*BaseService
}

func NewGoogleEmailService() *GoogleEmailService {
	repo := repositories.NewGoogleEmailRepo()
	return &GoogleEmailService{BaseService: NewBaseService(repo.BaseRepo)}
}

type GoogleMapService struct {
	*BaseService
}

func NewGoogleMapService() *GoogleMapService {
	repo := repositories.NewGoogleMapRepo()
	return &GoogleMapService{BaseService: NewBaseService(repo.BaseRepo)}
}

type GooglePhotoService struct {
	*BaseService
}

func NewGooglePhotoService() *GooglePhotoService {
	repo := repositories.NewGooglePhotoRepo()
	return &GooglePhotoService{BaseService: NewBaseService(repo.BaseRepo)}
}

type GpsFleetService struct {
	*BaseService
}

func NewGpsFleetService() *GpsFleetService {
	repo := repositories.NewGpsFleetRepo()
	return &GpsFleetService{BaseService: NewBaseService(repo.BaseRepo)}
}

type FuelTheftAuditService struct {
	*BaseService
	theftRepo *repositories.FuelTheftAuditRepo
}

func NewFuelTheftAuditService() *FuelTheftAuditService {
	repo := repositories.NewFuelTheftAuditRepo()
	return &FuelTheftAuditService{
		BaseService: NewBaseService(repo.BaseRepo),
		theftRepo:   repo,
	}
}

func (s *FuelTheftAuditService) GetByID(idOrCode string) (map[string]interface{}, error) {
	return s.theftRepo.GetByID(idOrCode)
}


type FuelNormConfigService struct {
	*BaseService
}

func NewFuelNormConfigService() *FuelNormConfigService {
	repo := repositories.NewFuelNormConfigRepo()
	return &FuelNormConfigService{BaseService: NewBaseService(repo.BaseRepo)}
}

type YardCheckService struct {
	*BaseService
}

func NewYardCheckService() *YardCheckService {
	repo := repositories.NewYardCheckRepo()
	return &YardCheckService{BaseService: NewBaseService(repo.BaseRepo)}
}
