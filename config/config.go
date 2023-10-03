package config

import (
	"gitlab.com/whoophy/privy/database"
	"gitlab.com/whoophy/privy/delivery/http"
	"gitlab.com/whoophy/privy/delivery/repository"
	"gitlab.com/whoophy/privy/delivery/usecase"
	"log"
)

func NewConfig() error {
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Printf("[Connect database] error while connect database %v", err)
		return err
	}

	_ = database.New(db)
	userRepo := repository.NewUserRepo(db)
	bankRepo := repository.NewBankRepo(db)
	userUsecase := usecase.NewUserUsecase(userRepo, bankRepo)
	bankUsecase := usecase.NewBankUsecase(userRepo, bankRepo)

	r := http.Router(userUsecase, bankUsecase)
	r.Run(":4000")
	return nil
}
