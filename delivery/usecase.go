package delivery

import "gitlab.com/whoophy/privy/model"

type UserUsecase interface {
	UserLogin(data model.UserLogin) (string, error)
	CreateUser(user model.User) error
	UserBalance(email string) (userBalance model.UserBalance, bankBalance []model.BankBalance, err error)
}

type BankUsecase interface {
	CreateVA(data model.BankBalance, email string) error
	TopupVA(data model.VATopup) (model.UserBalance, error)
	Transfer(data model.Transfer) (userSource model.UserBalance, err error)
}
