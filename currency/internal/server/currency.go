package server

import (
	"context"

	pb "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"go.uber.org/zap"
)

type CurrencyDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
}

type currency struct {
	pb.UnimplementedCurrencyServer
	log *zap.SugaredLogger
	ctx context.Context
}

func NewCurrency(deps CurrencyDeps) pb.CurrencyServer {
	return &currency{
		log: deps.Log,
		ctx: deps.Ctx,
	}
}

func (c *currency) GetRate(ctx context.Context, r *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Infof("Handle request for GetRate base %s - dest %s", r.GetBase(), r.GetDestination())
	return &pb.RateResponse{Rate: 0.5}, nil
}
