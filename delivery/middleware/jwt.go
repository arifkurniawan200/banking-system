package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/whoophy/privy/model"
	"os"
	"time"
)

func GenerateToken(user model.User) string {
	secretkey := os.Getenv("SECRETKEY")
	claims := jwt.MapClaims{
		"id_user": user.ID,
		"email":   user.Email,
		"expired": time.Now().Add(time.Minute * 15),
	}
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := parseToken.SignedString([]byte(secretkey))
	return signedToken
}

func VerifyToken(c *gin.Context) (interface{}, error) {
	cookies, err := c.Request.Cookie("token")
	if err != nil {
		return nil, errors.New("error while get cookies, please login again")
	}

	if cookies.Value == "" {
		return nil, errors.New("please login to access menu")
	}

	tokenValue := cookies.Value

	//setup process
	secretkey := os.Getenv("SECRETKEY")
	errResponse := errors.New("Sign in to proced")

	token, _ := jwt.Parse(tokenValue, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}
		return []byte(secretkey), nil
	})
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errResponse
	}
	return token.Claims.(jwt.MapClaims), nil
}
