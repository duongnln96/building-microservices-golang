package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/duongnln96/building-microservices-golang/auth-service/config"
	"github.com/duongnln96/building-microservices-golang/auth-service/middleware"
	"github.com/duongnln96/building-microservices-golang/auth-service/tools/mongodb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger
var globalContext context.Context

func prepareLogger() {
	logger, _ := zap.NewDevelopment()
	log = logger.Sugar()
	log.Info("Log is prepared in development mode")
}

func init() {
	prepareLogger()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func(log *zap.SugaredLogger) {
		sig := <-signals
		log.Infof("Got signal: %#v", sig)
		cancel()
		os.Exit(1)
	}(log)

	globalContext = ctx
}

func main() {
	appConfig := config.GetConfig()
	log.Debugf("Start application with config %+v", appConfig)

	mongodbConn := mongodb.NewMongoDBConnector(mongodb.MongoDBConnectorDeps{
		Log: log,
		Ctx: globalContext,
		Cfg: appConfig.MongoDBConfig,
	})
	mongodbConn.Start()
	defer mongodbConn.Stop()

	e := NewEchoRouters()
	if err := e.Start(fmt.Sprintf(":%d", appConfig.AuthConfig.Port)); err != nil {
		log.Fatalf("Error while starting Echo server: %+v", err)
	}
}

func NewEchoRouters() *echo.Echo {
	e := echo.New()

	e.Use(middleware.ZapLogger(log))
	e.Validator = middleware.NewValidation()

	return e
}
