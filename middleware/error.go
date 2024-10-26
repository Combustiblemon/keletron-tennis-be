package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Errors) > 0 {
			info, _ := ctx.Get("info")

			slog.Error("Errors detected in request %v", info, ctx.Errors)
		}

		ctx.Next()
	}
}
