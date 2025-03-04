package middleware

import (
	"api/utils/errs"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (mw *Mw) ValidateNumericParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, param := range c.Params {
			if _, err := strconv.Atoi(param.Value); err != nil {
				errs.HandleError(c, errs.InvalidID) // Используем обработчик ошибок
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
