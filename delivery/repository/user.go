package repository

import (
	"errors"
	"fmt"
	"gitlab.com/whoophy/privy/model"
	"gorm.io/gorm"
	"regexp"
	"time"
)

type Handler struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *Handler {
	return &Handler{db}
}

func (u *Handler) GetUser(filter model.FilterUser) (model.User, error) {
	var user model.User
	if filter.ID != 0 {
		result := u.db.Find(&user, "id = ?", filter.ID)
		if result.Error != nil {
			return user, errors.New("error while get specific user with filter id ")
		}
		return user, nil
	}
	if filter.Username != "" {
		fmt.Println(filter.Username)
		result := u.db.Find(&user, "username = ?", filter.Username)
		if result.Error != nil {
			return user, errors.New("error while get specific user with filter username")
		}
		return user, nil
	}

	if filter.Email != "" {
		result := u.db.Find(&user, "email = ?", filter.Email)
		if result.Error != nil {
			return user, errors.New("error while get specific user with filter email")
		}
		return user, nil
	}
	return user, nil
}

func (u *Handler) RegisterUser(user model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.CreatedBy = "SYSTEM"
	user.UpdatedBy = "SYSTEM"
	err := u.db.Create(&user).Error
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

func (u *Handler) CreateUserBalance(userBalance model.UserBalance) error {
	userBalance.CreatedAt = time.Now()
	userBalance.UpdatedAt = time.Now()
	userBalance.CreatedBy = "SYSTEM"
	userBalance.UpdatedBy = "SYSTEM"
	err := u.db.Create(&userBalance).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *Handler) GetUserBalance(filter model.FilterUserBalance) (model.UserBalance, error) {
	var user model.UserBalance
	if filter.ID != 0 {
		result := u.db.Find(&user, "id = ?", filter.ID)
		if result.Error != nil {
			return user, errors.New("error while get specific user balance, id is wrong ")
		}
		return user, nil
	}
	if filter.UserID != 0 {
		result := u.db.Find(&user, "user_id = ?", filter.UserID)
		if result.Error != nil {
			return user, errors.New("error while get specific user balance, user id is wrong")
		}
		return user, nil
	}

	return user, nil
}

func (u *Handler) GetUserBalanceTx(tx *gorm.DB, filter model.FilterUserBalance) (model.UserBalance, error) {
	var user model.UserBalance
	if filter.ID != 0 {
		result := tx.Find(&user, "id = ?", filter.ID)
		if result.Error != nil {
			return user, errors.New("error while get specific user balance, id is wrong ")
		}
		return user, nil
	}
	if filter.UserID != 0 {
		result := tx.Find(&user, "user_id = ?", filter.UserID)
		if result.Error != nil {
			return user, errors.New("error while get specific user balance, user id is wrong")
		}
		return user, nil
	}

	return user, nil
}

func (u *Handler) UpdateUserBalanceTx(tx *gorm.DB, userBalance model.UserBalance) error {
	userBalance.UpdatedAt = time.Now()
	err := tx.Model(&userBalance).Where("user_id = ?", userBalance.UserID).Updates(model.UserBalance{
		Balance:        userBalance.Balance,
		BalanceAchieve: userBalance.BalanceAchieve,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *Handler) CreateUserBalanceHistoryTx(tx *gorm.DB, history model.UserBalanceHistory) error {
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
