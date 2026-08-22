package models

type InventoryInbound struct {
	Code     string `json:"code"`
	Source   string `json:"source"`
	Loc      string `json:"loc"`
	Item     string `json:"item"`
	Qty      string `json:"qty"`
	Quantity string `json:"quantity"`
	Status   string `json:"status"`
}

type InventoryOutbound struct {
	Code     string `json:"code"`
	Customer string `json:"customer"`
	Dest     string `json:"dest"`
	Item     string `json:"item"`
	Qty      string `json:"qty"`
	Quantity string `json:"quantity"`
	Status   string `json:"status"`
}

type InventoryStocktake struct {
	Code     string `json:"code"`
	Zone     string `json:"zone"`
	Item     string `json:"item"`
	Volume   string `json:"volume"`
	Survey   string `json:"survey"`
	ERP      string `json:"erp"`
	Book     string `json:"book"`
	Actual   string `json:"actual"`
	Diff     string `json:"diff"`
	Quantity string `json:"quantity"`
	Status   string `json:"status"`
}

type InventoryMovement struct {
	Code     string `json:"code"`
	From     string `json:"from"`
	To       string `json:"to"`
	Item     string `json:"item"`
	Qty      string `json:"qty"`
	Quantity string `json:"quantity"`
	Status   string `json:"status"`
}
