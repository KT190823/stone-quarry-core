package models

// PrintTemplate is a configurable document/report layout used to render
// printable vouchers (tickets, invoices, warehouse slips, alerts, reports).
// The Layout field holds a designer JSON document produced by the frontend.
type PrintTemplate struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	DocType     string `json:"docType"`
	Size        string `json:"size"`
	Orientation string `json:"orientation"`
	Description string `json:"description"`
	Layout      string `json:"layout"`
	IsDefault   bool   `json:"isDefault"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
