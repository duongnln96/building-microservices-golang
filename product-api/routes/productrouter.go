package routes

import (
	"github.com/duongnln96/building-microservices-golang/product-api/controller"
	"github.com/duongnln96/building-microservices-golang/product-api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ProductRouterI interface {
	ProductRouterV1()
}

type ProductRouterDeps struct {
	Log         *zap.SugaredLogger
	Router      *echo.Echo
	Controllers *controller.Controllers
	JwtService  service.JWTSeriveI
}

type productRouter struct {
	log         *zap.SugaredLogger
	router      *echo.Echo
	controllers *controller.Controllers
	jwtService  service.JWTSeriveI
}

func NewProductRouter(deps ProductRouterDeps) ProductRouterI {
	return &productRouter{
		log:         deps.Log,
		router:      deps.Router,
		controllers: deps.Controllers,
		jwtService:  deps.JwtService,
	}
}
