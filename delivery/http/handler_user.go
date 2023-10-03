package http

import (
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/whoophy/privy/helper"
	"gitlab.com/whoophy/privy/model"
	"net/http"
	"time"
)

func (u Handler) LoginHandler(ctx *gin.Context) {
	user := model.UserLogin{}
	contentType := helper.GetContentType(ctx)

	if contentType == model.AppJson {
		ctx.ShouldBindJSON(&user)
	} else {
		ctx.ShouldBind(&user)
	}

	isValid, errx := govalidator.ValidateStruct(user)
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

	token, err := u.userUcase.UserLogin(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	cookie := &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
	}

	http.SetCookie(ctx.Writer, cookie)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success login, now you can access all menus",
	})
}

func (u Handler) RegisHandler(ctx *gin.Context) {
	user := model.User{}
	contentType := helper.GetContentType(ctx)

	if contentType == model.AppJson {
		ctx.ShouldBindJSON(&user)
	} else {
		ctx.ShouldBind(&user)
	}

	isValid, errx := govalidator.ValidateStruct(user)
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

	err := u.userUcase.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Message": "Success Registration, please login to access menu",
	})
}

func (u Handler) LogOutHandler(ctx *gin.Context) {
	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(ctx.Writer, &c)
	ctx.JSON(http.StatusOK, gin.H{
		"Message": "Success logout, please login to access menu",
	})
}

func (u Handler) ListVAHandler(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	email := userData["email"].(string)
	userBalance, bankBalance, err := u.userUcase.UserBalance(email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_balance":         userBalance,
		"user_virtual_account": bankBalance,
	})
}
