package handlers

import (
	"net/http"

	"mo-da-backend/internal/services"
)

type CatalogHandler struct {
	materialSvc    *services.CatalogMaterialService
	supplierSvc    *services.SupplierService
	customerSvc    *services.CustomerService
	ticketTypeSvc  *services.TicketTypeService
	vehicleCatSvc  *services.VehicleCatalogService
	stationSvc     *services.StationService
	equipmentSvc   *services.EquipmentService
}

func NewCatalogHandler() *CatalogHandler {
	return &CatalogHandler{
		materialSvc:   services.NewCatalogMaterialService(),
		supplierSvc:   services.NewSupplierService(),
		customerSvc:   services.NewCustomerService(),
		ticketTypeSvc: services.NewTicketTypeService(),
		vehicleCatSvc: services.NewVehicleCatalogService(),
		stationSvc:    services.NewStationService(),
		equipmentSvc:  services.NewEquipmentService(),
	}
}

func (h *CatalogHandler) ListMaterials(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.materialSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CatalogHandler) GetMaterial(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.materialSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *CatalogHandler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.materialSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *CatalogHandler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.materialSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *CatalogHandler) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.materialSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *CatalogHandler) ListTicketTypes(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.ticketTypeSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CatalogHandler) ListVehicleCatalogs(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.vehicleCatSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CatalogHandler) ListStations(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.stationSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CatalogHandler) ListEquipment(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.equipmentSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CatalogHandler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.supplierSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CatalogHandler) GetSupplier(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.supplierSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *CatalogHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.supplierSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *CatalogHandler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.supplierSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *CatalogHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.supplierSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *CatalogHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.customerSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CatalogHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.customerSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *CatalogHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.customerSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *CatalogHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.customerSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *CatalogHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.customerSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}
