package model

type (
	BankBalance struct {
		GormModel
		BankName           string               `gorm:"not null" json:"bank_name,omitempty" valid:"required"`
		Balance            int64                `gorm:"not null;default:0" json:"balance"`
		BalanceAchieve     int64                `gorm:"not null;default:0" json:"balance_achieve"`
		Code               string               `gorm:"not null;uniqueIndex"json:"code,omitempty" valid:"required~code should be filled,numeric,minstringlength(10),maxstringlength(12) "`
		Enable             bool                 `gorm:"not null;default:true"json:"enable"`
		UserID             int64                `gorm:"not null" json:"user_id"`
		User               *User                `json:"user,omitempty"`
		BankBalanceHistory []BankBalanceHistory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL,references:BankBalanceID" json:"user_balance_id,omitempty"`
	}

	BankBalanceHistory struct {
		GormModel
		BankBalanceID int64        `gorm:"not null"  json:"bank_balance_id"`
		BalanceBefore int64        `gorm:"not null;default:0" json:"balance_before"`
		BalanceAfter  int64        `gorm:"not null;default:0" json:"balance_after"`
		Activity      string       `gorm:"not null" json:"activity"`
		Type          string       `gorm:"not null" json:"type"`
		IP            string       `gorm:"not null" json:"ip"`
		Location      string       `gorm:"not null" json:"location"`
		UserAgent     string       `gorm:"not null" json:"user_agent"`
		Author        string       `gorm:"not null" json:"author"`
		BankBalance   *BankBalance `json:"bank,omitempty"`
	}

	VATopup struct {
		VaCode        string `json:"va_code" valid:"required~code should be filled,numeric,minstringlength(10),maxstringlength(12)"`
		DesiredAmount int64  `json:"desired_amount" valid:"required~Amount should be filled, numeric"`
		Email         string `json:"email"`
	}

	FilterVA struct {
		Code string `json:"code"`
		ID   int64  `json:"ID"`
	}

	Transfer struct {
		SourceCode string `json:"source_code" valid:"required~code should be filled,numeric,minstringlength(10),maxstringlength(12)"`
		TargetCode string `json:"target_code" valid:"required~code should be filled,numeric,minstringlength(10),maxstringlength(12)"`
		Amount     int64  `json:"amount" valid:"required~Amount should be filled, numeric"`
		UserEmail  string `json:"user_email"`
	}
)
