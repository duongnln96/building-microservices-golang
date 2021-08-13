package handler

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"github.com/duongnln96/building-microservices-golang/product-api/config"
	"github.com/duongnln96/building-microservices-golang/product-api/internal/data"
	"github.com/duongnln96/building-microservices-golang/product-api/internal/utils"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

type ProductHandlerI interface {
	StartServerLock()
}

type ProductHandlerDeps struct {
	Ctx            context.Context
	Log            *zap.SugaredLogger
	Cfg            *config.AppConfig
	CurrencyClient pb.CurrencyClient
}

type productHandler struct {
	ctx            context.Context
	log            *zap.SugaredLogger
	cfg            *config.AppConfig
	currencyClient pb.CurrencyClient
}

func NewProductHandler(deps ProductHandlerDeps) ProductHandlerI {
	return &productHandler{
		ctx:            deps.Ctx,
		log:            deps.Log,
		cfg:            deps.Cfg,
		currencyClient: deps.CurrencyClient,
	}
}

func (p *productHandler) StartServerLock() {
	e := echo.New()

	e.Use(utils.ZapLogger(p.log))

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

	err := p.getProductData(c, &prod)
	if err != nil {
		return p.responseErrCode(c, http.StatusInternalServerError, "Failed to read request body")
	}

	data.AddProduct(&prod)
	return p.responseStatusOK(c)
}

// GET /products
func (p *productHandler) getAllProducts(c echo.Context) error {
	products := data.GetProducts()
	return p.responseData(c, products)
}

// GET /product/:id
func (p *productHandler) getSingleProduct(c echo.Context) error {
	id := p.getProductIDParam(c)

	prod, err := data.GetProductByID(id)
	if err != nil {
		return p.responseErrCode(c, http.StatusNotFound, data.ErrProductNotFound.Error())
	}

	// Get the exchange
	rateRequest := pb.RateRequest{
		Base:        pb.Currencies(pb.Currencies_value["USA"]),
		Destination: pb.Currencies(pb.Currencies_value["VND"]),
	}
	cRate, err := p.currencyClient.GetRate(p.ctx, &rateRequest)
	if err != nil {
		p.log.Errorf("Cannot get rate %+v", err)
		return p.responseErrCode(c, http.StatusInternalServerError, err.Error())
	}

	p.log.Debugf("Currency Rate %+v", cRate)
	prod.Price = prod.Price * cRate.Rate

	return p.responseData(c, prod)
}

// PUT /product/:id
func (p *productHandler) updateProduct(c echo.Context) error {
	id := p.getProductIDParam(c)

	prod := data.Product{}
	err := p.getProductData(c, &prod)
	if err != nil {
		return p.responseErrCode(c, http.StatusInternalServerError, "Failed to read request body")
	}

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		return p.responseErrCode(c, http.StatusNotFound, data.ErrProductNotFound.Error())
	}

	return p.responseStatusOK(c)
}

// DELETE product/:id
func (p *productHandler) deleteProduct(c echo.Context) error {
	id := p.getProductIDParam(c)

	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		return p.responseErrCode(c, http.StatusNotFound, data.ErrProductNotFound.Error())
	} else if err != nil {
		return p.responseErrCode(c, http.StatusInternalServerError, "")
	}

	return p.responseStatusOK(c)
}
