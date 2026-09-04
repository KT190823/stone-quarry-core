package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// printSeedDir must match the print template directory used by the services.
const printSeedDir = "./data/print_templates"

func writePrintLayoutFile(id, layout string) {
	if layout == "" {
		return
	}
	if err := os.MkdirAll(printSeedDir, 0o755); err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(printSeedDir, fmt.Sprintf("%s.json", id)), []byte(layout), 0o644)
}

// printBackfillLayoutToDisk writes any legacy DB-stored layout to a disk file,
// then clears the column so the DB holds metadata only.
func printBackfillLayoutToDisk(ctx context.Context) {
	rows, err := Pool.Query(ctx, `SELECT id, COALESCE(layout,'') FROM print_templates WHERE COALESCE(layout,'') <> ''`)
	if err != nil {
		return
	}
	defer rows.Close()
	type rowT struct {
		id, layout string
	}
	var toClear []rowT
	for rows.Next() {
		var r rowT
		if rows.Scan(&r.id, &r.layout) == nil && r.layout != "" {
			// Only write when the disk file does not already exist.
			if _, statErr := os.Stat(filepath.Join(printSeedDir, fmt.Sprintf("%s.json", r.id))); os.IsNotExist(statErr) {
				writePrintLayoutFile(r.id, r.layout)
			}
			toClear = append(toClear, r)
		}
	}
	for _, r := range toClear {
		Pool.Exec(ctx, `UPDATE print_templates SET layout = '' WHERE id = $1`, r.id)
	}
}

func Seed() {
	if Pool == nil {
		return
	}
	ctx := context.Background()
	if err := Pool.Ping(ctx); err != nil {
		return
	}

	var empCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM hr_employees").Scan(&empCount)
	if empCount == 0 {
		fmt.Println("🌱 Seeding HR employees, contracts, attendance, leaves, timesheets & payslips...")

		employees := []struct {
			ID, Code, Name, Avatar, Gender, DOB, Phone, Email, IDCard, Address, Dept, Pos, Manager, JoinDate, ContractType, Loc, Status string
			Certificates                                                                                                                string
			BaseSalary, ContractedHours                                                                                                 float64
		}{
			{"EMP-01", "NV-001", "Nguyễn Đức Trường", "T", "Nam", "15/04/1982", "0912.345.678", "truongnd@ttcgroup.vn", "025082001234", "Thị xã Phú Thọ, Tỉnh Phú Thọ", "Ban Giám Đốc", "Giám Đốc Điều Hành Mỏ", "HĐQT Tập Đoàn TTC", "15/03/2018", "Chính thức vô thời hạn", "Văn phòng Điều Hành Mỏ Phú Thọ", "Đang làm việc", `["Chứng chỉ Chỉ huy Nổ mìn Hạng 1", "Chứng chỉ Giám đốc Điều hành Mỏ"]`, 38000000, 8.0},
			{"EMP-02", "NV-002", "Nguyễn Văn Dũng", "D", "Nam", "22/08/1988", "0984.556.789", "dungnv@ttcgroup.vn", "025088005678", "Huyện Thanh Ba, Tỉnh Phú Thọ", "Tổ Vận Hành Trạm Cân", "Trưởng Trạm Cân Điện Tử", "Nguyễn Đức Trường", "01/06/2020", "Chính thức vô thời hạn", "Trạm Cân Cổng 01 - Mỏ Phú Thọ", "Đang làm việc", `["Chứng chỉ Vận hành Cân điện tử & OCR AI", "Chứng nhận An toàn Lao động Nhóm 3"]`, 16500000, 8.0},
			{"EMP-03", "NV-003", "Trần Văn Kiên", "K", "Nam", "10/12/1985", "0988.341.992", "kientv@ttcgroup.vn", "025085009012", "Huyện Phù Ninh, Tỉnh Phú Thọ", "Tổ Khoan Nổ Mìn & VLNCN", "Chỉ Huy Nổ Mìn & An Toàn Mỏ", "Nguyễn Đức Trường", "12/08/2019", "Chính thức vô thời hạn", "Moong Khai Thác Tầng 3 (+45m)", "Đang làm việc", `["Chứng chỉ Chỉ huy Nổ mìn Cục CNQC", "Thẻ An toàn Hóa chất & VLNCN"]`, 28000000, 8.0},
			{"EMP-04", "NV-004", "Lê Hữu Thắng", "T", "Nam", "05/03/1991", "0977.654.321", "thanglh@ttcgroup.vn", "025091003456", "Huyện Đoan Hùng, Tỉnh Phú Thọ", "Xưởng Nghiền Sàng Đá", "Trưởng Ca Dây Chuyền Nghiền Sàng 01", "Nguyễn Đức Trường", "10/01/2021", "Hợp đồng 1 - 3 năm", "Khu Nghiền Sàng Đá 1x2, 4x6", "Đang làm việc", `["Bằng Kỹ thuật Cơ khí Chế tạo máy", "Chứng chỉ An toàn Vận hành Máy nghiền"]`, 18500000, 8.0},
			{"EMP-05", "NV-005", "Trần Đình Trọng", "T", "Nam", "18/07/1987", "0984.112.334", "trongtd@ttcgroup.vn", "025087007890", "Thành phố Việt Trì, Phú Thọ", "Đội Cơ Giới & Vận Tải Mỏ", "Tài Xế Xe Ben 4 Chân (HOWO 88H-042.27)", "Hoàng Minh Đức", "05/05/2022", "Khoán sản lượng", "Hành lang Vận tải QL2 & Mỏ", "Đang làm việc", `["Giấy phép Lái xe Hạng FC", "Chứng nhận Tập huấn Nghiệp vụ Lái xe Mỏ"]`, 14000000, 8.0},
			{"EMP-06", "NV-006", "Nguyễn Văn Mạnh", "M", "Nam", "29/11/1993", "0982.145.882", "manhnv@ttcgroup.vn", "025093001122", "Huyện Lâm Thao, Tỉnh Phú Thọ", "Đội Cơ Giới & Vận Tải Mỏ", "Tài Xế Xe Ben Chenglong (19H-056.22)", "Hoàng Minh Đức", "15/09/2022", "Khoán sản lượng", "Tuyến Mỏ TTC ➔ Cao Tốc Phú Thọ", "Đang làm việc", `["Giấy phép Lái xe Hạng E", "Chứng chỉ An toàn Cơ giới"]`, 14000000, 8.0},
			{"EMP-07", "NV-007", "Lê Văn Cường", "C", "Nam", "08/04/1989", "0977.890.123", "cuonglv@ttcgroup.vn", "025089004455", "Huyện Thanh Ba, Phú Thọ", "Đội Cơ Giới & Vận Tải Mỏ", "Thợ Vận Hành Máy Xúc Komatsu PC450", "Nguyễn Đức Trường", "01/03/2021", "Chính thức vô thời hạn", "Moong Khai Thác Tầng 3", "Đang làm việc", `["Chứng chỉ Thợ lái Máy xúc Bánh xích Bậc 4/7", "Chứng chỉ ATLĐ Mỏ lộ thiên"]`, 19000000, 8.0},
			{"EMP-08", "NV-008", "Nguyễn Thị Thủy", "T", "Nữ", "14/09/1992", "0915.678.901", "thuynt@ttcgroup.vn", "025092006677", "Thành phố Việt Trì, Tỉnh Phú Thọ", "Phòng Kế Toán & Vật Tư", "Kế Toán Trưởng & Đối Soát Doanh Thu", "Nguyễn Đức Trường", "01/10/2019", "Chính thức vô thời hạn", "Văn phòng Kế toán Mỏ", "Đang làm việc", `["Chứng chỉ Kế toán Trưởng Bộ Tài Chính", "Chứng chỉ Hóa Đơn Điện Tử"]`, 24000000, 8.0},
			{"EMP-09", "NV-009", "Phạm Hoàng Nam", "N", "Nam", "20/06/1990", "0983.221.445", "namph@ttcgroup.vn", "025090008899", "Huyện Tam Nông, Tỉnh Phú Thọ", "Phòng Kỹ Thuật & An Toàn", "Kỹ Sư Trắc Địa & Đo Đạc Trữ Lượng Mỏ", "Nguyễn Đức Trường", "15/04/2021", "Hợp đồng 1 - 3 năm", "Khai trường Toàn mỏ TTC", "Đang làm việc", `["Kỹ sư Mỏ Địa chất ĐH Mỏ Địa Chất", "Chứng chỉ Bay Flycam 3D Khảo sát trữ lượng"]`, 21000000, 8.0},
			{"EMP-10", "NV-010", "Hoàng Minh Đức", "D", "Nam", "02/02/1986", "0912.876.432", "duchm@ttcgroup.vn", "025086003344", "Huyện Phù Ninh, Tỉnh Phú Thọ", "Tổ Vận Hành Trạm Cân", "Điều Phối Bãi Tập Kết & Trạm Cân Cổng 2", "Nguyễn Đức Trường", "01/08/2020", "Chính thức vô thời hạn", "Bãi Xe & Trạm Cân Cổng 02", "Đang làm việc", `["Chứng chỉ Điều độ Giao thông Mỏ", "Thẻ An toàn Trạm Cân"]`, 16500000, 8.0},
		}

		for _, e := range employees {
			Pool.Exec(ctx, `
				INSERT INTO hr_employees (id, code, name, avatar, gender, dob, phone, email, id_card, address, department, job_position, manager, join_date, contract_type, contracted_hours, work_location, status, certificates, base_salary)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
				ON CONFLICT (id) DO NOTHING
			`, e.ID, e.Code, e.Name, e.Avatar, e.Gender, e.DOB, e.Phone, e.Email, e.IDCard, e.Address, e.Dept, e.Pos, e.Manager, e.JoinDate, e.ContractType, e.ContractedHours, e.Loc, e.Status, e.Certificates, e.BaseSalary)
		}

		// HR & Enterprise Business/Mining Contracts
		Pool.Exec(ctx, `
			INSERT INTO hr_contracts (id, code, employee_id, employee_name, department, contract_type, contracted_hours, start_date, end_date, base_salary, hazard_allowance, safety_allowance, meal_allowance, responsibility_allowance, attendance_allowance, social_insurance_base, status)
			VALUES 
			('CTR-KD-01', 'HĐMB-2026-ĐC', 'PARTNER-01', 'Tập đoàn Đèo Cả (Dự án Cao tốc Bắc - Nam)', 'Khối Khách hàng & Kinh doanh', 'Hợp đồng mua bán khoáng sản', 8.0, '01/10/2026', '31/12/2027', 33000000000, 0, 0, 0, 0, 0, 0, 'Hiệu lực'),
			('CTR-KT-02', 'HĐBX-2026-TK', 'PARTNER-02', 'Đội xe Cơ giới & Vận tải Tân Tiến', 'Ban Khai trường & Bốc xúc', 'Hợp đồng bốc xúc moong mỏ', 8.0, '15/09/2026', '15/09/2027', 14500000000, 0, 0, 0, 0, 0, 0, 'Hiệu lực'),
			('CTR-TB-03', 'HĐTB-2026-FL', 'PARTNER-03', 'Công ty CP Thiết bị Vinamac', 'Phòng Cơ điện & Thiết bị mỏ', 'Hợp đồng thuê dây chuyền nghiền', 8.0, '01/10/2026', '01/10/2027', 5400000000, 0, 0, 0, 0, 0, 0, 'Hiệu lực'),
			('CTR-NM-04', 'BBNT-2026-NM08', 'PARTNER-04', 'Công ty Hóa chất Mỏ MIC Miền Bắc', 'Tổ Kỹ thuật Khoan nổ mìn', 'Biên bản nghiệm thu nổ mìn đợt 8', 8.0, '20/10/2026', '28/10/2026', 2850000000, 0, 0, 0, 0, 0, 0, 'Hiệu lực'),
			('CTR-VT-05', 'HĐNL-2026-PLX', 'PARTNER-05', 'Chi nhánh Xăng dầu Petrolimex Phú Thọ', 'Phòng Cung ứng Nhiên liệu', 'Hợp đồng cung cấp dầu DO 0.05S', 8.0, '01/11/2026', '01/11/2027', 8200000000, 0, 0, 0, 0, 0, 0, 'Hiệu lực'),
			('CTR-KD-06', 'HĐMB-2026-BQP319', 'PARTNER-06', 'Tổng Công ty 319 - Bộ Quốc Phòng', 'Khối Khách hàng Dự án Quốc gia', 'Hợp đồng mua bán đá Hộc & Base', 8.0, '01/08/2026', '01/08/2027', 18800000000, 0, 0, 0, 0, 0, 0, 'Hiệu lực'),
			('CTR-LD-07', 'HĐLĐ-2026-KS', 'EMP-01', 'Nguyễn Đức Trường', 'Ban Giám Đốc Điều Hành', 'Chính thức vô thời hạn', 8.0, '15/03/2018', 'Vô thời hạn', 38000000, 1500000, 1500000, 1000000, 2000000, 1000000, 38000000, 'Hiệu lực'),
			('CTR-01', 'HĐLĐ-2018-01', 'EMP-01', 'Nguyễn Đức Trường', 'Ban Giám Đốc', 'Chính thức vô thời hạn', 8.0, '15/03/2018', 'Vô thời hạn', 38000000, 1500000, 1500000, 1000000, 2000000, 1000000, 38000000, 'Hiệu lực'),
			('CTR-02', 'HĐLĐ-2020-02', 'EMP-02', 'Nguyễn Văn Dũng', 'Tổ Vận Hành Trạm Cân', 'Chính thức vô thời hạn', 8.0, '01/06/2020', 'Vô thời hạn', 16500000, 800000, 600000, 800000, 1000000, 800000, 16500000, 'Hiệu lực')
			ON CONFLICT (id) DO NOTHING
		`)

		// HR Attendances (Hours based)
		attendances := []struct {
			ID, EmpName, Pos, Dept, EqType, EqCode, Plate, Date, CheckIn, CheckOut, Method, Loc, Area, Status, Notes string
			ContractedHours, WorkedHours, OTHours                                                                    float64
		}{
			{"ATT-2810-01", "Nguyễn Đức Trường", "Giám Đốc Điều Hành Mỏ", "Ban Giám Đốc", "", "", "", "28/10/2026", "07:25", "17:05", "Camera AI ANPR Cổng 1", "Văn phòng Điều Hành Mỏ", "Toàn khai trường", "Đủ giờ công", "Kiểm tra công trường moong tầng 3 (Vượt định mức 0.5h)", 8.0, 8.5, 0.5},
			{"ATT-2810-02", "Nguyễn Văn Dũng", "Trưởng Trạm Cân Điện Tử", "Tổ Vận Hành Trạm Cân", "Trạm cân 100T", "LOADCELL-01", "", "28/10/2026", "06:30", "15:00", "Vân tay Trạm cân", "Trạm Cân Cổng 01", "Cổng 1", "Đủ giờ công", "Chốt 48 phiếu cân xe xuất bãi", 8.0, 8.5, 0.5},
			{"ATT-2810-03", "Trần Văn Kiên", "Chỉ Huy Nổ Mìn & An Toàn Mỏ", "Tổ Khoan Nổ Mìn & VLNCN", "Dàn khoan khí nén", "DRILL-D45", "", "28/10/2026", "06:00", "15:30", "GPS Moong Mỏ", "Moong Khai Thác Tầng 3", "Tầng 3 (+45m)", "Tăng ca (OT)", "Chỉ huy vụ nổ mìn và giám sát an toàn (Tăng ca 1.5h)", 8.0, 9.5, 1.5},
			{"ATT-2810-04", "Lê Văn Cường", "Thợ Vận Hành Máy Xúc", "Đội Cơ Giới & Vận Tải Mỏ", "Máy xúc Komatsu", "EXCAV-PC450", "", "28/10/2026", "06:10", "16:10", "Cảm biến Telematics Komatsu", "Moong Tầng 3", "Tầng 3", "Tăng ca (OT)", "Bốc đá hộc xô bồ đạt 1.450 m3 (OT 2.0h)", 8.0, 10.0, 2.0},
			{"ATT-2810-05", "Trần Đình Trọng", "Tài Xế Xe Ben 4 Chân", "Đội Cơ Giới & Vận Tải Mỏ", "Xe ben HOWO", "HOWO-371", "88H-042.27", "28/10/2026", "07:00", "15:00", "Định vị GPS Hộp Đen", "Quốc Lộ 2 Km 68", "Tuyến QL2", "Đủ giờ công", "Đã hoàn thành 5 chuyến đá 1x2", 8.0, 8.0, 0.0},
			{"ATT-2810-06", "Nguyễn Văn Mạnh", "Tài Xế Xe Ben Chenglong", "Đội Cơ Giới & Vận Tải Mỏ", "Xe ben Chenglong", "CHENG-385", "19H-056.22", "28/10/2026", "08:00", "16:00", "Camera AI Cổng Trạm Cân", "Đường Nội Bộ Mỏ", "Mỏ ➔ Cảng Sông Lô", "Đủ giờ công", "Đã hoàn thành đủ 8.0h định mức HĐ", 8.0, 8.0, 0.0},
			{"ATT-2810-07", "Lê Hữu Thắng", "Trưởng Ca Máy Nghiền 01", "Xưởng Nghiền Sàng Đá", "Máy nghiền kẹp hàm", "CRUSH-01", "", "28/10/2026", "07:00", "13:30", "Vân tay Xưởng Nghiền", "Xưởng Nghiền Sàng 01", "Khu nghiền sàng", "Thiếu giờ", "Đã làm 6.5h / 8.0h định mức (Đang làm tiếp buổi chiều)", 8.0, 6.5, 0.0},
		}

		for _, a := range attendances {
			Pool.Exec(ctx, `
				INSERT INTO hr_attendances (id, employee_name, job_position, department, equipment_type, equipment_code, vehicle_plate, date, check_in_time, check_out_time, contracted_hours, worked_hours, ot_hours, method, location, work_area, status, notes)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
				ON CONFLICT (id) DO NOTHING
			`, a.ID, a.EmpName, a.Pos, a.Dept, a.EqType, a.EqCode, a.Plate, a.Date, a.CheckIn, a.CheckOut, a.ContractedHours, a.WorkedHours, a.OTHours, a.Method, a.Loc, a.Area, a.Status, a.Notes)
		}

		// HR Leaves
		Pool.Exec(ctx, `
			INSERT INTO hr_leave_requests (id, code, employee_name, department, leave_type, start_date, end_date, total_days, total_hours, reason, stage, approver)
			VALUES 
			('LEAVE-01', 'ĐNP-281025-01', 'Trần Đình Trọng', 'Đội Cơ Giới & Vận Tải Mỏ', 'Phép năm', '29/10/2026', '30/10/2026', 2.0, 16.0, 'Giải quyết công việc gia đình tại Vĩnh Phúc', 'Chờ GĐ Mỏ duyệt', 'Nguyễn Đức Trường'),
			('LEAVE-02', 'ĐNP-281025-02', 'Phạm Hoàng Nam', 'Phòng Kỹ Thuật & An Toàn', 'Công tác / Đào tạo', '25/10/2026', '26/10/2026', 2.0, 16.0, 'Tham dự Khóa huấn luyện An toàn VLNCN do Sở Công Thương tổ chức', 'Đã duyệt', 'Nguyễn Đức Trường'),
			('LEAVE-03', 'ĐNP-281025-03', 'Lê Hữu Thắng', 'Xưởng Nghiền Sàng Đá', 'Nghỉ ốm', '12/10/2026', '13/10/2026', 2.0, 16.0, 'Nghỉ ốm có giấy xác nhận của Bệnh viện Đa khoa Thanh Ba', 'Đã duyệt', 'Nguyễn Đức Trường')
			ON CONFLICT (id) DO NOTHING
		`)

		// HR Timesheets JSON
		timesheetsJSON := `[
			{"employeeId":"EMP-01","employeeCode":"NV-001","employeeName":"Nguyễn Đức Trường","department":"Ban Giám Đốc","jobPosition":"Giám Đốc Điều Hành Mỏ","contractedHoursPerDay":8.0,"standardDays":26,"actualDays":26,"totalContractHours":208.0,"totalWorkedHours":220.0,"paidLeaveDays":0,"unpaidLeaveDays":0,"otHours":12.0,"pieceworkTons":0,"qualifiedDays":26},
			{"employeeId":"EMP-02","employeeCode":"NV-002","employeeName":"Nguyễn Văn Dũng","department":"Tổ Vận Hành Trạm Cân","jobPosition":"Trưởng Trạm Cân Điện Tử","contractedHoursPerDay":8.0,"standardDays":26,"actualDays":25,"totalContractHours":208.0,"totalWorkedHours":218.0,"paidLeaveDays":1,"unpaidLeaveDays":0,"otHours":18.0,"pieceworkTons":2500,"qualifiedDays":25},
			{"employeeId":"EMP-03","employeeCode":"NV-003","employeeName":"Trần Văn Kiên","department":"Tổ Khoan Nổ Mìn & VLNCN","jobPosition":"Chỉ Huy Nổ Mìn & An Toàn Mỏ","contractedHoursPerDay":8.0,"standardDays":26,"actualDays":26,"totalContractHours":208.0,"totalWorkedHours":232.0,"paidLeaveDays":0,"unpaidLeaveDays":0,"otHours":24.0,"pieceworkTons":14500,"qualifiedDays":26},
			{"employeeId":"EMP-04","employeeCode":"NV-004","employeeName":"Lê Hữu Thắng","department":"Xưởng Nghiền Sàng Đá","jobPosition":"Trưởng Ca Máy Nghiền 01","contractedHoursPerDay":8.0,"standardDays":26,"actualDays":24,"totalContractHours":208.0,"totalWorkedHours":208.0,"paidLeaveDays":2,"unpaidLeaveDays":0,"otHours":16.0,"pieceworkTons":6800,"qualifiedDays":24},
			{"employeeId":"EMP-05","employeeCode":"NV-005","employeeName":"Trần Đình Trọng","department":"Đội Cơ Giới & Vận Tải Mỏ","jobPosition":"Tài Xế Xe Ben HOWO","contractedHoursPerDay":8.0,"standardDays":26,"actualDays":26,"totalContractHours":208.0,"totalWorkedHours":240.0,"paidLeaveDays":0,"unpaidLeaveDays":0,"otHours":32.0,"pieceworkTons":8200,"qualifiedDays":26},
			{"employeeId":"EMP-07","employeeCode":"NV-007","employeeName":"Lê Văn Cường","department":"Đội Cơ Giới & Vận Tải Mỏ","jobPosition":"Thợ Vận Hành Máy Xúc PC450","contractedHoursPerDay":8.0,"standardDays":26,"actualDays":26,"totalContractHours":208.0,"totalWorkedHours":236.0,"paidLeaveDays":0,"unpaidLeaveDays":0,"otHours":28.0,"pieceworkTons":12500,"qualifiedDays":26},
			{"employeeId":"EMP-08","employeeCode":"NV-008","employeeName":"Nguyễn Thị Thủy","department":"Phòng Kế Toán & Vật Tư","jobPosition":"Kế Toán Trưởng Mỏ","contractedHoursPerDay":8.0,"standardDays":26,"actualDays":26,"totalContractHours":208.0,"totalWorkedHours":216.0,"paidLeaveDays":0,"unpaidLeaveDays":0,"otHours":8.0,"pieceworkTons":0,"qualifiedDays":26}
		]`
		Pool.Exec(ctx, `INSERT INTO hr_timesheets (id, data) VALUES ('TS-102026', $1::jsonb) ON CONFLICT (id) DO NOTHING`, timesheetsJSON)

		// HR Payslips JSON
		payslipsJSON := `[
			{"id":"PS-1026-01","code":"PL-1026-01","employeeId":"EMP-01","employeeName":"Nguyễn Đức Trường","department":"Ban Giám Đốc","jobPosition":"Giám Đốc Điều Hành Mỏ","month":"10/2026","contractedHours":208.0,"workedHours":220.0,"actualDays":26,"baseSalary":38000000,"workingSalary":38000000,"pieceworkSalary":0,"allowances":6000000,"overtimeSalary":3500000,"bonus":5000000,"grossSalary":52500000,"totalDeductions":5200000,"netSalary":47300000,"status":"Đã thanh toán ngân hàng"},
			{"id":"PS-1026-02","code":"PL-1026-02","employeeId":"EMP-02","employeeName":"Nguyễn Văn Dũng","department":"Tổ Vận Hành Trạm Cân","jobPosition":"Trưởng Trạm Cân Điện Tử","month":"10/2026","contractedHours":208.0,"workedHours":218.0,"actualDays":25,"baseSalary":16500000,"workingSalary":15865000,"pieceworkSalary":2500000,"allowances":2200000,"overtimeSalary":2100000,"bonus":1500000,"grossSalary":24165000,"totalDeductions":2150000,"netSalary":22015000,"status":"Đã thanh toán ngân hàng"},
			{"id":"PS-1026-03","code":"PL-1026-03","employeeId":"EMP-03","employeeName":"Trần Văn Kiên","department":"Tổ Khoan Nổ Mìn & VLNCN","jobPosition":"Chỉ Huy Nổ Mìn & An Toàn Mỏ","month":"10/2026","contractedHours":208.0,"workedHours":232.0,"actualDays":26,"baseSalary":28000000,"workingSalary":28000000,"pieceworkSalary":8500000,"allowances":4500000,"overtimeSalary":4200000,"bonus":3000000,"grossSalary":48200000,"totalDeductions":4600000,"netSalary":43600000,"status":"Đã thanh toán ngân hàng"},
			{"id":"PS-1026-04","code":"PL-1026-04","employeeId":"EMP-04","employeeName":"Lê Hữu Thắng","department":"Xưởng Nghiền Sàng Đá","jobPosition":"Trưởng Ca Dây Chuyền Nghiền Sàng 01","month":"10/2026","contractedHours":208.0,"workedHours":208.0,"actualDays":24,"baseSalary":18500000,"workingSalary":17076000,"pieceworkSalary":4800000,"allowances":2800000,"overtimeSalary":2400000,"bonus":2000000,"grossSalary":29076000,"totalDeductions":2650000,"netSalary":26426000,"status":"Đã xác nhận"},
			{"id":"PS-1026-05","code":"PL-1026-05","employeeId":"EMP-05","employeeName":"Trần Đình Trọng","department":"Đội Cơ Giới & Vận Tải Mỏ","jobPosition":"Tài Xế Xe Ben 4 Chân (HOWO 88H-042.27)","month":"10/2026","contractedHours":208.0,"workedHours":240.0,"actualDays":26,"baseSalary":14000000,"workingSalary":14000000,"pieceworkSalary":11200000,"allowances":2500000,"overtimeSalary":3800000,"bonus":1500000,"grossSalary":33000000,"totalDeductions":2950000,"netSalary":30050000,"status":"Đã thanh toán ngân hàng"},
			{"id":"PS-1026-06","code":"PL-1026-06","employeeId":"EMP-06","employeeName":"Nguyễn Văn Mạnh","department":"Đội Cơ Giới & Vận Tải Mỏ","jobPosition":"Tài Xế Xe Ben Chenglong (19H-056.22)","month":"10/2026","contractedHours":208.0,"workedHours":235.0,"actualDays":26,"baseSalary":14000000,"workingSalary":14000000,"pieceworkSalary":10800000,"allowances":2500000,"overtimeSalary":3600000,"bonus":1500000,"grossSalary":32400000,"totalDeductions":2900000,"netSalary":29500000,"status":"Đã thanh toán ngân hàng"},
			{"id":"PS-1026-07","code":"PL-1026-07","employeeId":"EMP-07","employeeName":"Lê Văn Cường","department":"Đội Cơ Giới & Vận Tải Mỏ","jobPosition":"Thợ Vận Hành Máy Xúc Komatsu PC450","month":"10/2026","contractedHours":208.0,"workedHours":236.0,"actualDays":26,"baseSalary":19000000,"workingSalary":19000000,"pieceworkSalary":9500000,"allowances":3200000,"overtimeSalary":4100000,"bonus":2000000,"grossSalary":37800000,"totalDeductions":3450000,"netSalary":34350000,"status":"Đã xác nhận"},
			{"id":"PS-1026-08","code":"PL-1026-08","employeeId":"EMP-08","employeeName":"Nguyễn Thị Thủy","department":"Phòng Kế Toán & Vật Tư","jobPosition":"Kế Toán Trưởng Mỏ","month":"10/2026","contractedHours":208.0,"workedHours":216.0,"actualDays":26,"baseSalary":24000000,"workingSalary":24000000,"pieceworkSalary":0,"allowances":3500000,"overtimeSalary":1200000,"bonus":2500000,"grossSalary":31200000,"totalDeductions":3100000,"netSalary":28100000,"status":"Đã thanh toán ngân hàng"}
		]`
		Pool.Exec(ctx, `INSERT INTO hr_payslips (id, data) VALUES ('PS-102026', $1::jsonb) ON CONFLICT (id) DO NOTHING`, payslipsJSON)
	}

	// 2. Seed Tickets (Phiếu Cân)
	var ticketCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM tickets").Scan(&ticketCount)
	if ticketCount == 0 {
		fmt.Println("🌱 Seeding Weighing Tickets...")
		tickets := []struct {
			ID, BenBan, BenMua, BienSo, LoaiXe, LaiXe, Sdt, Loai, Stage, StageLabel, KL1, KL2, KLHang, KLTinhTien, ThanhTien, Time1, Time2, Date, NguoiCan1, MatHang, QuyCach, DOCode, TramCan, CongCan, GhiChu string
			CanL1, DonGia                                                                                                                                                                                       float64
		}{
			{"TK-20261028-001", "Công ty Cổ phần Mỏ Đá TTC", "Công ty CP Đầu Tư Xây Dựng 319", "88H-042.27", "Xe ben 4 chân", "Trần Đình Trọng", "0984.112.334", "Cân bán hàng", "done", "Hoàn thành", "15.42", "48.60", "33.18", "33.18", "7.963.200", "08:15", "08:42", "28/10/2026", "Nguyễn Văn Dũng", "Đá 1x2 Xây Dựng", "TCVN 7570:2006", "DO-319-082", "Trạm Cân 01 (100 Tấn)", "Cổng Cân Số 1", "Đủ điều kiện xuất mỏ", 15.42, 240000},
			{"TK-20261028-002", "Công ty Cổ phần Mỏ Đá TTC", "Tổng Công Ty XD Trường Sơn", "19H-056.22", "Xe ben Chenglong", "Nguyễn Văn Mạnh", "0982.145.882", "Cân bán hàng", "done", "Hoàn thành", "14.80", "46.20", "31.40", "31.40", "5.652.000", "08:30", "08:58", "28/10/2026", "Nguyễn Văn Dũng", "Đá Base Cấp Phối Dmax25", "TCVN 8859:2011", "DO-TS-045", "Trạm Cân 01 (100 Tấn)", "Cổng Cân Số 1", "Đá cấp phối chuẩn cao tốc", 14.80, 180000},
			{"TK-20261028-003", "Công ty Cổ phần Mỏ Đá TTC", "Công ty Bê Tông Việt Trì", "29H-882.19", "Xe bồn trộn 12m3", "Lê Văn Tuấn", "0912.445.667", "Cân bán hàng", "done", "Hoàn thành", "16.10", "42.50", "26.40", "26.40", "6.864.000", "09:05", "09:32", "28/10/2026", "Nguyễn Văn Dũng", "Cát Nghiền Nhân Tạo", "Mô đun độ lớn 2.6", "DO-VT-012", "Trạm Cân 01 (100 Tấn)", "Cổng Cân Số 2", "Cát nghiền mịn trạm trộn bê tông", 16.10, 260000},
			{"TK-20261028-004", "Công ty Cổ phần Mỏ Đá TTC", "Công ty CP Tập Đoàn Đèo Cả", "90C-123.45", "Xe đầu kéo Howo", "Vũ Quốc Đạt", "0978.556.223", "Cân bán hàng", "weighing_2", "Đang cân bì lần 2", "18.50", "—", "—", "—", "—", "09:40", "—", "28/10/2026", "Nguyễn Văn Dũng", "Đá 4x6 Kè Móng", "Đá 4x6 tuyển", "DO-DC-098", "Trạm Cân 02 (120 Tấn)", "Cổng Cân Số 1", "Đang bốc hàng moong tầng 3", 18.50, 220000},
		}

		for _, t := range tickets {
			Pool.Exec(ctx, `
				INSERT INTO tickets (id, ben_ban, ben_mua, bien_so, loai_xe, lai_xe, sdt_lai_xe, loai, stage, stage_label, can_l1, kl1, can_l2, kl2, kl_hang, kl_tinh_tien, don_gia, thanh_tien, time1, time2, date, nguoi_can1, mat_hang, quy_cach, do_code, tram_can, cong_can, ghi_chu)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)
				ON CONFLICT (id) DO NOTHING
			`, t.ID, t.BenBan, t.BenMua, t.BienSo, t.LoaiXe, t.LaiXe, t.Sdt, t.Loai, t.Stage, t.StageLabel, t.CanL1, t.KL1, t.KL2, t.KL2, t.KLHang, t.KLTinhTien, t.DonGia, t.ThanhTien, t.Time1, t.Time2, t.Date, t.NguoiCan1, t.MatHang, t.QuyCach, t.DOCode, t.TramCan, t.CongCan, t.GhiChu)
		}
	}

	// 3. Seed Vehicles
	// 3. Seed Vehicles
	fmt.Println("🌱 Seeding Fleet Vehicles...")
	vehicles := []struct {
		BS, Loai, RFID, Status, HanDangKiem, Date, ChuXe, Unit string
		Bi, TaiTrong                                          float64
		Count                                                 int
	}{
		{"88H-042.27", "Xe ben 4 chân HOWO", "RFID-88H-042", "Hoạt động", "15/12/2026", "28/10/2026", "Công ty CP Đầu Tư Xây Dựng 319", "tấn", 15.42, 30.0, 142},
		{"19H-056.22", "Xe ben Chenglong", "RFID-19H-056", "Hoạt động", "20/01/2027", "28/10/2026", "Tổng Công Ty XD Trường Sơn", "tấn", 14.80, 30.0, 98},
		{"29H-882.19", "Xe bồn trộn bê tông", "RFID-29H-882", "Hoạt động", "10/11/2026", "28/10/2026", "Công ty Bê Tông Việt Trì", "tấn", 16.10, 25.0, 76},
		{"90C-123.45", "Xe đầu kéo Mooc ben", "RFID-90C-123", "Hoạt động", "05/03/2027", "28/10/2026", "Công ty CP Tập Đoàn Đèo Cả", "tấn", 18.50, 45.0, 64},
		{"19C-098.76", "Xe tải ben 3 chân", "RFID-19C-098", "Bảo dưỡng", "28/10/2026", "28/10/2026", "Hợp tác xã Vận tải Hùng Vương", "tấn", 11.20, 20.0, 35},
	}

	for _, v := range vehicles {
		Pool.Exec(ctx, `
			INSERT INTO vehicles (bs, loai, bi, rfid, tai_trong, unit, status, count, han_dang_kiem, date, chu_xe)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT (bs) DO UPDATE SET bi = EXCLUDED.bi, tai_trong = EXCLUDED.tai_trong, unit = EXCLUDED.unit, loai = EXCLUDED.loai, status = EXCLUDED.status, count = EXCLUDED.count, han_dang_kiem = EXCLUDED.han_dang_kiem, chu_xe = EXCLUDED.chu_xe
		`, v.BS, v.Loai, v.Bi, v.RFID, v.TaiTrong, v.Unit, v.Status, v.Count, v.HanDangKiem, v.Date, v.ChuXe)
	}

	// 4. Seed Materials Catalog
	fmt.Println("🌱 Seeding Stone Materials Catalog...")
	materials := []struct {
		Code, Name, DVT, Unit, Kho, Standard, Status, Date string
		Density, DinhMuc, Price                             float64
	}{
		{"MAT-01", "Đá 1x2 Xây Dựng", "tấn", "tấn", "Bãi Đá Thành Phẩm 01", "TCVN 7570:2006", "Đang kinh doanh", "28/10/2026", 1.55, 1.55, 240000},
		{"MAT-02", "Đá 4x6 Kè Móng", "tấn", "tấn", "Bãi Đá Hộc & 4x6", "TCVN 1771:1987", "Đang kinh doanh", "28/10/2026", 1.60, 1.60, 220000},
		{"MAT-03", "Đá Base Cấp Phối Loại 1", "tấn", "tấn", "Bãi Cấp Phối Dmax25", "TCVN 8859:2011", "Đang kinh doanh", "28/10/2026", 1.00, 1.00, 180000},
		{"MAT-04", "Cát Nghiền Nhân Tạo (Mạt đá)", "tấn", "tấn", "Kho Mái Che Cát Nghiền", "TCVN 9205:2012", "Đang kinh doanh", "28/10/2026", 1.45, 1.45, 260000},
		{"MAT-05", "Đá Hộc Xô Bồ Nổ Mìn", "tấn", "tấn", "Moong Tầng 3 (+45m)", "Đá khai thác nguyên khai", "Đang kinh doanh", "28/10/2026", 1.65, 1.65, 160000},
		{"MAT-06", "Đá 2x4 Đổ Bê Tông", "tấn", "tấn", "Bãi Đá Thành Phẩm 02", "TCVN 7570:2006", "Đang kinh doanh", "28/10/2026", 1.58, 1.58, 235000},
	}

	for _, m := range materials {
		Pool.Exec(ctx, `
			INSERT INTO catalog_materials (code, ma, name, ten, dvt, unit, density, ty_trong, dinh_muc, price, gia, kho, standard, status, date)
			VALUES ($1, $1, $2, $2, $3, $4, $5, $5, $6, $7, $7, $8, $9, $10, $11)
			ON CONFLICT (code) DO UPDATE SET dvt = EXCLUDED.dvt, unit = EXCLUDED.unit, density = EXCLUDED.density, ty_trong = EXCLUDED.ty_trong, dinh_muc = EXCLUDED.dinh_muc, price = EXCLUDED.price, gia = EXCLUDED.gia, name = EXCLUDED.name, ten = EXCLUDED.ten, kho = EXCLUDED.kho, standard = EXCLUDED.standard
		`, m.Code, m.Name, m.DVT, m.Unit, m.Density, m.DinhMuc, m.Price, m.Kho, m.Standard, m.Status, m.Date)
	}

	// 5. Seed Customers & Suppliers
	var custCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM customers").Scan(&custCount)
	if custCount == 0 {
		fmt.Println("🌱 Seeding Customers & Partners...")
		Pool.Exec(ctx, `
			INSERT INTO customers (name, ten, tax_code, mst, contact, phone, sdt, group_name, doanh_so, cong_no, han_muc, status, date)
			VALUES 
			('Công ty CP Đầu Tư Xây Dựng 319', 'Công ty CP Đầu Tư Xây Dựng 319', '0100109988', '0100109988', 'Nguyễn Đức Hưng', '0913.224.556', '0913.224.556', 'Khách hàng VIP Doanh nghiệp', '4.850.000.000 đ', '320.000.000 đ', '1.000.000.000 đ', 'Hoạt động', '28/10/2026'),
			('Tổng Công Ty XD Trường Sơn', 'Tổng Công Ty XD Trường Sơn', '0100108877', '0100108877', 'Phạm Văn Thành', '0982.667.889', '0982.667.889', 'Nhà thầu Quốc gia', '8.200.000.000 đ', '540.000.000 đ', '2.000.000.000 đ', 'Hoạt động', '28/10/2026'),
			('Công ty Bê Tông Việt Trì', 'Công ty Bê Tông Việt Trì', '2600456789', '2600456789', 'Lê Hữu Đạt', '0915.889.001', '0915.889.001', 'Trạm trộn bê tông thương phẩm', '3.100.000.000 đ', '180.000.000 đ', '800.000.000 đ', 'Hoạt động', '28/10/2026')
			ON CONFLICT DO NOTHING;

			INSERT INTO suppliers (name, ten, tax_code, mst, phone, sdt, items, hang, no, rating, status, date)
			VALUES
			('Tổng Công Ty Kinh Tế Kỹ Thuật CNQP (GAET)', 'Tổng Công Ty GAET', '0100109911', '0100109911', '0243.856.789', '0243.856.789', 'Thuốc nổ Anfo, Kíp nổ vi sai & Dịch vụ nổ mìn', 'Thuốc nổ & Kíp nổ', '450.000.000 đ', '5 sao (Đối tác chiến lược)', 'Hoạt động', '28/10/2026'),
			('Công Ty Xăng Dầu Phú Thọ (Petrolimex)', 'Petrolimex Phú Thọ', '2600108822', '2600108822', '0210.384.666', '0210.384.666', 'Dầu Diesel DO 0.05S-II cho xe máy mỏ', 'Dầu DO 0.05S', '280.000.000 đ', '5 sao', 'Hoạt động', '28/10/2026'),
			('Công Ty TNHH Thiết Bị Nặng Marubeni (Komatsu)', 'Komatsu Việt Nam', '0101567890', '0101567890', '0243.765.432', '0243.765.432', 'Phụ tùng máy xúc PC450, gầu đào, răng gầu mỏ', 'Phụ tùng máy xúc', '120.000.000 đ', '5 sao', 'Hoạt động', '28/10/2026')
			ON CONFLICT DO NOTHING;
		`)
	}

	// 6. Seed Weigh Bridges & Cameras
	var bridgeCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM weigh_bridges").Scan(&bridgeCount)
	if bridgeCount == 0 {
		Pool.Exec(ctx, `
			INSERT INTO weigh_bridges (code, name, ten, location, ip, type, capacity, loadcell, status)
			VALUES 
			('WB-01', 'Trạm Cân Điện Tử 01 (100T)', 'Trạm Cân Điện Tử 01 (100T)', 'Cổng số 1 - Lối vào chính', '192.168.1.100', 'Cân bàn thép 18x3m', '100 Tấn', '6x Loadcell Zemic 30T', 'Online Live'),
			('WB-02', 'Trạm Cân Điện Tử 02 (120T)', 'Trạm Cân Điện Tử 02 (120T)', 'Cổng số 2 - Khu bãi tập kết', '192.168.1.101', 'Cân bàn thép 21x3.4m', '120 Tấn', '8x Loadcell Keli 30T', 'Online Live')
			ON CONFLICT (code) DO NOTHING
		`)

		Pool.Exec(ctx, `
			INSERT INTO camera_devices (code, name, ten, location, ip, res, status, type)
			VALUES 
			('CAM-01', 'Camera AI ANPR Cổng 1 (Biển số trước)', 'Camera AI ANPR Cổng 1', 'Trạm Cân 01 - Cổng vào', '192.168.1.201', '4MP 60FPS AI', 'Online 24/7', 'Camera AI Nhận diện biển số'),
			('CAM-02', 'Camera Toàn Cảnh Thùng Xe Cổng 1', 'Camera Toàn Cảnh Thùng Xe Cổng 1', 'Trạm Cân 01 - Cột cao 6m', '192.168.1.202', '4K UHD Zoom', 'Online 24/7', 'Camera Giám sát gian lận tải trọng'),
			('CAM-03', 'Camera AI ANPR Cổng 2', 'Camera AI ANPR Cổng 2', 'Trạm Cân 02 - Cổng ra', '192.168.1.203', '4MP 60FPS AI', 'Online 24/7', 'Camera AI Nhận diện biển số')
			ON CONFLICT (code) DO NOTHING
		`)
	}

	// 7. Seed GPS Fleet Telemetry
	var gpsCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM gps_fleet").Scan(&gpsCount)
	if gpsCount == 0 {
		fmt.Println("🌱 Seeding Live GPS Fleet Telemetry...")
		Pool.Exec(ctx, `
			INSERT INTO gps_fleet (code, data) VALUES
			('DRV-01', '{"id":"DRV-01","plate":"88H-042.27","driverName":"Trần Đình Trọng","driverPhone":"0984.112.334","rfidCode":"RFID-1008-TTC","vehicleType":"Xe ben HOWO 4 chân 371HP","status":"idling_alert","lat":21.3210,"lng":105.3280,"speed":0,"engineAcc":"OFF","engineRpm":0,"currentFuelLiters":182,"fuelTankPercent":45,"locationName":"Quán nước ven QL2 Km 72 (Tắt máy sụt dầu)","destination":"QL2 ➔ Trạm trộn Phù Ninh","cargo":"Đá 1x2 bê tông","updatedTime":"12:55 28/10/2026"}'::jsonb),
			('DRV-02', '{"id":"DRV-02","plate":"19H-056.22","driverName":"Nguyễn Văn Mạnh","driverPhone":"0982.145.882","rfidCode":"RFID-1015-TTC","vehicleType":"Xe ben Chenglong Hải Âu 385HP","status":"moving","lat":21.3240,"lng":105.3180,"speed":48,"engineAcc":"ON","engineRpm":1450,"currentFuelLiters":268,"fuelTankPercent":67,"locationName":"Đang chạy trên Quốc Lộ 2 Km 68","destination":"Mỏ Đá TTC ➔ Cảng Sông Lô","cargo":"Đá Base cấp phối","updatedTime":"13:02 28/10/2026"}'::jsonb),
			('DRV-03', '{"id":"DRV-03","plate":"29C-781.90","driverName":"Lê Văn Cường","driverPhone":"0977.890.123","rfidCode":"RFID-1022-TTC","vehicleType":"Xe ben Shacman 4 chân X3000","status":"loading","lat":21.3255,"lng":105.3005,"speed":0,"engineAcc":"ON","engineRpm":800,"currentFuelLiters":310,"fuelTankPercent":77,"locationName":"Bãi Xe Trung Tâm - Khai Trường Mỏ TTC","destination":"Chờ nhận phiếu cân","cargo":"Đá hộc xô bồ","updatedTime":"13:05 28/10/2026"}'::jsonb)
		`)
	}

	// 8. Seed Fuel Theft Audits
	var theftCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM fuel_theft_audits").Scan(&theftCount)
	if theftCount == 0 {
		fmt.Println("🌱 Seeding Fuel Theft Audits...")
		Pool.Exec(ctx, `
			INSERT INTO fuel_theft_audits (code, data) VALUES
			('AUD-01', '{"id":"AUD-01","code":"AUD-01","plate":"88H-042.27","driverName":"Trần Đình Trọng","vehicleType":"Xe ben HOWO 4 chân 371HP","shiftDate":"28/10/2026","totalKm":124,"totalKmLoaded":82,"totalKmEmpty":42,"idlingHours":2.1,"fuelInitialLiters":320,"fuelRefilledLiters":0,"fuelFinalLiters":182,"actualFuelConsumedLiters":138,"theoreticalNormLiters":52.7,"fuelVarianceLiters":85.3,"variancePercent":61.8,"theftCostLostVnd":1919250,"riskLevel":"CRITICAL_THEFT_DETECTED","anomalyDescription":"Tắt máy đứng yên 42 phút ven QL2 sụt 38L dầu bất thường","resolutionStatus":"Chờ đối soát kỷ luật","locationTheftDetected":"Km 72 QL2 (Quán nước ven đường)","investigator":"Tô Quốc Huy (Ban Thanh Tra Xăng Dầu)"}'::jsonb),
			('AUD-02', '{"id":"AUD-02","code":"AUD-02","plate":"19H-056.22","driverName":"Nguyễn Văn Mạnh","vehicleType":"Xe ben Chenglong Hải Âu 385HP","shiftDate":"28/10/2026","totalKm":156,"totalKmLoaded":110,"totalKmEmpty":46,"idlingHours":0.6,"fuelInitialLiters":312,"fuelRefilledLiters":0,"fuelFinalLiters":268,"actualFuelConsumedLiters":44,"theoreticalNormLiters":46.5,"fuelVarianceLiters":-2.5,"variancePercent":-5.4,"theftCostLostVnd":0,"riskLevel":"SAVING","anomalyDescription":"Vận hành tiết kiệm dầu, đúng lộ trình cao tốc","resolutionStatus":"Bình thường / Khen thưởng","locationTheftDetected":"Tuyến Mỏ TTC ➔ Cao Tốc","investigator":"Tô Quốc Huy"}'::jsonb),
			('AUD-03', '{"id":"AUD-03","code":"AUD-03","plate":"29C-781.90","driverName":"Lê Văn Cường","vehicleType":"Xe ben Shacman 4 chân X3000","shiftDate":"28/10/2026","totalKm":98,"totalKmLoaded":68,"totalKmEmpty":30,"idlingHours":0.8,"fuelInitialLiters":345,"fuelRefilledLiters":0,"fuelFinalLiters":310,"actualFuelConsumedLiters":35,"theoreticalNormLiters":36.2,"fuelVarianceLiters":-1.2,"variancePercent":-3.3,"theftCostLostVnd":0,"riskLevel":"SAVING","anomalyDescription":"Đạt định mức chuẩn nội bộ khai trường","resolutionStatus":"Bình thường","locationTheftDetected":"Khai trường Moong Tầng 3","investigator":"Tô Quốc Huy"}'::jsonb)
		`)
	}

	// 9. Seed Fuel Norms
	var normCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM fuel_norms").Scan(&normCount)
	if normCount == 0 {
		fmt.Println("🌱 Seeding Fleet Fuel Norms...")
		Pool.Exec(ctx, `
			INSERT INTO fuel_norms (code, data) VALUES
			('V-01', '{"id":"V-01","code":"V-01","plate":"88H-042.27","vehicleType":"Xe ben HOWO 4 chân 371HP","type":"Xe ben HOWO 4 chân 371HP","driverName":"Trần Đình Trọng","driverPhone":"0984.112.334","rfidCode":"RFID-1008-TTC","status":"theft_alert","loadedNormLitersPer100km":42.5,"emptyNormLitersPer100km":28.0,"idlingNormLitersPerHour":3.5,"quarryTerrainFactor":1.15,"normLoadedL100km":42.5,"normUnloadedL100km":28.0,"normIdlingLiterPerHour":3.5,"slopeCorrectionFactor":1.15,"currentFuelLiters":182,"tankCapacityLiters":400,"fuelCapacityLiters":400,"odometerCurrentKm":184520,"driverIntegrityScore":58,"lastMaintenanceDate":"15/10/2026","fuelSensorStatus":"ACTIVE_ANOMALY","notes":"Phát hiện sụt dầu bất thường 38 Lít khi dừng tắt máy ven QL2."}'::jsonb),
			('V-02', '{"id":"V-02","code":"V-02","plate":"19H-056.22","vehicleType":"Xe ben Chenglong Hải Âu 385HP","type":"Xe ben Chenglong Hải Âu 385HP","driverName":"Nguyễn Văn Mạnh","driverPhone":"0982.145.882","rfidCode":"RFID-1015-TTC","status":"active","loadedNormLitersPer100km":39.0,"emptyNormLitersPer100km":26.5,"idlingNormLitersPerHour":3.0,"quarryTerrainFactor":1.12,"normLoadedL100km":39.0,"normUnloadedL100km":26.5,"normIdlingLiterPerHour":3.0,"slopeCorrectionFactor":1.12,"currentFuelLiters":268,"tankCapacityLiters":400,"fuelCapacityLiters":400,"odometerCurrentKm":148580,"driverIntegrityScore":96,"lastMaintenanceDate":"20/10/2026","fuelSensorStatus":"ACTIVE_NORMAL","notes":"Đang vận chuyển đá 1x2 cho công trường Cao Tốc."}'::jsonb),
			('V-03', '{"id":"V-03","code":"V-03","plate":"29C-781.90","vehicleType":"Xe ben Shacman 4 chân X3000","type":"Xe ben Shacman 4 chân X3000","driverName":"Lê Văn Cường","driverPhone":"0977.890.123","rfidCode":"RFID-1022-TTC","status":"active","loadedNormLitersPer100km":41.0,"emptyNormLitersPer100km":27.0,"idlingNormLitersPerHour":3.2,"quarryTerrainFactor":1.15,"normLoadedL100km":41.0,"normUnloadedL100km":27.0,"normIdlingLiterPerHour":3.2,"slopeCorrectionFactor":1.15,"currentFuelLiters":310,"tankCapacityLiters":400,"fuelCapacityLiters":400,"odometerCurrentKm":210400,"driverIntegrityScore":94,"lastMaintenanceDate":"12/10/2026","fuelSensorStatus":"ACTIVE_NORMAL","notes":"Chạy tuyến khai trường mỏ về trạm nghiền trung tâm."}'::jsonb)
		`)
	}

	// 10. Seed Yard CheckInOut
	var yardCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM yard_checkinout").Scan(&yardCount)
	if yardCount == 0 {
		fmt.Println("🌱 Seeding Yard Check-in/out Logs...")
		Pool.Exec(ctx, `
			INSERT INTO yard_checkinout (code, data) VALUES
			('CHK-01', '{"id":"CHK-01","code":"CHK-01","yardId":"YARD-PT-01","yardName":"Bãi Xe Trung Tâm - Khai Trường Mỏ Phú Thọ","plate":"88H-042.27","driverName":"Trần Đình Trọng","driverPhone":"0984.112.334","rfidTag":"RFID-1008-TTC","actionType":"CHECK_OUT","timestamp":"07:05 28/10/2026","distanceMeters":12.4,"isWithin20mGeofence":true,"geofenceStatus":"VALID_WITHIN_20M","startOdometer":184396,"endOdometer":184520,"kmTraveled":124,"startFuelLiters":320,"endFuelLiters":182,"actualFuelConsumedLiters":138,"destination":"QL2 ➔ Trạm trộn Phù Ninh","cargo":"Đá 1x2 bê tông"}'::jsonb),
			('CHK-02', '{"id":"CHK-02","code":"CHK-02","yardId":"YARD-PT-01","yardName":"Bãi Xe Trung Tâm - Khai Trường Mỏ Phú Thọ","plate":"19H-056.22","driverName":"Nguyễn Văn Mạnh","driverPhone":"0982.145.882","rfidTag":"RFID-1015-TTC","actionType":"CHECK_IN","timestamp":"08:15 28/10/2026","distanceMeters":8.5,"isWithin20mGeofence":true,"geofenceStatus":"VALID_WITHIN_20M","startOdometer":148424,"endOdometer":148580,"kmTraveled":156,"startFuelLiters":312,"endFuelLiters":268,"actualFuelConsumedLiters":44,"destination":"Mỏ Đá TTC ➔ Cảng Sông Lô","cargo":"Đá Base cấp phối"}'::jsonb),
			('CHK-03', '{"id":"CHK-03","code":"CHK-03","yardId":"YARD-QL2-02","yardName":"Bãi Tập Kết Trung Chuyển QL2 - Phù Ninh","plate":"29C-781.90","driverName":"Lê Văn Cường","driverPhone":"0977.890.123","rfidTag":"RFID-1022-TTC","actionType":"CHECK_IN","timestamp":"09:30 28/10/2026","distanceMeters":15.1,"isWithin20mGeofence":true,"geofenceStatus":"VALID_WITHIN_20M","startOdometer":210302,"endOdometer":210400,"kmTraveled":98,"startFuelLiters":345,"endFuelLiters":310,"actualFuelConsumedLiters":35,"destination":"Khai trường Moong Tầng 3","cargo":"Đá hộc xô bồ"}'::jsonb)
		`)
	}

	// 11. Seed Users & System Settings
	var usrCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&usrCount)
		if usrCount == 0 {
			Pool.Exec(ctx, `
				INSERT INTO users (username, name, ten, role, dept, email, phone, sdt, password_hash, last_login, status) VALUES
				('admin', 'Nguyễn Đức Trường', 'Nguyễn Đức Trường', 'Giám Đốc Mỏ', 'Ban Giám Đốc', 'truongnd@ttcgroup.vn', '0912.345.678', '0912.345.678', 'f078d229547fc3001dba693baa1b39552b0a66677e1dcf08e8da9d0d2dedeb35', '28/10/2026 08:30', 'Hoạt động'),
				('dungnv', 'Nguyễn Văn Dũng', 'Nguyễn Văn Dũng', 'Trưởng Trạm Cân', 'Tổ Vận Hành Trạm Cân', 'dungnv@ttcgroup.vn', '0984.556.789', '0984.556.789', 'f078d229547fc3001dba693baa1b39552b0a66677e1dcf08e8da9d0d2dedeb35', '28/10/2026 09:15', 'Hoạt động'),
				('kientv', 'Trần Văn Kiên', 'Trần Văn Kiên', 'Chỉ Huy Nổ Mìn', 'Tổ Khoan Nổ Mìn & VLNCN', 'kientv@ttcgroup.vn', '0988.341.992', '0988.341.992', 'f078d229547fc3001dba693baa1b39552b0a66677e1dcf08e8da9d0d2dedeb35', '28/10/2026 07:45', 'Hoạt động'),
				('thuynt', 'Nguyễn Thị Thủy', 'Nguyễn Thị Thủy', 'Kế Toán Trưởng', 'Phòng Kế Toán & Vật Tư', 'thuynt@ttcgroup.vn', '0915.678.901', '0915.678.901', 'f078d229547fc3001dba693baa1b39552b0a66677e1dcf08e8da9d0d2dedeb35', '28/10/2026 08:00', 'Hoạt động')
				ON CONFLICT (username) DO NOTHING;

				INSERT INTO user_roles (role, "desc", description, users, level, permission, status) VALUES
				('Ban Giám Đốc', 'Toàn quyền điều hành và phê duyệt mỏ', 'Toàn quyền', '2', 'Cấp 1', 'Toàn quyền hệ thống', 'Hoạt động'),
				('Trưởng Trạm Cân', 'Vận hành cân, chốt phiếu cân, AI OCR ANPR', 'Vận hành cân', '4', 'Cấp 2', 'Quản lý cân & vé xuất bãi', 'Hoạt động'),
				('Chỉ Huy Nổ Mìn', 'Lập hộ chiếu nổ mìn, xuất kho VLNCN', 'Nổ mìn mỏ', '3', 'Cấp 2', 'Quản lý nổ mìn & an toàn', 'Hoạt động'),
				('Kế Toán Trưởng', 'Kế toán công nợ, đối soát hóa đơn VAT điện tử', 'Kế toán tài chính', '3', 'Cấp 2', 'Quản lý tài chính kế toán', 'Hoạt động');

				INSERT INTO user_logs (time, "user", username, name, action, target, ip, status) VALUES
				('28/10/2026 09:42', 'dungnv', 'dungnv', 'Nguyễn Văn Dũng', 'Chốt phiếu cân xuất mỏ Lần 2', 'TK-20261028-001 (88H-042.27)', '192.168.1.100', 'Thành công'),
				('28/10/2026 09:30', 'kientv', 'kientv', 'Trần Văn Kiên', 'Lập hộ chiếu nổ mìn Moong tầng 3', 'PASSPORT-BLAST-1026-01', '192.168.1.105', 'Thành công'),
			('28/10/2026 09:15', 'admin', 'admin', 'Nguyễn Đức Trường', 'Duyệt kế hoạch khai thác tháng 10', 'PLAN-Q4-2026-01', '192.168.1.10', 'Thành công'),
			('28/10/2026 08:45', 'thuynt', 'thuynt', 'Nguyễn Thị Thủy', 'Xuất hóa đơn điện tử VAT', 'INV-20261028-0082', '192.168.1.50', 'Thành công');

			INSERT INTO reports (id, name, item, type, period, plan, actual, diff, unit, status) VALUES
			('REP-01', 'Báo cáo sản lượng khai thác đá nguyên khai', 'Đá nguyên khai', 'Sản xuất', 'Tháng 10/2026', 120000, 128500, '+7.1%', 'tấn', 'Hoàn thành'),
			('REP-02', 'Báo cáo nghiền sàng đá thành phẩm', 'Đá 1x2, 4x6, Base', 'Chế biến', 'Tháng 10/2026', 95000, 98200, '+3.4%', 'tấn', 'Hoàn thành'),
			('REP-03', 'Báo cáo tiêu hao nhiên liệu cơ giới', 'Dầu DO 0.05S', 'Nhiên liệu', 'Tháng 10/2026', 42000, 44850, '+6.8%', 'Lít', 'Cần kiểm tra'),
			('REP-04', 'Báo cáo xuất bán & đối soát công nợ', 'Doanh thu xuất mỏ', 'Tài chính', 'Tháng 10/2026', 28500000000, 29800000000, '+4.6%', 'VNĐ', 'Hoàn thành')
			ON CONFLICT (id) DO UPDATE SET plan = EXCLUDED.plan, actual = EXCLUDED.actual, diff = EXCLUDED.diff, unit = EXCLUDED.unit;

			INSERT INTO settings (code, key, name, val, scope, status) VALUES
			('SET-01', 'WEIGH_TOLERANCE_KG', 'Dung sai bì trạm cân cho phép (kg)', '50', 'Trạm Cân', 'Hoạt động'),
			('SET-02', 'GEOFENCE_RADIUS_METERS', 'Bán kính Geofence trạm bãi xe (mét)', '20', 'Định Vị GPS', 'Hoạt động'),
			('SET-03', 'FUEL_THEFT_THRESHOLD_LITERS', 'Ngưỡng cảnh báo sụt dầu bất thường (lít)', '5.0', 'Nhiên Liệu', 'Hoạt động'),
			('SET-04', 'OCR_ANPR_CONFIDENCE_MIN', 'Độ tin cậy nhận diện biển số AI OCR (%)', '92', 'Camera AI', 'Hoạt động')
			ON CONFLICT (code) DO NOTHING;
		`)
	}

	// Backfill a default demo password (admin) for any accounts missing one.
	Pool.Exec(ctx, `UPDATE users SET password_hash = 'f078d229547fc3001dba693baa1b39552b0a66677e1dcf08e8da9d0d2dedeb35' WHERE password_hash IS NULL OR password_hash = ''`)

	// 12. Seed Catalogs (Ticket types, Vehicles, Stations, Equipment)
	var ttCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM ticket_types").Scan(&ttCount)
	if ttCount == 0 {
		Pool.Exec(ctx, `
			INSERT INTO ticket_types (code, name, ten, loai, description, status) VALUES
			('TT-01', 'Phiếu Cân Bán Hàng Xuất Mỏ', 'Phiếu Cân Bán Hàng', 'Bán hàng', 'Cân 2 lần xuất hàng ra khỏi mỏ cho khách', 'Hoạt động'),
			('TT-02', 'Phiếu Cân Nhập Kho Nguyên Vật Liệu', 'Phiếu Cân Nhập Kho', 'Nhập kho', 'Cân nhập vật tư, phụ gia, nhiên liệu', 'Hoạt động'),
			('TT-03', 'Phiếu Cân Trung Chuyển Nội Bộ', 'Phiếu Cân Nội Bộ', 'Nội bộ', 'Cân vận chuyển đá từ moong về trạm nghiền', 'Hoạt động')
			ON CONFLICT (code) DO NOTHING;

			INSERT INTO vehicle_catalogs (code, name, loai, tai_trong, unit, so_truc, status) VALUES
			('VC-01', 'Xe ben HOWO 371HP 4 Chân', 'Xe ben 4 chân', 30.0, 'tấn', '4 trục', 'Hoạt động'),
			('VC-02', 'Xe ben Chenglong Hải Âu 385HP', 'Xe ben 4 chân', 30.0, 'tấn', '4 trục', 'Hoạt động'),
			('VC-03', 'Xe ben Shacman X3000', 'Xe ben 4 chân', 30.0, 'tấn', '4 trục', 'Hoạt động'),
			('VC-04', 'Xe đầu kéo Howo Mooc Ben', 'Xe đầu kéo', 45.0, 'tấn', '6 trục', 'Hoạt động'),
			('VC-05', 'Xe bồn bê tông trộn 12m3', 'Xe bồn trộn', 25.0, 'tấn', '3 trục', 'Hoạt động')
			ON CONFLICT (code) DO UPDATE SET tai_trong = EXCLUDED.tai_trong, unit = EXCLUDED.unit;

			INSERT INTO stations (code, name, ten, location, ip, capacity, type, status) VALUES
			('ST-01', 'Trạm Cân Điện Tử Cổng 01 (100T)', 'Trạm Cân 01', 'Cổng số 1 - Khai trường mỏ', '192.168.1.100', '100 Tấn', 'Cân bàn thép 18x3m', 'Online Live'),
			('ST-02', 'Trạm Cân Điện Tử Cổng 02 (120T)', 'Trạm Cân 02', 'Cổng số 2 - Khu bãi trung chuyển', '192.168.1.101', '120 Tấn', 'Cân bàn thép 21x3.4m', 'Online Live')
			ON CONFLICT (code) DO NOTHING;

			INSERT INTO equipment (code, name, ten, loai, capacity, location, status) VALUES
			('EQ-01', 'Máy Xúc Bánh Xích Komatsu PC450', 'Máy xúc Komatsu PC450', 'Máy xúc đào', 'Gầu 2.1 m3', 'Moong Khai Thác Tầng 3', 'Hoạt động 100%'),
			('EQ-02', 'Dàn Khoan Khí Nén Thủy Lực D45', 'Dàn khoan khí nén D45', 'Dàn khoan nổ mìn', 'Đường kính lỗ khoan 110mm', 'Tầng 3 (+45m)', 'Hoạt động 100%'),
			('EQ-03', 'Dây Chuyền Nghiền Sàng Số 01 (350T/h)', 'Trạm Nghiền Sàng 01', 'Trạm nghiền sàng', '350 Tấn/giờ', 'Xưởng Chế Biến Đá TTC', 'Hoạt động 100%')
			ON CONFLICT (code) DO NOTHING;
		`)
	}

	// 13. Seed Mining Operations, Statutory Reports & Inventory/Payments
	fmt.Println("🌱 Seeding Mining Operations, Plans, Permits & Inventory/Payments...")
	miningQueries := []string{
		`DELETE FROM statutory_reports`,
		`INSERT INTO statutory_reports (id, code, title, recipient, period, date, mined_volume, unit, tax_amount, env_fee_amount, status, status_label) VALUES
		('STAT-01', 'BC-TK-2026-Q3', 'Báo cáo thống kê hoạt động khai thác khoáng sản Q3/2026 (Mẫu số 01/KTKS)', 'Sở Tài nguyên & Môi trường tỉnh Phú Thọ', 'Quý 3/2026', '15/10/2026', 362400, 'tấn', 4348800000, 1087200000, 'approved', 'Đã phê duyệt'),
		('STAT-02', 'BC-AT-2026-10', 'Báo cáo công tác an toàn lao động & VLNCN tháng 10/2026', 'Sở Công Thương tỉnh Phú Thọ', 'Tháng 10/2026', '25/10/2026', 124500, 'tấn', 1494000000, 373500000, 'approved', 'Đã thẩm định'),
		('STAT-03', 'BC-BVMT-2026-Q3', 'Báo cáo quan trắc môi trường định kỳ mỏ Q3/2026', 'Cục Địa chất & Khoáng sản Việt Nam', 'Quý 3/2026', '30/10/2026', 362400, 'tấn', 4348800000, 1087200000, 'pending', 'Chờ tiếp nhận')`,

		`DELETE FROM resource_taxes`,
		`INSERT INTO resource_taxes (code, mineral_type, mined_volume, unit, tax_price_per_unit, tax_rate, resource_tax_amount, environmental_fee, total_payable, status) VALUES
		('TAX-01', 'Đá vôi xây dựng nguyên khai (Đá hộc)', 362400, 'tấn', 120000, '10%', 4348800000, 1087200000, 5436000000, 'Đã nộp NSNN'),
		('TAX-02', 'Đá thành phẩm nghiền sàng (Đá 1x2, 2x4, 4x6)', 245000, 'tấn', 160000, '10%', 3920000000, 735000000, 4655000000, 'Đã nộp NSNN'),
		('TAX-03', 'Cát nghiền nhân tạo (Mạt đá mi)', 85000, 'tấn', 110000, '10%', 935000000, 340000000, 1275000000, 'Chờ quyết toán')`,

		`DELETE FROM production_stages`,
		`INSERT INTO production_stages (code, stage_number, stage_name, icon, volume_month, volume_ytd, unit, loss_rate, loss_status, measurement_method, description, status) VALUES
		('STG-01', 1, 'Khoan & Nổ mìn Cắt tầng', 'zap', 125000, 1120000, 'tấn', '0.5% (Định mức)', 'normal', 'Đo đạc 3D Laser Scanner & Hộ chiếu nổ mìn', 'Hộ chiếu nổ mìn điện tử phê duyệt Sở Công Thương', 'Đang vận hành'),
		('STG-02', 2, 'Bốc xúc & Vận tải Moong', 'truck', 124300, 1114000, 'tấn', '0.6% (Định mức)', 'normal', 'Camera AI đếm chuyến & Cảm biến gầu xúc', 'Giám sát GPS và hành trình xe ben nội bộ', 'Đang vận hành'),
		('STG-03', 3, 'Nghiền sàng & Phân loại', 'layers', 123500, 1107000, 'tấn', '0.8% (Định mức)', 'normal', 'Cân băng tải động & Ampe kế phụ tải nghiền', 'Tự động bù ẩm và chốt công tơ điện định mức', 'Đang vận hành'),
		('STG-04', 4, 'Tồn trữ Bãi thành phẩm', 'database', 85000, 85000, 'tấn', '0.2%', 'normal', 'Drone RTK quét địa hình tính khối lượng bãi', 'Bay quét Flycam 3D định kỳ 15 ngày/lần', 'Đang vận hành'),
		('STG-05', 5, 'Cân điện tử & Xuất bán', 'scale', 118200, 1058000, 'tấn', '0.0% (Chuẩn)', 'success', 'Trạm cân 120T Keli + Camera AI chụp 4 góc', 'Phiếu cân điện tử mã hóa QR, ký số HĐĐT tức thì', 'Đang vận hành')`,

		`INSERT INTO mining_plans (id, mine, item, annual_target, unit, q1_plan, q1_actual, q2_plan, q2_actual, q3_plan, q3_actual, q4_plan, q4_actual, ytd_actual, completion_rate, status, status_label) VALUES
		('PLAN-01', 'Khai trường Thung Kén (Tầng 1)', 'Đá 1x2 bê tông tiêu chuẩn', 450000, 'tấn', 110000, 114500, 115000, 118200, 115000, 116400, 110000, 38200, 387300, 86.1, 'active', 'Đang thực hiện'),
		('PLAN-02', 'Khai trường Thung Kén (Tầng 2)', 'Đá 4x6 móng công trình', 350000, 'tấn', 85000, 89200, 90000, 91500, 90000, 90800, 85000, 31200, 302700, 86.5, 'active', 'Đang thực hiện'),
		('PLAN-03', 'Khai trường Thung Kén (Bốc phủ)', 'Cát nghiền nhân tạo (Mạt đá)', 200000, 'tấn', 50000, 52100, 50000, 51400, 50000, 50900, 50000, 17800, 172200, 86.1, 'active', 'Đang thực hiện'),
		('PLAN-04', 'Mỏ Cẩm Khê (Vỉa Bắc)', 'Đá 1x2 bê tông tiêu chuẩn', 360000, 'tấn', 90000, 92400, 90000, 91800, 90000, 90200, 90000, 28500, 302900, 84.1, 'active', 'Đang thực hiện'),
		('PLAN-05', 'Mỏ Cẩm Khê (Vỉa Nam)', 'Đá 2x4 móng hạ tầng', 280000, 'tấn', 70000, 71500, 70000, 70800, 70000, 69900, 70000, 24100, 236300, 84.4, 'active', 'Đang thực hiện'),
		('PLAN-06', 'Mỏ Cẩm Khê (Moong)', 'Đá base cấp phối loại 1', 160000, 'tấn', 40000, 41200, 40000, 40500, 40000, 39800, 40000, 13900, 135400, 84.6, 'active', 'Đang thực hiện')
		ON CONFLICT (id) DO UPDATE SET mine=EXCLUDED.mine, item=EXCLUDED.item, annual_target=EXCLUDED.annual_target, unit=EXCLUDED.unit`,

		`INSERT INTO blasting_passports (id, code, mine_name, blast_date, blast_time, location, hole_count, hole_depth_meters, anfo_explosive_kg, emulsion_explosive_kg, detonator_count, designed_rock_volume_m3, powder_factor_kg_per_m3, actual_rock_mined_m3, safety_status, blaster_in_charge, certified_number) VALUES
		('BLAST-01', 'HC-20261028-01', 'Mỏ 1 (Thanh Ba)', '28/10/2026', '11:30', 'Tầng +45m Khai trường Tây', 42, 12.5, 3850, 450, 42, 9800, 0.438, 10250, 'Đã nghiệm thu', 'Nguyễn Văn Hùng', 'CH-VLNCN-2024-089'),
		('BLAST-02', 'HC-20261027-02', 'Mỏ 1 (Thanh Ba)', '27/10/2026', '11:30', 'Tầng +30m Khai trường Nam', 36, 11.8, 3200, 380, 36, 8200, 0.436, 8560, 'Đã nghiệm thu', 'Trần Đình Trọng', 'CH-VLNCN-2023-112'),
		('BLAST-03', 'HC-20261029-01', 'Mỏ 2 (Cẩm Khê)', '29/10/2026', '16:30', 'Tầng +60m Vỉa Bắc', 48, 13.0, 4500, 520, 48, 11500, 0.436, 0, 'Chờ kích nổ', 'Hoàng Minh Đức', 'CH-VLNCN-2025-045')
		ON CONFLICT (id) DO NOTHING`,

		`DELETE FROM crusher_plants`,
		`INSERT INTO crusher_plants (plant_code, plant_name, capacity_ton_per_hour, input_rock_today_tons, power_consumption_kwh, kwh_per_ton, yield_breakdown) VALUES
		('PLANT-01', 'Dây Chuyền Nghiền Sàng Số 01 (350T/h)', 350, 2680, 3216, 1.20, '[{"productName":"Đá 1x2 bê tông tiêu chuẩn","producedTons":1340,"yieldPercent":50,"standardRate":48},{"productName":"Đá 4x6 móng công trình","producedTons":804,"yieldPercent":30,"standardRate":32},{"productName":"Cát nghiền nhân tạo (Mạt đá)","producedTons":402,"yieldPercent":15,"standardRate":14},{"productName":"Đá mi bụi sàng tuyển","producedTons":134,"yieldPercent":5,"standardRate":6}]'::jsonb),
		('PLANT-02', 'Dây Chuyền Nghiền Sàng Số 02 (250T/h)', 250, 1850, 2312, 1.25, '[{"productName":"Đá 1x2 sàng tuyển mác cao","producedTons":925,"yieldPercent":50,"standardRate":48},{"productName":"Đá 2x4 đổ bê tông","producedTons":555,"yieldPercent":30,"standardRate":30},{"productName":"Cát nhân tạo VSI","producedTons":277,"yieldPercent":15,"standardRate":16},{"productName":"Bột đá phụ gia","producedTons":93,"yieldPercent":5,"standardRate":6}]'::jsonb)`,

		`INSERT INTO equipment_fuel_logs (id, equipment_code, equipment_name, category, operator_name, hours_worked_today, total_hours_meter, fuel_quota_liters_per_hour, actual_fuel_issued_liters, actual_fuel_consumed_liters, fuel_variance_liters, variance_status, location, maintenance_status) VALUES
		('FUEL-01', 'EQ-EXC-01', 'Máy Xúc Bánh Xích Komatsu PC450 #01', 'Máy xúc moong', 'Nguyễn Văn Dũng', 7.5, 4820.5, 32.0, 240, 234.5, -5.5, 'Tiết kiệm (-5.5L)', 'Moong Khai Thác Tầng 3', 'Bảo dưỡng định kỳ A'),
		('FUEL-02', 'EQ-EXC-02', 'Máy Xúc Bánh Xích Hitachi ZX350 #02', 'Máy xúc moong', 'Trần Văn Kiên', 8.0, 3915.0, 28.0, 224, 221.0, -3.0, 'Tiết kiệm (-3.0L)', 'Moong Khai Thác Tầng 2', 'Hoạt động tốt'),
		('FUEL-03', 'EQ-TRK-01', 'Xe Ben Howo 371HP 8x4 #01', 'Xe ben nội bộ', 'Lê Hữu Thắng', 8.0, 5210.0, 14.5, 120, 116.0, -4.0, 'Tiết kiệm (-4.0L)', 'Tuyến Moong ➔ Trạm Nghiền', 'Hoạt động tốt'),
		('FUEL-04', 'EQ-TRK-02', 'Xe Ben Howo 371HP 8x4 #02', 'Xe ben nội bộ', 'Hoàng Minh Đức', 8.0, 4980.0, 14.5, 120, 118.5, -1.5, 'Đạt định mức', 'Tuyến Moong ➔ Trạm Nghiền', 'Hoạt động tốt'),
		('FUEL-05', 'EQ-DRL-01', 'Dàn Khoan Thủy Lực Furukawa D45', 'Máy khoan tự hành', 'Phạm Văn Nam', 6.5, 2150.0, 22.0, 145, 142.0, -3.0, 'Tiết kiệm (-3.0L)', 'Tầng +45m Khai trường Tây', 'Hoạt động tốt')
		ON CONFLICT (id) DO NOTHING`,

		`INSERT INTO mining_permits (id, code, title, mine_name, category, category_label, issuer, license_number, issue_date, expiry_date, capacity, approved_reserve, mined_so_far, mined_percent, depth_level, area, coordinates, status, status_label, days_remaining, files, notes) VALUES
		('PERMIT-01', 'GP-28102018-BTNMT', 'Giấy phép khai thác khoáng sản Mỏ Đá Vôi Thanh Ba', 'Mỏ 1 (Thanh Ba)', 'mining_license', 'Giấy phép khai thác mỏ', 'Bộ Tài nguyên và Môi trường', 'Số 2810/GP-BTNMT', '28/10/2018', '28/10/2038', '1.000.000 Tấn/năm', '20.000.000 Tấn', '5.240.000 Tấn', 26, 'Mức cao +85m đến +15m', '48.5 Hecta', '21.3210°N, 105.3280°E', 'valid', 'Còn hiệu lực', 4380, '[{"name":"Quyet_Dinh_Cap_Phep_2810.pdf","size":"4.2 MB","type":"pdf","url":"#"}]'::jsonb, 'Được phép khai thác đá vôi làm VLXD thông thường'),
		('PERMIT-02', 'GP-15062020-UBND', 'Giấy phép khai thác đá xây dựng Mỏ Cẩm Khê', 'Mỏ 2 (Cẩm Khê)', 'mining_license', 'Giấy phép khai thác mỏ', 'UBND Tỉnh Phú Thọ', 'Số 1506/GP-UBND', '15/06/2020', '15/06/2035', '800.000 Tấn/năm', '12.000.000 Tấn', '3.120.000 Tấn', 26, 'Mức cao +90m đến +25m', '35.2 Hecta', '21.4120°N, 105.2150°E', 'valid', 'Còn hiệu lực', 3215, '[{"name":"Giay_Phep_Cam_Khe_1506.pdf","size":"3.8 MB","type":"pdf","url":"#"}]'::jsonb, 'Khai thác đá xây dựng và cát nhân tạo'),
		('PERMIT-03', 'GP-118-SCT', 'Giấy phép sử dụng Vật liệu nổ công nghiệp năm 2026', 'Mỏ 1 (Thanh Ba)', 'legal_profile', 'Hồ sơ pháp lý VLNCN', 'Sở Công Thương tỉnh Phú Thọ', 'Số 118/GP-SCT', '01/01/2026', '31/12/2026', '180 Tấn/năm', '180 Tấn ANFO', '142 Tấn', 78, 'Kho mìn cấp 1', 'Toàn mỏ', '21.3210°N, 105.3280°E', 'valid', 'Còn hiệu lực', 128, '[]'::jsonb, 'Định mức nổ mìn ca ngày'),
		('PERMIT-04', 'DTM-2018-BTNMT', 'Quyết định phê duyệt Báo cáo Đánh giá tác động môi trường (ĐTM)', 'Mỏ 1 (Thanh Ba)', 'environmental', 'Giấy phép môi trường & ĐTM', 'Bộ Tài nguyên và Môi trường', 'Số 315/QĐ-BTNMT', '15/08/2018', '28/10/2038', 'Hồ lắng 3 cấp 450 m³/ngày', 'Ký quỹ 4.2 Tỷ', 'Đã ký quỹ 100%', 100, 'Khu xử lý nước thải', '48.5 Hecta', '21.3210°N, 105.3280°E', 'valid', 'Còn hiệu lực', 4380, '[]'::jsonb, 'Báo cáo quan trắc định kỳ 4 đợt/năm'),
		('PERMIT-05', 'TL-2018-HDKS', 'Báo cáo kết quả thăm dò & Phê duyệt trữ lượng khoáng sản', 'Mỏ 1 (Thanh Ba)', 'geological', 'Hồ sơ trữ lượng & Bản đồ', 'Hội đồng Đánh giá Trữ lượng Khoáng sản Quốc gia', 'Số 89/HĐTL-QG', '10/05/2018', '28/10/2038', 'Cấp 121 + 122', '24.800.000 Tấn', '5.240.000 Tấn', 21, 'Độ sâu đến +15m', '48.5 Hecta', '21.3210°N, 105.3280°E', 'valid', 'Đã nghiệm thu đo đạc', 4380, '[]'::jsonb, 'Bản đồ hiện trạng mỏ cập nhật định kỳ 6 tháng/lần'),
		('PERMIT-06', 'HD-THUEDAT-2018', 'Hợp đồng thuê đất mở mỏ khai thác đá', 'Mỏ 1 (Thanh Ba)', 'contract', 'Hợp đồng mỏ & dịch vụ', 'Sở Tài nguyên & Môi trường tỉnh Phú Thọ', 'Số 48/HĐTĐ-STNMT', '15/11/2018', '15/11/2048', '48.5 Hecta', 'Thời hạn 30 năm', 'Thuê 48.5 ha', 100, 'Toàn bộ diện tích mỏ', '48.5 Hecta', '21.3210°N, 105.3280°E', 'valid', 'Đang thực hiện', 8120, '[]'::jsonb, 'Đã nộp tiền thuê đất hàng năm đầy đủ'),
		('PERMIT-07', 'HD-GAET-2026', 'Hợp đồng dịch vụ Khoan nổ mìn trọn gói 2026', 'Mỏ 1 & Mỏ 2', 'contract', 'Hợp đồng mỏ & dịch vụ', 'Tổng Công Ty Kinh Tế Kỹ Thuật CNQP (GAET)', 'Số 26/HĐ-GAET-TTC', '02/01/2026', '31/12/2026', '1.400.000 Tấn đá', '180 Tấn VLNCN', '142 Tấn', 78, 'Toàn bộ khai trường', '83.7 Hecta', '21.3210°N, 105.3280°E', 'valid', 'Đang thực hiện', 128, '[]'::jsonb, 'Dịch vụ nạp nổ mìn trọn gói theo hộ chiếu'),
		('PERMIT-08', 'HD-CIENCO4-2026', 'Hợp đồng cung cấp đá bê tông dự án Cao tốc', 'Mỏ 1 (Thanh Ba)', 'contract', 'Hợp đồng mỏ & dịch vụ', 'Tập Đoàn CIENCO 4', 'Số 112/HĐ-C4-TTC', '10/02/2026', '30/06/2027', '1.200.000 Tấn Đá 1x2 & 4x6', 'Giá trị 240 Tỷ đ', '620.000 Tấn', 52, 'Gói thầu XL-02 Cao tốc', 'Tuyến cao tốc', '21.3210°N, 105.3280°E', 'valid', 'Đang thực hiện', 310, '[]'::jsonb, 'Cung ứng đá cấp phối và đá bê tông mác cao')
		ON CONFLICT (id) DO NOTHING`,

		`DELETE FROM inventory_inbound`,
		`INSERT INTO inventory_inbound (code, source, loc, item, qty, quantity, unit, date, status) VALUES
		('NK-281025-01', 'Moong Khai Thác Tầng 3', 'Bãi Đá Hộc & Nguyên Khai', 'Đá hộc khai thác tầng', 450.00, 450.00, 'tấn', '28/10/2026', 'Đã nhập bãi'),
		('NK-281025-02', 'Trạm Nghiền Sàng Số 01', 'Bãi Đá Thành Phẩm 01', 'Đá 1x2 bê tông tiêu chuẩn', 380.00, 380.00, 'tấn', '28/10/2026', 'Đã nhập bãi'),
		('NK-281025-03', 'Trạm Nghiền Sàng Số 01', 'Bãi Đá 4x6 Kè Móng', 'Đá 4x6 móng công trình', 320.00, 320.00, 'tấn', '27/10/2026', 'Đã nhập bãi'),
		('NK-281025-04', 'Dây Chuyền Nghiền Sàng 02', 'Kho Cát Nghiền Mái Che', 'Cát nghiền nhân tạo (Mạt đá)', 280.50, 280.50, 'tấn', '27/10/2026', 'Đã nhập bãi')`,

		`DELETE FROM inventory_outbound`,
		`INSERT INTO inventory_outbound (code, customer, dest, item, qty, quantity, unit, date, status) VALUES
		('XK-281025-01', 'Công ty CP Đầu Tư Xây Dựng 319', 'Dự án KCN Phú Hà', 'Đá 1x2 bê tông tiêu chuẩn', 380.00, 380.00, 'tấn', '28/10/2026', 'Đã xuất bãi'),
		('XK-281025-02', 'Tập Đoàn CIENCO 4 (Cao Tốc)', 'Gói thầu XL-02 Cao tốc', 'Đá 4x6 móng công trình', 290.00, 290.00, 'tấn', '28/10/2026', 'Đã xuất bãi'),
		('XK-281025-03', 'Công ty Bê Tông Việt Trì', 'Trạm trộn Bê tông Việt Trì', 'Đá 1x2 bê tông mác 350', 180.00, 180.00, 'tấn', '27/10/2026', 'Đã xuất bãi'),
		('XK-281025-04', 'Tổng Công Ty XD Trường Sơn', 'Dự án Cầu Phong Châu mới', 'Đá base cấp phối loại 1', 140.00, 140.00, 'tấn', '27/10/2026', 'Đã xuất bãi')`,

		`DELETE FROM inventory_stocktake`,
		`INSERT INTO inventory_stocktake (code, zone, item, volume, survey, erp, book, actual, diff, quantity, unit, date, status) VALUES
		('KK-20261028-01', 'Bãi Đá 1x2 Thành Phẩm Số 1 (Lô A)', 'Đá 1x2 bê tông tiêu chuẩn', 7517.0, 7517.0, 7520.0, 7520.0, 7517.0, -3.0, 7517.0, 'tấn', '28/10/2026', 'Khớp 99.9% (Chiếm 55%)'),
		('KK-20261028-02', 'Bãi Đá 4x6 Móng Hạ Tầng (Lô B)', 'Đá 4x6 móng công trình', 5056.0, 5056.0, 5050.0, 5050.0, 5056.0, 6.0, 5056.0, 'tấn', '28/10/2026', 'Khớp 100% (Chiếm 32%)'),
		('KK-20261028-03', 'Kho Phụ Trợ 02 (Lô C)', 'Cát nghiền nhân tạo (Mạt đá)', 611.0, 611.0, 615.0, 615.0, 611.0, -4.0, 611.0, 'tấn', '27/10/2026', 'Khớp 99.4%')`,

		`DELETE FROM inventory_movements`,
		`INSERT INTO inventory_movements (code, from_loc, to_loc, item, qty, quantity, unit, date, status) VALUES
		('DC-281025-01', 'Bãi Nổ Mìn Tầng 2', 'Phễu Nghiền Thô 01', 'Đá hộc cấp liệu nghiền 1x2 & 4x6', 850.0, 850.0, 'tấn', '28/10/2026', 'Hoàn tất'),
		('DC-281025-02', 'Trạm Nghiền Sàng Số 01', 'Bãi Chứa Thành Phẩm Đá 1x2', 'Đá 1x2 sàng tuyển tiêu chuẩn', 480.0, 480.0, 'tấn', '28/10/2026', 'Hoàn tất'),
		('DC-281025-03', 'Trạm Nghiền Sàng Số 01', 'Bãi Chứa Móng Hạ Tầng Đá 4x6', 'Đá 4x6 móng đường cao tốc', 290.0, 290.0, 'tấn', '27/10/2026', 'Hoàn tất')`,

		`DELETE FROM payments_debt`,
		`INSERT INTO payments_debt (code, partner, customer, limit_amount, limit_val, balance, amount, debt, due, status) VALUES
		('CN-001', 'Công ty CP Đầu Tư Xây Dựng 319', 'Công ty CP Đầu Tư Xây Dựng 319', '3.0 Tỷ đ', '3.0 Tỷ đ', '320.000.000 đ', '320.000.000 đ', '320.000.000 đ', '28/11/2026', 'Trong hạn'),
		('CN-002', 'Tổng Công Ty XD Trường Sơn', 'Tổng Công Ty XD Trường Sơn', '2.0 Tỷ đ', '2.0 Tỷ đ', '540.000.000 đ', '540.000.000 đ', '540.000.000 đ', '12/11/2026', 'Trong hạn'),
		('CN-003', 'Tập Đoàn CIENCO 4 (Cao Tốc Phú Thọ)', 'Tập Đoàn CIENCO 4', '5.0 Tỷ đ', '5.0 Tỷ đ', '1.250.000.000 đ', '1.250.000.000 đ', '1.250.000.000 đ', '02/11/2026', 'Trong hạn'),
		('CN-004', 'Công ty Bê Tông Việt Trì', 'Công ty Bê Tông Việt Trì', '800 Triệu đ', '800 Triệu đ', '180.000.000 đ', '180.000.000 đ', '180.000.000 đ', '05/11/2026', 'Quá hạn')
		ON CONFLICT DO NOTHING`,

		`DELETE FROM payments_invoices`,
		`INSERT INTO payments_invoices (code, partner, customer, amount, balance, due, status, created_at) VALUES
		('HD-20260903-01', 'Công ty CP Đầu Tư Xây Dựng 319', 'Công ty CP Đầu Tư Xây Dựng 319', '145.000.000 đ', '0 đ', '28/10/2026', 'paid', NOW()),
		('HD-20260903-02', 'Tổng Công Ty XD Trường Sơn', 'Tổng Công Ty XD Trường Sơn', '85.500.000 đ', '0 đ', '28/10/2026', 'paid', NOW()),
		('HD-20260902-01', 'Tập Đoàn CIENCO 4', 'Tập Đoàn CIENCO 4', '320.000.000 đ', '0 đ', '27/10/2026', 'paid', NOW() - INTERVAL '1 day'),
		('HD-20260901-01', 'Công ty Bê Tông Việt Trì', 'Công ty Bê Tông Việt Trì', '210.000.000 đ', '0 đ', '26/10/2026', 'paid', NOW() - INTERVAL '2 day')
		ON CONFLICT DO NOTHING`,

		`INSERT INTO payments_reconcile (code, period, day, scale_rev, acc_rev, scale_revenue, erp_revenue, diff, status, date) VALUES
		('DS-281026-01', 'Ca 1 - 28/10/2026', '28/10/2026', '114.800.000 đ', '114.800.000 đ', '114.800.000 đ', '114.800.000 đ', '0 đ (Khớp 100%)', 'Khớp 100%', '28/10/2026'),
		('DS-271026-02', 'Ca 2 - 27/10/2026', '27/10/2026', '98.500.000 đ', '98.500.000 đ', '98.500.000 đ', '98.500.000 đ', '0 đ (Khớp 100%)', 'Khớp 100%', '27/10/2026'),
		('DS-261026-03', 'Ca 3 - 26/10/2026', '26/10/2026', '105.200.000 đ', '105.000.000 đ', '105.200.000 đ', '105.000.000 đ', '+200.000 đ (Chênh lệch)', 'Cần kiểm tra', '26/10/2026')
		ON CONFLICT DO NOTHING`,

		`INSERT INTO hr_shifts (code, name, start_time, end_time, hours, status) VALUES
		('SHIFT-01', 'Ca Sáng (Ca 1)', '06:00', '14:00', 8.0, 'Hoạt động'),
		('SHIFT-02', 'Ca Chiều (Ca 2)', '14:00', '22:00', 8.0, 'Hoạt động'),
		('SHIFT-03', 'Ca Hành Chính', '07:30', '16:30', 8.0, 'Hoạt động')
		ON CONFLICT (code) DO NOTHING`,
	}
	for i, q := range miningQueries {
		if _, err := Pool.Exec(ctx, q); err != nil {
			log.Printf("⚠️ Mining seed query %d error: %v", i+1, err)
		}
	}

	// 13.5 Enterprise Business & Quarry Mining E-Sign Contracts
	Pool.Exec(ctx, `DELETE FROM hr_esign_signers; DELETE FROM hr_esign_documents;`)
	Pool.Exec(ctx, `
		INSERT INTO hr_esign_documents (id, code, contract_id, contract_code, subject, status, document_type, partner_name, contract_value, volume) VALUES
		('ES-01', 'ES-2026-01', 'CTR-KD-01', 'HĐMB-2026-ĐC', 'Hợp đồng Cung cấp 150.000 Tấn Đá thương phẩm - Tập đoàn Đèo Cả', 'cho_ky', 'kinh_doanh', 'Công ty CP Tập Đoàn Đèo Cả (Dự án Cao tốc Bắc - Nam)', 33000000000, '150.000 Tấn Đá 1x2 & Base'),
		('ES-02', 'ES-2026-02', 'CTR-KT-02', 'HĐBX-2026-TK', 'Hợp đồng Khoán bốc xúc & Vận tải đất đá moong mỏ Thung Kén (200.000 Tấn)', 'cho_ky', 'khai_thac', 'Đội xe Cơ giới & Vận tải Mỏ Tân Tiến', 14500000000, '200.000 Tấn Đất đá moong'),
		('ES-03', 'ES-2026-03', 'CTR-TB-03', 'HĐTB-2026-FL', 'Hợp đồng Thuê dây chuyền nghiền sàng đá Terx Finlay di động 250 T/h', 'cho_ky', 'thiet_bi', 'Công ty CP Thiết bị Khai khoáng Vinamac', 5400000000, '12 tháng vận hành'),
		('ES-04', 'ES-2026-04', 'CTR-NM-04', 'BBNT-2026-NM08', 'Biên bản Nghiệm thu nổ mìn phá đá nguyên khai đợt 8 (35.000 Tấn)', 'cho_ky', 'no_min', 'Công ty Công nghiệp Hóa chất Mỏ MIC Miền Bắc', 2850000000, '35.000 Tấn Đá nguyên khai'),
		('ES-05', 'ES-2026-05', 'CTR-VT-05', 'HĐNL-2026-PLX', 'Hợp đồng Cung cấp nhiên liệu Dầu DO 0.05S cho dàn máy công trình mỏ', 'da_chuyen', 'vat_tu', 'Chi nhánh Xăng dầu Petrolimex Phú Thọ', 8200000000, '350.000 Lít Dầu DO'),
		('ES-06', 'ES-2026-06', 'CTR-KD-06', 'HĐMB-2026-BQP319', 'Hợp đồng Mua bán đá Hộc kè đê & Đá base 80.000 tấn - Tổng Cty 319 BQP', 'da_ky', 'kinh_doanh', 'Tổng Công ty 319 - Bộ Quốc Phòng', 18800000000, '80.000 Tấn Đá Base & Hộc'),
		('ES-07', 'ES-2026-07', 'CTR-LD-07', 'HĐLĐ-2026-KS', 'Hợp đồng Lao động Kỹ sư trưởng Giám sát An toàn & Nổ mìn mỏ', 'cho_ky', 'lao_dong', 'Khai trường Mỏ Đá Thung Kén', 380000000, 'Vô thời hạn')
		ON CONFLICT (id) DO NOTHING;
	`)

	Pool.Exec(ctx, `
		INSERT INTO hr_esign_signers (id, document_id, signer_type, employee_id, employee_name, role_title, sign_order, sign_status) VALUES
		('SG-01', 'ES-01', 'dai_dien_mo', 'EMP-01', 'Nguyễn Đức Trường', 'Tổng Giám Đốc Mỏ Đá', 1, 'cho_ky'),
		('SG-02', 'ES-01', 'dai_dien_doi_tac', NULL, 'Vũ Quốc Đạt (Đại diện Đèo Cả)', 'Giám đốc Điều hành Dự án', 2, 'cho_ky'),
		('SG-03', 'ES-02', 'dai_dien_mo', 'EMP-09', 'Phạm Hoàng Nam', 'Giám đốc Khai trường Mỏ', 1, 'cho_ky'),
		('SG-04', 'ES-02', 'nha_thau_van_tai', 'EMP-03', 'Trần Văn Kiên', 'Đội trưởng Vận tải Mỏ', 2, 'cho_ky'),
		('SG-05', 'ES-03', 'dai_dien_mo', 'EMP-04', 'Lê Hữu Thắng', 'Phó Giám đốc Kỹ thuật & Thiết bị', 1, 'cho_ky'),
		('SG-06', 'ES-03', 'nha_cung_cap', NULL, 'Hoàng Văn Hải (Đại diện Vinamac)', 'Giám đốc Kỹ thuật Thiết bị', 2, 'cho_ky'),
		('SG-07', 'ES-04', 'chi_huy_no_min', 'EMP-05', 'Trần Đình Trọng', 'Chỉ huy Nổ mìn Mỏ', 1, 'cho_ky'),
		('SG-08', 'ES-04', 'giam_sat_an_toan', 'EMP-02', 'Nguyễn Văn Dũng', 'Giám sát Kỹ thuật An toàn Mỏ', 2, 'cho_ky'),
		('SG-09', 'ES-05', 'truong_phong_vat_tu', NULL, 'Nguyễn Thị Thủy', 'Trưởng phòng Cung ứng Vật tư', 1, 'da_chuyen'),
		('SG-10', 'ES-05', 'ke_toan_kho_bon', 'EMP-06', 'Nguyễn Văn Mạnh', 'Kế toán trưởng & Quản lý Kho bồn', 2, 'cho_ky'),
		('SG-11', 'ES-06', 'dai_dien_mo', 'EMP-01', 'Nguyễn Đức Trường', 'Tổng Giám Đốc Mỏ Đá', 1, 'da_ky'),
		('SG-12', 'ES-06', 'dai_dien_doi_tac', NULL, 'Đại tá Lê Văn Long (Tổng Cty 319)', 'Chỉ huy trưởng Công trình', 2, 'da_ky'),
		('SG-13', 'ES-07', 'nhan_vien', 'EMP-03', 'Trần Văn Kiên', 'Kỹ sư trưởng Khai thác Mỏ', 1, 'cho_ky'),
		('SG-14', 'ES-07', 'quan_ly_nhan_su', NULL, 'Nguyễn Thị Thủy', 'Ban Giám Đốc Nhân sự', 2, 'cho_ky')
		ON CONFLICT (id) DO NOTHING;
	`)

	// 14. Seed HRM extended modules (recruitment, decisions, insurance, training, KPI, workflow, timesheet/payroll)
	var extCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM hr_recruitment_requests").Scan(&extCount)
	if extCount == 0 {
		fmt.Println("🌱 Seeding HRM extended modules (recruitment, decisions, onboarding, insurance, training, KPI/OKR, workflow)...")

		Pool.Exec(ctx, `
			INSERT INTO hr_recruitment_requests (id, code, title, department, position, headcount, reason, requested_by, approval_status, current_step) VALUES
			('RR-01', 'ĐXTĐ-1026-01', 'Tuyển 2 kỹ sư khai thác mỏ', 'Phòng Kỹ Thuật & An Toàn', 'Kỹ sư khai thác mỏ', 2, 'Thiếu nhân sự khai thác so với định biên Q4/2026', 'Nguyễn Đức Trường', 'da_duyet', 2),
			('RR-02', 'ĐXTĐ-1026-02', 'Tuyển 1 lái xe ben có kinh nghiệm mỏ', 'Đội Cơ Giới & Vận Tải Mỏ', 'Tài xế xe ben', 1, 'Bổ sung ca 2 vận tải nội bộ', 'Hoàng Minh Đức', 'cho_duyet', 1),
			('RR-03', 'ĐXTĐ-1026-03', 'Tuyển 2 nhân viên trạm cân', 'Tổ Vận Hành Trạm Cân', 'Nhân viên vận hành trạm cân', 2, 'Tăng ca trực 2 cổng cân liên tục', 'Nguyễn Văn Dũng', 'tu_choi', 1)
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_recruitment_campaigns (id, code, name, request_id, position, department, open_date, recruit_deadline, target_headcount, hired_count, source, status) VALUES
			('CAM-01', 'TD-1026-01', 'Chiến dịch tuyển dụng Kỹ sư khai thác mỏ Q4/2026', 'RR-01', 'Kỹ sư khai thác mỏ', 'Phòng Kỹ Thuật & An Toàn', '01/10/2026', '30/11/2026', 2, 0, 'Vietnamworks + Giới thiệu nội bộ', 'dang_tuyen')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_candidates (id, code, name, phone, email, position, source, pipeline_stage, note) VALUES
			('CAND-01', 'UV-001', 'Nguyễn Văn Hùng', '0901.111.222', 'hung.nv@gmail.com', 'Kỹ sư khai thác mỏ', 'Vietnamworks', 'pv_dat', 'Kinh nghiệm 5 năm tại mỏ Yên Bái'),
			('CAND-02', 'UV-002', 'Trần Thị Lan', '0902.333.444', 'lan.tt@gmail.com', 'Kỹ sư khai thác mỏ', 'Giới thiệu nội bộ', 'ung_tuyen_dat', 'Tốt nghiệp ĐH Mỏ Địa Chất'),
			('CAND-03', 'UV-003', 'Lê Minh Cường', '0903.555.666', 'cuong.lm@gmail.com', 'Kỹ sư khai thác mỏ', 'Vietnamworks', 'khong_trung_tuyen', 'Thiếu chứng chỉ an toàn mỏ'),
			('CAND-04', 'UV-004', 'Phạm Quốc Bảo', '0904.777.888', 'bao.pq@gmail.com', 'Lái xe ben', 'Facebook Jobs', 'len_lich_pv', 'Có bằng lái hạng FC')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_candidate_applications (id, candidate_id, campaign_id, pipeline_stage, applied_at) VALUES
			('CA-01', 'CAND-01', 'CAM-01', 'pv_dat', '05/10/2026'),
			('CA-02', 'CAND-02', 'CAM-01', 'ung_tuyen_dat', '08/10/2026'),
			('CA-03', 'CAND-03', 'CAM-01', 'khong_trung_tuyen', '12/10/2026')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_interview_schedules (id, candidate_id, candidate_name, position, interview_date, interview_time, interviewer, round, result, status) VALUES
			('IV-01', 'CAND-01', 'Nguyễn Văn Hùng', 'Kỹ sư khai thác mỏ', '20/10/2026', '09:00', 'Nguyễn Đức Trường', 1, 'dat', 'da_co_ket_qua'),
			('IV-02', 'CAND-04', 'Phạm Quốc Bảo', 'Lái xe ben', '25/10/2026', '10:30', 'Hoàng Minh Đức', 1, '', 'cho_phong_van')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_decisions (id, code, employee_id, employee_name, decision_type, effective_date, content_detail, status) VALUES
			('DEC-01', 'QD-1026-01', 'EMP-09', 'Phạm Hoàng Nam', 'tang_luong', '01/11/2026', '{"luong_cu":21000000,"luong_moi":23000000,"muc_do":"10%"}', 'cho_duyet'),
			('DEC-02', 'QD-1026-02', 'EMP-02', 'Nguyễn Văn Dũng', 'khen_thuong', '28/10/2026', '{"muc_do":"Công nhận xuất sắc","so_tien":3000000}', 'da_ban_hanh')
			ON CONFLICT (id) DO NOTHING;
		`)



		Pool.Exec(ctx, `
			INSERT INTO hr_insurance_records (id, employee_id, employee_name, insurance_number, type, participation_date, social_insurance_base, status) VALUES
			('INS-01', 'EMP-01', 'Nguyễn Đức Trường', 'SBH-025081234', 'BHXH/BHYT/BHTN', '15/04/2018', 38000000, 'dang_tham_gia'),
			('INS-02', 'EMP-02', 'Nguyễn Văn Dũng', 'SBH-025085678', 'BHXH/BHYT/BHTN', '01/07/2020', 16500000, 'dang_tham_gia')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_insurance_declarations (id, code, employee_id, employee_name, declaration_type, period, status) VALUES
			('INSDEC-01', 'TANG-1026-01', 'EMP-04', 'Lê Hữu Thắng', 'tang', '10/2026', 'cho_ke_khai'),
			('INSDEC-02', 'GIAM-1026-01', 'EMP-06', 'Nguyễn Văn Mạnh', 'giam', '10/2026', 'cho_ke_khai')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_training_courses (id, code, name, category, provider, start_date, end_date, venue, instructor, duration_hours, attendance_threshold, status) VALUES
			('CRS-01', 'DT-1026-01', 'An toàn lao động & VLNCN', 'An toàn', 'Sở Công Thương Phú Thọ', '25/10/2026', '26/10/2026', 'Hội trường Mỏ TTC', 'Nguyễn Văn Hòa', 16, 80, 'dang_hoc')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_training_sessions (id, course_id, session_name, session_date, start_time, end_time, instructor, status) VALUES
			('SES-01', 'CRS-01', 'Buổi 1 - Lý thuyết ATLĐ & VLNCN', '25/10/2026', '08:00', '11:30', 'Nguyễn Văn Hòa', 'da_dien_ra'),
			('SES-02', 'CRS-01', 'Buổi 2 - Thực hành & kiểm tra', '26/10/2026', '13:30', '17:00', 'Nguyễn Văn Hòa', 'da_dien_ra')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_training_participants (id, course_id, employee_id, employee_name, attendance_rate, completion_status) VALUES
			('PT-01', 'CRS-01', 'EMP-09', 'Phạm Hoàng Nam', 100, 'hoan_thanh'),
			('PT-02', 'CRS-01', 'EMP-03', 'Trần Văn Kiên', 50, 'khong_hoan_thanh')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_training_attendances (id, session_id, participant_id, employee_id, employee_name, attendance_status) VALUES
			('TA-01', 'SES-01', 'PT-01', 'EMP-09', 'Phạm Hoàng Nam', 'co_mat'),
			('TA-02', 'SES-02', 'PT-01', 'EMP-09', 'Phạm Hoàng Nam', 'co_mat'),
			('TA-03', 'SES-01', 'PT-02', 'EMP-03', 'Trần Văn Kiên', 'co_mat'),
			('TA-04', 'SES-02', 'PT-02', 'EMP-03', 'Trần Văn Kiên', 'vang_mat')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_kpi_targets (id, employee_id, employee_name, period, kpi_code, kpi_name, target_value, actual_value, unit, weight_percent, status) VALUES
			('KPI-01', 'EMP-01', 'Nguyễn Đức Trường', '10/2026', 'KPI-SL', 'Sản lượng khai thác', 120000, 128500, 'tấn', 40, 'da_ban_hanh'),
			('KPI-02', 'EMP-01', 'Nguyễn Đức Trường', '10/2026', 'KPI-AT', 'An toàn lao động', 100, 100, '%', 30, 'da_ban_hanh'),
			('KPI-03', 'EMP-01', 'Nguyễn Đức Trường', '10/2026', 'KPI-DT', 'Doanh thu', 28.5, 29.8, 'tỷ', 30, 'da_ban_hanh')
			ON CONFLICT (id) DO UPDATE SET unit = EXCLUDED.unit, target_value = EXCLUDED.target_value, actual_value = EXCLUDED.actual_value;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_okr_objectives (id, code, title, level, period, progress_percent, status) VALUES
			('OKR-01', 'OKR-Q4-01', 'Nâng sản lượng khai thác 10% quý 4/2026', 'company', 'Q4/2026', 65, 'dang_cho')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_okr_key_results (id, objective_id, key_result, start_value, target_value, current_value, unit, progress_percent) VALUES
			('KR-01', 'OKR-01', 'Tăng sản lượng đá nguyên khai', 100000, 120000, 113000, 'tấn', 65)
			ON CONFLICT (id) DO UPDATE SET unit = EXCLUDED.unit;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_workflow_templates (id, code, name, applies_to, definition, status) VALUES
			('WT-01', 'WFD-X-NGHI-PHEP', 'Workflow duyệt đơn nghỉ phép', 'request', '{"nodes":[{"id":"n1","type":"start"},{"id":"n2","type":"approval","approver_role":"truong_phong"},{"id":"n3","type":"approval","approver_role":"giam_doc"},{"id":"n4","type":"end"}],"edges":[{"from":"n1","to":"n2"},{"from":"n2","to":"n3"},{"from":"n3","to":"n4"}]}', 'hoat_dong'),
			('WT-02', 'WFD-X-DE-XUAT-TD', 'Workflow duyệt đề xuất tuyển dụng', 'recruitment_request', '{"nodes":[{"id":"n1","type":"start"},{"id":"n2","type":"approval","approver_role":"truong_phong"},{"id":"n3","type":"approval","approver_role":"giam_doc"},{"id":"n4","type":"end"}],"edges":[{"from":"n1","to":"n2"},{"from":"n2","to":"n3"},{"from":"n3","to":"n4"}]}', 'hoat_dong')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_workflow_instances (id, code, template_id, applies_to, object_id, object_ref, status, current_step, started_at) VALUES
			('WF-01', 'WF-1026-01', 'WT-02', 'recruitment_request', 'RR-02', 'ĐXTĐ-1026-02', 'dang_chay', 1, '15/10/2026')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_workflow_steps (id, instance_id, step_order, step_type, approver_role, step_status) VALUES
			('WFS-01', 'WF-01', 1, 'approval', 'truong_phong', 'da_duyet'),
			('WFS-02', 'WF-01', 2, 'approval', 'giam_doc', 'cho_duyet')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_onboarding_checklists (id, employee_id, employee_name, task_name, is_completed, status) VALUES
			('OBC-01', 'EMP-01', 'Nguyễn Đức Trường', 'Ảnh thẻ / Ảnh cá nhân 3x4', TRUE, 'da_bo_sung'),
			('OBC-02', 'EMP-01', 'Nguyễn Đức Trường', 'Bản sao Giấy khai sinh', TRUE, 'da_bo_sung'),
			('OBC-03', 'EMP-04', 'Lê Hữu Thắng', 'Bản sao Bằng cấp / Chứng chỉ', FALSE, 'cho_bo_sung')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_employee_work_history (id, employee_id, employee_name, event_type, title, detail, changed_by, changed_at) VALUES
			('WH-01', 'EMP-01', 'Nguyễn Đức Trường', 'bo_nhiem', 'Bổ nhiệm Giám đốc điều hành', 'Bổ nhiệm chính thức', 'HĐQT Tập Đoàn TTC', '15/03/2018')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_attendance_logs (id, employee_id, employee_name, date, check_in_time, check_out_time, work_value, status_code, overtime_hours) VALUES
			('AL-001', 'EMP-01', 'Nguyễn Đức Trường', '28/10/2026', '07:25', '17:05', 1, 'P', 0.5),
			('AL-002', 'EMP-02', 'Nguyễn Văn Dũng', '28/10/2026', '06:30', '15:00', 1, 'P', 0.5),
			('AL-003', 'EMP-03', 'Trần Văn Kiên', '28/10/2026', '06:00', '15:30', 1, 'TC', 1.5),
			('AL-004', 'EMP-04', 'Lê Hữu Thắng', '28/10/2026', '07:00', '13:30', 0.9, 'P', 0)
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_shift_assignments (id, employee_id, employee_name, shift_code, shift_name, date, time_slot, status) VALUES
			('SA-01', 'EMP-01', 'Nguyễn Đức Trường', 'SHIFT-03', 'Ca Hành Chính', '28/10/2026', '07:30-16:30', 'da_phan_ca'),
			('SA-02', 'EMP-02', 'Nguyễn Văn Dũng', 'SHIFT-01', 'Ca Sáng', '28/10/2026', '06:00-14:00', 'da_phan_ca')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_holidays (id, holiday_name, holiday_date, date_type, status) VALUES
			('HOL-01', 'Lễ Quốc khánh 2/9', '02/09/2026', 'le', 'hoat_dong')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_payroll_periods (id, code, month, status, is_locked, start_date, end_date) VALUES
			('PR-1026', 'PL-1026', '10/2026', 'dang_tinh', FALSE, '01/10/2026', '31/10/2026')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_payroll_formula_columns (id, code, name, col_type, formula_expression, sort_order, status) VALUES
			('FML-01', 'A', 'Lương cơ bản', 'income', 'base_salary', 1, 'hoat_dong'),
			('FML-02', 'B', 'Lương ngày công', 'income', 'A * actual_days / standard_days', 2, 'hoat_dong'),
			('FML-03', 'C', 'Lương tăng ca', 'income', 'overtime_hours * (A / 26 / 8) * 1.5', 3, 'hoat_dong')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_salary_advances (id, code, employee_id, employee_name, amount, request_date, deduction_period, reason, status) VALUES
			('ADV-01', 'UL-1026-01', 'EMP-05', 'Trần Đình Trọng', 3000000, '25/10/2026', '11/2026', 'Ứng trước chi phí gia đình', 'cho_duyet')
			ON CONFLICT (id) DO NOTHING;
		`)
	}


	// Seed Alerts (Cảnh báo gian lận / lệch bì trạm cân)
	var alertCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM alerts").Scan(&alertCount)
	if alertCount == 0 {
		fmt.Println("🌱 Seeding Fraud / Tare-Mismatch Alerts...")
		alerts := []struct {
			ID, Title, BS, Note, Time, Date, Status, Severity, Phieu, Cam string
			BiDangKy, BiThucTe, LechBi                                        float64
		}{
			{"ALT-2026-001", "Lệch bì +380 kg (Vượt ngưỡng)", "19H-056.22", "Bùn đất dính dày dưới gầm thùng xe sau mưa moong", "13:59 28/10", "28/10/2026", "Đang xử lý hiện trường", "danger", "TK-20261028-002", "Trạm Cân Cổng 01 - Phú Thọ", 26450, 26830, 380},
			{"ALT-2026-002", "Lệch bì -140 kg (Trong dung sai nhưng bất thường)", "88H-042.27", "Trừ bì lệch do dư lượng thùng sau bốc hàng", "11:20 28/10", "28/10/2026", "Đã phê duyệt xử lý xong", "warning", "TK-20261028-001", "Trạm Cân Cổng 01 - Phú Thọ", 15420, 15280, 140},
			{"ALT-2026-003", "Nhận diện biển số Camera AI thấp", "90C-123.45", "Camera ANPR độ tin cậy 82% do mưa lớn", "14:56 28/10", "28/10/2026", "Chuyển Thanh tra mỏ", "warning", "TK-20261028-004", "Trạm Cân Cổng 02 - Phú Thọ", 18500, 18500, 0},
			{"ALT-2026-004", "Bảng giá đơn giá bất thường so với hợp đồng", "19C-098.76", "Nghi ngờ chỉnh sửa giá trọng tải đơn giá thủ công", "09:45 27/10", "27/10/2026", "Chờ xử lý", "info", "TK-20261027-009", "Trạm Cân Cổng 01 - Phú Thọ", 11200, 11200, 0},
		}

		for _, a := range alerts {
			Pool.Exec(ctx, `
				INSERT INTO alerts (id, title, bs, note, time, date, status, severity, phieu, cam, bi_dang_ky, bi_thuc_te, lech_bi)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
				ON CONFLICT (id) DO NOTHING
			`, a.ID, a.Title, a.BS, a.Note, a.Time, a.Date, a.Status, a.Severity, a.Phieu, a.Cam, a.BiDangKy, a.BiThucTe, a.LechBi)
		}
	}

	// Seed Print Templates (configurable voucher / report layouts)
	var printCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM print_templates").Scan(&printCount)
	if printCount == 0 {
		fmt.Println("🌱 Seeding Configurable Print Templates...")
		prints := []struct {
			ID, Code, Name, DocType, Size, Orientation, Description, Layout, Status string
			IsDefault                                                                 bool
		}{
			{"PRT-001", "TPL-TICKET-A5", "Phiếu cân xe A5 - 3 liên", "ticket", "A5", "portrait", "Phiếu in 3 liên (Khách, Kế toán, Bảo vệ) cho trạm cân", `{"page":{"size":"A5","orientation":"portrait","paddingMm":8},"elements":[{"id":"e1","kind":"text","text":"PHIẾU CÂN XE","x":0,"y":0,"w":100,"fontSize":16,"bold":true,"align":"center"},{"id":"e2","kind":"field","field":"bienSo","label":"Biển số","x":0,"y":24,"w":50,"fontSize":12,"bold":true},{"id":"e3","kind":"field","field":"matHang","label":"Mặt hàng","x":0,"y":38,"w":50,"fontSize":11},{"id":"e4","kind":"field","field":"klHang","label":"Khối lượng (kg)","x":0,"y":52,"w":50,"fontSize":12,"bold":true},{"id":"e5","kind":"field","field":"khachHang","label":"Bên mua","x":0,"y":66,"w":100,"fontSize":11},{"id":"e6","kind":"field","field":"laiXe","label":"Tài xế","x":0,"y":80,"w":50,"fontSize":11},{"id":"e7","kind":"field","field":"date","label":"Ngày","x":50,"y":80,"w":50,"fontSize":11}]}`, "active", true},
			{"PRT-002", "TPL-INVOICE-A4", "Hóa đơn VAT - Mẫu A4", "invoice", "A4", "portrait", "Hóa đơn giá trị gia tăng bán hàng", `{"page":{"size":"A4","orientation":"portrait","paddingMm":12},"elements":[{"id":"e1","kind":"text","text":"HÓA ĐƠN GIÁ TRỊ GIA TĂNG","x":0,"y":0,"w":100,"fontSize":16,"bold":true,"align":"center"},{"id":"e2","kind":"text","text":"Ký hiệu: TT/26E  •  Số: 0002281","x":0,"y":20,"w":100,"fontSize":11,"align":"center"},{"id":"e3","kind":"field","field":"benMua","label":"Người mua","x":0,"y":36,"w":100,"fontSize":11},{"id":"e4","kind":"table","field":"items","label":"Danh mục","x":0,"y":60,"w":100,"fontSize":10}]}`, "active", true},
			{"PRT-003", "TPL-ALERT-EN", "Biên bản xử lý cảnh báo", "alert", "A4", "portrait", "Biên bản vi phạm / cảnh báo lệch bì", `{"page":{"size":"A4","orientation":"portrait","paddingMm":10},"elements":[{"id":"e1","kind":"text","text":"BIÊN BẢN XỬ LÝ VI PHẠM","x":0,"y":0,"w":100,"fontSize":15,"bold":true,"align":"center"},{"id":"e2","kind":"field","field":"bienSo","label":"Biển số","x":0,"y":24,"w":50,"fontSize":12},{"id":"e3","kind":"field","field":"lechBi","label":"Lệch bì (kg)","x":50,"y":24,"w":50,"fontSize":12,"bold":true},{"id":"e4","kind":"field","field":"cause","label":"Nguyên nhân","x":0,"y":44,"w":100,"fontSize":11}]}`, "active", true},
			{"PRT-004", "TPL-WH-SLIP", "Phiếu xuất kho đá", "warehouse", "A5", "portrait", "Phiếu xuất bãi / nhập kho đá", `{"page":{"size":"A5","orientation":"portrait","paddingMm":8},"elements":[{"id":"e1","kind":"text","text":"PHIẾU XUẤT KHO","x":0,"y":0,"w":100,"fontSize":15,"bold":true,"align":"center"},{"id":"e2","kind":"field","field":"item","label":"Mặt hàng","x":0,"y":24,"w":60,"fontSize":12},{"id":"e3","kind":"field","field":"quantity","label":"Số lượng (tấn)","x":0,"y":40,"w":40,"fontSize":12,"bold":true},{"id":"e4","kind":"field","field":"customer","label":"Khách hàng","x":0,"y":56,"w":100,"fontSize":11}]}`, "active", false},
			{"PRT-005", "TPL-REPORT-LTR", "Báo cáo tổng hợp", "report", "A4", "landscape", "Báo cáo sản lượng / tổng hợp", `{"page":{"size":"A4","orientation":"landscape","paddingMm":12},"elements":[{"id":"e1","kind":"text","text":"KẾ HOẠCH & SẢN LƯỢNG","x":0,"y":0,"w":100,"fontSize":15,"bold":true,"align":"center"},{"id":"e2","kind":"field","field":"period","label":"Kỳ báo cáo","x":0,"y":24,"w":40,"fontSize":11},{"id":"e3","kind":"field","field":"actual","label":"Thực tế","x":40,"y":24,"w":30,"fontSize":11},{"id":"e4","kind":"field","field":"diff","label":"Chênh lệch","x":70,"y":24,"w":30,"fontSize":11,"bold":true}]}`, "active", false},
		}

		for _, p := range prints {
			// Layout JSON is stored on disk, not in the DB (metadata only).
			writePrintLayoutFile(p.ID, p.Layout)
			Pool.Exec(ctx, `
				INSERT INTO print_templates (id, code, name, doc_type, size, orientation, description, is_default, status)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				ON CONFLICT (id) DO NOTHING
			`, p.ID, p.Code, p.Name, p.DocType, p.Size, p.Orientation, p.Description, p.IsDefault, p.Status)
		}
	}

	// Backfill: move any legacy layout stored in the DB column out to disk so
	// templates remain disk-backed regardless of the seeding above.
	printBackfillLayoutToDisk(ctx)

	// Seed Quarry Core Module (Quarries, Areas, 3D Survey Cycles, Reconciliations)
	SeedQuarryModule()

	// Seed Google Maps GPS & Telematics Data
	var mapCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM google_map_entries").Scan(&mapCount)
	if mapCount == 0 {
		fmt.Println("🌱 Seeding Google Maps GPS Telematics & Vehicle Fleet...")
		mapEntries := []struct {
			Code string
			Data string
		}{
			{
				Code: "MAP-DRV-01",
				Data: `{"id":"DRV-01","plate":"88H-042.27","vehicleType":"Xe ben 4 chân Shacman X3000 (371HP)","driverName":"Trần Đình Trọng","driverPhone":"0913.678.902","lat":21.3280,"lng":105.3020,"speed":0,"maxSpeedLimit":60,"status":"idling_alert","statusLabel":"Dừng nổ máy 28p (Idling Alert)","engineStatus":"ON","engineRpm":780,"fuelLevelPercent":68,"fuelTankLiters":272,"idlingDurationMinutes":28,"idlingFuelWastedLiters":2.10,"idlingCostVnd":46200,"locationName":"Cổng trạm cân số 1 (Km 68 QL2)","destination":"Bãi tuyển đá 1x2 Nam Phương","cargo":"Đá 1x2 bê tông tiêu chuẩn","loadWeight":"36.40 Tấn","todayTrips":6,"todayKm":142,"updatedTime":"09:48:12 (Live GPS)","ticketId":"TK-20261028-001"}`,
			},
			{
				Code: "MAP-DRV-02",
				Data: `{"id":"DRV-02","plate":"19H-088.19","vehicleType":"Xe ben 3 chân Howo 371","driverName":"Nguyễn Văn Hùng","driverPhone":"0988.123.456","lat":21.3240,"lng":105.3180,"speed":42,"maxSpeedLimit":60,"status":"moving","statusLabel":"Đang di chuyển giao đá","engineStatus":"ON","engineRpm":1450,"fuelLevelPercent":82,"fuelTankLiters":328,"idlingDurationMinutes":0,"idlingFuelWastedLiters":0,"idlingCostVnd":0,"locationName":"QL2 Km 68 hướng Nút giao IC8","destination":"Gói thầu Cao tốc Cienco 4","cargo":"Đá 4x6 móng công trình","loadWeight":"32.10 Tấn","todayTrips":8,"todayKm":198,"updatedTime":"09:48:10 (Live GPS)","ticketId":"TK-20261028-003"}`,
			},
			{
				Code: "MAP-DRV-03",
				Data: `{"id":"DRV-03","plate":"19C-142.88","vehicleType":"Xe đầu kéo Hyundai HD270","driverName":"Lê Tuấn Anh","driverPhone":"0977.654.321","lat":21.3325,"lng":105.2950,"speed":0,"maxSpeedLimit":40,"status":"loading","statusLabel":"Đang xúc bốc đá tại Moong","engineStatus":"ON","engineRpm":820,"fuelLevelPercent":54,"fuelTankLiters":216,"idlingDurationMinutes":8,"idlingFuelWastedLiters":0.60,"idlingCostVnd":13200,"locationName":"Moong Khai Thác Tầng 2 (+100m)","destination":"Trạm Nghiền Sàng Số 01","cargo":"Đá hộc nổ mìn cỡ lớn","loadWeight":"28.50 Tấn","todayTrips":5,"todayKm":68,"updatedTime":"09:48:08 (Live GPS)","ticketId":"TK-20261028-002"}`,
			},
			{
				Code: "MAP-DRV-04",
				Data: `{"id":"DRV-04","plate":"29H-921.35","vehicleType":"Xe ben Howo V7X 4 chân","driverName":"Đặng Quốc Việt","driverPhone":"0904.889.900","lat":21.3285,"lng":105.3015,"speed":5,"maxSpeedLimit":20,"status":"weighing","statusLabel":"Đang vào bàn cân điện tử 80T","engineStatus":"ON","engineRpm":900,"fuelLevelPercent":91,"fuelTankLiters":364,"idlingDurationMinutes":3,"idlingFuelWastedLiters":0.22,"idlingCostVnd":4840,"locationName":"Cầu Cân Điện Tử 80 Tấn Số 02","destination":"Kho VLXD Sông Lô","cargo":"Đá mi sàng 0-5mm","loadWeight":"38.20 Tấn","todayTrips":7,"todayKm":165,"updatedTime":"09:48:05 (Live GPS)","ticketId":"TK-20261028-004"}`,
			},
			{
				Code: "MAP-DRV-05",
				Data: `{"id":"DRV-05","plate":"19H-056.22","vehicleType":"Xe ben 4 chân Dongfeng KC420","driverName":"Bùi Văn Long","driverPhone":"0912.334.455","lat":21.3300,"lng":105.2985,"speed":35,"maxSpeedLimit":50,"status":"moving","statusLabel":"Vận chuyển nội bộ mỏ","engineStatus":"ON","engineRpm":1280,"fuelLevelPercent":62,"fuelTankLiters":248,"idlingDurationMinutes":0,"idlingFuelWastedLiters":0,"idlingCostVnd":0,"locationName":"Đường Hào Vận Chuyển Tây Bắc","destination":"Bãi chứa đá thành phẩm","cargo":"Đá 1x2 sàng tuyển","loadWeight":"35.80 Tấn","todayTrips":9,"todayKm":212,"updatedTime":"09:48:15 (Live GPS)","ticketId":"TK-20261028-002"}`,
			},
			{
				Code: "MAP-DRV-06",
				Data: `{"id":"DRV-06","plate":"88C-099.14","vehicleType":"Xe bồn tưới đường Isuzu Giga (15m³)","driverName":"Phạm Văn Tuấn","driverPhone":"0966.778.899","lat":21.3260,"lng":105.2960,"speed":0,"maxSpeedLimit":40,"status":"parked","statusLabel":"Dừng tắt máy nghỉ ca","engineStatus":"OFF","engineRpm":0,"fuelLevelPercent":75,"fuelTankLiters":300,"idlingDurationMinutes":0,"idlingFuelWastedLiters":0,"idlingCostVnd":0,"locationName":"Khu Nhà Xưởng Bảo Dưỡng & Gara","destination":"Sẵn sàng tưới dập bụi","cargo":"Nước dập bụi môi trường","loadWeight":"0 Tấn (Xe bồn)","todayTrips":3,"todayKm":45,"updatedTime":"09:40:00","ticketId":"TK-20261028-001"}`,
			},
		}

		for _, m := range mapEntries {
			Pool.Exec(ctx, `
				INSERT INTO google_map_entries (code, data)
				VALUES ($1, $2::jsonb)
			`, m.Code, m.Data)
		}
	}

	// Seed Google Drive Files
	var driveCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM google_drive_files").Scan(&driveCount)
	if driveCount == 0 {
		fmt.Println("🌱 Seeding Google Drive Documents & Survey Files...")
		driveFiles := []struct {
			Code string
			Data string
		}{
			{Code: "DOC-01", Data: `{"code":"DOC-01","name":"Ho_so_kiem_dinh_tram_can_2025.pdf","fileName":"Ho_so_kiem_dinh_tram_can_2025.pdf","folder":"Trạm Cân / Kiểm Định","size":"2.4 MB","syncedAt":"28/10/2026 14:20","updated":"28/10/2026 14:20","status":"Đã đồng bộ","mimeType":"application/pdf","owner":"Kỹ thuật trạm cân"}`},
			{Code: "DOC-02", Data: `{"code":"DOC-02","name":"Hop_dong_cung_ung_Nam_Phuong.pdf","fileName":"Hop_dong_cung_ung_Nam_Phuong.pdf","folder":"Đối Tác / Hợp Đồng","size":"1.8 MB","syncedAt":"28/10/2026 10:15","updated":"28/10/2026 10:15","status":"Đã đồng bộ","mimeType":"application/pdf","owner":"Phòng Kinh doanh"}`},
			{Code: "DOC-03", Data: `{"code":"DOC-03","name":"Bien_ban_ban_giao_mo_da_TTC.pdf","fileName":"Bien_ban_ban_giao_mo_da_TTC.pdf","folder":"Pháp Lý / Bàn Giao","size":"4.1 MB","syncedAt":"27/10/2026 16:45","updated":"27/10/2026 16:45","status":"Đã đồng bộ","mimeType":"application/pdf","owner":"Ban Giám Đốc"}`},
			{Code: "DOC-04", Data: `{"code":"DOC-04","name":"Bang_ke_xuat_da_ca_chieu_2810.xlsx","fileName":"Bang_ke_xuat_da_ca_chieu_2810.xlsx","folder":"Báo Cáo / Xuất Đá","size":"850 KB","syncedAt":"28/10/2026 17:30","updated":"28/10/2026 17:30","status":"Đã đồng bộ","mimeType":"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet","owner":"Kế toán kho"}`},
			{Code: "DOC-05", Data: `{"code":"DOC-05","name":"Anh_camera_xe_19H05622_4goc.zip","fileName":"Anh_camera_xe_19H05622_4goc.zip","folder":"Camera AI / Chụp Bì","size":"15.6 MB","syncedAt":"28/10/2026 14:02","updated":"28/10/2026 14:02","status":"Đã đồng bộ","mimeType":"application/zip","owner":"Hệ thống Camera AI"}`},
		}
		for _, f := range driveFiles {
			Pool.Exec(ctx, `INSERT INTO google_drive_files (code, data) VALUES ($1, $2::jsonb)`, f.Code, f.Data)
		}
	}

	// Seed Google Gmail Notification Log
	var gmailCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM google_emails").Scan(&gmailCount)
	if gmailCount == 0 {
		fmt.Println("🌱 Seeding Google Gmail Notifications...")
		emails := []struct {
			Code string
			Data string
		}{
			{Code: "MAIL-01", Data: `{"code":"MAIL-01","subject":"[SỞ TN&MT PHÚ THỌ] Thông báo tiếp nhận Báo cáo thống kê sản lượng khai thác Quý 3/2026","sender":"sotnmt.phutho@chinhquyen.vn","recipient":"phaply@ttcquarry.vn","date":"28/10/2026 08:30","status":"Đã gửi","category":"Pháp lý Sở TN&MT"}`},
			{Code: "MAIL-02", Data: `{"code":"MAIL-02","subject":"[CẢNH BÁO AI] Phát hiện xe 19H-056.22 lệch bì +380kg tại Cầu cân số 01","sender":"system-alerts@smartscale.ai","recipient":"giamdoc.mo@ttcquarry.vn","date":"28/10/2026 13:59","status":"Tự động","category":"Cảnh báo An ninh"}`},
			{Code: "MAIL-03", Data: `{"code":"MAIL-03","subject":"Biên bản đối soát khối lượng xuất đá Công ty Nam Phương đợt 2 tháng 10","sender":"ketoan.ttc@gmail.com","recipient":"thanhtoan@namphuongcorp.vn","date":"27/10/2026 17:15","status":"Đã gửi","category":"Đối soát Khách hàng"}`},
		}
		for _, m := range emails {
			Pool.Exec(ctx, `INSERT INTO google_emails (code, data) VALUES ($1, $2::jsonb)`, m.Code, m.Data)
		}
	}

	// Seed Google Photos AI Imagery
	var photoCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM google_photo_entries").Scan(&photoCount)
	if photoCount == 0 {
		fmt.Println("🌱 Seeding Google Photos AI Media...")
		photos := []struct {
			Code string
			Data string
		}{
			{Code: "PHT-01", Data: `{"code":"PHT-01","title":"Ảnh chụp nhận diện biển số 19H-056.22 vào cân","url":"/assets/images/cam-plate-19H05622.jpg","source":"Camera AI LPR-01","capturedAt":"28/10/2026 13:58:45","vehicle":"19H-056.22","confidence":0.998}`},
			{Code: "PHT-02", Data: `{"code":"PHT-02","title":"Ảnh chụp thùng xe đầy đá 1x2 - Cân xuất","url":"/assets/images/cam-cargo-19H05622.jpg","source":"Camera Góc Trên TOP-02","capturedAt":"28/10/2026 14:02:10","vehicle":"19H-056.22","confidence":0.992}`},
			{Code: "PHT-03", Data: `{"code":"PHT-03","title":"Ảnh Drone khảo sát 3D Moong tầng khai thác +120m","url":"/assets/images/drone-survey-oct2026.jpg","source":"DJI Matrice 300 RTK","capturedAt":"25/10/2026 09:30:00","vehicle":"DRONE-01","confidence":1.0}`},
		}
		for _, p := range photos {
			Pool.Exec(ctx, `INSERT INTO google_photo_entries (code, data) VALUES ($1, $2::jsonb)`, p.Code, p.Data)
		}
	}

	fmt.Println("✅ Complete mining domain datasets successfully verified and seeded in PostgreSQL!")

	// ===== Seed new modules from meeting requirements =====

	// Seed vehicle trips (Camera AI)
	var tripCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM vehicle_trips").Scan(&tripCount)
	if tripCount == 0 {
		fmt.Println("🌱 Seeding Vehicle Trips (Camera AI)...")
		trips := []struct {
			Plate    string
			Driver   string
			Camera   string
			Dir      string
			InTime   string
			OutTime  string
			TripNum  int
			Eqty     float64
			Aqty     float64
			Conf     float64
			Status   string
		}{
			{"29E-380.15", "Nguyễn Văn Mạnh", "CAM-01", "inbound", "2026-08-27 06:30:00", "2026-08-27 06:45:00", 1, 30.5, 30.2, 0.98, "completed"},
			{"29E-380.15", "Nguyễn Văn Mạnh", "CAM-01", "outbound", "2026-08-27 07:00:00", "2026-08-27 07:15:00", 2, 30.0, 29.8, 0.97, "completed"},
			{"88H-042.27", "Trần Văn Bình", "CAM-02", "inbound", "2026-08-27 06:45:00", "2026-08-27 07:00:00", 1, 28.0, 27.6, 0.96, "completed"},
			{"88H-042.27", "Trần Văn Bình", "CAM-02", "outbound", "2026-08-27 07:30:00", "2026-08-27 07:45:00", 2, 28.0, 27.9, 0.99, "completed"},
			{"19H-056.22", "Lê Đức Hùng", "CAM-01", "inbound", "2026-08-27 08:00:00", "2026-08-27 08:15:00", 1, 32.0, 31.5, 0.95, "completed"},
		}
		for _, t := range trips {
			Pool.Exec(ctx, `INSERT INTO vehicle_trips (license_plate, driver_name, camera_id, direction, check_in_time, check_out_time, trip_number, estimated_quantity, actual_quantity, ai_confidence, status) VALUES ($1,$2,$3,$4,$5::timestamptz,$6::timestamptz,$7,$8,$9,$10,$11)`,
				t.Plate, t.Driver, t.Camera, t.Dir, t.InTime, t.OutTime, t.TripNum, t.Eqty, t.Aqty, t.Conf, t.Status)
		}
	}

	// Seed cost norms
	var cnCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM cost_norms").Scan(&cnCount)
	if cnCount == 0 {
		fmt.Println("🌱 Seeding Cost Norms...")
		norms := []struct {
			Name string
			Type string
			Cost float64
			Unit string
			Mat  string
		}{
			{"Nhiên liệu diesel", "fuel", 42000, "đ/lít", "Đá 1x2"},
			{"Chi phí nhân công", "labor", 15000, "đ/tấn", "Tất cả"},
			{"Khấu hao thiết bị", "depreciation", 8000, "đ/tấn", "Tất cả"},
			{"Chi phí vận chuyển", "transport", 12000, "đ/tấn/km", "Tất cả"},
			{"Chi phí sản xuất", "production", 5000, "đ/tấn", "Đá nghiền"},
		}
		for _, n := range norms {
			Pool.Exec(ctx, `INSERT INTO cost_norms (norm_name, norm_type, unit_cost, unit, material_type, effective_date, status) VALUES ($1,$2,$3,$4,$5,'2026-01-01','active')`,
				n.Name, n.Type, n.Cost, n.Unit, n.Mat)
		}
	}

	// Seed production costs
	var pcCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM production_costs").Scan(&pcCount)
	if pcCount == 0 {
		fmt.Println("🌱 Seeding Production Costs for Executive Command Center...")
		currPeriod := time.Now().Format("2006-01")
		prevPeriod := time.Now().AddDate(0, -1, 0).Format("2006-01")

		costs := []struct {
			CostType string
			Category string
			Norm     float64
			Actual   float64
			Period   string
			MineArea string
			Desc     string
		}{
			// Current month
			{"Sản xuất", "Sản xuất", 220000000, 245000000, currPeriod, "Moong Tầng 3 (+45m)", "Chi phí nổ mìn và vận hành trạm nghiền sàng"},
			{"Nhiên liệu", "Nhiên liệu", 180000000, 210000000, currPeriod, "Toàn mỏ", "Dầu diesel máy xúc PC450 và xe vận chuyển moong"},
			{"Nhân công", "Nhân công", 150000000, 155000000, currPeriod, "Toàn mỏ", "Tiền lương ca mỏ, thợ khoan, lái xe ben"},
			{"Vận chuyển", "Vận chuyển", 90000000, 98000000, currPeriod, "Tuyến nội bộ", "Cước trung chuyển đá hộc ra bãi sơ chế"},
			{"Khấu hao", "Khấu hao", 75000000, 75000000, currPeriod, "Trạm nghiền 01", "Khấu hao tài sản thiết bị nghiền sàng"},
			{"Khác", "Khác", 30000000, 32000000, currPeriod, "Văn phòng & An toàn", "Bảo hộ lao động, quan trắc bụi mỏ"},

			// Previous month (baseline)
			{"Sản xuất", "Sản xuất", 210000000, 215000000, prevPeriod, "Moong Tầng 3 (+45m)", "Chi phí nổ mìn tháng trước"},
			{"Nhiên liệu", "Nhiên liệu", 175000000, 185000000, prevPeriod, "Toàn mỏ", "Nhiên liệu diesel tháng trước"},
			{"Nhân công", "Nhân công", 150000000, 150000000, prevPeriod, "Toàn mỏ", "Tiền lương tháng trước"},
			{"Vận chuyển", "Vận chuyển", 85000000, 88000000, prevPeriod, "Tuyến nội bộ", "Chi phí vận chuyển tháng trước"},
			{"Khấu hao", "Khấu hao", 75000000, 75000000, prevPeriod, "Trạm nghiền 01", "Khấu hao tháng trước"},
			{"Khác", "Khác", 28000000, 29000000, prevPeriod, "Văn phòng", "Chi phí khác tháng trước"},
		}

		for _, c := range costs {
			Pool.Exec(ctx, `
				INSERT INTO production_costs (cost_type, cost_category, norm_value, norm_unit, actual_value, actual_unit, period, mine_area, description)
				VALUES ($1, $2, $3, 'VNĐ', $4, 'VNĐ', $5, $6, $7)
			`, c.CostType, c.Category, c.Norm, c.Actual, c.Period, c.MineArea, c.Desc)
		}
	}

	// Seed tax records
	var taxCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM tax_records").Scan(&taxCount)
	if taxCount == 0 {
		fmt.Println("🌱 Seeding Tax Records...")
		taxes := []struct {
			Type   string
			Code   string
			Period string
			TaxAmt float64
			Paid   float64
			Due    string
			Status string
		}{
			{"VAT", "VAT-2026-08", "08/2026", 85000000, 0, "2026-09-20", "pending"},
			{"Thuế tài nguyên", "TR-2026-08", "08/2026", 120000000, 120000000, "2026-09-15", "paid"},
			{"Thuế thu nhập DN", "CIT-2026-Q3", "Q3/2026", 250000000, 0, "2026-10-30", "pending"},
		}
		for _, t := range taxes {
			Pool.Exec(ctx, `INSERT INTO tax_records (tax_type, tax_code, period, tax_amount, paid_amount, due_date, status, authority) VALUES ($1,$2,$3,$4,$5,$6,$7,'Cục Thuế Phú Thọ')`,
				t.Type, t.Code, t.Period, t.TaxAmt, t.Paid, t.Due, t.Status)
		}
	}

	// Seed risk alerts
	var riskCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM risk_alerts").Scan(&riskCount)
	if riskCount == 0 {
		fmt.Println("🌱 Seeding Risk Alerts...")
		alerts := []struct {
			Type    string
			Sev     string
			Title   string
			Desc    string
			Module  string
			Status  string
		}{
			{"vehicle", "high", "Xe 88H-042.27 đi sai tuyến", "Phát hiện xe ra khỏi vùng geofence allowed", "gps", "open"},
			{"permit", "medium", "Giấy phép mỏ sắp hết hạn", "Giấy phép 189/GP-BTNMT hết hạn sau 45 ngày", "permits", "open"},
			{"production", "low", "Sản lượng tháng 8 dưới kế hoạch", "Đạt 78% kế hoạch tháng", "mining-plan", "open"},
			{"hr", "medium", "Thiếu 3 nhân sự ca sáng", "Ca sáng thiếu 3 người so với lịch phân công", "hr", "open"},
		}
		for _, a := range alerts {
			Pool.Exec(ctx, `INSERT INTO risk_alerts (alert_type, severity, title, description, module_source, status) VALUES ($1,$2,$3,$4,$5,$6)`,
				a.Type, a.Sev, a.Title, a.Desc, a.Module, a.Status)
		}
	}

	// Seed delegations
	var delCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM delegations").Scan(&delCount)
	if delCount == 0 {
		fmt.Println("🌱 Seeding Delegations...")
		dels := []struct {
			Delegator string
			Delegate  string
			PermType  string
			Scope     string
			From      string
			To        string
			Status    string
		}{
			{"Giám đốc Nguyễn Văn A", "Phó Giám đốc Trần Văn B", "Ký hợp đồng", "Hợp đồng mua bán đá", "2026-08-01", "2026-08-31", "active"},
			{"Kế toán trưởng Lê Văn C", "Nhân viên kế toán Phạm Thị D", "Xác nhận hóa đơn", "Hóa đơn VAT dưới 500 triệu", "2026-08-15", "2026-09-15", "active"},
		}
		for _, d := range dels {
			Pool.Exec(ctx, `INSERT INTO delegations (delegator_name, delegate_name, permission_type, scope, start_date, end_date, status) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
				d.Delegator, d.Delegate, d.PermType, d.Scope, d.From, d.To, d.Status)
		}
	}

	fmt.Println("✅ New module seed data completed!")
}

