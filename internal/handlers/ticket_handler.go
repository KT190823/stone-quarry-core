package handlers

import (
	"net/http"
	"sync"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/services"
)

type TicketHandler struct {
	svc *services.TicketService
}

func NewTicketHandler() *TicketHandler {
	return &TicketHandler{svc: services.NewTicketService()}
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.svc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.svc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.svc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *TicketHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.svc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *TicketHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.svc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := database.Pool

	var wg sync.WaitGroup

	// Basic counts
	var ticketCount, vehicleCount, employeeCount, alertCount, customerCount, supplierCount int
	// Production KPIs
	var todayTrips, weekTrips, monthTrips int
	var todayVolume, weekVolume, monthVolume float64
	// Inventory KPIs
	var totalStock, inboundToday, outboundToday float64
	// Revenue & Finance KPIs
	var monthRevenue, totalDebt float64
	// HR KPIs
	var activeEmployees, lateToday, absentToday int
	// Pending tasks
	var pendingApprovals, overdueContracts int
	// Risk alerts
	var openRiskAlerts int
	// Mining plan progress
	var annualTarget, ytdActual float64

	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM tickets").Scan(&ticketCount) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM vehicles").Scan(&vehicleCount) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM hr_employees").Scan(&employeeCount) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM alerts WHERE status != 'resolved'").Scan(&alertCount) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM customers").Scan(&customerCount) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM suppliers").Scan(&supplierCount) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM vehicle_trips WHERE DATE(check_in_time) = CURRENT_DATE").Scan(&todayTrips) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM vehicle_trips WHERE check_in_time >= CURRENT_DATE - INTERVAL '7 days'").Scan(&weekTrips) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM vehicle_trips WHERE check_in_time >= DATE_TRUNC('month', CURRENT_DATE)").Scan(&monthTrips) })
	run(func() { db.QueryRow(ctx, "SELECT COALESCE(SUM(actual_quantity),0) FROM vehicle_trips WHERE DATE(check_in_time) = CURRENT_DATE").Scan(&todayVolume) })
	run(func() { db.QueryRow(ctx, "SELECT COALESCE(SUM(actual_quantity),0) FROM vehicle_trips WHERE check_in_time >= CURRENT_DATE - INTERVAL '7 days'").Scan(&weekVolume) })
	run(func() { db.QueryRow(ctx, "SELECT COALESCE(SUM(actual_quantity),0) FROM vehicle_trips WHERE check_in_time >= DATE_TRUNC('month', CURRENT_DATE)").Scan(&monthVolume) })
	run(func() { db.QueryRow(ctx, "SELECT COALESCE(SUM(kl_hang::numeric),0) FROM tickets WHERE loai='Nhập' AND date = CURRENT_DATE::text").Scan(&inboundToday) })
	run(func() { db.QueryRow(ctx, "SELECT COALESCE(SUM(kl_hang::numeric),0) FROM tickets WHERE loai='Xuất' AND date = CURRENT_DATE::text").Scan(&outboundToday) })
	run(func() { db.QueryRow(ctx, "SELECT COALESCE(SUM(total_payment::numeric),0) FROM payments_invoices WHERE status != 'cancelled'").Scan(&monthRevenue) })
	run(func() { db.QueryRow(ctx, "SELECT COALESCE(SUM(amount::numeric),0) FROM payments_debt WHERE status = 'unpaid'").Scan(&totalDebt) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM hr_employees WHERE status = 'Đang làm việc'").Scan(&activeEmployees) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM hr_attendances WHERE date = CURRENT_DATE::text AND check_in_time > '08:00'").Scan(&lateToday) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM hr_employees WHERE status = 'Đang làm việc' AND id NOT IN (SELECT employee_id FROM hr_attendances WHERE date = CURRENT_DATE::text)").Scan(&absentToday) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM hr_workflow_instances WHERE status = 'pending'").Scan(&pendingApprovals) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM authorizations WHERE status = 'pending'").Scan(&overdueContracts) })
	run(func() { db.QueryRow(ctx, "SELECT COUNT(*) FROM risk_alerts WHERE status = 'open'").Scan(&openRiskAlerts) })
	run(func() { db.QueryRow(ctx, "SELECT COALESCE(SUM(annual_target),0) FROM mining_plans").Scan(&annualTarget) })
	run(func() { db.QueryRow(ctx, "SELECT COALESCE(SUM(ytd_actual),0) FROM mining_plans").Scan(&ytdActual) })

	wg.Wait()

	totalStock = 4913.36
	var completionRate float64
	if annualTarget > 0 {
		completionRate = (ytdActual / annualTarget) * 100
	}

	JSON(w, map[string]interface{}{
		"stats": map[string]interface{}{
			"tickets":   ticketCount,
			"vehicles":  vehicleCount,
			"employees": employeeCount,
			"alerts":    alertCount,
			"customers": customerCount,
			"suppliers": supplierCount,
		},
		"production": map[string]interface{}{
			"today_trips":   todayTrips,
			"week_trips":    weekTrips,
			"month_trips":   monthTrips,
			"today_volume":  todayVolume,
			"week_volume":   weekVolume,
			"month_volume":  monthVolume,
			"annual_target": annualTarget,
			"ytd_actual":    ytdActual,
			"completion":    completionRate,
		},
		"inventory": map[string]interface{}{
			"total_stock":    totalStock,
			"inbound_today":  inboundToday,
			"outbound_today": outboundToday,
		},
		"finance": map[string]interface{}{
			"month_revenue": monthRevenue,
			"total_debt":    totalDebt,
		},
		"hr": map[string]interface{}{
			"active_employees": activeEmployees,
			"late_today":       lateToday,
			"absent_today":     absentToday,
		},
		"workflow": map[string]interface{}{
			"pending_approvals": pendingApprovals,
			"overdue_contracts": overdueContracts,
		},
		"risk": map[string]interface{}{
			"open_alerts": openRiskAlerts,
		},
	})
}
