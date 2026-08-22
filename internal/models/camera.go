package models

type CameraDevice struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Location string `json:"location"`
	IP       string `json:"ip"`
	Res      string `json:"res"`
	Status   string `json:"status"`
}

type GoogleDriveFile struct {
	Code     string `json:"code"`
	FileName string `json:"fileName"`
	Name     string `json:"name"`
	Folder   string `json:"folder"`
	Size     string `json:"size"`
	SyncedAt string `json:"syncedAt"`
	Updated  string `json:"updated"`
	Status   string `json:"status"`
}

type GoogleEmail struct {
	Code      string `json:"code"`
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Time      string `json:"time"`
	SentAt    string `json:"sentAt"`
	Status    string `json:"status"`
}

type GoogleMapEntry struct {
	Code      string `json:"code"`
	Plate     string `json:"plate"`
	Location  string `json:"location"`
	Route     string `json:"route"`
	Speed     string `json:"speed"`
	Time      string `json:"time"`
	UpdatedAt string `json:"updatedAt"`
	Status    string `json:"status"`
}

type GooglePhotoEntry struct {
	Code       string `json:"code"`
	Album      string `json:"album"`
	Name       string `json:"name"`
	Plate      string `json:"plate"`
	Count      string `json:"count"`
	Photos     string `json:"photos"`
	Time       string `json:"time"`
	CapturedAt string `json:"capturedAt"`
	Resolution string `json:"resolution"`
	Status     string `json:"status"`
}

type SystemAuditLog struct {
	ID              string   `json:"id"`
	Code            string   `json:"code"`
	Timestamp       string   `json:"timestamp"`
	Date            string   `json:"date"`
	Time            string   `json:"time"`
	UserID          string   `json:"userId"`
	Username        string   `json:"username"`
	FullName        string   `json:"fullName"`
	UserRole        string   `json:"userRole"`
	Department      string   `json:"department"`
	AvatarText      string   `json:"avatarText"`
	Module          string   `json:"module"`
	Action          string   `json:"action"`
	ActionType      string   `json:"actionType"`
	TargetObject    string   `json:"targetObject"`
	TargetID        string   `json:"targetId,omitempty"`
	IPAddress       string   `json:"ipAddress"`
	Location        string   `json:"location"`
	Device          string   `json:"device"`
	SecurityStatus  string   `json:"securityStatus"`
	Severity        string   `json:"severity"`
	Status          string   `json:"status"`
	ExecutionTimeMs int      `json:"executionTimeMs"`
	Details         string   `json:"details"`
	HashSha256      string   `json:"hashSha256"`
	Diffs           []string `json:"diffs,omitempty"`
}
