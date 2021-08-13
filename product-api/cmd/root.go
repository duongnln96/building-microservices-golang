package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"github.com/duongnln96/building-microservices-golang/product-api/config"
	"github.com/duongnln96/building-microservices-golang/product-api/internal/handler"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var log *zap.SugaredLogger
var globalContex context.Context

var rootCmd = &cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := config.GetConfig()
		log.Infof("Run Product REST APIs with config %+v", appConfig)

		// Create the connection to curreny service
		// TODO: get config
		conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
		if err != nil {
			log.Panicf("Cannot create gRPC connection %+v", err)
		}
		defer conn.Close()
		cc := pb.NewCurrencyClient(conn)

		ph := handler.NewProductHandler(
			handler.ProductHandlerDeps{
				Ctx:            globalContex,
				Log:            log,
				Cfg:            appConfig,
				CurrencyClient: cc,
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
