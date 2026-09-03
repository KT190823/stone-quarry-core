package main

import (
	"fmt"
	"log"
	"net/http"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/handlers"
	"mo-da-backend/internal/middleware"
)

func main() {
	database.Connect()
	defer database.Close()

	database.Migrate()
	database.Seed()

	mux := http.NewServeMux()

	ticketH := handlers.NewTicketHandler()
	vehicleH := handlers.NewVehicleHandler()
	alertH := handlers.NewAlertHandler()
	catalogH := handlers.NewCatalogHandler()
	hrH := handlers.NewHrHandler()
	hrXH := handlers.NewHrExtendedHandler()
	miningH := handlers.NewMiningHandler()
	inventoryH := handlers.NewInventoryHandler()
	paymentH := handlers.NewPaymentHandler()
	userH := handlers.NewUserHandler()
	cameraH := handlers.NewCameraHandler()
	printH := handlers.NewPrintHandler()
	quarryH := handlers.NewQuarryHandler()
	integrationH := handlers.NewIntegrationHandler()

	// Register Quarry Core & 3rd-Party Integration Gateway
	quarryH.Register(mux)
	integrationH.Register(mux)

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.JSON(w, map[string]string{"status": "ok", "service": "mo-da-backend"})
	})

	// Print templates (configurable document/report layouts)
	mux.HandleFunc("GET /api/print/templates", printH.List)
	mux.HandleFunc("GET /api/print/templates/{id}", printH.Get)
	mux.HandleFunc("POST /api/print/templates", printH.Create)
	mux.HandleFunc("PUT /api/print/templates/{id}", printH.Update)
	mux.HandleFunc("DELETE /api/print/templates/{id}", printH.Delete)
	mux.HandleFunc("POST /api/print/templates/{id}/default", printH.SetDefault)

	// Auth (login / logout / current user)
	authH := handlers.NewAuthHandler()
	mux.HandleFunc("POST /api/auth/login", authH.Login)
	mux.HandleFunc("POST /api/auth/logout", authH.Logout)
	mux.HandleFunc("GET /api/auth/me", authH.Me)

	mux.HandleFunc("GET /api/dashboard", handlers.DashboardHandler)

	mux.HandleFunc("GET /api/tickets", ticketH.List)
	mux.HandleFunc("GET /api/tickets/{id}", ticketH.Get)
	mux.HandleFunc("POST /api/tickets", ticketH.Create)
	mux.HandleFunc("PUT /api/tickets/{id}", ticketH.Update)
	mux.HandleFunc("DELETE /api/tickets/{id}", ticketH.Delete)

	mux.HandleFunc("GET /api/vehicles", vehicleH.List)
	mux.HandleFunc("GET /api/vehicles/{id}", vehicleH.Get)
	mux.HandleFunc("POST /api/vehicles", vehicleH.Create)
	mux.HandleFunc("PUT /api/vehicles/{id}", vehicleH.Update)
	mux.HandleFunc("DELETE /api/vehicles/{id}", vehicleH.Delete)

	mux.HandleFunc("GET /api/alerts", alertH.List)
	mux.HandleFunc("GET /api/alerts/{id}", alertH.Get)
	mux.HandleFunc("POST /api/alerts", alertH.Create)
	mux.HandleFunc("PUT /api/alerts/{id}", alertH.Update)
	mux.HandleFunc("DELETE /api/alerts/{id}", alertH.Delete)

	mux.HandleFunc("GET /api/catalog/materials", catalogH.ListMaterials)
	mux.HandleFunc("GET /api/catalog/materials/{id}", catalogH.GetMaterial)
	mux.HandleFunc("POST /api/catalog/materials", catalogH.CreateMaterial)
	mux.HandleFunc("PUT /api/catalog/materials/{id}", catalogH.UpdateMaterial)
	mux.HandleFunc("DELETE /api/catalog/materials/{id}", catalogH.DeleteMaterial)

	mux.HandleFunc("GET /api/catalog/ticket-types", catalogH.ListTicketTypes)
	mux.HandleFunc("GET /api/catalog/vehicles", catalogH.ListVehicleCatalogs)
	mux.HandleFunc("GET /api/catalog/stations", catalogH.ListStations)
	mux.HandleFunc("GET /api/catalog/equipment", catalogH.ListEquipment)

	mux.HandleFunc("GET /api/inventory/inbound", inventoryH.ListInbound)
	mux.HandleFunc("GET /api/inventory/outbound", inventoryH.ListOutbound)
	mux.HandleFunc("GET /api/inventory/stocktake", inventoryH.ListStocktake)
	mux.HandleFunc("GET /api/inventory/movement", inventoryH.ListMovements)

	mux.HandleFunc("GET /api/payments/invoices", paymentH.ListInvoices)
	mux.HandleFunc("GET /api/payments/debt", paymentH.ListDebt)
	mux.HandleFunc("GET /api/payments/reconcile", paymentH.ListReconcile)

	mux.HandleFunc("GET /api/users", userH.ListUsers)
	mux.HandleFunc("GET /api/users/{id}", userH.GetUser)
	mux.HandleFunc("POST /api/users", userH.CreateUser)
	mux.HandleFunc("PUT /api/users/{id}", userH.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", userH.DeleteUser)
	mux.HandleFunc("GET /api/users/roles", userH.ListRoles)
	mux.HandleFunc("GET /api/users/logs", userH.ListLogs)

	mux.HandleFunc("GET /api/reports", userH.ListReports)
	mux.HandleFunc("GET /api/reports/{id}", userH.GetReport)
	mux.HandleFunc("POST /api/reports", userH.CreateReport)
	mux.HandleFunc("PUT /api/reports/{id}", userH.UpdateReport)
	mux.HandleFunc("DELETE /api/reports/{id}", userH.DeleteReport)

	mux.HandleFunc("GET /api/settings", userH.ListSettings)
	mux.HandleFunc("GET /api/settings/{id}", userH.GetSetting)
	mux.HandleFunc("POST /api/settings", userH.CreateSetting)
	mux.HandleFunc("PUT /api/settings/{id}", userH.UpdateSetting)
	mux.HandleFunc("DELETE /api/settings/{id}", userH.DeleteSetting)

	mux.HandleFunc("GET /api/camera", cameraH.ListCameras)
	mux.HandleFunc("GET /api/camera/{id}", cameraH.GetCamera)
	mux.HandleFunc("POST /api/camera", cameraH.CreateCamera)
	mux.HandleFunc("PUT /api/camera/{id}", cameraH.UpdateCamera)
	mux.HandleFunc("DELETE /api/camera/{id}", cameraH.DeleteCamera)

	mux.HandleFunc("GET /api/permits", miningH.ListPermits)
	mux.HandleFunc("GET /api/permits/{id}", miningH.GetPermit)
	mux.HandleFunc("POST /api/permits", miningH.CreatePermit)
	mux.HandleFunc("PUT /api/permits/{id}", miningH.UpdatePermit)
	mux.HandleFunc("DELETE /api/permits/{id}", miningH.DeletePermit)

	mux.HandleFunc("GET /api/mining-plan", miningH.ListPlans)
	mux.HandleFunc("GET /api/mining-plan/{id}", miningH.GetPlan)
	mux.HandleFunc("POST /api/mining-plan", miningH.CreatePlan)
	mux.HandleFunc("PUT /api/mining-plan/{id}", miningH.UpdatePlan)
	mux.HandleFunc("DELETE /api/mining-plan/{id}", miningH.DeletePlan)
	mux.HandleFunc("GET /api/mining-plan/production-chain", miningH.ListProductionStages)
	mux.HandleFunc("GET /api/mining-plan/statutory-reports", miningH.ListStatutoryReports)
	mux.HandleFunc("GET /api/mining-plan/resource-tax", miningH.ListResourceTaxes)

	mux.HandleFunc("GET /api/mining-operations/blasting", miningH.ListBlasting)
	mux.HandleFunc("GET /api/mining-operations/blasting/{id}", miningH.GetBlasting)
	mux.HandleFunc("POST /api/mining-operations/blasting", miningH.CreateBlasting)
	mux.HandleFunc("PUT /api/mining-operations/blasting/{id}", miningH.UpdateBlasting)
	mux.HandleFunc("DELETE /api/mining-operations/blasting/{id}", miningH.DeleteBlasting)
	mux.HandleFunc("GET /api/mining-operations/crusher-plant", miningH.ListCrusherPlants)
	mux.HandleFunc("GET /api/mining-operations/equipment-fuel", miningH.ListEquipmentFuel)
	mux.HandleFunc("GET /api/mining-operations/equipment-fuel/{id}", miningH.GetEquipmentFuel)
	mux.HandleFunc("POST /api/mining-operations/equipment-fuel", miningH.CreateEquipmentFuel)
	mux.HandleFunc("PUT /api/mining-operations/equipment-fuel/{id}", miningH.UpdateEquipmentFuel)
	mux.HandleFunc("DELETE /api/mining-operations/equipment-fuel/{id}", miningH.DeleteEquipmentFuel)

	mux.HandleFunc("GET /api/partners/customers", catalogH.ListCustomers)
	mux.HandleFunc("GET /api/partners/customers/{id}", catalogH.GetCustomer)
	mux.HandleFunc("POST /api/partners/customers", catalogH.CreateCustomer)
	mux.HandleFunc("PUT /api/partners/customers/{id}", catalogH.UpdateCustomer)
	mux.HandleFunc("DELETE /api/partners/customers/{id}", catalogH.DeleteCustomer)

	mux.HandleFunc("GET /api/partners/suppliers", catalogH.ListSuppliers)
	mux.HandleFunc("GET /api/partners/suppliers/{id}", catalogH.GetSupplier)
	mux.HandleFunc("POST /api/partners/suppliers", catalogH.CreateSupplier)
	mux.HandleFunc("PUT /api/partners/suppliers/{id}", catalogH.UpdateSupplier)
	mux.HandleFunc("DELETE /api/partners/suppliers/{id}", catalogH.DeleteSupplier)

	mux.HandleFunc("GET /api/hr/dashboard", hrH.Dashboard)
	mux.HandleFunc("GET /api/hr/employees", hrH.ListEmployees)
	mux.HandleFunc("GET /api/hr/employees/{id}", hrH.GetEmployee)
	mux.HandleFunc("POST /api/hr/employees", hrH.CreateEmployee)
	mux.HandleFunc("PUT /api/hr/employees/{id}", hrH.UpdateEmployee)
	mux.HandleFunc("DELETE /api/hr/employees/{id}", hrH.DeleteEmployee)
	mux.HandleFunc("GET /api/hr/contracts", hrH.ListContracts)
	mux.HandleFunc("GET /api/hr/contracts/{id}", hrH.GetContract)
	mux.HandleFunc("POST /api/hr/contracts", hrH.CreateContract)
	mux.HandleFunc("PUT /api/hr/contracts/{id}", hrH.UpdateContract)
	mux.HandleFunc("DELETE /api/hr/contracts/{id}", hrH.DeleteContract)
	mux.HandleFunc("GET /api/hr/shifts", hrH.ListShifts)
	mux.HandleFunc("GET /api/hr/shifts/{id}", hrH.GetShift)
	mux.HandleFunc("POST /api/hr/shifts", hrH.CreateShift)
	mux.HandleFunc("PUT /api/hr/shifts/{id}", hrH.UpdateShift)
	mux.HandleFunc("DELETE /api/hr/shifts/{id}", hrH.DeleteShift)
	mux.HandleFunc("GET /api/hr/shift-schedule", hrH.ListShiftSchedules)
	mux.HandleFunc("GET /api/hr/shift-schedule/{id}", hrH.GetShiftSchedule)
	mux.HandleFunc("POST /api/hr/shift-schedule", hrH.CreateShiftSchedule)
	mux.HandleFunc("PUT /api/hr/shift-schedule/{id}", hrH.UpdateShiftSchedule)
	mux.HandleFunc("DELETE /api/hr/shift-schedule/{id}", hrH.DeleteShiftSchedule)
	mux.HandleFunc("GET /api/hr/attendance", hrH.ListAttendances)
	mux.HandleFunc("GET /api/hr/attendance/{id}", hrH.GetAttendance)
	mux.HandleFunc("POST /api/hr/attendance", hrH.CreateAttendance)
	mux.HandleFunc("PUT /api/hr/attendance/{id}", hrH.UpdateAttendance)
	mux.HandleFunc("DELETE /api/hr/attendance/{id}", hrH.DeleteAttendance)
	mux.HandleFunc("GET /api/hr/leaves", hrH.ListLeaves)
	mux.HandleFunc("GET /api/hr/leaves/{id}", hrH.GetLeave)
	mux.HandleFunc("POST /api/hr/leaves", hrH.CreateLeave)
	mux.HandleFunc("PUT /api/hr/leaves/{id}", hrH.UpdateLeave)
	mux.HandleFunc("DELETE /api/hr/leaves/{id}", hrH.DeleteLeave)
	mux.HandleFunc("GET /api/hr/leave-allocations", hrH.ListLeaveAllocations)
	mux.HandleFunc("GET /api/hr/leave-allocations/{id}", hrH.GetLeaveAllocation)
	mux.HandleFunc("POST /api/hr/leave-allocations", hrH.CreateLeaveAllocation)
	mux.HandleFunc("PUT /api/hr/leave-allocations/{id}", hrH.UpdateLeaveAllocation)
	mux.HandleFunc("DELETE /api/hr/leave-allocations/{id}", hrH.DeleteLeaveAllocation)
	mux.HandleFunc("GET /api/hr/timesheets", hrH.ListTimesheets)
	mux.HandleFunc("GET /api/hr/payroll", hrH.ListPayslips)

	// ===== HRM extended modules =====
	hrXH.Register(mux)

	mux.HandleFunc("GET /api/vehicles/gps-fleet", cameraH.ListGpsFleet)
	mux.HandleFunc("GET /api/vehicles/fuel-theft-audit", cameraH.ListFuelTheftAudits)
	mux.HandleFunc("GET /api/vehicles/fuel-norms", cameraH.ListFuelNorms)
	mux.HandleFunc("GET /api/vehicles/yard-checkinout", cameraH.ListYardCheckInOuts)

	mux.HandleFunc("GET /api/google/drive", cameraH.ListGoogleDrive)
	mux.HandleFunc("GET /api/google/gmail", cameraH.ListGoogleGmail)
	mux.HandleFunc("GET /api/google/maps", cameraH.ListGoogleMaps)
	mux.HandleFunc("GET /api/google/photos", cameraH.ListGooglePhotos)

	// ===== New modules from meeting requirements =====

	// Vehicle trips (Camera AI)
	mux.HandleFunc("GET /api/vehicle-trips", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListVehicleTrips(w, r)
	})
	mux.HandleFunc("GET /api/vehicle-trips/stats", func(w http.ResponseWriter, r *http.Request) {
		handlers.VehicleTripStats(w, r)
	})

	// Cost norms
	mux.HandleFunc("GET /api/cost-norms", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListCostNorms(w, r)
	})

	// Production costs
	mux.HandleFunc("GET /api/production-costs", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListProductionCosts(w, r)
	})
	mux.HandleFunc("GET /api/production-costs/summary", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProductionCostSummary(w, r)
	})

	// Delivery confirmations
	mux.HandleFunc("GET /api/deliveries", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListDeliveries(w, r)
	})

	// Tax records
	mux.HandleFunc("GET /api/tax", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListTaxRecords(w, r)
	})

	// Delegations
	mux.HandleFunc("GET /api/delegations", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListDelegations(w, r)
	})

	// Risk alerts
	mux.HandleFunc("GET /api/risk-alerts", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListRiskAlerts(w, r)
	})

	// Geofences
	mux.HandleFunc("GET /api/geofences", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListGeofences(w, r)
	})

	// Authorizations (e-sign)
	mux.HandleFunc("GET /api/authorizations", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListAuthorizations(w, r)
	})

	addr := ":8080"
	fmt.Printf("Backend API running at http://localhost%s\n", addr)
	fmt.Printf("Health check: http://localhost%s/api/health\n", addr)

	handler := middleware.Logger(middleware.CORS(mux))
	log.Fatal(http.ListenAndServe(addr, handler))
}
