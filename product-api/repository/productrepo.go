package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/duongnln96/building-microservices-golang/product-api/entity"
	tools "github.com/duongnln96/building-microservices-golang/product-api/tools/postgresql"
	"go.uber.org/zap"
)

type ProductsDBI interface {
	AllProducts() (entity.Products, error)
	FindProductByID(int) (*entity.Product, error)
	CreateProduct(*entity.Product) error
	UpdateProduct(*entity.Product) error
	DeleteProduct(*entity.Product) error
}

type ProductsDBDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
	DB  tools.PsqlConnectorI
}

type productsDB struct {
	log *zap.SugaredLogger
	ctx context.Context
	db  tools.PsqlConnectorI
}

func NewProductDB(deps ProductsDBDeps) ProductsDBI {
	return &productsDB{
		log: deps.Log,
		ctx: deps.Ctx,
		db:  deps.DB,
	}
}

var ErrProductNotFound error = fmt.Errorf("Product not found")

func (db *productsDB) AllProducts() (entity.Products, error) {
	return entity.ProductList, nil
}

func (db *productsDB) FindProductByID(id int) (*entity.Product, error) {
	idx := db.findProductIndex(id)
	if idx == -1 {
		return nil, ErrProductNotFound
	}

	return entity.ProductList[idx], nil
}

func (db *productsDB) CreateProduct(p *entity.Product) error {
	p.ID = db.getNextId()
	p.CreatedOn = time.Now().UTC().String()
	p.UpdatedOn = time.Now().UTC().String()
	entity.ProductList = append(entity.ProductList, p)

	return nil
}

func (db *productsDB) UpdateProduct(p *entity.Product) error {
	idx := db.findProductIndex(p.ID)
	if idx == -1 {
		return ErrProductNotFound
	}
	oldProd := entity.ProductList[idx]
	p.CreatedOn = oldProd.CreatedOn
	p.UpdatedOn = time.Now().UTC().String()
	entity.ProductList[idx] = p

	return nil
}

func (db *productsDB) DeleteProduct(p *entity.Product) error {
	idx := db.findProductIndex(p.ID)
	if idx == -1 {
		return ErrProductNotFound
	}

	entity.ProductList = append(entity.ProductList[:idx], entity.ProductList[idx+1:]...)
	return nil
}

func (db *productsDB) findProductIndex(id int) int {
	for i, p := range entity.ProductList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func (db *productsDB) getNextId() int {
	lp := entity.ProductList[len(entity.ProductList)-1]
	return lp.ID + 1
}
