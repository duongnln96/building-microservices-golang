package repository

import (
	"database/sql"

	"github.com/duongnln96/building-microservices-golang/product-api/entity"
	tools "github.com/duongnln96/building-microservices-golang/product-api/tools/postgresql"
	"go.uber.org/zap"
)

type UserRepoI interface {
	FindUserByEmail(string) (entity.User, error)
	FintUserByID(int) (entity.User, error)
	CreateUser(*entity.User) error
	UpdateUser(*entity.User) error
	DeleteUser(*entity.User) error
	CheckEmailExisted(string) (bool, error)
}

type UserRepoDeps struct {
	Log *zap.SugaredLogger
	DB  tools.PsqlConnectorI
}

type userRepo struct {
	log *zap.SugaredLogger
	db  tools.PsqlConnectorI
}

func NewUserRepo(deps UserRepoDeps) UserRepoI {
	return &userRepo{
		log: deps.Log,
		db:  deps.DB,
	}
}

func (ur *userRepo) FindUserByEmail(email string) (entity.User, error) {
	user := entity.User{}
	err := ur.db.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
			if err != nil {
				return err
			}
			return nil
		},
		"select id, name, email, password from users where email=$1",
		email,
	)
	return user, err
}

func (ur *userRepo) FintUserByID(id int) (entity.User, error) {
	user := entity.User{}
	err := ur.db.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
			if err != nil {
				return err
			}
			return nil
		},
		"select id, name, email, password from users where id=$1",
		id,
	)
	return user, err
}

func (ur *userRepo) CreateUser(user *entity.User) error {
	err := ur.db.DBExec(
		"insert into users (name, email, password) values ($1, $2, $3)",
		user.Name, user.Email, user.Password,
	)
	return err
}

func (ur *userRepo) UpdateUser(user *entity.User) error {
	err := ur.db.DBExec(
		"update users set name=$1, email=$2, password=$3 where id=$4",
		user.Name, user.Email, user.Password, user.Id,
	)
	return err
}

func (ur *userRepo) DeleteUser(user *entity.User) error {
	err := ur.db.DBExec(
		"delete from users where id=$1",
		user.Id,
	)
	return err
}

func (ur *userRepo) CheckEmailExisted(email string) (bool, error) {
	var existedEmail string
	err := ur.db.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&existedEmail)
			if err != nil {
				return err
			}
			return nil
		},
		"select email from users where email=$1",
		email,
	)

	if existedEmail == email {
		return true, err
	}

	return false, err
}
