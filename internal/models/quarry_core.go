package models

import (
	"encoding/json"
	"time"
)

// Organization represents the top-level mining company / corporation
type Organization struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	TaxCode   string    `json:"taxCode"`
	Address   string    `json:"address"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Quarry represents a specific quarry mine site
type Quarry struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Longitude   float64   `json:"longitude"`
	Latitude    float64   `json:"latitude"`
	TotalArea   float64   `json:"totalArea"`   // m2 or ha
	BoundaryGeo string    `json:"boundaryGeo"` // GeoJSON or WKT string
	Status      string    `json:"status"`      // active, maintenance, closed
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Derived / populated fields
	AreasCount  int `json:"areasCount,omitempty"`
	CyclesCount int `json:"cyclesCount,omitempty"`
}

// QuarryArea represents an excavation zone, moong, bench, or stockpile area
type QuarryArea struct {
	ID          string    `json:"id"`
	QuarryID    string    `json:"quarryId"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BoundaryGeo string    `json:"boundaryGeo"`
	AreaM2      float64   `json:"areaM2"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// SurveyCycle represents a single 3D survey run (e.g. Drone RTK, LiDAR, GPS, Total Station)
type SurveyCycle struct {
	ID                string     `json:"id"`
	QuarryID          string     `json:"quarryId"`
	CycleCode         string     `json:"cycleCode"`
	SurveyDate        string     `json:"surveyDate"`
	SurveyStartedAt   *time.Time `json:"surveyStartedAt,omitempty"`
	SurveyCompletedAt *time.Time `json:"surveyCompletedAt,omitempty"`
	PreviousCycleID   *string    `json:"previousCycleId,omitempty"`
	SurveyMethod      string     `json:"surveyMethod"` // drone, lidar, gps, total_station
	Status            string     `json:"status"`       // processing, completed, failed
	OperatorName      string     `json:"operatorName"`
	ExternalID        string     `json:"externalId,omitempty"`
	ExternalSource    string     `json:"externalSource,omitempty"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`

	// Joined details for API response
	QuarryName     string              `json:"quarryName,omitempty"`
	QuarryCode     string              `json:"quarryCode,omitempty"`
	SurfaceModels  []SurfaceModel      `json:"surfaceModels,omitempty"`
	Calculations   []VolumeCalculation `json:"calculations,omitempty"`
	AreaResults    []SurveyAreaResult  `json:"areaResults,omitempty"`
	Reconciliation *ProductionReconciliation `json:"reconciliation,omitempty"`
}

// SurveyAreaResult stores elevation statistics per area in a survey cycle
type SurveyAreaResult struct {
	ID            string    `json:"id"`
	SurveyCycleID string    `json:"surveyCycleId"`
	QuarryAreaID  string    `json:"quarryAreaId"`
	AreaName      string    `json:"areaName,omitempty"`
	AreaCode      string    `json:"areaCode,omitempty"`
	AreaM2        float64   `json:"areaM2"`
	MinElevation  float64   `json:"minElevation"`
	MaxElevation  float64   `json:"maxElevation"`
	AvgElevation  float64   `json:"avgElevation"`
	CreatedAt     time.Time `json:"createdAt"`
}

// VolumeCalculation represents 3D volume computation results
type VolumeCalculation struct {
	ID                    string     `json:"id"`
	SurveyCycleID         string     `json:"surveyCycleId"`
	QuarryAreaID          string     `json:"quarryAreaId"`
	AreaName              string     `json:"areaName,omitempty"`
	AreaCode              string     `json:"areaCode,omitempty"`
	PreviousSurveyCycleID *string    `json:"previousSurveyCycleId,omitempty"`
	CalculationType       string     `json:"calculationType"` // remaining, extracted, stockpile, cut_fill
	PreviousVolumeM3      float64    `json:"previousVolumeM3"`
	CurrentVolumeM3       float64    `json:"currentVolumeM3"`
	ExtractedVolumeM3     float64    `json:"extractedVolumeM3"`
	FillVolumeM3          float64    `json:"fillVolumeM3"`
	NetVolumeM3           float64    `json:"netVolumeM3"`
	CalculationMethod     string     `json:"calculationMethod"` // tin_surface, voxel, cross_sections, mesh_diff
	ExternalID            string     `json:"externalId,omitempty"`
	CalculatedAt          *time.Time `json:"calculatedAt,omitempty"`
	CreatedAt             time.Time  `json:"createdAt"`
}

// SurfaceModel stores metadata for 3D assets (Point Cloud, Mesh, DTM, DSM, DEM, Orthophoto)
type SurfaceModel struct {
	ID               string    `json:"id"`
	SurveyCycleID    string    `json:"surveyCycleId"`
	QuarryAreaID     *string   `json:"quarryAreaId,omitempty"`
	ModelType        string    `json:"modelType"` // point_cloud, mesh, dtm, dsm, dem, orthophoto
	FileFormat       string    `json:"fileFormat"` // laz, las, glb, obj, tif, tiff, xyz
	StoragePath      string    `json:"storagePath"`
	FileSize         int64     `json:"fileSize"`
	CoordinateSystem string    `json:"coordinateSystem"` // EPSG:4326, VN2000, UTM48N
	MinX             float64   `json:"minX"`
	MinY             float64   `json:"minY"`
	MinZ             float64   `json:"minZ"`
	MaxX             float64   `json:"maxX"`
	MaxY             float64   `json:"maxY"`
	MaxZ             float64   `json:"maxZ"`
	CreatedAt        time.Time `json:"createdAt"`
}

// QuarryMaterial represents mineral/stone type and its specific gravity
type QuarryMaterial struct {
	ID               string    `json:"id"`
	QuarryID         *string   `json:"quarryId,omitempty"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	DensityTonPerM3  float64   `json:"densityTonPerM3"` // e.g. 1.55 ton/m3
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
}

// QuarryVehicle represents trucks transporting stone from the quarry
type QuarryVehicle struct {
	ID           string    `json:"id"`
	QuarryID     *string   `json:"quarryId,omitempty"`
	LicensePlate string    `json:"licensePlate"`
	VehicleType  string    `json:"vehicleType"`
	CapacityTon  float64   `json:"capacityTon"`
	Status       string    `json:"status"`
	ExternalID   string    `json:"externalId,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// WeighingTransaction represents a single weight bridge ticket transaction
type WeighingTransaction struct {
	ID             string          `json:"id"`
	QuarryID       string          `json:"quarryId"`
	VehicleID      *string         `json:"vehicleId,omitempty"`
	LicensePlate   string          `json:"licensePlate,omitempty"`
	MaterialID     *string         `json:"materialId,omitempty"`
	MaterialName   string          `json:"materialName,omitempty"`
	SurveyCycleID  *string         `json:"surveyCycleId,omitempty"`
	TicketNumber   string          `json:"ticketNumber"`
	WeightInKg     float64         `json:"weightInKg"`
	WeightOutKg    float64         `json:"weightOutKg"`
	NetWeightKg    float64         `json:"netWeightKg"`
	NetWeightTon   float64         `json:"netWeightTon"`
	WeighedInAt    *time.Time      `json:"weighedInAt,omitempty"`
	WeighedOutAt   *time.Time      `json:"weighedOutAt,omitempty"`
	Status         string          `json:"status"` // pending, completed, cancelled
	ExternalID     string          `json:"externalId,omitempty"`
	ExternalSource string          `json:"externalSource,omitempty"`
	RawData        json.RawMessage `json:"rawData,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`

	Images []VehicleTripImage `json:"images,omitempty"`
}

// VehicleTripImage stores snapshot photos taken at weighbridge
type VehicleTripImage struct {
	ID                    string     `json:"id"`
	WeighingTransactionID string     `json:"weighingTransactionId"`
	ImageType             string     `json:"imageType"` // front, rear, side, scale, license_plate
	StoragePath           string     `json:"storagePath"`
	CapturedAt            *time.Time `json:"capturedAt,omitempty"`
	CreatedAt             time.Time  `json:"createdAt"`
}

// ProductionReconciliation is the core admin report comparing 3D Volume vs Scale Weight
type ProductionReconciliation struct {
	ID                string    `json:"id"`
	QuarryID          string    `json:"quarryId"`
	QuarryName        string    `json:"quarryName,omitempty"`
	QuarryAreaID      *string   `json:"quarryAreaId,omitempty"`
	AreaName          string    `json:"areaName,omitempty"`
	SurveyCycleID     string    `json:"surveyCycleId"`
	CycleCode         string    `json:"cycleCode,omitempty"`
	MaterialID        *string   `json:"materialId,omitempty"`
	MaterialName      string    `json:"materialName,omitempty"`
	VolumeM3          float64   `json:"volumeM3"`
	DensityTonPerM3   float64   `json:"densityTonPerM3"`
	ExpectedWeightTon float64   `json:"expectedWeightTon"` // volumeM3 * densityTonPerM3
	ActualWeightTon   float64   `json:"actualWeightTon"`   // sum(netWeightKg) / 1000
	DifferenceTon     float64   `json:"differenceTon"`     // expected - actual
	DifferencePercent float64   `json:"differencePercent"` // (differenceTon / expectedWeightTon) * 100
	Status            string    `json:"status"`            // matched, warning, discrepancy, pending
	CreatedAt         time.Time `json:"createdAt"`
}

// QuarryAlert represents automated discrepancy warnings
type QuarryAlert struct {
	ID                string     `json:"id"`
	QuarryID          string     `json:"quarryId"`
	QuarryAreaID      *string    `json:"quarryAreaId,omitempty"`
	SurveyCycleID     *string    `json:"surveyCycleId,omitempty"`
	AlertType         string     `json:"alertType"` // volume_weight_discrepancy, tare_anomaly, unauthorized_truck
	Severity          string     `json:"severity"`  // info, warning, critical
	Title             string     `json:"title"`
	Message           string     `json:"message"`
	ExpectedValue     float64    `json:"expectedValue"`
	ActualValue       float64    `json:"actualValue"`
	DifferencePercent float64    `json:"differencePercent"`
	IsRead            bool       `json:"isRead"`
	IsResolved        bool       `json:"isResolved"`
	CreatedAt         time.Time  `json:"createdAt"`
	ResolvedAt        *time.Time `json:"resolvedAt,omitempty"`
}

// IntegrationEvent logs every external 3rd-party webhook/API hit
type IntegrationEvent struct {
	ID           string          `json:"id"`
	Source       string          `json:"source"`    // MineSurveySoftware, DroneRTK, WeighbridgeScale, ERP
	EventType    string          `json:"eventType"` // survey_completed, ticket_created, volume_calculated
	ExternalID   string          `json:"externalId"`
	Payload      json.RawMessage `json:"payload"`
	Status       string          `json:"status"` // received, processing, completed, failed
	ErrorMessage string          `json:"errorMessage,omitempty"`
	ReceivedAt   time.Time       `json:"receivedAt"`
	ProcessedAt  *time.Time      `json:"processedAt,omitempty"`
}

// WebhookRequest is the payload received from 3rd party mine software or drone survey
type WebhookRequest struct {
	Source     string          `json:"source"`     // e.g. "MineSurveySoftware", "DroneRTK", "SmartScale"
	EventType  string          `json:"eventType"`  // "survey_completed", "scale_ticket", "volume_sync"
	ExternalID string          `json:"externalId"` // e.g. "SURVEY-20260826-001"
	Payload    json.RawMessage `json:"payload"`
}

// WebhookSurveyPayload represents a parsed survey cycle payload from 3rd party
type WebhookSurveyPayload struct {
	QuarryCode     string `json:"quarryCode"`
	CycleCode      string `json:"cycleCode"`
	SurveyDate     string `json:"surveyDate"`
	SurveyMethod   string `json:"surveyMethod"`
	OperatorName   string `json:"operatorName"`
	Notes          string `json:"notes"`
	Areas          []WebhookAreaVolume `json:"areas"`
}

type WebhookAreaVolume struct {
	AreaCode          string               `json:"areaCode"`
	AreaName          string               `json:"areaName"`
	AreaM2            float64              `json:"areaM2"`
	MinElevation      float64              `json:"minElevation"`
	MaxElevation      float64              `json:"maxElevation"`
	AvgElevation      float64              `json:"avgElevation"`
	PreviousVolumeM3  float64              `json:"previousVolumeM3"`
	CurrentVolumeM3   float64              `json:"currentVolumeM3"`
	ExtractedVolumeM3 float64              `json:"extractedVolumeM3"`
	FillVolumeM3      float64              `json:"fillVolumeM3"`
	CalculationMethod string               `json:"calculationMethod"`
	MaterialCode      string               `json:"materialCode"`
	Models            []WebhookSurfaceModel `json:"models"`
}

type WebhookSurfaceModel struct {
	ModelType        string  `json:"modelType"`
	FileFormat       string  `json:"fileFormat"`
	StoragePath      string  `json:"storagePath"`
	FileSize         int64   `json:"fileSize"`
	CoordinateSystem string  `json:"coordinateSystem"`
	MinX             float64 `json:"minX"`
	MinY             float64 `json:"minY"`
	MinZ             float64 `json:"minZ"`
	MaxX             float64 `json:"maxX"`
	MaxY             float64 `json:"maxY"`
	MaxZ             float64 `json:"maxZ"`
}

type WebhookScalePayload struct {
	QuarryCode   string                 `json:"quarryCode"`
	TicketNumber string                 `json:"ticketNumber"`
	LicensePlate string                 `json:"licensePlate"`
	MaterialCode string                 `json:"materialCode"`
	WeightInKg   float64                `json:"weightInKg"`
	WeightOutKg  float64                `json:"weightOutKg"`
	NetWeightKg  float64                `json:"netWeightKg"`
	WeighedInAt  string                 `json:"weighedInAt"`
	WeighedOutAt string                 `json:"weighedOutAt"`
	RawDetails   map[string]interface{} `json:"rawDetails"`
	Images       []WebhookTripImage     `json:"images"`
}

type WebhookTripImage struct {
	ImageType   string `json:"imageType"`
	StoragePath string `json:"storagePath"`
}

// QuarryDashboardSummary provides high-level metrics for the entire quarry module
type QuarryDashboardSummary struct {
	TotalQuarries          int                       `json:"totalQuarries"`
	TotalActiveAreas       int                       `json:"totalActiveAreas"`
	TotalSurveyCycles      int                       `json:"totalSurveyCycles"`
	TotalExtractedVolumeM3 float64                   `json:"totalExtractedVolumeM3"`
	TotalActualWeightTon   float64                   `json:"totalActualWeightTon"`
	TotalExpectedWeightTon float64                   `json:"totalExpectedWeightTon"`
	OverallDifferenceTon   float64                   `json:"overallDifferenceTon"`
	OverallDifferencePct   float64                   `json:"overallDifferencePct"`
	ActiveAlertsCount      int                       `json:"activeAlertsCount"`
	RecentCycles           []SurveyCycle             `json:"recentCycles"`
	RecentReconciliations  []ProductionReconciliation `json:"recentReconciliations"`
	RecentIntegrationEvents []IntegrationEvent       `json:"recentIntegrationEvents"`
}
