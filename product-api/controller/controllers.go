package controller

import (
	"context"

	protos "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"github.com/duongnln96/building-microservices-golang/product-api/service"
	"go.uber.org/zap"
)

type ControllersDeps struct {
	Log      *zap.SugaredLogger
	Ctx      context.Context
	Currency protos.CurrencyClient
	Svcs     service.Services
}

type Controllers struct {
	AuthenCtrl AuthenControllerI
	// UserCtrl    UserControllerI
	ProductCtrl ProductControllerI
}

func NewControllers(deps ControllersDeps) *Controllers {
	authenController := NewAuthorizeController(
		AuthenControllerDeps{
			Log:           deps.Log,
			AuthenService: deps.Svcs.AuthenSvc,
			JwtService:    deps.Svcs.JWTSvc,
		},
	)

	productController := NewProductController(
		ProductControllerDeps{
			Log: deps.Log,
			Ctx: deps.Ctx,
			Svc: deps.Svcs.ProductSvc,
			Cc:  deps.Currency,
		},
	)

	return &Controllers{
		AuthenCtrl:  authenController,
		ProductCtrl: productController,
	}
}
