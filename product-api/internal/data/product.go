package data

import (
	"fmt"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

var ErrProductNotFound error = fmt.Errorf("Product not found")

func GetProducts() Products {
	return productList
}

func GetProductByID(id int) (*Product, error) {
	idx := findProductIndex(id)
	if idx == -1 {
		return nil, ErrProductNotFound
	}

	return productList[idx], nil
}

func AddProduct(p *Product) {
	p.ID = getNextId()
	p.CreatedOn = time.Now().UTC().String()
	p.UpdatedOn = time.Now().UTC().String()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	idx := findProductIndex(id)
	if idx == -1 {
		return ErrProductNotFound
	}

	p.ID = id
	productList[idx] = p

	return nil
}

func DeleteProduct(id int) error {
	idx := findProductIndex(id)
	if idx == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:idx], productList[idx+1:]...)
	return nil
}

func findProductIndex(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
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
