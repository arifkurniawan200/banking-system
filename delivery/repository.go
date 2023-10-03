package delivery

import (
	"gitlab.com/whoophy/privy/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(user model.User) error
	CreateUserBalance(userBalance model.UserBalance) error
	GetUser(filter model.FilterUser) (model.User, error)
	GetUserBalance(filter model.FilterUserBalance) (model.UserBalance, error)

	// transaction
	GetUserBalanceTx(tx *gorm.DB, filter model.FilterUserBalance) (model.UserBalance, error)
	UpdateUserBalanceTx(tx *gorm.DB, userBalance model.UserBalance) error
	CreateUserBalanceHistoryTx(tx *gorm.DB, history model.UserBalanceHistory) error
}

type BankRepository interface {
	// InitTransaction transaction
	InitTransaction() *gorm.DB
	UpdateVATX(tx *gorm.DB, bank model.BankBalance) error
	CreateBalanceHistoryTx(tx *gorm.DB, history model.BankBalanceHistory) error
	GetBankBalanceTX(tx *gorm.DB, userid int64) ([]model.BankBalance, error)

	CreateVA(balance model.BankBalance) error
	GetBankBalance(filter model.FilterVA) (model.BankBalance, error)
	GetBankBalances(userid int64) ([]model.BankBalance, error)
}
