package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/duongnln96/building-microservices-golang/product-api/config"
	"github.com/duongnln96/building-microservices-golang/product-api/internal/handler"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger
var globalContex context.Context

var rootCmd = &cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := config.GetConfig()
		log.Infof("Run Product REST APIs with config %+v", appConfig)

		ph := handler.NewProductHandler(
			handler.ProductHandlerDeps{
				Ctx: globalContex,
				Log: log,
				Cfg: appConfig,
			},
		)

		ph.StartServerLock()
	},
}

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
		log.Infof("Got signal: %v+", sig)
		cancel()
		os.Exit(0)
	}(log)

	globalContex = ctx
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
}
