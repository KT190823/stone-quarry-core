package models

type InventoryInbound struct {
	Code     string  `json:"code"`
	Source   string  `json:"source"`
	Loc      string  `json:"loc"`
	Item     string  `json:"item"`
	Qty      float64 `json:"qty"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Date     string  `json:"date,omitempty"`
	Status   string  `json:"status"`
}

type InventoryOutbound struct {
	Code     string  `json:"code"`
	Customer string  `json:"customer"`
	Dest     string  `json:"dest"`
	Item     string  `json:"item"`
	Qty      float64 `json:"qty"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Date     string  `json:"date,omitempty"`
	Status   string  `json:"status"`
}

type InventoryStocktake struct {
	Code     string  `json:"code"`
	Zone     string  `json:"zone"`
	Item     string  `json:"item"`
	Volume   float64 `json:"volume"`
	Survey   float64 `json:"survey"`
	ERP      float64 `json:"erp"`
	Book     float64 `json:"book"`
	Actual   float64 `json:"actual"`
	Diff     float64 `json:"diff"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Date     string  `json:"date,omitempty"`
	Status   string  `json:"status"`
}

type InventoryMovement struct {
	Code     string  `json:"code"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	Item     string  `json:"item"`
	Qty      float64 `json:"qty"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Date     string  `json:"date,omitempty"`
	Status   string  `json:"status"`
}
