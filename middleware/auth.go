package middleware

import (
	"combustiblemon/keletron-tennis-be/database/models/UserModel"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func getUser(ctx *gin.Context) (*UserModel.User, error) {
	cookie, err := ctx.Request.Cookie("session")

	if err != nil {
		return nil, nil
	}

	return UserModel.FindOne(bson.D{{Key: "session", Value: cookie}})
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := getUser(ctx)

		if err != nil {
			helpers.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		if user == nil {
			helpers.ClearAuthCookie(ctx)
			ctx.JSON(http.StatusUnauthorized, map[string]any{})
			helpers.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.Set("User", user)
		ctx.Next()
	}
}

func Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := getUser(ctx)

		if err != nil {
			fmt.Printf("Error in Auth middleware: %v", err)
			ctx.JSON(http.StatusUnauthorized, map[string]any{})

			return
		}

		if user == nil || user.Role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, map[string]any{})

			return
		}

		ctx.Set("User", user)
		ctx.Set("Banana", user)
		ctx.Next()
	}
}
