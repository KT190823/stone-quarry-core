package models

type CatalogMaterial struct {
	Code     string `json:"code"`
	Ma       string `json:"ma"`
	Name     string `json:"name"`
	Ten      string `json:"ten"`
	DVT      string `json:"dvt"`
	Density  string `json:"density"`
	TyTrong  string `json:"tyTrong"`
	DinhMuc  string `json:"dinhMuc"`
	Price    string `json:"price"`
	Gia      string `json:"gia"`
	Kho      string `json:"kho"`
	Standard string `json:"standard"`
	Status   string `json:"status"`
	Date     string `json:"date"`
}

type Supplier struct {
	Name    string `json:"name"`
	Ten     string `json:"ten"`
	TaxCode string `json:"taxCode"`
	MST     string `json:"mst"`
	Phone   string `json:"phone"`
	SDT     string `json:"sdt"`
	Items   string `json:"items"`
	Hang    string `json:"hang"`
	No      string `json:"no"`
	Status  string `json:"status"`
	Rating  string `json:"rating"`
	Date    string `json:"date"`
}

type Customer struct {
	Name    string `json:"name"`
	Ten     string `json:"ten"`
	TaxCode string `json:"taxCode"`
	MST     string `json:"mst"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	SDT     string `json:"sdt"`
	Group   string `json:"group"`
	DoanhSo string `json:"doanhSo"`
	CongNo  string `json:"congNo"`
	HanMuc  string `json:"hanMuc"`
	Status  string `json:"status"`
	Date    string `json:"date"`
}
