package data

import (
	"encoding/json"
	"io"
	"time"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	SKU         string `json:"sku"`
	CreatedOn   string `json:"-"`
	UpdatedOn   string `json:"-"`
	DeletedOn   string `json:"-"`
}

type Products []*Product

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextId()
	p.CreatedOn = time.Now().UTC().String()
	p.UpdatedOn = time.Now().UTC().String()
	productList = append(productList, p)
}

func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

// Example data
var productList = []*Product{
	{
		ID:          1,
		Name:        "Sony Turntable - PSLX350H",
		Description: "Sony Turntable - PSLX350H/ Belt Drive System/ 33-1/3 and 45 RPM Speeds/ Servo Speed Control/ Supplied Moving Magnet Phono Cartridge/ Bonded Diamond Stylus/ Static Balance Tonearm/ Pitch Control",
		Price:       "$399.00",
		SKU:         "PSLX350H",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Bose Acoustimass 5 Series III Speaker System - AM53BK",
		Description: "Bose Acoustimass 5 Series III Speaker System - AM53BK/ 2 Dual Cube Speakers With Two 2-1/2' Wide-range Drivers In Each Speaker/ Powerful Bass Module With Two 5-1/2' Woofers/ 200 Watts Max Power/ Black Finish",
		Price:       "$399.00",
		SKU:         "AM53BK",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
