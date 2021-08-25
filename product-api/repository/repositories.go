package repository

import (
	tools "github.com/duongnln96/building-microservices-golang/product-api/tools/postgresql"
	"go.uber.org/zap"
)

type RepositoriesDeps struct {
	Log *zap.SugaredLogger
	DB  tools.PsqlConnectorI
}

type Repositories struct {
	UserRepo    UserRepoI
	ProductRepo ProductsRepoI
}

func NewRepositories(deps RepositoriesDeps) *Repositories {
	userRepo := NewUserRepo(
		UserRepoDeps{
			Log: deps.Log,
			DB:  deps.DB,
		},
	)

	productRepo := NewProductRepo(
		ProductsRepoDeps{
			Log: deps.Log,
			DB:  deps.DB,
		},
	)

	return &Repositories{
		UserRepo:    userRepo,
		ProductRepo: productRepo,
	}
}
