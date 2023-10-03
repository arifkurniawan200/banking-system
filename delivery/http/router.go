package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/whoophy/privy/delivery"
	"gitlab.com/whoophy/privy/delivery/middleware"
)

type Handler struct {
	userUcase   delivery.UserUsecase
	bankUsecase delivery.BankUsecase
}

func Router(u delivery.UserUsecase, b delivery.BankUsecase) *gin.Engine {
	r := gin.Default()

	handler := Handler{
		userUcase:   u,
		bankUsecase: b,
	}
	fmt.Println(handler)

	r.POST("/login", handler.LoginHandler)
	r.POST("/register", handler.RegisHandler)
	userRouter := r.Group("/users")
	{
		userRouter.Use(middleware.Authentication())
		userRouter.Use(middleware.Authorization())
		userRouter.GET("/virtual-account/list", handler.ListVAHandler)
		userRouter.POST("/virtual-account/topup", handler.TopupHandler)
		userRouter.POST("/virtual-account/transfer", handler.TransferHandler)
		userRouter.GET("/logout", handler.LogOutHandler)
		userRouter.POST("/virtual-account/create", handler.CreateVAHandler)
	}
	return r
}
