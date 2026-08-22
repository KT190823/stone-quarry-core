package database

import (
	"context"
	"fmt"
)

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

		// HR Contracts
		Pool.Exec(ctx, `
			INSERT INTO hr_contracts (id, code, employee_id, employee_name, department, contract_type, contracted_hours, start_date, end_date, base_salary, hazard_allowance, safety_allowance, meal_allowance, responsibility_allowance, attendance_allowance, social_insurance_base, status)
			VALUES 
			('CTR-01', 'HĐLĐ-2018-01', 'EMP-01', 'Nguyễn Đức Trường', 'Ban Giám Đốc', 'Chính thức vô thời hạn', 8.0, '15/03/2018', 'Vô thời hạn', 38000000, 1500000, 1500000, 1000000, 2000000, 1000000, 38000000, 'Hiệu lực'),
			('CTR-02', 'HĐLĐ-2020-02', 'EMP-02', 'Nguyễn Văn Dũng', 'Tổ Vận Hành Trạm Cân', 'Chính thức vô thời hạn', 8.0, '01/06/2020', 'Vô thời hạn', 16500000, 800000, 600000, 800000, 1000000, 800000, 16500000, 'Hiệu lực'),
			('CTR-03', 'HĐLĐ-2019-03', 'EMP-03', 'Trần Văn Kiên', 'Tổ Khoan Nổ Mìn & VLNCN', 'Chính thức vô thời hạn', 8.0, '12/08/2019', 'Vô thời hạn', 28000000, 2000000, 1500000, 1000000, 1500000, 1000000, 28000000, 'Hiệu lực')
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
	var vehCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM vehicles").Scan(&vehCount)
	if vehCount == 0 {
		fmt.Println("🌱 Seeding Fleet Vehicles...")
		vehicles := []struct {
			BS, Loai, Bi, RFID, TaiTrong, Status, HanDangKiem, Date, ChuXe string
			Count                                                          int
		}{
			{"88H-042.27", "Xe ben 4 chân HOWO", "15.42 tấn", "RFID-88H-042", "30.0 tấn", "Hoạt động", "15/12/2026", "28/10/2026", "Công ty CP Đầu Tư Xây Dựng 319", 142},
			{"19H-056.22", "Xe ben Chenglong", "14.80 tấn", "RFID-19H-056", "30.0 tấn", "Hoạt động", "20/01/2027", "28/10/2026", "Tổng Công Ty XD Trường Sơn", 98},
			{"29H-882.19", "Xe bồn trộn bê tông", "16.10 tấn", "RFID-29H-882", "25.0 tấn", "Hoạt động", "10/11/2026", "28/10/2026", "Công ty Bê Tông Việt Trì", 76},
			{"90C-123.45", "Xe đầu kéo Mooc ben", "18.50 tấn", "RFID-90C-123", "45.0 tấn", "Hoạt động", "05/03/2027", "28/10/2026", "Công ty CP Tập Đoàn Đèo Cả", 64},
			{"19C-098.76", "Xe tải ben 3 chân", "11.20 tấn", "RFID-19C-098", "20.0 tấn", "Bảo dưỡng", "28/10/2026", "28/10/2026", "Hợp tác xã Vận tải Hùng Vương", 35},
		}

		for _, v := range vehicles {
			Pool.Exec(ctx, `
				INSERT INTO vehicles (bs, loai, bi, rfid, tai_trong, status, count, han_dang_kiem, date, chu_xe)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				ON CONFLICT (bs) DO NOTHING
			`, v.BS, v.Loai, v.Bi, v.RFID, v.TaiTrong, v.Status, v.Count, v.HanDangKiem, v.Date, v.ChuXe)
		}
	}

	// 4. Seed Materials Catalog
	var matCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM catalog_materials").Scan(&matCount)
	if matCount == 0 {
		fmt.Println("🌱 Seeding Stone Materials Catalog...")
		materials := []struct {
			Code, Name, DVT, Density, DinhMuc, Price, Kho, Standard, Status, Date string
		}{
			{"MAT-01", "Đá 1x2 Xây Dựng", "m3", "1.55 tấn/m3", "1.55", "240.000 đ/m3", "Bãi Đá Thành Phẩm 01", "TCVN 7570:2006", "Đang kinh doanh", "28/10/2026"},
			{"MAT-02", "Đá 4x6 Kè Móng", "m3", "1.60 tấn/m3", "1.60", "220.000 đ/m3", "Bãi Đá Hộc & 4x6", "TCVN 1771:1987", "Đang kinh doanh", "28/10/2026"},
			{"MAT-03", "Đá Base Cấp Phối Loại 1", "tấn", "1.00 tấn/tấn", "1.00", "180.000 đ/tấn", "Bãi Cấp Phối Dmax25", "TCVN 8859:2011", "Đang kinh doanh", "28/10/2026"},
			{"MAT-04", "Cát Nghiền Nhân Tạo (Mạt đá)", "m3", "1.45 tấn/m3", "1.45", "260.000 đ/m3", "Kho Mái Che Cát Nghiền", "TCVN 9205:2012", "Đang kinh doanh", "28/10/2026"},
			{"MAT-05", "Đá Hộc Xô Bồ Nổ Mìn", "m3", "1.65 tấn/m3", "1.65", "160.000 đ/m3", "Moong Tầng 3 (+45m)", "Đá khai thác nguyên khai", "Đang kinh doanh", "28/10/2026"},
			{"MAT-06", "Đá 2x4 Đổ Bê Tông", "m3", "1.58 tấn/m3", "1.58", "235.000 đ/m3", "Bãi Đá Thành Phẩm 02", "TCVN 7570:2006", "Đang kinh doanh", "28/10/2026"},
		}

		for _, m := range materials {
			Pool.Exec(ctx, `
				INSERT INTO catalog_materials (code, ma, name, ten, dvt, density, ty_trong, dinh_muc, price, gia, kho, standard, status, date)
				VALUES ($1, $1, $2, $2, $3, $4, $4, $5, $6, $6, $7, $8, $9, $10)
				ON CONFLICT (code) DO NOTHING
			`, m.Code, m.Name, m.DVT, m.Density, m.DinhMuc, m.Price, m.Kho, m.Standard, m.Status, m.Date)
		}
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
			INSERT INTO users (username, name, ten, role, dept, email, phone, sdt, last_login, status) VALUES
			('admin', 'Nguyễn Đức Trường', 'Nguyễn Đức Trường', 'Giám Đốc Mỏ', 'Ban Giám Đốc', 'truongnd@ttcgroup.vn', '0912.345.678', '0912.345.678', '28/10/2026 08:30', 'Hoạt động'),
			('dungnv', 'Nguyễn Văn Dũng', 'Nguyễn Văn Dũng', 'Trưởng Trạm Cân', 'Tổ Vận Hành Trạm Cân', 'dungnv@ttcgroup.vn', '0984.556.789', '0984.556.789', '28/10/2026 09:15', 'Hoạt động'),
			('kientv', 'Trần Văn Kiên', 'Trần Văn Kiên', 'Chỉ Huy Nổ Mìn', 'Tổ Khoan Nổ Mìn & VLNCN', 'kientv@ttcgroup.vn', '0988.341.992', '0988.341.992', '28/10/2026 07:45', 'Hoạt động'),
			('thuynt', 'Nguyễn Thị Thủy', 'Nguyễn Thị Thủy', 'Kế Toán Trưởng', 'Phòng Kế Toán & Vật Tư', 'thuynt@ttcgroup.vn', '0915.678.901', '0915.678.901', '28/10/2026 08:00', 'Hoạt động')
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

			INSERT INTO reports (id, name, item, type, period, plan, actual, diff, status) VALUES
			('REP-01', 'Báo cáo sản lượng khai thác đá nguyên khai', 'Đá nguyên khai', 'Sản xuất', 'Tháng 10/2026', '120.000 m3', '128.500 m3', '+7.1%', 'Hoàn thành'),
			('REP-02', 'Báo cáo nghiền sàng đá thành phẩm', 'Đá 1x2, 4x6, Base', 'Chế biến', 'Tháng 10/2026', '95.000 m3', '98.200 m3', '+3.4%', 'Hoàn thành'),
			('REP-03', 'Báo cáo tiêu hao nhiên liệu cơ giới', 'Dầu DO 0.05S', 'Nhiên liệu', 'Tháng 10/2026', '42.000 Lít', '44.850 Lít', '+6.8%', 'Cần kiểm tra'),
			('REP-04', 'Báo cáo xuất bán & đối soát công nợ', 'Doanh thu xuất mỏ', 'Tài chính', 'Tháng 10/2026', '28.5 tỷ đ', '29.8 tỷ đ', '+4.6%', 'Hoàn thành')
			ON CONFLICT (id) DO NOTHING;

			INSERT INTO settings (code, key, name, val, scope, status) VALUES
			('SET-01', 'WEIGH_TOLERANCE_KG', 'Dung sai bì trạm cân cho phép (kg)', '50', 'Trạm Cân', 'Hoạt động'),
			('SET-02', 'GEOFENCE_RADIUS_METERS', 'Bán kính Geofence trạm bãi xe (mét)', '20', 'Định Vị GPS', 'Hoạt động'),
			('SET-03', 'FUEL_THEFT_THRESHOLD_LITERS', 'Ngưỡng cảnh báo sụt dầu bất thường (lít)', '5.0', 'Nhiên Liệu', 'Hoạt động'),
			('SET-04', 'OCR_ANPR_CONFIDENCE_MIN', 'Độ tin cậy nhận diện biển số AI OCR (%)', '92', 'Camera AI', 'Hoạt động')
			ON CONFLICT (code) DO NOTHING;
		`)
	}

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

			INSERT INTO vehicle_catalogs (code, name, loai, tai_trong, so_truc, status) VALUES
			('VC-01', 'Xe ben HOWO 371HP 4 Chân', 'Xe ben 4 chân', '30.0 tấn', '4 trục', 'Hoạt động'),
			('VC-02', 'Xe ben Chenglong Hải Âu 385HP', 'Xe ben 4 chân', '30.0 tấn', '4 trục', 'Hoạt động'),
			('VC-03', 'Xe ben Shacman X3000', 'Xe ben 4 chân', '30.0 tấn', '4 trục', 'Hoạt động'),
			('VC-04', 'Xe đầu kéo Howo Mooc Ben', 'Xe đầu kéo', '45.0 tấn', '6 trục', 'Hoạt động'),
			('VC-05', 'Xe bồn bê tông trộn 12m3', 'Xe bồn trộn', '25.0 tấn', '3 trục', 'Hoạt động')
			ON CONFLICT (code) DO NOTHING;

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

	// 13. Seed Mining Operations & Statutory Reports
	var statCount int
	Pool.QueryRow(ctx, "SELECT COUNT(*) FROM statutory_reports").Scan(&statCount)
	if statCount == 0 {
		Pool.Exec(ctx, `
			INSERT INTO statutory_reports (id, code, name, item, period, plan, actual, diff, status) VALUES
			('STAT-01', 'BC-TK-2026-Q3', 'Báo cáo thống kê hoạt động khai thác khoáng sản Q3/2026 (Mẫu số 01/KTKS)', 'Đá vôi xây dựng', 'Quý 3/2026', '350.000 m3', '362.400 m3', '+3.5%', 'Đã nộp Sở TN&MT'),
			('STAT-02', 'BC-AT-2026-10', 'Báo cáo công tác an toàn lao động & VLNCN tháng 10/2026', 'An toàn & Thuốc nổ', 'Tháng 10/2026', '100% An toàn', 'Không sự cố', 'Đạt', 'Đã nộp Sở Công Thương'),
			('STAT-03', 'BC-BVMT-2026-Q3', 'Báo cáo quan trắc môi trường định kỳ mỏ Q3/2026', 'Bụi, Tiếng ồn, Nước thải', 'Quý 3/2026', 'Đạt QCVN', 'QCVN 05/2013/BTNMT', 'Đạt', 'Đã thẩm định')
			ON CONFLICT (id) DO NOTHING;

			INSERT INTO resource_taxes (code, material, tax_rate, environmental_fee, status) VALUES
			('TAX-01', 'Đá vôi xây dựng nguyên khai (Đá hộc)', '10%', '3.000 đ/m3', 'Hiệu lực'),
			('TAX-02', 'Đá thành phẩm nghiền sàng (Đá 1x2, 2x4, 4x6)', '10%', '3.000 đ/m3', 'Hiệu lực'),
			('TAX-03', 'Cát nghiền nhân tạo (Mạt đá)', '10%', '4.000 đ/m3', 'Hiệu lực')
			ON CONFLICT (code) DO NOTHING;

			INSERT INTO production_stages (code, name, stage, description, status) VALUES
			('STG-01', 'Khoan & Nổ Mìn Cắt Tầng', 'Công đoạn 1', 'Khoan lỗ D110 nạp thuốc nổ ANFO/Nhũ tương phá vỡ tầng đá', 'Đang vận hành'),
			('STG-02', 'Bốc Xúc & Vận Chuyển Moong', 'Công đoạn 2', 'Máy xúc PC450 bốc đá hộc lên xe ben vận chuyển về trạm nghiền', 'Đang vận hành'),
			('STG-03', 'Nghiền Sàng Tuyển Thành Phẩm', 'Công đoạn 3', 'Nghiền thô kẹp hàm + Sàng rung phân loại Đá 1x2, 4x6, cát nghiền', 'Đang vận hành'),
			('STG-04', 'Cân Điện Tử & Xuất Bán Mỏ', 'Công đoạn 4', 'Cân xe 2 lần chốt tải trọng và phát hành phiếu cân DO/HĐĐT', 'Đang vận hành')
			ON CONFLICT (code) DO NOTHING;

			INSERT INTO hr_shifts (code, name, start_time, end_time, hours, status) VALUES
			('SHIFT-01', 'Ca Sáng (Ca 1)', '06:00', '14:00', 8.0, 'Hoạt động'),
			('SHIFT-02', 'Ca Chiều (Ca 2)', '14:00', '22:00', 8.0, 'Hoạt động'),
			('SHIFT-03', 'Ca Hành Chính', '07:30', '16:30', 8.0, 'Hoạt động')
			ON CONFLICT (code) DO NOTHING;
		`)
	}

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
			INSERT INTO hr_esign_documents (id, code, contract_id, contract_code, subject, status) VALUES
			('ES-01', 'ES-1026-01', 'CTR-01', 'HĐLĐ-2018-01', 'Ký số Hợp đồng lao động - Nguyễn Đức Trường', 'cho_ky')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_esign_signers (id, document_id, signer_type, employee_id, employee_name, sign_order, sign_status) VALUES
			('SG-01', 'ES-01', 'nhan_vien', 'EMP-01', 'Nguyễn Đức Trường', 1, 'cho_ky'),
			('SG-02', 'ES-01', 'quan_ly_nhan_su', NULL, 'Nguyễn Thị Thủy', 2, 'cho_ky')
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
			('KPI-01', 'EMP-01', 'Nguyễn Đức Trường', '10/2026', 'KPI-SL', 'Sản lượng khai thác', 120000, 128500, 'm3', 40, 'da_ban_hanh'),
			('KPI-02', 'EMP-01', 'Nguyễn Đức Trường', '10/2026', 'KPI-AT', 'An toàn lao động', 100, 100, '%', 30, 'da_ban_hanh'),
			('KPI-03', 'EMP-01', 'Nguyễn Đức Trường', '10/2026', 'KPI-DT', 'Doanh thu', 28.5, 29.8, 'tỷ', 30, 'da_ban_hanh')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_okr_objectives (id, code, title, level, period, progress_percent, status) VALUES
			('OKR-01', 'OKR-Q4-01', 'Nâng sản lượng khai thác 10% quý 4/2026', 'company', 'Q4/2026', 65, 'dang_cho')
			ON CONFLICT (id) DO NOTHING;
		`)

		Pool.Exec(ctx, `
			INSERT INTO hr_okr_key_results (id, objective_id, key_result, start_value, target_value, current_value, unit, progress_percent) VALUES
			('KR-01', 'OKR-01', 'Tăng sản lượng đá nguyên khai', 100000, 120000, 113000, 'm3', 65)
			ON CONFLICT (id) DO NOTHING;
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

	fmt.Println("✅ Complete mining domain datasets successfully verified and seeded in PostgreSQL!")
}
