package service

import (
	"context"

	"github.com/duongnln96/building-microservices-golang/product-api/dto"
	"github.com/duongnln96/building-microservices-golang/product-api/entity"
	"github.com/duongnln96/building-microservices-golang/product-api/repository"

	"go.uber.org/zap"
)

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

type ProductServiceI interface {
	All() (entity.Products, error)
	FindByID(int) (*entity.Product, error)
	Create(*dto.ProductDTO) error
	Update(*dto.ProductUpdateDTO) error
	Delete(*dto.ProductUpdateDTO) error
}

type ProductServiceDeps struct {
	Ctx  context.Context
	Log  *zap.SugaredLogger
	Repo repository.ProductsDBI
}

type productSerivce struct {
	ctx  context.Context
	log  *zap.SugaredLogger
	repo repository.ProductsDBI
}

func NewProductSerivce(deps ProductServiceDeps) ProductServiceI {
	return &productSerivce{
		ctx:  deps.Ctx,
		log:  deps.Log,
		repo: deps.Repo,
	}
}

func (svc *productSerivce) All() (entity.Products, error) {
	prods, err := svc.repo.AllProducts()
	if err != nil {
		svc.log.Debugf("SVC cannot fetch products %+v", err)
		return nil, err
	}
	return prods, nil
}

func (svc *productSerivce) FindByID(id int) (*entity.Product, error) {
	prod, err := svc.repo.FindProductByID(id)
	if err != nil {
		svc.log.Debugf("SVC cannot fetch product %+v", err)
		return nil, err
	}

	return prod, nil
}

func (svc *productSerivce) Create(p *dto.ProductDTO) error {
	prod := entity.Product{}
	prod.Name = p.Name
	prod.Description = p.Description
	prod.Price = p.Price
	prod.SKU = p.SKU

	err := svc.repo.CreateProduct(&prod)
	if err != nil {
		svc.log.Debugf("SVC Cannot insert product %+v", err)
		return err
	}

	return nil
}

func (svc *productSerivce) Update(p *dto.ProductUpdateDTO) error {
	prod := entity.Product{}
	prod.ID = p.ID
	prod.Name = p.Name
	prod.Description = p.Description
	prod.Price = p.Price
	prod.SKU = p.SKU

	err := svc.repo.UpdateProduct(&prod)
	if err != nil {
		svc.log.Debugf("SVC Cannot update product %+v", err)
		return err
	}

	return nil
}

func (svc *productSerivce) Delete(p *dto.ProductUpdateDTO) error {
	prod := entity.Product{}
	prod.ID = p.ID

	err := svc.repo.DeleteProduct(&prod)
	if err != nil {
		svc.log.Debugf("SVC Cannot delete product %+v", err)
		return err
	}

	return nil
}
