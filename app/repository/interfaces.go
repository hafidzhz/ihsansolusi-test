package repository

import (
	"github.com/hafidzhz/ihsansolusi-test/app/dto"
	"github.com/hafidzhz/ihsansolusi-test/app/entity"
)

type UserRepository interface {
	CreateUser(request *dto.CreateUserRequest) (*entity.User, error)
	DepositToUserAccount(request *dto.DepositRequest) (*entity.User, error)
	WithdrawFromUserAccount(request *dto.DepositRequest) (*entity.User, error)
	ResolveByAccountNumber(accountNumber string) (user entity.User, err error)
}
