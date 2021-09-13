package service

import (
	"github.com/duongnln96/building-microservices-golang/auth-service/src/domain/dto"
	"github.com/duongnln96/building-microservices-golang/auth-service/src/domain/entity"
	"github.com/duongnln96/building-microservices-golang/auth-service/src/domain/repository"
	"github.com/duongnln96/building-microservices-golang/auth-service/tools/utils"
	"go.uber.org/zap"
)

type AuthServiceI interface {
	FindUserByEmail(email string) (entity.User, error)
	CreateUser(u *dto.UserRegisterDTO) error
	IsDuplicateEmail(email string) (bool, error)
	ValidateCredential(userCredential *dto.UserLoginDTO) (bool, entity.User)
}

type AuthServiceDeps struct {
	Log      *zap.SugaredLogger
	AuthRepo repository.AuthRepoI
}

type authService struct {
	log      *zap.SugaredLogger
	authRepo repository.AuthRepoI
}

func NewAuthService(deps AuthServiceDeps) AuthServiceI {
	return &authService{
		log:      deps.Log,
		authRepo: deps.AuthRepo,
	}
}

func (as *authService) FindUserByEmail(email string) (entity.User, error) {
	return as.authRepo.FindUserByEmail(email)
}

func (as *authService) CreateUser(u *dto.UserRegisterDTO) error {
	user := entity.User{}

	existed, err := as.IsDuplicateEmail(u.Email)
	if err != nil || !existed {
		return err
	}

	user.Name = u.Name
	user.Email = u.Email
	user.Password = utils.Encrypt(u.Password)

	if err := as.authRepo.CreateUser(&user); err != nil {
		return err
	}
	return nil
}

func (as *authService) IsDuplicateEmail(email string) (bool, error) {
	return as.authRepo.CheckEmailExisted(email)
}

func (as *authService) ValidateCredential(userCredential *dto.UserLoginDTO) (bool, entity.User) {
	user, err := as.authRepo.FindUserByEmail(userCredential.Email)
	if err != nil {
		return false, entity.User{}
	}

	if user.Password != utils.Encrypt(userCredential.Password) {
		return false, entity.User{}
	}

	return true, user
}
