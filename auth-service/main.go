package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/duongnln96/building-microservices-golang/auth-service/config"
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
}
