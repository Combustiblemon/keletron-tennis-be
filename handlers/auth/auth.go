package auth

import (
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Session() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println(helpers.GetURL(ctx))

		helpers.SetAuthCookie(ctx, "bananana")
		ctx.Status(http.StatusOK)
	}
}

func Login() gin.HandlerFunc {
	return func(_ctx *gin.Context) {

	}
}

func Register() gin.HandlerFunc {
	return func(_ctx *gin.Context) {

	}
}
