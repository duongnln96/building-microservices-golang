package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDBConnectorI interface {
	Start()
	Stop() error
	OpenCollection(string) (*mongo.Collection, error)
}

type MongoDBConnectorDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
	Cfg MongoDBConfig
}

type mongodbConnector struct {
	log      *zap.SugaredLogger
	ctx      context.Context
	cfg      MongoDBConfig
	conn     *mongo.Client
	isOk     bool
	stopChan chan bool
}

func NewMongoDBConnector(deps MongoDBConnectorDeps) MongoDBConnectorI {
	mongodbInfo := deps.Cfg.GetInfo()

	ctx, cancel := context.WithTimeout(deps.Ctx, deps.Cfg.Timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbInfo))
	if err != nil {
		deps.Log.Fatalf("Cannot connect to mongodb %+v\n", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		deps.Log.Fatalf("Cannot ping to mongodb %+v\n", err)
	}

	return &mongodbConnector{
		log:      deps.Log,
		ctx:      deps.Ctx,
		cfg:      deps.Cfg,
		conn:     client,
		isOk:     true,
		stopChan: make(chan bool),
	}
}

func (mc *mongodbConnector) connIsOK() bool {
	return mc.isOk
}

func (mc *mongodbConnector) healthCheck() {
	ticker := time.NewTicker(mc.cfg.HealthCheckInterval)
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(mc.ctx, mc.cfg.Timeout)
			if err := mc.conn.Ping(ctx, nil); err != nil {
				mc.isOk = false
			} else {
				mc.isOk = true
			}
			cancel()
		case <-mc.stopChan:
			close(mc.stopChan)
			return
		}
	}
}

func (mc *mongodbConnector) Start() {
	go mc.healthCheck()
}

func (mc *mongodbConnector) Stop() error {
	mc.stopChan <- true
	ctx, cancel := context.WithTimeout(mc.ctx, mc.cfg.Timeout)
	defer cancel()
	if err := mc.conn.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}

func (mc *mongodbConnector) OpenCollection(collectionName string) (*mongo.Collection, error) {
	if !mc.connIsOK() {
		return nil, fmt.Errorf("The connection to mongodb is not ok")
	}
	return mc.conn.Database(mc.cfg.DBName).Collection(collectionName), nil
}
