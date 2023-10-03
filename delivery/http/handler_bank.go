package http

import (
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/whoophy/privy/helper"
	"gitlab.com/whoophy/privy/model"
	"net/http"
)

func (u Handler) CreateVAHandler(ctx *gin.Context) {
	bank := model.BankBalance{}
	contentType := helper.GetContentType(ctx)

	if contentType == model.AppJson {
		ctx.ShouldBindJSON(&bank)
	} else {
		ctx.ShouldBind(&bank)
	}

	isValid, errx := govalidator.ValidateStruct(bank)
	if errx != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": errx.Error(),
		})
		return
	}

	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "please input username and password",
		})
		return
	}

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	email := userData["email"].(string)
	err := u.bankUsecase.CreateVA(bank, email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "success to add virtual account",
	})
}

func (u Handler) TopupHandler(ctx *gin.Context) {
	data := model.VATopup{}
	contentType := helper.GetContentType(ctx)

	if contentType == model.AppJson {
		ctx.ShouldBindJSON(&data)
	} else {
		ctx.ShouldBind(&data)
	}

	isValid, errx := govalidator.ValidateStruct(data)
	if errx != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": errx.Error(),
		})
		return
	}

	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "please input username and password",
		})
		return
	}

	if data.DesiredAmount < 10000 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "minimum topup is 10000",
		})
		return
	}

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	email := userData["email"].(string)
	data.Email = email

	userBalance, err := u.bankUsecase.TopupVA(data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":   "success topup",
		"user_info": userBalance,
	})
}

func (u Handler) TransferHandler(ctx *gin.Context) {
	transfer := model.Transfer{}
	contentType := helper.GetContentType(ctx)

	if contentType == model.AppJson {
		ctx.ShouldBindJSON(&transfer)
	} else {
		ctx.ShouldBind(&transfer)
	}

	isValid, errx := govalidator.ValidateStruct(transfer)
	if errx != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": errx.Error(),
		})
		return
	}

	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "please input username and password",
		})
		return
	}

	if transfer.Amount < 10000 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "minimum transfer is 10000",
		})
		return
	}

	if transfer.TargetCode == transfer.SourceCode {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "code source and target cant be same",
		})
		return
	}

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	email := userData["email"].(string)
	transfer.UserEmail = email

	userBalance, err := u.bankUsecase.Transfer(transfer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success transfer",
		"data":    userBalance,
	})
}
