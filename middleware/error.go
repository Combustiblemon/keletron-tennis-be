package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if len(ctx.Errors) > 0 {
			fmt.Println("Error number is", len(ctx.Errors))
		}

		ctx.Next()
	}
}
