package handlers

import (
	"net/http"
	"strings"
	"time"
)

type ProductCostPnL struct {
	Product     string  `json:"product"`
	Tons        float64 `json:"tons"`
	Revenue     float64 `json:"revenue"`
	Cost        float64 `json:"cost"`
	Profit      float64 `json:"profit"`
	CostPerTon  float64 `json:"cost_per_ton"`
	PricePerTon float64 `json:"price_per_ton"`
	Margin      float64 `json:"margin"`
	TrendPct    float64 `json:"trend_pct"`
}

type CostCategoryItem struct {
	Category     string  `json:"category"`
	AmountPerTon float64 `json:"amount_per_ton"`
	Pct          float64 `json:"pct"`
	Color        string  `json:"color"`
	TrendPct     float64 `json:"trend_pct"`
}

type CostBreakdownResponse struct {
	Product        string             `json:"product"`
	Period         string             `json:"period"`
	TotalPerTon    float64            `json:"total_per_ton"`
	Breakdown      []CostCategoryItem `json:"breakdown"`
	VsLastMonthPct float64            `json:"vs_last_month_pct"`
	VsStandardPct  float64            `json:"vs_standard_pct"`
}

type VehicleProfitability struct {
	LicensePlate    string  `json:"license_plate"`
	Trips           int     `json:"trips"`
	Tons            float64 `json:"tons"`
	Revenue         float64 `json:"revenue"`
	FuelCost        float64 `json:"fuel_cost"`
	MaintenanceCost float64 `json:"maintenance_cost"`
	OtherCost       float64 `json:"other_cost"`
	Profit          float64 `json:"profit"`
	Margin          float64 `json:"margin"`
	CostPerTon      float64 `json:"cost_per_ton"`
	Flag            string  `json:"flag"` // green | yellow | red
}

type CustomerProfitability struct {
	CustomerName   string  `json:"customer_name"`
	Revenue        float64 `json:"revenue"`
	ProductionCost float64 `json:"production_cost"`
	TransportCost  float64 `json:"transport_cost"`
	Discount       float64 `json:"discount"`
	DebtCost       float64 `json:"debt_cost"`
	Profit         float64 `json:"profit"`
	Margin         float64 `json:"margin"`
	Flag           string  `json:"flag"`
}

func CostPerTon(w http.ResponseWriter, r *http.Request) {
	prodFilter := strings.TrimSpace(r.URL.Query().Get("product"))

	allProducts := []ProductCostPnL{
		{
			Product:     "Đá 1x2 Xây Dựng",
			Tons:        14200,
			Revenue:     3408000000,
			Cost:        2726400000,
			Profit:      681600000,
			CostPerTon:  192000,
			PricePerTon: 240000,
			Margin:      20.0,
			TrendPct:    4.2,
		},
		{
			Product:     "Đá 4x6 Kè Móng",
			Tons:        9800,
			Revenue:     2156000000,
			Cost:        1789000000,
			Profit:      367000000,
			CostPerTon:  182500,
			PricePerTon: 220000,
			Margin:      17.0,
			TrendPct:    -2.5,
		},
		{
			Product:     "Đá Base Cấp Phối Loại 1",
			Tons:        18500,
			Revenue:     3330000000,
			Cost:        2497500000,
			Profit:      832500000,
			CostPerTon:  135000,
			PricePerTon: 180000,
			Margin:      25.0,
			TrendPct:    6.8,
		},
		{
			Product:     "Cát Nghiền Nhân Tạo",
			Tons:        6400,
			Revenue:     1664000000,
			Cost:        1214720000,
			Profit:      449280000,
			CostPerTon:  189800,
			PricePerTon: 260000,
			Margin:      27.0,
			TrendPct:    8.5,
		},
		{
			Product:     "Đá Hộc Xô Bồ Nổ Mìn",
			Tons:        12100,
			Revenue:     1936000000,
			Cost:        1490720000,
			Profit:      445280000,
			CostPerTon:  123200,
			PricePerTon: 160000,
			Margin:      23.0,
			TrendPct:    1.4,
		},
	}

	if prodFilter == "" {
		JSON(w, allProducts)
		return
	}

	var filtered []ProductCostPnL
	for _, p := range allProducts {
		if strings.Contains(strings.ToLower(p.Product), strings.ToLower(prodFilter)) {
			filtered = append(filtered, p)
		}
	}
	JSON(w, filtered)
}

func CostBreakdown(w http.ResponseWriter, r *http.Request) {
	prod := r.URL.Query().Get("product")
	if prod == "" {
		prod = "Đá 1x2 Xây Dựng"
	}
	period := r.URL.Query().Get("period")
	if period == "" {
		period = time.Now().Format("2006-01")
	}

	breakdown := []CostCategoryItem{
		{
			Category:     "Sản xuất & Nổ mìn",
			AmountPerTon: 78000,
			Pct:          40.6,
			Color:        "#1976D2", // modaBlue
			TrendPct:     1.2,
		},
		{
			Category:     "Nhiên liệu Diesel",
			AmountPerTon: 45000,
			Pct:          23.4,
			Color:        "#FF9800", // warningOrange
			TrendPct:     12.5,
		},
		{
			Category:     "Tiền lương nhân công",
			AmountPerTon: 32000,
			Pct:          16.7,
			Color:        "#43A047", // successGreen
			TrendPct:     0.0,
		},
		{
			Category:     "Khấu hao máy móc",
			AmountPerTon: 21000,
			Pct:          10.9,
			Color:        "#8E24AA", // purpleAccent
			TrendPct:     0.0,
		},
		{
			Category:     "Vận chuyển nội bộ",
			AmountPerTon: 16000,
			Pct:          8.4,
			Color:        "#00ACC1", // teal
			TrendPct:     4.8,
		},
	}

	totalPerTon := 192000.0

	// Adjust for 4x6 if selected
	if strings.Contains(prod, "4x6") {
		totalPerTon = 182500.0
		breakdown = []CostCategoryItem{
			{Category: "Sản xuất & Nổ mìn", AmountPerTon: 71000, Pct: 38.9, Color: "#1976D2", TrendPct: -1.0},
			{Category: "Nhiên liệu Diesel", AmountPerTon: 44500, Pct: 24.4, Color: "#FF9800", TrendPct: 8.2},
			{Category: "Tiền lương nhân công", AmountPerTon: 31000, Pct: 17.0, Color: "#43A047", TrendPct: 0.0},
			{Category: "Khấu hao máy móc", AmountPerTon: 21000, Pct: 11.5, Color: "#8E24AA", TrendPct: 0.0},
			{Category: "Vận chuyển nội bộ", AmountPerTon: 15000, Pct: 8.2, Color: "#00ACC1", TrendPct: 2.1},
		}
	}

	resp := CostBreakdownResponse{
		Product:        prod,
		Period:         period,
		TotalPerTon:    totalPerTon,
		Breakdown:      breakdown,
		VsLastMonthPct: 4.8,
		VsStandardPct:  6.2,
	}

	JSON(w, resp)
}

func VehicleProfitabilityHandler(w http.ResponseWriter, r *http.Request) {
	vehicles := []VehicleProfitability{
		{
			LicensePlate:    "19H-056.22",
			Trips:           98,
			Tons:            3077.0,
			Revenue:         738480000,
			FuelCost:        162000000,
			MaintenanceCost: 45000000,
			OtherCost:       22000000,
			Profit:          509480000,
			Margin:          69.0,
			CostPerTon:      74423,
			Flag:            "green",
		},
		{
			LicensePlate:    "29H-882.19",
			Trips:           76,
			Tons:            2006.0,
			Revenue:         521560000,
			FuelCost:        125000000,
			MaintenanceCost: 38000000,
			OtherCost:       16000000,
			Profit:          342560000,
			Margin:          65.7,
			CostPerTon:      89232,
			Flag:            "green",
		},
		{
			LicensePlate:    "90C-123.45",
			Trips:           64,
			Tons:            2880.0,
			Revenue:         633600000,
			FuelCost:        195000000,
			MaintenanceCost: 52000000,
			OtherCost:       25000000,
			Profit:          361600000,
			Margin:          57.1,
			CostPerTon:      94444,
			Flag:            "green",
		},
		{
			LicensePlate:    "88H-042.27",
			Trips:           142,
			Tons:            4711.0,
			Revenue:         1130640000,
			FuelCost:        435000000,
			MaintenanceCost: 145000000,
			OtherCost:       75000000,
			Profit:          475640000,
			Margin:          42.1,
			CostPerTon:      139036,
			Flag:            "yellow",
		},
		{
			LicensePlate:    "19C-098.76",
			Trips:           35,
			Tons:            700.0,
			Revenue:         154000000,
			FuelCost:        65000000,
			MaintenanceCost: 38000000,
			OtherCost:       15000000,
			Profit:          36000000,
			Margin:          23.4,
			CostPerTon:      168571,
			Flag:            "red",
		},
	}

	JSON(w, vehicles)
}

func CustomerProfitabilityHandler(w http.ResponseWriter, r *http.Request) {
	customers := []CustomerProfitability{
		{
			CustomerName:   "Công ty CP Đầu Tư Xây Dựng 319",
			Revenue:        4850000000,
			ProductionCost: 3250000000,
			TransportCost:  450000000,
			Discount:       50000000,
			DebtCost:       12000000,
			Profit:         1088000000,
			Margin:         22.4,
			Flag:           "green",
		},
		{
			CustomerName:   "Công ty Bê Tông Việt Trì",
			Revenue:        3100000000,
			ProductionCost: 2100000000,
			TransportCost:  280000000,
			Discount:       20000000,
			DebtCost:       8000000,
			Profit:         692000000,
			Margin:         22.3,
			Flag:           "green",
		},
		{
			CustomerName:   "Tổng Công Ty XD Trường Sơn",
			Revenue:        8200000000,
			ProductionCost: 5900000000,
			TransportCost:  950000000,
			Discount:       120000000,
			DebtCost:       38000000,
			Profit:         1192000000,
			Margin:         14.5,
			Flag:           "yellow",
		},
	}

	JSON(w, customers)
}
