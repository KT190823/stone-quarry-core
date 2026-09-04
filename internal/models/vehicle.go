package models

type Vehicle struct {
	BS          string  `json:"bs"`
	Loai        string  `json:"loai"`
	BI          float64 `json:"bi"`
	RFID        string  `json:"rfid"`
	TaiTrong    float64 `json:"taiTrong"`
	Unit        string  `json:"unit,omitempty"`
	Status      string  `json:"status"`
	Count       int     `json:"count"`
	HanDangKiem string  `json:"hanDangKiem,omitempty"`
	Date        string  `json:"date,omitempty"`
	ChuXe       string  `json:"chuXe,omitempty"`
}

type Alert struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	BS         string  `json:"bs"`
	Note       string  `json:"note"`
	Time       string  `json:"time"`
	Date       string  `json:"date,omitempty"`
	Status     string  `json:"status"`
	Severity   string  `json:"severity"`
	Phieu      string  `json:"phieu,omitempty"`
	BiDangKy   float64 `json:"biDangKy,omitempty"`
	BiThucTe   float64 `json:"biThucTe,omitempty"`
	LechBi     float64 `json:"lechBi,omitempty"`
	Cam        string  `json:"cam,omitempty"`
}

type SpeedLogEntry struct {
	Time  string `json:"time"`
	Speed int    `json:"speed"`
	RPM   int    `json:"rpm"`
}

type HistoryEvent struct {
	Time string `json:"time"`
	Text string `json:"text"`
	Type string `json:"type"`
}

type DriverGpsTelemetry struct {
	ID                     string          `json:"id"`
	Plate                  string          `json:"plate"`
	DriverName             string          `json:"driverName"`
	DriverPhone            string          `json:"driverPhone"`
	Avatar                 string          `json:"avatar"`
	VehicleType            string          `json:"vehicleType"`
	VehicleModel           string          `json:"vehicleModel"`
	Lat                    float64         `json:"lat"`
	Lng                    float64         `json:"lng"`
	MapX                   int             `json:"mapX"`
	MapY                   int             `json:"mapY"`
	LocationName           string          `json:"locationName"`
	Destination            string          `json:"destination"`
	DestinationDistance    string          `json:"destinationDistance"`
	ETA                    string          `json:"eta"`
	Speed                  int             `json:"speed"`
	MaxSpeedLimit          int             `json:"maxSpeedLimit"`
	EngineStatus           string          `json:"engineStatus"`
	EngineRPM              int             `json:"engineRpm"`
	FuelLevelPercent       int             `json:"fuelLevelPercent"`
	FuelTankLiters         int             `json:"fuelTankLiters"`
	LoadStatus             string          `json:"loadStatus"`
	LoadWeight             string          `json:"loadWeight"`
	Cargo                  string          `json:"cargo"`
	Status                 string          `json:"status"`
	StatusLabel            string          `json:"statusLabel"`
	StatusDescription      string          `json:"statusDescription"`
	IdlingDurationMinutes  int             `json:"idlingDurationMinutes"`
	IdlingFuelWastedLiters float64         `json:"idlingFuelWastedLiters"`
	IdlingCostVnd          int             `json:"idlingCostVnd"`
	IdlingWarning          string          `json:"idlingWarning,omitempty"`
	GeofenceZone           string          `json:"geofenceZone"`
	TodayTrips             int             `json:"todayTrips"`
	TodayKm                float64         `json:"todayKm"`
	ActiveShift            string          `json:"activeShift"`
	UpdatedTime            string          `json:"updatedTime"`
	SpeedLog               []SpeedLogEntry `json:"speedLog"`
	HistoryEvents          []HistoryEvent  `json:"historyEvents"`
}

type FuelTheftAudit struct {
	Code            string `json:"code"`
	Vehicle         string `json:"vehicle"`
	Driver          string `json:"driver"`
	Shift           string `json:"shift"`
	QuotaLiters     string `json:"quotaLiters"`
	ActualIssued    string `json:"actualIssued"`
	ActualConsumed  string `json:"actualConsumed"`
	Variance        string `json:"variance"`
	VariancePercent string `json:"variancePercent"`
	ExpectedKM      string `json:"expectedKm"`
	ActualKM        string `json:"actualKm"`
	KMPerLiter      string `json:"kmPerLiter"`
	IdlingMinutes   string `json:"idlingMinutes"`
	IdlingWaste     string `json:"idlingWaste"`
	GPSAnomaly      string `json:"gpsAnomaly"`
	Status          string `json:"status"`
	Note            string `json:"note"`
}

type FuelNormConfig struct {
	Code                string `json:"code"`
	VehicleType         string `json:"vehicleType"`
	VehicleLabel        string `json:"vehicleLabel"`
	NormLitersPer100KM  string `json:"normLitersPer100Km"`
	NormLitersPerHour   string `json:"normLitersPerHour"`
	TankCapacity        string `json:"tankCapacity"`
	IdlingThreshold     string `json:"idlingThreshold"`
	AlertThreshold      string `json:"alertThreshold"`
	Status              string `json:"status"`
}

type YardCheckInOut struct {
	Code       string `json:"code"`
	Plate      string `json:"plate"`
	Driver     string `json:"driver"`
	VehicleType string `json:"vehicleType"`
	CheckType  string `json:"checkType"`
	Time       string `json:"time"`
	Gate       string `json:"gate"`
	Purpose    string `json:"purpose"`
	RFID       string `json:"rfid"`
	Method     string `json:"method"`
	Guard      string `json:"guard"`
	Status     string `json:"status"`
	Note       string `json:"note"`
}
