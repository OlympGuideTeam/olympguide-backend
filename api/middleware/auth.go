package middleware

import (
	"api/utils/constants"
	"api/utils/errs"
	"github.com/gin-gonic/gin"
)

func (mw *Mw) UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exist := c.Get(constants.ContextUserID)
		if !exist {
			errs.HandleError(c, errs.Unauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}

func (mw *Mw) AlreadyLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exist := c.Get(constants.ContextUserID)
		if exist {
			errs.HandleError(c, errs.UserAlreadyLoggedIn)
			c.Abort()
			return
		}
		c.Next()
	}
}
