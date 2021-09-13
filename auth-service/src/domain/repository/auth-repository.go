package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/duongnln96/building-microservices-golang/auth-service/config"
	"github.com/duongnln96/building-microservices-golang/auth-service/src/domain/entity"
	"github.com/duongnln96/building-microservices-golang/auth-service/tools/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type AuthRepoI interface {
	FindUserByEmail(string) (entity.User, error)
	// FintUserByID(int) (entity.User, error)
	CreateUser(*entity.User) error
	// UpdateUser(*entity.User) error
	// DeleteUser(*entity.User) error
	CheckEmailExisted(string) (bool, error)
}

type AuthRepoDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
	Cfg config.AuthServiceConfig
	Db  mongodb.MongoDBConnectorI
}

type authRepo struct {
	log *zap.SugaredLogger
	ctx context.Context
	cfg config.AuthServiceConfig
	db  mongodb.MongoDBConnectorI
}

func NewAuthRepo(deps AuthRepoDeps) AuthRepoI {
	return &authRepo{
		log: deps.Log,
		ctx: deps.Ctx,
		cfg: deps.Cfg,
		db:  deps.Db,
	}
}

func (ap *authRepo) FindUserByEmail(email string) (entity.User, error) {
	user := entity.User{}
	coll, err := ap.db.OpenCollection(ap.cfg.DBCollection)
	if err != nil {
		return user, err
	}
	ctx, cancel := context.WithTimeout(ap.ctx, 500*time.Millisecond)
	defer cancel()
	err = coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (ap *authRepo) CreateUser(user *entity.User) error {
	coll, err := ap.db.OpenCollection(ap.cfg.DBCollection)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ap.ctx, 500*time.Millisecond)
	defer cancel()
	user.ID = primitive.NewObjectID()
	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (ap *authRepo) CheckEmailExisted(email string) (bool, error) {
	coll, err := ap.db.OpenCollection(ap.cfg.DBCollection)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(ap.ctx, 500*time.Millisecond)
	defer cancel()
	count, err := coll.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, fmt.Errorf("Email is existed")
	}

	return true, nil
}
