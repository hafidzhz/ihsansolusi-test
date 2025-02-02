package dto

import "github.com/hafidzhz/ihsansolusi-test/app/entity"

type User struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	IdentityNumber string  `json:"idenity_number"`
	PhoneNumber    string  `json:"phone_number"`
	Balance        float64 `json:"balance"`
}

func ToUser(u *entity.User) *User {
	return &User{
		ID:             u.ID,
		Name:           u.Name,
		IdentityNumber: u.IdentityNumber,
		PhoneNumber:    u.PhoneNumber,
		Balance:        u.Balance,
	}
}

type CreateUserRequest struct {
	Name           string `json:"name" validate:"required,lte=50,gte=5"`
	IdentityNumber string `json:"identity_number" validate:"required,gte=16,numeric"`
	PhoneNumber    string `json:"phone_number" validate:"required,lte=15,gte=8,numeric"`
}

type DepositRequest struct {
	AccountNumber string  `json:"account_number" validate:"required,gte=17,numeric"`
	Amount        float64 `json:"amount" validate:"required,gte=1"`
}

type WithdrawRequest struct {
	DepositRequest
}

func ToUsers(users []*entity.User) []*User {
	res := make([]*User, len(users))
	for i, user := range users {
		res[i] = ToUser(user)
	}
	return res
}
