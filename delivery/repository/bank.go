package repository

import (
	"errors"
	"gitlab.com/whoophy/privy/model"
	"gorm.io/gorm"
	"regexp"
	"time"
)

func NewBankRepo(db *gorm.DB) *Handler {
	return &Handler{db}
}

func (u *Handler) InitTransaction() *gorm.DB {
	return u.db.Begin()
}

func (u *Handler) CreateVA(bank model.BankBalance) error {
	bank.CreatedAt = time.Now()
	bank.UpdatedAt = time.Now()
	bank.CreatedBy = "SYSTEM"
	bank.UpdatedBy = "SYSTEM"
	err := u.db.Create(&bank).Error
	if err != nil {
		if isError, _ := regexp.MatchString("username", err.Error()); isError {
			return errors.New("username already exist")
		}
		if isError, _ := regexp.MatchString("email", err.Error()); isError {
			return errors.New("email already exist")
		}
		return err
	}

	return nil
}

func (u *Handler) GetBankBalance(filter model.FilterVA) (model.BankBalance, error) {
	var bank model.BankBalance
	if filter.ID != 0 {
		result := u.db.Find(&bank, "id = ?", filter.ID)
		if result.Error != nil {
			return bank, errors.New("error while get specific bank balance, id is wrong ")
		}
		return bank, nil
	}
	if filter.Code != "" {
		result := u.db.Find(&bank, "code = ?", filter.Code)
		if result.Error != nil {
			return bank, errors.New("error while get specific bank balance, code is wrong")
		}
		return bank, nil
	}

	return bank, nil
}

func (u *Handler) UpdateVATX(tx *gorm.DB, bank model.BankBalance) error {
	bank.UpdatedAt = time.Now()
	err := tx.Model(&bank).Where("code = ?", bank.Code).Updates(model.BankBalance{
		Balance:        bank.Balance,
		BalanceAchieve: bank.BalanceAchieve,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *Handler) CreateBalanceHistoryTx(tx *gorm.DB, history model.BankBalanceHistory) error {
	history.CreatedAt = time.Now()
	history.UpdatedAt = time.Now()
	history.CreatedBy = "SYSTEM"
	history.UpdatedBy = "SYSTEM"
	err := tx.Create(&history).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *Handler) GetBankBalanceTX(tx *gorm.DB, userid int64) ([]model.BankBalance, error) {
	var bank []model.BankBalance
	result := tx.Find(&bank, "user_id = ?", userid)
	if result.Error != nil {
		return bank, errors.New("error while get specific bank balance, id is wrong ")
	}
	return bank, nil
}

func (u *Handler) GetBankBalances(userid int64) ([]model.BankBalance, error) {
	var bank []model.BankBalance
	result := u.db.Find(&bank, "user_id = ?", userid)
	if result.Error != nil {
		return bank, errors.New("error while get specific bank balance, id is wrong ")
	}
	return bank, nil
}
