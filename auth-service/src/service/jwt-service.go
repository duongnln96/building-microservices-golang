package service

import (
	"fmt"
	"time"

	"github.com/duongnln96/building-microservices-golang/auth-service/config"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type JWTServiceI interface {
	GenerateToken(string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type JWTServiceDeps struct {
	Log *zap.SugaredLogger
	Cfg config.AuthServiceConfig
}

type jwtService struct {
	log    *zap.SugaredLogger
	secret string
	issuer string
}

func NewJWTService(deps JWTServiceDeps) JWTServiceI {
	return &jwtService{
		log:    deps.Log,
		secret: deps.Cfg.JWTSercet,
		issuer: deps.Cfg.JWTSercet,
	}
}

type jwtCustomClaim struct {
	jwt.StandardClaims
	userID string
}

func (js *jwtService) GenerateToken(uID string) string {
	claim := jwtCustomClaim{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
			Issuer:    js.issuer,
			IssuedAt:  time.Now().Unix(),
		},
		uID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := token.SignedString([]byte(js.secret))
	if err != nil {
		js.log.Panic(err)
	}

	return t
}

func (js *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(js.secret), nil
		},
	)
}
