package repository

import (
	"context"
	"database/sql"

	"github.com/duongnln96/building-microservices-golang/product-api/entity"
	tools "github.com/duongnln96/building-microservices-golang/product-api/tools/postgresql"
	"go.uber.org/zap"
)

type ProductsRepoI interface {
	AllProducts() ([]entity.Product, error)
	FindProductByID(int) (entity.Product, error)
	CreateProduct(*entity.Product) error
	UpdateProduct(*entity.Product) error
	DeleteProduct(*entity.Product) error
}

type ProductsRepoDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
	DB  tools.PsqlConnectorI
}

type productsRepo struct {
	log *zap.SugaredLogger
	ctx context.Context
	db  tools.PsqlConnectorI
}

func NewProductDB(deps ProductsRepoDeps) ProductsRepoI {
	return &productsRepo{
		log: deps.Log,
		ctx: deps.Ctx,
		db:  deps.DB,
	}
}

func (pr *productsRepo) AllProducts() ([]entity.Product, error) {
	products := make([]entity.Product, 0)
	err := pr.db.DBQueryRows(
		func(r *sql.Rows) bool {
			prod := entity.Product{}
			err := r.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.SKU)
			if err != nil {
				return false
			}
			products = append(products, prod)
			return true
		},
		"select id, name, description, price, sku from products;",
	)
	return products, err
}

func (pr *productsRepo) FindProductByID(id int) (entity.Product, error) {
	prod := entity.Product{}
	err := pr.db.DBQueryRow(
		func(r *sql.Row) error {
			if err := r.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.SKU); err != nil {
				return err
			}
			return nil
		},
		"select id, name, description, price, sku from products where id=$1;",
		id,
	)

	return prod, err
}

func (pr *productsRepo) CreateProduct(p *entity.Product) error {
	err := pr.db.DBExec(
		"insert into products (name, description, price, sku) values ($1, $2, $3, $4);",
		p.Name, p.Description, p.Price, p.SKU,
	)
	return err
}

func (pr *productsRepo) UpdateProduct(p *entity.Product) error {
	err := pr.db.DBExec(
		"update products set name=$1, description=$2, price=$3, sku=$4 where id=$5;",
		p.Name, p.Description, p.Price, p.SKU, p.ID,
	)

	return err
}

func (pr *productsRepo) DeleteProduct(p *entity.Product) error {
	err := pr.db.DBExec(
		"delete from products where id=$1",
		p.ID,
	)
	return err
}
