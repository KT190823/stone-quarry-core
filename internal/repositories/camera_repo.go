package repositories

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
