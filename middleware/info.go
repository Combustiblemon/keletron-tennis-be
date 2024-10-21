package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Info() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("info", time.Now().UnixMicro())

		ctx.Next()
	}
}
