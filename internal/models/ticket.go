package models

type CameraInfo struct {
	Name          string `json:"name"`
	Location      string `json:"location"`
	Time          string `json:"time"`
	PlateDetected string `json:"plateDetected,omitempty"`
	Confidence    string `json:"confidence,omitempty"`
	Note          string `json:"note,omitempty"`
	Image         string `json:"image,omitempty"`
}

type TicketCameras struct {
	Cam1 CameraInfo `json:"cam1"`
	Cam2 CameraInfo `json:"cam2"`
	Cam3 CameraInfo `json:"cam3"`
	Cam4 CameraInfo `json:"cam4"`
}

type ChatterEntry struct {
	ID         string `json:"id"`
	Author     string `json:"author"`
	AvatarText string `json:"avatarText"`
	Time       string `json:"time"`
	Type       string `json:"type"`
	Content    string `json:"content"`
}

type Ticket struct {
	STT        int             `json:"stt"`
	ID         string          `json:"id"`
	BenBan     string          `json:"benBan"`
	BenMua     string          `json:"benMua"`
	BienSo     string          `json:"bienSo"`
	LoaiXe     string          `json:"loaiXe"`
	LaiXe      string          `json:"laiXe"`
	SdtLaiXe   string          `json:"sdtLaiXe"`
	RFID       string          `json:"rfid"`
	Loai       string          `json:"loai"`
	Stage      TicketStage     `json:"stage"`
	StageLabel string          `json:"stageLabel"`
	CanL1      float64         `json:"canL1"`
	KL1        interface{}     `json:"kl1"`
	CanL2      string          `json:"canL2"`
	KL2        interface{}     `json:"kl2"`
	KLHang     interface{}     `json:"klHang"`
	KLTapChat  float64         `json:"klTapChat"`
	KLTinhTien interface{}     `json:"klTinhTien"`
	DonGia     float64         `json:"donGia"`
	ThanhTien  interface{}     `json:"thanhTien"`
	Time1      string          `json:"time1"`
	Time2      string          `json:"time2"`
	Date       string          `json:"date"`
	NguoiCan1  string          `json:"nguoiCan1"`
	NguoiCan2  string          `json:"nguoiCan2"`
	MatHang    string          `json:"matHang"`
	QuyCach    string          `json:"quyCach"`
	Do         string          `json:"do"`
	TramCan    string          `json:"tramCan"`
	GhiChu           string          `json:"ghiChu"`
	HoaDonSo         string          `json:"hoaDonSo,omitempty"`
	VehicleOwnership string          `json:"vehicleOwnership,omitempty"` // "company", "customer", "contractor"
	DriverID         string          `json:"driverId,omitempty"`         // Mã nhân viên tài xế (e.g. EMP-DRV-01)
	DriverName       string          `json:"driverName,omitempty"`       // Họ tên tài xế
	DriverPhone      string          `json:"driverPhone,omitempty"`
	AssignmentCode   string          `json:"assignmentCode,omitempty"`   // Lệnh điều xe theo ca (e.g. DISPATCH-20261028-01)
	ShiftName        string          `json:"shiftName,omitempty"`        // Ca làm việc (e.g. Ca 1)
	AttendanceStatus string          `json:"attendanceStatus,omitempty"` // checked_in, in_progress, completed
	Cameras          TicketCameras   `json:"cameras"`
	Chatter          []ChatterEntry  `json:"chatter"`
}
