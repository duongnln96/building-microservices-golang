package data

import "time"

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

// Example data
var productList = []*Product{
	{
		ID:          1,
		Name:        "Sony Turntable - PSLX350H",
		Description: "Belt Drive System 33-1/3 and 45 RPM Speeds Servo Speed Control Supplied Moving Magnet Phono Cartridge Bonded Diamond Stylus Static Balance Tonearm Pitch Control",
		Price:       128.00,
		SKU:         "PSLX350H",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Bose Acoustimass 5 Series III Speaker System - AM53BK",
		Description: "2 Dual Cube Speakers With Two 2-1/2 Wide-range Drivers In Each Speaker Powerful Bass Module With Two 5-1/2 Woofers 200 Watts Max Power Black Finish",
		Price:       256.00,
		SKU:         "AM53BK",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
