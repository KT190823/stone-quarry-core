package models

type MiningPermit struct {
	ID              string       `json:"id"`
	Code            string       `json:"code"`
	Title           string       `json:"title"`
	MineName        string       `json:"mineName"`
	Category        string       `json:"category"`
	CategoryLabel   string       `json:"categoryLabel"`
	Issuer          string       `json:"issuer"`
	LicenseNumber   string       `json:"licenseNumber"`
	IssueDate       string       `json:"issueDate"`
	ExpiryDate      string       `json:"expiryDate"`
	Capacity        string       `json:"capacity"`
	ApprovedReserve string       `json:"approvedReserve"`
	MinedSoFar      string       `json:"minedSoFar"`
	MinedPercent    int          `json:"minedPercent"`
	DepthLevel      string       `json:"depthLevel"`
	Area            string       `json:"area"`
	Coordinates     string       `json:"coordinates"`
	Status          string       `json:"status"`
	StatusLabel     string       `json:"statusLabel"`
	DaysRemaining   int          `json:"daysRemaining"`
	Files           []PermitFile `json:"files"`
	Notes           string       `json:"notes"`
}

type PermitFile struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type MiningPlanItem struct {
	ID             string  `json:"id"`
	Mine           string  `json:"mine"`
	Item           string  `json:"item"`
	AnnualTarget   float64 `json:"annualTarget"`
	Unit           string  `json:"unit"`
	Q1Plan         float64 `json:"q1Plan"`
	Q1Actual       float64 `json:"q1Actual"`
	Q2Plan         float64 `json:"q2Plan"`
	Q2Actual       float64 `json:"q2Actual"`
	Q3Plan         float64 `json:"q3Plan"`
	Q3Actual       float64 `json:"q3Actual"`
	Q4Plan         float64 `json:"q4Plan"`
	Q4Actual       float64 `json:"q4Actual"`
	YTDActual      float64 `json:"ytdActual"`
	CompletionRate float64 `json:"completionRate"`
	Status         string  `json:"status"`
	StatusLabel    string  `json:"statusLabel"`
}

type ProductionStage struct {
	StageNumber       int     `json:"stageNumber"`
	StageName         string  `json:"stageName"`
	Icon              string  `json:"icon"`
	VolumeMonth       float64 `json:"volumeMonth"`
	VolumeYTD         float64 `json:"volumeYtd"`
	Unit              string  `json:"unit"`
	LossRate          string  `json:"lossRate"`
	LossStatus        string  `json:"lossStatus"`
	MeasurementMethod string  `json:"measurementMethod"`
	Description       string  `json:"description"`
}

type StatutoryReport struct {
	ID           string  `json:"id"`
	Code         string  `json:"code"`
	Title        string  `json:"title"`
	Recipient    string  `json:"recipient"`
	Period       string  `json:"period"`
	Date         string  `json:"date"`
	MinedVolume  float64 `json:"minedVolume"`
	Unit         string  `json:"unit"`
	TaxAmount    float64 `json:"taxAmount"`
	EnvFeeAmount float64 `json:"envFeeAmount"`
	Status       string  `json:"status"`
	StatusLabel  string  `json:"statusLabel"`
}

type NaturalResourceTax struct {
	ID                string  `json:"id"`
	Code              string  `json:"code"`
	MineralType       string  `json:"mineralType"`
	MinedVolume       float64 `json:"minedVolume"`
	Unit              string  `json:"unit"`
	TaxPricePerUnit   float64 `json:"taxPricePerUnit"`
	TaxRate           string  `json:"taxRate"`
	ResourceTaxAmount float64 `json:"resourceTaxAmount"`
	EnvironmentalFee  float64 `json:"environmentalFee"`
	TotalPayable      float64 `json:"totalPayable"`
	Status            string  `json:"status"`
}

type BlastingPassport struct {
	ID                   string              `json:"id"`
	Code                 string              `json:"code"`
	MineName             string              `json:"mineName"`
	BlastDate            string              `json:"blastDate"`
	BlastTime            string              `json:"blastTime"`
	Location             string              `json:"location"`
	BenchLevel           string              `json:"benchLevel"`
	HoleCount            int                 `json:"holeCount"`
	HoleDepthMeters      float64             `json:"holeDepthMeters"`
	HoleDiameterMm       float64             `json:"holeDiameterMm"`
	BurdenM              float64             `json:"burdenM"`
	SpacingM             float64             `json:"spacingM"`
	RowSpacingM          float64             `json:"rowSpacingM"`
	StemmingLengthM      float64             `json:"stemmingLengthM"`
	RockHardnessF        string              `json:"rockHardnessF"`
	FiringMethod         string              `json:"firingMethod"`
	DelayIntervalMs      int                 `json:"delayIntervalMs"`
	AnfoExplosiveKg      float64             `json:"anfoExplosiveKg"`
	EmulsionExplosiveKg  float64             `json:"emulsionExplosiveKg"`
	DetonatorCount       int                 `json:"detonatorCount"`
	DesignedRockVolumeM3 float64             `json:"designedRockVolumeM3"`
	PowderFactorKgPerM3  float64             `json:"powderFactorKgPerM3"`
	ActualRockMinedM3    float64             `json:"actualRockMinedM3"`
	SafetyStatus         string              `json:"safetyStatus"`
	StatusCode           string              `json:"statusCode"`
	BlasterInCharge      string              `json:"blasterInCharge"`
	CertifiedNumber      string              `json:"certifiedNumber"`
	CertificateExpiry    string              `json:"certificateExpiry"`
	ApprovedBy           string              `json:"approvedBy"`
	ApprovedAt           string              `json:"approvedAt"`
	
	// Pháp lý & Quota Sở Công Thương
	LicenseNumber        string              `json:"licenseNumber"`
	QuotaAnnualLimitTons float64             `json:"quotaAnnualLimitTons"`
	QuotaUsedYtdTons     float64             `json:"quotaUsedYtdTons"`
	QuotaRemainingTons   float64             `json:"quotaRemainingTons"`
	
	// An toàn cảnh giới & Thông báo chính quyền (QCVN 01:2019/BCT)
	SafetyPerimeterM     int                 `json:"safetyPerimeterM"`
	PoliceNotified       bool                `json:"policeNotified"`
	CommuneNotified      bool                `json:"communeNotified"`
	NotificationDocRef   string              `json:"notificationDocRef"`
	AllGuardsConfirmed   bool                `json:"allGuardsConfirmed"`
	EvacuationConfirmed  bool                `json:"evacuationConfirmed"`
	SirenAlertsCompleted bool                `json:"sirenAlertsCompleted"`
	GuardPosts           []BlastGuardPost    `json:"guardPosts"`
	
	// Quản lý vật liệu nổ công nghiệp (Theo Lô & Hoàn trả trong ngày)
	Materials            []BlastMaterialItem `json:"materials"`
	SameDayReturnClosed  bool                `json:"sameDayReturnClosed"`
	ReturnCompletedTime  string              `json:"returnCompletedTime"`
	
	// Kiểm tra sau nổ & Mìn câm (Misfire)
	PostBlastClearance   bool                `json:"postBlastClearance"`
	SmokeClearingMinutes int                 `json:"smokeClearingMinutes"`
	MisfireReported      bool                `json:"misfireReported"`
	MisfireCount         int                 `json:"misfireCount"`
	MisfireDetails       string              `json:"misfireDetails"`
	MisfireResolution    string              `json:"misfireResolution"`
	
	// Bảng tính giá thành đợt nổ (Cost Sheet)
	ExplosiveCost        float64             `json:"explosiveCost"`
	DrillingCost         float64             `json:"drillingCost"`
	BlastingServiceFee   float64             `json:"blastingServiceFee"`
	GuardLaborCost       float64             `json:"guardLaborCost"`
	TotalBlastCost       float64             `json:"totalBlastCost"`
	CostPerM3Rock        float64             `json:"costPerM3Rock"`
	CostPerTonRock       float64             `json:"costPerTonRock"`
	Notes                string              `json:"notes"`
}

type BlastGuardPost struct {
	PostName         string `json:"postName"`
	GuardPerson      string `json:"guardPerson"`
	Phone            string `json:"phone"`
	DistanceFromPitM int    `json:"distanceFromPitM"`
	Status           string `json:"status"`
	CheckedInAt      string `json:"checkedInAt"`
}

type BlastMaterialItem struct {
	MaterialName  string  `json:"materialName"`
	Category      string  `json:"category"`
	BatchLotNo    string  `json:"batchLotNo"`
	Supplier      string  `json:"supplier"`
	PlannedQty    float64 `json:"plannedQty"`
	IssuedQty     float64 `json:"issuedQty"`
	ActualUsedQty float64 `json:"actualUsedQty"`
	ReturnedQty   float64 `json:"returnedQty"`
	Discrepancy   float64 `json:"discrepancy"`
	Unit          string  `json:"unit"`
	ExpiryDate    string  `json:"expiryDate"`
	ReturnStatus  string  `json:"returnStatus"`
}

type YieldBreakdown struct {
	ProductName  string  `json:"productName"`
	ProducedTons float64 `json:"producedTons"`
	YieldPercent float64 `json:"yieldPercent"`
	StandardRate float64 `json:"standardRate"`
}

type CrusherPlantMatrix struct {
	PlantCode           string           `json:"plantCode"`
	PlantName           string           `json:"plantName"`
	CapacityTonPerHour  int              `json:"capacityTonPerHour"`
	InputRockTodayTons  int              `json:"inputRockTodayTons"`
	PowerConsumptionKwh int              `json:"powerConsumptionKwh"`
	KwhPerTon           float64          `json:"kwhPerTon"`
	YieldBreakdown      []YieldBreakdown `json:"yieldBreakdown"`
}

type HeavyEquipmentFuelLog struct {
	ID                       string  `json:"id"`
	EquipmentCode            string  `json:"equipmentCode"`
	EquipmentName            string  `json:"equipmentName"`
	Category                 string  `json:"category"`
	OperatorName             string  `json:"operatorName"`
	HoursWorkedToday         float64 `json:"hoursWorkedToday"`
	TotalHoursMeter          float64 `json:"totalHoursMeter"`
	FuelQuotaLitersPerHour   float64 `json:"fuelQuotaLitersPerHour"`
	ActualFuelIssuedLiters   float64 `json:"actualFuelIssuedLiters"`
	ActualFuelConsumedLiters float64 `json:"actualFuelConsumedLiters"`
	FuelVarianceLiters       float64 `json:"fuelVarianceLiters"`
	VarianceStatus           string  `json:"varianceStatus"`
	Location                 string  `json:"location"`
	MaintenanceStatus        string  `json:"maintenanceStatus"`
}
