package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	protos "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"github.com/duongnln96/building-microservices-golang/product-api/config"
	"github.com/duongnln96/building-microservices-golang/product-api/controller"
	"github.com/duongnln96/building-microservices-golang/product-api/middleware"
	"github.com/duongnln96/building-microservices-golang/product-api/repository"
	"github.com/duongnln96/building-microservices-golang/product-api/routes"
	"github.com/duongnln96/building-microservices-golang/product-api/service"
	tools "github.com/duongnln96/building-microservices-golang/product-api/tools/postgresql"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var log *zap.SugaredLogger
var globalContex context.Context

func prepareLogging() {
	logger, _ := zap.NewDevelopment()
	log = logger.Sugar()
	log.Info("Log is prepared in development mode")
}

func init() {
	prepareLogging()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func(log *zap.SugaredLogger) {
		sig := <-signals
		log.Infof("Got signal: %+v", sig)
		cancel()
		os.Exit(0)
	}(log)

	globalContex = ctx
}

func main() {
	appConfig := config.GetConfig()
	log.Infof("Starting Product Service with config %+v", appConfig)

	// Create the connection to curreny service
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", appConfig.CurrService.Host, appConfig.CurrService.Port), grpc.WithInsecure())
	if err != nil {
		log.Panic("Cannot create gRPC connection")
	}
	defer conn.Close()
	currency := protos.NewCurrencyClient(conn)

	// Product APIs service
	psql := tools.NewPsqlConnector(
		tools.PsqlConnectorDeps{
			Log: log,
			Ctx: globalContex,
			Cfg: appConfig.Psql,
		},
	)
	psql.Start()
	defer psql.Close()

	repo := repository.NewProductRepo(
		repository.ProductsRepoDeps{
			Log: log,
			DB:  psql,
		},
	)

	service := service.NewProductSerivce(
		service.ProductServiceDeps{
			Log:  log,
			Repo: repo,
		},
	)

	controller := controller.NewProductController(
		controller.ProductControllerDeps{
			Log: log,
			Ctx: globalContex,
			Svc: service,
			Cc:  currency,
		},
	)

	e := echoRouter()
	router := routes.NewProductRouter(
		routes.ProductRouterDeps{
			Log:    log,
			Router: e,
			Ctrler: controller,
		},
	)

	router.ProductRouterV1()
	if err := e.Start(fmt.Sprintf(":%d", appConfig.Server.Port)); err != nil {
		log.Fatalf("Error while starting Echo server: %+v", err)
	}
}

func echoRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.ZapLogger(log))
	e.Validator = middleware.NewValidation()

	return e
}
