package routes

import (
	"github.com/duongnln96/building-microservices-golang/product-api/src/controller"
	"github.com/duongnln96/building-microservices-golang/product-api/src/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ProductRouterI interface {
	ProductRouterV1()
}

type ProductRouterDeps struct {
	Log        *zap.SugaredLogger
	Router     *echo.Echo
	Controller controller.ProductControllerI
	JwtService service.JWTSeriveI
}

type productRouter struct {
	log        *zap.SugaredLogger
	router     *echo.Echo
	controller controller.ProductControllerI
	jwtService service.JWTSeriveI
}

func NewProductRouter(deps ProductRouterDeps) ProductRouterI {
	return &productRouter{
		log:        deps.Log,
		router:     deps.Router,
		controller: deps.Controller,
		jwtService: deps.JwtService,
	}
}
