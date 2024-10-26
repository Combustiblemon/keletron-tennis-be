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
	token, err := ctx.Request.Cookie("session")

	if err != nil {
		return nil, err
	}

	userData, err := helpers.ParseToken(token.Value)

	if err != nil {
		return nil, err
	}

	return UserModel.FindOne(bson.D{{Key: "session", Value: userData.Session}})
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := getUser(ctx)

		if err != nil {
			helpers.SendError(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		if user == nil {
			helpers.ClearAuthCookie(ctx)
			helpers.SendError(ctx, http.StatusUnauthorized, fmt.Errorf("forbidden"))
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

func Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := getUser(ctx)

		if err != nil {
			fmt.Printf("Error in Auth middleware: %v", err)
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}

		if user == nil || user.Role != "ADMIN" {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
