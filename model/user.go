package model

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"gitlab.com/whoophy/privy/helper"
	"gorm.io/gorm"
)

type (
	User struct {
		GormModel
		Username    string        `gorm:"not null;uniqueIndex" json:"username,omitempty" valid:"required~username is required"`
		Password    string        `gorm:"not null" json:"password,omitempty" valid:"required~password is required,minstringlength(6)~Password has to have a minimum length of 6 characters"`
		Email       string        `gorm:"not null;uniqueIndex" json:"email,omitempty" valid:"required~email is required,email~invalid email format"`
		UserBalance []UserBalance `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL,references:UserID" json:"user_balance,omitempty"`
		Bank        []BankBalance `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL,references:UserID" json:"bank,omitempty"`
	}

	UserBalance struct {
		GormModel
		UserID             int64                `gorm:"not null" json:"user_id"`
		Balance            int64                `gorm:"not null;default:0" json:"user_balance"`
		BalanceAchieve     int64                `gorm:"not null;default:0" json:"balance_achieve"`
		User               *User                `json:"user,omitempty"`
		UserBalanceHistory []UserBalanceHistory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL,references:UserBalanceID" json:"user_balance_id,omitempty"`
	}

	UserBalanceHistory struct {
		GormModel
		UserBalanceID int64        `gorm:"not null"  json:"user_balance_id"`
		BalanceBefore int64        `gorm:"not null" json:"balance_before"`
		BalanceAfter  int64        `gorm:"not null" json:"balance_after"`
		Activity      string       `gorm:"not null" json:"activity"`
		Type          string       `gorm:"not null"  json:"type"`
		IP            string       `gorm:"not null" json:"ip"`
		Location      string       `gorm:"not null" json:"location"`
		UserAgent     string       `gorm:"not null" json:"user_agent"`
		Author        string       `gorm:"not null" json:"author"`
		UserBalance   *UserBalance `gorm:"not null" json:"user_balance,omitempty"`
	}

	FilterUser struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	FilterUserBalance struct {
		ID     int64 `json:"id"`
		UserID int64 `json:"user_id"`
	}

	UserLogin struct {
		Username string `json:"username,omitempty" valid:"required~username is required"`
		Password string `json:"password,omitempty" valid:"required~password is required,minstringlength(6)~Password has to have a minimum length of 6 characters"`
	}
)

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)
	if errCreate != nil {
		err = errCreate
		return
	}
	u.Password = helper.HashPass(u.Password)
	fmt.Println(u.Password)
	err = nil
	return
}
