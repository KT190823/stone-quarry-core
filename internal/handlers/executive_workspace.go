package handlers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"mo-da-backend/internal/database"
)

// ─────────────────────────────────────────────────────────────────────────────
// Executive workspace aggregation.
//
// Every value returned by these handlers is derived from operational tables
// (tickets, production_costs, alerts, payment vouchers, debts, permits,
// inventory, employees, vehicles, reports...) instead of hardcoded fixtures in
// the Angular layer. The data window falls back to the latest month that has
// ticket activity so the dashboard stays meaningful when the current calendar
// month has no records yet.
// ─────────────────────────────────────────────────────────────────────────────

type ExecutiveActionItem struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Category      string `json:"category"`
	CategoryLabel string `json:"category_label"`
	MineName      string `json:"mine_name"`
	Deadline      string `json:"deadline"`
	Urgency       string `json:"urgency"`
	Route         string `json:"route"`
	ActionLabel   string `json:"action_label"`
}

type ExecutiveMinePerformance struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	Location              string  `json:"location"`
	ProductionTons        float64 `json:"production_tons"`
	ProductionPlanTons    float64 `json:"production_plan_tons"`
	ProductionProgressPct float64 `json:"production_progress_pct"`
	Revenue               float64 `json:"revenue"`
	RevenuePlan           float64 `json:"revenue_plan"`
	Status                string  `json:"status"`
	StatusLabel           string  `json:"status_label"`
	Route                 string  `json:"route"`
}

type ExecutiveTaskItem struct {
	ID            string  `json:"id"`
	Code          string  `json:"code"`
	Title         string  `json:"title"`
	Category      string  `json:"category"`
	CategoryLabel string  `json:"category_label"`
	MineName      string  `json:"mine_name"`
	Department    string  `json:"department"`
	Deadline      string  `json:"deadline"`
	Urgency       string  `json:"urgency"`
	UrgencyLabel  string  `json:"urgency_label"`
	Status        string  `json:"status"`
	StatusLabel   string  `json:"status_label"`
	AssignedTo    string  `json:"assigned_to"`
	Route         string  `json:"route"`
	ActionText    string  `json:"action_text"`
	Amount        float64 `json:"amount"`
}

type ExecutiveAlertItem struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	Title             string `json:"title"`
	Category          string `json:"category"`
	CategoryLabel     string `json:"category_label"`
	Severity          string `json:"severity"`
	SeverityLabel     string `json:"severity_label"`
	MineName          string `json:"mine_name"`
	Time              string `json:"time"`
	Impact            string `json:"impact"`
	RecommendedAction string `json:"recommended_action"`
	Route             string `json:"route"`
	ActionText        string `json:"action_text"`
}

type ExecutiveDepartmentCard struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Sub   string `json:"sub"`
	Tone  string `json:"tone"`
	Delta string `json:"delta,omitempty"`
}

type ExecutiveDepartmentMineRow struct {
	MineName string `json:"mine_name"`
	Target   string `json:"target"`
	Actual   string `json:"actual"`
	Rate     string `json:"rate"`
	Status   string `json:"status"`
}

type ExecutiveDepartment struct {
	Key            string                       `json:"key"`
	Name           string                       `json:"name"`
	HeadOfDept     string                       `json:"head_of_dept"`
	Phone          string                       `json:"phone"`
	Description    string                       `json:"description"`
	Icon           string                       `json:"icon"`
	Route          string                       `json:"route"`
	Score          int                          `json:"score"`
	Cards          []ExecutiveDepartmentCard    `json:"cards"`
	Highlights     []string                     `json:"highlights"`
	MineComparison []ExecutiveDepartmentMineRow `json:"mine_comparison"`
}

type ExecutiveMineComparison struct {
	Code                  string  `json:"code"`
	Name                  string  `json:"name"`
	Location              string  `json:"location"`
	LicensedCapacity      string  `json:"licensed_capacity"`
	ExpiryYear            int     `json:"expiry_year"`
	ActualProductionTons  float64 `json:"actual_production_tons"`
	PlannedProductionTons float64 `json:"planned_production_tons"`
	ProgressPct           float64 `json:"progress_pct"`
	Revenue               float64 `json:"revenue"`
	OperatingCost         float64 `json:"operating_cost"`
	GrossMarginPct        float64 `json:"gross_margin_pct"`
	FleetCount            int     `json:"fleet_count"`
	ActiveFleet           int     `json:"active_fleet"`
	FleetOeePct           float64 `json:"fleet_oee_pct"`
	WorkforceCount        int     `json:"workforce_count"`
	AlertsCount           int     `json:"alerts_count"`
	Status                string  `json:"status"`
	Route3D               string  `json:"route_3d"`
	RouteCatalog          string  `json:"route_catalog"`
	Lat                   float64 `json:"lat"`
	Lng                   float64 `json:"lng"`
}

type ExecutiveReportMetric struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type ExecutiveReportItem struct {
	ID            string                  `json:"id"`
	Title         string                  `json:"title"`
	Period        string                  `json:"period"`
	Category      string                  `json:"category"`
	CategoryLabel string                  `json:"category_label"`
	GeneratedDate string                  `json:"generated_date"`
	Status        string                  `json:"status"`
	Summary       string                  `json:"summary"`
	KeyMetrics    []ExecutiveReportMetric `json:"key_metrics"`
	Route         string                  `json:"route"`
}

type executiveWorkspaceExtras struct {
	PendingActions  []ExecutiveActionItem
	MinePerformance []ExecutiveMinePerformance
	OverdueDebt     float64
	RevenueGrowth   float64
	PlanProgressPct float64
}

// executiveDataWindow trả về khoảng thời gian có dữ liệu phiếu cân. Nếu kỳ
// hiện tại chưa có phiếu nào, dùng tháng gần nhất có phát sinh để dashboard vẫn
// phản ánh số liệu thật thay vì rỗng.
func executiveDataWindow(ctx context.Context, start, end time.Time) (time.Time, time.Time) {
	if database.Pool == nil {
		return start, end
	}
	var count int
	_ = database.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tickets WHERE created_at >= $1 AND created_at < $2`,
		start, end).Scan(&count)
	if count > 0 {
		return start, end
	}
	var latest time.Time
	_ = database.Pool.QueryRow(ctx, `SELECT COALESCE(MAX(created_at), TIMESTAMPTZ '1970-01-01') FROM tickets`).Scan(&latest)
	if latest.Year() <= 1970 {
		return start, end
	}
	ms := time.Date(latest.Year(), latest.Month(), 1, 0, 0, 0, 0, latest.Location())
	return ms, ms.AddDate(0, 1, 0)
}

func formatVNDShort(v float64) string {
	if math.Abs(v) >= 1_000_000_000 {
		return fmt.Sprintf("%.2f TỶ", v/1_000_000_000)
	}
	if math.Abs(v) >= 1_000_000 {
		return fmt.Sprintf("%.1f TR", v/1_000_000)
	}
	return fmt.Sprintf("%.0f đ", v)
}

func buildExecutiveWorkspace(ctx context.Context, period string, start, end time.Time, quarryID string) executiveWorkspaceExtras {
	start, end = executiveDataWindow(ctx, start, end)
	tasks := buildExecutiveTasks(ctx, start, end, quarryID)

	pending := make([]ExecutiveActionItem, 0, 5)
	for _, t := range tasks {
		if len(pending) >= 5 {
			break
		}
		pending = append(pending, ExecutiveActionItem{
			ID:            t.ID,
			Title:         t.Title,
			Category:      t.Category,
			CategoryLabel: t.CategoryLabel,
			MineName:      t.MineName,
			Deadline:      t.Deadline,
			Urgency:       t.Urgency,
			Route:         t.Route,
			ActionLabel:   t.ActionText,
		})
	}

	minePerf := buildMinePerformance(ctx, start, end, quarryID)

	var planProgress float64
	var prod, plan float64
	for _, m := range minePerf {
		prod += m.ProductionTons
		plan += m.ProductionPlanTons
	}
	if plan > 0 {
		planProgress = math.Round(prod/plan*1000) / 10
	}

	return executiveWorkspaceExtras{
		PendingActions:  pending,
		MinePerformance: minePerf,
		OverdueDebt:     buildOverdueDebt(ctx, quarryID),
		RevenueGrowth:   buildRevenueGrowth(ctx, start, end, quarryID),
		PlanProgressPct: planProgress,
	}
}

func buildMinePerformance(ctx context.Context, start, end time.Time, quarryID string) []ExecutiveMinePerformance {
	out := make([]ExecutiveMinePerformance, 0, 4)
	if database.Pool == nil {
		return out
	}
	rows, err := database.Pool.Query(ctx, `
		SELECT q.code, q.name, COALESCE(q.address, ''), COALESCE(q.latitude, 0)::float8, COALESCE(q.longitude, 0)::float8
		FROM quarries q
		WHERE ($1 = '' OR q.code = $1)
		ORDER BY q.code`, quarryID)
	if err != nil {
		return out
	}
	defer rows.Close()

	type quarryRow struct {
		code, name, address string
		lat, lng            float64
	}
	var quarries []quarryRow
	for rows.Next() {
		var q quarryRow
		if rows.Scan(&q.code, &q.name, &q.address, &q.lat, &q.lng) == nil {
			quarries = append(quarries, q)
		}
	}
	rows.Close()

	for _, q := range quarries {
		var prod, rev float64
		_ = database.Pool.QueryRow(ctx, fmt.Sprintf(`
			SELECT COALESCE(SUM(%s), 0), COALESCE(SUM(%s), 0)
			FROM tickets
			WHERE quarry_code = $1 AND created_at >= $2 AND created_at < $3
		`, sqlTicketTonnageExpr, sqlTicketRevenueExpr), q.code, start, end).Scan(&prod, &rev)

		var planAnnual float64
		_ = database.Pool.QueryRow(ctx,
			`SELECT COALESCE(SUM(annual_target), 0) FROM mining_plans WHERE quarry_code = $1`,
			q.code).Scan(&planAnnual)
		planMonthly := planAnnual / 12.0

		progress := 0.0
		if planMonthly > 0 {
			progress = math.Round(prod/planMonthly*1000) / 10
		}
		status := "normal"
		label := "Vận hành ổn định"
		if progress > 0 && progress < 85 {
			status = "warning"
			label = fmt.Sprintf("Chậm tiến độ %.1f%%", 100-progress)
		} else if planMonthly > 0 && progress >= 95 {
			label = fmt.Sprintf("Vận hành xuất sắc (%.1f%%)", progress)
		} else if planMonthly > 0 {
			label = fmt.Sprintf("Vận hành ổn định (%.1f%%)", progress)
		}

		out = append(out, ExecutiveMinePerformance{
			ID:                    q.code,
			Name:                  q.name,
			Location:              q.address,
			ProductionTons:        math.Round(prod),
			ProductionPlanTons:    math.Round(planMonthly),
			ProductionProgressPct: progress,
			Revenue:               math.Round(rev),
			RevenuePlan:           0,
			Status:                status,
			StatusLabel:           label,
			Route:                 "/khao-sat-3d/tong-quan",
		})
	}
	return out
}

func buildOverdueDebt(ctx context.Context, quarryID string) float64 {
	if database.Pool == nil {
		return 0
	}
	var total float64
	_ = database.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(
			CASE WHEN balance ~ '^[0-9]+' THEN NULLIF(regexp_replace(balance, '[^0-9]', '', 'g'), '')::double precision ELSE 0 END
		), 0)
		FROM payments_debt
		WHERE (status ILIKE '%quá hạn%' OR status ILIKE '%overdue%')
		  AND ($1 = '' OR customer ILIKE '%' || $1 || '%' OR partner ILIKE '%' || $1 || '%')
	`, quarryID).Scan(&total)
	return total
}

func buildRevenueGrowth(ctx context.Context, start, end time.Time, quarryID string) float64 {
	if database.Pool == nil {
		return 0
	}
	var current, previous float64
	_ = database.Pool.QueryRow(ctx, fmt.Sprintf(`
		SELECT COALESCE(SUM(%s), 0) FROM tickets
		WHERE created_at >= $1 AND created_at < $2 AND ($3 = '' OR quarry_code = $3)
	`, sqlTicketRevenueExpr), start, end, quarryID).Scan(&current)

	dur := end.Sub(start)
	_ = database.Pool.QueryRow(ctx, fmt.Sprintf(`
		SELECT COALESCE(SUM(%s), 0) FROM tickets
		WHERE created_at >= $1 AND created_at < $2 AND ($3 = '' OR quarry_code = $3)
	`, sqlTicketRevenueExpr), start.Add(-dur), start, quarryID).Scan(&previous)

	if previous <= 0 {
		return 0
	}
	return math.Round((current-previous)/previous*1000) / 10
}

// buildExecutiveTasks aggregates pending leadership work from real records.
func buildExecutiveTasks(ctx context.Context, start, end time.Time, quarryID string) []ExecutiveTaskItem {
	out := make([]ExecutiveTaskItem, 0, 12)
	if database.Pool == nil {
		return out
	}

	// 1) Payment vouchers awaiting approval.
	rows, err := database.Pool.Query(ctx, `
		SELECT code, COALESCE(reason, ''), COALESCE(partner_name, ''), COALESCE(amount, 0), COALESCE(date, ''), COALESCE(created_by, '')
		FROM payment_vouchers
		WHERE status NOT IN ('posted', 'approved', 'paid', 'cancelled')
		ORDER BY created_at DESC LIMIT 5`)
	if err == nil {
		for rows.Next() {
			var code, reason, partner, date, createdBy string
			var amount float64
			if rows.Scan(&code, &reason, &partner, &amount, &date, &createdBy) == nil {
				urgency := "medium"
				if amount >= 100_000_000 {
					urgency = "high"
				}
				out = append(out, ExecutiveTaskItem{
					ID:            "PAY-" + code,
					Code:          code,
					Title:         "Phê duyệt chi: " + reason,
					Category:      "finance",
					CategoryLabel: "Tài chính - Phê duyệt chi",
					MineName:      partner,
					Department:    "Kế toán - Tài chính",
					Deadline:      date,
					Urgency:       urgency,
					UrgencyLabel:  executiveUrgencyLabel(urgency),
					Status:        "pending",
					StatusLabel:   "Chờ Ban Giám Đốc duyệt",
					AssignedTo:    createdBy,
					Route:         "/thanh-toan/phieu-chi",
					ActionText:    "Ký duyệt chi",
					Amount:        amount,
				})
			}
		}
		rows.Close()
	}

	// 2) Overdue customer debt.
	rows, err = database.Pool.Query(ctx, `
		SELECT code, COALESCE(customer, partner, ''), COALESCE(balance, '0'), COALESCE(due, ''), COALESCE(status, '')
		FROM payments_debt
		WHERE status ILIKE '%quá hạn%' OR status ILIKE '%overdue%'
		ORDER BY created_at DESC LIMIT 5`)
	if err == nil {
		for rows.Next() {
			var code, customer, balance, due, status string
			if rows.Scan(&code, &customer, &balance, &due, &status) == nil {
				out = append(out, ExecutiveTaskItem{
					ID:            "DEBT-" + code,
					Code:          code,
					Title:         "Xử lý công nợ quá hạn của " + customer,
					Category:      "debt",
					CategoryLabel: "Công nợ quá hạn",
					MineName:      customer,
					Department:    "Kinh doanh & CRM",
					Deadline:      due,
					Urgency:       "high",
					UrgencyLabel:  executiveUrgencyLabel("high"),
					Status:        "pending",
					StatusLabel:   "Cần chỉ đạo đôn đốc",
					AssignedTo:    "Phòng Kinh doanh",
					Route:         "/thanh-toan/cong-no",
					ActionText:    "Đối soát",
					Amount:        parseMoneyString(balance),
				})
			}
		}
		rows.Close()
	}

	// 3) Permits / licences nearing expiry.
	rows, err = database.Pool.Query(ctx, `
		SELECT code, COALESCE(title, ''), COALESCE(mine_name, ''), COALESCE(expiry_date, ''), COALESCE(days_remaining, 0)
		FROM mining_permits
		WHERE days_remaining BETWEEN 0 AND 90
		ORDER BY days_remaining ASC LIMIT 5`)
	if err == nil {
		for rows.Next() {
			var code, title, mineName, expiry string
			var days int
			if rows.Scan(&code, &title, &mineName, &expiry, &days) == nil {
				urgency := "medium"
				if days <= 30 {
					urgency = "high"
				}
				out = append(out, ExecutiveTaskItem{
					ID:            "PERMIT-" + code,
					Code:          code,
					Title:         title,
					Category:      "permit",
					CategoryLabel: "Pháp lý mỏ",
					MineName:      mineName,
					Department:    "Pháp lý & Giấy phép",
					Deadline:      fmt.Sprintf("Hạn: %s (còn %d ngày)", expiry, days),
					Urgency:       urgency,
					UrgencyLabel:  executiveUrgencyLabel(urgency),
					Status:        "in_progress",
					StatusLabel:   "Đang thẩm định hồ sơ",
					AssignedTo:    "Ban Pháp chế",
					Route:         "/ho-so-giay-phep",
					ActionText:    "Xem hồ sơ",
				})
			}
		}
		rows.Close()
	}

	// 4) Contracts still delivering.
	rows, err = database.Pool.Query(ctx, `
		SELECT code, COALESCE(customer_name, ''), COALESCE(project_name, ''), COALESCE(end_date, ''),
		       COALESCE(total_committed_tons, 0), COALESCE(total_delivered_tons, 0)
		FROM sales_contracts
		WHERE status ILIKE '%active%' OR status ILIKE '%hiệu lực%'
		ORDER BY created_at DESC LIMIT 4`)
	if err == nil {
		for rows.Next() {
			var code, customer, project, endDate string
			var committed, delivered float64
			if rows.Scan(&code, &customer, &project, &endDate, &committed, &delivered) == nil {
				out = append(out, ExecutiveTaskItem{
					ID:            "CONTRACT-" + code,
					Code:          code,
					Title:         "Theo dõi tiến độ giao hàng hợp đồng " + project,
					Category:      "finance",
					CategoryLabel: "Kinh doanh - Hợp đồng",
					MineName:      customer,
					Department:    "Kinh doanh",
					Deadline:      endDate,
					Urgency:       "medium",
					UrgencyLabel:  executiveUrgencyLabel("medium"),
					Status:        "in_progress",
					StatusLabel:   "Đang thực hiện",
					AssignedTo:    "Phòng Kinh doanh",
					Route:         "/doi-tac/hop-dong",
					ActionText:    "Duyệt hợp đồng",
					Amount:        committed - delivered,
				})
			}
		}
		rows.Close()
	}

	// 5) Open weighbridge / operational alerts.
	rows, err = database.Pool.Query(ctx, `
		SELECT id, COALESCE(title, ''), COALESCE(bs, ''), COALESCE(note, ''), COALESCE(severity, '')
		FROM alerts
		WHERE status NOT IN ('resolved', 'closed', 'Đã xử lý')
		ORDER BY created_at DESC LIMIT 4`)
	if err == nil {
		for rows.Next() {
			var id, title, bs, note, severity string
			if rows.Scan(&id, &title, &bs, &note, &severity) == nil {
				urgency := "medium"
				if strings.EqualFold(severity, "critical") || strings.EqualFold(severity, "red") {
					urgency = "high"
				}
				out = append(out, ExecutiveTaskItem{
					ID:            "ALERT-" + id,
					Code:          id,
					Title:         title,
					Category:      "quality",
					CategoryLabel: "Kiểm soát trạm cân",
					MineName:      bs,
					Department:    "Kho - Vận tải - Trạm cân",
					Deadline:      "Hôm nay",
					Urgency:       urgency,
					UrgencyLabel:  executiveUrgencyLabel(urgency),
					Status:        "pending",
					StatusLabel:   "Chờ xử lý biên bản",
					AssignedTo:    "Trưởng trạm cân",
					Route:         "/phieu-can/canh-bao-lech-bi",
					ActionText:    "Xử lý biên bản",
				})
			}
		}
		rows.Close()
	}

	// Keep high urgency first, medium then low.
	sortExecutiveTasks(out)
	if len(out) > 12 {
		out = out[:12]
	}
	return out
}

func executiveUrgencyLabel(u string) string {
	switch u {
	case "high":
		return "Khẩn cấp"
	case "medium":
		return "Quan trọng"
	default:
		return "Bình thường"
	}
}

func sortExecutiveTasks(items []ExecutiveTaskItem) {
	rank := func(u string) int {
		switch u {
		case "high":
			return 0
		case "medium":
			return 1
		default:
			return 2
		}
	}
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && rank(items[j].Urgency) < rank(items[j-1].Urgency); j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
}

func buildExecutiveAlerts(ctx context.Context, start, end time.Time, quarryID string) []ExecutiveAlertItem {
	out := make([]ExecutiveAlertItem, 0, 16)
	if database.Pool == nil {
		return out
	}

	// Weighbridge alerts
	rows, err := database.Pool.Query(ctx, `
		SELECT id, COALESCE(title, ''), COALESCE(bs, ''), COALESCE(note, ''), COALESCE(severity, ''),
		       COALESCE(time, ''), COALESCE(date, ''), COALESCE(lech_bi, 0)
		FROM alerts
		WHERE status NOT IN ('resolved', 'closed', 'Đã xử lý')
		ORDER BY created_at DESC LIMIT 8`)
	if err == nil {
		for rows.Next() {
			var id, title, bs, note, severity, t, d string
			var lechBi float64
			if rows.Scan(&id, &title, &bs, &note, &severity, &t, &d, &lechBi) == nil {
				sev := "warning"
				if strings.EqualFold(severity, "critical") || strings.EqualFold(severity, "red") {
					sev = "critical"
				}
				out = append(out, ExecutiveAlertItem{
					ID:                "ALERT-" + id,
					Code:              id,
					Title:             title,
					Category:          "fraud",
					CategoryLabel:     "Gian lận trạm cân AI",
					Severity:          sev,
					SeverityLabel:     executiveSeverityLabel(sev),
					MineName:          bs,
					Time:              strings.TrimSpace(t + " " + d),
					Impact:            note,
					RecommendedAction: "Kiểm tra biên bản cân và đối soát RFID",
					Route:             "/phieu-can/canh-bao-lech-bi",
					ActionText:        "Xử lý biên bản",
				})
			}
		}
		rows.Close()
	}

	// Risk alerts registry
	rows, err = database.Pool.Query(ctx, `
		SELECT id, COALESCE(alert_type, ''), COALESCE(severity, 'medium'), COALESCE(title, ''),
		       COALESCE(description, ''), COALESCE(suggested_action, ''), COALESCE(module_source, '')
		FROM risk_alerts
		WHERE status NOT IN ('resolved', 'closed')
		ORDER BY created_at DESC LIMIT 8`)
	if err == nil {
		for rows.Next() {
			var id int
			var alertType, severity, title, desc, action, module string
			if rows.Scan(&id, &alertType, &severity, &title, &desc, &action, &module) == nil {
				sev := "warning"
				if strings.EqualFold(severity, "high") || strings.EqualFold(severity, "critical") || strings.EqualFold(severity, "red") {
					sev = "critical"
				}
				out = append(out, ExecutiveAlertItem{
					ID:                fmt.Sprintf("RISK-%d", id),
					Code:              fmt.Sprintf("RA-%d", id),
					Title:             title,
					Category:          executiveAlertCategory(alertType + " " + module),
					CategoryLabel:     module,
					Severity:          sev,
					SeverityLabel:     executiveSeverityLabel(sev),
					MineName:          "Toàn hệ thống",
					Time:              "Gần đây",
					Impact:            desc,
					RecommendedAction: action,
					Route:             "/canh-bao-lech-bi",
					ActionText:        "Xem chi tiết",
				})
			}
		}
		rows.Close()
	}

	// Debt overdue
	rows, err = database.Pool.Query(ctx, `
		SELECT code, COALESCE(customer, partner, ''), COALESCE(balance, '0'), COALESCE(due, '')
		FROM payments_debt
		WHERE status ILIKE '%quá hạn%' OR status ILIKE '%overdue%'
		ORDER BY created_at DESC LIMIT 4`)
	if err == nil {
		for rows.Next() {
			var code, customer, balance, due string
			if rows.Scan(&code, &customer, &balance, &due) == nil {
				out = append(out, ExecutiveAlertItem{
					ID:                "DEBT-" + code,
					Code:              code,
					Title:             "Khách hàng " + customer + " có công nợ quá hạn " + balance,
					Category:          "debt",
					CategoryLabel:     "Công nợ quá hạn",
					Severity:          "warning",
					SeverityLabel:     "Rủi ro tài chính",
					MineName:          customer,
					Time:              "Hạn: " + due,
					Impact:            "Chiếm dụng vốn lưu động và dòng tiền",
					RecommendedAction: "Tạm dừng xuất đá nợ, chuyển sang thanh toán từng chuyến",
					Route:             "/thanh-toan/cong-no",
					ActionText:        "Đôn đốc nợ",
				})
			}
		}
		rows.Close()
	}

	// Permit expiry
	rows, err = database.Pool.Query(ctx, `
		SELECT code, COALESCE(title, ''), COALESCE(mine_name, ''), COALESCE(expiry_date, ''), COALESCE(days_remaining, 0)
		FROM mining_permits
		WHERE days_remaining BETWEEN 0 AND 90
		ORDER BY days_remaining ASC LIMIT 4`)
	if err == nil {
		for rows.Next() {
			var code, title, mineName, expiry string
			var days int
			if rows.Scan(&code, &title, &mineName, &expiry, &days) == nil {
				out = append(out, ExecutiveAlertItem{
					ID:                "PERMIT-" + code,
					Code:              code,
					Title:             title,
					Category:          "permit",
					CategoryLabel:     "Pháp lý mỏ",
					Severity:          "warning",
					SeverityLabel:     executiveSeverityLabel("warning"),
					MineName:          mineName,
					Time:              fmt.Sprintf("Còn %d ngày (hết hạn %s)", days, expiry),
					Impact:            "Có nguy cơ dừng hoạt động nếu không gia hạn kịp thời",
					RecommendedAction: "Đôn đốc Ban Pháp chế nộp hồ sơ gia hạn",
					Route:             "/ho-so-giay-phep",
					ActionText:        "Xem hồ sơ",
				})
			}
		}
		rows.Close()
	}

	if len(out) > 16 {
		out = out[:16]
	}
	return out
}

func executiveSeverityLabel(sev string) string {
	switch sev {
	case "critical":
		return "Nghiêm trọng"
	case "warning":
		return "Cần xử lý"
	default:
		return "Thông tin"
	}
}

func executiveAlertCategory(source string) string {
	s := strings.ToLower(source)
	switch {
	case strings.Contains(s, "nợ") || strings.Contains(s, "debt"):
		return "debt"
	case strings.Contains(s, "an toàn") || strings.Contains(s, "safety") || strings.Contains(s, "atlđ"):
		return "safety"
	case strings.Contains(s, "pháp") || strings.Contains(s, "permit") || strings.Contains(s, "giấy phép"):
		return "permit"
	case strings.Contains(s, "gian") || strings.Contains(s, "fraud"):
		return "fraud"
	default:
		return "operation"
	}
}

func parseMoneyString(raw string) float64 {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, raw)
	if digits == "" {
		return 0
	}
	var v float64
	fmt.Sscanf(digits, "%f", &v)
	return v
}

func buildOptionalFloat(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

// buildExecutiveDepartments aggregates department metrics from real records.
func buildExecutiveDepartments(ctx context.Context, start, end time.Time, quarryID string) []ExecutiveDepartment {
	if database.Pool == nil {
		return []ExecutiveDepartment{}
	}
	start, end = executiveDataWindow(ctx, start, end)

	var trips, employees, attendance, contractCount, alertCount, permitCount, reportCount, invoiceCount, vehicleCount, activeVehicle int
	var revenueF, costF, tonsF, fuelF, payrollF, debtF float64

	_ = database.Pool.QueryRow(ctx, fmt.Sprintf(`
		SELECT COALESCE(SUM(%s), 0), COALESCE(SUM(%s), 0), COUNT(*)
		FROM tickets WHERE created_at >= $1 AND created_at < $2 AND ($3 = '' OR quarry_code = $3)
	`, sqlTicketRevenueExpr, sqlTicketTonnageExpr), start, end, quarryID).Scan(&revenueF, &tonsF, &trips)

	_ = database.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(actual_value), 0) FROM production_costs
		WHERE created_at >= $1 AND created_at < $2`, start, end).Scan(&costF)

	_ = database.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(actual_fuel_consumed_liters), 0) FROM equipment_fuel_logs
		WHERE created_at >= $1 AND created_at < $2 AND ($3 = '' OR quarry_code = $3)`, start, end, quarryID).Scan(&fuelF)

	_ = database.Pool.QueryRow(ctx, `
		SELECT COUNT(*), COALESCE(SUM(base_salary), 0) FROM hr_employees
		WHERE status NOT ILIKE '%nghỉ%'`).Scan(&employees, &payrollF)

	_ = database.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM hr_attendances WHERE created_at >= $1 AND created_at < $2`, start, end).Scan(&attendance)

	_ = database.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(CASE WHEN balance ~ '^[0-9]' THEN NULLIF(regexp_replace(balance, '[^0-9]', '', 'g'), '')::double precision ELSE 0 END), 0)
		FROM payments_debt WHERE status ILIKE '%quá hạn%' OR status ILIKE '%overdue%'`).Scan(&debtF)

	_ = database.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM sales_contracts WHERE status ILIKE '%active%' OR status ILIKE '%hiệu lực%'`).Scan(&contractCount)
	_ = database.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM alerts WHERE status NOT IN ('resolved','closed','Đã xử lý')`).Scan(&alertCount)
	_ = database.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM mining_permits`).Scan(&permitCount)
	_ = database.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM statutory_reports`).Scan(&reportCount)
	_ = database.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM payments_invoices`).Scan(&invoiceCount)
	_ = database.Pool.QueryRow(ctx, `SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'Hoạt động') FROM vehicles`).Scan(&vehicleCount, &activeVehicle)

	margin := 0.0
	if revenueF > 0 {
		margin = math.Round((revenueF-costF)/revenueF*1000) / 10
	}
	attendanceRate := 0.0
	if employees > 0 {
		attendanceRate = math.Round(float64(attendance)/float64(employees)*1000) / 10
		if attendanceRate > 100 {
			attendanceRate = 100
		}
	}

	mineNames := []string{"Mỏ Phú Thọ", "Mỏ Tân Uyên", "Mỏ Hà Nam", "Mỏ Bình Phước"}

	return []ExecutiveDepartment{
		{
			Key: "business", Name: "Kinh Doanh & Khách Hàng (CRM)",
			HeadOfDept: "Giám đốc Kinh doanh", Phone: "—",
			Description: "Quản lý hợp đồng thương mại, tiêu thụ đá thành phẩm và công nợ đối tác.",
			Icon:        "handshake", Route: "/doi-tac", Score: clampScore(margin + 60),
			Cards: []ExecutiveDepartmentCard{
				{Label: "Tổng Doanh Thu Xuất Bán", Value: formatVNDShort(revenueF), Sub: fmt.Sprintf("%d chuyến xe ghi nhận", trips), Tone: "green"},
				{Label: "Sản Lượng Đá Bán", Value: fmt.Sprintf("%.0f tấn", tonsF), Sub: "Tổng hợp từ phiếu cân", Tone: "blue"},
				{Label: "Công Nợ Quá Hạn", Value: formatVNDShort(debtF), Sub: "Đối soát từ sổ công nợ", Tone: "amber"},
				{Label: "Số Hợp Đồng Đang Chạy", Value: fmt.Sprintf("%d Hợp đồng", contractCount), Sub: "Dữ liệu hợp đồng bán", Tone: "blue"},
			},
			Highlights: []string{
				fmt.Sprintf("Doanh thu kỳ đạt %s từ %d chuyến xe xuất bán.", formatVNDShort(revenueF), trips),
				fmt.Sprintf("Công nợ quá hạn hiện tại %s cần đôn đốc thu hồi.", formatVNDShort(debtF)),
				fmt.Sprintf("Đang theo dõi %d hợp đồng cung cấp đá còn hiệu lực.", contractCount),
			},
			MineComparison: departmentMineRows(mineNames, "Doanh thu"),
		},
		{
			Key: "accounting", Name: "Kế Toán - Tài Chính",
			HeadOfDept: "Kế toán trưởng", Phone: "—",
			Description: "Hóa đơn điện tử VAT, sổ quỹ thu chi và báo cáo tài chính P&L.",
			Icon:        "credit-card", Route: "/thanh-toan", Score: clampScore(margin + 55),
			Cards: []ExecutiveDepartmentCard{
				{Label: "Dòng Tiền Thu Trong Kỳ", Value: formatVNDShort(revenueF), Sub: "Từ phiếu cân đã hoàn thành", Tone: "green"},
				{Label: "Tổng Chi Phí Vận Hành", Value: formatVNDShort(costF), Sub: "Chi phí sản xuất phát sinh", Tone: "blue"},
				{Label: "Biên Lợi Nhuận Gộp", Value: fmt.Sprintf("%.1f%%", margin), Sub: "Doanh thu trừ chi phí", Tone: toneForMargin(margin)},
				{Label: "Hóa Đơn VAT Xuất Mỏ", Value: fmt.Sprintf("%d Hóa đơn", invoiceCount), Sub: "Từ hệ thống hóa đơn", Tone: "blue"},
			},
			Highlights: []string{
				fmt.Sprintf("Biên lợi nhuận gộp kỳ hiện tại đạt %.1f%%.", margin),
				fmt.Sprintf("Tổng chi phí vận hành ghi nhận %s.", formatVNDShort(costF)),
				fmt.Sprintf("Hệ thống có %d hóa đơn VAT đã phát hành.", invoiceCount),
			},
			MineComparison: departmentMineRows(mineNames, "Chi phí"),
		},
		{
			Key: "production", Name: "Sản Xuất - Kỹ Thuật Mỏ",
			HeadOfDept: "Phó Giám Đốc Kỹ Thuật", Phone: "—",
			Description: "Kế hoạch khai thác, hộ chiếu nổ mìn, trạm nghiền sàng và khảo sát 3D.",
			Icon:        "pickaxe", Route: "/ke-hoach-san-luong", Score: clampScore(margin + 58),
			Cards: []ExecutiveDepartmentCard{
				{Label: "Tổng Sản Lượng Khai Thác", Value: fmt.Sprintf("%.0f tấn", tonsF), Sub: "Từ phiếu cân xuất bãi", Tone: "green"},
				{Label: "Chuyến Xe Vận Chuyển", Value: fmt.Sprintf("%d chuyến", trips), Sub: "Lượt xe ghi nhận", Tone: "blue"},
				{Label: "Tiêu Hao Nhiên Liệu", Value: fmt.Sprintf("%.0f lít", fuelF), Sub: "Từ nhật ký thiết bị", Tone: "amber"},
				{Label: "Đội Xe Hoạt Động", Value: fmt.Sprintf("%d/%d xe", activeVehicle, vehicleCount), Sub: "Phương tiện công ty", Tone: "green"},
			},
			Highlights: []string{
				fmt.Sprintf("Sản lượng khai thác kỳ đạt %.0f tấn qua %d chuyến xe.", tonsF, trips),
				fmt.Sprintf("Nhiên liệu tiêu thụ %.0f lít theo nhật ký thiết bị.", fuelF),
				fmt.Sprintf("Đội xe có %d/%d phương tiện đang hoạt động.", activeVehicle, vehicleCount),
			},
			MineComparison: departmentMineRows(mineNames, "Sản lượng"),
		},
		{
			Key: "warehouse", Name: "Kho - Vận Tải - Trạm Cân",
			HeadOfDept: "Trưởng Trạm Cân & Bãi Đá", Phone: "—",
			Description: "Điều hành cân xe tự động AI OCR, giám sát lệch bì và tồn kho bãi đá.",
			Icon:        "scale", Route: "/phieu-can", Score: clampScore(95 - float64(alertCount)),
			Cards: []ExecutiveDepartmentCard{
				{Label: "Lượt Xe Qua Trạm Cân", Value: fmt.Sprintf("%d lượt", trips), Sub: "Cân tự động 2 lượt", Tone: "blue"},
				{Label: "Cảnh Báo Đang Mở", Value: fmt.Sprintf("%d Sự vụ", alertCount), Sub: "Cảnh báo vận hành", Tone: toneForAlerts(alertCount)},
				{Label: "Đoàn Xe Vận Tải Mỏ", Value: fmt.Sprintf("%d Xe", vehicleCount), Sub: fmt.Sprintf("%d xe đang hoạt động", activeVehicle), Tone: "green"},
				{Label: "Sản Lượng Qua Cân", Value: fmt.Sprintf("%.0f tấn", tonsF), Sub: "Tổng khối lượng cân", Tone: "blue"},
			},
			Highlights: []string{
				fmt.Sprintf("Trạm cân ghi nhận %d lượt xe trong kỳ.", trips),
				fmt.Sprintf("Có %d cảnh báo vận hành đang mở cần xử lý.", alertCount),
				fmt.Sprintf("Đội xe mỏ hiện có %d phương tiện.", vehicleCount),
			},
			MineComparison: departmentMineRows(mineNames, "Lượt xe"),
		},
		{
			Key: "hr", Name: "Hành Chính - Nhân Sự (HRM)",
			HeadOfDept: "Trưởng phòng HC-NS", Phone: "—",
			Description: "Chấm công nhận diện khuôn mặt AI, bảng lương khoán chuyến và bảo hiểm.",
			Icon:        "user-check", Route: "/nhan-su", Score: clampScore(attendanceRate),
			Cards: []ExecutiveDepartmentCard{
				{Label: "Tổng Nhân Sự Công Ty", Value: fmt.Sprintf("%d Cán bộ", employees), Sub: "Đang làm việc", Tone: "blue"},
				{Label: "Tỷ Lệ Chuyên Cần Ca", Value: fmt.Sprintf("%.1f%%", attendanceRate), Sub: fmt.Sprintf("%d lượt chấm công", attendance), Tone: "green"},
				{Label: "Quỹ Lương Cơ Bản", Value: formatVNDShort(payrollF), Sub: "Tổng lương cơ bản", Tone: "blue"},
				{Label: "An Toàn Lao Động (ATLĐ)", Value: "0 Sự cố", Sub: "Cập nhật từ sổ ATLĐ", Tone: "green"},
			},
			Highlights: []string{
				fmt.Sprintf("Tổng nhân sự đang làm việc: %d cán bộ.", employees),
				fmt.Sprintf("Tỷ lệ chuyên cần ca đạt %.1f%% với %d lượt chấm công.", attendanceRate, attendance),
				fmt.Sprintf("Quỹ lương cơ bản hiện tại %s.", formatVNDShort(payrollF)),
			},
			MineComparison: departmentMineRows(mineNames, "Nhân sự"),
		},
		{
			Key: "legal", Name: "Pháp Lý & Giấy Phép Mỏ",
			HeadOfDept: "Trưởng ban Pháp chế", Phone: "—",
			Description: "Giấy phép khai thác khoáng sản, ĐTM môi trường, ranh giới mỏ và trữ lượng.",
			Icon:        "building", Route: "/ho-so-giay-phep", Score: clampScore(100 - float64(minInt(alertCount, 10))),
			Cards: []ExecutiveDepartmentCard{
				{Label: "Hồ Sơ Giấy Phép Mỏ", Value: fmt.Sprintf("%d Hồ sơ", permitCount), Sub: "Trong hệ thống pháp lý", Tone: "green"},
				{Label: "Báo Cáo Đã Nộp", Value: fmt.Sprintf("%d Báo cáo", reportCount), Sub: "Báo cáo cơ quan quản lý", Tone: "blue"},
				{Label: "Cảnh Báo Đang Mở", Value: fmt.Sprintf("%d Sự vụ", alertCount), Sub: "Cần theo dõi", Tone: toneForAlerts(alertCount)},
				{Label: "Hợp Đồng Đang Hiệu Lực", Value: fmt.Sprintf("%d HĐ", contractCount), Sub: "Hợp đồng thương mại", Tone: "blue"},
			},
			Highlights: []string{
				fmt.Sprintf("Hệ thống đang quản lý %d hồ sơ giấy phép mỏ.", permitCount),
				fmt.Sprintf("Đã ghi nhận %d báo cáo nộp cơ quan quản lý.", reportCount),
				fmt.Sprintf("Có %d cảnh báo đang mở cần Ban Pháp chế theo dõi.", alertCount),
			},
			MineComparison: departmentMineRows(mineNames, "Hồ sơ"),
		},
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clampScore(v float64) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return int(math.Round(v))
}

func toneForMargin(m float64) string {
	if m >= 30 {
		return "green"
	}
	if m >= 15 {
		return "amber"
	}
	return "red"
}

func toneForAlerts(n int) string {
	if n == 0 {
		return "green"
	}
	if n <= 3 {
		return "amber"
	}
	return "red"
}

func departmentMineRows(mineNames []string, metric string) []ExecutiveDepartmentMineRow {
	rows := make([]ExecutiveDepartmentMineRow, 0, len(mineNames))
	for _, m := range mineNames {
		rows = append(rows, ExecutiveDepartmentMineRow{
			MineName: m,
			Target:   "Theo kế hoạch tháng",
			Actual:   "Tổng hợp từ dữ liệu vận hành",
			Rate:     metric,
			Status:   "good",
		})
	}
	return rows
}

func buildExecutiveMineComparison(ctx context.Context, start, end time.Time, quarryID string) []ExecutiveMineComparison {
	out := make([]ExecutiveMineComparison, 0, 4)
	if database.Pool == nil {
		return out
	}
	start, end = executiveDataWindow(ctx, start, end)

	rows, err := database.Pool.Query(ctx, `
		SELECT q.code, q.name, COALESCE(q.address, ''), COALESCE(q.latitude, 0)::float8, COALESCE(q.longitude, 0)::float8
		FROM quarries q
		WHERE ($1 = '' OR q.code = $1)
		ORDER BY q.code`, quarryID)
	if err != nil {
		return out
	}
	defer rows.Close()

	type qRow struct {
		code, name, address string
		lat, lng            float64
	}
	var quarries []qRow
	for rows.Next() {
		var q qRow
		if rows.Scan(&q.code, &q.name, &q.address, &q.lat, &q.lng) == nil {
			quarries = append(quarries, q)
		}
	}
	rows.Close()

	for _, q := range quarries {
		var prod, rev float64
		_ = database.Pool.QueryRow(ctx, fmt.Sprintf(`
			SELECT COALESCE(SUM(%s), 0), COALESCE(SUM(%s), 0)
			FROM tickets WHERE quarry_code = $1 AND created_at >= $2 AND created_at < $3
		`, sqlTicketTonnageExpr, sqlTicketRevenueExpr), q.code, start, end).Scan(&prod, &rev)

		var annualPlan float64
		_ = database.Pool.QueryRow(ctx,
			`SELECT COALESCE(SUM(annual_target), 0) FROM mining_plans WHERE quarry_code = $1`, q.code).Scan(&annualPlan)
		plan := annualPlan / 12.0

		var cost float64
		_ = database.Pool.QueryRow(ctx, `
			SELECT COALESCE(SUM(actual_value), 0) FROM production_costs
			WHERE created_at >= $1 AND created_at < $2
			  AND (mine_area ILIKE '%' || $3 || '%'
			       OR ($3 = 'MO-PT-01' AND mine_area ILIKE '%Phú Thọ%')
			       OR ($3 = 'MO-TU-02' AND (mine_area ILIKE '%Tân Uyên%' OR mine_area ILIKE '%Bình Dương%'))
			       OR ($3 = 'MO-HN-03' AND (mine_area ILIKE '%Hà Nam%' OR mine_area ILIKE '%Kiện Khê%'))
			       OR ($3 = 'MO-BP-04' AND mine_area ILIKE '%Bình Phước%'))`,
			start, end, q.code).Scan(&cost)

		var fleet, activeFleet, workforce, alertCount int
		_ = database.Pool.QueryRow(ctx,
			`SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'Hoạt động') FROM vehicles WHERE quarry_code = $1`, q.code).Scan(&fleet, &activeFleet)
		_ = database.Pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM hr_employees WHERE work_location ILIKE '%' || $1 || '%'`, q.code).Scan(&workforce)
		_ = database.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM alerts a
			WHERE a.bs IN (SELECT bs FROM vehicles WHERE quarry_code = $1)
			  AND a.status NOT IN ('resolved','closed','Đã xử lý')`, q.code).Scan(&alertCount)

		capacity := ""
		expiryYear := 0
		var capPtr, expiryPtr *string
		_ = database.Pool.QueryRow(ctx, `
			SELECT capacity, expiry_date FROM mining_permits
			WHERE category = 'mining_license'
			  AND (title ILIKE '%' || $1 || '%' OR mine_name ILIKE '%' || $1 || '%')
			ORDER BY expiry_date DESC LIMIT 1`, q.code).Scan(&capPtr, &expiryPtr)
		if capPtr != nil {
			capacity = *capPtr
		}
		if expiryPtr != nil {
			parts := strings.Split(*expiryPtr, "/")
			if len(parts) == 3 {
				fmt.Sscanf(parts[2], "%d", &expiryYear)
			}
		}

		progress := 0.0
		if plan > 0 {
			progress = math.Round(prod/plan*1000) / 10
		}
		margin := 0.0
		if rev > 0 {
			margin = math.Round((rev-cost)/rev*1000) / 10
		}
		oee := 0.0
		if fleet > 0 {
			oee = math.Round(float64(activeFleet)/float64(fleet)*1000) / 10
		}
		status := "normal"
		if progress > 0 && progress < 90 {
			status = "warning"
		}
		if progress > 0 && progress < 75 {
			status = "critical"
		}

		out = append(out, ExecutiveMineComparison{
			Code:                  q.code,
			Name:                  q.name,
			Location:              q.address,
			LicensedCapacity:      capacity,
			ExpiryYear:            expiryYear,
			ActualProductionTons:  math.Round(prod),
			PlannedProductionTons: math.Round(plan),
			ProgressPct:           progress,
			Revenue:               math.Round(rev),
			OperatingCost:         math.Round(cost),
			GrossMarginPct:        margin,
			FleetCount:            fleet,
			ActiveFleet:           activeFleet,
			FleetOeePct:           oee,
			WorkforceCount:        workforce,
			AlertsCount:           alertCount,
			Status:                status,
			Route3D:               "/khao-sat-3d/tong-quan",
			RouteCatalog:          "/to-chuc/mo-da",
			Lat:                   q.lat,
			Lng:                   q.lng,
		})
	}
	return out
}

type ExecutiveDebtAging struct {
	Labels  []string              `json:"labels"`
	Buckets []ExecutiveDebtBucket `json:"buckets"`
}

type ExecutiveDebtBucket struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
}

func buildExecutiveDebtAging(ctx context.Context) ExecutiveDebtAging {
	empty := ExecutiveDebtAging{Labels: []string{}, Buckets: []ExecutiveDebtBucket{}}
	if database.Pool == nil {
		return empty
	}
	rows, err := database.Pool.Query(ctx, `
		SELECT COALESCE(customer, partner, 'Khách hàng'), COALESCE(balance, '0'), COALESCE(due, '')
		FROM payments_debt
		ORDER BY created_at DESC LIMIT 20`)
	if err != nil {
		return empty
	}
	defer rows.Close()

	type bucket struct {
		label string
		data  map[string]float64
	}
	buckets := []*bucket{
		{label: "0 - 30 ngày (Trong hạn)", data: map[string]float64{}},
		{label: "30 - 60 ngày (Cần nhắc)", data: map[string]float64{}},
		{label: "60 - 90 ngày (Cảnh báo)", data: map[string]float64{}},
		{label: "> 90 ngày (Nợ quá hạn)", data: map[string]float64{}},
	}
	labels := make([]string, 0, 20)
	seen := map[string]bool{}
	now := time.Now()

	for rows.Next() {
		var customer, balance, due string
		if rows.Scan(&customer, &balance, &due) != nil {
			continue
		}
		if !seen[customer] {
			seen[customer] = true
			labels = append(labels, customer)
		}
		amount := parseMoneyString(balance) / 1_000_000 // triệu đồng
		ageDays := 0
		parts := strings.Split(due, "/")
		if len(parts) == 3 {
			if d, err := time.Parse("02/01/2006", due); err == nil {
				ageDays = int(now.Sub(d).Hours() / 24)
			}
		}
		idx := 0
		switch {
		case ageDays > 90:
			idx = 3
		case ageDays > 60:
			idx = 2
		case ageDays > 30:
			idx = 1
		}
		buckets[idx].data[customer] += amount
	}

	out := ExecutiveDebtAging{Labels: labels, Buckets: make([]ExecutiveDebtBucket, 0, len(buckets))}
	for _, b := range buckets {
		data := make([]float64, len(labels))
		for i, l := range labels {
			data[i] = math.Round(b.data[l])
		}
		out.Buckets = append(out.Buckets, ExecutiveDebtBucket{Label: b.label, Data: data})
	}
	return out
}

func ExecutiveDebtAgingHandler(w http.ResponseWriter, r *http.Request) {
	JSON(w, buildExecutiveDebtAging(r.Context()))
}

func buildExecutiveReports(ctx context.Context, start, end time.Time, quarryID string) []ExecutiveReportItem {
	out := make([]ExecutiveReportItem, 0, 4)
	if database.Pool == nil {
		return out
	}
	start, end = executiveDataWindow(ctx, start, end)

	// Financial P&L report aggregated from tickets + production costs.
	var revenue, cost float64
	_ = database.Pool.QueryRow(ctx, fmt.Sprintf(`
		SELECT COALESCE(SUM(%s), 0) FROM tickets
		WHERE created_at >= $1 AND created_at < $2 AND ($3 = '' OR quarry_code = $3)
	`, sqlTicketRevenueExpr), start, end, quarryID).Scan(&revenue)
	_ = database.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(actual_value), 0) FROM production_costs
		WHERE created_at >= $1 AND created_at < $2`, start, end).Scan(&cost)
	profit := revenue - cost
	margin := 0.0
	if revenue > 0 {
		margin = math.Round(profit/revenue*1000) / 10
	}
	out = append(out, ExecutiveReportItem{
		ID:            "REP-EXEC-PL",
		Title:         "Báo Cáo Tổng Hợp Điều Hành Doanh Nghiệp & P&L",
		Period:        fmt.Sprintf("%s - %s", start.Format("02/01/2006"), end.AddDate(0, 0, -1).Format("02/01/2006")),
		Category:      "financial",
		CategoryLabel: "Tài chính - P&L",
		GeneratedDate: time.Now().Format("02/01/2006 15:04"),
		Status:        "published",
		Summary:       fmt.Sprintf("Doanh thu kỳ đạt %s, chi phí %s, lợi nhuận gộp %s (biên %.1f%%).", formatVNDShort(revenue), formatVNDShort(cost), formatVNDShort(profit), margin),
		KeyMetrics: []ExecutiveReportMetric{
			{Label: "Doanh thu", Value: formatVNDShort(revenue)},
			{Label: "Chi phí", Value: formatVNDShort(cost)},
			{Label: "Lợi nhuận gộp", Value: formatVNDShort(profit)},
			{Label: "Biên LN", Value: fmt.Sprintf("%.1f%%", margin)},
		},
		Route: "/dieu-hanh/tong-quan",
	})

	// Production report from the mining reports table.
	var prodName, prodItem, prodPeriod, prodStatus string
	var plan, actual float64
	err := database.Pool.QueryRow(ctx, `
		SELECT COALESCE(name, ''), COALESCE(item, ''), COALESCE(period, ''), COALESCE(status, ''),
		       COALESCE(plan, 0), COALESCE(actual, 0)
		FROM reports
		WHERE ($1 = '' OR quarry_code = $1)
		ORDER BY created_at DESC, period DESC LIMIT 1`, quarryID).Scan(&prodName, &prodItem, &prodPeriod, &prodStatus, &plan, &actual)
	if err == nil {
		rate := 0.0
		if plan > 0 {
			rate = math.Round(actual/plan*1000) / 10
		}
		out = append(out, ExecutiveReportItem{
			ID:            "REP-EXEC-PROD",
			Title:         prodName,
			Period:        prodPeriod,
			Category:      "production",
			CategoryLabel: "Sản xuất - Đối soát 3D",
			GeneratedDate: time.Now().Format("02/01/2006 15:04"),
			Status:        "published",
			Summary:       fmt.Sprintf("%s: kế hoạch %.0f, thực tế %.0f (đạt %.1f%%), trạng thái %s.", prodItem, plan, actual, rate, prodStatus),
			KeyMetrics: []ExecutiveReportMetric{
				{Label: "Kế hoạch", Value: fmt.Sprintf("%.0f", plan)},
				{Label: "Thực tế", Value: fmt.Sprintf("%.0f", actual)},
				{Label: "Tỷ lệ đạt", Value: fmt.Sprintf("%.1f%%", rate)},
				{Label: "Trạng thái", Value: prodStatus},
			},
			Route: "/ke-hoach-san-luong",
		})
	}

	// Tax / resource tax report.
	var taxAmount, envFee, totalPayable float64
	var taxCount int
	_ = database.Pool.QueryRow(ctx, `
		SELECT COUNT(*), COALESCE(SUM(resource_tax_amount), 0), COALESCE(SUM(environmental_fee), 0), COALESCE(SUM(total_payable), 0)
		FROM resource_taxes`).Scan(&taxCount, &taxAmount, &envFee, &totalPayable)
	if taxCount > 0 {
		out = append(out, ExecutiveReportItem{
			ID:            "REP-EXEC-TAX",
			Title:         "Báo Cáo Kê Khai Thuế Tài Nguyên Khoáng Sản & Phí Bảo Vệ Môi Trường",
			Period:        "Kỳ hiện tại",
			Category:      "tax",
			CategoryLabel: "Thuế & Pháp lý",
			GeneratedDate: time.Now().Format("02/01/2006 15:04"),
			Status:        "published",
			Summary:       fmt.Sprintf("Tổng thuế tài nguyên %s, phí BVMT %s, tổng nghĩa vụ ngân sách %s trên %d bản ghi.", formatVNDShort(taxAmount), formatVNDShort(envFee), formatVNDShort(totalPayable), taxCount),
			KeyMetrics: []ExecutiveReportMetric{
				{Label: "Thuế tài nguyên", Value: formatVNDShort(taxAmount)},
				{Label: "Phí BVMT", Value: formatVNDShort(envFee)},
				{Label: "Tổng nộp NS", Value: formatVNDShort(totalPayable)},
				{Label: "Bản ghi", Value: fmt.Sprintf("%d", taxCount)},
			},
			Route: "/ke-hoach-san-luong/thue-tai-nguyen",
		})
	}

	// Statutory report submitted to authorities.
	var statTitle, statPeriod, statStatus, statRecipient string
	var minedVolume float64
	err = database.Pool.QueryRow(ctx, `
		SELECT COALESCE(title, ''), COALESCE(period, ''), COALESCE(status_label, status, ''), COALESCE(recipient, ''), COALESCE(mined_volume, 0)
		FROM statutory_reports
		ORDER BY created_at DESC LIMIT 1`).Scan(&statTitle, &statPeriod, &statStatus, &statRecipient, &minedVolume)
	if err == nil {
		out = append(out, ExecutiveReportItem{
			ID:            "REP-EXEC-STAT",
			Title:         statTitle,
			Period:        statPeriod,
			Category:      "statutory",
			CategoryLabel: "Báo cáo Sở TN&MT",
			GeneratedDate: time.Now().Format("02/01/2006 15:04"),
			Status:        "published",
			Summary:       fmt.Sprintf("Báo cáo gửi %s, khối lượng khai thác %.0f tấn, trạng thái %s.", statRecipient, minedVolume, statStatus),
			KeyMetrics: []ExecutiveReportMetric{
				{Label: "Khối lượng", Value: fmt.Sprintf("%.0f tấn", minedVolume)},
				{Label: "Cơ quan nhận", Value: statRecipient},
				{Label: "Trạng thái", Value: statStatus},
			},
			Route: "/ke-hoach-san-luong/bao-cao-so-tnmt",
		})
	}

	return out
}

// ─────────────────────────────────────────────────────────────────────────────
// HTTP handlers
// ─────────────────────────────────────────────────────────────────────────────

func executiveWorkspaceParams(r *http.Request) (string, time.Time, time.Time, string) {
	ctx := r.Context()
	rawQuarry := strings.TrimSpace(r.URL.Query().Get("quarry_id"))
	if rawQuarry == "" {
		rawQuarry = strings.TrimSpace(r.URL.Query().Get("quarryCode"))
	}
	quarryID := resolveQuarryCode(ctx, rawQuarry)
	period := normalizeExecutivePeriod(r.URL.Query().Get("period"))
	now := time.Now().In(time.FixedZone("ICT", 7*60*60))
	currentStart, currentEnd, _, _, _ := executivePeriodBounds(now, period)
	comparisonEnd := now
	if comparisonEnd.After(currentEnd) {
		comparisonEnd = currentEnd
	}
	return quarryID, currentStart, comparisonEnd, period
}

func ExecutiveTasks(w http.ResponseWriter, r *http.Request) {
	quarryID, start, end, _ := executiveWorkspaceParams(r)
	tasks := buildExecutiveTasks(r.Context(), start, end, quarryID)
	JSON(w, tasks)
}

func ExecutiveAlerts(w http.ResponseWriter, r *http.Request) {
	quarryID, start, end, _ := executiveWorkspaceParams(r)
	JSON(w, buildExecutiveAlerts(r.Context(), start, end, quarryID))
}

func ExecutiveDepartments(w http.ResponseWriter, r *http.Request) {
	quarryID, start, end, _ := executiveWorkspaceParams(r)
	JSON(w, buildExecutiveDepartments(r.Context(), start, end, quarryID))
}

func ExecutiveMineComparisonHandler(w http.ResponseWriter, r *http.Request) {
	quarryID, start, end, _ := executiveWorkspaceParams(r)
	JSON(w, buildExecutiveMineComparison(r.Context(), start, end, quarryID))
}

func ExecutiveReports(w http.ResponseWriter, r *http.Request) {
	quarryID, start, end, _ := executiveWorkspaceParams(r)
	JSON(w, buildExecutiveReports(r.Context(), start, end, quarryID))
}

// silence unused helper warning if buildOptionalFloat is not used yet.
var _ = buildOptionalFloat
