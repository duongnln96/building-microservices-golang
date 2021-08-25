package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/duongnln96/building-microservices-golang/product-api/config"
	"go.uber.org/zap"
)

type JWTSeriveI interface {
	GenerateToken(userID string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type JWTSerivceDeps struct {
	Log *zap.SugaredLogger
	Cfg config.ServerConfig
}

type jwtService struct {
	log    *zap.SugaredLogger
	secret string
	issuer string
}

type jwtCustomClaim struct {
	UserID string
	jwt.StandardClaims
}

var CustomClaim *jwtCustomClaim

func NewJWTService(deps JWTSerivceDeps) JWTSeriveI {
	return &jwtService{
		log:    deps.Log,
		secret: deps.Cfg.JwtSecret,
		issuer: deps.Cfg.JwtSecret,
	}
}

func (svc *jwtService) GenerateToken(userID string) string {
	CustomClaim = &jwtCustomClaim{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			Issuer:    svc.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaim)
	t, err := token.SignedString([]byte(svc.secret))
	if err != nil {
		svc.log.Panic(err)
	}

	return t
}

func (svc *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t.Header["alg"])
		}
		return []byte(svc.secret), nil
	})
}
