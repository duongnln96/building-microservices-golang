package routes

import (
	"github.com/duongnln96/building-microservices-golang/product-api/controller"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ProductRouterI interface {
	ProductRouterV1()
}

type ProductRouterDeps struct {
	Log    *zap.SugaredLogger
	Router *echo.Echo
	Ctrler controller.ProductControllerI
}

type productRouter struct {
	log    *zap.SugaredLogger
	router *echo.Echo
	ctrler controller.ProductControllerI
}

func NewProductRouter(deps ProductRouterDeps) ProductRouterI {
	return &productRouter{
		log:    deps.Log,
		router: deps.Router,
		ctrler: deps.Ctrler,
	}
}