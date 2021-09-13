package routes

import (
	"github.com/duongnln96/building-microservices-golang/auth-service/src/controller"
	"github.com/duongnln96/building-microservices-golang/auth-service/src/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthRouterI interface {
	Version1()
}

type AuthRouterDeps struct {
	Log        *zap.SugaredLogger
	Router     *echo.Echo
	Controller controller.AuthControllerI
	JWTSvc     service.JWTServiceI
}

type authRouter struct {
	log         *zap.SugaredLogger
	router      *echo.Echo
	controllers controller.AuthControllerI
	jwtService  service.JWTServiceI
}

func NewProductRouter(deps AuthRouterDeps) AuthRouterI {
	return &authRouter{
		log:         deps.Log,
		router:      deps.Router,
		controllers: deps.Controller,
		jwtService:  deps.JWTSvc,
	}
}

func (ar *authRouter) Version1() {
	v1 := ar.router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", ar.controllers.Login)
		auth.POST("/register", ar.controllers.Register)
	}
}
