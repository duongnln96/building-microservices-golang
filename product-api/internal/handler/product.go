package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/duongnln96/building-microservices-golang/product-api/config"
	"github.com/duongnln96/building-microservices-golang/product-api/internal/data"
	"github.com/duongnln96/building-microservices-golang/product-api/internal/logger"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type ProductHandlerI interface {
	StartServerLock()
	CreateProduct(echo.Context) error
	GetProducts(echo.Context) error
	UpdateProduct(echo.Context) error
}

type ProductHandlerDeps struct {
	Ctx context.Context
	Log *zap.SugaredLogger
	Cfg *config.AppConfig
}

type productHandler struct {
	ctx context.Context
	log *zap.SugaredLogger
	cfg *config.AppConfig
}

func NewProductHandler(deps ProductHandlerDeps) ProductHandlerI {
	return &productHandler{
		ctx: deps.Ctx,
		log: deps.Log,
		cfg: deps.Cfg,
	}
}

func (p *productHandler) StartServerLock() {
	e := echo.New()

	e.Use(logger.ZapLogger(p.log))

	e.POST("/products", p.CreateProduct)
	e.GET("/products", p.GetProducts)
	e.PUT("/products/:id", p.UpdateProduct)

	err := e.Start(fmt.Sprintf(":%d", p.cfg.Server.Port))
	if err != nil {
		p.log.Fatalf("Error while starting Echo server: %s", err)
	}
}

func (p *productHandler) CreateProduct(c echo.Context) error {
	product := data.Product{}

	defer c.Request().Body.Close()
	err := product.FromJSON(c.Request().Body)
	if err != nil {
		if err == io.EOF {
			return c.String(http.StatusBadRequest, "")
		} else {
			p.log.Infof("Failed to read request body for product %v", err)
			return c.String(http.StatusInternalServerError, "")
		}
	}

	data.AddProduct(&product)

	return c.String(http.StatusOK, "")
}

func (p *productHandler) GetProducts(c echo.Context) error {
	products := data.GetProducts()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return products.ToJSON(c.Response())
}

func (p *productHandler) UpdateProduct(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		p.log.Errorf("Cannot convert param string %+v", err)
	}

	prod := data.Product{}
	err = prod.FromJSON(c.Request().Body)
	if err != nil {
		if err == io.EOF {
			return c.String(http.StatusBadRequest, "")
		} else {
			p.log.Infof("Failed to read request body for product %v", err)
			return c.String(http.StatusInternalServerError, "")
		}
	}
	defer c.Request().Body.Close()

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		return c.String(http.StatusNotFound, data.ErrProductNotFound.Error())
	}

	return c.String(http.StatusOK, "")
}
