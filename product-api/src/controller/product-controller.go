package controller

import (
	"context"
	"net/http"
	"strconv"

	protos "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"github.com/duongnln96/building-microservices-golang/product-api/dto"
	"github.com/duongnln96/building-microservices-golang/product-api/entity"
	"github.com/duongnln96/building-microservices-golang/product-api/helper"
	"github.com/duongnln96/building-microservices-golang/product-api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ProductControllerI interface {
	AllProducts(echo.Context) error
	FindProductByID(echo.Context) error
	CreateProduct(echo.Context) error
	UpdateProduct(echo.Context) error
	DeleteProduct(echo.Context) error
}

type ProductControllerDeps struct {
	Ctx context.Context
	Log *zap.SugaredLogger
	Svc service.ProductServiceI
	Cc  protos.CurrencyClient
}

type productController struct {
	ctx   context.Context
	log   *zap.SugaredLogger
	svc   service.ProductServiceI
	cc    protos.CurrencyClient
	ccSub protos.Currency_SubscribeRatesClient
	rates map[string]float64
}

func NewProductController(deps ProductControllerDeps) ProductControllerI {
	pc := productController{
		ctx:   deps.Ctx,
		log:   deps.Log,
		svc:   deps.Svc,
		cc:    deps.Cc,
		ccSub: nil,
		rates: make(map[string]float64),
	}
	go pc.handleRateUpdateStream()
	return &pc
}

func (pc *productController) handleRateUpdateStream() {
	sub, err := pc.cc.SubscribeRates(pc.ctx)
	if err != nil {
		pc.log.Error("Unable to subscribe for rates", "error", err)
	}

	pc.ccSub = sub

	for {
		rateRcv, err := pc.ccSub.Recv()
		pc.log.Info("Recieved updated rate from server", "dest", rateRcv.GetDestination().String())
		if err != nil {
			pc.log.Error("Error receiving message", "error", err)
			return
		}

		// Update to mem-cached
		pc.rates[rateRcv.Destination.String()] = rateRcv.Rate
	}
}

func (pc *productController) AllProducts(c echo.Context) error {
	currency := pc.getProductQuery(c)

	products, err := pc.svc.All()
	if err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, res)
	}

	if currency == "" {
		res := helper.BuildResponse("OK", products)
		return c.JSON(http.StatusOK, res)
	}

	rate, err := pc.getRate(currency)
	if err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, res)
	}

	pc.log.Infof("Currency Rate %+v", rate)

	prodsWithRate := make([]entity.Product, 0)
	for _, prod := range products {
		newProd := prod
		newProd.Price = newProd.Price * rate
		prodsWithRate = append(prodsWithRate, newProd)
	}

	res := helper.BuildResponse("OK", prodsWithRate)
	return c.JSON(http.StatusOK, res)
}

func (pc *productController) FindProductByID(c echo.Context) error {
	id := pc.getProductIDParam(c)
	currency := pc.getProductQuery(c)

	prod, err := pc.svc.FindByID(id)
	if err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, res)
	}

	if currency == "" {
		res := helper.BuildResponse("OK", prod)
		return c.JSON(http.StatusOK, res)
	}

	rate, err := pc.getRate(currency)
	if err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, res)
	}
	pc.log.Infof("Currency Rate %+v", rate)

	prodWithRate := prod
	prodWithRate.Price = prodWithRate.Price * rate

	res := helper.BuildResponse("OK", prodWithRate)
	return c.JSON(http.StatusOK, res)
}

func (pc *productController) CreateProduct(c echo.Context) error {
	prodCreateDto := dto.ProductDTO{}
	if err := c.Bind(&prodCreateDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, res)
	}

	if err := c.Validate(prodCreateDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, res)
	}

	err := pc.svc.Create(&prodCreateDto)
	if err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, res)
	}

	res := helper.BuildResponse("OK", prodCreateDto)
	return c.JSON(http.StatusOK, res)
}

func (pc *productController) UpdateProduct(c echo.Context) error {
	id := pc.getProductIDParam(c)

	prodUpdateDto := dto.ProductUpdateDTO{
		ID: id,
	}
	if err := c.Bind(&prodUpdateDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, res)
	}

	if err := c.Validate(prodUpdateDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, res)
	}

	if err := pc.svc.Update(&prodUpdateDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, res)
	}

	res := helper.BuildResponse("OK", prodUpdateDto)
	return c.JSON(http.StatusOK, res)
}

func (pc *productController) DeleteProduct(c echo.Context) error {
	id := pc.getProductIDParam(c)

	prodUpdateDto := dto.ProductUpdateDTO{
		ID: id,
	}

	if err := pc.svc.Delete(&prodUpdateDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, res)
	}

	res := helper.BuildResponse("OK", nil)
	return c.JSON(http.StatusOK, res)
}

// getProductIDParam returns the product ID from the URL
func (pc *productController) getProductIDParam(c echo.Context) int {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		pc.log.Panicf("Cannot convert param string %+v", err)
	}

	return id
}

// getProductQuery returns the product currency from the URL
func (pc *productController) getProductQuery(c echo.Context) string {
	currency := c.QueryParam("currency")
	return currency
}

// getRate return currency rate from currency client
func (pc *productController) getRate(currency string) (float64, error) {
	// if cached return
	if r, ok := pc.rates[currency]; ok {
		return r, nil
	}

	rateRequest := protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[currency]),
	}

	rate, err := pc.cc.GetRate(pc.ctx, &rateRequest)
	pc.rates[currency] = rate.Rate

	pc.ccSub.Send(&rateRequest)

	return float64(rate.Rate), err
}
