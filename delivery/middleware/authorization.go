package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		userData := c.MustGet("userData").(jwt.MapClaims)
		expired := userData["expired"].(string)

		date, errx := time.Parse(time.RFC3339, expired)
		if errx != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "failed when mapclaims",
				"message": "please retry again",
			})
			return
		}

		if date.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "token has ben expired",
				"message": "please login again to get new token",
			})
			return
		}
		c.Next()
	}
}
