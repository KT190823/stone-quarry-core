package models

import "time"

// InventoryProduct represents a declared rock / mineral / warehouse item
type InventoryProduct struct {
	ID            int       `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	Unit          string    `json:"unit"`
	Density       float64   `json:"density"`
	SalePrice     float64   `json:"salePrice"`
	PurchasePrice float64   `json:"purchasePrice"`
	StorageLoc    string    `json:"storageLoc"`
	Standard      string    `json:"standard"`
	MinStock      float64   `json:"minStock"`
	CurrentStock  float64   `json:"currentStock"`
	Status        string    `json:"status"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// PurchaseVoucher represents a goods purchase order / voucher
type PurchaseVoucher struct {
	ID            int            `json:"id"`
	Code          string         `json:"code"`
	SupplierCode  string         `json:"supplierCode"`
	SupplierName  string         `json:"supplierName"`
	Date          string         `json:"date"`
	WarehouseLoc  string         `json:"warehouseLoc"`
	TotalAmount   float64        `json:"totalAmount"`
	VatAmount     float64        `json:"vatAmount"`
	GrandTotal    float64        `json:"grandTotal"`
	PaymentStatus string         `json:"paymentStatus"`
	Status        string         `json:"status"`
	CreatedBy     string         `json:"createdBy"`
	Notes         string         `json:"notes"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	Items         []PurchaseItem `json:"items,omitempty"`
}

type PurchaseItem struct {
	ID          int     `json:"id"`
	VoucherID   int     `json:"voucherId"`
	VoucherCode string  `json:"voucherCode"`
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	Unit        string  `json:"unit"`
	Density     float64 `json:"density"`
	UnitPrice   float64 `json:"unitPrice"`
	Quantity    float64 `json:"quantity"`
	TotalAmount float64 `json:"totalAmount"`
	WeightTon   float64 `json:"weightTon"`
	Standard    string  `json:"standard"`
	StorageLoc  string  `json:"storageLoc"`
	Notes       string  `json:"notes"`
}

// SalesVoucher represents a sales order / delivery voucher
type SalesVoucher struct {
	ID            int         `json:"id"`
	Code          string      `json:"code"`
	CustomerCode  string      `json:"customerCode"`
	CustomerName  string      `json:"customerName"`
	Date          string      `json:"date"`
	WarehouseLoc  string      `json:"warehouseLoc"`
	LicensePlate  string      `json:"licensePlate"`
	TicketCode    string      `json:"ticketCode"`
	TotalAmount   float64     `json:"totalAmount"`
	VatAmount     float64     `json:"vatAmount"`
	GrandTotal    float64     `json:"grandTotal"`
	PaidAmount    float64     `json:"paidAmount"`
	DebtAmount    float64     `json:"debtAmount"`
	PaymentStatus string      `json:"paymentStatus"`
	Status        string      `json:"status"`
	CreatedBy     string      `json:"createdBy"`
	Notes         string      `json:"notes"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	Items         []SalesItem `json:"items,omitempty"`
}

type SalesItem struct {
	ID          int     `json:"id"`
	VoucherID   int     `json:"voucherId"`
	VoucherCode string  `json:"voucherCode"`
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	Unit        string  `json:"unit"`
	Density     float64 `json:"density"`
	UnitPrice   float64 `json:"unitPrice"`
	Quantity    float64 `json:"quantity"`
	TotalAmount float64 `json:"totalAmount"`
	WeightTon   float64 `json:"weightTon"`
	Standard    string  `json:"standard"`
	StorageLoc  string  `json:"storageLoc"`
	Notes       string  `json:"notes"`
}

// ReturnVoucher represents goods return (sales return or purchase return)
type ReturnVoucher struct {
	ID             int          `json:"id"`
	Code           string       `json:"code"`
	ReturnType     string       `json:"returnType"` // 'sales_return' or 'purchase_return'
	PartnerCode    string       `json:"partnerCode"`
	PartnerName    string       `json:"partnerName"`
	RefVoucherCode string       `json:"refVoucherCode"`
	Date           string       `json:"date"`
	WarehouseLoc   string       `json:"warehouseLoc"`
	Reason         string       `json:"reason"`
	TotalAmount    float64      `json:"totalAmount"`
	Status         string       `json:"status"`
	CreatedBy      string       `json:"createdBy"`
	Notes          string       `json:"notes"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	Items          []ReturnItem `json:"items,omitempty"`
}

type ReturnItem struct {
	ID          int     `json:"id"`
	VoucherID   int     `json:"voucherId"`
	VoucherCode string  `json:"voucherCode"`
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	Unit        string  `json:"unit"`
	UnitPrice   float64 `json:"unitPrice"`
	Quantity    float64 `json:"quantity"`
	TotalAmount float64 `json:"totalAmount"`
	Notes       string  `json:"notes"`
}

// ReceiptVoucher represents a cash / bank receipt voucher (Phiếu thu)
type ReceiptVoucher struct {
	ID            int       `json:"id"`
	Code          string    `json:"code"`
	PartnerCode   string    `json:"partnerCode"`
	PartnerName   string    `json:"partnerName"`
	PartnerType   string    `json:"partnerType"` // 'customer', 'other'
	PayerName     string    `json:"payerName"`
	Reason        string    `json:"reason"`
	RefCode       string    `json:"refCode"`
	Amount        float64   `json:"amount"`
	AmountInWords string    `json:"amountInWords"`
	PaymentMethod string    `json:"paymentMethod"` // 'cash', 'bank_transfer'
	FundAccount   string    `json:"fundAccount"`
	Status        string    `json:"status"`
	Date          string    `json:"date"`
	CreatedBy     string    `json:"createdBy"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// PaymentVoucher represents a cash / bank payment voucher (Phiếu chi)
type PaymentVoucher struct {
	ID            int       `json:"id"`
	Code          string    `json:"code"`
	PartnerCode   string    `json:"partnerCode"`
	PartnerName   string    `json:"partnerName"`
	PartnerType   string    `json:"partnerType"` // 'supplier', 'customer_refund', 'other'
	ReceiverName  string    `json:"receiverName"`
	Reason        string    `json:"reason"`
	RefCode       string    `json:"refCode"`
	Amount        float64   `json:"amount"`
	AmountInWords string    `json:"amountInWords"`
	PaymentMethod string    `json:"paymentMethod"` // 'cash', 'bank_transfer'
	FundAccount   string    `json:"fundAccount"`
	Status        string    `json:"status"`
	Date          string    `json:"date"`
	CreatedBy     string    `json:"createdBy"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
