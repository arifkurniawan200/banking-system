package usecase

import (
	"errors"
	"gitlab.com/whoophy/privy/delivery"
	"gitlab.com/whoophy/privy/model"
)

type bankUsecase struct {
	userRepo delivery.UserRepository
	bankRepo delivery.BankRepository
}

func NewBankUsecase(userRepo delivery.UserRepository, bankRepo delivery.BankRepository) *bankUsecase {
	return &bankUsecase{
		userRepo: userRepo,
		bankRepo: bankRepo,
	}
}

func (u bankUsecase) CreateVA(bank model.BankBalance, email string) error {
	user, err := u.userRepo.GetUser(model.FilterUser{
		Email: email,
	})
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errors.New("cant find user/email in database")
	}
	bank.UserID = user.ID
	err = u.bankRepo.CreateVA(bank)
	if err != nil {
		return err
	}
	return nil
}

func (u bankUsecase) TopupVA(data model.VATopup) (model.UserBalance, error) {
	var user model.User
	var userBalance model.UserBalance
	user, err := u.userRepo.GetUser(model.FilterUser{
		Email: data.Email,
	})
	if err != nil {
		return userBalance, err
	}

	//check user is existed
	if user.ID == 0 {
		return userBalance, errors.New("cant find user/email in database")
	}

	bankBalance, err := u.bankRepo.GetBankBalance(model.FilterVA{Code: data.VaCode})
	if err != nil {
		return userBalance, err
	}

	if bankBalance.ID == 0 {
		return userBalance, errors.New("cant find virtual account")
	}

	if bankBalance.Enable != true {
		return userBalance, errors.New("virtual account is inactive")
	}

	//init transaction
	tx := u.bankRepo.InitTransaction()
	err = u.bankRepo.UpdateVATX(tx, model.BankBalance{
		Code:           bankBalance.Code,
		Balance:        bankBalance.Balance + data.DesiredAmount,
		BalanceAchieve: bankBalance.BalanceAchieve + data.DesiredAmount,
	})
	if err != nil {
		tx.Rollback()
		return userBalance, errors.New("failed to update bank balance")
	}

	//create bank balance history
	history := model.BankBalanceHistory{
		BankBalanceID: bankBalance.ID,
		BalanceBefore: bankBalance.Balance,
		BalanceAfter:  bankBalance.Balance + data.DesiredAmount,
		Activity:      model.ActivityTopup,
		Type:          model.Debit,
		IP:            "123.160.018",
		Location:      "Indonesia",
		UserAgent:     "API",
		Author:        user.Email,
	}

	err = u.bankRepo.CreateBalanceHistoryTx(tx, history)
	if err != nil {
		tx.Rollback()
		return userBalance, errors.New("failed to create bank balance history")
	}

	bankBalances, err := u.bankRepo.GetBankBalanceTX(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return userBalance, errors.New("failed to get bank balance")
	}

	var sumBalance int64

	for _, balance := range bankBalances {
		sumBalance = sumBalance + balance.Balance
	}

	userBalance, err = u.userRepo.GetUserBalanceTx(tx, model.FilterUserBalance{
		UserID: user.ID,
	})
	if err != nil {
		tx.Rollback()
		return userBalance, errors.New("failed to get user balance")
	}

	err = u.userRepo.UpdateUserBalanceTx(tx, model.UserBalance{
		UserID:         userBalance.UserID,
		Balance:        sumBalance,
		BalanceAchieve: sumBalance,
	})
	if err != nil {
		tx.Rollback()
		return userBalance, errors.New("failed to update user balance")
	}

	historyUser := model.UserBalanceHistory{
		UserBalanceID: userBalance.ID,
		BalanceBefore: userBalance.Balance,
		BalanceAfter:  sumBalance,
		Activity:      model.ActivityTopup,
		Type:          model.Debit,
		IP:            "123.160.018",
		Location:      "Boyolali",
		UserAgent:     "API",
		Author:        user.Email,
	}
	err = u.userRepo.CreateUserBalanceHistoryTx(tx, historyUser)
	if err != nil {
		tx.Rollback()
		return userBalance, errors.New("failed to create user balance history")
	}

	userBalance, err = u.userRepo.GetUserBalanceTx(tx, model.FilterUserBalance{
		UserID: user.ID,
	})
	if err != nil {
		tx.Rollback()
		return userBalance, errors.New("failed to get user balance")
	}

	if tx.Commit().Error != nil {
		return userBalance, errors.New("failed to commit database")
	}

	return userBalance, nil
}

func (u bankUsecase) Transfer(data model.Transfer) (userSource model.UserBalance, err error) {

	// check source VA
	sourceVA, err := u.bankRepo.GetBankBalance(model.FilterVA{
		Code: data.SourceCode,
	})
	if err != nil {
		return userSource, err
	}
	if sourceVA.ID == 0 {
		return userSource, errors.New("source va code is notfound")
	}

	if sourceVA.Balance < data.Amount {
		return userSource, errors.New("balance not enough, please topup or change another VA")
	}

	// check target VA
	targetVA, err := u.bankRepo.GetBankBalance(model.FilterVA{
		Code: data.TargetCode,
	})
	if err != nil {
		return userSource, err
	}
	if targetVA.ID == 0 {
		return userSource, errors.New("target va code is notfound")
	}

	//init transaction
	tx := u.bankRepo.InitTransaction()

	//update bank balance source
	err = u.bankRepo.UpdateVATX(tx, model.BankBalance{
		Code:           sourceVA.Code,
		Balance:        sourceVA.Balance - data.Amount,
		BalanceAchieve: sourceVA.Balance - data.Amount,
	})
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to update transfer, failed while update bank balance source")
	}

	//create history bank balance source
	history := model.BankBalanceHistory{
		BankBalanceID: sourceVA.ID,
		BalanceBefore: sourceVA.Balance,
		BalanceAfter:  sourceVA.Balance - data.Amount,
		Activity:      model.ActivityTransfer,
		Type:          model.Credit,
		IP:            "123.160.018",
		Location:      "Boyolali",
		UserAgent:     "API",
		Author:        data.UserEmail,
	}
	err = u.bankRepo.CreateBalanceHistoryTx(tx, history)
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to update transfer, failed while create bank balance history source")
	}

	//get current user bank balance
	bankBalancesSource, err := u.bankRepo.GetBankBalanceTX(tx, sourceVA.UserID)
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to update transfer, failed while update user balance source")
	}

	//temporary variable for totaling balance
	var sumBalance int64

	for _, balance := range bankBalancesSource {
		sumBalance = sumBalance + balance.Balance
	}

	// get current source balance
	sourceBalance, err := u.userRepo.GetUserBalanceTx(tx, model.FilterUserBalance{
		UserID: sourceVA.UserID,
	})
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to update transfer, failed while get user balance source")
	}

	//update source balance after transfer
	err = u.userRepo.UpdateUserBalanceTx(tx, model.UserBalance{
		UserID:         sourceBalance.UserID,
		Balance:        sumBalance,
		BalanceAchieve: sumBalance,
	})
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to update transfer, failed while update user balance source")
	}

	//create history for user balance
	historyUser := model.UserBalanceHistory{
		UserBalanceID: sourceBalance.ID,
		BalanceBefore: sourceBalance.Balance,
		BalanceAfter:  sumBalance,
		Activity:      model.ActivityTransfer,
		Type:          model.Credit,
		IP:            "123.160.018",
		Location:      "Boyolali",
		UserAgent:     "API",
		Author:        data.UserEmail,
	}
	err = u.userRepo.CreateUserBalanceHistoryTx(tx, historyUser)
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to update transfer, failed while create user balance source")
	}

	// update bank balance target
	err = u.bankRepo.UpdateVATX(tx, model.BankBalance{
		Code:           targetVA.Code,
		Balance:        targetVA.Balance + data.Amount,
		BalanceAchieve: targetVA.Balance + data.Amount,
	})
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to update transfer, failed while update bank balance target")
	}

	// create history for bank balance
	historyTarget := model.BankBalanceHistory{
		BankBalanceID: targetVA.ID,
		BalanceBefore: targetVA.Balance,
		BalanceAfter:  targetVA.Balance + data.Amount,
		Activity:      model.ActivityTransfered,
		Type:          model.Debit,
		IP:            "123.160.018",
		Location:      "Boyolali",
		UserAgent:     "API",
		Author:        data.UserEmail,
	}
	err = u.bankRepo.CreateBalanceHistoryTx(tx, historyTarget)
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to update transfer, failed while create bank balance history target")
	}

	//get current bank balance target
	bankBalanceTarget, err := u.bankRepo.GetBankBalanceTX(tx, targetVA.UserID)
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to update transfer, failed while get user balance target")
	}

	// temporary variable for getting total target balance
	sumBalance = 0

	for _, balance := range bankBalanceTarget {
		sumBalance = sumBalance + balance.Balance
	}

	// get current target balance
	targetBalance, err := u.userRepo.GetUserBalanceTx(tx, model.FilterUserBalance{
		UserID: targetVA.UserID,
	})
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to update transfer, failed while get user balance history")
	}

	//update bank balance target after transferred
	err = u.userRepo.UpdateUserBalanceTx(tx, model.UserBalance{
		UserID:         targetBalance.UserID,
		Balance:        sumBalance,
		BalanceAchieve: sumBalance,
	})
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to transfer, failed while update user balance")
	}

	// create user balance history
	historyBalanceTarget := model.UserBalanceHistory{
		UserBalanceID: targetBalance.ID,
		BalanceBefore: targetBalance.Balance,
		BalanceAfter:  sumBalance,
		Activity:      model.ActivityTransfered,
		Type:          model.Debit,
		IP:            "123.160.018",
		Location:      "Boyolali",
		UserAgent:     "API",
		Author:        data.UserEmail,
	}
	err = u.userRepo.CreateUserBalanceHistoryTx(tx, historyBalanceTarget)
	if err != nil {
		tx.Rollback()
		return userSource, errors.New("failed to transfer, failed while create user balance history")
	}

	// commit all transaction
	if tx.Commit().Error != nil {
		return userSource, errors.New("failed to commit database")
	}

	// get latest source balance
	sourceBalance, err = u.userRepo.GetUserBalance(model.FilterUserBalance{
		UserID: sourceVA.UserID,
	})
	if err != nil {
		return userSource, errors.New("failed to update transfer, failed while get user balance source")
	}
	return sourceBalance, nil
}
