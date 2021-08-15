package server

import (
	"context"

	"github.com/duongnln96/building-microservices-golang/currency/internal/data"
	pb "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"go.uber.org/zap"
)

type CurrencyDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
	Db  data.CurrencyRates
}

type currency struct {
	pb.UnimplementedCurrencyServer
	log *zap.SugaredLogger
	ctx context.Context
	db  data.CurrencyRates
}

func NewCurrency(deps CurrencyDeps) pb.CurrencyServer {
	return &currency{
		log: deps.Log,
		ctx: deps.Ctx,
		db:  deps.Db,
	}
}

func (c *currency) GetRate(ctx context.Context, r *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Infof("Handle request for GetRate base %s - dest %s", r.GetBase(), r.GetDestination())

	rate, err := c.db.GetRate(r.GetBase().String(), r.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &pb.RateResponse{Rate: rate}, nil
}
