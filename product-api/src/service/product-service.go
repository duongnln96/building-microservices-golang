package service

import (
	"github.com/duongnln96/building-microservices-golang/product-api/src/domain/dto"
	"github.com/duongnln96/building-microservices-golang/product-api/src/domain/entity"
	"github.com/duongnln96/building-microservices-golang/product-api/src/domain/repository"
	"go.uber.org/zap"
)

type ProductServiceI interface {
	All() ([]entity.Product, error)
	FindByID(int) (entity.Product, error)
	Create(*dto.ProductDTO) error
	Update(*dto.ProductUpdateDTO) error
	Delete(*dto.ProductUpdateDTO) error
}

type ProductServiceDeps struct {
	Log  *zap.SugaredLogger
	Repo repository.ProductsRepoI
}

type productSerivce struct {
	log  *zap.SugaredLogger
	repo repository.ProductsRepoI
}

func NewProductSerivce(deps ProductServiceDeps) ProductServiceI {
	return &productSerivce{
		log:  deps.Log,
		repo: deps.Repo,
	}
}

func (svc *productSerivce) All() ([]entity.Product, error) {
	prods, err := svc.repo.AllProducts()
	if err != nil {
		svc.log.Debugf("SVC cannot fetch all products %+v", err)
		return nil, err
	}
	return prods, nil
}

func (svc *productSerivce) FindByID(id int) (entity.Product, error) {
	prod, err := svc.repo.FindProductByID(id)
	if err != nil {
		svc.log.Debugf("SVC cannot fetch product %+v", err)
		return entity.Product{}, err
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
