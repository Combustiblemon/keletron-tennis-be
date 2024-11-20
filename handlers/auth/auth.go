package auth

import (
	"combustiblemon/keletron-tennis-be/database/models/UserModel"
	"combustiblemon/keletron-tennis-be/modules/errorHandler"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"combustiblemon/keletron-tennis-be/modules/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Session() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, exists := helpers.GetUser(ctx)

		if exists {
			ctx.Status(http.StatusOK)
		} else {
			ctx.Status(http.StatusForbidden)
		}
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyAsByteArray, err := io.ReadAll(ctx.Request.Body)

		if err != nil {
			ctx.Status(http.StatusBadRequest)
			logger.Debug(ctx, err.Error())
			return
		}

		var data map[string]any
		err = json.Unmarshal(bodyAsByteArray, &data)

		if err != nil {
			errorHandler.SendError(ctx, http.StatusBadRequest, err)
			return
		}

		email, ok := data["email"].(string)

		if !ok {
			errorHandler.SendError(ctx, http.StatusBadRequest, errors.New("email required"))
			return
		}

		_, err = mail.ParseAddress(email)

		if err != nil {
			errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("email invalid"))
			return
		}

		password, ok := data["password"].(string)

		if !ok {
			errorHandler.SendError(ctx, http.StatusBadRequest, errors.New("password required"))
			return
		}
		user, err := UserModel.FindOne(bson.D{{Key: "email", Value: email}})

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		newSession, err := uuid.NewV7()

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		helpers.SetAuthCookie(ctx, helpers.Condition(user.Session == "", newSession.String(), user.Session))
		ctx.JSON(http.StatusOK, user.Sanitize())
	}
}

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyAsByteArray, err := io.ReadAll(ctx.Request.Body)

		if err != nil {
			ctx.Status(http.StatusBadRequest)
			logger.Debug(ctx, err.Error())
			return
		}

		var data map[string]any
		err = json.Unmarshal(bodyAsByteArray, &data)

		if err != nil {
			errorHandler.SendError(ctx, http.StatusBadRequest, err)
			return
		}

		email, ok := data["email"].(string)

		if !ok {
			errorHandler.SendError(ctx, http.StatusBadRequest, errors.New("email required"))
			return
		}

		_, err = mail.ParseAddress(email)

		if err != nil {
			errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("email invalid"))
			return
		}

		password, ok := data["password"].(string)

		if !ok {
			errorHandler.SendError(ctx, http.StatusBadRequest, errors.New("password required"))
			return
		}

		if len(password) < 6 {
			errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("password too short"))
		}

		name, ok := data["name"].(string)

		if !ok {
			errorHandler.SendError(ctx, http.StatusBadRequest, errors.New("name required"))
			return
		}

		usr, err := UserModel.FindOne(bson.D{{Key: "email", Value: strings.ToLower(strings.TrimSpace(email))}})

		if usr != nil {
			errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("email exists"))
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		session, err := uuid.NewV7()

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		user := UserModel.User{
			Name:        name,
			Role:        "USER",
			Email:       email,
			Password:    string(hash),
			FCMTokens:   []string{},
			ResetKey:    "",
			Session:     session.String(),
			AccountType: "PASSWORD",
		}

		err = UserModel.Create(user)

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
		} else {
			ctx.JSON(http.StatusCreated, user)
		}
	}
}
