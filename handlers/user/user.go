package user

import (
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		user, exists := helpers.GetUser(ctx)

		if exists {
			ctx.JSON(http.StatusOK, user.Sanitize())
			return
		}

		helpers.SendError(ctx, http.StatusNotFound, fmt.Errorf("userNotFound"))
	}
}
