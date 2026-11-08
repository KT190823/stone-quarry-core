package database

import (
	"context"
	"encoding/json"
	"log"
)

type menuItemSeed struct {
	Label string `json:"label"`
	Route string `json:"route"`
	Badge string `json:"badge,omitempty"`
}

type menuModuleSeed struct {
	SectionKey   string
	SectionTitle string
	SectionOrder int
	ModuleKey    string
	Title        string
	Short        string
	Subtitle     string
	Icon         string
	Tone         string
	Route        string
	Dashboard    string
	ModuleOrder  int
	Children     []menuItemSeed
}

// menuModuleSeeds mirrors the former hard-coded frontend MODULE_SECTIONS.
// The database is now the single source of truth; frontend only renders the API.
var menuModuleSeeds = []menuModuleSeed{
	{
		SectionKey: "ban-dieu-hanh", SectionTitle: "BAN ĐIỀU HÀNH DOANH NGHIỆP (EXECUTIVE)", SectionOrder: 1,
		ModuleKey: "dieu-hanh", Title: "Cổng Điều Hành Doanh Nghiệp (Executive Portal)", Short: "Điều hành",
		Subtitle: "Trung tâm chỉ huy dành cho Ban Giám đốc: Chỉ số sức khỏe doanh nghiệp, phê duyệt chỉ đạo và giám sát liên phòng ban.",
		Icon:     "dashboard", Tone: "indigo", Route: "/dieu-hanh", Dashboard: "/dieu-hanh/tong-quan", ModuleOrder: 1,
		Children: []menuItemSeed{
			{Label: "Tổng quan điều hành", Route: "/dieu-hanh/tong-quan"},
			{Label: "Việc cần xử lý", Route: "/dieu-hanh/viec-can-xu-ly"},
			{Label: "Chỉ số các phòng ban", Route: "/dieu-hanh/phong-ban"},
			{Label: "Cảnh báo trọng yếu", Route: "/dieu-hanh/canh-bao"},
			{Label: "Báo cáo điều hành", Route: "/dieu-hanh/bao-cao"},
			{Label: "So sánh các mỏ", Route: "/dieu-hanh/so-sanh-mo"},
		},
	},
	{
		SectionKey: "kinh-doanh-ke-toan-nhan-su", SectionTitle: "KINH DOANH · KẾ TOÁN · NHÂN SỰ", SectionOrder: 2,
		ModuleKey: "doi-tac", Title: "Quản trị CRM & Kinh doanh", Short: "Kinh doanh",
		Subtitle: "Quản lý quan hệ khách hàng tiêu thụ đá, nhà cung cấp, đối tác vận tải và hợp đồng thương mại.",
		Icon:     "handshake", Tone: "rose", Route: "/doi-tac", Dashboard: "/doi-tac/tong-quan", ModuleOrder: 1,
		Children: []menuItemSeed{
			{Label: "Tổng quan", Route: "/doi-tac/tong-quan"},
			{Label: "Khách hàng tiêu thụ", Route: "/doi-tac/khach-hang"},
			{Label: "Nhà cung cấp & Đội xe", Route: "/doi-tac/nha-cung-cap"},
			{Label: "Hợp đồng thương mại", Route: "/doi-tac/hop-dong"},
		},
	},
	{
		SectionKey: "kinh-doanh-ke-toan-nhan-su", SectionTitle: "KINH DOANH · KẾ TOÁN · NHÂN SỰ", SectionOrder: 2,
		ModuleKey: "thanh-toan", Title: "Kế toán - Tài chính", Short: "Kế toán",
		Subtitle: "Hóa đơn VAT điện tử, chứng từ bán hàng, mua hàng, trả hàng, sổ quỹ thu chi và công nợ đối soát.",
		Icon:     "credit-card", Tone: "emerald", Route: "/thanh-toan", Dashboard: "/thanh-toan/tong-quan", ModuleOrder: 2,
		Children: []menuItemSeed{
			{Label: "Tổng quan", Route: "/thanh-toan/tong-quan"},
			{Label: "Phiếu bán hàng", Route: "/thanh-toan/phieu-ban"},
			{Label: "Phiếu mua hàng", Route: "/thanh-toan/phieu-mua"},
			{Label: "Phiếu trả hàng", Route: "/thanh-toan/phieu-tra"},
			{Label: "Phiếu thu tiền", Route: "/thanh-toan/phieu-thu"},
			{Label: "Phiếu chi tiền", Route: "/thanh-toan/phieu-chi"},
			{Label: "Hóa đơn điện tử VAT", Route: "/thanh-toan/hoa-don-vat"},
			{Label: "Sổ theo dõi công nợ", Route: "/thanh-toan/cong-no"},
			{Label: "Đối soát & Sổ quỹ", Route: "/thanh-toan/doi-soat"},
		},
	},
	{
		SectionKey: "kinh-doanh-ke-toan-nhan-su", SectionTitle: "KINH DOANH · KẾ TOÁN · NHÂN SỰ", SectionOrder: 2,
		ModuleKey: "nhan-su", Title: "Quản trị Nhân sự (HRM)", Short: "Nhân sự",
		Subtitle: "Hồ sơ nhân sự, chấm công máy nhận diện, bảng lương theo chuyến xe khoán và KPI toàn mỏ.",
		Icon:     "user-check", Tone: "purple", Route: "/nhan-su", Dashboard: "/nhan-su/tong-quan", ModuleOrder: 3,
		Children: []menuItemSeed{
			{Label: "Tổng quan", Route: "/nhan-su/tong-quan"},
			{Label: "Hồ sơ nhân viên", Route: "/nhan-su/nhan-vien"},
			{Label: "Chấm công & Ca trực", Route: "/nhan-su/cham-cong"},
			{Label: "Đơn từ & Nghỉ phép", Route: "/nhan-su/don-tu", Badge: "1"},
			{Label: "Bảng chấm công", Route: "/nhan-su/bang-cham-cong"},
			{Label: "Bảng tính lương", Route: "/nhan-su/bang-luong"},
			{Label: "Tuyển dụng nhân sự", Route: "/nhan-su/tuyen-dung"},
			{Label: "Hồ sơ ứng viên", Route: "/nhan-su/ung-vien"},
			{Label: "Quyết định nhân sự", Route: "/nhan-su/quyet-dinh"},
			{Label: "Bảo hiểm xã hội", Route: "/nhan-su/bao-hiem"},
			{Label: "Đào tạo an toàn mỏ", Route: "/nhan-su/dao-tao"},
			{Label: "Đánh giá KPI / OKR", Route: "/nhan-su/kpi-okr"},
			{Label: "Vòng đời nhân sự", Route: "/nhan-su/vong-doi-nhan-vien"},
			{Label: "Quy trình Workflow", Route: "/nhan-su/workflow"},
		},
	},
	{
		SectionKey: "phap-ly", SectionTitle: "PHÁP LÝ DOANH NGHIỆP", SectionOrder: 3,
		ModuleKey: "ho-so-giay-phep", Title: "Pháp lý & Giấy phép Mỏ", Short: "Pháp lý",
		Subtitle: "Giấy phép khai thác khoáng sản, phê duyệt ĐTM môi trường, trữ lượng mỏ và mốc giới quy hoạch.",
		Icon:     "building", Tone: "red", Route: "/ho-so-giay-phep", Dashboard: "/ho-so-giay-phep/tong-quan", ModuleOrder: 1,
		Children: []menuItemSeed{
			{Label: "Tổng quan", Route: "/ho-so-giay-phep/tong-quan"},
			{Label: "Giấy phép khai thác", Route: "/ho-so-giay-phep/giay-phep"},
			{Label: "Hồ sơ pháp lý mỏ", Route: "/ho-so-giay-phep/phap-ly"},
			{Label: "Báo cáo ĐTM & Môi trường", Route: "/ho-so-giay-phep/moi-truong"},
			{Label: "Bản đồ mốc giới & Trữ lượng", Route: "/ho-so-giay-phep/tru-luong-ban-do"},
		},
	},
	{
		SectionKey: "san-xuat-ky-thuat", SectionTitle: "SẢN XUẤT - KỸ THUẬT", SectionOrder: 4,
		ModuleKey: "ke-hoach-san-luong", Title: "Kế hoạch Kỹ thuật Mỏ", Short: "Kế hoạch",
		Subtitle: "Lập kế hoạch khai thác, chuỗi sản lượng, báo cáo Sở TN&MT và thuế tài nguyên khoáng sản.",
		Icon:     "trending-up", Tone: "green", Route: "/ke-hoach-san-luong", Dashboard: "/ke-hoach-san-luong/tong-quan", ModuleOrder: 1,
		Children: []menuItemSeed{
			{Label: "Tổng quan", Route: "/ke-hoach-san-luong/tong-quan"},
			{Label: "Kế hoạch khai thác", Route: "/ke-hoach-san-luong"},
			{Label: "Chuỗi sản lượng mỏ", Route: "/ke-hoach-san-luong/chuoi-san-luong"},
			{Label: "Báo cáo Sở TN&MT", Route: "/ke-hoach-san-luong/bao-cao-so-tnmt"},
			{Label: "Thuế & Phí tài nguyên", Route: "/ke-hoach-san-luong/thue-tai-nguyen"},
		},
	},
	{
		SectionKey: "san-xuat-ky-thuat", SectionTitle: "SẢN XUẤT - KỸ THUẬT", SectionOrder: 4,
		ModuleKey: "khai-thac-co-gioi", Title: "Khai thác & Chế biến Đá", Short: "Khai thác",
		Subtitle: "Quản lý hộ chiếu nổ mìn VLNCN, định mức dầu máy xúc/khoan và dây chuyền trạm nghiền sàng.",
		Icon:     "pickaxe", Tone: "amber", Route: "/khai-thac-co-gioi", Dashboard: "/khai-thac-co-gioi/tong-quan", ModuleOrder: 2,
		Children: []menuItemSeed{
			{Label: "Tổng quan", Route: "/khai-thac-co-gioi/tong-quan"},
			{Label: "Nổ mìn & VLNCN", Route: "/khai-thac-co-gioi/no-min-vlncn"},
			{Label: "Cơ giới mỏ & Dầu DO", Route: "/khai-thac-co-gioi/co-gioi-dau-do"},
			{Label: "Dây chuyền nghiền sàng", Route: "/khai-thac-co-gioi/tram-nghien-sang"},
		},
	},
	{
		SectionKey: "san-xuat-ky-thuat", SectionTitle: "SẢN XUẤT - KỸ THUẬT", SectionOrder: 4,
		ModuleKey: "khao-sat-3d", Title: "Trắc địa & Khảo sát 3D Mỏ", Short: "Khảo sát 3D",
		Subtitle: "Quét đám mây điểm 3D Drone RTK, LiDAR, tính toán thể tích m³ và đối soát cân xe.",
		Icon:     "layers", Tone: "indigo", Route: "/khao-sat-3d", Dashboard: "/khao-sat-3d/tong-quan", ModuleOrder: 3,
		Children: []menuItemSeed{
			{Label: "Tổng quan", Route: "/khao-sat-3d/tong-quan"},
			{Label: "Chu kỳ đo đạc", Route: "/khao-sat-3d/chu-ky-do-dac"},
			{Label: "Mô hình 3D mỏ", Route: "/khao-sat-3d/mo-hinh-3d"},
			{Label: "Đối soát sản lượng", Route: "/khao-sat-3d/doi-soat-san-luong"},
			{Label: "Trạm cân AI", Route: "/khao-sat-3d/tram-can-ai"},
		},
	},
	{
		SectionKey: "san-xuat-ky-thuat", SectionTitle: "SẢN XUẤT - KỸ THUẬT", SectionOrder: 4,
		ModuleKey: "mo-da", Title: "Danh mục Mỏ Đá Khai Thác", Short: "Mỏ đá",
		Subtitle: "Hồ sơ các mỏ đá trực thuộc, vị trí địa lý, trữ lượng thiết kế và công suất cấp phép.",
		Icon:     "map-pin", Tone: "cyan", Route: "/danh-muc/mo-da", Dashboard: "/danh-muc/mo-da", ModuleOrder: 4,
		Children: []menuItemSeed{
			{Label: "Danh sách mỏ đá", Route: "/danh-muc/mo-da"},
		},
	},
	{
		SectionKey: "kho-van-tai-tram-can", SectionTitle: "KHO - VẬN TẢI - TRẠM CÂN", SectionOrder: 5,
		ModuleKey: "quan-ly-kho", Title: "Quản lý Kho bãi & Bãi Đá", Short: "Kho bãi",
		Subtitle: "Theo dõi tồn kho bãi đá thành phẩm, nhập – xuất kho, kiểm kê và điều chuyển giữa các kho bãi.",
		Icon:     "layers", Tone: "cyan", Route: "/kho", Dashboard: "/kho/tong-quan", ModuleOrder: 1,
		Children: []menuItemSeed{
			{Label: "Tổng quan", Route: "/kho/tong-quan"},
			{Label: "Sản phẩm", Route: "/kho/san-pham"},
			{Label: "Phiếu nhập kho", Route: "/kho/nhap"},
			{Label: "Phiếu xuất kho", Route: "/kho/xuat"},
			{Label: "Kiểm kê bãi đá", Route: "/kho/kiem-ke"},
			{Label: "Dịch chuyển kho", Route: "/kho/dieu-chuyen"},
		},
	},
	{
		SectionKey: "kho-van-tai-tram-can", SectionTitle: "KHO - VẬN TẢI - TRẠM CÂN", SectionOrder: 5,
		ModuleKey: "quan-ly-xe-nhien-lieu", Title: "Đội xe & Hậu cần Vận tải", Short: "Đội xe",
		Subtitle: "Giám sát GPS live, phân công tài xế theo ca, đối soát tiêu hao và kiểm soát sai lệch nhiên liệu theo định mức.",
		Icon:     "truck", Tone: "indigo", Route: "/quan-ly-xe", Dashboard: "/quan-ly-xe/tong-quan", ModuleOrder: 2,
		Children: []menuItemSeed{
			{Label: "Tổng quan", Route: "/quan-ly-xe/tong-quan"},
			{Label: "GPS Giám sát Live", Route: "/quan-ly-xe/gps-live"},
			{Label: "Chấm công ca", Route: "/quan-ly-xe/check-in-out"},
			{Label: "Phân công xe", Route: "/quan-ly-xe/phan-cong-xe"},
			{Label: "Đối soát dầu", Route: "/quan-ly-xe/doi-soat-nhien-lieu", Badge: "1"},
			{Label: "Định mức tiêu hao xe", Route: "/quan-ly-xe/dinh-muc-xe"},
			{Label: "Hồ sơ xe & Thẻ RFID", Route: "/quan-ly-xe/danh-sach"},
		},
	},
	{
		SectionKey: "kho-van-tai-tram-can", SectionTitle: "KHO - VẬN TẢI - TRẠM CÂN", SectionOrder: 5,
		ModuleKey: "phieu-can", Title: "Trạm Cân Điện Tử AI", Short: "Trạm cân",
		Subtitle: "Điều hành trạm cân tự động, chụp ảnh biển số OCR AI, cân tổng bì và cân bì rỗng theo thời gian thực.",
		Icon:     "scale", Tone: "blue", Route: "/phieu-can", Dashboard: "/phieu-can", ModuleOrder: 3,
		Children: []menuItemSeed{
			{Label: "Danh sách phiếu cân", Route: "/phieu-can"},
			{Label: "Tạo phiếu cân mới", Route: "/phieu-can/form"},
			{Label: "Báo cáo trạm cân", Route: "/bao-cao"},
		},
	},
	{
		SectionKey: "kho-van-tai-tram-can", SectionTitle: "KHO - VẬN TẢI - TRẠM CÂN", SectionOrder: 5,
		ModuleKey: "canh-bao-lech-bi", Title: "Giám Sát Cảnh Báo Lệch Bì AI", Short: "Lệch bì",
		Subtitle: "Hệ thống AI tự động phát hiện xe cơi nới thùng gian lận, sai lệch tự trọng bì xe và sai số bàn cân.",
		Icon:     "alert-triangle", Tone: "amber", Route: "/canh-bao-lech-bi", Dashboard: "/canh-bao-lech-bi", ModuleOrder: 4,
		Children: []menuItemSeed{
			{Label: "Danh sách cảnh báo", Route: "/canh-bao-lech-bi"},
		},
	},
	{
		SectionKey: "he-thong-cong-cu", SectionTitle: "DỮ LIỆU - HỆ THỐNG", SectionOrder: 6,
		ModuleKey: "dich-vu-google", Title: "Dịch vụ Số & Bản đồ GPS", Short: "Google & GPS",
		Subtitle: "Google Maps lộ trình xe tải mỏ, Drive hồ sơ đám mây, Gmail thông báo và Photos AI trạm cân.",
		Icon:     "map-pin", Tone: "blue", Route: "/google", Dashboard: "/google/maps", ModuleOrder: 1,
		Children: []menuItemSeed{
			{Label: "Google Maps GPS", Route: "/google/maps"},
			{Label: "Google Drive hồ sơ", Route: "/google/drive"},
			{Label: "Hòm thư Gmail tự động", Route: "/google/gmail"},
			{Label: "Kho ảnh Photos AI", Route: "/google/photos"},
		},
	},
	{
		SectionKey: "he-thong-cong-cu", SectionTitle: "DỮ LIỆU - HỆ THỐNG", SectionOrder: 6,
		ModuleKey: "nguoi-dung", Title: "Người dùng & Phân quyền", Short: "Người dùng",
		Subtitle: "Tài khoản cán bộ nhân viên, nhóm quyền vai trò (RBAC) và nhật ký an ninh hệ thống.",
		Icon:     "users", Tone: "indigo", Route: "/nguoi-dung", Dashboard: "/nguoi-dung/danh-sach", ModuleOrder: 2,
		Children: []menuItemSeed{
			{Label: "Danh sách người dùng", Route: "/nguoi-dung/danh-sach"},
			{Label: "Nhóm quyền & Vai trò", Route: "/nguoi-dung/nhom-quyen"},
			{Label: "Nhật ký kiểm toán", Route: "/nguoi-dung/nhat-ky"},
		},
	},
	{
		SectionKey: "he-thong-cong-cu", SectionTitle: "DỮ LIỆU - HỆ THỐNG", SectionOrder: 6,
		ModuleKey: "danh-muc", Title: "Danh mục Dùng chung", Short: "Danh mục",
		Subtitle: "Chủng loại đá thương phẩm, loại phiếu cân, danh mục xe cơ giới và thiết bị trạm cân.",
		Icon:     "folder", Tone: "amber", Route: "/danh-muc", Dashboard: "/danh-muc/mo-da", ModuleOrder: 3,
		Children: []menuItemSeed{
			{Label: "Danh mục mỏ đá", Route: "/danh-muc/mo-da"},
			{Label: "Chủng loại hàng đá", Route: "/danh-muc/loai-hang-da"},
			{Label: "Loại chứng từ phiếu cân", Route: "/danh-muc/loai-phieu-can"},
			{Label: "Chủng loại xe cơ giới", Route: "/danh-muc/loai-xe"},
			{Label: "Danh mục trạm cân", Route: "/danh-muc/tram-can"},
			{Label: "Thiết bị & Barrier IoT", Route: "/danh-muc/thiet-bi"},
		},
	},
	{
		SectionKey: "he-thong-cong-cu", SectionTitle: "DỮ LIỆU - HỆ THỐNG", SectionOrder: 6,
		ModuleKey: "cau-hinh", Title: "Cấu hình & Tích hợp IoT", Short: "Cấu hình",
		Subtitle: "Tham số hệ thống trạm cân, camera AI biển số, đầu đọc thẻ RFID và barrier tự động.",
		Icon:     "settings", Tone: "slate", Route: "/cau-hinh", Dashboard: "/cau-hinh", ModuleOrder: 4,
		Children: []menuItemSeed{
			{Label: "Tham số hệ thống", Route: "/cau-hinh"},
			{Label: "Giám sát Camera AI", Route: "/camera"},
		},
	},
	{
		SectionKey: "he-thong-cong-cu", SectionTitle: "DỮ LIỆU - HỆ THỐNG", SectionOrder: 6,
		ModuleKey: "mau-in", Title: "Mẫu in & Biểu mẫu Chứng từ", Short: "Mẫu in",
		Subtitle: "Thiết kế và chuẩn hóa mẫu in phiếu cân 3 liên, hóa đơn VAT, phiếu xuất kho và biên bản giao nhận.",
		Icon:     "printer", Tone: "cyan", Route: "/mau-in", Dashboard: "/mau-in", ModuleOrder: 5,
		Children: []menuItemSeed{
			{Label: "Quản lý mẫu in", Route: "/mau-in"},
		},
	},
}

const (
	modExecutive = "dieu-hanh"
	modCRM       = "doi-tac"
	modFinance   = "thanh-toan"
	modHR        = "nhan-su"
	modLegal     = "ho-so-giay-phep"
	modPlan      = "ke-hoach-san-luong"
	modMining    = "khai-thac-co-gioi"
	modSurvey    = "khao-sat-3d"
	modQuarry    = "mo-da"
	modWarehouse = "quan-ly-kho"
	modFleet     = "quan-ly-xe-nhien-lieu"
	modScale     = "phieu-can"
	modAlerts    = "canh-bao-lech-bi"
	modGoogle    = "dich-vu-google"
	modUsers     = "nguoi-dung"
	modCatalog   = "danh-muc"
	modConfig    = "cau-hinh"
	modPrint     = "mau-in"
)

var allMenuModules = []string{
	modExecutive, modCRM, modFinance, modHR, modLegal, modPlan, modMining, modSurvey, modQuarry,
	modWarehouse, modFleet, modScale, modAlerts, modGoogle, modUsers, modCatalog, modConfig, modPrint,
}

// roleModuleSeeds is the default RBAC matrix. Existing rows are never overwritten,
// so administrators can tune access directly in role_menu_permissions.
var roleModuleSeeds = map[string][]string{
	"Ban Lãnh Đạo": allMenuModules,
	"Ban Giám Đốc": allMenuModules,
	"Giám Đốc Mỏ":  allMenuModules,
	"Phó Giám Đốc Kỹ Thuật": {
		modExecutive, modLegal, modPlan, modMining, modSurvey, modQuarry,
		modWarehouse, modFleet, modAlerts, modCatalog, modPrint,
	},
	"Trưởng Trạm Cân": {
		modWarehouse, modFleet, modScale, modAlerts, modCatalog, modPrint,
	},
	"Nhân viên vận hành bàn cân": {
		modScale, modAlerts, modPrint,
	},
	"Kỹ Sư Trắc Địa 3D": {
		modPlan, modSurvey, modQuarry, modCatalog,
	},
	"Chỉ Huy Nổ Mìn": {
		modMining, modWarehouse, modAlerts, modCatalog, modPrint,
	},
	"Thợ Khoan & An Toàn Mỏ": {
		modMining, modAlerts, modCatalog,
	},
	"Kế Toán Trưởng": {
		modCRM, modFinance, modHR, modWarehouse, modAlerts, modCatalog, modPrint,
	},
	"Kế toán công nợ & Hóa đơn": {
		modCRM, modFinance, modPrint,
	},
	"Thủ kho bãi khoáng sản": {
		modWarehouse, modScale, modAlerts, modCatalog, modPrint,
	},
	"Trưởng Phòng Kinh Doanh": {
		modCRM, modFinance, modWarehouse, modScale, modAlerts, modCatalog, modPrint,
	},
	"Điều Phối Vận Tải & Giám Sát GPS": {
		modWarehouse, modFleet, modScale, modAlerts, modGoogle, modCatalog,
	},
	"Cán Bộ Pháp Lý Mỏ & BVMT": {
		modLegal, modPlan, modQuarry, modAlerts, modCatalog, modPrint,
	},
}

// SeedMenuConfig provisions the navigation registry and default RBAC matrix.
func SeedMenuConfig(ctx context.Context) {
	if Pool == nil {
		return
	}

	_, _ = Pool.Exec(ctx, `UPDATE users SET client_id = 'TTC' WHERE client_id IS NULL OR client_id = ''`)
	_, _ = Pool.Exec(ctx, `UPDATE user_roles SET client_id = 'TTC' WHERE client_id IS NULL OR client_id = ''`)

	for _, m := range menuModuleSeeds {
		children, err := json.Marshal(m.Children)
		if err != nil {
			log.Printf("menu seed: marshal children for %s: %v", m.ModuleKey, err)
			continue
		}
		_, err = Pool.Exec(ctx, `
			INSERT INTO app_modules (
				client_id, section_key, section_title, section_order, module_key,
				title, short, subtitle, icon, tone, route, dashboard_route, module_order, children, status
			) VALUES ('TTC', $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, 'active')
			ON CONFLICT (client_id, module_key) DO NOTHING
		`, m.SectionKey, m.SectionTitle, m.SectionOrder, m.ModuleKey, m.Title, m.Short,
			m.Subtitle, m.Icon, m.Tone, m.Route, m.Dashboard, m.ModuleOrder, children)
		if err != nil {
			log.Printf("menu seed: upsert module %s: %v", m.ModuleKey, err)
		}
	}

	for role, moduleKeys := range roleModuleSeeds {
		for _, key := range moduleKeys {
			_, err := Pool.Exec(ctx, `
				INSERT INTO role_menu_permissions (client_id, role, module_key, can_view)
				VALUES ('TTC', $1, $2, TRUE)
				ON CONFLICT (client_id, role, module_key) DO NOTHING
			`, role, key)
			if err != nil {
				log.Printf("menu seed: permission %s/%s: %v", role, key, err)
			}
		}
	}
}
