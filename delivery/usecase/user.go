package usecase

import (
	"errors"
	"gitlab.com/whoophy/privy/delivery"
	"gitlab.com/whoophy/privy/delivery/middleware"
	"gitlab.com/whoophy/privy/helper"
	"gitlab.com/whoophy/privy/model"
)

type userUsecase struct {
	userRepo delivery.UserRepository
	bankRepo delivery.BankRepository
}

func NewUserUsecase(userRepo delivery.UserRepository, bankRepo delivery.BankRepository) *userUsecase {
	return &userUsecase{
		userRepo: userRepo,
		bankRepo: bankRepo,
	}
}

func (u userUsecase) UserLogin(data model.UserLogin) (string, error) {
	var err error
	var user model.User

	filter := model.FilterUser{
		Username: data.Username,
	}

	user, err = u.userRepo.GetUser(filter)
	if err != nil {
		return "", err
	}

	if user.ID == 0 {
		return "", errors.New("cant find username in database")
	}

	comparePass := helper.ComparePass([]byte(user.Password), []byte(data.Password))
	if !comparePass {
		return "", errors.New("incorrect password")
	}

	return middleware.GenerateToken(user), nil
}

func (u userUsecase) CreateUser(user model.User) error {
	err := u.userRepo.RegisterUser(user)
	if err != nil {
		return err
	}

	userInfo, err := u.userRepo.GetUser(model.FilterUser{Username: user.Username})
	if err != nil {
		return err
	}

	userBalance := model.UserBalance{
		UserID: userInfo.ID,
	}

	err = u.userRepo.CreateUserBalance(userBalance)
	if err != nil {
		return err
	}
	return nil
}

func (u userUsecase) UserBalance(email string) (userBalance model.UserBalance, bankBalance []model.BankBalance, err error) {
	user, err := u.userRepo.GetUser(model.FilterUser{Email: email})
	if err != nil {
		return userBalance, bankBalance, err
	}

	userBalance, err = u.userRepo.GetUserBalance(model.FilterUserBalance{UserID: user.ID})
	if err != nil {
		return userBalance, bankBalance, err
	}

	bankBalance, err = u.bankRepo.GetBankBalances(user.ID)
	if err != nil {
		return userBalance, bankBalance, err
	}
	return userBalance, bankBalance, err

}
