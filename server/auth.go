package server

import (
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	HeaderAuthorization = "Authorization"

	AuthBearer = "Bearer "

	tokenKey = "token"
)

func authorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(HeaderAuthorization)
		before, after, ok := strings.Cut(authHeader, AuthBearer)
		token := strings.TrimSpace(after)
		if !ok || before != "" || token == "" {
			c.JSON(401, map[string]string{"response": "authorization required"})
			return
		}
		c.Set(tokenKey, token)
	}
}

func GetAccessToken(c *gin.Context) string {
	val, exists := c.Get(tokenKey)
	if !exists {
		return ""
	}

	return val.(string)
}
