package handlers

import (
	"net/http"
	"strconv"

	"mo-da-backend/internal/database"
)

func ListVehicleTrips(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit := 50
	offset := (page - 1) * limit

	rows, err := database.Pool.Query(r.Context(),
		`SELECT id, license_plate, driver_name, camera_id, direction,
		        check_in_time, check_out_time, trip_number,
		        estimated_quantity, actual_quantity, ai_confidence, status
		 FROM vehicle_trips ORDER BY check_in_time DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		JSONError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	var trips []map[string]interface{}
	for rows.Next() {
		var id, plate, driver, cam, dir, status string
		var inTime, outTime interface{}
		var tripNum int
		var eqty, aqty, conf float64
		rows.Scan(&id, &plate, &driver, &cam, &dir, &inTime, &outTime, &tripNum, &eqty, &aqty, &conf, &status)
		trips = append(trips, map[string]interface{}{
			"id": id, "license_plate": plate, "driver_name": driver,
			"camera_id": cam, "direction": dir,
			"check_in_time": inTime, "check_out_time": outTime,
			"trip_number": tripNum, "estimated_quantity": eqty,
			"actual_quantity": aqty, "ai_confidence": conf, "status": status,
		})
	}
	JSON(w, map[string]interface{}{"data": trips, "page": page, "limit": limit})
}

func VehicleTripStats(w http.ResponseWriter, r *http.Request) {
	var todayTrips, monthTrips int
	var todayVolume, monthVolume float64
	db := database.Pool
	ctx := r.Context()
	db.QueryRow(ctx, "SELECT COUNT(*) FROM vehicle_trips WHERE DATE(check_in_time) = CURRENT_DATE").Scan(&todayTrips)
	db.QueryRow(ctx, "SELECT COUNT(*) FROM vehicle_trips WHERE check_in_time >= DATE_TRUNC('month', CURRENT_DATE)").Scan(&monthTrips)
	db.QueryRow(ctx, "SELECT COALESCE(SUM(actual_quantity),0) FROM vehicle_trips WHERE DATE(check_in_time) = CURRENT_DATE").Scan(&todayVolume)
	db.QueryRow(ctx, "SELECT COALESCE(SUM(actual_quantity),0) FROM vehicle_trips WHERE check_in_time >= DATE_TRUNC('month', CURRENT_DATE)").Scan(&monthVolume)
	JSON(w, map[string]interface{}{
		"today_trips":  todayTrips,
		"month_trips":  monthTrips,
		"today_volume": todayVolume,
		"month_volume": monthVolume,
	})
}

func ListCostNorms(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(r.Context(),
		`SELECT id, norm_name, norm_type, unit_cost, unit, material_type, status
		 FROM cost_norms ORDER BY norm_name`)
	if err != nil {
		JSONError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var norms []map[string]interface{}
	for rows.Next() {
		var id int
		var name, ntype, unit, mat, status string
		var cost float64
		rows.Scan(&id, &name, &ntype, &cost, &unit, &mat, &status)
		norms = append(norms, map[string]interface{}{
			"id": id, "norm_name": name, "norm_type": ntype,
			"unit_cost": cost, "unit": unit, "material_type": mat, "status": status,
		})
	}
	JSON(w, map[string]interface{}{"data": norms})
}

func ListProductionCosts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(r.Context(),
		`SELECT id, cost_type, cost_category, norm_value, norm_unit, actual_value, actual_unit, period, mine_area, description
		 FROM production_costs ORDER BY period DESC`)
	if err != nil {
		JSONError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var costs []map[string]interface{}
	for rows.Next() {
		var id int
		var ctype, cat, unit, aunit, period, mine, desc string
		var normVal, actVal float64
		rows.Scan(&id, &ctype, &cat, &normVal, &unit, &actVal, &aunit, &period, &mine, &desc)
		costs = append(costs, map[string]interface{}{
			"id": id, "cost_type": ctype, "cost_category": cat,
			"norm_value": normVal, "norm_unit": unit,
			"actual_value": actVal, "actual_unit": aunit,
			"period": period, "mine_area": mine, "description": desc,
		})
	}
	JSON(w, map[string]interface{}{"data": costs})
}

func ProductionCostSummary(w http.ResponseWriter, r *http.Request) {
	db := database.Pool
	ctx := r.Context()
	var totalCost, totalTons float64
	db.QueryRow(ctx, "SELECT COALESCE(SUM(actual_value),0) FROM production_costs").Scan(&totalCost)
	db.QueryRow(ctx, "SELECT COALESCE(SUM(kl_hang::numeric),0) FROM tickets WHERE loai='Xuất'").Scan(&totalTons)
	costPerTon := 0.0
	if totalTons > 0 {
		costPerTon = totalCost / totalTons
	}
	JSON(w, map[string]interface{}{
		"total_cost":   totalCost,
		"total_tons":   totalTons,
		"cost_per_ton": costPerTon,
	})
}

func ListDeliveries(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(r.Context(),
		`SELECT id, order_code, customer_name, vehicle_plate, product_type,
		        quantity_ordered, quantity_delivered, quantity_confirmed,
		        confirmation_status, confirmed_by
		 FROM delivery_confirmations ORDER BY created_at DESC LIMIT 50`)
	if err != nil {
		JSONError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var deliveries []map[string]interface{}
	for rows.Next() {
		var id int
		var order, cust, plate, prod, status, confirmedBy string
		var qtyOrd, qtyDel, qtyConf float64
		rows.Scan(&id, &order, &cust, &plate, &prod, &qtyOrd, &qtyDel, &qtyConf, &status, &confirmedBy)
		deliveries = append(deliveries, map[string]interface{}{
			"id": id, "order_code": order, "customer_name": cust,
			"vehicle_plate": plate, "product_type": prod,
			"quantity_ordered": qtyOrd, "quantity_delivered": qtyDel,
			"quantity_confirmed": qtyConf, "confirmation_status": status,
			"confirmed_by": confirmedBy,
		})
	}
	JSON(w, map[string]interface{}{"data": deliveries})
}

func ListTaxRecords(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(r.Context(),
		`SELECT id, tax_type, tax_code, period, tax_amount, paid_amount,
		        due_date, paid_date, status, authority
		 FROM tax_records ORDER BY due_date`)
	if err != nil {
		JSONError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var taxes []map[string]interface{}
	for rows.Next() {
		var id int
		var ttype, tcode, period, due, paidDate, status, auth string
		var taxAmt, paidAmt float64
		rows.Scan(&id, &ttype, &tcode, &period, &taxAmt, &paidAmt, &due, &paidDate, &status, &auth)
		taxes = append(taxes, map[string]interface{}{
			"id": id, "tax_type": ttype, "tax_code": tcode,
			"period": period, "tax_amount": taxAmt, "paid_amount": paidAmt,
			"due_date": due, "paid_date": paidDate, "status": status, "authority": auth,
		})
	}
	JSON(w, map[string]interface{}{"data": taxes})
}

func ListDelegations(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(r.Context(),
		`SELECT id, delegator_name, delegate_name, permission_type, scope,
		        start_date, end_date, status
		 FROM delegations ORDER BY created_at DESC`)
	if err != nil {
		JSONError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var dels []map[string]interface{}
	for rows.Next() {
		var id int
		var delegator, delegate, perm, scope, from, to, status string
		rows.Scan(&id, &delegator, &delegate, &perm, &scope, &from, &to, &status)
		dels = append(dels, map[string]interface{}{
			"id": id, "delegator_name": delegator, "delegate_name": delegate,
			"permission_type": perm, "scope": scope,
			"start_date": from, "end_date": to, "status": status,
		})
	}
	JSON(w, map[string]interface{}{"data": dels})
}

func ListRiskAlerts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(r.Context(),
		`SELECT id, alert_type, severity, title, description, module_source, status, created_at
		 FROM risk_alerts ORDER BY
			CASE severity WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 ELSE 4 END,
			created_at DESC`)
	if err != nil {
		JSONError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var alerts []map[string]interface{}
	for rows.Next() {
		var id int
		var atype, sev, title, desc, module, status string
		var createdAt interface{}
		rows.Scan(&id, &atype, &sev, &title, &desc, &module, &status, &createdAt)
		alerts = append(alerts, map[string]interface{}{
			"id": id, "alert_type": atype, "severity": sev,
			"title": title, "description": desc,
			"module_source": module, "status": status, "created_at": createdAt,
		})
	}
	JSON(w, map[string]interface{}{"data": alerts})
}

func ListGeofences(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(r.Context(),
		`SELECT id, name, fence_type, center_lat, center_lng, radius, status
		 FROM geofences ORDER BY name`)
	if err != nil {
		JSONError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var fences []map[string]interface{}
	for rows.Next() {
		var id int
		var name, ftype, status string
		var lat, lng, radius float64
		rows.Scan(&id, &name, &ftype, &lat, &lng, &radius, &status)
		fences = append(fences, map[string]interface{}{
			"id": id, "name": name, "fence_type": ftype,
			"center_lat": lat, "center_lng": lng, "radius": radius, "status": status,
		})
	}
	JSON(w, map[string]interface{}{"data": fences})
}

func ListAuthorizations(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(r.Context(),
		`SELECT id, auth_type, document_type, document_code, authorizer_name,
		        signer_name, signature_method, signed_at, valid_from, valid_to, status
		 FROM authorizations ORDER BY created_at DESC`)
	if err != nil {
		JSONError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var auths []map[string]interface{}
	for rows.Next() {
		var id int
		var atype, dtype, dcode, authName, signer, method, from, to, status string
		var signedAt interface{}
		rows.Scan(&id, &atype, &dtype, &dcode, &authName, &signer, &method, &signedAt, &from, &to, &status)
		auths = append(auths, map[string]interface{}{
			"id": id, "auth_type": atype, "document_type": dtype,
			"document_code": dcode, "authorizer_name": authName,
			"signer_name": signer, "signature_method": method,
			"signed_at": signedAt, "valid_from": from, "valid_to": to, "status": status,
		})
	}
	JSON(w, map[string]interface{}{"data": auths})
}
