package database

import (
	"gitlab.com/whoophy/privy/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type Handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}

func ConnectDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Printf("Error %s when preparing database", err)
		return nil, err
	}
	db.AutoMigrate(&model.User{}, &model.UserBalance{}, &model.UserBalanceHistory{}, &model.BankBalance{}, &model.BankBalanceHistory{})
	return db, nil
}
