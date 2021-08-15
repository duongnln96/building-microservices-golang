package data

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ProductsDBI interface {
	GetProducts() (Products, error)
	GetProductByID(int) (*Product, error)
	AddProduct(*Product)
	UpdateProduct(int, *Product) error
	DeleteProduct(int) error
}

type ProductsDBDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
}

type productsDB struct {
	log *zap.SugaredLogger
	ctx context.Context
}

func NewProductDB(deps ProductsDBDeps) ProductsDBI {
	return &productsDB{
		log: deps.Log,
		ctx: deps.Ctx,
	}
}

var ErrProductNotFound error = fmt.Errorf("Product not found")

func (db *productsDB) GetProducts() (Products, error) {
	return productList, nil
}

func (db *productsDB) GetProductByID(id int) (*Product, error) {
	idx := db.findProductIndex(id)
	if idx == -1 {
		return nil, ErrProductNotFound
	}

	return productList[idx], nil
}

func (db *productsDB) AddProduct(p *Product) {
	p.ID = db.getNextId()
	p.CreatedOn = time.Now().UTC().String()
	p.UpdatedOn = time.Now().UTC().String()
	productList = append(productList, p)
}

func (db *productsDB) UpdateProduct(id int, p *Product) error {
	idx := db.findProductIndex(id)
	if idx == -1 {
		return ErrProductNotFound
	}

	p.ID = id
	productList[idx] = p

	return nil
}

func (db *productsDB) DeleteProduct(id int) error {
	idx := db.findProductIndex(id)
	if idx == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:idx], productList[idx+1:]...)
	return nil
}

func (db *productsDB) findProductIndex(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func (db *productsDB) getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}
