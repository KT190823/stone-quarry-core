package models

type TicketStage string

const (
	StageDraft     TicketStage = "draft"
	StageWeigh1    TicketStage = "weigh1"
	StageWeigh2    TicketStage = "weigh2"
	StageConfirmed TicketStage = "confirmed"
	StageInvoiced  TicketStage = "invoiced"
)

type DashboardStat struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Sub   string `json:"sub"`
	Trend string `json:"trend"`
	Icon  string `json:"icon"`
}

type DashboardLog struct {
	Time string `json:"time"`
	User string `json:"user"`
	Text string `json:"text"`
}

type Report struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Item      string `json:"item"`
	Type      string `json:"type"`
	Period    string `json:"period"`
	Plan      string `json:"plan"`
	Actual    string `json:"actual"`
	Diff      string `json:"diff"`
	UpdatedAt string `json:"updatedAt"`
	Status    string `json:"status"`
}

type SettingConfig struct {
	Code      string `json:"code"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	Val       string `json:"val"`
	UpdatedAt string `json:"updatedAt"`
	Scope     string `json:"scope"`
	Status    string `json:"status"`
}

type UserRow struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	Dept      string `json:"dept"`
	Email     string `json:"email"`
	LastLogin string `json:"lastLogin"`
	Status    string `json:"status"`
}

type RoleRow struct {
	Role       string `json:"role"`
	Desc       string `json:"desc"`
	Users      string `json:"users"`
	Level      string `json:"level"`
	Permission string `json:"permission"`
	Status     string `json:"status"`
}

type LogRow struct {
	Time   string `json:"time"`
	User   string `json:"user"`
	Name   string `json:"name"`
	Action string `json:"action"`
	Target string `json:"target"`
	IP     string `json:"ip"`
	Status string `json:"status"`
}

type InvoiceItem struct {
	Name      string  `json:"name"`
	Unit      string  `json:"unit"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
	VatRate   int     `json:"vatRate"`
}

type InvoiceData struct {
	InvoiceNumber string        `json:"invoiceNumber"`
	InvoiceSeries string        `json:"invoiceSeries"`
	InvoiceDate   string        `json:"invoiceDate"`
	Seller        InvoiceParty  `json:"seller"`
	Buyer         InvoiceParty  `json:"buyer"`
	Vehicle       string        `json:"vehicle"`
	WeighBridge   string        `json:"weighBridge"`
	Warehouse     string        `json:"warehouse"`
	DeliveryNote  string        `json:"deliveryNote"`
	PaymentMethod string        `json:"paymentMethod"`
	Items         []InvoiceItem `json:"items"`
	Subtotal      float64       `json:"subtotal"`
	VATAmount     float64       `json:"vatAmount"`
	GrandTotal    float64       `json:"grandTotal"`
	AmountInWords string        `json:"amountInWords"`
	Notes         []string      `json:"notes"`
}

type InvoiceParty struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	TaxCode string `json:"taxCode"`
	Bank    string `json:"bank,omitempty"`
	Contact string `json:"contact,omitempty"`
}

type TicketType struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Process  string `json:"process"`
	Template string `json:"template"`
	Status   string `json:"status"`
}

type VehicleCatalog struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Axles        string `json:"axles"`
	StandardTare string `json:"standardTare"`
	MaxPayload   string `json:"maxPayload"`
	Status       string `json:"status"`
}

type Station struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Size      string `json:"size"`
	Capacity  string `json:"capacity"`
	Indicator string `json:"indicator"`
	Status    string `json:"status"`
}

type Equipment struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location string `json:"location"`
	IP       string `json:"ip"`
	Res      string `json:"res"`
	Status   string `json:"status"`
}

type DebtRow struct {
	Code     string `json:"code"`
	Partner  string `json:"partner"`
	Customer string `json:"customer"`
	Limit    string `json:"limit"`
	Balance  string `json:"balance"`
	Amount   string `json:"amount"`
	Due      string `json:"due"`
	Status   string `json:"status"`
}

type ReconcileRow struct {
	Period      string `json:"period"`
	Day         string `json:"day"`
	ScaleRevenue string `json:"scaleRevenue"`
	ERPRevenue  string `json:"erpRevenue"`
	ERP         string `json:"erp"`
	Cash        string `json:"cash"`
	Diff        string `json:"diff"`
	Status      string `json:"status"`
}
