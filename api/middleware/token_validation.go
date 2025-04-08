package middleware

import (
	"api/utils/constants"
	"api/utils/errs"
	"github.com/gin-gonic/gin"
	"strings"
)

func (mw *Mw) IDTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errs.HandleError(c, errs.Unauthorized)
			c.Abort()
			return
		}

		tempToken := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := mw.tokenService.ValidateIDToken(tempToken)
		if err != nil {
			errs.HandleError(c, err)
			c.Abort()
			return
		}
		c.Set(constants.ContextUserID, userID)
		c.Next()
	}
}

func (mw *Mw) EmailTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errs.HandleError(c, errs.Unauthorized)
			c.Abort()
			return
		}

		tempToken := strings.TrimPrefix(authHeader, "Bearer ")
		email, err := mw.tokenService.ValidateEmailToken(tempToken)
		if err != nil {
			errs.HandleError(c, err)
			c.Abort()
			return
		}

		c.Set(constants.ContextEmail, email)
		c.Next()
	}
}
