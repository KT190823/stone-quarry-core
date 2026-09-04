package intelligence

import (
	"fmt"
	"strings"
)

type CopilotCause struct {
	Factor    string  `json:"factor"`
	ImpactPct float64 `json:"impact_pct"`
	Evidence  string  `json:"evidence"`
}

type CopilotRecommendation struct {
	Action      string  `json:"action"`
	ActionType  string  `json:"action_type"` // navigate | show_data | export
	TargetURL   *string `json:"target_url,omitempty"`
	TargetLabel *string `json:"target_label,omitempty"`
}

type CopilotResponse struct {
	Answer            string                  `json:"answer"`
	RootCauses        []CopilotCause          `json:"root_causes,omitempty"`
	Recommendations   []CopilotRecommendation `json:"recommendations,omitempty"`
	FollowUpQuestions []string                `json:"follow_up_questions,omitempty"`
}

type CopilotContext struct {
	Period string `json:"period"`
	Domain string `json:"domain"`
}

func Analyze(question string, ctx *CopilotContext) CopilotResponse {
	q := strings.ToLower(question)

	if strings.Contains(q, "lệch bì") || strings.Contains(q, "bì") || strings.Contains(q, "gian lận") {
		return analyzeTareAlerts()
	}

	if strings.Contains(q, "phiếu cân") || strings.Contains(q, "trạm cân") || strings.Contains(q, "chuyến xe") || strings.Contains(q, "xuất bán") {
		return analyzeTicketOperations()
	}

	if strings.Contains(q, "tồn kho") || strings.Contains(q, "kho") || strings.Contains(q, "bãi") {
		return analyzeInventoryStock()
	}

	if strings.Contains(q, "88h") || strings.Contains(q, "19h") || strings.Contains(q, "29c") {
		return analyzeSpecificTruck(q)
	}

	if strings.Contains(q, "lợi nhuận") && (strings.Contains(q, "giảm") || strings.Contains(q, "tại sao") || strings.Contains(q, "kém") || strings.Contains(q, "thấp")) {
		return analyzeProfitDecline()
	}

	if (strings.Contains(q, "chi phí") || strings.Contains(q, "giá")) && (strings.Contains(q, "tăng") || strings.Contains(q, "tại sao") || strings.Contains(q, "nhiều")) {
		return analyzeCostIncrease()
	}

	if (strings.Contains(q, "xe") || strings.Contains(q, "ben") || strings.Contains(q, "vận chuyển")) && (strings.Contains(q, "dầu") || strings.Contains(q, "ăn") || strings.Contains(q, "tiền") || strings.Contains(q, "hiệu quả")) {
		return analyzeVehicleEfficiency()
	}

	if strings.Contains(q, "sản lượng") && (strings.Contains(q, "giảm") || strings.Contains(q, "thấp") || strings.Contains(q, "kém")) {
		return analyzeProductionDrop()
	}

	if strings.Contains(q, "khách hàng") || strings.Contains(q, "đối tác") || strings.Contains(q, "công nợ") {
		return analyzeCustomerProfitability()
	}

	if strings.Contains(q, "dự báo") || strings.Contains(q, "tháng tới") || strings.Contains(q, "kịch bản") {
		return analyzeForecast()
	}

	return generalRecommendation()
}

func analyzeProfitDecline() CopilotResponse {
	targetUrlSim := "/simulator"
	labelSim := "Mở What-If Simulator"
	targetUrlFleet := "/profitability?tab=vehicles"
	labelFleet := "Kiểm tra Đội Xe"

	return CopilotResponse{
		Answer: "Lợi nhuận mỏ tháng này ước tính 635 triệu VNĐ, giảm 5.8% so với cùng kỳ tháng trước (674 triệu VNĐ). Nguyên nhân xuất phát từ giá nhiên liệu dầu diesel tăng và năng suất máy nghiền 01 sụt giảm trong tuần bảo trì.",
		RootCauses: []CopilotCause{
			{
				Factor:    "Chi phí nhiên liệu diesel tăng vọt",
				ImpactPct: 12.5,
				Evidence:  "Giá dầu tăng điều chỉnh và xe ben chạy cự ly moong phụ phát sinh thêm 1.2km.",
			},
			{
				Factor:    "Sản lượng khai thác tầng 3 giảm",
				ImpactPct: -14.9,
				Evidence:  "Máy nghiền 01 dừng máy 14 giờ để thay tấm lót hàm và lưới sàng.",
			},
			{
				Factor:    "Xe 88H-042.27 tiêu hao dầu cao bất thường",
				ImpactPct: 6.2,
				Evidence:  "Mức tiêu hao 38.5 L/100km, cao hơn định mức 31%.",
			},
		},
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Kiểm tra bảo dưỡng khẩn cấp xe ben 88H-042.27",
				ActionType:  "navigate",
				TargetURL:   &targetUrlFleet,
				TargetLabel: &labelFleet,
			},
			{
				Action:      "Chạy mô phỏng What-If tăng giá bán bù đắp chi phí nhiên liệu",
				ActionType:  "navigate",
				TargetURL:   &targetUrlSim,
				TargetLabel: &labelSim,
			},
		},
		FollowUpQuestions: []string{
			"Xe nào đang tiêu hao nhiên liệu cao nhất?",
			"Nên tăng giá bán bao nhiêu để bù đắp chi phí dầu?",
			"Chi phí sản xuất 1 tấn đá 1x2 hiện tại là bao nhiêu?",
		},
	}
}

func analyzeCostIncrease() CopilotResponse {
	targetUrlCost := "/costIntel"
	labelCost := "Xem Chi Tiết Chi Phí"

	return CopilotResponse{
		Answer: "Chi phí sản xuất toàn mỏ tăng 8.2% trong tháng này, chủ yếu do giá dầu diesel tăng 12.5% và chi phí trung chuyển đá hộc. Định mức giá thành đá 1x2 hiện tại là 192.000 đ/tấn (vượt định mức 6.2%).",
		RootCauses: []CopilotCause{
			{
				Factor:    "Nhiên liệu Diesel (chiếm 23.4% tổng chi phí)",
				ImpactPct: 12.5,
				Evidence:  "Chi phí đạt 210 triệu VNĐ so với định mức 180 triệu VNĐ.",
			},
			{
				Factor:    "Vận chuyển cự ly dài nội bộ",
				ImpactPct: 4.8,
				Evidence:  "Cung đường moong tầng 3 ra bãi đá thành phẩm kéo dài do mở rộng gương khai thác.",
			},
		},
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Xem phân rã chi phí chi tiết theo 5 danh mục",
				ActionType:  "navigate",
				TargetURL:   &targetUrlCost,
				TargetLabel: &labelCost,
			},
		},
		FollowUpQuestions: []string{
			"Giá bán đá 4x6 có bù đắp được giá thành không?",
			"Chi phí nhân công tháng này có tăng không?",
		},
	}
}

func analyzeVehicleEfficiency() CopilotResponse {
	targetUrl := "/profitability?tab=vehicles"
	label := "Xếp Hạng Đội Xe"

	return CopilotResponse{
		Answer: "Đội xe hiện có 5 xe chính hoạt động. Đa số đạt hiệu quả tốt (biên lợi nhuận trên 55%), ngoại trừ 2 xe cần can thiệp: Xe 88H-042.27 (margin 42.1%, ăn dầu) và Xe 19C-098.76 (margin 23.4%, xe cũ đang bảo dưỡng).",
		RootCauses: []CopilotCause{
			{
				Factor:    "Xe 88H-042.27 chạy nhiều chuyến nhất nhưng chi phí dầu chiếm 38.5%",
				ImpactPct: 18.0,
				Evidence:  "Đã chạy 142 chuyến, doanh thu 1.13 tỷ nhưng tốn 435 triệu tiền dầu.",
			},
			{
				Factor:    "Xe 19C-098.76 chi phí sửa chữa cao",
				ImpactPct: 24.6,
				Evidence:  "Chi phí bảo trì 38 triệu trên 154 triệu doanh thu (chiếm gần 25%).",
			},
		},
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Mở danh sách xếp hạng sinh lời đội xe để xem chi tiết PnL",
				ActionType:  "navigate",
				TargetURL:   &targetUrl,
				TargetLabel: &label,
			},
		},
		FollowUpQuestions: []string{
			"Lái xe nào điều khiển xe 88H-042.27?",
			"Xe 19H-056.22 chạy có hiệu quả không?",
		},
	}
}

func analyzeProductionDrop() CopilotResponse {
	return CopilotResponse{
		Answer: "Sản lượng 7 ngày qua đạt 2.680 tấn, giảm 14.9% so với định mức tuần (3.150 tấn). Trạm nghiền sàng 01 đã hoạt động trở lại bình thường từ hôm qua, dự kiến sản lượng sẽ phục hồi 100% trong 48 giờ tới.",
		RootCauses: []CopilotCause{
			{
				Factor:    "Bảo trì định kỳ cụm sàng rung & máy nghiền nón",
				ImpactPct: -14.9,
				Evidence:  "Dừng máy kế hoạch 14 giờ để thay tấm lót hàm.",
			},
		},
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Tăng cường 1 ca phụ vào buổi tối để bù đắp sản lượng thiếu hụt",
				ActionType:  "show_data",
				TargetURL:   nil,
				TargetLabel: nil,
			},
		},
		FollowUpQuestions: []string{
			"Tồn kho đá 1x2 còn đủ bán cho các dự án không?",
			"Hôm nay sản lượng trạm cân đã đạt bao nhiêu tấn?",
		},
	}
}

func analyzeCustomerProfitability() CopilotResponse {
	targetUrl := "/profitability?tab=customers"
	label := "Xem Phân Tích Khách Hàng"

	return CopilotResponse{
		Answer: "Tổng Công Ty XD Trường Sơn là khách hàng lớn nhất (8.2 tỷ VNĐ) nhưng biên lợi nhuận chỉ đạt 14.5% và đang có 540 triệu công nợ quá hạn. Khách hàng Công ty 319 và Bê Tông Việt Trì mang lại biên lợi nhuận cao nhất (>22%).",
		RootCauses: []CopilotCause{
			{
				Factor:    "Chiết khấu và vận chuyển đường dài cho Trường Sơn",
				ImpactPct: 14.6,
				Evidence:  "Chiết khấu thương mại 120 triệu + chi phí tài chính công nợ 38 triệu.",
			},
		},
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Xem danh sách đối tác và chính sách chiết khấu",
				ActionType:  "navigate",
				TargetURL:   &targetUrl,
				TargetLabel: &label,
			},
		},
		FollowUpQuestions: []string{
			"Hạn mức công nợ của Trường Sơn là bao nhiêu?",
			"Khách hàng nào thanh toán nhanh nhất?",
		},
	}
}

func analyzeForecast() CopilotResponse {
	targetUrl := "/simulator"
	label := "Mở What-If Simulator"

	return CopilotResponse{
		Answer: "Dự báo tháng tới nhu cầu đá cấp phối Base và đá 1x2 sẽ tăng mạnh do dự án đường cao tốc đẩy nhanh tiến độ thi công. Nếu nâng sản lượng thêm 10% và giữ ổn định giá, lợi nhuận dự kiến tăng thêm 210 triệu VNĐ/tháng.",
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Thử nghiệm mô phỏng tăng sản lượng và điều chỉnh giá bán",
				ActionType:  "navigate",
				TargetURL:   &targetUrl,
				TargetLabel: &label,
			},
		},
		FollowUpQuestions: []string{
			"Nếu tăng giá bán 20k/tấn thì lợi nhuận thế nào?",
			"Đội xe hiện tại có đủ năng lực vận chuyển thêm 10% không?",
		},
	}
}

func generalRecommendation() CopilotResponse {
	targetUrl := "/chi-huy-dieu-hanh/mo-phong"
	label := "Khám Phá Kịch Bản What-If"

	return CopilotResponse{
		Answer: "Xin chào! Tôi là Trợ Lý AI Mỏ Đá SmartScale. Tôi luôn túc trực hỗ trợ nhân viên và ban điều hành tra cứu nhanh thông tin phiếu cân, kiểm soát lệch bì, đối soát định mức dầu đội xe và phân tích giá thành sản phẩm. Bạn muốn kiểm tra số liệu nào hôm nay?",
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Chạy What-If Simulator để xem dự báo lợi nhuận",
				ActionType:  "navigate",
				TargetURL:   &targetUrl,
				TargetLabel: &label,
			},
		},
		FollowUpQuestions: []string{
			"Hôm nay có bao nhiêu phiếu cân xuất bán?",
			"Có cảnh báo lệch bì xe nào vượt ngưỡng không?",
			"Xe nào đang tiêu hao dầu vượt định mức?",
			"Chi phí 1 tấn đá 1x2 hiện tại là bao nhiêu?",
		},
	}
}

func analyzeTareAlerts() CopilotResponse {
	targetUrlAlert := "/canh-bao-lech-bi/chi-tiet"
	labelAlert := "Xem Chi Tiết Cảnh Báo Lệch Bì"
	targetUrlTruck := "/quan-ly-xe/19H-056.22"
	labelTruck := "Kiểm Tra Xe 19H-056.22"

	return CopilotResponse{
		Answer: "Phát hiện 1 cảnh báo lệch bì nghiêm trọng trong ca trực: Xe 19H-056.22 tại Cổng 1 có khối lượng bì thực tế 26.830 kg, vượt định mức chuẩn 26.450 kg (+380 kg). Nguy cơ tài xế thay đổi kết cấu hoặc chứa nước gian lận.",
		RootCauses: []CopilotCause{
			{
				Factor:    "Lệch bì vượt ngưỡng an toàn cho phép",
				ImpactPct: 1.44,
				Evidence:  "Bì thực tế 26.830 kg vs Chuẩn 26.450 kg (+380 kg / ngưỡng cho phép 200 kg).",
			},
			{
				Factor:    "Thời gian dừng xe tại trạm cân kéo dài",
				ImpactPct: 3.2,
				Evidence:  "Thời gian cân lâu hơn 45 giây so với tiêu chuẩn, nghi vấn thao tác trên bàn cân.",
			},
		},
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Mở biên bản sự cố và camera đối soát lệch bì",
				ActionType:  "navigate",
				TargetURL:   &targetUrlAlert,
				TargetLabel: &labelAlert,
			},
			{
				Action:      "Xem lịch sử đăng kiểm & cân bì xe 19H-056.22",
				ActionType:  "navigate",
				TargetURL:   &targetUrlTruck,
				TargetLabel: &labelTruck,
			},
		},
		FollowUpQuestions: []string{
			"Xe 19H-056.22 thuộc đơn vị vận tải nào?",
			"Lịch sử cân bì 5 chuyến gần nhất của xe này?",
			"Tổng số vụ lệch bì phát hiện trong tuần qua?",
		},
	}
}

func analyzeTicketOperations() CopilotResponse {
	targetUrlTickets := "/phieu-can"
	labelTickets := "Mở Danh Sách Phiếu Cân"
	targetUrlCost := "/chi-huy-dieu-hanh/gia-thanh"
	labelCost := "Xem Giá Thành 1 Tấn Đá"

	return CopilotResponse{
		Answer: "Hôm nay hệ thống trạm cân đã xử lý 142 lượt xe ra vào mỏ. Tổng khối lượng xuất bán thành phẩm đạt 4.860 tấn, doanh thu trong ngày ước tính 685 triệu VNĐ. Mặt hàng xuất nhiều nhất là Đá 1x2 (2.150 tấn) và Đá 4x6 (1.420 tấn).",
		RootCauses: []CopilotCause{
			{
				Factor:    "Nhu cầu công trình hạ tầng tăng đột biến",
				ImpactPct: 24.5,
				Evidence:  "Tăng 18 chuyến xe so với ngày hôm qua từ Tổng Cty Xây Dựng Trường Sơn.",
			},
		},
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Xem danh sách 142 phiếu cân ngày hôm nay",
				ActionType:  "navigate",
				TargetURL:   &targetUrlTickets,
				TargetLabel: &labelTickets,
			},
			{
				Action:      "Xem giá thành sản xuất và biên lợi nhuận đá 1x2",
				ActionType:  "navigate",
				TargetURL:   &targetUrlCost,
				TargetLabel: &labelCost,
			},
		},
		FollowUpQuestions: []string{
			"Có phiếu cân nào bị quá tải trọng không?",
			"Khách hàng nào lấy nhiều đá nhất hôm nay?",
			"Giá bán bình quân đá 1x2 hôm nay là bao nhiêu?",
		},
	}
}

func analyzeInventoryStock() CopilotResponse {
	targetUrlStock := "/kho/tong-quan"
	labelStock := "Xem Tổng Quan Bãi Kho"

	return CopilotResponse{
		Answer: "Trữ lượng tồn kho thành phẩm toàn mỏ hiện còn 42.500 tấn đá các loại. Trong đó: Đá 1x2 còn 14.200 tấn (đủ cấp 7 ngày), Đá 4x6 còn 18.100 tấn, Đá mạt 0x4 còn 7.600 tấn. Riêng đá hộc nguyên khai tại phễu cấp liệu máy nghiền 01 chỉ còn 2.600 tấn — cần đẩy mạnh nổ mìn bốc xúc tầng 3.",
		RootCauses: []CopilotCause{
			{
				Factor:    "Phễu cấp liệu máy nghiền 01 có nguy cơ thiếu nguyên liệu",
				ImpactPct: -15.0,
				Evidence:  "Đá hộc tồn bãi chỉ còn cung cấp được trong 1.5 ca nghiền tiếp theo.",
			},
		},
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Mở quản lý kho bãi và kế hoạch điều phối nguyên liệu",
				ActionType:  "navigate",
				TargetURL:   &targetUrlStock,
				TargetLabel: &labelStock,
			},
		},
		FollowUpQuestions: []string{
			"Bao giờ có đợt nổ mìn tiếp theo tại tầng 3?",
			"Máy nghiền 01 có đang hoạt động hết công suất?",
		},
	}
}

func analyzeSpecificTruck(q string) CopilotResponse {
	plate := "88H-042.27"
	if strings.Contains(q, "19h") {
		plate = "19H-056.22"
	} else if strings.Contains(q, "29c") {
		plate = "29C-789.12"
	}

	targetUrlTruck := "/quan-ly-xe/" + plate
	labelTruck := "Hồ Sơ Xe " + plate

	if plate == "19H-056.22" {
		targetUrlAlert := "/canh-bao-lech-bi/chi-tiet"
		labelAlert := "Xem Biên Bản Lệch Bì"

		return CopilotResponse{
			Answer: fmt.Sprintf("Xe ben %s (Đội xe Tân Uyên - Lái xe Hoàng Đình Trọng). Xe này đang có 1 cảnh báo lệch bì (+380 kg) lúc 13:59 hôm nay tại trạm cân Cổng 1. Mức tiêu hao nhiên liệu trung bình 32.1 L/100km (đạt chuẩn).", plate),
			Recommendations: []CopilotRecommendation{
				{
					Action:      "Xem chi tiết cảnh báo lệch bì vừa phát sinh",
					ActionType:  "navigate",
					TargetURL:   &targetUrlAlert,
					TargetLabel: &labelAlert,
				},
				{
					Action:      "Xem hồ sơ phương tiện & định mức xe " + plate,
					ActionType:  "navigate",
					TargetURL:   &targetUrlTruck,
					TargetLabel: &labelTruck,
				},
			},
			FollowUpQuestions: []string{
				"Lịch sử vi phạm của xe " + plate + " có nhiều không?",
				"Đăng kiểm xe " + plate + " còn bao nhiêu ngày?",
			},
		}
	}

	return CopilotResponse{
		Answer: fmt.Sprintf("Xe ben %s (Đội xe Phú Thọ - Tải trọng 30 tấn). Xe này có biên lợi nhuận thấp (42.1%%) do mức tiêu hao dầu thực tế đạt 38.5 L/100km, cao hơn định mức 31%%. Ngoài ra xe sắp hết hạn đăng kiểm trong 15 ngày tới.", plate),
		RootCauses: []CopilotCause{
			{
				Factor:    "Tiêu hao nhiên liệu bất thường",
				ImpactPct: 31.0,
				Evidence:  "Kim phun và lọc gió có dấu hiệu bẩn, gây hao hụt dầu trên cung đường dốc tầng 2.",
			},
		},
		Recommendations: []CopilotRecommendation{
			{
				Action:      "Mở hồ sơ & lịch sử bảo dưỡng xe " + plate,
				ActionType:  "navigate",
				TargetURL:   &targetUrlTruck,
				TargetLabel: &labelTruck,
			},
		},
		FollowUpQuestions: []string{
			"Chi phí bảo dưỡng định kỳ xe " + plate + " là bao nhiêu?",
			"Lái xe phụ trách xe " + plate + " là ai?",
		},
	}
}

