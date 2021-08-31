package server

import (
	"context"
	"io"
	"time"

	"github.com/duongnln96/building-microservices-golang/currency/internal/data"
	protos "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"go.uber.org/zap"
)

type CurrencyDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
	Db  data.CurrencyRatesI
}

type currency struct {
	protos.UnimplementedCurrencyServer
	log           *zap.SugaredLogger
	ctx           context.Context
	db            data.CurrencyRatesI
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrency(deps CurrencyDeps) protos.CurrencyServer {
	c := currency{
		log:           deps.Log,
		ctx:           deps.Ctx,
		db:            deps.Db,
		subscriptions: make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest),
	}

	go c.handleUpdate()

	return &c
}

func (c *currency) handleUpdate() {
	rateUpdateChan := c.db.MonitorRate(10 * time.Second)
	for range rateUpdateChan {
		c.log.Info("Get update rates")

		for k, v := range c.subscriptions {
			for _, rateReq := range v {
				rate, err := c.db.GetRate(rateReq.GetBase().String(), rateReq.Destination.String())
				if err != nil {
					c.log.Error("Unable to get update rate", " base ", rateReq.GetBase().String(), " destination ", rateReq.GetDestination().String())
				}

				err = k.Send(&protos.RateResponse{Base: rateReq.Base, Destination: rateReq.Destination, Rate: rate})
				if err != nil {
					c.log.Error("Unable to send updated rate, ", err)
				}
			}
		}
	}
}

func (c *currency) GetRate(ctx context.Context, r *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Infof("Handle request for GetRate base %s - dest %s", r.GetBase(), r.GetDestination())

	rate, err := c.db.GetRate(r.GetBase().String(), r.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Base: r.Base, Destination: r.Destination, Rate: rate}, nil
}

func (c *currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
	for {
		rr, err := src.Recv()
		if err == io.EOF {
			c.log.Infof("Client has closed the connection %+v", err)
			break
		}

		if err != nil {
			c.log.Error("Unable to read from client ", err)
			return err
		}

		c.log.Info("Hanle client request: ", "{ base: ", rr.GetBase(), ", destination: ", rr.GetDestination(), " }")

		rrs, ok := c.subscriptions[src]
		if !ok {
			c.log.Debug("Subscriber is not ok")
			rrs = []*protos.RateRequest{}
		}
		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}
	return nil
}
