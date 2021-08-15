package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/duongnln96/building-microservices-golang/currency/internal/data"
	"github.com/duongnln96/building-microservices-golang/currency/internal/server"
	pb "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var log *zap.SugaredLogger
var globalContex context.Context

var rootCmd = &cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		cr, err := data.NewCurrencyData(
			data.CurrencyRatesDeps{
				Log: log,
				Ctx: globalContex,
			},
		)
		if err != nil {
			log.Panicf("Cannot collect the data; %+v", err)
			os.Exit(1)
		}

		log.Info("Start gRPC application")
		gs := grpc.NewServer()
		c := server.NewCurrency(
			server.CurrencyDeps{
				Log: log,
				Ctx: globalContex,
				Db:  cr,
			},
		)

		pb.RegisterCurrencyServer(gs, c)
		// register the reflection service which allows clients to determine the methods
		// for this gRPC service
		reflection.Register(gs)

		// create a TCP socket for inbound server connections
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))
		if err != nil {
			log.Error("Unable to create listener", "error", err)
			os.Exit(1)
		}

		// listen for requests
		gs.Serve(l)
	},
}

func prepareLogging() {
	l, _ := zap.NewDevelopment()
	log = l.Sugar()
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
