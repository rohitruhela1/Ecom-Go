package middleware

import (
	token "github.com/Ecom-go/tokens"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClintToken := c.Request.Header.Get("token")
		if ClintToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "No authorization token"})
			c.Abort()
		}
		claims, err := token.ValidateToken(ClintToken)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
