package service

import (
	"github.com/duongnln96/building-microservices-golang/product-api/config"
	"github.com/duongnln96/building-microservices-golang/product-api/repository"
	"go.uber.org/zap"
)

type ServicesDeps struct {
	Log  *zap.SugaredLogger
	Cfg  config.ServerConfig
	Repo repository.Repositories
}

type Services struct {
	JWTSvc     JWTSeriveI
	AuthenSvc  AuthenServiceI
	ProductSvc ProductServiceI
}

func NewServices(deps ServicesDeps) *Services {
	jwtService := NewJWTService(
		JWTSerivceDeps{
			Log: deps.Log,
			Cfg: deps.Cfg,
		},
	)

	authenService := NewAuthenService(
		AuthenServiceDeps{
			Log:      deps.Log,
			UserRepo: deps.Repo.UserRepo,
		},
	)

	productService := NewProductSerivce(
		ProductServiceDeps{
			Log:  deps.Log,
			Repo: deps.Repo.ProductRepo,
		},
	)

	return &Services{
		JWTSvc:     jwtService,
		AuthenSvc:  authenService,
		ProductSvc: productService,
	}
}
