package repositories

import (
	"context"
	"encoding/json"
	"mo-da-backend/internal/database"
)

type CameraDeviceRepo struct {
	*BaseRepo
}

func NewCameraDeviceRepo() *CameraDeviceRepo {
	return &CameraDeviceRepo{BaseRepo: NewBaseRepo("camera_devices", "id")}
}

type GoogleDriveRepo struct {
	*BaseRepo
}

func NewGoogleDriveRepo() *GoogleDriveRepo {
	return &GoogleDriveRepo{BaseRepo: NewBaseRepo("google_drive_files", "id")}
}

type GoogleEmailRepo struct {
	*BaseRepo
}

func NewGoogleEmailRepo() *GoogleEmailRepo {
	return &GoogleEmailRepo{BaseRepo: NewBaseRepo("google_emails", "id")}
}

type GoogleMapRepo struct {
	*BaseRepo
}

func NewGoogleMapRepo() *GoogleMapRepo {
	return &GoogleMapRepo{BaseRepo: NewBaseRepo("google_map_entries", "id")}
}

type GooglePhotoRepo struct {
	*BaseRepo
}

func NewGooglePhotoRepo() *GooglePhotoRepo {
	return &GooglePhotoRepo{BaseRepo: NewBaseRepo("google_photo_entries", "id")}
}

type GpsFleetRepo struct {
	*BaseRepo
}

func NewGpsFleetRepo() *GpsFleetRepo {
	return &GpsFleetRepo{BaseRepo: NewBaseRepo("gps_fleet", "id")}
}

type FuelTheftAuditRepo struct {
	*BaseRepo
}

func NewFuelTheftAuditRepo() *FuelTheftAuditRepo {
	return &FuelTheftAuditRepo{BaseRepo: NewBaseRepo("fuel_theft_audits", "id")}
}

func (r *FuelTheftAuditRepo) GetByID(idOrCode string) (map[string]interface{}, error) {
	ctx := context.Background()
	query := `SELECT COALESCE(data, row_to_json(t)::jsonb) FROM (
		SELECT * FROM fuel_theft_audits WHERE id::text = $1 OR code ILIKE $1 OR data->>'code' ILIKE $1 OR data->>'plate' ILIKE $1 ORDER BY id DESC LIMIT 1
	) t`

	var data []byte
	err := database.Pool.QueryRow(ctx, query, idOrCode).Scan(&data)
	if err != nil {
		return nil, err
	}

	var item map[string]interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}

	return item, nil
}


type FuelNormConfigRepo struct {
	*BaseRepo
}

func NewFuelNormConfigRepo() *FuelNormConfigRepo {
	return &FuelNormConfigRepo{BaseRepo: NewBaseRepo("fuel_norms", "id")}
}

type YardCheckRepo struct {
	*BaseRepo
}

func NewYardCheckRepo() *YardCheckRepo {
	return &YardCheckRepo{BaseRepo: NewBaseRepo("yard_checkinout", "id")}
}
