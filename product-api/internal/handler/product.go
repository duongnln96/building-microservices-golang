package handler

import (
	"context"
	"fmt"
	"net/http"

	protos "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"github.com/duongnln96/building-microservices-golang/product-api/config"
	"github.com/duongnln96/building-microservices-golang/product-api/internal/data"
	"github.com/duongnln96/building-microservices-golang/product-api/internal/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

type ProductHandlerI interface {
	StartServerLock()
}

type ProductHandlerDeps struct {
	Ctx context.Context
	Log *zap.SugaredLogger
	Cfg *config.AppConfig
	Db  data.ProductsDBI
	Cc  protos.CurrencyClient
}

type productHandler struct {
	ctx context.Context
	log *zap.SugaredLogger
	cfg *config.AppConfig
	db  data.ProductsDBI
	cc  protos.CurrencyClient
}

func NewProductHandler(deps ProductHandlerDeps) ProductHandlerI {
	return &productHandler{
		ctx: deps.Ctx,
		log: deps.Log,
		cfg: deps.Cfg,
		db:  deps.Db,
		cc:  deps.Cc,
	}
}

func (p *productHandler) StartServerLock() {
	e := echo.New()

	e.Use(utils.ZapLogger(p.log))
	e.Validator = data.NewValidation()

	e.POST("/products", p.createProduct)
	e.GET("/products", p.getAllProducts)
	e.GET("/product/:id", p.getSingleProduct)
	e.PUT("/product/:id", p.updateProduct)
	e.DELETE("/product/:id", p.deleteProduct)

	err := e.Start(fmt.Sprintf(":%d", p.cfg.Server.Port))
	if err != nil {
		p.log.Fatalf("Error while starting Echo server: %s", err)
	}
}

// POST /product
func (p *productHandler) createProduct(c echo.Context) error {
	prod := data.Product{}

	if err := c.Bind(&prod); err != nil {
		return p.responseErrCode(c, http.StatusBadRequest, fmt.Sprintf("Bind: %s", err.Error()))
	}

	if err := c.Validate(prod); err != nil {
		return p.responseErrCode(c, http.StatusBadRequest, fmt.Sprintf("Validation: %s", err.Error()))
	}

	p.db.AddProduct(&prod)
	return p.responseStatusOK(c)
}

// GET /products
func (p *productHandler) getAllProducts(c echo.Context) error {
	currency := p.getProductQuery(c)

	products, err := p.db.GetProducts()
	if err != nil {
		return p.responseErrCode(c, http.StatusInternalServerError, err.Error())
	}

	if currency == "" {
		return p.responseData(c, products)
	}

	rate, err := p.getRate(currency)
	if err != nil {
		return p.responseErrCode(c, http.StatusInternalServerError, err.Error())
	}

	p.log.Infof("Currency Rate %+v", rate)

	prodsWithRate := data.Products{}
	for _, prod := range products {
		newProd := *prod
		newProd.Price = newProd.Price * rate
		prodsWithRate = append(prodsWithRate, &newProd)
	}

	return p.responseData(c, prodsWithRate)
}

// GET /product/:id
func (p *productHandler) getSingleProduct(c echo.Context) error {
	id := p.getProductIDParam(c)
	currency := p.getProductQuery(c)

	prod, err := p.db.GetProductByID(id)
	if err != nil {
		return p.responseErrCode(c, http.StatusNotFound, data.ErrProductNotFound.Error())
	}

	if currency == "" {
		return p.responseData(c, prod)
	}

	rate, err := p.getRate(currency)
	if err != nil {
		return p.responseErrCode(c, http.StatusInternalServerError, err.Error())
	}

	p.log.Infof("Currency Rate %+v", rate)

	prodWithRate := *prod
	prodWithRate.Price = prodWithRate.Price * rate

	return p.responseData(c, prodWithRate)
}

// PUT /product/:id
func (p *productHandler) updateProduct(c echo.Context) error {
	id := p.getProductIDParam(c)

	prod := data.Product{}
	if err := c.Bind(&prod); err != nil {
		return p.responseErrCode(c, http.StatusBadRequest, fmt.Sprintf("Bind: %s", err.Error()))
	}

	if err := c.Validate(prod); err != nil {
		return p.responseErrCode(c, http.StatusBadRequest, fmt.Sprintf("Validation: %s", err.Error()))
	}

	if err := p.db.UpdateProduct(id, &prod); err == data.ErrProductNotFound {
		return p.responseErrCode(c, http.StatusNotFound, data.ErrProductNotFound.Error())
	}

	return p.responseStatusOK(c)
}

// DELETE product/:id
func (p *productHandler) deleteProduct(c echo.Context) error {
	id := p.getProductIDParam(c)

	err := p.db.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		return p.responseErrCode(c, http.StatusNotFound, data.ErrProductNotFound.Error())
	} else if err != nil {
		return p.responseErrCode(c, http.StatusInternalServerError, "")
	}

	return p.responseStatusOK(c)
}
